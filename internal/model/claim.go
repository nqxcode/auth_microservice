package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/auth_v1.AuthV1/Get"
)

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
