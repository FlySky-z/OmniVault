package handlers

import (
	"net/http"
	"omnivault/config"
	"omnivault/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				ErrorCode: http.StatusNotFound,
				ErrorMsg:  "User not found",
				Details:   "User with the given ID does not exist",
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

	c.JSON(http.StatusOK, user)
}

// UpdateUser 用户修改信息
// @Summary 修改用户信息
// @Description 修改当前用户的信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param userInfo body models.User true "用户信息"
// @Success 200 {object} successResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put] [operationId updateUserInfo]
func UpdateUser(c *gin.Context) {
	var userInfo models.User
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Invalid input",
			Details:   err.Error(),
		})
		return
	}

	if err := config.DB.Model(&models.User{}).Where("id = ?", userInfo.ID).Updates(userInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Failed to update user info",
			Details:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Message: "User info updated successfully",
	})
}
