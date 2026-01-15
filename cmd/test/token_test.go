package main

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int
	Email  string
	AppId  int
}

func parseToken(token, signingKey string) (claims *tokenClaims, err error) {

	ParsedToken, err := jwt.ParseWithClaims(token, &tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Invalid signing method")
			}
			return []byte(signingKey), nil
		})
	if err != nil {
		fmt.Printf("token: %s", ParsedToken.Raw)
		return nil, errors.New("Error parsing token")
	}

	claims, ok := ParsedToken.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("TokenClaims are not of type *TokenClaims")
	}

	return claims, nil
}
