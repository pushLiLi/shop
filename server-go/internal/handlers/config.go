package handlers

import (
	"net/http"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetConfig(c *gin.Context) {
	var configs []models.SiteConfig
	database.DB.Find(&configs)

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}

	c.JSON(http.StatusOK, result)
}

func UpdateConfig(c *gin.Context) {
	key := c.Param("key")

	var input struct {
		Value string `json:"value" binding:"required"`
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
