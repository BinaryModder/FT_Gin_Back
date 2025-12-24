package model

import "gorm.io/gorm"

type User struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Age   int    `json:"age"`
}

func (User) TableName() string {
	return "users"
}
