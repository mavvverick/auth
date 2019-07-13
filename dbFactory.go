package main

import (
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	nano "github.com/matoous/go-nanoid"
	auth "gitlab.com/go-pher/go-auth/proto"
	prvd "gitlab.com/go-pher/go-auth/providers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getOrCreateUser(ctx context.Context, coll *mongo.Collection, in *auth.FBAccountKitUserData, user prvd.Me) (userFromDB User, err error) {
	var unm string
	newID, _ := nano.Generate(alphabet, 15)
	// To check whether the username generated, exists in DB
	for true {
		unm = genUsername(user.Phone.NationalNumber)
		err := coll.FindOne(ctx, bson.D{
			{"username", unm},
		}).Decode(&userFromDB)

		if err != nil {
			break
		}
		continue
	}
	// Only used when a new user is created
	setOnInsert := bson.D{
		{"userId", newID},
		{"firstTimeUser", false},
		{"username", unm},
		{"androidId", in.AndroidId},
		{"imeiNumber", in.ImeiNumber},
		{"linked", true},
		{"state", "NA"},
		{"dob", "NA"},
		{"beneAdded", false},
		{"sub", "subToWrite"},
		{"referCode", "NA"},
		{"phoneNumber", user.Phone.NationalNumber},
		{"paytmNumber", user.Phone.NationalNumber},
		{"countryCode", user.Phone.CountryPrefix},
		{"vpa", "null"},
		{"unmUpdated", false},
		{"picUpdated", false},
		{"refUpdated", false},
		{"avatarUrl", "0"},
		{"isActive", true},
		{"isBlocked", false},
		{"providers", bson.A{
			bson.D{
				{"providerId", user.Application.ID},
				{"provider", "fbAccountKit"},
			},
		},
		},
		{"createdAt", time.Now()},
		{"updatedAt", time.Now()},
		{"chest", bson.D{
			{"HRS", bson.D{
				{"count", 0},
			}},
			{"PRG", bson.D{
				{"num", 0},
				{"count", 0},
				{"lastAward", 0},
			}},
		}}}
	err = coll.FindOneAndUpdate(ctx, bson.D{
		{"phoneNumber", user.Phone.NationalNumber},
	}, bson.D{{"$setOnInsert", setOnInsert}}, options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)).Decode(&userFromDB)
	return userFromDB, err
}

func genUsername(phn string) string {
	a := []rune(phn)
	unm := string(a[0:2]) + randomChars(5) + string(a[7:10])
	return unm
}

func randomChars(length int) string {
	var result string
	rand.Seed(time.Now().UnixNano())
	var characters = "*****#####@@@@@&&&&&"
	var charactersLength = len(characters)
	for i := 0; i < length; i++ {
		result += string(characters[rand.Intn(charactersLength)])
	}
	return result
}
