// Package global storage the global context likes database
package global

import "gorm.io/gorm"

var (
	DB *gorm.DB
)
