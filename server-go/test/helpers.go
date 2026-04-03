package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"

	"bycigar-server/internal/config"
	"bycigar-server/internal/database"
	"bycigar-server/internal/handlers"
	"bycigar-server/internal/middleware"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/storage"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var AdminToken string
var CustomerToken string
var Customer2Token string

var AdminUser models.User
var CustomerUser models.User
var Customer2User models.User
var ServiceUser models.User

type SeedData struct {
	Categories []models.Category
	Products   []models.Product
	Banners    []models.Banner
	Pages      []models.Page
	Addresses  []models.Address
	CartItems  []models.CartItem
	Favorites  []models.Favorite
	Orders     []models.Order
}

var Data SeedData

func SetupTestConfig() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_NAME", "bycigar_test")
	os.Setenv("JWT_SECRET", "test-jwt-secret-key")
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("MINIO_ENDPOINT", "localhost:9000")
	os.Setenv("MINIO_ACCESS_KEY", "minioadmin")
	os.Setenv("MINIO_SECRET_KEY", "minioadmin123")
	os.Setenv("MINIO_BUCKET", "bycigar")
	os.Setenv("MINIO_USE_SSL", "false")
	config.LoadConfig()
}

func SetupTestDB() {
	database.Connect()
	database.Migrate()
	utils.InitSnowflake(1)
	storage.InitStorage("")
	CleanDB()
	SeedTestData()
}

func CleanDB() {
	db := database.DB
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE order_items")
	db.Exec("TRUNCATE TABLE orders")
	db.Exec("TRUNCATE TABLE payment_proofs")
	db.Exec("TRUNCATE TABLE payment_methods")
	db.Exec("TRUNCATE TABLE favorites")
	db.Exec("TRUNCATE TABLE cart_items")
	db.Exec("TRUNCATE TABLE addresses")
	db.Exec("TRUNCATE TABLE products")
	db.Exec("TRUNCATE TABLE categories")
	db.Exec("TRUNCATE TABLE banners")
	db.Exec("TRUNCATE TABLE pages")
	db.Exec("TRUNCATE TABLE settings")
	db.Exec("TRUNCATE TABLE site_configs")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("TRUNCATE TABLE notifications")
	db.Exec("TRUNCATE TABLE messages")
	db.Exec("TRUNCATE TABLE conversations")
	db.Exec("TRUNCATE TABLE quick_replies")
	db.Exec("TRUNCATE TABLE ratings")
	db.Exec("TRUNCATE TABLE contact_methods")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func findCategoryIDBySlug(db *gorm.DB, slug string) uint {
	var cat models.Category
	db.Where("slug = ?", slug).First(&cat)
	return cat.ID
}

func uintPtr(v uint) *uint {
	return &v
}

func SeedTestData() {
	db := database.DB

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	AdminUser = models.User{
		Email: "admin@test.com", Password: string(hashedPassword),
		Name: "TestAdmin", Role: "admin",
	}
	db.Create(&AdminUser)

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("user1234"), bcrypt.DefaultCost)
	CustomerUser = models.User{
		Email: "user1@test.com", Password: string(hashedPassword),
		Name: "TestUser1", Role: "customer",
	}
	db.Create(&CustomerUser)

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("user1234"), bcrypt.DefaultCost)
	Customer2User = models.User{
		Email: "user2@test.com", Password: string(hashedPassword),
		Name: "TestUser2", Role: "customer",
	}
	db.Create(&Customer2User)

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("service123"), bcrypt.DefaultCost)
	ServiceUser = models.User{
		Email: "service@test.com", Password: string(hashedPassword),
		Name: "TestService", Role: "service",
	}
	db.Create(&ServiceUser)

	topCategories := []models.Category{
		{Name: "古巴雪茄", Slug: "cohiba"},
		{Name: "多米尼加", Slug: "dominican"},
		{Name: "配件", Slug: "accessories"},
	}
	db.Create(&topCategories)

	cohibaID := findCategoryIDBySlug(db, "cohiba")
	dominicanID := findCategoryIDBySlug(db, "dominican")

	subCategories := []models.Category{
		{Name: "高希霸", Slug: "cohiba-classic", ParentID: uintPtr(cohibaID)},
		{Name: "蒙特", Slug: "montecristo", ParentID: uintPtr(cohibaID)},
		{Name: "大卫杜夫", Slug: "davidoff", ParentID: uintPtr(dominicanID)},
	}
	db.Create(&subCategories)

	Data.Categories = append(topCategories, subCategories...)

	classicID := findCategoryIDBySlug(db, "cohiba-classic")
	monteID := findCategoryIDBySlug(db, "montecristo")
	davidoffID := findCategoryIDBySlug(db, "davidoff")

	Data.Products = []models.Product{
		{Name: "高希霸世纪一号", Slug: "cohiba-siglo-i", Price: 120.00, CategoryID: classicID, Stock: 50, IsActive: true, IsFeatured: true, Image: "/test/cohiba-1.jpg", Description: "经典古巴雪茄"},
		{Name: "高希霸世纪二号", Slug: "cohiba-siglo-ii", Price: 150.00, CategoryID: classicID, Stock: 30, IsActive: true, IsFeatured: false, Image: "/test/cohiba-2.jpg", Description: "中等浓郁"},
		{Name: "高希霸世纪三号", Slug: "cohiba-siglo-iii", Price: 180.00, CategoryID: classicID, Stock: 0, IsActive: true, IsFeatured: false, Image: "/test/cohiba-3.jpg", Description: "缺货商品"},
		{Name: "高希霸短号", Slug: "cohiba-short", Price: 60.00, CategoryID: classicID, Stock: 100, IsActive: false, IsFeatured: false, Image: "/test/cohiba-short.jpg", Description: "已下架"},
		{Name: "高希霸马杜罗", Slug: "cohiba-maduro", Price: 250.00, CategoryID: classicID, Stock: 20, IsActive: true, IsFeatured: true, Image: "/test/cohiba-maduro.jpg", Description: "浓郁口感"},
		{Name: "蒙特2号", Slug: "montecristo-no2", Price: 100.00, CategoryID: monteID, Stock: 40, IsActive: true, IsFeatured: false, Image: "/test/monte-2.jpg", Description: "经典鱼雷"},
		{Name: "蒙特4号", Slug: "montecristo-no4", Price: 80.00, CategoryID: monteID, Stock: 60, IsActive: true, IsFeatured: false, Image: "/test/monte-4.jpg", Description: "最畅销"},
		{Name: "蒙特埃德蒙多", Slug: "montecristo-edmundo", Price: 130.00, CategoryID: monteID, Stock: 25, IsActive: true, IsFeatured: false, Image: "/test/monte-edmundo.jpg", Description: "丰富层次"},
		{Name: "大卫杜夫2000", Slug: "davidoff-2000", Price: 200.00, CategoryID: davidoffID, Stock: 35, IsActive: true, IsFeatured: true, Image: "/test/davidoff-2000.jpg", Description: "瑞士品质"},
		{Name: "大卫杜夫千年", Slug: "davidoff-millennium", Price: 300.00, CategoryID: davidoffID, Stock: 15, IsActive: true, IsFeatured: false, Image: "/test/davidoff-millennium.jpg", Description: "限量版"},
		{Name: "大卫杜夫温斯顿", Slug: "davidoff-winston", Price: 450.00, CategoryID: davidoffID, Stock: 10, IsActive: true, IsFeatured: true, Image: "/test/davidoff-winston.jpg", Description: "顶级"},
	}
	db.Create(&Data.Products)
	db.Model(&Data.Products[3]).Update("is_active", false)

	db.Where("1=1").Find(&Data.Categories)

	Data.Banners = []models.Banner{
		{Title: "主推古巴精选", Image: "/test/banner1.jpg", SortOrder: 1, IsActive: true},
		{Title: "新品上市", Image: "/test/banner2.jpg", SortOrder: 2, IsActive: true},
		{Title: "已过期活动", Image: "/test/banner3.jpg", SortOrder: 3, IsActive: false},
	}
	db.Create(&Data.Banners)
	db.Model(&Data.Banners[2]).Update("is_active", false)

	Data.Pages = []models.Page{
		{Slug: "about", Title: "关于我们", Content: "# 关于我们\n\n测试内容"},
		{Slug: "services", Title: "服务条款", Content: "# 服务条款\n\n测试内容"},
		{Slug: "privacy-policy", Title: "隐私政策", Content: "# 隐私政策\n\n测试内容"},
		{Slug: "statement", Title: "严正声明", Content: "# 严正声明\n\n测试内容"},
	}
	db.Create(&Data.Pages)

	Data.Addresses = []models.Address{
		{UserID: CustomerUser.ID, FullName: "张三", AddressLine1: "123 Main St", City: "Beijing", State: "Beijing", ZipCode: "100000", Phone: "13800138001", IsDefault: true},
		{UserID: CustomerUser.ID, FullName: "张三", AddressLine1: "456 Side St", City: "Shanghai", State: "Shanghai", ZipCode: "200000", Phone: "13800138002", IsDefault: false},
	}
	db.Create(&Data.Addresses)

	Data.CartItems = []models.CartItem{
		{UserID: CustomerUser.ID, ProductID: Data.Products[0].ID, Quantity: 2},
		{UserID: CustomerUser.ID, ProductID: Data.Products[5].ID, Quantity: 1},
	}
	db.Create(&Data.CartItems)

	Data.Favorites = []models.Favorite{
		{UserID: CustomerUser.ID, ProductID: Data.Products[0].ID},
		{UserID: CustomerUser.ID, ProductID: Data.Products[5].ID},
		{UserID: CustomerUser.ID, ProductID: Data.Products[8].ID},
	}
	db.Create(&Data.Favorites)

	orderNo := utils.GenerateOrderNo()
	Data.Orders = []models.Order{
		{
			OrderNo: orderNo, UserID: CustomerUser.ID, AddressID: Data.Addresses[0].ID,
			Total: 340.00, Status: "completed",
			Items: []models.OrderItem{
				{ProductID: Data.Products[0].ID, Quantity: 2, Price: 120.00},
				{ProductID: Data.Products[5].ID, Quantity: 1, Price: 100.00},
			},
		},
	}
	db.Create(&Data.Orders)

	SeedBulkOrders(3000)
}

func SeedBulkOrders(count int) {
	db := database.DB

	statuses := []string{"pending", "processing", "shipped", "completed", "cancelled"}
	userIDs := []uint{CustomerUser.ID, Customer2User.ID}
	addressIDs := []uint{Data.Addresses[0].ID, Data.Addresses[1].ID}

	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))

	for i := 0; i < count; i++ {
		status := statuses[r.Intn(len(statuses))]
		userID := userIDs[r.Intn(len(userIDs))]
		addressID := addressIDs[r.Intn(len(addressIDs))]

		itemCount := r.Intn(5) + 1
		var total float64
		items := make([]models.OrderItem, itemCount)

		for j := 0; j < itemCount; j++ {
			product := Data.Products[r.Intn(len(Data.Products))]
			quantity := r.Intn(10) + 1
			price := product.Price
			total += price * float64(quantity)
			items[j] = models.OrderItem{
				ProductID: product.ID,
				Quantity:  quantity,
				Price:     price,
			}
		}

		daysAgo := r.Intn(90)
		createdAt := now.AddDate(0, 0, -daysAgo).Add(time.Duration(r.Intn(86400)) * time.Second)

		order := models.Order{
			OrderNo:   utils.GenerateOrderNo(),
			UserID:    userID,
			AddressID: addressID,
			Total:     total,
			Status:    status,
			Remark:    fmt.Sprintf("测试订单 %d", i+1),
			Items:     items,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}
		db.Create(&order)
	}

	var totalOrders int64
	db.Model(&models.Order{}).Count(&totalOrders)
	fmt.Printf("已生成 %d 笔测试订单（当前总计: %d）\n", count, totalOrders)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)
	r.GET("/api/auth/me", handlers.GetProfile)
	r.PUT("/api/auth/profile", handlers.UpdateProfile)
	r.GET("/api/auth/captcha", handlers.GetCaptcha)
	r.PUT("/api/auth/change-password", middleware.RequireAuth(), handlers.ChangePassword)

	r.GET("/api/products", handlers.GetProducts)
	r.GET("/api/products/suggest", handlers.GetProductSuggestions)
	r.GET("/api/products/:id", handlers.GetProduct)
	r.GET("/api/categories", handlers.GetCategories)
	r.GET("/api/config", handlers.GetConfig)
	r.GET("/api/site-identity", handlers.GetSiteIdentity)
	r.GET("/api/banners", handlers.GetBanners)
	r.GET("/api/pages/:slug", handlers.GetPage)
	r.GET("/api/settings", handlers.GetSettings)

	r.GET("/api/cart", handlers.GetCart)
	r.POST("/api/cart", handlers.AddToCart)
	r.PUT("/api/cart/:id", handlers.UpdateCartItem)
	r.DELETE("/api/cart/:id", handlers.DeleteCartItem)

	r.GET("/api/favorites", handlers.GetFavorites)
	r.POST("/api/favorites", handlers.AddFavorite)
	r.DELETE("/api/favorites/:productId", handlers.DeleteFavorite)

	r.GET("/api/addresses", handlers.GetAddresses)
	r.POST("/api/addresses", handlers.CreateAddress)
	r.PUT("/api/addresses/:id", handlers.UpdateAddress)
	r.PUT("/api/addresses/:id/default", handlers.SetDefaultAddress)
	r.DELETE("/api/addresses/:id", handlers.DeleteAddress)

	r.GET("/api/notifications", middleware.RequireAuth(), handlers.GetNotifications)
	r.GET("/api/notifications/unread-count", middleware.RequireAuth(), handlers.GetUnreadCount)
	r.GET("/api/notifications/:id", middleware.RequireAuth(), handlers.GetNotification)
	r.PUT("/api/notifications/:id/read", middleware.RequireAuth(), handlers.MarkAsRead)
	r.PUT("/api/notifications/read-all", middleware.RequireAuth(), handlers.MarkAllRead)
	r.DELETE("/api/notifications/:id", middleware.RequireAuth(), handlers.DeleteNotification)
	r.DELETE("/api/notifications/read", middleware.RequireAuth(), handlers.DeleteReadNotifications)

	r.GET("/api/orders", handlers.GetOrders)
	r.POST("/api/orders", handlers.CreateOrder)
	r.GET("/api/orders/:id", handlers.GetOrder)

	r.GET("/api/payment-methods", handlers.GetPaymentMethods)
	r.POST("/api/orders/:id/payment-proof", middleware.RequireAuth(), handlers.UploadPaymentProof)
	r.GET("/api/orders/:id/payment-proof", middleware.RequireAuth(), handlers.GetOrderPaymentProof)

	r.POST("/api/chat/conversations", middleware.RequireAuth(), handlers.CreateConversation)
	r.GET("/api/chat/conversations", middleware.RequireAuth(), handlers.GetConversations)
	r.GET("/api/chat/conversations/:id/messages", middleware.RequireAuth(), handlers.GetMessages)
	r.POST("/api/chat/conversations/:id/messages", middleware.RequireAuth(), handlers.SendMessage)
	r.PUT("/api/chat/conversations/:id/close", middleware.RequireAuth(), handlers.CustomerCloseConversation)
	r.POST("/api/chat/conversations/:id/rate", middleware.RequireAuth(), handlers.RateConversation)
	r.GET("/api/chat/unread-count", middleware.RequireAuth(), handlers.GetChatUnreadCount)
	r.POST("/api/chat/upload-image", middleware.RequireAuth(), handlers.UploadChatImage)
	r.GET("/api/chat/service-status", handlers.GetServiceStatus)

	admin := r.Group("/api/admin")
	admin.Use(middleware.AdminOnly)
	{
		admin.POST("/upload", handlers.UploadImage)

		admin.GET("/products", handlers.GetAdminProducts)
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
		admin.PUT("/products/batch/status", handlers.BatchUpdateProductStatus)
		admin.DELETE("/products/batch", handlers.BatchDeleteProducts)

		admin.GET("/categories", handlers.GetAdminCategories)
		admin.POST("/categories", handlers.CreateCategory)
		admin.PUT("/categories/:id", handlers.UpdateCategory)
		admin.DELETE("/categories/:id", handlers.DeleteCategory)

		admin.GET("/orders", handlers.GetAdminOrders)
		admin.GET("/orders/export", handlers.ExportAdminOrders)
		admin.GET("/orders/:id", handlers.GetAdminOrder)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)

		admin.GET("/dashboard/stats", handlers.GetDashboardStats)
		admin.GET("/dashboard/recent-orders", handlers.GetDashboardRecentOrders)
		admin.GET("/dashboard/low-stock", handlers.GetDashboardLowStock)
		admin.GET("/dashboard/top-products", handlers.GetDashboardTopProducts)
		admin.GET("/stats/revenue", handlers.GetRevenueByDate)

		admin.GET("/users", handlers.GetAdminUsers)
		admin.GET("/users/:id", handlers.GetAdminUser)
		admin.POST("/users/:id/reset-password", handlers.ResetUserPassword)

		admin.GET("/chat/conversations", handlers.GetAdminConversations)
		admin.GET("/chat/conversations/:id/messages", handlers.GetAdminMessages)
		admin.POST("/chat/conversations/:id/messages", handlers.AdminSendMessage)
		admin.PUT("/chat/conversations/:id/close", handlers.CloseConversation)
		admin.PUT("/chat/conversations/:id/assign", handlers.AssignConversation)
		admin.POST("/chat/conversations/:id/messages/:msgId/recall", handlers.RecallMessage)
		admin.GET("/chat/unread-stats", handlers.GetAdminUnreadStats)
		admin.POST("/chat/service-status", handlers.SetServiceStatus)
		admin.GET("/quick-replies", handlers.GetQuickReplies)
		admin.POST("/quick-replies", handlers.CreateQuickReply)
		admin.PUT("/quick-replies/:id", handlers.UpdateQuickReply)
		admin.DELETE("/quick-replies/:id", handlers.DeleteQuickReply)
		admin.GET("/stats/satisfaction", handlers.GetSatisfactionStats)
		admin.GET("/chat/agent-stats", handlers.GetAgentStats)

		admin.PUT("/payment-proofs/batch-review", handlers.BatchReviewPaymentProofs)
		admin.PUT("/payment-proofs/:id/review", handlers.ReviewPaymentProof)
	}

	superAdmin := r.Group("/api/admin")
	superAdmin.Use(middleware.SuperAdminOnly)
	{
		superAdmin.GET("/banners", handlers.GetAdminBanners)
		superAdmin.POST("/banners", handlers.CreateBanner)
		superAdmin.PUT("/banners/:id", handlers.UpdateBanner)
		superAdmin.DELETE("/banners/:id", handlers.DeleteBanner)

		superAdmin.GET("/pages", handlers.GetAdminPages)
		superAdmin.PUT("/pages/:slug", handlers.UpdatePage)

		superAdmin.PUT("/config/:key", handlers.UpdateConfig)
		superAdmin.PUT("/settings/:key", handlers.UpdateSetting)

		superAdmin.PUT("/users/:id/role", handlers.UpdateUserRole)
		superAdmin.PUT("/users/:id/ban", handlers.BanUser)
		superAdmin.PUT("/users/:id/unban", handlers.UnbanUser)
		superAdmin.DELETE("/users/:id", handlers.DeleteUser)

		superAdmin.GET("/payment-methods", handlers.GetAdminPaymentMethods)
		superAdmin.POST("/payment-methods", handlers.CreatePaymentMethod)
		superAdmin.PUT("/payment-methods/:id", handlers.UpdatePaymentMethod)
		superAdmin.DELETE("/payment-methods/:id", handlers.DeletePaymentMethod)

		superAdmin.GET("/contact-methods", handlers.GetAdminContactMethods)
		superAdmin.POST("/contact-methods", handlers.CreateContactMethod)
		superAdmin.PUT("/contact-methods/:id", handlers.UpdateContactMethod)
		superAdmin.DELETE("/contact-methods/:id", handlers.DeleteContactMethod)

		superAdmin.POST("/email/test", handlers.TestEmail)
		superAdmin.POST("/cleanup", handlers.BatchCleanup)
	}

	return r
}

func MakeRequest(r *gin.Engine, method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func MakeFormRequest(r *gin.Engine, method, path string, formData map[string]string, fileFields map[string]string, headers map[string]string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for key, value := range formData {
		writer.WriteField(key, value)
	}

	for field, filePath := range fileFields {
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		defer file.Close()
		part, _ := writer.CreateFormFile(field, filepath.Base(filePath))
		io.Copy(part, file)
	}

	writer.Close()

	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetAdminAuthHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("user-%d", AdminUser.ID),
	}
}

func GetCustomerAuthHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("user-%d", CustomerUser.ID),
	}
}

func GetCustomer2AuthHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("user-%d", Customer2User.ID),
	}
}

func GetServiceAuthHeader() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("user-%d", ServiceUser.ID),
	}
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func CreateTestImageFile(ext string) string {
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("test_upload_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(tmpDir, fileName)

	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(0xAB)
	}
	os.WriteFile(filePath, data, 0644)
	return filePath
}

func CreateLargeTestImageFile(ext string) string {
	tmpDir := os.TempDir()
	fileName := fmt.Sprintf("test_large_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(tmpDir, fileName)

	data := make([]byte, 11*1024*1024)
	os.WriteFile(filePath, data, 0644)
	return filePath
}
