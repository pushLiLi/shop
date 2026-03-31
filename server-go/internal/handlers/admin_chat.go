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
	assignedTo := c.Query("assignedTo")

	query := database.DB.Model(&models.Conversation{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if assignedTo == "me" {
		userID, _ := c.Get("userID")
		query = query.Where("assigned_to = ?", userID)
	} else if assignedTo == "unassigned" {
		query = query.Where("assigned_to IS NULL")
	}

	var conversations []models.Conversation
	query.Preload("User").Preload("AssignedUser").
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
	ws.DefaultHub.SendToUser(conversation.UserID, WSResponse{
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

func RecallMessage(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
		return
	}
	messageID, err := strconv.Atoi(c.Param("msgId"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的消息ID")
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)

	var message models.Message
	if err := database.DB.Where("id = ? AND conversation_id = ? AND sender_type = ? AND sender_id = ? AND recalled_at IS NULL",
		messageID, conversationID, "service", uid).First(&message).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "消息不存在或无法撤回")
		return
	}

	now := time.Now()
	database.DB.Model(&message).Update("recalled_at", now)
	database.DB.Where("id = ?", message.ID).First(&message)

	ws.DefaultHub.SendToUser(message.ConversationID, WSResponse{
		Type:           "message_recalled",
		ConversationID: uint(conversationID),
		Message:        message,
	})
	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:           "message_recalled",
		ConversationID: uint(conversationID),
		Message:        message,
	})

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func AssignConversation(c *gin.Context) {
	conversationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的对话ID")
		return
	}

	var input struct {
		AssignedTo *uint `json:"assignedTo"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	var conversation models.Conversation
	if err := database.DB.First(&conversation, conversationID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "对话不存在")
		return
	}

	updates := map[string]interface{}{"assigned_to": input.AssignedTo}
	if conversation.Status == "open" && input.AssignedTo != nil {
		updates["status"] = "active"
	}
	database.DB.Model(&conversation).Updates(updates)

	ws.DefaultHub.SendToAdmins(WSResponse{
		Type:         "conversation_updated",
		Conversation: buildConversationDetail(uint(conversationID)),
	})

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetQuickReplies(c *gin.Context) {
	var replies []models.QuickReply
	database.DB.Preload("User").Order("sort_order asc, created_at desc").Find(&replies)
	c.JSON(http.StatusOK, gin.H{"quickReplies": replies})
}

func CreateQuickReply(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		SortOrder int    `json:"sortOrder"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "标题和内容不能为空")
		return
	}

	reply := models.QuickReply{
		Title:     input.Title,
		Content:   input.Content,
		CreatedBy: userID.(uint),
		SortOrder: input.SortOrder,
	}
	if err := database.DB.Create(&reply).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建失败")
		return
	}

	database.DB.Preload("User").First(&reply, reply.ID)
	c.JSON(http.StatusOK, gin.H{"quickReply": reply})
}

func UpdateQuickReply(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的ID")
		return
	}

	var reply models.QuickReply
	if err := database.DB.First(&reply, id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "快捷回复不存在")
		return
	}

	var input struct {
		Title     *string `json:"title"`
		Content   *string `json:"content"`
		SortOrder *int    `json:"sortOrder"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if input.Title != nil {
		updates["title"] = *input.Title
	}
	if input.Content != nil {
		updates["content"] = *input.Content
	}
	if input.SortOrder != nil {
		updates["sort_order"] = *input.SortOrder
	}
	database.DB.Model(&reply).Updates(updates)
	database.DB.Preload("User").First(&reply, reply.ID)

	c.JSON(http.StatusOK, gin.H{"quickReply": reply})
}

func DeleteQuickReply(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的ID")
		return
	}
	database.DB.Delete(&models.QuickReply{}, id)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetSatisfactionStats(c *gin.Context) {
	var ratings []models.Rating
	database.DB.Preload("Conversation.User").Order("created_at desc").Find(&ratings)

	var totalScore int
	scoreDistribution := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	for _, r := range ratings {
		totalScore += r.Score
		scoreDistribution[r.Score]++
	}

	avgScore := float64(0)
	if len(ratings) > 0 {
		avgScore = float64(totalScore) / float64(len(ratings))
	}

	c.JSON(http.StatusOK, gin.H{
		"total":              len(ratings),
		"averageScore":       avgScore,
		"scoreDistribution":  scoreDistribution,
		"ratings":            ratings,
	})
}

func GetAgentStats(c *gin.Context) {
	type AgentStat struct {
		UserID       uint   `json:"userId"`
		UserName     string `json:"userName"`
		TotalMsgs    int64  `json:"totalMessages"`
		AvgRating    float64
		RatingCount  int64 `json:"ratingCount"`
	}

	var agents []models.User
	database.DB.Where("role IN ?", []string{"admin", "service"}).Find(&agents)

	var result []AgentStat
	for _, agent := range agents {
		var totalMsgs int64
		database.DB.Model(&models.Message{}).
			Where("sender_type = ? AND sender_id = ?", "service", agent.ID).
			Count(&totalMsgs)

		var avgRating float64
		var ratingCount int64
		database.DB.Model(&models.Rating{}).
			Joins("JOIN conversations ON conversations.id = ratings.conversation_id").
			Where("conversations.assigned_to = ?", agent.ID).
			Count(&ratingCount)

		if ratingCount > 0 {
			var totalScore float64
			database.DB.Model(&models.Rating{}).
				Joins("JOIN conversations ON conversations.id = ratings.conversation_id").
				Where("conversations.assigned_to = ?", agent.ID).
				Select("COALESCE(SUM(score), 0)").
				Scan(&totalScore)
			avgRating = totalScore / float64(ratingCount)
		}

		result = append(result, AgentStat{
			UserID:      agent.ID,
			UserName:    agent.Name,
			TotalMsgs:   totalMsgs,
			AvgRating:   avgRating,
			RatingCount: ratingCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{"agents": result})
}
