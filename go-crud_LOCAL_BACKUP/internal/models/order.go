package models

import "time"

type OrderStatus string

const (
	OrderPending OrderStatus = "pending"
	OrderPaid    OrderStatus = "paid"
	OrderFailed  OrderStatus = "failed"
)

type Order struct {
	ID         uint        `gorm:"primaryKey"`
	UserID     uint        `gorm:"index;not null"`
	Status     OrderStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	TotalCents int64       `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Items []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID         uint   `gorm:"primaryKey"`
	OrderID    uint   `gorm:"index;not null"`
	ProductID  uint   `gorm:"index;not null"`
	Qty        int    `gorm:"not null"`
	PriceCents int64  `gorm:"not null"` // snapshot del precio
	NameSnap   string `gorm:"not null"` // snapshot del nombre
	ImageSnap  string `gorm:"type:text"`
}
