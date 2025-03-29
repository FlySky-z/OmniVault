package middleware

import (
	"fmt"
	"net/http"
	"omnivault/config"
	"omnivault/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode: http.StatusUnauthorized,
				ErrorMsg:  "Authorization header is missing",
				Details:   "Please provide a valid token in the Authorization header",
			})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 检查token是否在黑名单
		if config.RedisClient.SIsMember("blacklist", tokenString).Val() {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode: http.StatusUnauthorized,
				ErrorMsg:  "Token is blacklisted",
				Details:   "This token has been revoked and cannot be used",
			})
			c.Abort()
			return
		}
		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证 token 的签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 返回用于签名的密钥
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode: http.StatusUnauthorized,
				ErrorMsg:  "Invalid token",
				Details:   err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
