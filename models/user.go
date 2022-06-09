package models

import "time"

type User struct {
	ID        uint `json:"id" gorm:"primaryKey" `
	CreatedAt time.Time
	FirstName string `json:"first_name" validate:"required,min=3,max=32"`
	LastName  string `json:"last_name" validate:"required,min=3,max=32"`
}
