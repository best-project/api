package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"time"
)

type UserClaim struct {
	ID    uint
	Email string
}

type CustomClaims struct {
	jwt.StandardClaims
	User *UserClaim `json:"user"`
}

var signingKey = []byte("jwt")

func NewCustomPayload(user *UserClaim, expire int64) *CustomClaims {
	now := time.Now()
	return &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
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

func ParseJWT(token string) (*UserClaim, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "while parsing jwt")
	}
	var user UserClaim

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		if err := mapstructure.Decode(claims["user"], &user); err != nil {
			return nil, errors.Wrap(err, "while decoding user")
		}
	} else {
		return nil, errors.Wrap(err, "while parsing a JWT token")
	}

	return &user, nil
}
