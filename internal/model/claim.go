package model

import "github.com/dgrijalva/jwt-go"

// UserClaims user claims
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
