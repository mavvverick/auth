package main

import (
	"time"

	"github.com/shopspring/decimal"
)

// User struct is model struct for user document
type User struct {
	ID              string          `gorm:"column:id;primary_key;size:20;unique_index:idx_uid_cur"`
	PhoneNumber     string          `gorm:"column:phoneNumber;not null"`
	Username        string          `gorm:"column:username;not null"`
	FirstTimeUser   bool            `gorm:"column:firstTimeUser;DEFAULT:true"`
	CountryCode     string          `gorm:"column:countryCode;not null"`
	IsActive        bool            `gorm:"column:isActive;DEFAULT:false"`
	IsBlocked       bool            `gorm:"column:isBlocked;DEFAULT:false"`
	UsernameUpdated bool            `gorm:"column:unmUpdt;DEFAULT:false"`
	ACL             int             `gorm:"column:acl;DEFAULT:0"`
	Scope           string          `gorm:"column:scope;DEFAULT:'na'"`
	Providers       []Provider      `gorm:"foreignkey:UserID;association_foreignkey:ID"`
	RefCode         string          `gorm:"column:refCode"`
	AnmCoinsAdded   bool            `gorm:"column:aca;DEFAULT:false"`
	IsRef           bool            `gorm:"column:isRef;DEFAULT:false"`
	Coins           int             `gorm:"column:coins;DEFAULT:0;"`
	Inr             decimal.Decimal `gorm:"column:inr;not null;DEFAULT:0;" sql:"type:decimal(10,2);"`
	CreatedAt       time.Time       `gorm:"column:createdAt"`
	UpdatedAt       time.Time       `gorm:"column:updatedAt"`
}

// TableName sets the table name in the DB
func (User) TableName() string {
	return "Users"
}

// Provider struct is model struct for provider document
type Provider struct {
	ID        string    `gorm:"column:id"`
	Provider  string    `gorm:"column:provider"`
	UserID    string    `gorm:"column:userId"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

// TableName sets the table name in the DB
func (Provider) TableName() string {
	return "Providers"
}
