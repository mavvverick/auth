package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	auth "github.com/YOVO-LABS/auth/proto"
	"github.com/jinzhu/gorm"
	nano "github.com/matoous/go-nanoid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getOrCreateUser(ctx context.Context, db *gorm.DB, in *auth.SendOTPInput) (userFromDB User, err error) {
	var unm string

	db.Where(&User{PhoneNumber: in.Phone}).Find(&userFromDB)

	if userFromDB.PhoneNumber == in.Phone {
		return userFromDB, err
	}

	newID, _ := nano.Generate(alphabet, 15)
	// To check whether the username generated, exists in DB
	for true {
		unm = genUsername(in.Phone)
		err := db.Where("username = ?", unm).Find(&userFromDB)

		if err != nil {
			break
		}
		continue
	}
	newUser := User{
		ID:          newID,
		PhoneNumber: in.Phone,
		CountryCode: "91",
		Username:    unm,
	}
	// newProvider := Provider{
	// 	ID:       user.ID,
	// 	Provider: "fbAccountKit",
	// 	UserID:   newID,
	// }
	err = db.Create(&newUser).Error
	if err != nil {
		fmt.Println("ERROR ---------- ", err)
		return newUser, err
	}
	// err = db.Create(&newProvider).Error
	// if err != nil {
	// 	fmt.Println("ERROR ---------- ", err)
	// 	return newUser, err
	// }
	return newUser, err
}

func updateAndReturnUser(ctx context.Context, db *gorm.DB, in *auth.VerifyOTPInput) (userFromDB User, ftu bool, err error) {
	err = db.Where(&User{PhoneNumber: in.Phone}).Find(&userFromDB).Error
	if err != nil {
		return userFromDB, false, status.Error(codes.Internal, "Internal Error. Contact Support")
	}

	ftu = userFromDB.FirstTimeUser

	if !userFromDB.IsActive {
		userFromDB.IsActive = true
		userFromDB.FirstTimeUser = false
		err = db.Save(&userFromDB).Error
		if err != nil {
			return userFromDB, false, status.Error(codes.Internal, "Internal Error. Contact Support")
		}
	}

	return userFromDB, ftu, nil
}

func getUserByID(ctx context.Context, db *gorm.DB, userID string) (User, error) {
	var user User

	if db.Where(&User{ID: userID}).First(&user).Error != nil {
		return user, status.Error(codes.InvalidArgument, "UserId doesn't exists")
	}

	return user, nil
}

func genUsername(phn string) string {
	a := []rune(phn)
	unm := string(a[0:2]) + randomChars(5) + string(a[7:10])
	return unm
}

func randomChars(length int) string {
	var result string
	rand.Seed(time.Now().UnixNano())
	var characters = "*****_____-----"
	var charactersLength = len(characters)
	for i := 0; i < length; i++ {
		result += string(characters[rand.Intn(charactersLength)])
	}
	return result
}
