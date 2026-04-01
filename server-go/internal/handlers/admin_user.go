package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	userSortColumnMap := map[string]string{
		"id":        "id",
		"email":     "email",
		"name":      "name",
		"role":      "role",
		"createdAt": "created_at",
	}
	sortColumn, ok := userSortColumnMap[sortBy]
	if !ok {
		sortColumn = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	query := database.DB.Model(&models.User{})

	if search != "" {
		query = query.Where("email LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if role != "" {
		roles := strings.Split(role, ",")
		query = query.Where("role IN ?", roles)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * limit
	query.Order(sortColumn + " " + sortOrder).Offset(offset).Limit(limit).Find(&users)

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
	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "密码重置失败")
		return
	}

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

// BanUser godoc
// @Summary 封禁用户
// @Description 超级管理员封禁指定用户，该用户将无法登录和访问任何接口
// @Tags admin-users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/users/{id}/ban [put]
func BanUser(c *gin.Context) {
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

	if user.Role == "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "不能封禁超级管理员")
		return
	}

	user.IsBanned = true
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "用户已封禁", "user": user})
}

// UnbanUser godoc
// @Summary 解封用户
// @Description 超级管理员解封指定用户
// @Tags admin-users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/users/{id}/unban [put]
func UnbanUser(c *gin.Context) {
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

	user.IsBanned = false
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "用户已解封", "user": user})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 超级管理员删除指定用户，级联删除其所有关联数据（订单、购物车、收藏、地址、通知等）
// @Tags admin-users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
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

	if user.Role == "admin" {
		utils.ErrorResponse(c, http.StatusForbidden, "不能删除超级管理员")
		return
	}

	db := database.DB

	var orderIDs []uint
	db.Model(&models.Order{}).Where("user_id = ?", id).Pluck("id", &orderIDs)
	if len(orderIDs) > 0 {
		db.Where("order_id IN ?", orderIDs).Delete(&models.PaymentProof{})
		db.Where("order_id IN ?", orderIDs).Delete(&models.OrderItem{})
		db.Where("order_id IN ?", orderIDs).Delete(&models.OrderSummary{})
	}
	db.Where("user_id = ?", id).Delete(&models.Order{})

	db.Where("user_id = ?", id).Delete(&models.CartItem{})
	db.Where("user_id = ?", id).Delete(&models.Favorite{})
	db.Where("user_id = ?", id).Delete(&models.Address{})
	db.Where("user_id = ?", id).Delete(&models.Notification{})

	var conversationIDs []uint
	db.Model(&models.Conversation{}).Where("user_id = ?", id).Pluck("id", &conversationIDs)
	if len(conversationIDs) > 0 {
		db.Where("conversation_id IN ?", conversationIDs).Delete(&models.Rating{})
		db.Where("conversation_id IN ?", conversationIDs).Delete(&models.Message{})
	}
	db.Where("user_id = ?", id).Delete(&models.Conversation{})

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}
