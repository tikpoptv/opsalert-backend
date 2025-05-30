package jwt

import (
	"opsalert/config"
	"time"
)

type Settings struct {
	SecretKey string

	ExpirationTime time.Duration

	Issuer string

	Audience string
}

func DefaultSettings() *Settings {
	return &Settings{
		SecretKey:      config.Get().JWTSecret,
		ExpirationTime: time.Duration(config.Get().JWTExpirationHours) * time.Hour,
		Issuer:         "opsalert",
		Audience:       "opsalert-api",
	}
}

func NewSettings(secretKey string, expirationTime time.Duration, issuer, audience string) *Settings {
	return &Settings{
		SecretKey:      secretKey,
		ExpirationTime: expirationTime,
		Issuer:         issuer,
		Audience:       audience,
	}
}
