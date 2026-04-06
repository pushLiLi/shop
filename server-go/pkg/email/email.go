package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"bycigar-server/internal/database"
	"bycigar-server/internal/models"
)

var (
	emailCache     map[string]string
	emailCacheTime time.Time
	emailCacheMu   sync.RWMutex
	emailCacheTTL  = 5 * time.Minute
)

func getEmailSettingsCached() (map[string]string, error) {
	emailCacheMu.RLock()
	if emailCache != nil && time.Since(emailCacheTime) < emailCacheTTL {
		result := emailCache
		emailCacheMu.RUnlock()
		return result, nil
	}
	emailCacheMu.RUnlock()

	emailCacheMu.Lock()
	defer emailCacheMu.Unlock()

	if emailCache != nil && time.Since(emailCacheTime) < emailCacheTTL {
		return emailCache, nil
	}

	var settings []models.Setting
	if err := database.DB.Where("`key` LIKE ?", "email_%").Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	emailCache = result
	emailCacheTime = time.Now()
	return result, nil
}

func InvalidateEmailCache() {
	emailCacheMu.Lock()
	emailCache = nil
	emailCacheMu.Unlock()
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromName string
	Enabled  bool
}

func GetEmailConfig() *EmailConfig {
	settings, err := getEmailSettingsCached()
	if err != nil {
		log.Printf("获取邮件配置失败: %v", err)
		return nil
	}

	if settings["email_enabled"] != "true" {
		return nil
	}

	host := settings["email_smtp_host"]
	port := settings["email_smtp_port"]
	username := settings["email_smtp_username"]
	password := settings["email_smtp_password"]
	fromName := settings["email_from_name"]

	if host == "" || port == "" || username == "" || password == "" {
		return nil
	}

	return &EmailConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		FromName: fromName,
		Enabled:  true,
	}
}

func RenderTemplate(templateStr string, data map[string]string) string {
	result := templateStr
	for key, value := range data {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

func getCurrencySymbol(currency string) string {
	if currency == "USD" {
		return "$"
	}
	return "¥"
}

func BuildOrderItemsHTML(items []models.OrderItem) string {
	var rows strings.Builder
	for _, item := range items {
		productName := item.Product.Name
		if productName == "" {
			productName = fmt.Sprintf("商品 #%d", item.ProductID)
		}
		symbol := getCurrencySymbol(item.Currency)
		rows.WriteString(fmt.Sprintf(`<tr>
			<td style="padding:10px 16px;border-bottom:1px solid #eee;font-size:14px;">%s</td>
			<td style="padding:10px 16px;border-bottom:1px solid #eee;font-size:14px;text-align:center;">%d</td>
			<td style="padding:10px 16px;border-bottom:1px solid #eee;font-size:14px;text-align:right;">%s%.2f</td>
		</tr>`, productName, item.Quantity, symbol, item.Price))
	}

	return fmt.Sprintf(`<table style="width:100%%;border-collapse:collapse;margin:16px 0;">
		<thead>
			<tr style="background:#f7f7f7;">
				<th style="padding:10px 16px;text-align:left;font-size:13px;color:#666;border-bottom:2px solid #ddd;">商品名称</th>
				<th style="padding:10px 16px;text-align:center;font-size:13px;color:#666;border-bottom:2px solid #ddd;">数量</th>
				<th style="padding:10px 16px;text-align:right;font-size:13px;color:#666;border-bottom:2px solid #ddd;">单价</th>
			</tr>
		</thead>
		<tbody>%s</tbody>
	</table>`, rows.String())
}

func buildShippingEmailHTML(order models.Order, user models.User, siteTitle string) (string, string) {
	subjectTemplate := `您的订单 {{orderNo}} 已发货`
	bodyTemplate := `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="margin:0;padding:0;background:#f5f5f5;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;">
<div style="max-width:600px;margin:40px auto;background:#fff;border-radius:8px;overflow:hidden;box-shadow:0 2px 8px rgba(0,0,0,0.08);">
	<div style="background:linear-gradient(135deg,#d4a574,#b08968);padding:32px;text-align:center;">
		<h1 style="margin:0;color:#fff;font-size:22px;">{{siteTitle}}</h1>
	</div>
	<div style="padding:32px;">
		<h2 style="margin:0 0 8px;font-size:18px;color:#333;">尊敬的 {{customerName}}，您好！</h2>
		<p style="margin:0 0 24px;font-size:14px;color:#666;line-height:1.6;">您的订单已发货，以下是物流信息：</p>
		<div style="background:#fdf8f3;border:1px solid #f0e0cc;border-radius:6px;padding:16px 20px;margin-bottom:24px;">
			<p style="margin:0 0 8px;font-size:14px;color:#333;"><strong>订单号：</strong>{{orderNo}}</p>
			<p style="margin:0 0 8px;font-size:14px;color:#333;"><strong>物流公司：</strong>{{trackingCompany}}</p>
			<p style="margin:0;font-size:14px;color:#333;"><strong>快递单号：</strong>{{trackingNumber}}</p>
		</div>
		<h3 style="margin:0 0 12px;font-size:15px;color:#333;">订单商品</h3>
		{{orderItems}}
		<div style="text-align:right;margin-top:16px;padding-top:16px;border-top:1px solid #eee;">
			<p style="margin:0;font-size:16px;color:#d4a574;"><strong>订单总额：¥{{orderTotal}}</strong></p>
		</div>
	</div>
	<div style="background:#f7f7f7;padding:20px 32px;text-align:center;">
		<p style="margin:0;font-size:12px;color:#999;">此邮件由系统自动发送，请勿直接回复。</p>
	</div>
	</div>
</body>
</html>`

	data := map[string]string{
		"orderNo":         order.OrderNo,
		"trackingCompany": order.TrackingCompany,
		"trackingNumber":  order.TrackingNumber,
		"customerName":    user.Name,
		"siteTitle":       siteTitle,
		"orderTotal":      fmt.Sprintf("%.2f", order.Total),
		"orderItems":      BuildOrderItemsHTML(order.Items),
	}

	subject := RenderTemplate(subjectTemplate, data)
	body := RenderTemplate(bodyTemplate, data)
	return subject, body
}

func SendMail(config *EmailConfig, to, subject, htmlBody string) error {
	addr := net.JoinHostPort(config.Host, config.Port)
	from := fmt.Sprintf("%s <%s>", config.FromName, config.Username)
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		from, to, subject,
	)
	msg := headers + htmlBody

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS连接失败: %w", err)
	}

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("创建SMTP客户端失败: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %w", err)
	}

	if err = client.Mail(config.Username); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取写入器失败: %w", err)
	}

	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("关闭写入器失败: %w", err)
	}

	return client.Quit()
}

func SendShippingNotification(order models.Order) {
	config := GetEmailConfig()
	if config == nil {
		return
	}

	var user models.User
	if err := database.DB.First(&user, order.UserID).Error; err != nil {
		log.Printf("邮件通知: 查询用户失败 userId=%d err=%v", order.UserID, err)
		return
	}

	if user.Email == "" {
		log.Printf("邮件通知: 用户无邮箱 userId=%d", order.UserID)
		return
	}

	siteTitle := "商城"
	if s, err := getEmailSettingsCached(); err == nil {
		if t, ok := s["site_title"]; ok && t != "" {
			siteTitle = t
		}
	}

	subject, body := buildShippingEmailHTML(order, user, siteTitle)

	if err := SendMail(config, user.Email, subject, body); err != nil {
		log.Printf("邮件通知发送失败: userId=%d email=%s err=%v", user.ID, user.Email, err)
	} else {
		log.Printf("邮件通知发送成功: userId=%d email=%s orderNo=%s", user.ID, user.Email, order.OrderNo)
	}
}
