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
	MessageType    string `json:"messageType,omitempty"`
	ThumbnailURL   string `json:"thumbnailUrl,omitempty"`
	MessageID      uint   `json:"messageId,omitempty"`
}

type WSResponse struct {
	Type           string      `json:"type"`
	Message        interface{} `json:"message,omitempty"`
	Conversation   interface{} `json:"conversation,omitempty"`
	ConversationID uint        `json:"conversationId,omitempty"`
	MessageID      uint        `json:"messageId,omitempty"`
	Status         string      `json:"status,omitempty"`
	Count          int         `json:"count,omitempty"`
	Stats          interface{} `json:"stats,omitempty"`
	TotalUnread    int64       `json:"totalUnread,omitempty"`
	ServiceOnline  bool        `json:"serviceOnline,omitempty"`
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

	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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
	userVal, _ := c.Get("user")
	user := userVal.(models.User)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		UserID: user.ID,
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
	case "typing":
		handleCustomerTyping(client, msg)
	case "close_conversation":
		handleCustomerCloseConvWS(client, msg)
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
	case "typing":
		handleAdminTyping(client, msg)
	case "recall_message":
		handleAdminRecallMessage(client, msg)
	case "service_online":
		wasEmpty := ws.DefaultHub.SetServiceOnline(client.UserID)
		if wasEmpty {
			ws.DefaultHub.SendToAllCustomers(WSResponse{
				Type:          "service_status",
				ServiceOnline: true,
			})
		}
	case "service_offline":
		nowEmpty := ws.DefaultHub.SetServiceOffline(client.UserID)
		if nowEmpty {
			ws.DefaultHub.SendToAllCustomers(WSResponse{
				Type:          "service_status",
				ServiceOnline: false,
			})
		}
	}
}

func handleCustomerSendMessage(client *ws.Client, msg WSMessage) {
	if msg.Content == "" {
		return
	}
	msgType := msg.MessageType
	if msgType == "" {
		msgType = "text"
	}
	if msgType == "text" && len(msg.Content) > 500 {
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
		MessageType:    msgType,
		Content:        msg.Content,
		ThumbnailURL:   msg.ThumbnailURL,
		Status:         "sent",
	}
	if err := database.DB.Create(&message).Error; err != nil {
		return
	}

	now := time.Now()
	database.DB.Model(&conversation).Update("last_message_at", now)

	database.DB.Where("id = ?", msg.ConversationID).Limit(1).Find(&conversation)

	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: msg.ConversationID,
	})

	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:           "ack",
		ConversationID: msg.ConversationID,
		MessageID:      message.ID,
		Status:         "sent",
	})

	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})

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
	if msg.Content == "" {
		return
	}
	msgType := msg.MessageType
	if msgType == "" {
		msgType = "text"
	}
	if msgType == "text" && len(msg.Content) > 500 {
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
		MessageType:    msgType,
		Content:        msg.Content,
		ThumbnailURL:   msg.ThumbnailURL,
		Status:         "sent",
	}
	if err := database.DB.Create(&message).Error; err != nil {
		return
	}

	now := time.Now()
	database.DB.Model(&conversation).Update("last_message_at", now)

	database.DB.Where("id = ?", msg.ConversationID).Limit(1).Find(&conversation)

	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:           "new_message",
		Message:        message,
		ConversationID: msg.ConversationID,
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
		Type:           "ack",
		ConversationID: msg.ConversationID,
		MessageID:      message.ID,
		Status:         "sent",
	})

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleCustomerMarkRead(client *ws.Client, msg WSMessage) {
	var msgs []models.Message
	database.DB.Where("conversation_id = ? AND sender_type = ? AND is_read = ?", msg.ConversationID, "service", false).Find(&msgs)

	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", msg.ConversationID, "service", false).
		Updates(map[string]interface{}{"is_read": true, "status": "read"})

	sendUnreadCountToCustomer(client.UserID)

	for _, m := range msgs {
		ws.DefaultHub.SendToAdmins(WSResponse{
			Type:           "message_status",
			ConversationID: msg.ConversationID,
			MessageID:      m.ID,
			Status:         "read",
		})
	}

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleCustomerTyping(client *ws.Client, msg WSMessage) {
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "typing",
		ConversationID: msg.ConversationID,
	})
}

func handleCustomerCloseConvWS(client *ws.Client, msg WSMessage) {
	var conversation models.Conversation
	if err := database.DB.Where("id = ? AND user_id = ?", msg.ConversationID, client.UserID).First(&conversation).Error; err != nil {
		return
	}
	if conversation.Status == "closed" {
		return
	}

	convID := msg.ConversationID

	systemMsg := models.Message{
		ConversationID: convID,
		SenderType:     "system",
		SenderID:       0,
		MessageType:    "text",
		Content:        "客户已结束对话",
	}
	database.DB.Create(&systemMsg)

	database.DB.Model(&conversation).Update("status", "closed")

	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:           "new_message",
		Message:        systemMsg,
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:           "conversation_closed",
		ConversationID: convID,
	})
	ws.DefaultHub.SendToUser(client.UserID, WSResponse{
		Type:           "rating_request",
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
}

func handleAdminTyping(client *ws.Client, msg WSMessage) {
	var conversation models.Conversation
	if err := database.DB.First(&conversation, msg.ConversationID).Error; err != nil {
		return
	}
	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:           "typing",
		ConversationID: msg.ConversationID,
	})
}

func handleAdminMarkRead(client *ws.Client, msg WSMessage) {
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", msg.ConversationID, "customer", false).
		Updates(map[string]interface{}{"is_read": true, "status": "read"})

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(msg.ConversationID),
	})
}

func handleAdminRecallMessage(client *ws.Client, msg WSMessage) {
	var message models.Message
	if err := database.DB.Where("id = ? AND conversation_id = ? AND sender_type = ? AND sender_id = ? AND recalled_at IS NULL",
		msg.MessageID, msg.ConversationID, "service", client.UserID).First(&message).Error; err != nil {
		return
	}

	now := time.Now()
	database.DB.Model(&message).Update("recalled_at", now)
	database.DB.Where("id = ?", message.ID).First(&message)

	var conversation models.Conversation
	database.DB.First(&conversation, msg.ConversationID)

	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
		Type:           "message_recalled",
		ConversationID: msg.ConversationID,
		Message:        message,
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "message_recalled",
		ConversationID: msg.ConversationID,
		Message:        message,
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
	if err := database.DB.Preload("User").Preload("AssignedUser").Preload("Rating").First(&conv, convID).Error; err != nil {
		return ConversationWithDetails{}
	}

	var unreadCount int64
	database.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_type = ? AND is_read = ?", convID, "customer", false).
		Count(&unreadCount)

	var lastMessages []models.Message
	database.DB.Where("conversation_id = ?", convID).
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
	return item
}
