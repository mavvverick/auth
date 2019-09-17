package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	nano "github.com/matoous/go-nanoid"
	auth "gitlab.com/go-pher/go-auth/proto"
	"gitlab.com/go-pher/go-auth/providers/accountkit"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getOrCreateUser(ctx context.Context, db *gorm.DB, in *auth.FBAccountKitUserData, user accountkit.Me) (userFromDB User, err error) {
	var unm string

	db.Where(&User{PhoneNumber: user.Phone.NationalNumber}).Find(&userFromDB)

	if userFromDB.PhoneNumber == user.Phone.NationalNumber {
		return userFromDB, err
	}

	newID, _ := nano.Generate(alphabet, 15)
	// To check whether the username generated, exists in DB
	for true {
		unm = genUsername(user.Phone.NationalNumber)
		err := db.Where("username = ?", unm).Find(&userFromDB)

		if err != nil {
			break
		}
		continue
	}
	newUser := User{
		ID:          newID,
		PhoneNumber: user.Phone.NationalNumber,
		CountryCode: user.Phone.CountryPrefix,
		Username:    unm,
	}
	newProvider := Provider{
		ID:       user.ID,
		Provider: "fbAccountKit",
		UserID:   newID,
	}
	err = db.Create(&newUser).Error
	if err != nil {
		fmt.Println("ERROR ---------- ", err)
		return newUser, err
	}
	err = db.Create(&newProvider).Error
	if err != nil {
		fmt.Println("ERROR ---------- ", err)
		return newUser, err
	}
	return newUser, err
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
