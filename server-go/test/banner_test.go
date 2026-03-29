package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type BannerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *BannerTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestBannerSuite(t *testing.T) {
	suite.Run(t, new(BannerTestSuite))
}

func (s *BannerTestSuite) TestGetBannersOnlyActive() {
	w := MakeRequest(s.router, "GET", "/api/banners", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	for _, banner := range resp {
		s.Equal(true, banner["isActive"])
	}
}

func (s *BannerTestSuite) TestGetBannersSortOrder() {
	w := MakeRequest(s.router, "GET", "/api/banners", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) >= 2 {
		first := resp[0]["sortOrder"].(float64)
		second := resp[1]["sortOrder"].(float64)
		s.True(first <= second, "banners should be sorted by sortOrder ASC")
	}
}

func (s *BannerTestSuite) TestAdminGetBannersAll() {
	w := MakeRequest(s.router, "GET", "/api/admin/banners", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(len(resp) >= 3, "should include inactive banners")

	hasInactive := false
	for _, banner := range resp {
		if banner["isActive"] == false {
			hasInactive = true
		}
	}
	s.True(hasInactive, "admin should see inactive banners")
}

func (s *BannerTestSuite) TestAdminCreateBanner() {
	body := map[string]interface{}{
		"title":     "测试轮播图",
		"imageUrl":  "/test/new-banner.jpg",
		"link":      "/test",
		"sortOrder": 10,
		"isActive":  true,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/banners", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("测试轮播图", resp["title"])
}

func (s *BannerTestSuite) TestAdminUpdateBanner() {
	bannerID := Data.Banners[0].ID
	body := map[string]interface{}{
		"title":     "更新标题",
		"imageUrl":  "/test/updated.jpg",
		"link":      "/updated",
		"sortOrder": 99,
		"isActive":  false,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/banners/%d", bannerID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("更新标题", resp["title"])
}

func (s *BannerTestSuite) TestAdminDeleteBanner() {
	body := map[string]interface{}{
		"title":     "待删除轮播图",
		"imageUrl":  "/test/delete-me.jpg",
		"sortOrder": 50,
		"isActive":  false,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/banners", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	bannerID := int(created["id"].(float64))

	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/banners/%d", bannerID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *BannerTestSuite) TestAdminUpdateBannerNotFound() {
	body := map[string]interface{}{
		"title": "test",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/banners/99999", body, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *BannerTestSuite) TestAdminBannerForbidden() {
	w := MakeRequest(s.router, "GET", "/api/admin/banners", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}
