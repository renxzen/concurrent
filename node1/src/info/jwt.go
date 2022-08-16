package info

import (
	"errors"
	"fmt"
	"strings"

	jwtGo "github.com/golang-jwt/jwt"
)

var JWT Provider = NewProvider()

type provider struct {
	SecretKey []byte
}

type Provider interface {
	CheckToken(string) (string, error)
}

func NewProvider() Provider {
	return &provider{
		SecretKey: []byte("vTNUeCApemwHoA38gQ05btE7F8jh3HboPnWw6Lxy"),
	}
}

func (p provider) CheckToken(authorizationHeader string) (string, error) {
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", errors.New("invalid token")
	}

	token := strings.TrimPrefix(authorizationHeader, "Bearer ")

	claims := &jwtGo.StandardClaims{}
	decodedToken, err := jwtGo.ParseWithClaims(token, claims, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, %v", token.Header["alg"])
		}

		return p.SecretKey, nil
	})

	if err != nil {
		return "", errors.New("error parsing token")
	}

	if !decodedToken.Valid {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil
}
