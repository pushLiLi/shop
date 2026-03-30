package models

import "time"

type ContactMethod struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"size:50;not null"`
	Label     string    `json:"label" gorm:"size:100;not null"`
	Value     string    `json:"value" gorm:"size:500;not null"`
	QRCodeUrl string    `json:"qrCodeUrl" gorm:"size:500"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	SortOrder int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateContactMethodInput struct {
	Type      string `json:"type" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Value     string `json:"value" binding:"required"`
	QRCodeUrl string `json:"qrCodeUrl"`
	IsActive  *bool  `json:"isActive"`
	SortOrder int    `json:"sortOrder"`
}

type UpdateContactMethodInput struct {
	Type      string `json:"type"`
	Label     string `json:"label"`
	Value     string `json:"value"`
	QRCodeUrl string `json:"qrCodeUrl"`
	IsActive  *bool  `json:"isActive"`
	SortOrder int    `json:"sortOrder"`
}
