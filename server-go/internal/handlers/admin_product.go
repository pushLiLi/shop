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
	Name           string  `json:"name"`
	Slug           string  `json:"slug"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	Image          string  `json:"imageUrl"`
	Images         string  `json:"images"`
	ThumbnailImage string  `json:"thumbnailUrl"`
	CategoryID     uint    `json:"categoryId"`
	Stock          int     `json:"stock"`
	IsActive       bool    `json:"isActive"`
	IsFeatured     bool    `json:"isFeatured"`
}

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ".", "")
	return slug
}

// CreateProduct godoc
// @Summary 创建商品
// @Description 创建新的商品
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body ProductInput true "商品信息"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/products [post]
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
		Name:           input.Name,
		Slug:           input.Slug,
		Description:    input.Description,
		Price:          input.Price,
		Image:          input.Image,
		Images:         input.Images,
		ThumbnailImage: input.ThumbnailImage,
		CategoryID:     input.CategoryID,
		Stock:          input.Stock,
		IsActive:       input.IsActive,
		IsFeatured:     input.IsFeatured,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary 更新商品
// @Description 更新指定商品信息
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Param input body ProductInput true "商品信息"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/products/{id} [put]
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

	oldStock := product.Stock
	oldPrice := product.Price

	product.Name = input.Name
	if input.Slug != "" {
		product.Slug = input.Slug
	}
	product.Description = input.Description
	product.Price = input.Price
	product.Image = input.Image
	product.Images = input.Images
	product.ThumbnailImage = input.ThumbnailImage
	product.CategoryID = input.CategoryID
	product.Stock = input.Stock
	product.IsActive = input.IsActive
	product.IsFeatured = input.IsFeatured

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新商品失败"})
		return
	}

	if oldStock <= 0 && input.Stock > 0 {
		var favorites []models.Favorite
		database.DB.Where("product_id = ?", uint(id)).Find(&favorites)
		notifications := make([]models.Notification, 0, len(favorites))
		for _, fav := range favorites {
			notifications = append(notifications, models.Notification{
				UserID:    fav.UserID,
				Type:      models.NotificationTypeBackInStock,
				Title:     "商品到货通知",
				Content:   "您收藏的商品「" + product.Name + "」已到货",
				Link:      "/products/" + strconv.Itoa(int(product.ID)),
				ProductID: &product.ID,
			})
		}
		if len(notifications) > 0 {
			database.DB.Create(&notifications)
		}
	}

	if input.Price < oldPrice {
		var favorites []models.Favorite
		database.DB.Where("product_id = ?", uint(id)).Find(&favorites)
		notifications := make([]models.Notification, 0, len(favorites))
		for _, fav := range favorites {
			notifications = append(notifications, models.Notification{
				UserID:    fav.UserID,
				Type:      models.NotificationTypePriceDrop,
				Title:     "商品降价通知",
				Content:   "您收藏的商品「" + product.Name + "」价格已下调",
				Link:      "/products/" + strconv.Itoa(int(product.ID)),
				ProductID: &product.ID,
			})
		}
		if len(notifications) > 0 {
			database.DB.Create(&notifications)
		}
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary 删除商品
// @Description 删除指定商品
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/products/{id} [delete]
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

// GetAdminProducts godoc
// @Summary 获取商品列表(管理端)
// @Description 获取商品列表，支持分页和筛选
// @Tags admin-products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(20)
// @Param search query string false "搜索关键词"
// @Param categoryId query int false "分类ID"
// @Param featured query string false "是否推荐"
// @Param active query string false "是否上架"
// @Success 200 {object} map[string]interface{}
// @Router /admin/products [get]
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

type BatchStatusInput struct {
	IDs      []uint `json:"ids" binding:"required"`
	IsActive bool   `json:"isActive"`
}

func BatchUpdateProductStatus(c *gin.Context) {
	var input BatchStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择商品"})
		return
	}

	result := database.DB.Model(&models.Product{}).Where("id IN ?", input.IDs).Update("is_active", input.IsActive)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量更新成功", "updated": result.RowsAffected})
}

type BatchDeleteInput struct {
	IDs []uint `json:"ids" binding:"required"`
}

func BatchDeleteProducts(c *gin.Context) {
	var input BatchDeleteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择商品"})
		return
	}

	result := database.DB.Where("id IN ?", input.IDs).Delete(&models.Product{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量删除成功", "deleted": result.RowsAffected})
}
