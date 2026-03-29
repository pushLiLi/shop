package models

import (
	"time"
)

const (
	NotificationTypeOrderStatus = "order_status"
	NotificationTypeBackInStock = "back_in_stock"
	NotificationTypePriceDrop   = "price_drop"
)

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null;index"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Title     string    `json:"title" gorm:"size:255;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	IsRead    bool      `json:"isRead" gorm:"default:false"`
	Link      string    `json:"link" gorm:"size:500"`
	ProductID *uint     `json:"productId"`
	OrderID   *uint     `json:"orderId"`
	CreatedAt time.Time `json:"createdAt"`
}
