package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := database.DB.Where("slug = ?", slug).First(&page).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "页面不存在")
		return
	}

	c.JSON(http.StatusOK, page)
}

func GetAdminPages(c *gin.Context) {
	var pages []models.Page
	database.DB.Order("id asc").Find(&pages)

	c.JSON(http.StatusOK, pages)
}

func UpdatePage(c *gin.Context) {
	slug := c.Param("slug")

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的请求数据")
		return
	}

	var page models.Page
	if err := database.DB.Where("slug = ?", slug).First(&page).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "页面不存在")
		return
	}

	page.Title = input.Title
	page.Content = input.Content

	if err := database.DB.Save(&page).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "保存失败")
		return
	}

	c.JSON(http.StatusOK, page)
}
