package models

import "gorm.io/gorm"

type Setting struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Key       string         `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string         `gorm:"type:text" json:"value"`
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
