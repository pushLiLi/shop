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

type OrderStateTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *OrderStateTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func (s *OrderStateTestSuite) createPendingOrder() (uint, uint) {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var cartResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cartResp)
	items := cartResp["items"].([]interface{})

	var productID uint
	for _, item := range items {
		it := item.(map[string]interface{})
		productID = uint(it["productId"].(float64))
		break
	}

	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"addressId": addrID,
		"remark":    "state test order",
	}
	w = MakeRequest(s.router, "POST", "/api/orders", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	orderID := uint(resp["orderId"].(float64))

	return orderID, productID
}

func TestOrderStateSuite(t *testing.T) {
	suite.Run(t, new(OrderStateTestSuite))
}

func (s *OrderStateTestSuite) TestPendingToProcessing() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": models.OrderStatusProcessing,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(models.OrderStatusProcessing, resp["status"])
}

func (s *OrderStateTestSuite) TestPendingToCancelled() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": models.OrderStatusCancelled,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(models.OrderStatusCancelled, resp["status"])
}

func (s *OrderStateTestSuite) TestPendingToShipped_illegal() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status":           models.OrderStatusShipped,
		"trackingCompany":  "SF",
		"trackingNumber":   "SF123456",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestPendingToCompleted_illegal() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": models.OrderStatusCompleted,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestProcessingToShipped() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany":  "SF Express",
			"trackingNumber":   "SF123456789",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(models.OrderStatusShipped, resp["status"])
}

func (s *OrderStateTestSuite) TestProcessingToCancelled() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCancelled}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(models.OrderStatusCancelled, resp["status"])
}

func (s *OrderStateTestSuite) TestProcessingToPending_illegal() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusPending}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestShippedToCompleted() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany":  "SF",
			"trackingNumber":   "SF123",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCompleted}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(models.OrderStatusCompleted, resp["status"])
}

func (s *OrderStateTestSuite) TestShippedToCancelled_illegal() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany":  "SF",
			"trackingNumber":   "SF123",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCancelled}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestShipOrderMissingTracking() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":          models.OrderStatusShipped,
			"trackingCompany": "",
			"trackingNumber":  "SF123",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany":  "SF",
			"trackingNumber":   "",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestCompletedTerminal() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany": "SF",
			"trackingNumber":  "SF123",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCompleted}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCancelled}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusPending}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestCancelledTerminal() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusCancelled}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusPending}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestInvalidStatusTransition() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": "invalid_status",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *OrderStateTestSuite) TestUpdateStatusUnauthorized() {
	body := map[string]interface{}{
		"status": models.OrderStatusProcessing,
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/orders/1/status", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *OrderStateTestSuite) TestUpdateStatusCustomerForbidden() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": models.OrderStatusProcessing,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *OrderStateTestSuite) TestServiceCanUpdateStatus() {
	orderID, _ := s.createPendingOrder()

	body := map[string]interface{}{
		"status": models.OrderStatusProcessing,
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID), body, GetServiceAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *OrderStateTestSuite) TestOrderStatusChangeTriggersNotification() {
	orderID, _ := s.createPendingOrder()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{
			"status":           models.OrderStatusShipped,
			"trackingCompany":  "SF",
			"trackingNumber":   "SF123",
		}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/notifications", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	notifs := resp["notifications"].([]interface{})
	s.True(len(notifs) >= 1, "should have at least one notification")
}
