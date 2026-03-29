package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SecurityTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *SecurityTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestSecuritySuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}

func (s *SecurityTestSuite) TestAdminEndpointNoAuth() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/admin/products"},
		{"GET", "/api/admin/categories"},
		{"GET", "/api/admin/banners"},
		{"GET", "/api/admin/pages"},
		{"POST", "/api/admin/products"},
		{"POST", "/api/admin/categories"},
	}
	for _, ep := range endpoints {
		w := MakeRequest(s.router, ep.method, ep.path, nil, nil)
		s.True(w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden,
			"%s %s should return 401 or 403 without auth, got %d", ep.method, ep.path, w.Code)
	}
}

func (s *SecurityTestSuite) TestAdminEndpointCustomerForbidden() {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/admin/products"},
		{"GET", "/api/admin/categories"},
		{"GET", "/api/admin/banners"},
		{"GET", "/api/admin/pages"},
		{"POST", "/api/admin/products"},
		{"POST", "/api/admin/categories"},
	}
	headers := GetCustomerAuthHeader()
	for _, ep := range endpoints {
		w := MakeRequest(s.router, ep.method, ep.path, nil, headers)
		s.Equal(http.StatusForbidden, w.Code,
			"%s %s should return 403 for customer", ep.method, ep.path)
	}
}

func (s *SecurityTestSuite) TestConfigEndpointNoAuth() {
	body := map[string]interface{}{
		"value": "hacked",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/config/security_test_key", body, nil)
	s.True(w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden,
		"PUT /api/admin/config/:key should require admin auth, got %d", w.Code)
}

func (s *SecurityTestSuite) TestJWTForgery() {
	headers := map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.fake",
	}
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *SecurityTestSuite) TestCrossUserCartAccess() {
	cartItemID := Data.CartItems[0].ID
	body := map[string]interface{}{
		"quantity": 99,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/cart/%d", cartItemID), body, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *SecurityTestSuite) TestCrossUserOrderAccess() {
	order := Data.Orders[0]
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d", order.ID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *SecurityTestSuite) TestCrossUserAddressAccess() {
	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"fullName": "Hacker",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/addresses/%d", addrID), body, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *SecurityTestSuite) TestSQLInjection() {
	w := MakeRequest(s.router, "GET", "/api/products?search=%27%3B+DROP+TABLE+products%3B--", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/products", nil, nil)
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["total"].(float64) > 0, "products table should still exist")
}

func (s *SecurityTestSuite) TestXSSInProductName() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	body := map[string]interface{}{
		"name":       "<script>alert(1)</script>",
		"price":      99,
		"categoryId": classicID,
		"stock":      1,
		"isActive":   true,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("<script>alert(1)</script>", resp["name"])
}

func (s *SecurityTestSuite) TestDevBypass() {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("user-%d", CustomerUser.ID),
	}
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusOK, w.Code)
}

func (s *SecurityTestSuite) TestCORSHeaders() {
	w := MakeRequest(s.router, "GET", "/api/products", nil, nil)
	s.Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
}
