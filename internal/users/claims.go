package users

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

func NewUserClaims() {
	return
}
