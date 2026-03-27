package models

import (
	"time"
)

type SiteConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ConfigKey   string    `json:"configKey" gorm:"type:varchar(191);uniqueIndex;not null"`
	ConfigValue string    `json:"configValue"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
