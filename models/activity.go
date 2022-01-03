package models

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID        uint           `json:"id"`
	Email     *string        `gorm:"varchar(50)" json:"email"`
	Title     *string        `gorm:"varchar(50)" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
