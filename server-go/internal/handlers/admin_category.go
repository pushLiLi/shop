package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CategoryInput struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID *uint  `json:"parentId"`
}

// CreateCategory godoc
// @Summary 创建分类
// @Description 创建新的商品分类
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body CategoryInput true "分类信息"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/categories [post]
func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Slug == "" {
		input.Slug = generateSlug(input.Name)
	}

	category := models.Category{
		Name:     input.Name,
		Slug:     input.Slug,
		ParentID: input.ParentID,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	InvalidateCategoriesCache()
	c.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc
// @Summary 更新分类
// @Description 更新指定分类信息
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param input body CategoryInput true "分类信息"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name
	if input.Slug != "" {
		category.Slug = input.Slug
	}
	category.ParentID = input.ParentID

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	InvalidateCategoriesCache()
	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary 删除分类
// @Description 删除指定分类
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/categories/{id} [delete]
func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var productCount int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在商品，无法删除"})
		return
	}

	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下存在子分类，无法删除"})
		return
	}

	result := database.DB.Delete(&models.Category{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	InvalidateCategoriesCache()
	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}

// GetAdminCategories godoc
// @Summary 获取分类列表(管理端)
// @Description 获取所有分类列表
// @Tags admin-categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Category
// @Router /admin/categories [get]
func GetAdminCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("id ASC").Find(&categories)

	c.JSON(http.StatusOK, categories)
}
