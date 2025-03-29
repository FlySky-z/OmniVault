package config

import (
	"fmt"
	"log"
	"omnivault/models"
	"strings"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func InitDB() {
	if GlobalConfig == nil {
		log.Fatalf("GlobalConfig is not initialized")
	}

	var err error
	switch GlobalConfig.Database.Type {
	case "mysql":
		mysqlConfig := GlobalConfig.Database.MySQL
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port, mysqlConfig.Database)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			// 如果数据库不存在，尝试创建数据库
			if strings.Contains(err.Error(), fmt.Sprintf("Unknown database '%s'", mysqlConfig.Database)) {
				createDB(mysqlConfig)
				DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			} else {
				log.Fatalf("failed to connect to database: %v", err)
			}
		}
		autoMigrate()
	case "sqlite":
		sqliteConfig := GlobalConfig.Database.SQLite
		dsn := sqliteConfig.File
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("unsupported database type: %s", GlobalConfig.Database.Type)
	}

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Database connected successfully")
}

func createDB(mysqlConfig MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConfig.Username, mysqlConfig.Password, mysqlConfig.Host, mysqlConfig.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to MySQL server: %v", err)
	}
	defer func(db *gorm.DB) {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to close database connection: %v", err)
		}
		sqlDB.Close()
	}(db)

	// 创建数据库
	// CREATE DATABASE IF NOT EXISTS `test` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
	// 创建数据库 mysqlConfig.Database，字符集为 utf8mb4，校对规则为 utf8mb4_general_ci
	// utf8mb4：指定字符集。
	// general：一种常见的排序规则，适用于大部分语言。
	// ci：表示“Case Insensitive”，即不区分大小写的排序规则。比如，“a”和“A”会被视为相等。
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", mysqlConfig.Database)
	if err := db.Exec(createDBSQL).Error; err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	fmt.Printf("Database %s created successfully\n", mysqlConfig.Database)
}

// autoMigrate初始化数据库并自动迁移
func autoMigrate() {
	err := DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.UserRole{},
		&models.RolePermission{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	log.Println("Database auto-migration completed successfully")
}

func InitCache() {
	if GlobalConfig == nil {
		log.Fatalf("GlobalConfig is not initialized")
	}

	switch GlobalConfig.Cache.Type {
	case "redis":
		redisConfig := GlobalConfig.Cache.Redis
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.Database,
		})

		_, err := RedisClient.Ping().Result()
		if err != nil {
			log.Fatalf("failed to connect to Redis: %v", err)
		}
		fmt.Println("Redis connected successfully")
	case "memory":
		log.Fatal("Memory cache is not supported yet")
	default:
		log.Fatalf("unsupported cache type: %s", GlobalConfig.Cache.Type)
	}
}
