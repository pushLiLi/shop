package models

import "time"

type OrderSummary struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Date            string    `json:"date" gorm:"type:varchar(10);uniqueIndex;not null"`
	TotalOrders     int       `json:"totalOrders" gorm:"default:0"`
	TotalRevenue    float64   `json:"totalRevenue" gorm:"default:0"`
	TotalItems      int       `json:"totalItems" gorm:"default:0"`
	PendingOrders   int       `json:"pendingOrders" gorm:"default:0"`
	CompletedOrders int       `json:"completedOrders" gorm:"default:0"`
	CancelledOrders int       `json:"cancelledOrders" gorm:"default:0"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
