package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/internal/ws"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PaymentProofSummary struct {
	ID              uint   `json:"id"`
	Status          string `json:"status"`
	ImageUrl        string `json:"imageUrl"`
	PaymentMethod   string `json:"paymentMethod"`
	RejectReason    string `json:"rejectReason,omitempty"`
	CreatedAt       string `json:"createdAt"`
}

func GetAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	search := c.Query("search")
	proofStatus := c.Query("proof_status")

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
	if proofStatus != "" {
		subQuery := database.DB.Model(&models.PaymentProof{}).
			Select("order_id").
			Where("status = ?", proofStatus)
		query = query.Where("id IN ?", subQuery)
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

	proofMap := make(map[uint]models.PaymentProof)
	if len(orders) > 0 {
		orderIDs := make([]uint, len(orders))
		for i, o := range orders {
			orderIDs[i] = o.ID
		}
		var proofs []models.PaymentProof
		database.DB.Where("order_id IN ?", orderIDs).
			Preload("PaymentMethod").
			Order("created_at desc").
			Find(&proofs)
		for _, p := range proofs {
			if _, exists := proofMap[p.OrderID]; !exists {
				proofMap[p.OrderID] = p
			}
		}
	}

	type OrderWithUser struct {
		models.Order
		User         interface{}           `json:"user"`
		PaymentProof *PaymentProofSummary  `json:"paymentProof"`
	}
	var result []OrderWithUser
	for _, o := range orders {
		u, ok := userMap[o.UserID]
		userInfo := gin.H{"id": o.UserID}
		if ok {
			userInfo = gin.H{"id": u.ID, "name": u.Name, "email": u.Email}
		}
		var proofSummary *PaymentProofSummary
		if p, exists := proofMap[o.ID]; exists {
			proofSummary = &PaymentProofSummary{
				ID:            p.ID,
				Status:        p.Status,
				ImageUrl:      p.ImageUrl,
				PaymentMethod: p.PaymentMethod.Name,
				RejectReason:  p.RejectReason,
				CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}
		result = append(result, OrderWithUser{Order: o, User: userInfo, PaymentProof: proofSummary})
	}

	var pendingProofCount int64
	database.DB.Model(&models.PaymentProof{}).Where("status = ?", "pending").Count(&pendingProofCount)

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":            result,
		"total":             total,
		"page":              page,
		"limit":             limit,
		"totalPages":        totalPages,
		"pendingProofCount": pendingProofCount,
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

	var proofSummary *PaymentProofSummary
	var proof models.PaymentProof
	if err := database.DB.Where("order_id = ?", order.ID).
		Preload("PaymentMethod").
		Order("created_at desc").
		First(&proof).Error; err == nil {
		proofSummary = &PaymentProofSummary{
			ID:            proof.ID,
			Status:        proof.Status,
			ImageUrl:      proof.ImageUrl,
			PaymentMethod: proof.PaymentMethod.Name,
			RejectReason:  proof.RejectReason,
			CreatedAt:     proof.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order":        order,
		"user":         gin.H{"id": user.ID, "name": user.Name, "email": user.Email},
		"paymentProof": proofSummary,
	})
}

type UpdateOrderStatusInput struct {
	Status          string `json:"status" binding:"required"`
	TrackingCompany string `json:"trackingCompany"`
	TrackingNumber  string `json:"trackingNumber"`
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

	if input.Status == models.OrderStatusShipped {
		if strings.TrimSpace(input.TrackingCompany) == "" || strings.TrimSpace(input.TrackingNumber) == "" {
			utils.ErrorResponse(c, http.StatusBadRequest, "发货时请填写物流平台和快递单号")
			return
		}
	}

	oldStatus := order.Status
	order.Status = input.Status
	order.TrackingCompany = input.TrackingCompany
	order.TrackingNumber = input.TrackingNumber
	database.DB.Save(&order)

	database.DB.Preload("Items.Product").Preload("Address").First(&order, id)

	if oldStatus != input.Status {
		notification := models.Notification{
			UserID:  order.UserID,
			Type:    models.NotificationTypeOrderStatus,
			Link:    "/orders",
			OrderID: &order.ID,
		}
		if input.Status == models.OrderStatusShipped {
			notification.Title = "订单已发货"
			notification.Content = "您的订单 " + order.OrderNo + " 已发货，物流平台：" + order.TrackingCompany + "，快递单号：" + order.TrackingNumber
		} else {
			statusMap := map[string]string{
				"pending":    "待处理",
				"paid":       "已支付",
				"processing": "处理中",
				"completed":  "已完成",
				"cancelled":  "已取消",
			}
			statusText := statusMap[input.Status]
			if statusText == "" {
				statusText = input.Status
			}
			notification.Title = "订单状态更新"
			notification.Content = "您的订单 " + order.OrderNo + " 状态已更新为「" + statusText + "」"
		}
		database.DB.Create(&notification)
		ws.DefaultHub.SendToUser(notification.UserID, gin.H{
			"type":         "notification",
			"notification": notification,
		})
	}

	c.JSON(http.StatusOK, order)
}
