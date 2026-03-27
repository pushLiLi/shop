package models

import (
	"time"
)

type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	ProductID uint      `json:"productId" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity" gorm:"default:1"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CartResponse struct {
	Items []CartItem `json:"items"`
	Total float64    `json:"total"`
}

type AddToCartInput struct {
	ProductID uint `json:"productId" binding:"required"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartInput struct {
	Quantity int `json:"quantity" binding:"required"`
}
