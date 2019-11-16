package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	auth "gitlab.com/go-pher/go-auth/proto"
	"gitlab.com/go-pher/go-auth/providers/accountkit"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

// FBAccountKitLogin verifies Account kit token and returns access and refresh tokens
func (s *Server) FBAccountKitLogin(ctx context.Context, in *auth.FBAccountKitUserData) (*auth.UserTokenResponse, error) {
	var userCreated bool

	// Get the details of user from FB
	userFromFB, err := accountkit.GetDetailsFromFB(in.Token)
	if err != nil {
		return nil, err
	}
	// Current time for checking whether user was created
	currTime := time.Now()

	// Create/Get user
	userFromDB, err := getOrCreateUser(ctx, s.db, in, userFromFB)
	if err != nil {
		return nil, err
	}

	// Check if new user is created.
	if currTime.Before(userFromDB.CreatedAt) {
		userCreated = true
	}

	// Get details of the user from cache
	userFromCache, err := getUserFromCache(ctx, s.redis, userFromDB.ID)
	fmt.Println("CACHE--", userFromCache)
	if userFromCache[0] == nil || userFromCache[3] == nil {
		// Update cache of the user.
		updateUserInCache(ctx, s.redis, userFromDB)
	}

	fmt.Println("ACL---", userFromDB.ACL)
	// Create payload struct for token generation.
	payload := Pld{
		Sub:      userFromDB.ID,
		Username: userFromDB.PhoneNumber,
		Code:     in.Code,
		Provider: "fbAccountKit",
		ACL:      strconv.Itoa(userFromDB.ACL),
	}

	// Get Access and Refresh tokens for the user
	tokens, err := token(s.redis, payload)
	if err != nil {
		return nil, err
	}

	// Response Struct
	re := auth.UserTokenResponse{
		AccessToken:       tokens[0],
		RefreshToken:      tokens[1],
		UserCreated:       strings.Join([]string{strconv.FormatBool(userCreated), userFromDB.ID}, ":"),
		FirstTimeUser:     "false",
		IsProfileComplete: "false",
	}

	// fmt.Println("ERROR: ", err)
	return &re, err
}

// RefreshToken returns access token from refresh token
func (s *Server) RefreshToken(ctx context.Context, in *auth.RefreshTokenInput) (resp *auth.RefreshTokenResponse, err error) {
	accessToken, err := refresh(s.redis, in)
	if err != nil {
		return resp, err
	}
	res := auth.RefreshTokenResponse{
		Token: accessToken,
	}
	return &res, err
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
