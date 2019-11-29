package server

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
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

func ParseJWT(token string) (*internal.User, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "while parsing jwt")
	}
	var user internal.User

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		if err := mapstructure.Decode(claims["user"], &user); err != nil {
			return nil, errors.Wrap(err, "while decoding user")
		}
	} else {
		return nil, errors.Wrap(err, "while parsing a JWT token")
	}

	return &user, nil
}
