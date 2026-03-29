package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ProductTestSuite struct {
	suite.Suite
	router    *gin.Engine
	prodCount int
}

func (s *ProductTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

func (s *ProductTestSuite) TestGetProductsDefaultPagination() {
	w := MakeRequest(s.router, "GET", "/api/products", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	products := resp["products"].([]interface{})
	s.True(len(products) <= 12)
	s.Equal(float64(1), resp["page"])
	s.Equal(float64(12), resp["limit"])
	s.True(resp["total"].(float64) > 0)
}

func (s *ProductTestSuite) TestGetProductsCustomPagination() {
	w := MakeRequest(s.router, "GET", "/api/products?page=2&limit=3", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(2), resp["page"])
	s.Equal(float64(3), resp["limit"])

	products := resp["products"].([]interface{})
	s.True(len(products) <= 3)
}

func (s *ProductTestSuite) TestGetProductsSearch() {
	w := MakeRequest(s.router, "GET", "/api/products?search=高希霸", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestGetProductsFilterByCategorySlug() {
	w := MakeRequest(s.router, "GET", "/api/products?category=cohiba-classic", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestGetProductsFilterByCategoryId() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/products?categoryId=%d", classicID), nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestGetProductsFeaturedOnly() {
	w := MakeRequest(s.router, "GET", "/api/products?featured=true", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	for _, p := range products {
		product := p.(map[string]interface{})
		s.Equal(true, product["isFeatured"])
	}
}

func (s *ProductTestSuite) TestGetProductsSortByPrice() {
	w := MakeRequest(s.router, "GET", "/api/products?sortBy=price&sortOrder=asc&limit=100", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	if len(products) >= 2 {
		first := products[0].(map[string]interface{})
		second := products[1].(map[string]interface{})
		s.True(first["price"].(float64) <= second["price"].(float64))
	}
}

func (s *ProductTestSuite) TestGetProductsSortByName() {
	w := MakeRequest(s.router, "GET", "/api/products?sortBy=name&sortOrder=desc&limit=100", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestGetProductsSortAliases() {
	w := MakeRequest(s.router, "GET", "/api/products?sort=price&order=asc&limit=100", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.True(resp["total"].(float64) > 0)
}

func (s *ProductTestSuite) TestGetProductsStockPriority() {
	w := MakeRequest(s.router, "GET", "/api/products?limit=100", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})

	foundZeroStock := false
	for _, p := range products {
		product := p.(map[string]interface{})
		if product["stock"].(float64) == 0 {
			foundZeroStock = true
		} else if foundZeroStock {
			s.Fail("out-of-stock product should come after in-stock products")
		}
	}
}

func (s *ProductTestSuite) TestGetProductsOnlyActive() {
	w := MakeRequest(s.router, "GET", "/api/products?limit=100", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	for _, p := range products {
		product := p.(map[string]interface{})
		s.Equal(true, product["isActive"])
	}
}

func (s *ProductTestSuite) TestGetProductsInvalidPagination() {
	w := MakeRequest(s.router, "GET", "/api/products?page=-1&limit=0", nil, nil)
	s.Equal(http.StatusOK, w.Code)
}

func (s *ProductTestSuite) TestGetProductsEmptyResult() {
	w := MakeRequest(s.router, "GET", "/api/products?search=nonexistent_product_xyz", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(0), resp["total"])
	products := resp["products"].([]interface{})
	s.Equal(0, len(products))
}

func (s *ProductTestSuite) TestGetProductsIncludesCategory() {
	w := MakeRequest(s.router, "GET", "/api/products?limit=1", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	if len(products) > 0 {
		product := products[0].(map[string]interface{})
		s.NotNil(product["category"])
	}
}

func (s *ProductTestSuite) TestGetProductExisting() {
	product := Data.Products[0]
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/products/%d", product.ID), nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(product.Name, resp["name"])
	s.NotNil(resp["category"])
}

func (s *ProductTestSuite) TestGetProductNotFound() {
	w := MakeRequest(s.router, "GET", "/api/products/99999", nil, nil)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *ProductTestSuite) TestGetProductNonNumericID() {
	w := MakeRequest(s.router, "GET", "/api/products/abc", nil, nil)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *ProductTestSuite) TestGetProductInactiveAccessible() {
	inactiveProduct := Data.Products[3]
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/products/%d", inactiveProduct.ID), nil, nil)
	s.Equal(http.StatusOK, w.Code)
}

func (s *ProductTestSuite) TestAdminGetProductsAll() {
	w := MakeRequest(s.router, "GET", "/api/admin/products", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})

	hasInactive := false
	for _, p := range products {
		product := p.(map[string]interface{})
		if val, ok := product["isActive"]; ok {
			if boolVal, ok2 := val.(bool); ok2 && !boolVal {
				hasInactive = true
				break
			}
		}
	}
	s.True(hasInactive, "admin should see inactive products")
}

func (s *ProductTestSuite) TestAdminGetProductsSearch() {
	w := MakeRequest(s.router, "GET", "/api/admin/products?search=蒙特", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestAdminGetProductsFilterActive() {
	w := MakeRequest(s.router, "GET", "/api/admin/products?active=false", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0, "should have at least one inactive product")
	for _, p := range products {
		product := p.(map[string]interface{})
		s.Equal(false, product["isActive"])
	}
}

func (s *ProductTestSuite) TestAdminGetProductsFilterFeatured() {
	w := MakeRequest(s.router, "GET", "/api/admin/products?featured=true", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	for _, p := range products {
		product := p.(map[string]interface{})
		s.Equal(true, product["isFeatured"])
	}
}

func (s *ProductTestSuite) TestAdminGetProductsFilterByCategory() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/admin/products?categoryId=%d", classicID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	products := resp["products"].([]interface{})
	s.True(len(products) > 0)
}

func (s *ProductTestSuite) TestAdminGetProductsPagination() {
	w := MakeRequest(s.router, "GET", "/api/admin/products?page=1&limit=5", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(float64(1), resp["page"])
	s.Equal(float64(5), resp["limit"])
	s.True(resp["total"].(float64) > 0)
}

func (s *ProductTestSuite) TestAdminGetProductsUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/admin/products", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *ProductTestSuite) TestAdminGetProductsForbidden() {
	w := MakeRequest(s.router, "GET", "/api/admin/products", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *ProductTestSuite) TestAdminCreateProduct() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	body := map[string]interface{}{
		"name":       "test product",
		"price":      99.99,
		"categoryId": classicID,
		"stock":      10,
		"isActive":   true,
		"isFeatured": false,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("test product", resp["name"])
	s.Equal("test-product", resp["slug"])
}

func (s *ProductTestSuite) TestAdminCreateProductCustomSlug() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	body := map[string]interface{}{
		"name":       "自定义Slug商品",
		"slug":       "custom-slug-test",
		"price":      88.88,
		"categoryId": classicID,
		"stock":      5,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", body, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("custom-slug-test", resp["slug"])
}

func (s *ProductTestSuite) TestAdminCreateProductMissingFields() {
	w := MakeRequest(s.router, "POST", "/api/admin/products", map[string]interface{}{}, GetAdminAuthHeader())
	s.True(w.Code == http.StatusBadRequest || w.Code == http.StatusInternalServerError)
}

func (s *ProductTestSuite) TestAdminCreateProductForbidden() {
	body := map[string]interface{}{
		"name":  "test",
		"price": 10,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", body, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *ProductTestSuite) createTestProduct() int {
	s.prodCount++
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	body := map[string]interface{}{
		"name":       fmt.Sprintf("test product %d", s.prodCount),
		"price":      88.88,
		"categoryId": classicID,
		"stock":      10,
		"isActive":   true,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", body, GetAdminAuthHeader())
	s.Require().Equal(http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return int(resp["id"].(float64))
}

func (s *ProductTestSuite) TestAdminUpdateProduct() {
	productID := s.createTestProduct()
	body := map[string]interface{}{
		"name":        "更新后名称",
		"slug":        "updated-slug",
		"description": "更新描述",
		"price":       199.99,
		"imageUrl":    "/test/updated.jpg",
		"images":      "",
		"categoryId":  findCategoryID(Data.Categories, "cohiba-classic"),
		"stock":       55,
		"isActive":    true,
		"isFeatured":  true,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/products/%d", productID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal("更新后名称", resp["name"])
	s.Equal("updated-slug", resp["slug"])
}

func (s *ProductTestSuite) TestAdminUpdateProductPartial() {
	productID := s.createTestProduct()
	body := map[string]interface{}{
		"name":        "partial update product",
		"description": "original desc",
		"price":       999.99,
		"imageUrl":    "/test/original.jpg",
		"images":      "",
		"categoryId":  findCategoryID(Data.Categories, "cohiba-classic"),
		"stock":       10,
		"isActive":    true,
		"isFeatured":  false,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/products/%d", productID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(999.99, resp["price"])
}

func (s *ProductTestSuite) TestAdminUpdateProductEmptySlugPreserved() {
	productID := s.createTestProduct()

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/products/%d", productID), nil, nil)
	s.Equal(http.StatusOK, w.Code)
	var origResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &origResp)
	originalSlug := origResp["slug"].(string)

	body := map[string]interface{}{
		"name":        "名称不改slug",
		"slug":        "",
		"description": "some description",
		"price":       99.99,
		"imageUrl":    "/test/img.jpg",
		"images":      "",
		"categoryId":  findCategoryID(Data.Categories, "cohiba-classic"),
		"stock":       10,
		"isActive":    true,
		"isFeatured":  false,
	}
	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/products/%d", productID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(originalSlug, resp["slug"])
}

func (s *ProductTestSuite) TestAdminUpdateProductNotFound() {
	body := map[string]interface{}{
		"name": "test",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/products/99999", body, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *ProductTestSuite) TestAdminDeleteProduct() {
	classicID := findCategoryID(Data.Categories, "cohiba-classic")
	newProduct := map[string]interface{}{
		"name":       "待删除商品",
		"price":      50,
		"categoryId": classicID,
		"stock":      1,
		"isActive":   true,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/products", newProduct, GetAdminAuthHeader())
	s.Equal(http.StatusCreated, w.Code)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	createdID := int(created["id"].(float64))

	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/products/%d", createdID), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", fmt.Sprintf("/api/products/%d", createdID), nil, nil)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *ProductTestSuite) TestAdminDeleteProductNotFound() {
	w := MakeRequest(s.router, "DELETE", "/api/admin/products/99999", nil, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func findCategoryID(categories []models.Category, slug string) uint {
	for _, c := range categories {
		if c.Slug == slug {
			return c.ID
		}
	}
	return 0
}
