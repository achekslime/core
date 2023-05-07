package rest_api_utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ExtractToken(context *gin.Context) (string, error) {
	token := context.Query("token")
	if token != "" {
		return token, nil
	}
	bearerToken := context.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1], nil
	}
	return "", errors.New("invalid token")
}

func BindBadRequest(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	context.Abort()
	return
}

func BindInternalError(context *gin.Context, err error) {
	context.JSON(http.StatusInternalServerError, gin.H{"db error": err.Error()})
	context.Abort()
	return
}
