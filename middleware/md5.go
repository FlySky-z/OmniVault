package middleware

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"omnivault/models"

	"github.com/gin-gonic/gin"
)

func MD5Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentMD5 := c.GetHeader("Content-MD5")
		if contentMD5 == "" {
			c.Next()
			return
		}

		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				ErrorCode: http.StatusBadRequest,
				ErrorMsg:  "Unable to retrieve file",
				Details:   err.Error(),
			})
			c.Abort()
			return
		}
		defer file.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode: http.StatusInternalServerError,
				ErrorMsg:  "Unable to compute MD5 hash",
				Details:   err.Error(),
			})
			c.Abort()
			return
		}
		fileMD5 := hex.EncodeToString(hash.Sum(nil))

		if contentMD5 != fileMD5 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				ErrorCode: http.StatusBadRequest,
				ErrorMsg:  "MD5 checksum mismatch",
				Details:   "Expected: " + contentMD5 + ", Got: " + fileMD5,
			})
			c.Abort()
			return
		}

		// 重置文件指针
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				ErrorCode: http.StatusInternalServerError,
				ErrorMsg:  "Unable to reset file pointer",
				Details:   err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
