package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AddressTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *AddressTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (s *AddressTestSuite) TestGetAddresses() {
	w := MakeRequest(s.router, "GET", "/api/addresses", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	addresses := resp["addresses"].([]interface{})
	s.True(len(addresses) >= 2, "seeded user should have at least 2 addresses")
}

func (s *AddressTestSuite) TestGetAddressesDefaultFirst() {
	w := MakeRequest(s.router, "GET", "/api/addresses", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	addresses := resp["addresses"].([]interface{})
	if len(addresses) > 0 {
		first := addresses[0].(map[string]interface{})
		s.Equal(true, first["isDefault"])
	}
}

func (s *AddressTestSuite) TestGetAddressesUnauthorized() {
	w := MakeRequest(s.router, "GET", "/api/addresses", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AddressTestSuite) TestCreateAddress() {
	body := map[string]interface{}{
		"fullName":     "李四",
		"addressLine1": "789 Oak St",
		"addressLine2": "Apt 5",
		"city":         "Guangzhou",
		"state":        "Guangdong",
		"zipCode":      "510000",
		"phone":        "13900139001",
		"isDefault":    false,
	}
	w := MakeRequest(s.router, "POST", "/api/addresses", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	address := resp["address"].(map[string]interface{})
	s.Equal("李四", address["fullName"])
}

func (s *AddressTestSuite) TestCreateAddressMissingFields() {
	body := map[string]interface{}{
		"fullName": "Test",
	}
	w := MakeRequest(s.router, "POST", "/api/addresses", body, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *AddressTestSuite) TestCreateAddressSetsDefault() {
	body := map[string]interface{}{
		"fullName":     "王五",
		"addressLine1": "321 Pine St",
		"city":         "Shenzhen",
		"state":        "Guangdong",
		"zipCode":      "518000",
		"phone":        "13700137001",
		"isDefault":    true,
	}
	w := MakeRequest(s.router, "POST", "/api/addresses", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/addresses", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	addresses := resp["addresses"].([]interface{})
	defaultCount := 0
	for _, a := range addresses {
		addr := a.(map[string]interface{})
		if addr["isDefault"] == true {
			defaultCount++
		}
	}
	s.Equal(1, defaultCount, "only one default address should exist")
}

func (s *AddressTestSuite) TestCreateAddressUnauthorized() {
	body := map[string]interface{}{
		"fullName": "Test",
	}
	w := MakeRequest(s.router, "POST", "/api/addresses", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *AddressTestSuite) TestUpdateAddress() {
	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"fullName":     "张三更新",
		"addressLine1": "999 Updated St",
		"city":         "Beijing",
		"state":        "Beijing",
		"zipCode":      "100001",
		"phone":        "13800138000",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/addresses/%d", addrID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	address := resp["address"].(map[string]interface{})
	s.Equal("张三更新", address["fullName"])
}

func (s *AddressTestSuite) TestUpdateAddressWrongUser() {
	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"fullName": "Hacker",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/addresses/%d", addrID), body, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *AddressTestSuite) TestUpdateAddressNotFound() {
	body := map[string]interface{}{
		"fullName": "Test",
	}
	w := MakeRequest(s.router, "PUT", "/api/addresses/99999", body, GetCustomerAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *AddressTestSuite) TestSetDefaultAddress() {
	addrID := Data.Addresses[1].ID
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/addresses/%d/default", addrID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/addresses", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	addresses := resp["addresses"].([]interface{})
	for _, a := range addresses {
		addr := a.(map[string]interface{})
		if int(addr["id"].(float64)) == int(addrID) {
			s.Equal(true, addr["isDefault"])
		} else {
			s.Equal(false, addr["isDefault"])
		}
	}
}

func (s *AddressTestSuite) TestSetDefaultAddressWrongUser() {
	addrID := Data.Addresses[0].ID
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/addresses/%d/default", addrID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *AddressTestSuite) TestDeleteAddress() {
	body := map[string]interface{}{
		"fullName":     "删除测试",
		"addressLine1": "000 Delete St",
		"city":         "Test",
		"state":        "Test",
		"zipCode":      "000000",
		"phone":        "13000130001",
		"isDefault":    false,
	}
	w := MakeRequest(s.router, "POST", "/api/addresses", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	address := created["address"].(map[string]interface{})
	newAddrID := int(address["id"].(float64))

	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/addresses/%d", newAddrID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *AddressTestSuite) TestDeleteAddressWrongUser() {
	addrID := Data.Addresses[0].ID
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/addresses/%d", addrID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *AddressTestSuite) TestDeleteAddressUnauthorized() {
	w := MakeRequest(s.router, "DELETE", "/api/addresses/1", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}
