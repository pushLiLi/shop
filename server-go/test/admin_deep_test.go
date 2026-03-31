package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AdminDeepTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *AdminDeepTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestAdminDeepSuite(t *testing.T) {
	suite.Run(t, new(AdminDeepTestSuite))
}

func (s *AdminDeepTestSuite) TestAdminGetConversations() {
	w := MakeRequest(s.router, "GET", "/api/admin/chat/conversations", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["conversations"])
}

func (s *AdminDeepTestSuite) TestAdminGetConversationsFilterByStatus() {
	MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())

	w := MakeRequest(s.router, "GET", "/api/admin/chat/conversations?status=open", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/chat/conversations?status=closed", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetConversationsFilterAssignedTo() {
	MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())

	w := MakeRequest(s.router, "GET", "/api/admin/chat/conversations?assignedTo=unassigned", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/chat/conversations?assignedTo=me", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetConversationsSort() {
	MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())

	sortCases := []struct {
		sortBy    string
		sortOrder string
	}{
		{"lastMessageAt", "desc"},
		{"lastMessageAt", "asc"},
		{"createdAt", "desc"},
		{"createdAt", "asc"},
		{"status", "asc"},
	}
	for _, c := range sortCases {
		w := MakeRequest(s.router,
			fmt.Sprintf("/api/admin/chat/conversations?sortBy=%s&sortOrder=%s", c.sortBy, c.sortOrder),
			nil, GetAdminAuthHeader())
		s.Equal(http.StatusOK, w.Code, "sortBy=%s sortOrder=%s should return 200", c.sortBy, c.sortOrder)
	}
}

func (s *AdminDeepTestSuite) TestAdminSendMessage() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	body := map[string]interface{}{"content": "客服回复您好"}
	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages", convID),
		body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var msgResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msgResp)
	msg := msgResp["message"].(map[string]interface{})
	s.Equal("客服回复您好", msg["content"])
	s.Equal("service", msg["senderType"])
}

func (s *AdminDeepTestSuite) TestAdminSendMessageClosedConv() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "PUT",
		fmt.Sprintf("/api/chat/conversations/%d/close", convID),
		nil, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "closed reply"}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminSendMessageTooLong() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	longContent := make([]byte, 501)
	for i := range longContent {
		longContent[i] = 'x'
	}
	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": string(longContent)}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AdminDeepTestSuite) TestCloseConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT",
		fmt.Sprintf("/api/admin/chat/conversations/%d/close", convID),
		nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var closeResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &closeResp)
	s.Equal(true, closeResp["success"])
}

func (s *AdminDeepTestSuite) TestRecallMessage() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "recall me"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	msgID := int(resp["message"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages/%d/recall", convID, msgID),
		nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestRecallMessageNotOwn() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "customer msg"}, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
	msgID := int(resp["message"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST",
		fmt.Sprintf("/api/admin/chat/conversations/%d/messages/%d/recall", convID, msgID),
		nil, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *AdminDeepTestSuite) TestAssignConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	uid := AdminUser.ID
	body := map[string]interface{}{"assignedTo": uid}
	w = MakeRequest(s.router, "PUT",
		fmt.Sprintf("/api/admin/chat/conversations/%d/assign", convID),
		body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAssignConversationToNull() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	body := map[string]interface{}{"assignedTo": nil}
	w = MakeRequest(s.router, "PUT",
		fmt.Sprintf("/api/admin/chat/conversations/%d/assign", convID),
		body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestGetQuickReplies() {
	w := MakeRequest(s.router, "GET", "/api/admin/quick-replies", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestCreateQuickReply() {
	body := map[string]interface{}{
		"title":     "测试快捷回复",
		"content":   "感谢您的咨询",
		"sortOrder": 10,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/quick-replies", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("测试快捷回复", resp["quickReply"].(map[string]interface{})["title"])
}

func (s *AdminDeepTestSuite) TestUpdateQuickReply() {
	body := map[string]interface{}{
		"title":     "原始标题",
		"content":   "原始内容",
		"sortOrder": 1,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/quick-replies", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	qrID := int(resp["quickReply"].(map[string]interface{})["id"].(float64))

	newTitle := "更新标题"
	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/quick-replies/%d", qrID),
		map[string]interface{}{"title": &newTitle}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/quick-replies", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &resp)
}

func (s *AdminDeepTestSuite) TestDeleteQuickReply() {
	body := map[string]interface{}{
		"title":   "待删除",
		"content": "内容",
	}
	w := MakeRequest(s.router, "POST", "/api/admin/quick-replies", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	qrID := int(resp["quickReply"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/quick-replies/%d", qrID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestGetSatisfactionStats() {
	w := MakeRequest(s.router, "GET", "/api/admin/stats/satisfaction", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["ratings"])
}

func (s *AdminDeepTestSuite) TestGetAgentStats() {
	w := MakeRequest(s.router, "GET", "/api/admin/chat/agent-stats", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestGetAdminUnreadStats() {
	w := MakeRequest(s.router, "GET", "/api/admin/chat/unread-stats", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["totalUnread"])
}

func (s *AdminDeepTestSuite) TestSetServiceStatus() {
	body := map[string]interface{}{"online": true}
	w := MakeRequest(s.router, "POST", "/api/admin/chat/service-status", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	body = map[string]interface{}{"online": false}
	w = MakeRequest(s.router, "POST", "/api/admin/chat/service-status", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetDashboardStats() {
	w := MakeRequest(s.router, "GET", "/api/admin/dashboard/stats", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["stats"])
}

func (s *AdminDeepTestSuite) TestAdminGetDashboardTopProducts() {
	w := MakeRequest(s.router, "GET", "/api/admin/dashboard/top-products", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetDashboardLowStock() {
	w := MakeRequest(s.router, "GET", "/api/admin/dashboard/low-stock", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetDashboardRecentOrders() {
	w := MakeRequest(s.router, "GET", "/api/admin/dashboard/recent-orders", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetRevenue() {
	w := MakeRequest(s.router, "GET", "/api/admin/stats/revenue", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCanAccessDashboard() {
	w := MakeRequest(s.router, "GET", "/api/admin/dashboard/stats", nil, GetServiceAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotAccessUsers() {
	w := MakeRequest(s.router, "GET", "/api/admin/users", nil, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotAccessConfig() {
	w := MakeRequest(s.router, "PUT", "/api/admin/config/test_key",
		map[string]interface{}{"value": "hacked"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotChangeUserRole() {
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/users/%d/role", CustomerUser.ID),
		map[string]interface{}{"role": "admin"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminCanChangeUserRole() {
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/users/%d/role", CustomerUser.ID),
		map[string]interface{}{"role": "service"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminOrdersPagination() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders?page=1&limit=5", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(5), resp["limit"])
	s.Equal(float64(1), resp["page"])
}

func (s *AdminDeepTestSuite) TestAdminOrdersFilterByStatus() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders?status=pending", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminOrdersQuickFilter() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders?quick_filter=pending_proof", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/orders?quick_filter=shipped", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/orders?quick_filter=completed", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminOrdersSearch() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders?search=test", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminOrdersSort() {
	sortCases := []struct {
		sortBy    string
		sortOrder string
	}{
		{"id", "desc"},
		{"orderNo", "desc"},
		{"total", "asc"},
		{"status", "desc"},
		{"createdAt", "asc"},
		{"id", "invalid"},
		{"invalid", "desc"},
		{"", ""},
	}
	for _, c := range sortCases {
		url := fmt.Sprintf("/api/admin/orders?sortBy=%s&sortOrder=%s", c.sortBy, c.sortOrder)
		w := MakeRequest(s.router, url, nil, GetAdminAuthHeader())
		s.Equal(http.StatusOK, w.Code, "sortBy=%s sortOrder=%s should return 200", c.sortBy, c.sortOrder)
	}
}

func (s *AdminDeepTestSuite) TestAdminOrdersExport() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders/export", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Header().Get("Content-Type"), "csv")
}

func (s *AdminDeepTestSuite) TestServiceCanAccessOrders() {
	w := MakeRequest(s.router, "GET", "/api/admin/orders", nil, GetServiceAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotAccessBanners() {
	w := MakeRequest(s.router, "GET", "/api/admin/banners", nil, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminCanAccessBanners() {
	w := MakeRequest(s.router, "GET", "/api/admin/banners", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestCustomerCannotAccessAdmin() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/admin/dashboard/stats"},
		{"GET", "/api/admin/orders"},
		{"GET", "/api/admin/products"},
		{"GET", "/api/admin/categories"},
		{"GET", "/api/admin/chat/conversations"},
		{"GET", "/api/admin/quick-replies"},
		{"GET", "/api/admin/users"},
	}
	for _, ep := range endpoints {
		w := MakeRequest(s.router, ep.method, ep.path, nil, GetCustomerAuthHeader())
		s.Equal(http.StatusForbidden, w.Code,
			"%s %s should return 403 for customer", ep.method, ep.path)
	}
}

func (s *AdminDeepTestSuite) TestAdminProductsSort() {
	sortCases := []struct {
		sortBy    string
		sortOrder string
	}{
		{"price", "asc"},
		{"price", "desc"},
		{"createdAt", "asc"},
		{"createdAt", "desc"},
		{"stock", "asc"},
		{"stock", "desc"},
		{"invalid", "desc"},
	}
	for _, c := range sortCases {
		url := fmt.Sprintf("/api/admin/products?sortBy=%s&sortOrder=%s", c.sortBy, c.sortOrder)
		w := MakeRequest(s.router, url, nil, GetAdminAuthHeader())
		s.Equal(http.StatusOK, w.Code)
	}
}

func (s *AdminDeepTestSuite) TestServiceCanManageProducts() {
	w := MakeRequest(s.router, "GET", "/api/admin/products", nil, GetServiceAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotManageBanners() {
	w := MakeRequest(s.router, "POST", "/api/admin/banners",
		map[string]interface{}{"title": "x", "image": "x"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotManagePaymentMethods() {
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods",
		map[string]interface{}{"name": "x"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestServiceCannotManageContactMethods() {
	w := MakeRequest(s.router, "POST", "/api/admin/contact-methods",
		map[string]interface{}{"type": "whatsapp", "label": "x", "value": "x"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *AdminDeepTestSuite) TestAdminGetUsers() {
	w := MakeRequest(s.router, "GET", "/api/admin/users", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["users"])
}

func (s *AdminDeepTestSuite) TestServiceCanChat() {
	w := MakeRequest(s.router, "GET", "/api/admin/chat/conversations", nil, GetServiceAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}
