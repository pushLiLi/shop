package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
	"bycigar-server/internal/ws"
	imagepkg "bycigar-server/pkg/image"
	miniopkg "bycigar-server/pkg/minio"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadPaymentProof(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	orderParam := c.Param("id")
	var order models.Order
	if id, err := strconv.Atoi(orderParam); err == nil {
		if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
	} else {
		if err := database.DB.Where("order_no = ? AND user_id = ?", orderParam, userID).First(&order).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
	}

	if order.Status != models.OrderStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前订单状态不允许上传付款凭证"})
		return
	}

	paymentMethodIDStr := c.PostForm("paymentMethodId")
	paymentMethodID, err := strconv.Atoi(paymentMethodIDStr)
	if err != nil || paymentMethodID == 0 {
		var existingProof models.PaymentProof
		if err := database.DB.Where("order_id = ?", order.ID).Order("created_at desc").First(&existingProof).Error; err == nil {
			paymentMethodID = int(existingProof.PaymentMethodID)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请选择付款方式"})
			return
		}
	}

	var paymentMethod models.PaymentMethod
	if err := database.DB.Where("id = ? AND is_active = ?", paymentMethodID, true).First(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "付款方式不存在或已停用"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的付款截图"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 jpg, jpeg, png, gif, webp 格式的图片"})
		return
	}
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 10MB"})
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	result, _ := imagepkg.Process(fileBytes, ext)
	baseName := fmt.Sprintf("%d_%s", time.Now().Unix(), uuid.New().String())
	origName := baseName + result.OrigExt

	_, err = miniopkg.Client.PutObject(
		context.Background(),
		miniopkg.Bucket,
		origName,
		bytes.NewReader(result.Original),
		int64(len(result.Original)),
		minio.PutObjectOptions{ContentType: result.ContentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败"})
		return
	}

	imageUrl := fmt.Sprintf("/media/%s/%s", miniopkg.Bucket, origName)

	var existingProof models.PaymentProof
	hasExisting := database.DB.Where("order_id = ? AND status = ?", order.ID, models.PaymentProofStatusPending).First(&existingProof).Error == nil

	if hasExisting {
		oldImageURL := existingProof.ImageUrl
		database.DB.Model(&existingProof).Updates(map[string]interface{}{
			"payment_method_id": paymentMethodID,
			"image_url":         imageUrl,
		})
		if oldImageURL != "" && oldImageURL != imageUrl {
			miniopkg.DeleteObjects([]string{oldImageURL})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "paymentProof": existingProof})
	} else {
		proof := models.PaymentProof{
			OrderID:         order.ID,
			UserID:          userID.(uint),
			PaymentMethodID: uint(paymentMethodID),
			ImageUrl:        imageUrl,
			Status:          models.PaymentProofStatusPending,
		}
		database.DB.Create(&proof)
		c.JSON(http.StatusOK, gin.H{"success": true, "paymentProof": proof})
	}
}

func ReviewPaymentProof(c *gin.Context) {
	id := c.Param("id")

	var proof models.PaymentProof
	if err := database.DB.First(&proof, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "付款凭证不存在"})
		return
	}

	if proof.Status != models.PaymentProofStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该凭证已被审核"})
		return
	}

	var input models.ReviewPaymentProofInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviewerID, _ := c.Get("userID")
	now := time.Now()

	if input.Action == "approve" {
		proof.Status = models.PaymentProofStatusApproved
		proof.ReviewerID = func() *uint { v := reviewerID.(uint); return &v }()
		proof.ReviewedAt = &now
		database.DB.Save(&proof)

		database.DB.Model(&models.Order{}).Where("id = ?", proof.OrderID).Update("status", models.OrderStatusProcessing)

		var order models.Order
		database.DB.First(&order, proof.OrderID)

		notification := models.Notification{
			UserID:  order.UserID,
			Type:    models.NotificationTypeOrderStatus,
			Title:   "订单状态更新",
			Content: "您的订单 " + order.OrderNo + " 付款已确认，订单开始处理",
			Link:    "/orders",
			OrderID: &order.ID,
		}
		database.DB.Create(&notification)
		ws.DefaultHub.SendToUser(notification.UserID, gin.H{
			"type":         "notification",
			"notification": notification,
		})

	} else if input.Action == "reject" {
		proof.Status = models.PaymentProofStatusRejected
		proof.ReviewerID = func() *uint { v := reviewerID.(uint); return &v }()
		proof.ReviewedAt = &now
		proof.RejectReason = input.RejectReason
		database.DB.Save(&proof)

		var order models.Order
		database.DB.First(&order, proof.OrderID)

		content := "您的订单 " + order.OrderNo + " 付款凭证未通过审核"
		if input.RejectReason != "" {
			content += "，原因：" + input.RejectReason
		}
		content += "。请重新上传付款凭证或联系客服处理退款。"

		notification := models.Notification{
			UserID:  order.UserID,
			Type:    models.NotificationTypeOrderStatus,
			Title:   "付款凭证审核未通过",
			Content: content,
			Link:    "/orders",
			OrderID: &order.ID,
		}
		database.DB.Create(&notification)
		ws.DefaultHub.SendToUser(notification.UserID, gin.H{
			"type":         "notification",
			"notification": notification,
		})

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审核操作"})
		return
	}

	database.DB.Preload("PaymentMethod").First(&proof, id)
	c.JSON(http.StatusOK, gin.H{"success": true, "paymentProof": proof})
}

func BatchReviewPaymentProofs(c *gin.Context) {
	var input struct {
		IDs          []uint `json:"ids" binding:"required"`
		Action       string `json:"action" binding:"required"`
		RejectReason string `json:"rejectReason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Action != "approve" && input.Action != "reject" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的审核操作"})
		return
	}
	if len(input.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要审核的凭证"})
		return
	}

	var proofs []models.PaymentProof
	database.DB.Where("id IN ? AND status = ?", input.IDs, models.PaymentProofStatusPending).Find(&proofs)
	if len(proofs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有可审核的凭证"})
		return
	}

	reviewerID, _ := c.Get("userID")
	now := time.Now()
	reviewed := 0

	for _, p := range proofs {
		p.Status = models.PaymentProofStatusApproved
		if input.Action == "reject" {
			p.Status = models.PaymentProofStatusRejected
			p.RejectReason = input.RejectReason
		}
		p.ReviewerID = func() *uint { v := reviewerID.(uint); return &v }()
		p.ReviewedAt = &now
		database.DB.Save(&p)

		var order models.Order
		database.DB.First(&order, p.OrderID)

		if input.Action == "approve" {
			database.DB.Model(&models.Order{}).Where("id = ?", p.OrderID).Update("status", models.OrderStatusProcessing)
			notification := models.Notification{
				UserID:  order.UserID,
				Type:    models.NotificationTypeOrderStatus,
				Title:   "订单状态更新",
				Content: "您的订单 " + order.OrderNo + " 付款已确认，订单开始处理",
				Link:    "/orders",
				OrderID: &order.ID,
			}
			database.DB.Create(&notification)
		} else {
			content := "您的订单 " + order.OrderNo + " 付款凭证未通过审核"
			if input.RejectReason != "" {
				content += "，原因：" + input.RejectReason
			}
			content += "。请重新上传付款凭证或联系客服处理退款。"
			notification := models.Notification{
				UserID:  order.UserID,
				Type:    models.NotificationTypeOrderStatus,
				Title:   "付款凭证审核未通过",
				Content: content,
				Link:    "/orders",
				OrderID: &order.ID,
			}
			database.DB.Create(&notification)
		}
		reviewed++
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "reviewed": reviewed})
}

func GetOrderPaymentProof(c *gin.Context) {
	orderID := c.Param("id")

	var proof models.PaymentProof
	if err := database.DB.Where("order_id = ?", orderID).Preload("PaymentMethod").Order("created_at desc").First(&proof).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"paymentProof": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"paymentProof": proof})
}
