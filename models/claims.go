package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId uint "json:userId"
	jwt.StandardClaims
}
