package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary 获取订单列表
// @Description 获取当前用户的订单列表
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var orders []models.Order
	database.DB.Where("user_id = ?", userID).
		Preload("Items.Product").
		Preload("Address").
		Order("created_at desc").
		Find(&orders)

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 从购物车商品创建订单
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CreateOrderInput true "订单信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.AddressID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "请选择收货地址")
		return
	}

	var address models.Address
	if err := database.DB.Where("id = ? AND user_id = ?", input.AddressID, userID).First(&address).Error; err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "地址不存在或无权使用")
		return
	}

	var cartItems []models.CartItem
	database.DB.Where("user_id = ?", userID).Preload("Product").Find(&cartItems)

	if len(cartItems) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Cart is empty")
		return
	}

	var total float64
	for _, item := range cartItems {
		if item.Product.ID > 0 {
			total += item.Product.Price * float64(item.Quantity)
		}
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	order := models.Order{
		OrderNo:   utils.GenerateOrderNo(),
		UserID:    uid,
		AddressID: input.AddressID,
		Total:     total,
		Remark:    input.Remark,
		Status:    "pending",
	}

	if err := database.DB.Create(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create order")
		return
	}

	for _, item := range cartItems {
		if item.Product.ID > 0 {
			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			}
			database.DB.Create(&orderItem)
		}
	}

	database.DB.Where("user_id = ?", userID).Delete(&models.CartItem{})

	c.JSON(http.StatusOK, gin.H{"success": true, "orderId": order.ID, "orderNo": order.OrderNo})
}

// GetOrder godoc
// @Summary 获取订单详情
// @Description 获取指定订单的详细信息
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	param := c.Param("id")

	var order models.Order
	query := database.DB.Where("user_id = ?", userID)

	if id, err := strconv.Atoi(param); err == nil {
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("order_no = ?", param)
	}

	if err := query.Preload("Items.Product").Preload("Address").First(&order).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Order not found")
		return
	}

	var proofSummary map[string]interface{}
	var proof models.PaymentProof
	if err := database.DB.Where("order_id = ?", order.ID).
		Preload("PaymentMethod").
		Order("created_at desc").
		First(&proof).Error; err == nil {
		proofSummary = map[string]interface{}{
			"id":              proof.ID,
			"status":          proof.Status,
			"imageUrl":        proof.ImageUrl,
			"paymentMethodId": proof.PaymentMethodID,
			"paymentMethod":   proof.PaymentMethod.Name,
			"rejectReason":    proof.RejectReason,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"order":        order,
		"paymentProof": proofSummary,
	})
}
