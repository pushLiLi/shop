package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetPaymentMethods(c *gin.Context) {
	var methods []models.PaymentMethod
	database.DB.Where("is_active = ?", true).Order("sort_order asc, id asc").Find(&methods)

	c.JSON(http.StatusOK, gin.H{"paymentMethods": methods})
}

func GetAdminPaymentMethods(c *gin.Context) {
	var methods []models.PaymentMethod
	database.DB.Order("sort_order asc, id asc").Find(&methods)

	c.JSON(http.StatusOK, gin.H{"paymentMethods": methods})
}

func CreatePaymentMethod(c *gin.Context) {
	var input models.CreatePaymentMethodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	method := models.PaymentMethod{
		Name:         input.Name,
		QRCodeUrl:    input.QRCodeUrl,
		Instructions: input.Instructions,
		IsActive:     isActive,
		SortOrder:    input.SortOrder,
	}

	if err := database.DB.Create(&method).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建付款方式失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "paymentMethod": method})
}

func UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var method models.PaymentMethod
	if err := database.DB.First(&method, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "付款方式不存在"})
		return
	}

	var input models.UpdatePaymentMethodInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.QRCodeUrl != "" {
		updates["qr_code_url"] = input.QRCodeUrl
	}
	if input.Instructions != "" {
		updates["instructions"] = input.Instructions
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	updates["sort_order"] = input.SortOrder

	database.DB.Model(&method).Updates(updates)
	database.DB.First(&method, id)

	c.JSON(http.StatusOK, gin.H{"success": true, "paymentMethod": method})
}

func DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")

	var method models.PaymentMethod
	if err := database.DB.First(&method, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "付款方式不存在"})
		return
	}

	database.DB.Delete(&method)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
