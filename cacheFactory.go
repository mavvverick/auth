package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

func getUserFromCache(ctx context.Context, redCli *redis.Client, uid string) ([]interface{}, error) {
	user, err := redCli.HMGet(strings.Join([]string{"u", uid}, ":"), "mgpl", "username", "firstTimeUser", "isBlocked", "acl").Result()
	if err == redis.Nil {
		fmt.Println("Not Found", err)
		panic(err)
	}
	return user, err
}

func updateUserInCache(ctx context.Context, redCli *redis.Client, user User) (status string, err error) {
	var m = make(map[string]interface{})
	m["username"] = user.Username
	m["phoneNumber"] = user.PhoneNumber
	m["firstTimeUser"] = user.FirstTimeUser
	m["isBlocked"] = strconv.FormatBool(user.IsBlocked)
	m["isActive"] = strconv.FormatBool(user.IsActive)
	m["acl"] = user.ACL

	status, err = redCli.HMSet(strings.Join([]string{"u", user.ID}, ":"), m).Result()
	return status, err
}

func setUserAccessCode(redCli *redis.Client, pld Pld, accessKey string) {
	redCli.HSet(strings.Join([]string{"u", pld.Sub}, ":"), pld.Code, accessKey)
}

func getUserAccessCode(redCli *redis.Client, userID string, code string) (string, error) {
	accessCode, err := redCli.HGet(strings.Join([]string{"u", userID}, ":"), code).Result()
	return accessCode, err
}
