package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Type   string       `mapstructure:"type"`
		SQLite SQLiteConfig `mapstructure:"sqlite"`
		MySQL  MySQLConfig  `mapstructure:"mysql"`
	} `mapstructure:"database"`
	Cache struct {
		Type  string      `mapstructure:"type"`
		Redis RedisConfig `mapstructure:"redis"`
	} `mapstructure:"cache"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     int    `mapstructure:"port"`
}

type SQLiteConfig struct {
	File string `mapstructure:"file"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Database int    `mapstructure:"database"`
}

// GlobalConfig 全局配置变量
var GlobalConfig *Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	GlobalConfig = &config

	// 初始化数据库和 Redis
	InitDB()
	InitCache()

	return nil
}
