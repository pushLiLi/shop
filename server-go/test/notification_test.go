package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type NotificationTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *NotificationTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func (s *NotificationTestSuite) SetupTest() {
	database.DB.Exec("DELETE FROM notifications")
}

func TestNotificationSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (s *NotificationTestSuite) seedNotifications() []models.Notification {
	notifs := []models.Notification{
		{UserID: CustomerUser.ID, Type: "order_status", Title: "订单已发货", Content: "您的订单已发出", IsRead: false},
		{UserID: CustomerUser.ID, Type: "back_in_stock", Title: "商品到货", Content: "您关注的商品已到货", IsRead: false},
		{UserID: CustomerUser.ID, Type: "price_drop", Title: "降价提醒", Content: "商品降价了", IsRead: true},
		{UserID: CustomerUser.ID, Type: "order_status", Title: "订单已完成", Content: "订单已签收", IsRead: true},
	}
	database.DB.Create(&notifs)
	return notifs
}

func (s *NotificationTestSuite) TestGetNotifications() {
	s.seedNotifications()

	w := MakeRequest(s.router, "GET", "/api/notifications?page=1&limit=20", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	notifs := resp["notifications"].([]interface{})
	s.True(len(notifs) >= 4)
}

func (s *NotificationTestSuite) TestGetNotificationsUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/notifications", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *NotificationTestSuite) TestGetUnreadCount() {
	s.seedNotifications()

	w := MakeRequest(s.router, "GET", "/api/notifications/unread-count", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(2), resp["count"])
}

func (s *NotificationTestSuite) TestGetNotification() {
	notifs := s.seedNotifications()

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/notifications/%d", notifs[0].ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var notif map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notif)
	s.Equal(notifs[0].Title, notif["title"])
}

func (s *NotificationTestSuite) TestGetNotificationAutoRead() {
	notifs := s.seedNotifications()
	unreadNotif := notifs[0]
	s.False(unreadNotif.IsRead)

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/notifications/%d", unreadNotif.ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var notif map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notif)
	s.True(notif["isRead"].(bool))
}

func (s *NotificationTestSuite) TestMarkAsRead() {
	notifs := s.seedNotifications()
	unreadNotif := notifs[0]

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/notifications/%d/read", unreadNotif.ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["success"].(bool))
}

func (s *NotificationTestSuite) TestMarkAllRead() {
	s.seedNotifications()

	w := MakeRequest(s.router, "PUT", "/api/notifications/read-all", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["success"].(bool))

	var count int64
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", CustomerUser.ID, false).Count(&count)
	s.Equal(int64(0), count)
}

func (s *NotificationTestSuite) TestDeleteNotification() {
	notifs := s.seedNotifications()

	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/notifications/%d", notifs[0].ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["success"].(bool))

	var count int64
	database.DB.Model(&models.Notification{}).Where("id = ?", notifs[0].ID).Count(&count)
	s.Equal(int64(0), count)
}

func (s *NotificationTestSuite) TestDeleteNotificationNotOwned() {
	notifs := s.seedNotifications()

	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/notifications/%d", notifs[0].ID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *NotificationTestSuite) TestDeleteReadNotifications() {
	s.seedNotifications()

	w := MakeRequest(s.router, "DELETE", "/api/notifications/read", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["success"].(bool))
	s.Equal(float64(2), resp["deleted"].(float64))

	var count int64
	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", CustomerUser.ID, true).Count(&count)
	s.Equal(int64(0), count)

	database.DB.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", CustomerUser.ID, false).Count(&count)
	s.Equal(int64(2), count)
}
