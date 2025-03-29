package config

import "time"

var (
	// JWT 过期时间 7 天
	JWTExpire = time.Hour * 24 * 7
	JWTSecret = "your_jwt_secret_key"
)
