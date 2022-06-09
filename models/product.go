package models

import "time"

type Product struct {
	ID           uint `json:"id" gorm:"primaryKey"`
	CreatedAt    time.Time
	Name         string `json:"name" validate:"required,min=3,max=32"`
	SerialNumber string `json:"serial_number" validate:"required,min=3,max=32"`
}
