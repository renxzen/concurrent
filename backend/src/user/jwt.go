package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtGo "github.com/golang-jwt/jwt"
)

var JWT Provider = NewProvider()

type provider struct {
	SecretKey      []byte
	UserRepository UserRepository
}

type Provider interface {
	GetToken(User) (Token, error)
	CheckToken(string) (string, error)
}

func NewProvider() Provider {
	return &provider{
		SecretKey:      []byte("vTNUeCApemwHoA38gQ05btE7F8jh3HboPnWw6Lxy"),
		UserRepository: NewUserRepository(),
	}
}

func (p *provider) GetToken(user User) (Token, error) {
	claims := jwtGo.StandardClaims{
		Audience:  "concurrency",
		Subject:   user.Email,
		Issuer:    "GO API",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
	}

	jwtToken := jwtGo.NewWithClaims(jwtGo.SigningMethodHS512, claims)
	tokenStr, err := jwtToken.SignedString(p.SecretKey)
	if err != nil {
		return Token{}, err
	}

	token := Token{}
	token.Token = tokenStr

	return token, nil
}

func (p *provider) CheckToken(authorizationHeader string) (string, error) {
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
