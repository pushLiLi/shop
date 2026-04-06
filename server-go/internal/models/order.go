package models

import (
	"time"
)

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

var ValidOrderStatusTransitions = map[string][]string{
	OrderStatusPending:    {OrderStatusProcessing, OrderStatusCancelled},
	OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
	OrderStatusShipped:    {OrderStatusCompleted},
	OrderStatusCompleted:  {},
	OrderStatusCancelled:  {},
}

type Order struct {
	ID              uint        `json:"id" gorm:"primaryKey"`
	OrderNo         string      `json:"orderNo" gorm:"uniqueIndex;size:64;not null"`
	UserID          uint        `json:"userId" gorm:"not null;index"`
	AddressID       uint        `json:"addressId" gorm:"not null"`
	Address         Address     `json:"address" gorm:"foreignKey:AddressID"`
	Total           float64     `json:"total" gorm:"not null"`
	Status          string      `json:"status" gorm:"default:'pending'"`
	Remark          string      `json:"remark"`
	TrackingCompany string      `json:"trackingCompany"`
	TrackingNumber  string      `json:"trackingNumber"`
	Items           []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"orderId" gorm:"not null;index"`
	ProductID uint    `json:"productId" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Currency  string  `json:"currency" gorm:"size:10;default:'CNY'"`
}

type OrderResponse struct {
	Orders []Order `json:"orders"`
}

type CreateOrderInput struct {
	AddressID uint   `json:"addressId" binding:"required"`
	Remark    string `json:"remark"`
}
