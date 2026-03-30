package models

import "time"

const (
	PaymentProofStatusPending  = "pending"
	PaymentProofStatusApproved = "approved"
	PaymentProofStatusRejected = "rejected"
)

type PaymentProof struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	OrderID         uint          `json:"orderId" gorm:"not null;index"`
	UserID          uint          `json:"userId" gorm:"not null;index"`
	PaymentMethodID uint          `json:"paymentMethodId" gorm:"not null"`
	PaymentMethod   PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodID"`
	ImageUrl        string        `json:"imageUrl" gorm:"size:500;not null"`
	Status          string        `json:"status" gorm:"size:20;default:'pending'"`
	ReviewerID      *uint         `json:"reviewerId"`
	ReviewedAt      *time.Time    `json:"reviewedAt"`
	RejectReason    string        `json:"rejectReason" gorm:"size:500"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type ReviewPaymentProofInput struct {
	Action       string `json:"action" binding:"required"`
	RejectReason string `json:"rejectReason"`
}
