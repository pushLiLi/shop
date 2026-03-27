package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

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
