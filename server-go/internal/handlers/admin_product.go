package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Image       string  `json:"imageUrl"`
	Images      string  `json:"images"`
	CategoryID  uint    `json:"categoryId"`
	Stock       int     `json:"stock"`
	IsActive    bool    `json:"isActive"`
	IsFeatured  bool    `json:"isFeatured"`
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	return slug
}

func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Slug == "" {
		input.Slug = generateSlug(input.Name)
	}

	product := models.Product{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Price:       input.Price,
		Image:       input.Image,
		Images:      input.Images,
		CategoryID:  input.CategoryID,
		Stock:       input.Stock,
		IsActive:    input.IsActive,
		IsFeatured:  input.IsFeatured,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found")
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Name = input.Name
	if input.Slug != "" {
		product.Slug = input.Slug
	}
	product.Description = input.Description
	product.Price = input.Price
	product.Image = input.Image
	product.Images = input.Images
	product.CategoryID = input.CategoryID
	product.Stock = input.Stock
	product.IsActive = input.IsActive
	product.IsFeatured = input.IsFeatured

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	result := database.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除商品失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "商品已删除"})
}

func GetAdminProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	categoryIDStr := c.Query("categoryId")
	featured := c.Query("featured")
	active := c.Query("active")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	query := database.DB.Model(&models.Product{})

	if categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			query = query.Where("category_id = ?", categoryID)
		}
	}

	if featured == "true" {
		query = query.Where("is_featured = ?", true)
	} else if featured == "false" {
		query = query.Where("is_featured = ?", false)
	}

	if active == "true" {
		query = query.Where("is_active = ?", true)
	} else if active == "false" {
		query = query.Where("is_active = ?", false)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	offset := (page - 1) * limit
	query.Preload("Category").Order("id DESC").Offset(offset).Limit(limit).Find(&products)

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
