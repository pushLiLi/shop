package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

var validContactTypes = []string{"whatsapp", "wechat", "qq", "telegram", "phone", "email", "custom"}

func isValidContactType(t string) bool {
	for _, v := range validContactTypes {
		if v == t {
			return true
		}
	}
	return false
}

func GetContactMethods(c *gin.Context) {
	var methods []models.ContactMethod
	database.DB.Where("is_active = ?", true).Order("sort_order asc, id asc").Find(&methods)

	c.JSON(http.StatusOK, gin.H{"contactMethods": methods})
}

func GetAdminContactMethods(c *gin.Context) {
	var methods []models.ContactMethod
	database.DB.Order("sort_order asc, id asc").Find(&methods)

	c.JSON(http.StatusOK, gin.H{"contactMethods": methods})
}

func CreateContactMethod(c *gin.Context) {
	var input models.CreateContactMethodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidContactType(input.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的联系方式类型"})
		return
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	method := models.ContactMethod{
		Type:      input.Type,
		Label:     input.Label,
		Value:     input.Value,
		QRCodeUrl: input.QRCodeUrl,
		IsActive:  isActive,
		SortOrder: input.SortOrder,
	}

	if err := database.DB.Create(&method).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建联系方式失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "contactMethod": method})
}

func UpdateContactMethod(c *gin.Context) {
	id := c.Param("id")

	var method models.ContactMethod
	if err := database.DB.First(&method, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "联系方式不存在"})
		return
	}

	var input models.UpdateContactMethodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Type != "" && !isValidContactType(input.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的联系方式类型"})
		return
	}

	updates := map[string]interface{}{}
	if input.Type != "" {
		updates["type"] = input.Type
	}
	if input.Label != "" {
		updates["label"] = input.Label
	}
	if input.Value != "" {
		updates["value"] = input.Value
	}
	if input.QRCodeUrl != "" {
		updates["qr_code_url"] = input.QRCodeUrl
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	updates["sort_order"] = input.SortOrder

	database.DB.Model(&method).Updates(updates)
	database.DB.First(&method, id)

	c.JSON(http.StatusOK, gin.H{"success": true, "contactMethod": method})
}

func DeleteContactMethod(c *gin.Context) {
	id := c.Param("id")

	var method models.ContactMethod
	if err := database.DB.First(&method, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "联系方式不存在"})
		return
	}

	database.DB.Delete(&method)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
