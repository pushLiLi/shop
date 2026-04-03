package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetAddresses godoc
// @Summary 获取地址列表
// @Description 获取当前用户的收货地址列表
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /addresses [get]
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

// CreateAddress godoc
// @Summary 创建地址
// @Description 创建新的收货地址
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CreateAddressInput true "地址信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /addresses [post]
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

	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	address := models.Address{
		UserID:       uid,
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

// UpdateAddress godoc
// @Summary 更新地址
// @Description 更新指定的收货地址
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "地址ID"
// @Param input body models.UpdateAddressInput true "地址信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /addresses/{id} [put]
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

// DeleteAddress godoc
// @Summary 删除地址
// @Description 删除指定的收货地址
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "地址ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /addresses/{id} [delete]
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

	if err := database.DB.Delete(&existing).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除失败，该地址可能已被订单引用")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// SetDefaultAddress godoc
// @Summary 设置默认地址
// @Description 将指定地址设为默认收货地址
// @Tags addresses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "地址ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /addresses/{id}/default [put]
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
