package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ChatTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *ChatTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestChatSuite(t *testing.T) {
	suite.Run(t, new(ChatTestSuite))
}

func (s *ChatTestSuite) TestCreateConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	conv := resp["conversation"].(map[string]interface{})
	s.Equal("open", conv["status"])
	s.NotZero(conv["id"])
}

func (s *ChatTestSuite) TestCreateConversationTwice() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp1 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp1)
	id1 := resp1["conversation"].(map[string]interface{})["id"].(float64)

	w2 := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w2.Code)

	var resp2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp2)
	id2 := resp2["conversation"].(map[string]interface{})["id"].(float64)

	s.Equal(id1, id2, "second create should return existing conversation")
}

func (s *ChatTestSuite) TestCreateConversationUnauthorized() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ChatTestSuite) TestGetConversations() {
	MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())

	w := MakeRequest(s.router, "GET", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convs := resp["conversations"].([]interface{})
	s.True(len(convs) >= 1)
}

func (s *ChatTestSuite) TestGetConversationsEmpty() {
	w := MakeRequest(s.router, "GET", "/api/chat/conversations", nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convs := resp["conversations"].([]interface{})
	s.Equal(0, len(convs))
}

func (s *ChatTestSuite) TestGetConversationsUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/chat/conversations", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ChatTestSuite) TestSendMessage() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	body := map[string]interface{}{"content": "测试消息"}
	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var msgResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msgResp)
	msg := msgResp["message"].(map[string]interface{})
	s.Equal("测试消息", msg["content"])
	s.Equal("customer", msg["senderType"])
}

func (s *ChatTestSuite) TestSendMessageTooLong() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	longContent := make([]byte, 501)
	for i := range longContent {
		longContent[i] = 'a'
	}
	body := map[string]interface{}{"content": string(longContent)}
	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestSendMessageEmptyContent() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	body := map[string]interface{}{"content": ""}
	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestSendMessageUnauthorized() {
	body := map[string]interface{}{"content": "test"}
	w := MakeRequest(s.router, "POST", "/api/chat/conversations/1/messages", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ChatTestSuite) TestSendMessageImageType() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	body := map[string]interface{}{
		"content":      "/media/test.jpg",
		"messageType":  "image",
		"thumbnailUrl": "/media/test_thumb.jpg",
	}
	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var msgResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msgResp)
	msg := msgResp["message"].(map[string]interface{})
	s.Equal("image", msg["messageType"])
}

func (s *ChatTestSuite) TestGetMessages() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "hello"}, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "GET", fmt.Sprintf("/api/chat/conversations/%d/messages", convID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var msgResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msgResp)
	msgs := msgResp["messages"].([]interface{})
	s.True(len(msgs) >= 2, "should have greeting + our message")
}

func (s *ChatTestSuite) TestGetMessagesWithAfterID() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "msg1"}, GetCustomerAuthHeader())
	var msg1Resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msg1Resp)
	msg1ID := int(msg1Resp["message"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "msg2"}, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "GET", fmt.Sprintf("/api/chat/conversations/%d/messages?after=%d", convID, msg1ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var msgsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &msgsResp)
	msgs := msgsResp["messages"].([]interface{})
	s.Equal(1, len(msgs), "should only return message after the given ID")
}

func (s *ChatTestSuite) TestGetMessagesNotOwnConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/chat/conversations/99999/messages", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *ChatTestSuite) TestCloseConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var closeResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &closeResp)
	s.Equal(true, closeResp["success"])

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/messages", convID),
		map[string]interface{}{"content": "after close"}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestCloseConversationAlreadyClosed() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestRateConversation() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())

	body := map[string]interface{}{"score": 5, "comment": "很好"}
	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var rateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &rateResp)
	rating := rateResp["rating"].(map[string]interface{})
	s.Equal(float64(5), rating["score"])
	s.Equal("很好", rating["comment"])
}

func (s *ChatTestSuite) TestRateConversationNotClosed() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID),
		map[string]interface{}{"score": 5}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestRateConversationAlreadyRated() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())
	MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID),
		map[string]interface{}{"score": 5}, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID),
		map[string]interface{}{"score": 4}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestRateConversationInvalidScore() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	convID := int(resp["conversation"].(map[string]interface{})["id"].(float64))

	MakeRequest(s.router, "PUT", fmt.Sprintf("/api/chat/conversations/%d/close", convID), nil, GetCustomerAuthHeader())

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID),
		map[string]interface{}{"score": 0}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	w = MakeRequest(s.router, "POST", fmt.Sprintf("/api/chat/conversations/%d/rate", convID),
		map[string]interface{}{"score": 6}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestGetChatUnreadCount() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/chat/unread-count", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["count"])
}

func (s *ChatTestSuite) TestGetChatUnreadCountUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/chat/unread-count", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ChatTestSuite) TestGetServiceStatus() {
	w := MakeRequest(s.router, "GET", "/api/chat/service-status", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["online"])
}

func (s *ChatTestSuite) TestUploadChatImage() {
	w := MakeRequest(s.router, "POST", "/api/chat/conversations", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	imgFile := CreateTestImageFile(".jpg")
	w = MakeFormRequest(s.router, "POST", "/api/chat/upload-image", map[string]string{}, map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	s.NotEmpty(resp["url"])
}

func (s *ChatTestSuite) TestUploadChatImageInvalidExt() {
	w := MakeFormRequest(s.router, "POST", "/api/chat/upload-image", map[string]string{}, map[string]string{"file": CreateTestImageFile(".exe")}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestUploadChatImageTooLarge() {
	imgFile := CreateLargeTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", "/api/chat/upload-image", map[string]string{}, map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ChatTestSuite) TestUploadChatImageUnauthorized() {
	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", "/api/chat/upload-image", map[string]string{}, map[string]string{"file": imgFile}, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}
