package models

import "time"

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"index;not null"`
	Qty       int  `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Product Product `gorm:"foreignKey:ProductID"`
}
