package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetProducts godoc
// @Summary 获取产品列表
// @Description 获取产品列表，支持分页、搜索、分类筛选和排序
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(12)
// @Param search query string false "搜索关键词"
// @Param category query string false "分类 slug"
// @Param categoryId query int false "分类 ID"
// @Param featured query string false "是否精选"
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
			query = query.Where("category_id = ?", category.ID)
		}
	}

	if categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			query = query.Where("category_id = ?", categoryID)
		}
	}

	if featured == "true" {
		query = query.Where("is_featured = ?", true)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	orderClause := sortColumn + " " + sortOrder
	if sortOrder != "asc" {
		orderClause = sortColumn + " DESC"
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
	var categories []models.Category
	database.DB.Where("parent_id IS NULL").Preload("Children").Find(&categories)

	type CategoryWithCount struct {
		models.Category
		ProductCount int64 `json:"_count"`
	}

	var result []CategoryWithCount
	for _, cat := range categories {
		var count int64
		database.DB.Model(&models.Product{}).Where("category_id = ?", cat.ID).Count(&count)
		result = append(result, CategoryWithCount{
			Category:     cat,
			ProductCount: count,
		})
	}

	c.JSON(http.StatusOK, result)
}
