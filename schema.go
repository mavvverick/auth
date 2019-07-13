package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct is model struct for user document
type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Coupon        Coupon             `bson:"coupon"`
	PhoneNumber   string             `bson:"phoneNumber"`
	PaytmNumber   string             `bson:"paytmNumber"`
	Vpa           string             `bson:"vpa"`
	UnmUpdated    bool               `bson:"unmUpdated"`
	PicUpdated    bool               `bson:"picUpdated"`
	RefUpdated    bool               `bson:"refUpdated"`
	Dob           string             `bson:"dob"`
	Username      string             `bson:"username"`
	FirstTimeUser bool               `bson:"firstTimeUser"`
	CountryCode   string             `bson:"countryCode"`
	AvatarURL     string             `bson:"avatarUrl"`
	IsActive      bool               `bson:"isActive"`
	IsBlocked     bool               `bson:"isBlocked"`
	AndroidID     string             `bson:"androidId"`
	State         string             `bson:"state"`
	ImeiNumber    string             `bson:"imeiNumber"`
	ReferCode     string             `bson:"referCode"`
	BeneAdded     bool               `bson:"beneAdded"`
	Linked        bool               `bson:"linked"`
	Boost         int64              `bson:"boost"`
	Sub           string             `bson:"sub"`
	Ftue          int64              `bson:"ftue"`
	UserID        string             `bson:"userId"`
	Providers     []Provider         `bson:"providers"`
	Chest         Chest              `bson:"chest"`
	CreatedAt     time.Time          `bson:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt"`
}

// Chest struct is model struct for chest(inner user) document
type Chest struct {
	Hrs  Hrs  `bson:"HRS" json:"HRS"`
	Win  Win  `bson:"WIN" json:"WIN"`
	Mgpl Mgpl `bson:"MGPL" json:"MGPL"`
	Star Mgpl `bson:"STAR" json:"STAR"`
	Prg  Prg  `bson:"PRG" json:"PRG"`
}

// Hrs struct is model struct for HRS(inner user) document
type Hrs struct {
	Count int64 `bson:"count" json:"sub"`
	After int64 `bson:"after" json:"after"`
}

// Mgpl struct is model struct for MGPL(inner user) document
type Mgpl struct {
	Count  int64 `bson:"count" json:"count"`
	Points int64 `bson:"points" json:"points"`
	Rem    int64 `bson:"rem" json:"rem"`
}

// Prg struct is model struct for PRG(inner user) document
type Prg struct {
	Num       int64 `bson:"num" json:"num"`
	Count     int64 `bson:"count" json:"count"`
	LastAward int64 `bson:"lastAward" json:"lastAward"`
}

// Win struct is model struct for WIN(inner user) document
type Win struct {
	Count int64 `bson:"count" json:"count"`
	Num   int64 `bson:"num" json:"num"`
	After int64 `bson:"after" json:"after"`
}

// Coupon struct is model struct for coupon document
type Coupon struct {
	IsRedeem bool `bson:"isRedeem"`
}

// Provider struct is model struct for provider document
type Provider struct {
	ProviderID string `bson:"providerId"`
	Provider   string `bson:"provider"`
	ID         string `bson:"_id"`
}

// Pld struct is used for payload data
type Pld struct {
	Sub      string `json:"sub,omitempty"`
	Username string `json:"username,omitempty"`
	Code     string `json:"code,omitempty"`
	Provider string `json:"provider,omitempty"`
}

// ClPld struct is used for access token data
type ClPld struct {
	Payload Pld    `json:"payload,omitempty"`
	Aud     string `json:"aud,omitempty"`
}

// RefPld struct is used for refresh token data
type RefPld struct {
	Jti     string    `json:"jti,omitempty"`
	Expiry  time.Time `json:"exp,omitempty"`
	Payload Pld       `json:"payload,omitempty"`
}
