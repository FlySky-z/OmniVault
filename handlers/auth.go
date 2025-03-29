package handlers

import (
	"net/http"
	"omnivault/config"
	"omnivault/models"
	"omnivault/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type successResponse struct {
	Message string `json:"message"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Code    int    `json:"code"`
}

// RegisterHandler 用户注册
// @Summary 用户注册
// @Description 注册新用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} successResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /authorize/register [post]
func RegisterHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Invalid input",
			Details:   err.Error(),
		})
		return
	}
	// 对密码进行哈希处理（使用 bcrypt 算法）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Failed to hash password",
			Details:   err.Error(),
		})
		return
	}
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Failed to create user",
			Details:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Message: "Registration successful",
	})
}

// LoginHandler 用户登录
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param loginRequest body loginRequest true "登录信息"
// @Success 200 {object} loginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /authorize/login [post]
func LoginHandler(c *gin.Context) {
	var loginData loginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Invalid input",
			Details:   err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				ErrorCode: http.StatusUnauthorized,
				ErrorMsg:  "Invalid username or password",
				Details:   "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode: http.StatusInternalServerError,
				ErrorMsg:  "Failed to query user",
				Details:   err.Error(),
			})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			ErrorCode: http.StatusUnauthorized,
			ErrorMsg:  "Invalid username or password",
			Details:   "Password mismatch",
		})
		return
	}

	// 生成 JWT 令牌
	jwtToken, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Failed to generate token",
			Details:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Code:    http.StatusOK,
		Message: "Login successful",
		Token:   jwtToken,
	})
}

// LogoutHandler 用户注销
// @Summary 用户注销
// @Description 注销当前用户
// @Tags 用户
// @Produce json
// @Success 200 {object} successResponse
// @Router /authorize/logout [post]
func LogoutHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Missing token",
			Details:   "Authorization header is required",
		})
		return
	}
	// 将 token 添加到黑名单
	if err := utils.AddTokenToBlacklist(token); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Failed to blacklist token",
			Details:   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, successResponse{
		Message: "Logout successful",
	})
}
