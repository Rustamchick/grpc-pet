package AuthService

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"grpc-pet/pkg/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserId int
	Email  string
	AppId  int
}

// GenerateToken return token signed by user and app
func GenerateToken(user models.User, app models.App, TokenTTL time.Duration) (string, error) {

	claims := &tokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now())},
		UserId: int(user.ID),
		Email:  user.Email,
		AppId:  app.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(app.Token))
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
		return nil, errors.New("Error parsing token")
	}

	claims, ok := ParsedToken.Claims.(*tokenClaims)
	if !ok {
		return nil, errors.New("TokenClaims are not of type *TokenClaims")
	}

	return claims, nil
}

func generatePasswordHash(password string) string {
	const salt = "r342fdsg8j8t4g20"

	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
