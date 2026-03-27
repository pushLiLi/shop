package models

import (
	"time"

	"gorm.io/gorm"
)

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`
	Image     string         `json:"imageUrl" gorm:"not null"`
	Link      string         `json:"link"`
	SortOrder int            `json:"sortOrder" gorm:"default:0"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
