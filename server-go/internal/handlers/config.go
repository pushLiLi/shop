package handlers

import (
	"net/http"
	"sync"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
)

var (
	configCache     map[string]string
	configCacheTime time.Time
	configCacheMu   sync.RWMutex
	configCacheTTL  = 5 * time.Minute
)

func getConfigCached() map[string]string {
	configCacheMu.RLock()
	if configCache != nil && time.Since(configCacheTime) < configCacheTTL {
		result := configCache
		configCacheMu.RUnlock()
		return result
	}
	configCacheMu.RUnlock()

	configCacheMu.Lock()
	defer configCacheMu.Unlock()

	if configCache != nil && time.Since(configCacheTime) < configCacheTTL {
		return configCache
	}

	var configs []models.SiteConfig
	database.DB.Find(&configs)

	result := make(map[string]string)
	for _, config := range configs {
		result[config.ConfigKey] = config.ConfigValue
	}
	configCache = result
	configCacheTime = time.Now()
	return result
}

// GetConfig godoc
// @Summary 获取网站配置
// @Description 获取所有网站配置项
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /config [get]
func GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, getConfigCached())
}

type SiteIdentity struct {
	Title           string `json:"title"`
	MetaDescription string `json:"metaDescription"`
	FaviconUrl      string `json:"faviconUrl"`
}

// GetSiteIdentity godoc
// @Summary 获取网站标识
// @Description 获取网站标题、META描述和图标URL
// @Tags site
// @Produce json
// @Success 200 {object} SiteIdentity
// @Router /site-identity [get]
func GetSiteIdentity(c *gin.Context) {
	config := getConfigCached()
	identity := SiteIdentity{
		Title:           config["site_title"],
		MetaDescription: config["site_meta_description"],
		FaviconUrl:      config["favicon_url"],
	}
	if identity.Title == "" {
		identity.Title = "BYCIGAR"
	}
	if identity.FaviconUrl == "" {
		identity.FaviconUrl = "/favicon.png"
	}
	c.JSON(http.StatusOK, identity)
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

	configCacheMu.Lock()
	configCache = nil
	configCacheMu.Unlock()

	c.JSON(http.StatusOK, config)
}
