package handlers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"omnivault/models"

	"github.com/gin-gonic/gin"
)

type UploadFileBody struct {
	File       *multipart.FileHeader `form:"file" binding:"required"`
	UploadPath string                `json:"uploadpath"`
}

type UploadFileResponse struct {
	Message    string `json:"message"`
	UploadPath string `json:"uploadpath"`
}

// UploadFile godoc
// @Summary Upload a file to the server (simulating object storage)
// @Description Uploads a file to the server using the PUT method and simulates object storage
// @Tags file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param Authorization header string true "Bearer token"
// @Param Content-MD5 header string false "Base64 encoded MD5 hash of the file"
// @Success 200 {object} UploadFileResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /upload [put]
func UploadFile(c *gin.Context) {
	// 从请求中获取文件
	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			ErrorMsg:  "Unable to retrieve file",
			Details:   err.Error(),
		})
		return
	}
	defer file.Close()

	// 生成唯一的文件名
	ext := filepath.Ext(handler.Filename)
	uniqueFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// 模拟对象存储路径（例如：存储在 `uploads` 文件夹下）
	uploadDir := "uploads"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Unable to create upload directory",
			Details:   err.Error(),
		})
		return
	}

	// 保存文件到本地目录，使用唯一文件名
	filePath := filepath.Join(uploadDir, uniqueFileName)
	tempFile, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Unable to save file",
			Details:   err.Error(),
		})
		return
	}
	defer tempFile.Close()

	// 将文件内容拷贝到临时文件
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Unable to reset file pointer",
			Details:   err.Error(),
		})
		return
	}
	if _, err := io.Copy(tempFile, file); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			ErrorMsg:  "Unable to copy file",
			Details:   err.Error(),
		})
		log.Printf("failed to copy file to uploads")
		return
	}

	// 构造返回的模拟对象存储 URL（实际是本地路径）
	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", uniqueFileName)

	// 返回上传成功的响应
	c.JSON(http.StatusOK, UploadFileResponse{
		Message:    "File uploaded successfully",
		UploadPath: fileURL, // 返回文件的模拟存储路径
	})
}
