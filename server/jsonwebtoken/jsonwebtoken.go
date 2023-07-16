package jsonwebtoken

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hn275/envhub/server/lib"
)

var (
	secret string
)

func init() {
	secret = lib.Getenv("JWT_SECRET")
}

type GithubUser struct {
	Token     string `json:"token,omitempty"`
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type JwtToken struct {
	GithubUser `json:",inline"`
	jwt.RegisteredClaims
}

func Sign(data jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString([]byte(secret))
}

func Decode(t string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(t, &JwtToken{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("invalid signing algo")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtToken)
	if !ok || !token.Valid {
		return nil, errors.New("invalid jwt type")
	}

	return claims, nil
}
