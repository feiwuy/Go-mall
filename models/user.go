// Package models 存储了各种服务所需的model
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;size:64;not null"`
	Password string	`gorm:"size:60;not null"` // bcrypt固定60
}
