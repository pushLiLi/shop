package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetCart godoc
// @Summary 获取购物车
// @Description 获取当前用户的购物车内容
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /cart [get]
func GetCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}, "total": 0})
		return
	}

	var items []models.CartItem
	database.DB.Where("user_id = ?", userID).Preload("Product.Category").Find(&items)

	var total float64
	for _, item := range items {
		if item.Product.ID > 0 {
			total += item.Product.Price * float64(item.Quantity)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"total": total,
	})
}

// AddToCart godoc
// @Summary 添加商品到购物车
// @Description 添加商品到当前用户的购物车
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.AddToCartInput true "商品信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /cart [post]
func AddToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.AddToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	quantity := input.Quantity
	if quantity < 1 {
		quantity = 1
	}

	var existing models.CartItem
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing)

	var item models.CartItem
	if result.Error == nil {
		existing.Quantity += quantity
		database.DB.Save(&existing)
		item = existing
	} else {
		item = models.CartItem{
			UserID:    userID.(uint),
			ProductID: input.ProductID,
			Quantity:  quantity,
		}
		database.DB.Create(&item)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "item": item})
}

// UpdateCartItem godoc
// @Summary 更新购物车商品数量
// @Description 更新购物车中指定商品的数量
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车项ID"
// @Param input body models.UpdateCartInput true "数量信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var input models.UpdateCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Item not found")
		return
	}

	if input.Quantity <= 0 {
		database.DB.Delete(&item)
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	item.Quantity = input.Quantity
	database.DB.Save(&item)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteCartItem godoc
// @Summary 删除购物车商品
// @Description 从购物车中删除指定商品
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车项ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /cart/{id} [delete]
func DeleteCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var item models.CartItem
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Item not found")
		return
	}

	database.DB.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
