package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query := database.DB.Model(&models.Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("order_no LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * limit
	query.Preload("Items.Product").Preload("Address").Preload("Items").
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&orders)

	var userIDs []uint
	for _, o := range orders {
		userIDs = append(userIDs, o.UserID)
	}
	var users []models.User
	if len(userIDs) > 0 {
		database.DB.Where("id IN ?", userIDs).Find(&users)
	}
	userMap := make(map[uint]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	type OrderWithUser struct {
		models.Order
		User interface{} `json:"user"`
	}
	var result []OrderWithUser
	for _, o := range orders {
		u, ok := userMap[o.UserID]
		userInfo := gin.H{"id": o.UserID}
		if ok {
			userInfo = gin.H{"id": u.ID, "name": u.Name, "email": u.Email}
		}
		result = append(result, OrderWithUser{Order: o, User: userInfo})
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":     result,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func GetAdminOrder(c *gin.Context) {
	param := c.Param("id")

	var order models.Order
	query := database.DB.Preload("Items.Product").Preload("Address")

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("order_no = ?", param)
	}

	if err := query.First(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "订单不存在")
		return
	}

	var user models.User
	database.DB.First(&user, order.UserID)

	c.JSON(http.StatusOK, gin.H{
		"order": order,
		"user":  gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
	})
}

type UpdateOrderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的订单ID")
		return
	}

	var input UpdateOrderStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "订单不存在")
		return
	}

	validNext, ok := models.ValidOrderStatusTransitions[order.Status]
	if !ok {
		utils.ErrorResponse(c, http.StatusBadRequest, "当前订单状态无效")
		return
	}

	found := false
	for _, s := range validNext {
		if s == input.Status {
			found = true
			break
		}
	}
	if !found {
		utils.ErrorResponse(c, http.StatusBadRequest, "不允许从 "+order.Status+" 变更为 "+input.Status)
		return
	}

	oldStatus := order.Status
	order.Status = input.Status
	database.DB.Save(&order)

	database.DB.Preload("Items.Product").Preload("Address").First(&order, id)

	if oldStatus != input.Status {
		statusMap := map[string]string{
			"pending":    "待处理",
			"paid":       "已支付",
			"processing": "处理中",
			"shipped":    "已发货",
			"completed":  "已完成",
			"cancelled":  "已取消",
		}
		statusText := statusMap[input.Status]
		if statusText == "" {
			statusText = input.Status
		}
		notification := models.Notification{
			UserID:  order.UserID,
			Type:    models.NotificationTypeOrderStatus,
			Title:   "订单状态更新",
			Content: "您的订单 " + order.OrderNo + " 状态已更新为「" + statusText + "」",
			Link:    "/orders",
			OrderID: &order.ID,
		}
		database.DB.Create(&notification)
	}

	c.JSON(http.StatusOK, order)
}
