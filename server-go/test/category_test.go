package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type CategoryTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *CategoryTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}

func (s *CategoryTestSuite) TestGetCategoriesTree() {
	w := MakeRequest(s.router, "GET", "/api/categories", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(len(resp) >= 3, "should have at least 3 top-level categories")
}

func (s *CategoryTestSuite) TestGetCategoriesChildren() {
	w := MakeRequest(s.router, "GET", "/api/categories", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	for _, cat := range resp {
		name := cat["name"].(string)
		if name == "古巴雪茄" {
			children := cat["children"].([]interface{})
			s.True(len(children) >= 2, "古巴雪茄 should have at least 2 children")
		}
	}
}

func (s *CategoryTestSuite) TestGetCategoriesProductCount() {
	w := MakeRequest(s.router, "GET", "/api/categories", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	for _, cat := range resp {
		s.NotNil(cat["_count"], "each category should have _count")
	}
}

func (s *CategoryTestSuite) TestAdminCreateTopLevelCategory() {
	body := map[string]interface{}{
		"name": "新顶级分类",
		"slug": "new-top-cat",
	}
	w := MakeRequest(s.router, "POST", "/api/admin/categories", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("新顶级分类", resp["name"])
	s.Equal("new-top-cat", resp["slug"])
}

func (s *CategoryTestSuite) TestAdminCreateChildCategory() {
	parentID := findCategoryID(Data.Categories, "cohiba")
	body := map[string]interface{}{
		"name":     "子分类",
		"slug":     "child-cat",
		"parentId": parentID,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/categories", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(parentID), resp["parentId"])
}

func (s *CategoryTestSuite) TestAdminCreateCategoryAutoSlug() {
	body := map[string]interface{}{
		"name": "Auto Slug Category",
	}
	w := MakeRequest(s.router, "POST", "/api/admin/categories", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("auto-slug-category", resp["slug"])
}

func (s *CategoryTestSuite) TestAdminUpdateCategory() {
	catID := findCategoryID(Data.Categories, "accessories")
	body := map[string]interface{}{
		"name": "配件更新",
		"slug": "accessories-updated",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/categories/%d", catID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("配件更新", resp["name"])
	s.Equal("accessories-updated", resp["slug"])
}

func (s *CategoryTestSuite) TestAdminDeleteCategoryNoDependencies() {
	body := map[string]interface{}{
		"name": "待删除分类",
		"slug": "to-delete-cat",
	}
	w := MakeRequest(s.router, "POST", "/api/admin/categories", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	catID := int(created["id"].(float64))

	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/categories/%d", catID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *CategoryTestSuite) TestAdminDeleteCategoryWithProducts() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/categories/%d", classicID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "商品")
}

func (s *CategoryTestSuite) TestAdminDeleteCategoryWithChildren() {
	cohibaID := findCategoryID(Data.Categories, "cohiba")
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/categories/%d", cohibaID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Contains(resp["error"], "子分类")
}

func (s *CategoryTestSuite) TestAdminDeleteCategoryNotFound() {
	w := MakeRequest(s.router, "DELETE", "/api/admin/categories/99999", nil, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *CategoryTestSuite) TestAdminCategoryForbidden() {
	w := MakeRequest(s.router, "POST", "/api/admin/categories", map[string]interface{}{"name": "test"}, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}
