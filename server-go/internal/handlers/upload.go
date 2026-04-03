package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	imagepkg "bycigar-server/pkg/image"
	"bycigar-server/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadImage godoc
// @Summary 上传图片
// @Description 上传图片文件到MinIO，自动生成缩略图
// @Tags admin-upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "图片文件"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /admin/upload [post]
func UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 jpg, jpeg, png, gif, webp 格式的图片"})
		return
	}
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 10MB"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	result, _ := imagepkg.Process(fileBytes, ext)

	baseName := fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	origName := baseName + result.OrigExt

	err = storage.SaveFile(origName, result.Original)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败"})
		return
	}

	url := storage.URLPrefix + origName
	thumbnailUrl := ""

	if len(result.Thumbnail) > 0 {
		thumbName := baseName + "_thumb.jpg"
		err = storage.SaveFile(thumbName, result.Thumbnail)
		if err == nil {
			thumbnailUrl = storage.URLPrefix + thumbName
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"url":          url,
		"thumbnailUrl": thumbnailUrl,
		"success":      true,
	})
}
