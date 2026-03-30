package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type PaymentTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *PaymentTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestPaymentSuite(t *testing.T) {
	suite.Run(t, new(PaymentTestSuite))
}

func (s *PaymentTestSuite) createPaymentMethod(name, qrCode string) uint {
	body := map[string]interface{}{
		"name":         name,
		"qrCodeUrl":    qrCode,
		"instructions": "请扫码支付",
		"isActive":     true,
		"sortOrder":    0,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	pm := resp["paymentMethod"].(map[string]interface{})
	return uint(pm["id"].(float64))
}

func (s *PaymentTestSuite) createPendingOrder() uint {
	w := MakeRequest(s.router, "GET", "/api/cart", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var cartResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &cartResp)
	cartItems := cartResp["items"].([]interface{})
	s.True(len(cartItems) > 0, "need cart items")

	w = MakeRequest(s.router, "POST", "/api/orders", map[string]interface{}{
		"addressId": Data.Addresses[0].ID,
	}, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var orderResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &orderResp)
	return uint(orderResp["orderId"].(float64))
}

func (s *PaymentTestSuite) uploadProof(orderID, pmID uint) uint {
	imagePath := CreateTestImageFile(".png")
	defer os.Remove(imagePath)

	formData := map[string]string{
		"paymentMethodId": fmt.Sprintf("%d", pmID),
	}
	fileFields := map[string]string{
		"file": imagePath,
	}
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), formData, fileFields, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	proof := resp["paymentProof"].(map[string]interface{})
	return uint(proof["id"].(float64))
}

func (s *PaymentTestSuite) TestGetPublicPaymentMethodsEmpty() {
	w := MakeRequest(s.router, "GET", "/api/payment-methods", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	methods := resp["paymentMethods"].([]interface{})
	s.Equal(0, len(methods))
}

func (s *PaymentTestSuite) TestAdminCreatePaymentMethod() {
	body := map[string]interface{}{
		"name":         "微信支付",
		"qrCodeUrl":    "/test/wechat-qr.png",
		"instructions": "请使用微信扫码支付",
		"isActive":     true,
		"sortOrder":    1,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["success"])
	pm := resp["paymentMethod"].(map[string]interface{})
	s.Equal("微信支付", pm["name"])
	s.Equal("/test/wechat-qr.png", pm["qrCodeUrl"])
	s.Equal(float64(1), pm["sortOrder"])
}

func (s *PaymentTestSuite) TestAdminCreatePaymentMethodDuplicate() {
	body := map[string]interface{}{
		"name":      "微信支付",
		"isActive":  true,
		"sortOrder": 2,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods", body, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentTestSuite) TestAdminGetPaymentMethods() {
	w := MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	methods := resp["paymentMethods"].([]interface{})
	s.Equal(1, len(methods))
}

func (s *PaymentTestSuite) TestAdminCreateSecondPaymentMethod() {
	body := map[string]interface{}{
		"name":         "支付宝",
		"qrCodeUrl":    "/test/alipay-qr.png",
		"instructions": "请使用支付宝扫码支付",
		"isActive":     true,
		"sortOrder":    2,
	}
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods", body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *PaymentTestSuite) TestGetPublicPaymentMethods() {
	w := MakeRequest(s.router, "GET", "/api/payment-methods", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	methods := resp["paymentMethods"].([]interface{})
	s.Equal(2, len(methods))
}

func (s *PaymentTestSuite) TestAdminUpdatePaymentMethod() {
	w := MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	methods := listResp["paymentMethods"].([]interface{})
	first := methods[0].(map[string]interface{})
	pmID := first["id"].(float64)

	body := map[string]interface{}{
		"name":         "微信支付更新",
		"instructions": "请使用微信扫码，备注订单号",
	}
	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-methods/%d", int(pmID)), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	pm := resp["paymentMethod"].(map[string]interface{})
	s.Equal("微信支付更新", pm["name"])
}

func (s *PaymentTestSuite) TestAdminDeactivatePaymentMethod() {
	w := MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	methods := listResp["paymentMethods"].([]interface{})
	s.Equal(2, len(methods))

	second := methods[1].(map[string]interface{})
	pmID := second["id"].(float64)

	body := map[string]interface{}{
		"isActive": false,
	}
	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-methods/%d", int(pmID)), body, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/payment-methods", nil, nil)
	s.Equal(http.StatusOK, w.Code)
	var pubResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &pubResp)
	pubMethods := pubResp["paymentMethods"].([]interface{})
	s.Equal(1, len(pubMethods), "inactive method should not appear in public list")
}

func (s *PaymentTestSuite) TestAdminDeletePaymentMethod() {
	w := MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var listResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResp)
	methods := listResp["paymentMethods"].([]interface{})
	s.Equal(2, len(methods))

	pmID := methods[1].(map[string]interface{})["id"].(float64)
	w = MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/admin/payment-methods/%d", int(pmID)), nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &listResp)
	methods = listResp["paymentMethods"].([]interface{})
	s.Equal(1, len(methods))
}

func (s *PaymentTestSuite) TestPaymentMethodAdminOnly() {
	body := map[string]interface{}{
		"name": "非法创建",
	}
	w := MakeRequest(s.router, "POST", "/api/admin/payment-methods", body, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)

	w = MakeRequest(s.router, "GET", "/api/admin/payment-methods", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)

	w = MakeRequest(s.router, "DELETE", "/api/admin/payment-methods/1", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}

func (s *PaymentTestSuite) TestUploadPaymentProof() {
	pmID := s.createPaymentMethod("测试支付", "/test/qr.png")
	orderID := s.createPendingOrder()
	proofID := s.uploadProof(orderID, pmID)
	s.True(proofID > 0)
}

func (s *PaymentTestSuite) TestGetOrderPaymentProof() {
	pmID := s.createPaymentMethod("查询测试支付", "/test/qr2.png")
	orderID := s.createPendingOrder()
	s.uploadProof(orderID, pmID)

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.NotNil(resp["paymentProof"])
	proof := resp["paymentProof"].(map[string]interface{})
	s.Equal("pending", proof["status"])
}

func (s *PaymentTestSuite) TestGetPaymentProofWrongUser() {
	pmID := s.createPaymentMethod("越权测试", "/test/qr3.png")
	orderID := s.createPendingOrder()
	s.uploadProof(orderID, pmID)

	w := MakeRequest(s.router, "GET", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), nil, GetCustomer2AuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Nil(resp["paymentProof"])
}

func (s *PaymentTestSuite) TestUploadProofWithoutFile() {
	pmID := s.createPaymentMethod("无文件测试", "/test/qr4.png")
	orderID := s.createPendingOrder()

	formData := map[string]string{
		"paymentMethodId": fmt.Sprintf("%d", pmID),
	}
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), formData, nil, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentTestSuite) TestUploadProofNonPendingOrder() {
	pmID := s.createPaymentMethod("非pending测试", "/test/qr5.png")

	imagePath := CreateTestImageFile(".png")
	defer os.Remove(imagePath)

	formData := map[string]string{
		"paymentMethodId": fmt.Sprintf("%d", pmID),
	}
	fileFields := map[string]string{
		"file": imagePath,
	}
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", Data.Orders[0].ID), formData, fileFields, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentTestSuite) TestReviewPaymentProofApprove() {
	pmID := s.createPaymentMethod("审核通过测试", "/test/qr6.png")
	orderID := s.createPendingOrder()
	proofID := s.uploadProof(orderID, pmID)

	reviewBody := map[string]interface{}{
		"action": "approve",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID), reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var reviewResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &reviewResp)
	s.Equal(true, reviewResp["success"])
	reviewedProof := reviewResp["paymentProof"].(map[string]interface{})
	s.Equal("approved", reviewedProof["status"])

	var updatedOrder models.Order
	database.DB.First(&updatedOrder, orderID)
	s.Equal("processing", updatedOrder.Status)
}

func (s *PaymentTestSuite) TestReviewPaymentProofReject() {
	pmID := s.createPaymentMethod("审核驳回测试", "/test/qr7.png")
	orderID := s.createPendingOrder()
	proofID := s.uploadProof(orderID, pmID)

	reviewBody := map[string]interface{}{
		"action":       "reject",
		"rejectReason": "截图不清晰",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID), reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var reviewResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &reviewResp)
	reviewedProof := reviewResp["paymentProof"].(map[string]interface{})
	s.Equal("rejected", reviewedProof["status"])
	s.Equal("截图不清晰", reviewedProof["rejectReason"])

	var updatedOrder models.Order
	database.DB.First(&updatedOrder, orderID)
	s.Equal("pending", updatedOrder.Status)
}

func (s *PaymentTestSuite) TestReuploadReplacesPendingProof() {
	pmID := s.createPaymentMethod("重新上传测试", "/test/qr8.png")
	orderID := s.createPendingOrder()
	proofID1 := s.uploadProof(orderID, pmID)

	imagePath := CreateTestImageFile(".png")
	defer os.Remove(imagePath)

	formData := map[string]string{
		"paymentMethodId": fmt.Sprintf("%d", pmID),
	}
	fileFields := map[string]string{
		"file": imagePath,
	}
	w := MakeFormRequest(s.router, "POST", fmt.Sprintf("/api/orders/%d/payment-proof", orderID), formData, fileFields, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	proofID2 := uint(resp["paymentProof"].(map[string]interface{})["id"].(float64))
	s.Equal(proofID1, proofID2, "re-uploading should update existing proof, not create new")
}

func (s *PaymentTestSuite) TestReviewAlreadyReviewed() {
	pmID := s.createPaymentMethod("重复审核测试", "/test/qr9.png")
	orderID := s.createPendingOrder()
	proofID := s.uploadProof(orderID, pmID)

	reviewBody := map[string]interface{}{
		"action": "approve",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID), reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID), reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentTestSuite) TestReviewProofNotFound() {
	reviewBody := map[string]interface{}{
		"action": "approve",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/99999/review", reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *PaymentTestSuite) TestReviewInvalidAction() {
	pmID := s.createPaymentMethod("无效操作测试", "/test/qr10.png")
	orderID := s.createPendingOrder()
	proofID := s.uploadProof(orderID, pmID)

	reviewBody := map[string]interface{}{
		"action": "invalid",
	}
	w := MakeRequest(s.router, "PUT", fmt.Sprintf("/api/admin/payment-proofs/%d/review", proofID), reviewBody, GetAdminAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *PaymentTestSuite) TestReviewProofAdminOnly() {
	reviewBody := map[string]interface{}{
		"action": "approve",
	}
	w := MakeRequest(s.router, "PUT", "/api/admin/payment-proofs/1/review", reviewBody, GetCustomerAuthHeader())
	s.Equal(http.StatusForbidden, w.Code)
}
