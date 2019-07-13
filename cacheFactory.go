package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

func getUserFromCache(ctx context.Context, redCli *redis.Client, uid string) ([]interface{}, error) {
	user, err := redCli.HMGet(strings.Join([]string{"u", uid}, ":"), "mgpl", "username", "firstTimeUser", "isBlocked").Result()
	if err == redis.Nil {
		fmt.Println("Not Found", err)
		panic(err)
	}
	return user, err
}

func updateUserInCache(ctx context.Context, redCli *redis.Client, user User) (status string, err error) {
	var m = make(map[string]interface{})
	chest, err := json.Marshal(user.Chest)
	m["username"] = user.Username
	m["phoneNumber"] = user.PhoneNumber
	m["firstTimeUser"] = user.FirstTimeUser
	m["paytmNumber"] = user.PaytmNumber
	m["isBlocked"] = strconv.FormatBool(user.IsBlocked)
	m["androidId"] = user.AndroidID
	m["imeiNumber"] = user.ImeiNumber
	m["referCode"] = user.ReferCode
	m["avatarUrl"] = user.AvatarURL
	m["unmUpdated"] = strconv.FormatBool(user.UnmUpdated)
	m["picUpdated"] = strconv.FormatBool(user.PicUpdated)
	m["refUpdated"] = strconv.FormatBool(user.RefUpdated)
	m["isActive"] = strconv.FormatBool(user.IsActive)
	m["sub"] = user.Sub
	m["chest"] = string(chest)
	m["ftue"] = user.Ftue

	status, err = redCli.HMSet(strings.Join([]string{"u", user.UserID}, ":"), m).Result()
	fmt.Println("REDIS ERROR:  ", err)
	return status, err
}

func setUserAccessCode(redCli *redis.Client, pld Pld, accessKey string) {
	redCli.HSet(strings.Join([]string{"u", pld.Sub}, ":"), pld.Code, accessKey)
}

func getUserAccessCode(redCli *redis.Client, userID string, code string) (string, error) {
	accessCode, err := redCli.HGet(strings.Join([]string{"u", userID}, ":"), code).Result()
	return accessCode, err
}
