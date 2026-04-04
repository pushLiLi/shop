package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/internal/ws"
	"bycigar-server/pkg/email"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentProofSummary struct {
	ID            uint   `json:"id"`
	Status        string `json:"status"`
	ImageUrl      string `json:"imageUrl"`
	PaymentMethod string `json:"paymentMethod"`
	RejectReason  string `json:"rejectReason,omitempty"`
	CreatedAt     string `json:"createdAt"`
}

func buildOrderQuery(c *gin.Context) (*gorm.DB, string) {
	query := database.DB.Model(&models.Order{})
	quickFilter := c.Query("quick_filter")
	status := c.Query("status")
	search := c.Query("search")
	proofStatus := c.Query("proof_status")
	sortBy := c.DefaultQuery("sortBy", "createdAt")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	orderSortColumnMap := map[string]string{
		"id":        "id",
		"orderNo":   "order_no",
		"total":     "total",
		"status":    "status",
		"createdAt": "created_at",
	}
	sortColumn, ok := orderSortColumnMap[sortBy]
	if !ok {
		sortColumn = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	if search != "" {
		query = query.Where("order_no LIKE ?", "%"+search+"%")
	}

	if quickFilter != "" {
		switch quickFilter {
		case "pending_proof":
			var orderIDs []uint
			database.DB.Model(&models.PaymentProof{}).
				Where("status = ?", "pending").Pluck("order_id", &orderIDs)
			if len(orderIDs) > 0 {
				query = query.Where("id IN ?", orderIDs)
			} else {
				query = query.Where("1 = 0")
			}
		case "to_ship":
			var orderIDs []uint
			database.DB.Model(&models.PaymentProof{}).
				Where("status = ?", "approved").Pluck("order_id", &orderIDs)
			if len(orderIDs) > 0 {
				query = query.Where("id IN ? AND status = ?", orderIDs, models.OrderStatusProcessing)
			} else {
				query = query.Where("1 = 0")
			}
		case "shipped":
			query = query.Where("status = ?", models.OrderStatusShipped)
		case "completed":
			query = query.Where("status = ?", models.OrderStatusCompleted)
		}
	} else {
		if status != "" {
			query = query.Where("status = ?", status)
		}
		if proofStatus != "" {
			var orderIDs []uint
			database.DB.Model(&models.PaymentProof{}).
				Where("status = ?", proofStatus).Pluck("order_id", &orderIDs)
			if len(orderIDs) > 0 {
				query = query.Where("id IN ?", orderIDs)
			} else {
				query = query.Where("1 = 0")
			}
		}
	}
	return query, sortColumn + " " + sortOrder
}

func GetAdminOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	query, orderClause := buildOrderQuery(c)

	var total int64
	query.Count(&total)

	var orders []models.Order
	offset := (page - 1) * limit
	query.Preload("Items.Product").Preload("Address").Preload("Items").
		Order(orderClause).
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
		User         interface{}          `json:"user"`
		PaymentProof *PaymentProofSummary `json:"paymentProof"`
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
		if ws.DefaultHub != nil {
			ws.DefaultHub.SendToUser(notification.UserID, gin.H{
				"type":         "notification",
				"notification": notification,
			})
		}

		if input.Status == models.OrderStatusShipped {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("SendShippingNotification panic recovered (orderNo=%s): %v", order.OrderNo, r)
					}
				}()
				email.SendShippingNotification(order)
			}()
		}
	}

	c.JSON(http.StatusOK, order)
}

func ExportAdminOrders(c *gin.Context) {
	query, orderClause := buildOrderQuery(c)

	var total int64
	query.Count(&total)
	if total > 10000 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Too many records to export (max: 10000). Please use filters to narrow down the results.")
		return
	}

	var orders []models.Order
	query.Preload("Items.Product").Preload("Address").
		Order(orderClause).
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

	filename := fmt.Sprintf("orders_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	headers := []string{"订单号", "用户姓名", "用户邮箱", "商品明细", "总金额", "订单状态", "备注", "收件人", "电话", "收货地址", "物流平台", "快递单号", "付款凭证状态", "付款方式", "下单时间"}
	writer.Write(headers)

	statusLabels := map[string]string{
		"pending":    "待处理",
		"processing": "处理中",
		"shipped":    "已发货",
		"completed":  "已完成",
		"cancelled":  "已取消",
	}
	proofStatusLabels := map[string]string{
		"pending":  "待审核",
		"approved": "已通过",
		"rejected": "已驳回",
	}

	for _, order := range orders {
		var items []string
		for _, item := range order.Items {
			items = append(items, fmt.Sprintf("%s x%d", item.Product.Name, item.Quantity))
		}
		itemStr := strings.Join(items, "; ")

		statusText := statusLabels[order.Status]
		if statusText == "" {
			statusText = order.Status
		}

		var userName, userEmail string
		if u, ok := userMap[order.UserID]; ok {
			userName = u.Name
			userEmail = u.Email
		}

		addr := order.Address
		addrParts := []string{addr.AddressLine1}
		if addr.AddressLine2 != "" {
			addrParts = append(addrParts, addr.AddressLine2)
		}
		addrParts = append(addrParts, addr.City, addr.State, addr.ZipCode)
		fullAddr := strings.Join(addrParts, " ")

		proofStatusText := "未提交"
		var paymentMethod string
		if p, exists := proofMap[order.ID]; exists {
			proofStatusText = proofStatusLabels[p.Status]
			if proofStatusText == "" {
				proofStatusText = p.Status
			}
			paymentMethod = p.PaymentMethod.Name
		}
		if paymentMethod == "" {
			paymentMethod = "-"
		}

		row := []string{
			order.OrderNo,
			userName,
			userEmail,
			itemStr,
			fmt.Sprintf("%.2f", order.Total),
			statusText,
			order.Remark,
			addr.FullName,
			addr.Phone,
			fullAddr,
			order.TrackingCompany,
			order.TrackingNumber,
			proofStatusText,
			paymentMethod,
			order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(row)
	}
}
