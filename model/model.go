package model

import "time"

type Model struct {
	ID uint32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}