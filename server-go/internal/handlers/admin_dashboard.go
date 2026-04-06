package handlers

import (
	"net/http"
	"strconv"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	type CurrencyRevenue struct {
		Currency string  `json:"currency"`
		Revenue  float64 `json:"revenue"`
	}

	var currentRevenues []CurrencyRevenue
	database.DB.Table("order_items").
		Select("order_items.currency, SUM(order_items.quantity * order_items.price) as revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status IN ?", []string{"completed", "shipped", "processing"}).
		Group("order_items.currency").
		Scan(&currentRevenues)

	var summaryRevenue float64
	database.DB.Model(&models.OrderSummary{}).
		Select("COALESCE(SUM(total_revenue), 0)").
		Scan(&summaryRevenue)

	totalRevenueByCurrency := make(map[string]float64)
	for _, cr := range currentRevenues {
		totalRevenueByCurrency[cr.Currency] = cr.Revenue
	}
	if summaryRevenue > 0 {
		totalRevenueByCurrency["CNY"] += summaryRevenue
	}

	var totalOrders int64
	database.DB.Model(&models.Order{}).Count(&totalOrders)
	var summaryOrders int64
	database.DB.Model(&models.OrderSummary{}).Select("COALESCE(SUM(total_orders), 0)").Scan(&summaryOrders)
	totalOrders += summaryOrders

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var totalProducts int64
	database.DB.Model(&models.Product{}).Count(&totalProducts)

	today := time.Now().Truncate(24 * time.Hour)
	var todayOrders int64
	database.DB.Model(&models.Order{}).Where("created_at >= ?", today).Count(&todayOrders)

	var todayRevenues []CurrencyRevenue
	database.DB.Table("order_items").
		Select("order_items.currency, SUM(order_items.quantity * order_items.price) as revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.created_at >= ? AND orders.status IN ?", today, []string{"completed", "shipped", "processing"}).
		Group("order_items.currency").
		Scan(&todayRevenues)

	todayRevenueByCurrency := make(map[string]float64)
	for _, tr := range todayRevenues {
		todayRevenueByCurrency[tr.Currency] = tr.Revenue
	}

	var pendingOrders int64
	database.DB.Model(&models.Order{}).Where("status = ?", "pending").Count(&pendingOrders)

	c.JSON(http.StatusOK, gin.H{
		"totalRevenueByCurrency": totalRevenueByCurrency,
		"totalOrders":            totalOrders,
		"totalUsers":             totalUsers,
		"totalProducts":          totalProducts,
		"todayOrders":            todayOrders,
		"todayRevenueByCurrency": todayRevenueByCurrency,
		"pendingOrders":          pendingOrders,
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
		Currency  string    `json:"currency"`
		Status    string    `json:"status"`
		ItemCount int       `json:"itemCount"`
		CreatedAt time.Time `json:"createdAt"`
	}

	var result []RecentOrder
	for _, o := range orders {
		u := userMap[o.UserID]
		currency := "CNY"
		if len(o.Items) > 0 {
			currency = o.Items[0].Currency
			for _, item := range o.Items {
				if item.Currency != currency {
					currency = "mixed"
					break
				}
			}
		}
		result = append(result, RecentOrder{
			ID:        o.ID,
			OrderNo:   o.OrderNo,
			UserName:  u.Name,
			Total:     o.Total,
			Currency:  currency,
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

func GetRevenueByDate(c *gin.Context) {
	days := 7
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	startDate := time.Now().AddDate(0, 0, -days+1).Truncate(24 * time.Hour)

	type DailyRevenue struct {
		Date              string             `json:"date"`
		RevenueByCurrency map[string]float64 `json:"revenueByCurrency"`
		Orders            int                `json:"orders"`
	}

	var results []DailyRevenue
	summaryMap := make(map[string]models.OrderSummary)
	var summaries []models.OrderSummary
	database.DB.Where("date >= ?", startDate.Format("2006-01-02")).Find(&summaries)
	for _, s := range summaries {
		summaryMap[s.Date] = s
	}

	for i := 0; i < days; i++ {
		day := startDate.AddDate(0, 0, i)
		nextDay := day.AddDate(0, 0, 1)
		dateStr := day.Format("2006-01-02")

		type CurrencyRevenue struct {
			Currency string  `json:"currency"`
			Revenue  float64 `json:"revenue"`
		}
		var dayRevenues []CurrencyRevenue
		database.DB.Table("order_items").
			Select("order_items.currency, SUM(order_items.quantity * order_items.price) as revenue").
			Joins("JOIN orders ON orders.id = order_items.order_id").
			Where("orders.created_at >= ? AND orders.created_at < ? AND orders.status IN ?", day, nextDay, []string{"completed", "shipped", "processing"}).
			Group("order_items.currency").
			Scan(&dayRevenues)

		var orders int64
		database.DB.Model(&models.Order{}).
			Where("created_at >= ? AND created_at < ? AND status IN ?", day, nextDay, []string{"completed", "shipped", "processing"}).
			Count(&orders)

		revenueByCurrency := make(map[string]float64)
		for _, dr := range dayRevenues {
			revenueByCurrency[dr.Currency] = dr.Revenue
		}

		if summary, ok := summaryMap[dateStr]; ok {
			revenueByCurrency["CNY"] += summary.TotalRevenue
			orders += int64(summary.CompletedOrders)
		}

		results = append(results, DailyRevenue{
			Date:              dateStr,
			RevenueByCurrency: revenueByCurrency,
			Orders:            int(orders),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func GetDashboardTopProducts(c *gin.Context) {
	type TopProduct struct {
		ProductID         uint               `json:"productId"`
		ProductName       string             `json:"productName"`
		Image             string             `json:"imageUrl"`
		TotalSold         int                `json:"totalSold"`
		RevenueByCurrency map[string]float64 `json:"revenueByCurrency"`
	}

	type ProductRevenue struct {
		ProductID   uint    `json:"productId"`
		ProductName string  `json:"productName"`
		Image       string  `json:"imageUrl"`
		TotalSold   int     `json:"totalSold"`
		Currency    string  `json:"currency"`
		Revenue     float64 `json:"revenue"`
	}

	var rawResults []ProductRevenue
	database.DB.Table("order_items").
		Select("order_items.product_id, products.name as product_name, products.image as image, SUM(order_items.quantity) as total_sold, order_items.currency, SUM(order_items.quantity * order_items.price) as revenue").
		Joins("JOIN products ON products.id = order_items.product_id").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status IN ?", []string{"completed", "shipped", "processing"}).
		Group("order_items.product_id, products.name, products.image, order_items.currency").
		Order("total_sold desc").
		Scan(&rawResults)

	productMap := make(map[uint]*TopProduct)
	var productOrder []uint
	for _, r := range rawResults {
		if _, exists := productMap[r.ProductID]; !exists {
			productMap[r.ProductID] = &TopProduct{
				ProductID:         r.ProductID,
				ProductName:       r.ProductName,
				Image:             r.Image,
				TotalSold:         r.TotalSold,
				RevenueByCurrency: make(map[string]float64),
			}
			productOrder = append(productOrder, r.ProductID)
		}
		productMap[r.ProductID].RevenueByCurrency[r.Currency] += r.Revenue
	}

	var results []TopProduct
	for _, id := range productOrder {
		results = append(results, *productMap[id])
	}

	if len(results) > 10 {
		results = results[:10]
	}

	c.JSON(http.StatusOK, gin.H{"products": results})
}
