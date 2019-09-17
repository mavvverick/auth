package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	auth "gitlab.com/go-pher/go-auth/proto"
	"gitlab.com/go-pher/go-auth/providers/accountkit"
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
	if userFromCache[0] == nil || userFromCache[3] == nil {
		// Update cache of the user.
		updateUserInCache(ctx, s.redis, userFromDB)
	}

	// Create payload struct for token generation.
	payload := Pld{
		Sub:      userFromDB.ID,
		Username: userFromDB.Username,
		Code:     in.Code,
		Provider: "fbAccountKit",
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
