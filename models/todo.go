package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID              int            `gorm:"int" json:"id"`
	Title           *string        `gorm:"varchar(50)" json:"title"`
	ActivityGroupId *uint          `gorm:"index" json:"activity_group_id"`
	IsActive        bool           `json:"is_active"`
	Priority        string         `gorm:"varchar(10)" json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
