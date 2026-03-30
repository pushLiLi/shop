package models

import "time"

type PaymentMethod struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:100;not null"`
	QRCodeUrl    string    `json:"qrCodeUrl" gorm:"size:500"`
	PaymentUrl   string    `json:"paymentUrl" gorm:"size:500"`
	Instructions string    `json:"instructions" gorm:"size:500"`
	IsActive     bool      `json:"isActive" gorm:"default:true"`
	SortOrder    int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreatePaymentMethodInput struct {
	Name         string `json:"name" binding:"required"`
	QRCodeUrl    string `json:"qrCodeUrl"`
	PaymentUrl   string `json:"paymentUrl"`
	Instructions string `json:"instructions"`
	IsActive     *bool  `json:"isActive"`
	SortOrder    int    `json:"sortOrder"`
}

type UpdatePaymentMethodInput struct {
	Name         string `json:"name"`
	QRCodeUrl    string `json:"qrCodeUrl"`
	PaymentUrl   string `json:"paymentUrl"`
	Instructions string `json:"instructions"`
	IsActive     *bool  `json:"isActive"`
	SortOrder    int    `json:"sortOrder"`
}
