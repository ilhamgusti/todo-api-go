package models

type Activity struct {
	ID    string  `json:"id"`
	Email *string `gorm:"varchar(45)" json:"email"`
	Title *string `gorm:"varchar(30)" json:"title"`
}
