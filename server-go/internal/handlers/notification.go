package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	var total int64
	database.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&total)

	var notifications []models.Notification
	offset := (page - 1) * limit
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&notifications)

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"limit":         limit,
		"totalPages":    totalPages,
	})
}

func GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}

	var count int64
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	result := database.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "通知不存在")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetNotification(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	var notification models.Notification
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "通知不存在")
		return
	}

	if !notification.IsRead {
		database.DB.Model(&notification).Update("is_read", true)
		notification.IsRead = true
	}

	c.JSON(http.StatusOK, notification)
}

func DeleteNotification(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的通知ID")
		return
	}

	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "通知不存在")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteReadNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result := database.DB.Where("user_id = ? AND is_read = ?", userID, true).Delete(&models.Notification{})
	c.JSON(http.StatusOK, gin.H{"success": true, "deleted": result.RowsAffected})
}

func MarkAllRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
