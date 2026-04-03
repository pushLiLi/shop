package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	uid, ok := userID.(uint)
	return uid, ok
}

func MustGetUserID(c *gin.Context) (uint, bool) {
	uid, ok := GetUserID(c)
	if !ok {
		ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return 0, false
	}
	return uid, true
}
