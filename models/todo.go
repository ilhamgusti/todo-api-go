package models

type Todo struct {
	ID              string  `json:"id"`
	Title           *string `gorm:"varchar(30)" json:"title"`
	ActivityGroupId *uint   `json:"activity_group_id"`
	IsActive        bool    `json:"is_active"`
	Priority        string  `gorm:"varchar(10)" json:"priority"`
}
