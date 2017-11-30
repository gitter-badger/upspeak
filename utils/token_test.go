package utils

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	jwt "github.com/dgrijalva/jwt-go"
)

// TestTokenGenerateAndVerify generated a token with pre-defined values and expects to verify the token
func TestTokenGenerateAndVerify(t *testing.T) {
	issuer := "upspeak"
	audience := "upspeak.in"
	username := "foo"
	secret := "bar"

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

	token, err := NewUserToken(issuer, audience, username, secret)
	if err != nil {
		t.Error("Error generating user token")
	}

	expectedClaims, err := ParseToken(audience, token, secret)
	if err != nil {
		t.Error("Error parsing user token")
	}

	tokenEquality := cmp.Equal(claims, *expectedClaims)
	if tokenEquality != true {
		t.Error("Token doesn't match")
	}

}
