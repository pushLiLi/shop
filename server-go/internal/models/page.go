package models

import (
	"time"
)

type Page struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Slug      string    `json:"slug" gorm:"uniqueIndex;not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	UpdatedAt time.Time `json:"updatedAt"`
}
