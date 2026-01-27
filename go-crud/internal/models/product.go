package models

import "time"

type Product struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	PriceCents  int64  `gorm:"not null"`
	ImagePath   string `gorm:"type:text"` // /static/uploads/x.jpg
	Stock       int    `gorm:"not null;default:0"`
	Active      bool   `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
