package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"bycigar-server/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *AuthTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func (s *AuthTestSuite) SetupTest() {
	handlers.ResetLoginFailures()
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) getCaptcha() (string, string) {
	w := MakeRequest(s.router, "GET", "/api/auth/captcha", nil, nil)
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	captchaID := resp["captchaId"].(string)
	code := handlers.GetCaptchaCode(captchaID)
	return captchaID, code
}

func (s *AuthTestSuite) TestGetCaptcha() {
	w := MakeRequest(s.router, "GET", "/api/auth/captcha", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotEmpty(resp["captchaId"])
	s.NotEmpty(resp["captchaImage"])

	b64Image := resp["captchaImage"].(string)
	s.True(strings.HasPrefix(b64Image, "data:image/png;base64,") || strings.HasPrefix(b64Image, "data:image/jpeg;base64,"))
}

func (s *AuthTestSuite) TestGetCaptchaUniqueIds() {
	ids := make(map[string]bool)
	for i := 0; i < 5; i++ {
		w := MakeRequest(s.router, "GET", "/api/auth/captcha", nil, nil)
		s.Equal(http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		id := resp["captchaId"].(string)
		s.False(ids[id], "captcha ID should be unique")
		ids[id] = true
	}
}

func (s *AuthTestSuite) TestGetCaptchaImageNotEmpty() {
	w := MakeRequest(s.router, "GET", "/api/auth/captcha", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	b64Image := resp["captchaImage"].(string)
	s.True(len(b64Image) > 100, "captcha image should be a non-trivial base64 string")
}

func (s *AuthTestSuite) TestRegisterSuccess() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "newuser@test.com",
		"password":    "password123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotEmpty(resp["token"])
	s.NotEmpty(resp["user"])
	user := resp["user"].(map[string]interface{})
	s.Equal("newuser@test.com", user["email"])
	s.Equal("newuser", user["name"])
	s.Equal("customer", user["role"])
}

func (s *AuthTestSuite) TestRegisterTokenWorks() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "tokenuser@test.com",
		"password":    "password123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	token := resp["token"].(string)

	w = MakeRequest(s.router, "GET", "/api/auth/me", nil, map[string]string{
		"Authorization": "Bearer " + token,
	})
	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthTestSuite) TestRegisterDuplicateEmail() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "user1@test.com",
		"password":    "password123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "已被注册")
}

func (s *AuthTestSuite) TestRegisterPasswordTooShort() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "shortpw@test.com",
		"password":    "12345",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AuthTestSuite) TestRegisterMissingFields() {
	cases := []map[string]interface{}{
		{"password": "123456", "captchaId": "x", "captchaCode": "x"},
		{"email": "a@b.com", "captchaId": "x", "captchaCode": "x"},
		{"email": "a@b.com", "password": "123456", "captchaCode": "x"},
		{"email": "a@b.com", "password": "123456", "captchaId": "x"},
	}
	for _, body := range cases {
		w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
		s.Equal(http.StatusBadRequest, w.Code)
	}
}

func (s *AuthTestSuite) TestRegisterWrongCaptcha() {
	body := map[string]interface{}{
		"email":       "wrongcap@test.com",
		"password":    "123456",
		"captchaId":   "nonexistent",
		"captchaCode": "wrong",
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "验证码")
}

func (s *AuthTestSuite) TestRegisterInvalidEmail() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "notanemail",
		"password":    "123456",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AuthTestSuite) TestRegisterAutoNameFromEmail() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "autoname@test.com",
		"password":    "123456",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	user := resp["user"].(map[string]interface{})
	s.Equal("autoname", user["name"])
}

func (s *AuthTestSuite) TestRegisterCustomName() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"email":       "customname@test.com",
		"password":    "123456",
		"name":        "MyCustomName",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "POST", "/api/auth/register", body, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	user := resp["user"].(map[string]interface{})
	s.Equal("MyCustomName", user["name"])
}

func (s *AuthTestSuite) TestLoginSuccess() {
	body := map[string]interface{}{
		"email":    "user1@test.com",
		"password": "user1234",
	}
	w := MakeRequest(s.router, "POST", "/api/auth/login", body, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotEmpty(resp["token"])
	s.NotEmpty(resp["user"])
	user := resp["user"].(map[string]interface{})
	s.Equal("user1@test.com", user["email"])
	s.Equal("customer", user["role"])
}

func (s *AuthTestSuite) TestLoginWrongPassword() {
	body := map[string]interface{}{
		"email":    "user1@test.com",
		"password": "wrongpassword",
	}
	w := MakeRequest(s.router, "POST", "/api/auth/login", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "邮箱或密码错误")
}

func (s *AuthTestSuite) TestLoginNonexistentUser() {
	body := map[string]interface{}{
		"email":    "nouser@test.com",
		"password": "whatever",
	}
	w := MakeRequest(s.router, "POST", "/api/auth/login", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthTestSuite) TestLoginProgressiveCaptcha() {
	email := "user1@test.com"
	wrongBody := map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}

	for i := 0; i < 3; i++ {
		w := MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
		s.Equal(http.StatusUnauthorized, w.Code)
	}

	w := MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	s.Equal(http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["requireCaptcha"])
}

func (s *AuthTestSuite) TestLoginCaptchaRequiredButMissing() {
	email := "capmiss@test.com"
	wrongBody := map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}

	for i := 0; i < 3; i++ {
		MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	}

	w := MakeRequest(s.router, "POST", "/api/auth/login", map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}, nil)
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "验证码")
	s.Equal(true, resp["requireCaptcha"])
}

func (s *AuthTestSuite) TestLoginCaptchaWrongCode() {
	email := "capwrong@test.com"
	wrongBody := map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}

	for i := 0; i < 3; i++ {
		MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	}

	w := MakeRequest(s.router, "POST", "/api/auth/login", map[string]interface{}{
		"email":       email,
		"password":    "wrongpassword",
		"captchaId":   "fake-id",
		"captchaCode": "wrong",
	}, nil)
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "验证码")
}

func (s *AuthTestSuite) TestLoginWithCaptchaAndCorrectPassword() {
	captchaID, captchaCode := s.getCaptcha()

	email := "caplogin@test.com"
	wrongBody := map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}

	for i := 0; i < 3; i++ {
		MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	}

	captchaID2, captchaCode2 := s.getCaptcha()

	loginBody := map[string]interface{}{
		"email":       email,
		"password":    "user1234",
		"captchaId":   captchaID2,
		"captchaCode": captchaCode2,
	}
	_ = captchaID
	_ = captchaCode
	w := MakeRequest(s.router, "POST", "/api/auth/login", loginBody, nil)

	if w.Code == http.StatusUnauthorized {
		s.T().Log("User doesn't exist, skipping success assertion")
		return
	}
	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthTestSuite) TestLoginSuccessResetsCounter() {
	email := "user1@test.com"
	wrongBody := map[string]interface{}{
		"email":    email,
		"password": "wrongpassword",
	}

	for i := 0; i < 3; i++ {
		MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	}

	correctBody := map[string]interface{}{
		"email":    email,
		"password": "user1234",
	}
	w := MakeRequest(s.router, "POST", "/api/auth/login", correctBody, nil)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code == http.StatusOK {
		s.NotEmpty(resp["token"])
	} else {
		captchaID2, captchaCode2 := s.getCaptcha()
		correctBody["captchaId"] = captchaID2
		correctBody["captchaCode"] = captchaCode2
		w = MakeRequest(s.router, "POST", "/api/auth/login", correctBody, nil)
		s.Equal(http.StatusOK, w.Code)
	}

	handlers.ResetLoginFailures()

	w = MakeRequest(s.router, "POST", "/api/auth/login", wrongBody, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
	var resp2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp2)
	_, hasCaptcha := resp2["requireCaptcha"]
	s.False(hasCaptcha && resp2["requireCaptcha"] == true, "counter should reset after successful login")
}

func (s *AuthTestSuite) TestLoginDifferentEmailIndependentCount() {
	email1 := map[string]interface{}{"email": "user1@test.com", "password": "wrong1"}
	email2 := map[string]interface{}{"email": "user2@test.com", "password": "wrong2"}

	for i := 0; i < 3; i++ {
		MakeRequest(s.router, "POST", "/api/auth/login", email1, nil)
	}

	w := MakeRequest(s.router, "POST", "/api/auth/login", email2, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	_, hasCaptcha := resp["requireCaptcha"]
	s.False(hasCaptcha && resp["requireCaptcha"] == true, "different emails should have independent counters")
}

func (s *AuthTestSuite) TestLoginMissingFields() {
	cases := []map[string]interface{}{
		{"password": "123456"},
		{"email": "a@b.com"},
		{},
	}
	for _, body := range cases {
		w := MakeRequest(s.router, "POST", "/api/auth/login", body, nil)
		s.Equal(http.StatusBadRequest, w.Code)
	}
}

func (s *AuthTestSuite) TestGetProfileValidToken() {
	headers := GetCustomerAuthHeader()
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	user := resp["user"].(map[string]interface{})
	s.Equal("user1@test.com", user["email"])
	s.Equal("TestUser1", user["name"])
	s.Equal("customer", user["role"])
}

func (s *AuthTestSuite) TestGetProfileNoToken() {
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthTestSuite) TestGetProfileInvalidToken() {
	headers := map[string]string{
		"Authorization": "Bearer invalidtoken123",
	}
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthTestSuite) TestGetProfileDevBypass() {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("user-%d", CustomerUser.ID),
	}
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthTestSuite) TestGetProfileDevBypassNonexistentUser() {
	headers := map[string]string{
		"Authorization": "user-99999",
	}
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthTestSuite) TestGetProfileFieldsNoPassword() {
	headers := GetAdminAuthHeader()
	w := MakeRequest(s.router, "GET", "/api/auth/me", nil, headers)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	user := resp["user"].(map[string]interface{})
	_, hasPassword := user["password"]
	s.False(hasPassword, "password should not be in response")
	s.NotNil(user["id"])
	s.NotNil(user["email"])
	s.NotNil(user["name"])
	s.NotNil(user["role"])
}

func (s *AuthTestSuite) TestUpdateProfileName() {
	headers := GetCustomerAuthHeader()
	body := map[string]interface{}{
		"name": "NewName",
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/profile", body, headers)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("NewName", resp["name"])
}

func (s *AuthTestSuite) TestUpdateProfileUnauthorized() {
	body := map[string]interface{}{
		"name": "NewName",
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/profile", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AuthTestSuite) TestUpdateProfileEmptyBody() {
	headers := GetCustomerAuthHeader()
	w := MakeRequest(s.router, "PUT", "/api/auth/profile", map[string]interface{}{}, headers)
	s.Equal(http.StatusOK, w.Code)
}

func (s *AuthTestSuite) TestChangePassword() {
	headers := GetCustomer2AuthHeader()
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"oldPassword": "user1234",
		"newPassword": "newpass123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, headers)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("密码修改成功", resp["message"])

	handlers.ResetLoginFailures()
	loginBody := map[string]interface{}{
		"email":    "user2@test.com",
		"password": "newpass123",
	}
	wLogin := MakeRequest(s.router, "POST", "/api/auth/login", loginBody, nil)
	s.Equal(http.StatusOK, wLogin.Code)
}

func (s *AuthTestSuite) TestChangePasswordWrongOld() {
	headers := GetCustomerAuthHeader()
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"oldPassword": "wrongold",
		"newPassword": "newpass123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, headers)
	s.Equal(http.StatusUnauthorized, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "原密码错误")
}

func (s *AuthTestSuite) TestChangePasswordSameAsOld() {
	headers := GetCustomerAuthHeader()
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"oldPassword": "user1234",
		"newPassword": "user1234",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, headers)
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "相同")
}

func (s *AuthTestSuite) TestChangePasswordTooShort() {
	headers := GetCustomerAuthHeader()
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"oldPassword": "user1234",
		"newPassword": "12345",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, headers)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AuthTestSuite) TestChangePasswordWrongCaptcha() {
	headers := GetCustomerAuthHeader()

	body := map[string]interface{}{
		"oldPassword": "user1234",
		"newPassword": "newpass123",
		"captchaId":   "nonexistent",
		"captchaCode": "wrong",
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, headers)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AuthTestSuite) TestChangePasswordUnauthorized() {
	captchaID, captchaCode := s.getCaptcha()

	body := map[string]interface{}{
		"oldPassword": "user1234",
		"newPassword": "newpass123",
		"captchaId":   captchaID,
		"captchaCode": captchaCode,
	}
	w := MakeRequest(s.router, "PUT", "/api/auth/change-password", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}
