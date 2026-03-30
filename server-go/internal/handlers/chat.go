package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	imagepkg "bycigar-server/pkg/image"
	miniopkg "bycigar-server/pkg/minio"
	"bycigar-server/pkg/utils"
	"bycigar-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func CreateConversation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var conversation models.Conversation
	result := database.DB.Where("user_id = ? AND status = ?", userID, "open").First(&conversation)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"conversation": conversation})
		return
	}

	conversation = models.Conversation{
		UserID:        userID.(uint),
		Status:        "open",
		LastMessageAt: time.Now(),
	}
	if err := database.DB.Create(&conversation).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建对话失败")
		return
	}

	greeting := models.Message{
		ConversationID: conversation.ID,
		SenderType:     "service",
		SenderID:       0,
		MessageType:    "text",
		Content:        "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？",
	}
	database.DB.Create(&greeting)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(conversation.ID),
	})

	c.JSON(http.StatusOK, gin.H{"conversation": conversation})
}

func GetConversations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var conversations []models.Conversation
	database.DB.Where("user_id = ?", userID).
		Order("last_message_at desc").
		Find(&conversations)

	c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

func GetMessages(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
		return
	}

	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "对话不存在")
		return
	}

	afterID, _ := strconv.Atoi(c.Query("after"))

	var messages []models.Message
	query := database.DB.Where("conversation_id = ?", conversationID)
	if afterID > 0 {
		query = query.Where("id > ?", afterID)
	} else {
		query = query.Order("created_at desc").Limit(50)
	}
	query.Find(&messages)

	if afterID == 0 && len(messages) > 0 {
		reversed := make([]models.Message, len(messages))
		for i, m := range messages {
			reversed[len(messages)-1-i] = m
		}
		messages = reversed
	}

	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", conversationID, "service", false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func SendMessage(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
		return
	}

	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "对话不存在")
		return
	}

	if conversation.Status == "closed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "对话已关闭")
		return
	}

	var input models.SendMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "消息内容不能为空")
		return
	}

	msgType := input.MessageType
	if msgType == "" {
		msgType = "text"
	}
	if msgType == "text" && len(input.Content) > 500 {
		utils.ErrorResponse(c, http.StatusBadRequest, "消息内容不能超过500字")
		return
	}

	message := models.Message{
		ConversationID: uint(conversationID),
		SenderType:     "customer",
		SenderID:       userID.(uint),
		MessageType:    msgType,
		Content:        input.Content,
		ThumbnailURL:   input.ThumbnailURL,
	}
	if err := database.DB.Create(&message).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "发送失败")
		return
	}

	now := time.Now()
	database.DB.Model(&conversation).Update("last_message_at", now)

	convID := uint(conversationID)
	ws.DefaultHub.SendToUser(userID.(uint), WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(userID.(uint), WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(convID),
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(convID),
	})

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func GetChatUnreadCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}

	var count int64
	database.DB.Model(&models.Message{}).
		Joins("JOIN conversations ON conversations.id = messages.conversation_id").
		Where("conversations.user_id = ? AND messages.sender_type = ? AND messages.is_read = ?", userID, "service", false).
		Count(&count)

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func UploadChatImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 jpg, png, webp 格式的图片"})
		return
	}
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 5MB"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	result, _ := imagepkg.Process(fileBytes, ext)

	baseName := fmt.Sprintf("chat_%d_%s", time.Now().Unix(), uuid.New().String())
	origName := baseName + result.OrigExt

	_, err = miniopkg.Client.PutObject(
		context.Background(),
		miniopkg.Bucket,
		origName,
		bytes.NewReader(result.Original),
		int64(len(result.Original)),
		minio.PutObjectOptions{ContentType: result.ContentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传图片失败"})
		return
	}

	url := fmt.Sprintf("/media/%s/%s", miniopkg.Bucket, origName)
	thumbnailUrl := ""

	if len(result.Thumbnail) > 0 {
		thumbName := baseName + "_thumb.jpg"
		_, err = miniopkg.Client.PutObject(
			context.Background(),
			miniopkg.Bucket,
			thumbName,
			bytes.NewReader(result.Thumbnail),
			int64(len(result.Thumbnail)),
			minio.PutObjectOptions{ContentType: "image/jpeg"},
		)
		if err == nil {
			thumbnailUrl = fmt.Sprintf("/media/%s/%s", miniopkg.Bucket, thumbName)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"url":          url,
		"thumbnailUrl": thumbnailUrl,
		"success":      true,
	})
}
