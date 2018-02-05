package rpc

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// Claims defines a common structure used for JWT claims
type Claims struct {
	IssuerType string `json:"ist,omitempty"`
	jwt.StandardClaims
}

// NewUserToken generates a JWT for users and signs with given secret
func NewUserToken(username string, issuer string, audience string, secret string) (string, error) {
	claims := Claims{
		"user",
		jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   username,
			Audience:  audience,
			ExpiresAt: time.Now().AddDate(0, 0, 28).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// NewSignupToken generates a JWT for user signups
func NewSignupToken(email string, issuer string, audience string, secret string) (string, error) {
	claims := Claims{
		"signup",
		jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   email,
			Audience:  audience,
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	utoken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return utoken.SignedString([]byte(secret))
}

// ParseToken attempts to verify a signed JWT issued for user auth
func ParseToken(token string, secret string, audience string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(parsed *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid && claims.VerifyAudience(audience, true) {
		return claims, nil
	}
	return nil, err
}
