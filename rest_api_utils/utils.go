package rest_api_utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
