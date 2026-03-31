package handlers

import (
	"net/http"

	"bycigar-server/pkg/email"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

type TestEmailRequest struct {
	To string `json:"to" binding:"required,email"`
}

// TestEmail godoc
// @Summary 发送测试邮件
// @Description 使用当前SMTP配置发送一封测试邮件
// @Tags admin
// @Accept json
// @Produce json
// @Param body body TestEmailRequest true "收件邮箱"
// @Success 200 {object} map[string]interface{}
// @Router /admin/email/test [post]
func TestEmail(c *gin.Context) {
	var req TestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请输入有效的邮箱地址")
		return
	}

	config := email.GetEmailConfig()
	if config == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "邮件功能未启用或配置不完整，请先完成SMTP配置")
		return
	}

	subject := "邮件配置测试 - " + config.FromName
	body := `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="margin:0;padding:0;background:#f5f5f5;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<div style="max-width:600px;margin:40px auto;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
	<div style="background:linear-gradient(135deg,#d4a574,#b08968);padding:32px;text-align:center;">
		<h1 style="margin:0;color:#fff;font-size:22px;">` + config.FromName + `</h1>
	</div>
	<div style="padding:32px;">
		<h2 style="margin:0 0 16px;font-size:18px;color:#333;">邮件配置测试成功</h2>
		<p style="margin:0 0 16px;font-size:14px;color:#666;line-height:1.6;">如果您收到此邮件，说明您的SMTP邮件配置已正确设置。</p>
		<div style="background:#f0faf0;border:1px solid #c8e6c8;border-radius:6px;padding:16px 20px;">
			<p style="margin:0;font-size:14px;color:#2e7d32;">SMTP服务器：` + config.Host + `</p>
			<p style="margin:4px 0 0;font-size:14px;color:#2e7d32;">发件邮箱：` + config.Username + `</p>
		</div>
	</div>
	<div style="background:#f7f7f7;padding:20px 32px;text-align:center;">
		<p style="margin:0;font-size:12px;color:#999;">此邮件由系统自动发送，请勿直接回复。</p>
	</div>
	</div>
</body>
</html>`

	if err := email.SendMail(config, req.To, subject, body); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "测试邮件发送失败："+err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"success": true, "message": "测试邮件发送成功"})
}
