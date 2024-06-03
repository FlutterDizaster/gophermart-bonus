package models

import "github.com/golang-jwt/jwt/v4"

type Token struct {
	jwt.RegisteredClaims
	UserName string
}
