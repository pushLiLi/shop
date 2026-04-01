package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

type CleanupInput struct {
	Orders        bool `json:"orders"`
	Users         bool `json:"users"`
	Conversations bool `json:"conversations"`
	Products      bool `json:"products"`
}

type CleanupResult struct {
	OrdersDeleted         int64 `json:"ordersDeleted"`
	OrderItemsDeleted     int64 `json:"orderItemsDeleted"`
	PaymentProofsDeleted  int64 `json:"paymentProofsDeleted"`
	OrderSummariesDeleted int64 `json:"orderSummariesDeleted"`
	UsersDeleted          int64 `json:"usersDeleted"`
	CartItemsDeleted      int64 `json:"cartItemsDeleted"`
	FavoritesDeleted      int64 `json:"favoritesDeleted"`
	AddressesDeleted      int64 `json:"addressesDeleted"`
	NotificationsDeleted  int64 `json:"notificationsDeleted"`
	ConversationsDeleted  int64 `json:"conversationsDeleted"`
	MessagesDeleted       int64 `json:"messagesDeleted"`
	RatingsDeleted        int64 `json:"ratingsDeleted"`
	QuickRepliesDeleted   int64 `json:"quickRepliesDeleted"`
	ProductsDeleted       int64 `json:"productsDeleted"`
	CategoriesDeleted     int64 `json:"categoriesDeleted"`
}

// BatchCleanup godoc
// @Summary 批量清理数据
// @Description 超级管理员一键清理测试数据（订单、用户、会话、商品等）
// @Tags admin-cleanup
// @Accept json
// @Produce json
// @Param body body CleanupInput true "清理选项"
// @Success 200 {object} CleanupResult
// @Router /api/admin/cleanup [post]
func BatchCleanup(c *gin.Context) {
	var input CleanupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if !input.Orders && !input.Users && !input.Conversations && !input.Products {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请至少选择一项清理内容"})
		return
	}

	var result CleanupResult
	db := database.DB

	if input.Orders {
		var orderIDs []uint
		db.Model(&models.Order{}).Pluck("id", &orderIDs)
		if len(orderIDs) > 0 {
			result.PaymentProofsDeleted = db.Where("order_id IN ?", orderIDs).Delete(&models.PaymentProof{}).RowsAffected
			result.OrderItemsDeleted = db.Where("order_id IN ?", orderIDs).Delete(&models.OrderItem{}).RowsAffected
		}
		result.OrdersDeleted = db.Where("1 = 1").Delete(&models.Order{}).RowsAffected
		result.OrderSummariesDeleted = db.Where("1 = 1").Delete(&models.OrderSummary{}).RowsAffected
	}

	if input.Users {
		var customerIDs []uint
		db.Model(&models.User{}).Where("role = ?", "customer").Pluck("id", &customerIDs)
		if len(customerIDs) > 0 {
			result.CartItemsDeleted = db.Where("user_id IN ?", customerIDs).Delete(&models.CartItem{}).RowsAffected
			result.FavoritesDeleted = db.Where("user_id IN ?", customerIDs).Delete(&models.Favorite{}).RowsAffected
			result.AddressesDeleted = db.Where("user_id IN ?", customerIDs).Delete(&models.Address{}).RowsAffected
			result.NotificationsDeleted = db.Where("user_id IN ?", customerIDs).Delete(&models.Notification{}).RowsAffected
		}
		result.UsersDeleted = db.Where("role = ?", "customer").Delete(&models.User{}).RowsAffected
	}

	if input.Conversations {
		var conversationIDs []uint
		db.Model(&models.Conversation{}).Pluck("id", &conversationIDs)
		if len(conversationIDs) > 0 {
			result.RatingsDeleted = db.Where("conversation_id IN ?", conversationIDs).Delete(&models.Rating{}).RowsAffected
			result.MessagesDeleted = db.Where("conversation_id IN ?", conversationIDs).Delete(&models.Message{}).RowsAffected
		}
		result.ConversationsDeleted = db.Where("1 = 1").Delete(&models.Conversation{}).RowsAffected
		result.QuickRepliesDeleted = db.Where("1 = 1").Delete(&models.QuickReply{}).RowsAffected
	}

	if input.Products {
		var productIDs []uint
		db.Model(&models.Product{}).Unscoped().Pluck("id", &productIDs)
		if len(productIDs) > 0 {
			db.Where("product_id IN ?", productIDs).Delete(&models.CartItem{})
			db.Where("product_id IN ?", productIDs).Delete(&models.Favorite{})
		}
		result.ProductsDeleted = db.Exec("DELETE FROM products").RowsAffected
		result.CategoriesDeleted = db.Exec("DELETE FROM categories").RowsAffected
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "清理完成",
		"result":  result,
	})
}
