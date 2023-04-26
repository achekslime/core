package jwt

import "github.com/dgrijalva/jwt-go"

type JWTClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
