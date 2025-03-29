package utils

import (
	"fmt"
	"omnivault/config"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT 生成 JWT token
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(config.JWTExpire)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 生成 JWT token, 使用 HS256 算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 将 token 添加到黑名单
func AddTokenToBlacklist(token string) error {
	// 去除 token 中的 Bearer 前缀
	token = strings.TrimPrefix(token, "Bearer ")
	// 使用 Redis 的 Set 方法存储 token，并设置过期时间
	expirationTime := time.Duration(config.JWTExpire.Seconds()) * time.Second
	err := config.RedisClient.Set(token, "blacklisted", expirationTime).Err()
	if err != nil {
		return fmt.Errorf("failed to add token to blacklist: %w", err)
	}
	return nil
}
