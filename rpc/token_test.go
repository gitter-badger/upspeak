package rpc

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	jwt "github.com/dgrijalva/jwt-go"
)

// TestTokenGenerateAndVerify generated a token with pre-defined values and expects to verify the token
func TestTokenGenerateAndVerify(t *testing.T) {
	issuer := "upspeak"
	audience := "upspeak.net"
	username := "foo"
	email := "foo@bar.com"
	secret := "bar"

	// User token generate and verify
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

	token, err := NewUserToken(username, issuer, audience, secret)
	if err != nil {
		t.Error("Error generating user token")
	}

	expectedClaims, err := ParseToken(token, secret, audience)
	if err != nil {
		t.Error("Error parsing user token")
	}

	tokenEquality := cmp.Equal(claims, *expectedClaims)
	if tokenEquality != true {
		t.Error("User token doesn't match")
	}

	// Signup token generate and verify
	claims = Claims{
		"signup",
		jwt.StandardClaims{
			Issuer:    issuer,
			Subject:   email,
			Audience:  audience,
			ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token, err = NewSignupToken(email, issuer, audience, secret)
	if err != nil {
		t.Error("Error generating signup token")
	}

	expectedClaims, err = ParseToken(token, secret, audience)
	if err != nil {
		t.Error("Error parsing signup token")
	}

	tokenEquality = cmp.Equal(claims, *expectedClaims)
	if tokenEquality != true {
		t.Error("Singup token doesn't match")
	}

}
