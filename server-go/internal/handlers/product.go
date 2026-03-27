package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	categorySlug := c.Query("category")
	categoryIDStr := c.Query("categoryId")
	featured := c.Query("featured")
	sortBy := c.DefaultQuery("sortBy", c.DefaultQuery("sort", "createdAt"))
	sortOrder := c.DefaultQuery("sortOrder", c.DefaultQuery("order", "desc"))

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

	orderClause := sortBy + " " + sortOrder
	if sortOrder != "asc" {
		orderClause = sortBy + " DESC"
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

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Preload("Children").Find(&categories)

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
