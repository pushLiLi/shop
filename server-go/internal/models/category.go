package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Slug      string         `json:"slug" gorm:"type:varchar(191);uniqueIndex;not null"`
	ParentID  *uint          `json:"parentId"`
	Parent    *Category      `json:"parent" gorm:"foreignKey:ParentID"`
	Children  []Category     `json:"children" gorm:"foreignKey:ParentID"`
	Products  []Product      `json:"products" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
