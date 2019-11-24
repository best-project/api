package server

import (
	"github.com/best-project/api/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type CustomClaims struct {
	jwt.StandardClaims
	User *internal.User `json:"user,omitempty"`
}

var signingKey = []byte("jwt")

func NewCustomPayload(user *internal.User) *CustomClaims {
	now := time.Now()
	return &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(24 * time.Hour).Unix(),
			IssuedAt:  now.Unix(),
		},
		User: user,
	}
}

func NewJWT(claims *CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tkn, err := token.SignedString(signingKey)
	if err != nil {
		return "", errors.Wrap(err, "while creating a token")
	}
	return tkn, nil
}

func ParseJWT(token string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "while parsing a JWT token")
	}

	return parsed, nil
}

func IsValid(token string) bool {
	tkn, err := ParseJWT(token)
	if err != nil {
		return false
	}
	return tkn.Valid
}
