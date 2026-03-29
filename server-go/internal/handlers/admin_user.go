package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAdminUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	role := c.Query("role")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := database.DB.Model(&models.User{})

	if search != "" {
		query = query.Where("email LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		query = query.Where("role = ?", role)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * limit
	query.Order("created_at desc").Offset(offset).Limit(limit).Find(&users)

	type UserWithStats struct {
		models.User
		OrderCount int64   `json:"orderCount"`
		TotalSpent float64 `json:"totalSpent"`
	}

	var result []UserWithStats
	for _, u := range users {
		var orderCount int64
		var totalSpent float64
		database.DB.Model(&models.Order{}).Where("user_id = ?", u.ID).Count(&orderCount)
		database.DB.Model(&models.Order{}).
			Where("user_id = ? AND status IN ?", u.ID, []string{"completed", "shipped", "processing"}).
			Select("COALESCE(SUM(total), 0)").
			Scan(&totalSpent)

		result = append(result, UserWithStats{
			User:       u,
			OrderCount: orderCount,
			TotalSpent: totalSpent,
		})
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"users":      result,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func GetAdminUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	var orders []models.Order
	database.DB.Where("user_id = ?", id).
		Preload("Items.Product").
		Order("created_at desc").
		Limit(20).
		Find(&orders)

	var addresses []models.Address
	database.DB.Where("user_id = ?", id).
		Order("is_default desc, created_at desc").
		Find(&addresses)

	c.JSON(http.StatusOK, gin.H{
		"user":      user,
		"orders":    orders,
		"addresses": addresses,
	})
}

func ResetUserPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	currentUser, _ := c.Get("user")
	currentRole := currentUser.(models.User).Role
	if currentRole == "service" && user.Role == "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "无权修改超级管理员密码")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码重置失败")
		return
	}

	user.Password = string(hashedPassword)
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "密码已重置为 123456，请通知用户尽快修改"})
}

type UpdateUserRoleInput struct {
	Role string `json:"role" binding:"required"`
}

func UpdateUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.Role != "admin" && input.Role != "service" && input.Role != "customer" {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "用户不存在")
		return
	}

	user.Role = input.Role
	database.DB.Save(&user)

	c.JSON(http.StatusOK, user)
}
