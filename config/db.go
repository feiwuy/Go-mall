package config

import (
	"fmt"
	"gomall/global"
	"gomall/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", AppConfig.DataBase.User, AppConfig.DataBase.Password, AppConfig.DataBase.Host, AppConfig.DataBase.Port, AppConfig.DataBase.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initial database, error: %v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to configure database, error: %v", err)
	}

	sqlDB.SetMaxIdleConns(AppConfig.DataBase.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.DataBase.MaxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Hour)

	// 数据库交给全局变量
	global.DB = db

	// 自动迁移, 若数据库无表格->创建
	if err := global.DB.AutoMigrate(&models.User{}, &models.TodoList{}); err != nil {
		log.Fatalf("Failed to create database, error: %v", err)
		return
	}
}
