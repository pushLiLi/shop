package models

import (
	"time"
)

type Address struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"userId" gorm:"not null;index"`
	FullName     string    `json:"fullName" gorm:"not null"`
	AddressLine1 string    `json:"addressLine1" gorm:"not null"`
	AddressLine2 string    `json:"addressLine2"`
	City         string    `json:"city" gorm:"not null"`
	State        string    `json:"state" gorm:"not null"`
	ZipCode      string    `json:"zipCode" gorm:"not null"`
	Phone        string    `json:"phone" gorm:"not null"`
	IsDefault    bool      `json:"isDefault" gorm:"default:false"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AddressResponse struct {
	Addresses []Address `json:"addresses"`
}

type CreateAddressInput struct {
	FullName     string `json:"fullName" binding:"required"`
	AddressLine1 string `json:"addressLine1" binding:"required"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city" binding:"required"`
	State        string `json:"state" binding:"required"`
	ZipCode      string `json:"zipCode" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	IsDefault    bool   `json:"isDefault"`
}

type UpdateAddressInput struct {
	FullName     string `json:"fullName"`
	AddressLine1 string `json:"addressLine1"`
	AddressLine2 string `json:"addressLine2"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
	Phone        string `json:"phone"`
}
