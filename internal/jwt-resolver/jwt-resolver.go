package jwtresolver

import (
	"errors"

	"github.com/FlutterDizaster/gophermart-bonus/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

type JWTResolver struct {
	secret string
}

func New(secret string) *JWTResolver {
	return &JWTResolver{
		secret: secret,
	}
}

func (res *JWTResolver) DecryptToken(tokenString string) (*models.Claims, error) {
	// Создание структуры models.Token
	claims := &models.Claims{}

	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error unexpected signing method")
		}
		return []byte(res.secret), nil
	})

	// Проверка токена на валидность
	if !token.Valid {
		return claims, errors.New("error invalid token")
	}

	return claims, err
}

func (res *JWTResolver) CreateToken(claims models.Claims) (string, error) {
	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Создание строки токена
	tokenString, err := token.SignedString([]byte(res.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
