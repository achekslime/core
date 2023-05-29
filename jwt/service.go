package jwt

import (
	"github.com/achekslime/core/rest_api_utils"
	"github.com/gin-gonic/gin"
)

func GetTokenClaims(context *gin.Context, jwtService *JwtService) (*JWTClaim, error) {
	token, err := rest_api_utils.ExtractToken(context)
	if err != nil {
		return nil, err
	}
	// validate token.
	userClaims, err := jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return userClaims, nil
}
