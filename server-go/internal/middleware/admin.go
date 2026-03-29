package middleware

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func AdminOnly(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		c.Abort()
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		c.Abort()
		return
	}

	if user.Role != "admin" && user.Role != "service" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}

func SuperAdminOnly(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		c.Abort()
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		c.Abort()
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Next()
}
