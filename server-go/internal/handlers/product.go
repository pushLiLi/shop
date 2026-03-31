package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

var (
	categoriesCache     json.RawMessage
	categoriesCacheTime time.Time
	categoriesCacheMu   sync.RWMutex
	categoriesCacheTTL  = 5 * time.Minute
)

// GetProducts godoc
// @Summary 获取产品列表
// @Description 获取产品列表，支持分页、搜索、分类筛选、价格区间和排序
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(12)
// @Param search query string false "搜索关键词（空格分隔多词）"
// @Param category query string false "分类 slug"
// @Param categoryId query int false "分类 ID"
// @Param featured query string false "是否精选"
// @Param minPrice query number false "最低价格"
// @Param maxPrice query number false "最高价格"
// @Param sortBy query string false "排序字段" default(createdAt)
// @Param sortOrder query string false "排序方向 (asc/desc)" default(desc)
// @Success 200 {object} map[string]interface{}
// @Router /products [get]
func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	categorySlug := c.Query("category")
	categoryIDStr := c.Query("categoryId")
	featured := c.Query("featured")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	sortBy := c.DefaultQuery("sortBy", c.DefaultQuery("sort", "createdAt"))
	sortOrder := c.DefaultQuery("sortOrder", c.DefaultQuery("order", "desc"))

	sortColumnMap := map[string]string{
		"createdAt": "created_at",
		"price":     "price",
		"name":      "name",
		"id":        "id",
	}
	sortColumn, ok := sortColumnMap[sortBy]
	if !ok {
		sortColumn = "created_at"
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 12
	}

	query := database.DB.Model(&models.Product{}).Where("is_active = ?", true)

	if categorySlug != "" {
		var category models.Category
		if database.DB.Where("slug = ?", categorySlug).First(&category).Error == nil {
			var childCats []models.Category
			database.DB.Where("parent_id = ?", category.ID).Find(&childCats)
			ids := []uint{category.ID}
			for _, child := range childCats {
				ids = append(ids, child.ID)
			}
			query = query.Where("category_id IN ?", ids)
		}
	}

	if categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			var childCats []models.Category
			database.DB.Where("parent_id = ?", categoryID).Find(&childCats)
			ids := []uint{uint(categoryID)}
			for _, child := range childCats {
				ids = append(ids, child.ID)
			}
			query = query.Where("category_id IN ?", ids)
		}
	}

	if featured == "true" {
		query = query.Where("is_featured = ?", true)
	}

	if search != "" {
		keywords := strings.Fields(search)
		for _, kw := range keywords {
			query = query.Where("name LIKE ? OR description LIKE ?", "%"+kw+"%", "%"+kw+"%")
		}
	}

	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("price >= ?", minPrice)
		}
	}

	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("price <= ?", maxPrice)
		}
	}

	var total int64
	query.Count(&total)

	stockPriority := "(CASE WHEN stock > 0 THEN 1 ELSE 0 END) DESC, "
	orderClause := stockPriority + sortColumn + " DESC"
	if sortOrder == "asc" {
		orderClause = stockPriority + sortColumn + " ASC"
	}

	var products []models.Product
	offset := (page - 1) * limit
	query.Preload("Category").Order(orderClause).Offset(offset).Limit(limit).Find(&products)

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"products":   products,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

// GetProduct godoc
// @Summary 获取单个产品
// @Description 根据ID获取产品详情
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetCategories godoc
// @Summary 获取分类列表
// @Description 获取所有分类及其子分类
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	categoriesCacheMu.RLock()
	if categoriesCache != nil && time.Since(categoriesCacheTime) < categoriesCacheTTL {
		cached := categoriesCache
		categoriesCacheMu.RUnlock()
		c.Data(http.StatusOK, "application/json; charset=utf-8", cached)
		return
	}
	categoriesCacheMu.RUnlock()

	categoriesCacheMu.Lock()
	defer categoriesCacheMu.Unlock()

	if categoriesCache != nil && time.Since(categoriesCacheTime) < categoriesCacheTTL {
		c.Data(http.StatusOK, "application/json; charset=utf-8", categoriesCache)
		return
	}

	type CategoryWithCount struct {
		models.Category
		ProductCount int64 `json:"_count"`
	}

	var categories []models.Category
	database.DB.Where("parent_id IS NULL").Preload("Children").Find(&categories)

	result := make([]CategoryWithCount, 0)
	for _, cat := range categories {
		ids := []uint{cat.ID}
		for _, child := range cat.Children {
			ids = append(ids, child.ID)
		}
		var count int64
		database.DB.Model(&models.Product{}).Where("category_id IN ?", ids).Count(&count)
		result = append(result, CategoryWithCount{
			Category:     cat,
			ProductCount: count,
		})
	}

	data, _ := json.Marshal(result)
	categoriesCache = json.RawMessage(data)
	categoriesCacheTime = time.Now()

	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
}

func InvalidateCategoriesCache() {
	categoriesCacheMu.Lock()
	defer categoriesCacheMu.Unlock()
	categoriesCache = nil
	categoriesCacheTime = time.Time{}
}

// GetTopSelling godoc
// @Summary 获取热销商品
// @Description 根据订单项聚合销量，返回销量最高的商品
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "数量" default(8)
// @Success 200 {object} map[string]interface{}
// @Router /products/top-selling [get]
func GetTopSelling(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))
	if limit < 1 {
		limit = 8
	}

	type ProductWithSales struct {
		models.Product
		SalesCount int64 `json:"salesCount"`
	}

	var results []ProductWithSales
	database.DB.Model(&models.OrderItem{}).
		Select("product_id, SUM(quantity) as sales_count").
		Group("product_id").
		Order("sales_count DESC").
		Limit(limit).
		Scan(&results)

	if len(results) == 0 {
		c.JSON(http.StatusOK, gin.H{"products": []interface{}{}})
		return
	}

	var productIDs []uint
	for _, r := range results {
		productIDs = append(productIDs, r.ID)
	}

	var products []models.Product
	database.DB.Where("id IN ? AND is_active = ?", productIDs, true).Preload("Category").Find(&products)

	productMap := make(map[uint]models.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	type TopProduct struct {
		models.Product
		SalesCount int64 `json:"salesCount"`
	}

	var topProducts []TopProduct
	for _, r := range results {
		if p, ok := productMap[r.ID]; ok {
			topProducts = append(topProducts, TopProduct{
				Product:    p,
				SalesCount: r.SalesCount,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"products": topProducts})
}

// GetProductSuggestions godoc
// @Summary 搜索建议
// @Description 根据关键词返回匹配的商品名称建议列表
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "搜索关键词"
// @Param limit query int false "返回数量" default(6)
// @Success 200 {object} map[string]interface{}
// @Router /products/suggest [get]
func GetProductSuggestions(c *gin.Context) {
	q := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))
	if limit < 1 {
		limit = 6
	}
	if limit > 10 {
		limit = 10
	}

	if q == "" || len(q) < 1 {
		c.JSON(http.StatusOK, gin.H{"suggestions": []interface{}{}})
		return
	}

	type SuggestionItem struct {
		ID           uint    `json:"id"`
		Name         string  `json:"name"`
		Price        float64 `json:"price"`
		ThumbnailURL string  `json:"thumbnailUrl"`
	}

	var suggestions []SuggestionItem
	database.DB.Model(&models.Product{}).
		Select("id, name, price, thumbnail_image").
		Where("is_active = ? AND name LIKE ?", true, "%"+q+"%").
		Order("stock > 0 DESC, created_at DESC").
		Limit(limit).
		Scan(&suggestions)

	if suggestions == nil {
		suggestions = []SuggestionItem{}
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
