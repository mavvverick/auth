package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-redis/redis"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/teris-io/shortid"
	auth "gitlab.com/go-pher/go-auth/proto"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func token(redis *redis.Client, payload Pld) ([]string, error) {
	// Generate accessKey for user.
	accessKey := getAccessKey()

	// Create JWT Access Token
	accessToken, err := getJWTToken(payload, accessKey)
	if err != nil {
		return []string{""}, err
	}

	// // Create JWE Refresh Token
	// refreshToken, err := getJWEToken(payload, accessKey)
	// if err != nil {
	// 	return []string{""}, err
	// }

	// Set the Access code for the user in redis
	setUserAccessCode(redis, payload, accessKey)
	return []string{accessToken}, err
}

func refresh(redis *redis.Client, in *auth.RefreshTokenInput) (string, error) {
	// Verify Refresh Token.
	refTok, err := jose.ParseEncrypted(in.Token)
	if err != nil {
		return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	privateKey, err := getPrivateKey()
	if err != nil {
		return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	// Decrypt the token.
	decrypted, err := refTok.Decrypt(privateKey)
	if err != nil {
		return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	// Get the decrypted data in a Map.
	bytesDecData := []byte(string(decrypted))
	var mapDecData map[string]interface{}
	json.Unmarshal(bytesDecData, &mapDecData)

	// Get the payload in a Map
	var mapPayload map[string]interface{}
	mapPayload = mapDecData["payload"].(map[string]interface{})

	// Get expiry of the token.
	exp, err := time.Parse(time.RFC3339, mapDecData["exp"].(string))

	// Check for sanity of the token data and payload
	if time.Now().After(exp) {
		return "", status.Error(codes.Unauthenticated, "Token Expired")
	} else if mapPayload["code"] == nil {
		return "", status.Error(codes.Unauthenticated, "Invalid token, client code missing")
	} else if mapPayload["code"] != in.Code {
		return "", status.Error(codes.Unauthenticated, "Request made from invalid client")
	} else if mapDecData["jti"] == nil {
		return "", status.Error(codes.Unauthenticated, "Invalid token, access key missing")
	}
	// Get access code of user from cache
	userAccessCode, err := getUserAccessCode(redis, mapPayload["sub"].(string), mapPayload["code"].(string))
	if userAccessCode != mapDecData["jti"].(string) {
		return "", status.Error(codes.Unauthenticated, "User already has active session on other device, logout and try again")
	}

	// Payload struct for token generation
	accessPld := Pld{
		mapPayload["sub"].(string),
		mapPayload["username"].(string),
		mapPayload["code"].(string),
		mapPayload["provider"].(string),
		mapPayload["acl"].(string),
	}
	accessToken, err := getJWTToken(accessPld, mapDecData["jti"].(string))
	return accessToken, err

}

func getJWTToken(clPay Pld, accessKey string) (string, error) {
	tokenOptions := getTokenOptions()
	// Get the private key for signing.
	// privateKey, err := getPrivateKey()
	// if err != nil {
	// 	return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	// }

	// Key ID for token header
	var kid jose.HeaderKey
	kid = "kid"
	key := tokenOptions["kid"].(string)

	// Create new signer with JWT type.
	sig, _ := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: getSharedKey()},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader(kid, key))

	// Set Token expiry.
	tokenExpDays, err := strconv.Atoi(strings.TrimSuffix(tokenOptions["tokenExp"].(string), "d"))
	tokenExpHours := time.Hour * 24 * time.Duration(tokenExpDays)
	//tokenExpHours := time.Minute * 5

	// Claims to add in the token.
	cl := jwt.Claims{
		Issuer:   tokenOptions["iss"].(string),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(tokenExpHours)),
		ID:       accessKey,
		//Audience: jwt.Audience{"playy"},
	}

	// Claims Payload struct containing all the data.
	clPld := ClPld{
		clPay,
		tokenOptions["aud"].(string),
	}

	// Get the signed JWT token.
	raw, err := jwt.Signed(sig).Claims(cl).Claims(clPld).CompactSerialize()

	// fmt.Println("RAW", raw)
	return raw, err
}

func getJWEToken(clPay Pld, accessKey string) (string, error) {
	// Create new Encrypter for encrypting data
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.RSA_OAEP, Key: getPublicKey()},
		(&jose.EncrypterOptions{}).WithType("JWE"))
	if err != nil {
		return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	// Payload struct for refresh token excryption
	tokenOptions := getTokenOptions()
	refTokenExpDays, err := strconv.Atoi(strings.TrimSuffix(tokenOptions["refreshExp"].(string), "d"))
	refTokenExpHours := time.Hour * 24 * time.Duration(refTokenExpDays)
	refPld := RefPld{
		accessKey,
		time.Now().Add(refTokenExpHours),
		clPay,
	}

	// Marshalling the struct to convert it to string
	pld, err := json.Marshal(refPld)
	var plaintext = []byte(string(pld))

	// Encryt the message
	object, err := encrypter.Encrypt(plaintext)
	if err != nil {
		return "", status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	// Get the JWE token in string form
	cs, err := object.CompactSerialize()
	// object.
	// fmt.Println("SER", cs)
	return cs, err
}

func getAccessKey() string {
	id, _ := shortid.Generate()
	return id[:4]
}

func getPublicKey() interface{} {
	pub := os.Getenv("CERT_PUB")
	set, err := jwk.ParseString(pub)
	if err != nil {
		panic(err)
	}

	key := set.Keys[0]
	// fmt.Println(key.KeyID())
	pubKey, err := key.Materialize()
	// fmt.Println("key", pubKey)
	return pubKey
}

func getSharedKey() interface{} {
	pub := os.Getenv("SHARED_KEY")
	set, err := jwk.ParseString(pub)
	if err != nil {
		panic(err)
	}

	key := set.Keys[0]
	// fmt.Println(key.KeyID())
	pubKey, err := key.Materialize()
	// fmt.Println("key", pubKey)
	return pubKey
}

func getPrivateKey() (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("CERT"))
	block, _ := pem.Decode(bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privateKey, err
}

func getTokenOptions() (tokenOptions map[string]interface{}) {
	// Get Token options and map them.
	tk := []byte(os.Getenv("TOKEN_OPTIONS"))
	json.Unmarshal(tk, &tokenOptions)
	return tokenOptions
}

// getRandNum returns a random number of size four
func getRandNum() (string, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(nBig.Int64()+1000, 10), nil
}
