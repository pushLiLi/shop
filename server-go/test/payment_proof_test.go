package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type PaymentProofTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *PaymentProofTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestPaymentProofSuite(t *testing.T) {
	suite.Run(t, new(PaymentProofTestSuite))
}

func (s *PaymentProofTestSuite) createPendingOrderWithPaymentMethod() (uint, uint) {
	method := models.PaymentMethod{Name: "测试支付", IsActive: true}
	database.DB.Create(&method)

	addrID := Data.Addresses[0].ID
	body := map[string]interface{}{
		"addressId": addrID,
		"remark":    "payment proof test",
	}
	w := MakeRequest(s.router, "POST", "/api/orders", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	orderID := uint(resp["orderId"].(float64))

	return orderID, method.ID
}

func (s *PaymentProofTestSuite) TestUploadPaymentProof() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofNoFile() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		nil, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofInvalidExt() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".exe")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofTooLarge() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateLargeTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofNonPendingOrder() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/orders/%d/status", orderID),
		map[string]interface{}{"status": models.OrderStatusProcessing}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	imgFile := CreateTestImageFile(".jpg")
	w = MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofOtherUserOrder() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomer2AuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofInactiveMethod() {
	orderID, _ := s.createPendingOrderWithPaymentMethod()

	inactiveMethod := models.PaymentMethod{Name: "已停用", IsActive: false}
	database.DB.Create(&inactiveMethod)

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", inactiveMethod.ID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestUploadPaymentProofReplacesExisting() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	imgFile2 := CreateTestImageFile(".png")
	w = MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile2},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofApprove() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := uploadResp["paymentProof"].(map[string]interface{})["id"].(float64)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "approve"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d", orderID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var orderResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &orderResp)
	s.Equal(models.OrderStatusProcessing, orderResp["order"].(map[string]interface{})["status"])
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofReject() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := uploadResp["paymentProof"].(map[string]interface{})["id"].(float64)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "reject", "rejectReason": "图片不清晰"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/notifications", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var notifResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notifResp)
	notifs := notifResp["notifications"].([]interface{})
	s.True(len(notifs) >= 1, "should have rejection notification")
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofAlreadyReviewed() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := uploadResp["paymentProof"].(map[string]interface{})["id"].(float64)

	MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "approve"}, GetAdminAuthHeader())

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "reject"}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofUnauthorized() {
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/1/review",
		map[string]interface{}{"action": "approve"}, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofCustomerForbidden() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := uploadResp["paymentProof"].(map[string]interface{})["id"].(float64)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "approve"}, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *PaymentProofTestSuite) TestReviewPaymentProofServiceForbidden() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile},
		GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := uploadResp["paymentProof"].(map[string]interface{})["id"].(float64)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", int(proofID)),
		map[string]interface{}{"action": "approve"}, GetServiceAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *PaymentProofTestSuite) TestBatchReviewPaymentProofsApprove() {
	orderID1, methodID := s.createPendingOrderWithPaymentMethod()
	orderID2, _ := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID1),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	var resp1 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp1)
	proofID1 := int(resp1["paymentProof"].(map[string]interface{})["id"].(float64))

	imgFile2 := CreateTestImageFile(".jpg")
	w = MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID2),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile2}, GetCustomerAuthHeader())
	var resp2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp2)
	proofID2 := int(resp2["paymentProof"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/batch-review",
		map[string]interface{}{"ids": []int{proofID1, proofID2}, "action": "approve"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var batchResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &batchResp)
	s.Equal(float64(2), batchResp["reviewed"])
}

func (s *PaymentProofTestSuite) TestBatchReviewPaymentProofsReject() {
	orderID1, methodID := s.createPendingOrderWithPaymentMethod()
	orderID2, _ := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID1),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	var resp1 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp1)
	proofID1 := int(resp1["paymentProof"].(map[string]interface{})["id"].(float64))

	imgFile2 := CreateTestImageFile(".jpg")
	w = MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID2),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile2}, GetCustomerAuthHeader())
	var resp2 map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp2)
	proofID2 := int(resp2["paymentProof"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/batch-review",
		map[string]interface{}{"ids": []int{proofID1, proofID2}, "action": "reject", "rejectReason": "批量驳回"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var batchResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &batchResp)
	s.Equal(float64(2), batchResp["reviewed"])
}

func (s *PaymentProofTestSuite) TestBatchReviewPaymentProofsEmpty() {
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/batch-review",
		map[string]interface{}{"ids": []int{}, "action": "approve"}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestBatchReviewPaymentProofsInvalidAction() {
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/batch-review",
		map[string]interface{}{"ids": []int{1}, "action": "invalid"}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestBatchReviewPaymentProofsNoPending() {
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/batch-review",
		map[string]interface{}{"ids": []int{99999}, "action": "approve"}, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentProofTestSuite) TestGetOrderPaymentProof() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["paymentProof"])
}

func (s *PaymentProofTestSuite) TestApproveProofTriggersNotification() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := int(uploadResp["paymentProof"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID),
		map[string]interface{}{"action": "approve"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/notifications", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var notifResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notifResp)
	notifs := notifResp["notifications"].([]interface{})
	s.True(len(notifs) >= 1)
}

func (s *PaymentProofTestSuite) TestRejectProofTriggersNotification() {
	orderID, methodID := s.createPendingOrderWithPaymentMethod()

	imgFile := CreateTestImageFile(".jpg")
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID),
		map[string]string{"paymentMethodId": fmt.Sprintf("%d", methodID)},
		map[string]string{"file": imgFile}, GetCustomerAuthHeader())
	var uploadResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &uploadResp)
	proofID := int(uploadResp["paymentProof"].(map[string]interface{})["id"].(float64))

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID),
		map[string]interface{}{"action": "reject", "rejectReason": "截图不清晰"}, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/notifications", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var notifResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &notifResp)
	notifs := notifResp["notifications"].([]interface{})
	s.True(len(notifs) >= 1)
}
