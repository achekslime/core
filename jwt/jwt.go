package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const jwtLifeTime = 24 * time.Hour

type JwtService struct {
	jwtKey []byte
}

func ConfigureJWT(jwtKey string) *JwtService {
	return &JwtService{jwtKey: []byte(jwtKey)}
}

func (service *JwtService) GenerateJWT(email string) (string, error) {
	expiredTime := time.Now().Add(jwtLifeTime)

	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(service.jwtKey)
}

func (service *JwtService) ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return service.jwtKey, nil
		},
	)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("token expired")
	}
	return nil
}
