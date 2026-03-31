package main

import (
	"log"

	_ "bycigar-server/docs"
	"bycigar-server/internal/config"
	"bycigar-server/internal/database"
	"bycigar-server/internal/handlers"
	"bycigar-server/internal/middleware"
	"bycigar-server/internal/ws"
	miniopkg "bycigar-server/pkg/minio"
	pkgutils "bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title BYCIGAR API
// @version 1.0
// @description BYCIGAR 电商后台 API 接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@bycigar.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @host localhost:3000
// @BasePath /api
// @schemes http
func main() {
	godotenv.Load()

	config.LoadConfig()
	database.Connect()
	database.Migrate()
	database.Seed()
	miniopkg.InitMinio()
	miniopkg.EnsureBucket(config.AppConfig.MinioBucket)
	pkgutils.InitSnowflake(1)
	database.BackfillOrderNo()
	pkgutils.StartChatCleanup(database.DB)
	pkgutils.StartNotificationCleanup(database.DB)

	ws.DefaultHub = ws.NewHub()
	go ws.DefaultHub.Run()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.Static("/static", "./static")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json"), ginSwagger.DefaultModelsExpandDepth(-1)))

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	r.GET("/api/products", handlers.GetProducts)
	r.GET("/api/products/suggest", handlers.GetProductSuggestions)
	r.GET("/api/products/top-selling", handlers.GetTopSelling)
	r.GET("/api/products/:id", handlers.GetProduct)
	r.GET("/api/categories", handlers.GetCategories)
	r.GET("/api/config", handlers.GetConfig)
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

	r.GET("/api/orders", handlers.GetOrders)
	r.POST("/api/orders", handlers.CreateOrder)
	r.GET("/api/orders/:id", handlers.GetOrder)

	r.GET("/api/payment-methods", handlers.GetPaymentMethods)
	r.GET("/api/contact-methods", handlers.GetContactMethods)
	r.POST("/api/orders/:id/payment-proof", middleware.RequireAuth(), handlers.UploadPaymentProof)
	r.GET("/api/orders/:id/payment-proof", middleware.RequireAuth(), handlers.GetOrderPaymentProof)

	r.GET("/api/notifications", middleware.RequireAuth(), handlers.GetNotifications)
	r.GET("/api/notifications/unread-count", middleware.RequireAuth(), handlers.GetUnreadCount)
	r.GET("/api/notifications/:id", middleware.RequireAuth(), handlers.GetNotification)
	r.PUT("/api/notifications/:id/read", middleware.RequireAuth(), handlers.MarkAsRead)
	r.PUT("/api/notifications/read-all", middleware.RequireAuth(), handlers.MarkAllRead)
	r.DELETE("/api/notifications/:id", middleware.RequireAuth(), handlers.DeleteNotification)
	r.DELETE("/api/notifications/read", middleware.RequireAuth(), handlers.DeleteReadNotifications)

	r.POST("/api/chat/conversations", middleware.RequireAuth(), handlers.CreateConversation)
	r.GET("/api/chat/conversations", middleware.RequireAuth(), handlers.GetConversations)
	r.GET("/api/chat/conversations/:id/messages", middleware.RequireAuth(), handlers.GetMessages)
	r.POST("/api/chat/conversations/:id/messages", middleware.RequireAuth(), handlers.SendMessage)
	r.PUT("/api/chat/conversations/:id/close", middleware.RequireAuth(), handlers.CustomerCloseConversation)
	r.POST("/api/chat/conversations/:id/rate", middleware.RequireAuth(), handlers.RateConversation)
	r.GET("/api/chat/unread-count", middleware.RequireAuth(), handlers.GetChatUnreadCount)
	r.POST("/api/chat/upload-image", middleware.RequireAuth(), handlers.UploadChatImage)
	r.GET("/api/chat/service-status", handlers.GetServiceStatus)
	r.GET("/api/chat/ws", middleware.RequireAuth(), handlers.HandleCustomerWS)

	r.GET("/api/auth/me", handlers.GetProfile)
	r.PUT("/api/auth/profile", handlers.UpdateProfile)
	r.GET("/api/auth/captcha", handlers.GetCaptcha)
	r.PUT("/api/auth/change-password", middleware.RequireAuth(), handlers.ChangePassword)

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
		admin.GET("/chat/ws", handlers.HandleAdminWS)
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

		superAdmin.GET("/payment-methods", handlers.GetAdminPaymentMethods)
		superAdmin.POST("/payment-methods", handlers.CreatePaymentMethod)
		superAdmin.PUT("/payment-methods/:id", handlers.UpdatePaymentMethod)
		superAdmin.DELETE("/payment-methods/:id", handlers.DeletePaymentMethod)

		superAdmin.GET("/contact-methods", handlers.GetAdminContactMethods)
		superAdmin.POST("/contact-methods", handlers.CreateContactMethod)
		superAdmin.PUT("/contact-methods/:id", handlers.UpdateContactMethod)
		superAdmin.DELETE("/contact-methods/:id", handlers.DeleteContactMethod)
	}

	log.Printf("Server running at http://localhost:%s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
