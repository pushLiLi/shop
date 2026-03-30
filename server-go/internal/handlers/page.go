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

	var pages []models.Page
	result := database.DB.Where("slug = ?", slug).Limit(1).Find(&pages)
	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "页面不存在")
		return
	}

	c.JSON(http.StatusOK, pages[0])
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

	var pages []models.Page
	result := database.DB.Where("slug = ?", slug).Limit(1).Find(&pages)
	if result.RowsAffected == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "页面不存在")
		return
	}
	page := pages[0]

	page.Title = input.Title
	page.Content = input.Content

	if err := database.DB.Save(&page).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "保存失败")
		return
	}

	c.JSON(http.StatusOK, page)
}
