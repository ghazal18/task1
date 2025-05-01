package controller

import (
	"strings"
	"task1/entity"
	"time"
	"fmt"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Config struct {
	SignKey              string
	AccessExpirationTime time.Duration
	AccessSubject        string
}

type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {

	return s.createToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)

}

func (s Service) createToken(userID int, subject string, expireDuration time.Duration) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {

	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("userId: %v, expire at :%v", claims.UserID, claims.RegisteredClaims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
	}
}