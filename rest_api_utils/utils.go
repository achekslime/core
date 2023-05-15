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

func BindNoContent(context *gin.Context) {
	context.JSON(http.StatusNoContent, nil)
	context.Abort()
	return
}

func BindBadRequest(context *gin.Context, err error) {
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	context.Abort()
	return
}

func BindUnauthorized(context *gin.Context, err error) {
	context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	context.Abort()
	return
}

func BindUnprocessableEntity(context *gin.Context, err error) {
	context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	context.Abort()
	return
}

func BindInternalError(context *gin.Context, err error) {
	context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	context.Abort()
	return
}
