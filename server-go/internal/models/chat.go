package models

import (
	"time"
)

type Conversation struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"userId" gorm:"index;not null"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	Status        string    `json:"status" gorm:"size:20;default:'open'"`
	LastMessageAt time.Time `json:"lastMessageAt" gorm:"index"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Message struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	ConversationID uint         `json:"conversationId" gorm:"index;not null"`
	Conversation   Conversation `json:"-" gorm:"foreignKey:ConversationID"`
	SenderType     string       `json:"senderType" gorm:"size:20;not null"`
	SenderID       uint         `json:"senderId" gorm:"not null"`
	MessageType    string       `json:"messageType" gorm:"size:20;default:text"`
	Content        string       `json:"content" gorm:"type:text;not null"`
	ThumbnailURL   string       `json:"thumbnailUrl" gorm:"size:500"`
	IsRead         bool         `json:"isRead" gorm:"default:false"`
	CreatedAt      time.Time    `json:"createdAt"`
}

type SendMessageInput struct {
	Content      string `json:"content" binding:"required"`
	MessageType  string `json:"messageType"`
	ThumbnailURL string `json:"thumbnailUrl"`
}
