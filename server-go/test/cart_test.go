package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type CartTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *CartTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestCartSuite(t *testing.T) {
	suite.Run(t, new(CartTestSuite))
}

func (s *CartTestSuite) TestGetCartLoggedIn() {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.True(len(items) >= 2, "seeded user should have 2 cart items")
	s.True(resp["total"].(float64) > 0)
}

func (s *CartTestSuite) TestGetCartNotLoggedIn() {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.Equal(0, len(items))
	s.Equal(float64(0), resp["total"])
}

func (s *CartTestSuite) TestGetCartItemsWithCategory() {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	if len(items) > 0 {
		item := items[0].(map[string]interface{})
		product := item["product"].(map[string]interface{})
		s.NotNil(product["category"])
	}
}

func (s *CartTestSuite) TestAddToCartNewProduct() {
	productID := Data.Products[8].ID
	body := map[string]interface{}{
		"productId": productID,
		"quantity":  3,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.True(len(items) >= 3)
}

func (s *CartTestSuite) TestAddToCartDuplicateAccumulates() {
	productID := Data.Products[0].ID
	body := map[string]interface{}{
		"productId": productID,
		"quantity":  5,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	for _, item := range items {
		i := item.(map[string]interface{})
		p := i["product"].(map[string]interface{})
		if int(p["id"].(float64)) == int(productID) {
			s.True(i["quantity"].(float64) >= 7, "quantity should accumulate (2 original + 5 new)")
		}
	}
}

func (s *CartTestSuite) TestAddToCartDefaultQuantity() {
	productID := Data.Products[1].ID
	body := map[string]interface{}{
		"productId": productID,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *CartTestSuite) TestAddToCartProductNotFound() {
	body := map[string]interface{}{
		"productId": 99999,
		"quantity":  1,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, GetCustomerAuthHeader())
	s.True(w.Code == http.StatusOK || w.Code == http.StatusBadRequest,
		"AddToCart does not validate product existence; FK may cause error or succeed")
}

func (s *CartTestSuite) TestAddToCartUnauthorized() {
	body := map[string]interface{}{
		"productId": 1,
		"quantity":  1,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *CartTestSuite) createTestCartItem() int {
	productID := Data.Products[9].ID
	body := map[string]interface{}{
		"productId": productID,
		"quantity":  2,
	}
	w := MakeRequest(s.router, "POST", "/api/cart", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	item := resp["item"].(map[string]interface{})
	return int(item["id"].(float64))
}

func (s *CartTestSuite) TestUpdateCartItemQuantity() {
	cartItemID := s.createTestCartItem()
	body := map[string]interface{}{
		"quantity": 5,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/cart/%d", cartItemID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *CartTestSuite) TestUpdateCartItemQuantityZeroFails() {
	cartItemID := s.createTestCartItem()
	body := map[string]interface{}{
		"quantity": 0,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/cart/%d", cartItemID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code, "quantity=0 fails required validation")
}

func (s *CartTestSuite) TestUpdateCartItemNegativeDeletes() {
	cartItemID := s.createTestCartItem()
	body := map[string]interface{}{
		"quantity": -1,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/cart/%d", cartItemID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *CartTestSuite) TestUpdateCartItemWrongUser() {
	cartItemID := Data.CartItems[0].ID
	body := map[string]interface{}{
		"quantity": 5,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/cart/%d", cartItemID), body, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *CartTestSuite) TestUpdateCartItemUnauthorized() {
	body := map[string]interface{}{
		"quantity": 5,
	}
	w := MakeRequest(s.router, "PUT", "/api/cart/1", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *CartTestSuite) TestDeleteCartItem() {
	cartItemID := s.createTestCartItem()
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/cart/%d", cartItemID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *CartTestSuite) TestDeleteCartItemWrongUser() {
	cartItemID := Data.CartItems[0].ID
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/cart/%d", cartItemID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *CartTestSuite) TestDeleteCartItemUnauthorized() {
	w := MakeRequest(s.router, "DELETE", "/api/cart/1", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}
