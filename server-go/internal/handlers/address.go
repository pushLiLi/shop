package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetAddresses(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户不存在，请重新登录")
		return
	}

	var addresses []models.Address
	database.DB.Where("user_id = ?", userID).Order("is_default desc").Find(&addresses)

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func CreateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户不存在，请重新登录")
		return
	}

	var count int64
	database.DB.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)
	if count >= 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "最多只能保存5个地址")
		return
	}

	var input models.CreateAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.FullName == "" || input.AddressLine1 == "" || input.City == "" || input.State == "" || input.ZipCode == "" || input.Phone == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "请填写完整的地址信息")
		return
	}

	if input.IsDefault {
		database.DB.Model(&models.Address{}).Where("user_id = ? AND is_default = ?", userID, true).Update("is_default", false)
	}

	address := models.Address{
		UserID:       userID.(uint),
		FullName:     input.FullName,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		ZipCode:      input.ZipCode,
		Phone:        input.Phone,
		IsDefault:    input.IsDefault || count == 0,
	}

	database.DB.Create(&address)

	c.JSON(http.StatusOK, gin.H{"success": true, "address": address})
}

func UpdateAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户不存在，请重新登录")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var existing models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "地址不存在")
		return
	}

	var input models.UpdateAddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.FullName != "" {
		existing.FullName = input.FullName
	}
	if input.AddressLine1 != "" {
		existing.AddressLine1 = input.AddressLine1
	}
	if input.AddressLine2 != "" {
		existing.AddressLine2 = input.AddressLine2
	}
	if input.City != "" {
		existing.City = input.City
	}
	if input.State != "" {
		existing.State = input.State
	}
	if input.ZipCode != "" {
		existing.ZipCode = input.ZipCode
	}
	if input.Phone != "" {
		existing.Phone = input.Phone
	}

	database.DB.Save(&existing)

	c.JSON(http.StatusOK, gin.H{"success": true, "address": existing})
}

func DeleteAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户不存在，请重新登录")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var existing models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "地址不存在")
		return
	}

	database.DB.Delete(&existing)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func SetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var existing models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&existing).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "地址不存在")
		return
	}

	database.DB.Model(&models.Address{}).Where("user_id = ? AND is_default = ?", userID, true).Update("is_default", false)
	existing.IsDefault = true
	database.DB.Save(&existing)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
