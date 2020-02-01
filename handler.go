package main

import (
	"context"
	"strconv"

	auth "github.com/YOVO-LABS/auth/proto"
	"github.com/go-redis/redis/v7"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

//FBAccountKitLogin verifies Account kit token and returns access and refresh tokens
func (s *Server) FBAccountKitLogin(ctx context.Context, in *auth.FBAccountKitUserData) (*auth.UserTokenResponse, error) {
	// var userCreated bool

	// // Get the details of user from FB
	// userFromFB, err := accountkit.GetDetailsFromFB(in.Token)
	// if err != nil {
	// 	return nil, err
	// }
	// // Current time for checking whether user was created
	// currTime := time.Now()

	// // Create/Get user
	// userFromDB, err := getOrCreateUser(ctx, s.db, in, userFromFB)
	// if err != nil {
	// 	return nil, err
	// }

	// // Check if new user is created.
	// if currTime.Before(userFromDB.CreatedAt) {
	// 	userCreated = true
	// }

	// // Get details of the user from cache
	// userFromCache, err := getUserFromCache(ctx, s.redis, userFromDB.ID)
	// if userFromCache[0] == nil || userFromCache[3] == nil {
	// 	// Update cache of the user.
	// 	updateUserInCache(ctx, s.redis, userFromDB)
	// }

	// // Create payload struct for token generation.
	// payload := Pld{
	// 	Sub:      userFromDB.ID,
	// 	Username: userFromDB.PhoneNumber,
	// 	Code:     in.Code,
	// 	Provider: "fbAccountKit",
	// 	ACL:      strconv.Itoa(userFromDB.ACL),
	// }

	// // Get Access and Refresh tokens for the user
	// tokens, err := token(s.redis, payload)
	// if err != nil {
	// 	return nil, err
	// }

	// // Response Struct
	// re := auth.UserTokenResponse{
	// 	AccessToken:       tokens[0],
	// 	RefreshToken:      tokens[1],
	// 	UserCreated:       strings.Join([]string{strconv.FormatBool(userCreated), userFromDB.ID}, ":"),
	// 	FirstTimeUser:     "false",
	// 	IsProfileComplete: "false",
	// }

	// // fmt.Println("ERROR: ", err)
	return &auth.UserTokenResponse{}, nil
}

// RefreshToken returns access token from refresh token
func (s *Server) RefreshToken(ctx context.Context, in *auth.RefreshTokenInput) (resp *auth.RefreshTokenResponse, err error) {
	// accessToken, err := refresh(s.redis, in)
	// if err != nil {
	// 	return resp, err
	// }
	// res := auth.RefreshTokenResponse{
	// 	Token: accessToken,
	// }
	return resp, err
}

// SendOTP sends OTP to the user's phone number
func (s *Server) SendOTP(ctx context.Context, in *auth.SendOTPInput) (resp *auth.SendOTPResponse, err error) {
	//		1. Get/Create user from/in the DB
	//		2. If OTP exists in cache, send that through the partner
	//		3. If OTP doesn't exists in cache, generate new OTP, save in cache and send through partner
	//		4. Send response
	var otp string

	// Create/Get user
	if !in.Resend {
		_, err := getOrCreateUser(ctx, s.db, in)
		if err != nil {
			return nil, err
		}
	}

	otp, err = getUserOTP(s.redis, in.Phone)
	if err == redis.Nil {
		otp, err = getRandNum()
		if err != nil {
			return resp, status.Error(codes.Internal, "Internal Error. Contact Support")
		}
		setUserOTP(s.redis, in.Phone, otp)
	} else if err != nil {
		return resp, status.Error(codes.Internal, "Internal Error. Contact Support")
	}
	// fmt.Println("OTP---", otp)
	//TODO: Send OTP via SMS provider
	resp = &auth.SendOTPResponse{
		Status: otp,
	}

	return resp, nil
}

// VerifyOTP verifies the OTP sent by the user
func (s *Server) VerifyOTP(ctx context.Context, in *auth.VerifyOTPInput) (resp *auth.VerifyOTPResponse, err error) {
	// 		1. OTP doesn't exists, send error, expired OTP
	//		2. Check in.OTP with saved OTP
	//			a. False, wrong OTP, try again
	//			b. True, next step
	//		3. Generate access token and return
	var user User
	user.ID = in.Phone
	var ftu bool
	if in.Otp != "true" {

		otp, err := getUserOTP(s.redis, in.Phone)
		if err == redis.Nil {
			return resp, status.Error(codes.InvalidArgument, "OTP expired. Resend")
		} else if err != nil {
			return resp, status.Error(codes.Internal, "Internal Error. Contact Support")
		}

		if in.Otp != otp {
			return resp, status.Error(codes.InvalidArgument, "Wrong OTP. Try again!")
		}

		user, ftu, err = updateAndReturnUser(ctx, s.db, in)
		if err != nil {
			return resp, err
		}
	} else {
		user, err = getUserByID(ctx, s.db, user.ID)
		if err != nil {
			return resp, err
		}
	}

	// Get details of the user from cache
	userFromCache, err := getUserFromCache(ctx, s.redis, user.ID)
	if err != nil || userFromCache[0] == nil || userFromCache[3] == nil || in.Otp == "true" {
		// Update cache of the user.
		updateUserInCache(ctx, s.redis, user)
	}

	// Create payload struct for token generation.
	payload := Pld{
		Sub:      user.ID,
		Username: user.Username,
		Code:     "yovo",
		Scope:    user.Scope,
		ACL:      strconv.Itoa(user.ACL),
	}

	// Get Access and Refresh tokens for the user
	tokens, err := token(s.redis, payload)
	if err != nil {
		return nil, err
	}

	resp = &auth.VerifyOTPResponse{
		Token: tokens[0],
		Ftu:   ftu,
	}

	return resp, nil
}

// Check handler is for health checking the gRPC service
func (s Server) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{
		Status: 1,
	}, nil
}

// Watch handler is for health checking the gRPC service
func (s Server) Watch(req *health.HealthCheckRequest, srv health.Health_WatchServer) error {
	// fmt.Println("WATCH")
	return nil
}
