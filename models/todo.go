package models

import (
	"gorm.io/gorm"
)

type TodoList struct {
	gorm.Model
	Title     string
	Content   string
	Completed bool
	UserID    uint `gorm:"index"` // 外键. 落在todo_lists上
}
