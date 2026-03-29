package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"type:varchar(191);uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name"`
	Role      string         `json:"role" gorm:"default:'customer'"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type RegisterInput struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	Name        string `json:"name"`
	CaptchaId   string `json:"captchaId" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

type LoginInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaId   string `json:"captchaId"`
	CaptchaCode string `json:"captchaCode"`
}

type UpdateProfileInput struct {
	Name string `json:"name"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
	CaptchaId   string `json:"captchaId" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}
