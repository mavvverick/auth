package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

func getUserFromCache(ctx context.Context, redCli *redis.Client, uid string) ([]interface{}, error) {
	return redCli.HMGet(strings.Join([]string{"u", uid}, ":"), "yovo", "username", "firstTimeUser", "isBlocked", "acl", "scope").Result()
}

func updateUserInCache(ctx context.Context, redCli *redis.Client, user User) (int64, error) {
	var m = make(map[string]interface{})
	m["username"] = user.Username
	m["phoneNumber"] = user.PhoneNumber
	m["firstTimeUser"] = user.FirstTimeUser
	m["isBlocked"] = strconv.FormatBool(user.IsBlocked)
	m["isActive"] = strconv.FormatBool(user.IsActive)
	m["unmUpdt"] = strconv.FormatBool(user.UsernameUpdated)
	m["acl"] = user.ACL
	m["scope"] = user.Scope
	m["coins"] = user.Coins
	m["inr"] = fmt.Sprintf("%v", user.Inr)
	m["picUpdt"] = strconv.FormatBool(user.ProfilePicUpdated)

	status, err := redCli.HMSet(strings.Join([]string{"u", user.ID}, ":"), m).Result()
	return status, err
}

func setUserAccessCode(redCli *redis.Client, pld Pld, accessKey string) {
	redCli.HSet(strings.Join([]string{"u", pld.Sub}, ":"), pld.Code, accessKey)
}

func getUserAccessCode(redCli *redis.Client, userID string, code string) (string, error) {
	return redCli.HGet(strings.Join([]string{"u", userID}, ":"), code).Result()
	//return accessCode, err
}

func getUserOTP(redCli *redis.Client, phone string) (string, error) {
	return redCli.Get(phone).Result()
}

func setUserOTP(redCli *redis.Client, phone, otp string) {
	redCli.Set(phone, otp, 5*time.Minute)
}

func isNumberBlocked(redCli *redis.Client, phone string) (bool, error) {
	return redCli.SIsMember("blockNum", phone).Result()
}
