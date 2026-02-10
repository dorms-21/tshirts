package models

import "time"

type PaymentStatus string

const (
	PayApproved PaymentStatus = "approved"
	PayDeclined PaymentStatus = "declined"
)

type Payment struct {
	ID        uint          `gorm:"primaryKey"`
	OrderID   uint          `gorm:"uniqueIndex;not null"`
	Status    PaymentStatus `gorm:"type:varchar(20);not null"`
	Method    string        `gorm:"type:varchar(50);not null;default:'simulated'"`
	Ref       string        `gorm:"type:varchar(100)"`
	CreatedAt time.Time
}
