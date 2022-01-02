package models

import (
	"time"
)

type Activity struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Email     *string   `gorm:"varchar(50)" json:"email"`
	Title     *string   `gorm:"varchar(50)" json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}
