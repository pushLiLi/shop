package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ConfigSettingTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *ConfigSettingTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestConfigSettingSuite(t *testing.T) {
	suite.Run(t, new(ConfigSettingTestSuite))
}

func (s *ConfigSettingTestSuite) TestGetConfig() {
	w := MakeRequest(s.router, "GET", "/api/config", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp)
}

func (s *ConfigSettingTestSuite) TestGetSettings() {
	w := MakeRequest(s.router, "GET", "/api/settings", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	s.NotNil(resp["data"])
}

func (s *ConfigSettingTestSuite) TestUpdateSetting() {
	body := map[string]interface{}{
		"value": "test value",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/settings/test_key", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *ConfigSettingTestSuite) TestUpdateSettingCreateNew() {
	body := map[string]interface{}{
		"value": "new setting value",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/settings/brand_new_setting", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *ConfigSettingTestSuite) TestUpdateSettingForbidden() {
	body := map[string]interface{}{
		"value": "test",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/settings/test", body, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *ConfigSettingTestSuite) TestUpdateConfig() {
	body := map[string]interface{}{
		"value": "updated config value",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/config/test_config_key", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *ConfigSettingTestSuite) TestUpdateConfigNoAuth() {
	body := map[string]interface{}{
		"value": "hacked value",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/config/unprotected_key", body, nil)
	s.True(w.Code == http.StatusUnauthorized || w.Code == http.StatusForbidden,
		"PUT /api/admin/config/:key should require auth, got %d", w.Code)
}
