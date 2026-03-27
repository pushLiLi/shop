package models

import (
	"time"
)

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null;uniqueIndex:idx_user_product"`
	ProductID uint      `json:"productId" gorm:"not null;uniqueIndex:idx_user_product"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	CreatedAt time.Time `json:"createdAt"`
}

type FavoriteResponse struct {
	Items []Favorite `json:"items"`
}

type AddFavoriteInput struct {
	ProductID uint `json:"productId" binding:"required"`
}
