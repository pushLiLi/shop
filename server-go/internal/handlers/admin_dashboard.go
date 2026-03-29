package handlers

import (
	"net/http"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var totalRevenue float64
	database.DB.Model(&models.Order{}).
		Where("status IN ?", []string{"completed", "shipped", "processing"}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalRevenue)

	var totalOrders int64
	database.DB.Model(&models.Order{}).Count(&totalOrders)

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var totalProducts int64
	database.DB.Model(&models.Product{}).Count(&totalProducts)

	today := time.Now().Truncate(24 * time.Hour)
	var todayOrders int64
	database.DB.Model(&models.Order{}).Where("created_at >= ?", today).Count(&todayOrders)

	var todayRevenue float64
	database.DB.Model(&models.Order{}).
		Where("created_at >= ? AND status IN ?", today, []string{"completed", "shipped", "processing"}).
		Select("COALESCE(SUM(total), 0)").
		Scan(&todayRevenue)

	var pendingOrders int64
	database.DB.Model(&models.Order{}).Where("status = ?", "pending").Count(&pendingOrders)

	c.JSON(http.StatusOK, gin.H{
		"totalRevenue":  totalRevenue,
		"totalOrders":   totalOrders,
		"totalUsers":    totalUsers,
		"totalProducts": totalProducts,
		"todayOrders":   todayOrders,
		"todayRevenue":  todayRevenue,
		"pendingOrders": pendingOrders,
	})
}

func GetDashboardRecentOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("Items").Order("created_at desc").Limit(10).Find(&orders)

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

	type RecentOrder struct {
		ID        uint      `json:"id"`
		OrderNo   string    `json:"orderNo"`
		UserName  string    `json:"userName"`
		Total     float64   `json:"total"`
		Status    string    `json:"status"`
		ItemCount int       `json:"itemCount"`
		CreatedAt time.Time `json:"createdAt"`
	}

	var result []RecentOrder
	for _, o := range orders {
		u := userMap[o.UserID]
		result = append(result, RecentOrder{
			ID:        o.ID,
			OrderNo:   o.OrderNo,
			UserName:  u.Name,
			Total:     o.Total,
			Status:    o.Status,
			ItemCount: len(o.Items),
			CreatedAt: o.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": result})
}

func GetDashboardLowStock(c *gin.Context) {
	var products []models.Product
	database.DB.Where("stock <= ? AND is_active = ?", 10, true).
		Order("stock asc").
		Limit(10).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetDashboardTopProducts(c *gin.Context) {
	type TopProduct struct {
		ProductID   uint    `json:"productId"`
		ProductName string  `json:"productName"`
		Image       string  `json:"imageUrl"`
		TotalSold   int     `json:"totalSold"`
		Revenue     float64 `json:"revenue"`
	}

	var results []TopProduct
	database.DB.Table("order_items").
		Select("order_items.product_id, products.name as product_name, products.image as image, SUM(order_items.quantity) as total_sold, SUM(order_items.quantity * order_items.price) as revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status IN ?", []string{"completed", "shipped", "processing"}).
		Group("order_items.product_id, products.name, products.image").
		Order("total_sold desc").
		Limit(10).
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"products": results})
}
