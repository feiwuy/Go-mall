// Package config storages the project configs
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	DataBase struct {
		Host         string
		Port         string
		User         string
		Password     string
		Name         string
		MaxIdleConns int
		MaxOpenConns int
	}
}

var AppConfig *Config

// InitConfig 获取并解析配置文件
func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
