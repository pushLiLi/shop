package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetBanners(c *gin.Context) {
	var banners []models.Banner
	database.DB.Where("is_active = ?", true).Order("sort_order asc, id desc").Find(&banners)

	c.JSON(http.StatusOK, banners)
}

func CreateBanner(c *gin.Context) {
	var input models.Banner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&input)
	c.JSON(http.StatusCreated, input)
}

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

func DeleteBanner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid banner ID"})
		return
	}

	database.DB.Delete(&models.Banner{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Banner deleted"})
}
