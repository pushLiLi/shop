package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/internal/ws"
	imagepkg "bycigar-server/pkg/image"
	"bycigar-server/pkg/storage"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	conversation = models.Conversation{
		UserID:        uid,
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
		Content:        "您好！欢迎来到 HUAUHE，有什么可以帮助您的吗？",
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

	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	message := models.Message{
		ConversationID: uint(conversationID),
		SenderType:     "customer",
		SenderID:       uid,
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
	ws.DefaultHub.SendToUser(uid, WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(uid, WSResponse{
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

	err = storage.SaveFile(origName, result.Original)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传图片失败"})
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

func CustomerCloseConversation(c *gin.Context) {
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

	convID := uint(conversationID)

	uid, ok := userID.(uint)
	if !ok {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	systemMsg := models.Message{
		ConversationID: convID,
		SenderType:     "system",
		SenderID:       0,
		MessageType:    "text",
		Content:        "客户已结束对话",
	}
	database.DB.Create(&systemMsg)

	database.DB.Model(&conversation).Update("status", "closed")

	ws.DefaultHub.SendToUser(uid, WSResponse{
		Type:           "new_message",
		Message:        systemMsg,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(uid, WSResponse{
		Type:           "conversation_closed",
		ConversationID: convID,
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "new_message",
		Message:        systemMsg,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(convID),
	})

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetServiceStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"online": ws.DefaultHub.IsServiceOnline(),
	})
}

func RateConversation(c *gin.Context) {
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

	if conversation.Status != "closed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "只能评价已关闭的对话")
		return
	}

	var existing models.Rating
	if err := database.DB.Where("conversation_id = ?", conversationID).First(&existing).Error; err == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "该对话已评价")
		return
	}

	var input struct {
		Score   int    `json:"score" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请提供评分")
		return
	}
	if input.Score < 1 || input.Score > 5 {
		utils.ErrorResponse(c, http.StatusBadRequest, "评分范围为1-5")
		return
	}

	rating := models.Rating{
		ConversationID: uint(conversationID),
		Score:          input.Score,
		Comment:        input.Comment,
	}
	if err := database.DB.Create(&rating).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "评价失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"rating": rating})
}
