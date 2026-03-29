package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type PageTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *PageTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(PageTestSuite))
}

func (s *PageTestSuite) TestGetPage() {
	w := MakeRequest(s.router, "GET", "/api/pages/services", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("服务条款", resp["title"])
	s.NotNil(resp["content"])
}

func (s *PageTestSuite) TestGetPageNotFound() {
	w := MakeRequest(s.router, "GET", "/api/pages/nonexistent", nil, nil)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *PageTestSuite) TestGetCommonSlugs() {
	slugs := []string{"about", "services", "privacy-policy", "statement"}
	for _, slug := range slugs {
		w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/pages/%s", slug), nil, nil)
		s.Equal(http.StatusOK, w.Code, "page %s should be accessible", slug)
	}
}

func (s *PageTestSuite) TestAdminGetPages() {
	w := MakeRequest(s.router, "GET", "/api/admin/pages", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(len(resp) >= 4, "should have at least 4 pages")
}

func (s *PageTestSuite) TestAdminUpdatePage() {
	body := map[string]interface{}{
		"title":   "更新标题",
		"content": "# 更新内容\n\n新内容",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/pages/about", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("更新标题", resp["title"])
	s.Equal("# 更新内容\n\n新内容", resp["content"])
}

func (s *PageTestSuite) TestAdminUpdatePageNotFound() {
	body := map[string]interface{}{
		"title":   "test",
		"content": "test",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/pages/nonexistent-page-xyz", body, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *PageTestSuite) TestAdminPageForbidden() {
	w := MakeRequest(s.router, "GET", "/api/admin/pages", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}
