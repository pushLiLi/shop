package handlers

import (
	"net/http"
	"strconv"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"
	"bycigar-server/internal/ws"

	"github.com/gin-gonic/gin"
)

type ConversationWithDetails struct {
	models.Conversation
	UnreadCount int64           `json:"unreadCount"`
	LastMessage *models.Message `json:"lastMessage"`
}

func GetAdminConversations(c *gin.Context) {
	status := c.Query("status")

	query := database.DB.Model(&models.Conversation{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var conversations []models.Conversation
	query.Preload("User").
		Order("last_message_at desc").
		Find(&conversations)

	result := make([]ConversationWithDetails, 0, len(conversations))
	for _, conv := range conversations {
		var unreadCount int64
		database.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_type = ? AND is_read = ?", conv.ID, "customer", false).
			Count(&unreadCount)

		var lastMessages []models.Message
		database.DB.Where("conversation_id = ?", conv.ID).
			Order("created_at desc").
			Limit(1).
			Find(&lastMessages)

		item := ConversationWithDetails{
			Conversation: conv,
			UnreadCount:  unreadCount,
		}
		if len(lastMessages) > 0 {
			item.LastMessage = &lastMessages[0]
		}
		result = append(result, item)
	}

	c.JSON(http.StatusOK, gin.H{"conversations": result})
}

func GetAdminMessages(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
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
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", conversationID, "customer", false).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func AdminSendMessage(c *gin.Context) {
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
	if err := database.DB.First(&conversation, conversationID).Error; err != nil {
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
		SenderType:     "service",
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
	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(convID),
	})
	sendUnreadCountToCustomer(conversation.UserID)
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

func CloseConversation(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
		return
	}

	var conversation models.Conversation
	if err := database.DB.First(&conversation, conversationID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "对话不存在")
		return
	}

	if conversation.Status == "closed" {
		utils.ErrorResponse(c, http.StatusBadRequest, "对话已关闭")
		return
	}

	convID := uint(conversationID)

	systemMsg := models.Message{
		ConversationID: convID,
		SenderType:     "system",
		SenderID:       0,
		MessageType:    "text",
		Content:        "客服已结束对话",
	}
	database.DB.Create(&systemMsg)

	database.DB.Model(&conversation).Update("status", "closed")

	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:           "new_message",
		Message:        systemMsg,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
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

func GetAdminUnreadStats(c *gin.Context) {
	type UnreadStat struct {
		ConversationID uint  `json:"conversationId"`
		UnreadCount    int64 `json:"unreadCount"`
	}

	var stats []UnreadStat
	database.DB.Model(&models.Message{}).
		Select("conversation_id, count(*) as unread_count").
		Where("sender_type = ? AND is_read = ?", "customer", false).
		Group("conversation_id").
		Scan(&stats)

	totalUnread := int64(0)
	for _, s := range stats {
		totalUnread += s.UnreadCount
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":       stats,
		"totalUnread": totalUnread,
	})
}

func SetServiceStatus(c *gin.Context) {
	var input struct {
		Online bool `json:"online"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)

	if input.Online {
		wasEmpty := ws.DefaultHub.SetServiceOnline(uid)
		if wasEmpty {
			ws.DefaultHub.SendToAllCustomers(WSResponse{
				Type:          "service_status",
				ServiceOnline: true,
			})
		}
	} else {
		nowEmpty := ws.DefaultHub.SetServiceOffline(uid)
		if nowEmpty {
			ws.DefaultHub.SendToAllCustomers(WSResponse{
				Type:          "service_status",
				ServiceOnline: false,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "online": input.Online})
}
