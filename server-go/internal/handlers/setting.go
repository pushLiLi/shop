package handlers

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/email"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

var (
	settingCache     map[string]string
	settingCacheTime time.Time
	settingCacheMu   sync.RWMutex
	settingCacheTTL  = 5 * time.Minute
)

func getSettingsCached() (map[string]string, error) {
	settingCacheMu.RLock()
	if settingCache != nil && time.Since(settingCacheTime) < settingCacheTTL {
		result := settingCache
		settingCacheMu.RUnlock()
		return result, nil
	}
	settingCacheMu.RUnlock()

	settingCacheMu.Lock()
	defer settingCacheMu.Unlock()

	if settingCache != nil && time.Since(settingCacheTime) < settingCacheTTL {
		return settingCache, nil
	}

	var settings []models.Setting
	if err := database.DB.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	settingCache = result
	settingCacheTime = time.Now()
	return result, nil
}

type UpdateSettingRequest struct {
	Value string `json:"value"`
}

// GetSettings godoc
// @Summary 获取所有设置
// @Description 获取所有站点设置（公开）
// @Tags settings
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /settings [get]
func GetSettings(c *gin.Context) {
	result, err := getSettingsCached()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取设置失败")
		return
	}

	if result["email_smtp_password"] != "" {
		result["email_smtp_password"] = "****"
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// UpdateSetting godoc
// @Summary 更新设置
// @Description 更新指定 key 的设置值（管理员）
// @Tags admin
// @Accept json
// @Produce json
// @Param key path string true "设置键名"
// @Param body body UpdateSettingRequest true "设置值"
// @Success 200 {object} map[string]interface{}
// @Router /admin/settings/{key} [put]
func UpdateSetting(c *gin.Context) {
	key := c.Param("key")
	var req UpdateSettingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	var setting models.Setting
	result := database.DB.Where("`key` = ?", key).First(&setting)

	if result.Error != nil {
		setting = models.Setting{
			Key:   key,
			Value: req.Value,
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "创建设置失败")
			return
		}
	} else {
		setting.Value = req.Value
		if err := database.DB.Save(&setting).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "更新设置失败")
			return
		}
	}

	settingCacheMu.Lock()
	settingCache = nil
	settingCacheMu.Unlock()

	if strings.HasPrefix(key, "email_") {
		email.InvalidateEmailCache()
	}

	utils.SuccessResponse(c, setting)
}
