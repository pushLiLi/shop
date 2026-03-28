package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

// GetConfig godoc
// @Summary 获取网站配置
// @Description 获取所有网站配置项
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /config [get]
func GetConfig(c *gin.Context) {
	var configs []models.SiteConfig
	database.DB.Find(&configs)

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}

	c.JSON(http.StatusOK, result)
}

// UpdateConfig godoc
// @Summary 更新网站配置
// @Description 更新指定的网站配置项
// @Tags admin-config
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param key path string true "配置键名"
// @Param input body map[string]string true "配置值"
// @Success 200 {object} models.SiteConfig
// @Failure 400 {object} map[string]interface{}
// @Router /admin/config/{key} [put]
func UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var input struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.SiteConfig
	result := database.DB.Where("config_key = ?", key).First(&config)
	if result.Error != nil {
		config = models.SiteConfig{
			ConfigKey:   key,
			ConfigValue: input.Value,
		}
		database.DB.Create(&config)
	} else {
		config.ConfigValue = input.Value
		database.DB.Save(&config)
	}

	c.JSON(http.StatusOK, config)
}
