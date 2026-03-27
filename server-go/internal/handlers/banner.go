package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

// GetBanners godoc
// @Summary 获取轮播图列表
// @Description 获取所有激活状态的轮播图
// @Tags banners
// @Accept json
// @Produce json
// @Success 200 {array} models.Banner
// @Router /banners [get]
func GetBanners(c *gin.Context) {
	var banners []models.Banner
	database.DB.Where("is_active = ?", true).Order("sort_order asc, id desc").Find(&banners)

	c.JSON(http.StatusOK, banners)
}

// CreateBanner godoc
// @Summary 创建轮播图
// @Description 创建新的轮播图
// @Tags admin-banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.Banner true "轮播图信息"
// @Success 201 {object} models.Banner
// @Failure 400 {object} map[string]interface{}
// @Router /admin/banners [post]
func CreateBanner(c *gin.Context) {
	var input models.Banner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

// UpdateBanner godoc
// @Summary 更新轮播图
// @Description 更新指定的轮播图
// @Tags admin-banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "轮播图ID"
// @Param input body models.Banner true "轮播图信息"
// @Success 200 {object} models.Banner
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/banners/{id} [put]
func UpdateBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	var banner models.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	var input models.Banner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banner.Title = input.Title
	banner.Image = input.Image
	banner.Link = input.Link
	banner.SortOrder = input.SortOrder
	banner.IsActive = input.IsActive

	database.DB.Save(&banner)
	c.JSON(http.StatusOK, banner)
}

// DeleteBanner godoc
// @Summary 删除轮播图
// @Description 删除指定的轮播图
// @Tags admin-banners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "轮播图ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /admin/banners/{id} [delete]
func DeleteBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}
	database.DB.Delete(&models.Banner{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Banner deleted"})
}
