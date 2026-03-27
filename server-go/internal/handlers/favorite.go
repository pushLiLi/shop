package handlers

import (
	"net/http"
	"strconv"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"items": []interface{}{}})
		return
	}

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).Preload("Product.Category").Find(&favorites)

	c.JSON(http.StatusOK, gin.H{"items": favorites})
}

func AddFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.AddFavoriteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var existing models.Favorite
	result := database.DB.Where("user_id = ? AND product_id = ?", userID, input.ProductID).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "exists": true})
		return
	}

	favorite := models.Favorite{
		UserID:    userID.(uint),
		ProductID: input.ProductID,
	}
	database.DB.Create(&favorite)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func DeleteFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID")
		return
	}

	database.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Favorite{})

	c.JSON(http.StatusOK, gin.H{"success": true})
}
