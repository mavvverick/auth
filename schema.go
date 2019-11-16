package main

import (
	"time"
)

// Pld struct is used for payload data
type Pld struct {
	Sub      string `json:"sub,omitempty"`
	Username string `json:"username,omitempty"`
	Code     string `json:"code,omitempty"`
	Provider string `json:"provider,omitempty"`
	ACL      string `json:"acl,omitempty"`
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
