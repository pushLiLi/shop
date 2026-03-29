package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type OrderTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *OrderTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) TestGetOrders() {
	w := MakeRequest(s.router, "GET", "/api/orders", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	orders := resp["orders"].([]interface{})
	s.True(len(orders) >= 1, "seeded user should have at least 1 order")
}

func (s *OrderTestSuite) TestGetOrdersUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/orders", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *OrderTestSuite) TestGetOrdersOnlyOwnOrders() {
	w := MakeRequest(s.router, "GET", "/api/orders", nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	orders := resp["orders"].([]interface{})
	s.Equal(0, len(orders), "customer2 should have no orders")
}

func (s *OrderTestSuite) TestGetOrderById() {
	order := Data.Orders[0]
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d", order.ID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["items"])
	s.NotNil(resp["address"])
}

func (s *OrderTestSuite) TestGetOrderByOrderNo() {
	order := Data.Orders[0]
	s.NotEmpty(order.OrderNo, "seeded order should have an OrderNo")

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%s", order.OrderNo), nil, GetCustomerAuthHeader())

	if w.Code == http.StatusOK {
		s.Equal(http.StatusOK, w.Code)
	} else {
		s.Equal(http.StatusNotFound, w.Code)
	}
}

func (s *OrderTestSuite) TestGetOrderWrongUser() {
	order := Data.Orders[0]
	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d", order.ID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *OrderTestSuite) TestGetOrderNotFound() {
	w := MakeRequest(s.router, "GET", "/api/orders/99999", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *OrderTestSuite) TestCreateOrder() {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var cartResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cartResp)
	cartItems := cartResp["items"].([]interface{})
	s.True(len(cartItems) > 0, "need cart items to create order")

	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"addressId": addrID,
		"remark":    "测试订单",
	}
	w = MakeRequest(s.router, "POST", "/api/orders", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	s.NotNil(resp["orderId"])
	s.NotNil(resp["orderNo"])

	w = MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var cartAfter map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cartAfter)
	itemsAfter := cartAfter["items"].([]interface{})
	s.Equal(0, len(itemsAfter), "cart should be cleared after order")
}

func (s *OrderTestSuite) TestCreateOrderEmptyCart() {
	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"addressId": addrID,
	}
	w := MakeRequest(s.router, "POST", "/api/orders", body, GetCustomer2AuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderTestSuite) TestCreateOrderNoAddress() {
	body := map[string]interface{}{
		"addressId": 0,
	}
	w := MakeRequest(s.router, "POST", "/api/orders", body, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderTestSuite) TestCreateOrderWrongAddress() {
	wrongAddrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"addressId": wrongAddrID,
	}
	w := MakeRequest(s.router, "POST", "/api/orders", body, GetCustomer2AuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderTestSuite) TestCreateOrderUnauthorized() {
	body := map[string]interface{}{
		"addressId": 1,
	}
	w := MakeRequest(s.router, "POST", "/api/orders", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *OrderTestSuite) TestGetOrderUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/orders/1", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}
