package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	Type           string `json:"type"`
	ConversationID uint   `json:"conversationId,omitempty"`
	Content        string `json:"content,omitempty"`
}

type WSResponse struct {
	Type           string      `json:"type"`
	Message        interface{} `json:"message,omitempty"`
	Conversation   interface{} `json:"conversation,omitempty"`
	ConversationID uint        `json:"conversationId,omitempty"`
	Count          int         `json:"count,omitempty"`
	Stats          interface{} `json:"stats,omitempty"`
	TotalUnread    int64       `json:"totalUnread,omitempty"`
}

func HandleCustomerWS(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	uid := userID.(uint)

	client := &ws.Client{
		UserID: uid,
		Role:   "customer",
		Conn:   conn,
		Hub:    ws.DefaultHub,
		Send:   make(chan []byte, 64),
	}

	ws.DefaultHub.Register <- client

	go client.WritePump()

	sendUnreadCountToCustomer(uid)

	client.ReadPump(handleCustomerMessage)
}

func HandleAdminWS(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uid := userID.(uint)
	var user models.User
	database.DB.First(&user, uid)
	if user.Role != "admin" && user.Role != "service" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		UserID: uid,
		Role:   user.Role,
		Conn:   conn,
		Hub:    ws.DefaultHub,
		Send:   make(chan []byte, 64),
	}

	ws.DefaultHub.Register <- client

	go client.WritePump()

	sendUnreadStatsToAdmins()

	client.ReadPump(handleAdminMessage)
}

func handleCustomerMessage(client *ws.Client, raw []byte) {
	var msg WSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case "send_message":
		handleCustomerSendMessage(client, msg)
	case "mark_read":
		handleCustomerMarkRead(client, msg)
	}
}

func handleAdminMessage(client *ws.Client, raw []byte) {
	var msg WSMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return
	}

	switch msg.Type {
	case "send_message":
		handleAdminSendMessage(client, msg)
	case "mark_read":
		handleAdminMarkRead(client, msg)
	}
}

func handleCustomerSendMessage(client *ws.Client, msg WSMessage) {
	if msg.Content == "" || len(msg.Content) > 500 {
		return
	}

	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", msg.ConversationID, client.UserID).First(&conversation).Error; err != nil {
		return
	}
	if conversation.Status == "closed" {
		return
	}

	message := models.Message{
		ConversationID: msg.ConversationID,
		SenderType:     "customer",
		SenderID:       client.UserID,
		Content:        msg.Content,
	}
	if err := database.DB.Create(&message).Error; err != nil {
		return
	}

	now := time.Now()
	database.DB.Model(&conversation).Update("last_message_at", now)

	database.DB.Where("id = ?", msg.ConversationID).First(&conversation)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: msg.ConversationID,
	})

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleAdminSendMessage(client *ws.Client, msg WSMessage) {
	if msg.Content == "" || len(msg.Content) > 500 {
		return
	}

	var conversation models.Conversation
	if err := database.DB.First(&conversation, msg.ConversationID).Error; err != nil {
		return
	}
	if conversation.Status == "closed" {
		return
	}

	message := models.Message{
		ConversationID: msg.ConversationID,
		SenderType:     "service",
		SenderID:       client.UserID,
		Content:        msg.Content,
	}
	if err := database.DB.Create(&message).Error; err != nil {
		return
	}

	now := time.Now()
	database.DB.Model(&conversation).Update("last_message_at", now)

	database.DB.Where("id = ?", msg.ConversationID).First(&conversation)

	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:    "new_message",
		Message: message,
	})

	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})

	sendUnreadCountToCustomer(conversation.UserID)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: msg.ConversationID,
	})

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleCustomerMarkRead(client *ws.Client, msg WSMessage) {
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", msg.ConversationID, "service", false).
		Update("is_read", true)

	sendUnreadCountToCustomer(client.UserID)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleAdminMarkRead(client *ws.Client, msg WSMessage) {
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", msg.ConversationID, "customer", false).
		Update("is_read", true)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func sendUnreadCountToCustomer(userID uint) {
	var count int64
	database.DB.Model(&models.Message{}).
		Joins("JOIN conversations ON conversations.id = messages.conversation_id").
		Where("conversations.user_id = ? AND messages.sender_type = ? AND messages.is_read = ?", userID, "service", false).
		Count(&count)

	ws.DefaultHub.SendToUser(userID, WSResponse{
		Type:  "unread_count",
		Count: int(count),
	})
}

func sendUnreadStatsToAdmins() {
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

	var totalUnread int64
	for _, s := range stats {
		totalUnread += s.UnreadCount
	}

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:        "unread_stats",
		Stats:       stats,
		TotalUnread: totalUnread,
	})
}

func buildConversationDetail(convID uint) ConversationWithDetails {
	var conv models.Conversation
	if err := database.DB.Preload("User").First(&conv, convID).Error; err != nil {
		return ConversationWithDetails{}
	}

	var unreadCount int64
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", convID, "customer", false).
		Count(&unreadCount)

	var lastMessage models.Message
	hasLast := database.DB.Where("conversation_id = ?", convID).
		Order("created_at desc").
		First(&lastMessage).Error == nil

	item := ConversationWithDetails{
		Conversation: conv,
		UnreadCount:  unreadCount,
	}
	if hasLast {
		item.LastMessage = &lastMessage
	}
	return item
}
