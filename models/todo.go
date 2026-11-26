package models

import (
	"gorm.io/gorm"
)

type TodoList struct {
	gorm.Model
	Title     string	`gorm:"not null"`
	Content   string	`gorm:"not null"`
	Completed bool		`gorm:"not null"`
	UserID    uint		`gorm:"index;not null"` // 外键. 落在todo_lists上
}
