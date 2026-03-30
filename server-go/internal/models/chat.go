package models

import (
	"time"
)

type Conversation struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"userId" gorm:"index;not null"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	Status        string    `json:"status" gorm:"size:20;default:'open'"`
	AssignedTo    *uint     `json:"assignedTo" gorm:"index"`
	AssignedUser  *User     `json:"assignedUser" gorm:"foreignKey:AssignedTo"`
	LastMessageAt time.Time `json:"lastMessageAt" gorm:"index"`
	Rating        *Rating   `json:"rating" gorm:"foreignKey:ConversationID"`
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
	Status         string       `json:"status" gorm:"size:20;default:'sent'"`
	RecalledAt     *time.Time   `json:"recalledAt" gorm:"index"`
	IsRead         bool         `json:"isRead" gorm:"default:false"`
	CreatedAt      time.Time    `json:"createdAt"`
}

type SendMessageInput struct {
	Content      string `json:"content" binding:"required"`
	MessageType  string `json:"messageType"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type QuickReply struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:100;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedBy uint      `json:"createdBy" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:CreatedBy"`
	SortOrder int       `json:"sortOrder" gorm:"default:0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Rating struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	ConversationID uint         `json:"conversationId" gorm:"uniqueIndex;not null"`
	Conversation   Conversation `json:"conversation" gorm:"foreignKey:ConversationID"`
	Score          int          `json:"score" gorm:"not null"`
	Comment        string       `json:"comment" gorm:"type:text"`
	CreatedAt      time.Time    `json:"createdAt"`
}
