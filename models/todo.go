package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID              uint           `json:"id"`
	Title           string         `json:"title" validate:"required"`
	ActivityGroupId uint64         `json:"activity_group_id" validate:"required"`
	IsActive        bool           `json:"is_active"`
	Priority        string         `json:"priority"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
