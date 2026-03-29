package test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UploadTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *UploadTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestUploadSuite(t *testing.T) {
	suite.Run(t, new(UploadTestSuite))
}

func (s *UploadTestSuite) TestUploadImageJPG() {
	filePath := CreateTestImageFile(".jpg")
	defer os.Remove(filePath)

	files := map[string]string{"file": filePath}
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, files, GetAdminAuthHeader())

	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	s.Contains(resp["url"], "/media/")
}

func (s *UploadTestSuite) TestUploadImagePNG() {
	filePath := CreateTestImageFile(".png")
	defer os.Remove(filePath)

	files := map[string]string{"file": filePath}
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, files, GetAdminAuthHeader())

	s.Equal(http.StatusOK, w.Code)
}

func (s *UploadTestSuite) TestUploadNonImageFile() {
	filePath := CreateTestImageFile(".txt")
	defer os.Remove(filePath)

	files := map[string]string{"file": filePath}
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, files, GetAdminAuthHeader())

	s.Equal(http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "jpg")
}

func (s *UploadTestSuite) TestUploadFileTooLarge() {
	filePath := CreateLargeTestImageFile(".jpg")
	defer os.Remove(filePath)

	files := map[string]string{"file": filePath}
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, files, GetAdminAuthHeader())

	s.Equal(http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "10MB")
}

func (s *UploadTestSuite) TestUploadNoFile() {
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, nil, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "文件")
}

func (s *UploadTestSuite) TestUploadForbidden() {
	filePath := CreateTestImageFile(".jpg")
	defer os.Remove(filePath)

	files := map[string]string{"file": filePath}
	w := MakeFormRequest(s.router, "POST", "/api/admin/upload", nil, files, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}
