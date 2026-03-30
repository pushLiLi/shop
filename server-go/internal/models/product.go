package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"not null"`
	Slug           string         `json:"slug" gorm:"type:varchar(191);uniqueIndex"`
	Description    string         `json:"description"`
	Price          float64        `json:"price" gorm:"not null"`
	Image          string         `json:"imageUrl"`
	Images         string         `json:"images"`
	ThumbnailImage string         `json:"thumbnailUrl" gorm:"column:thumbnail_image"`
	CategoryID     uint           `json:"categoryId" gorm:"index:idx_category_active"`
	Category       Category       `json:"category" gorm:"foreignKey:CategoryID"`
	Stock          int            `json:"stock" gorm:"default:0"`
	IsActive       bool           `json:"isActive" gorm:"default:true;index:idx_category_active;index:idx_active_featured"`
	IsFeatured     bool           `json:"isFeatured" gorm:"default:false;index:idx_active_featured"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type ProductListResponse struct {
	Products   []Product `json:"products"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}
