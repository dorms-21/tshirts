package models

import "time"

type UserProfile struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex;not null"`

	Gender   string `gorm:"not null"` // "hombre" | "mujer"
	Location string `gorm:"not null;default:''"`

	Style    string `gorm:"not null;default:''"`
	Colors   string `gorm:"not null;default:''"`
	Fit      string `gorm:"not null;default:''"`
	Occasion string `gorm:"not null;default:''"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
