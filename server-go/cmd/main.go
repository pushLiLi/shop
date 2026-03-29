package main

import (
	"log"

	_ "bycigar-server/docs"
	"bycigar-server/internal/config"
	"bycigar-server/internal/database"
	"bycigar-server/internal/handlers"
	"bycigar-server/internal/middleware"
	miniopkg "bycigar-server/pkg/minio"
	"bycigar-server/pkg/utils"

	"github.com/gin-gonic/gin"
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
	config.LoadConfig()
	database.Connect()
	database.Migrate()
	database.Seed()
	miniopkg.InitMinio()
	miniopkg.EnsureBucket(config.AppConfig.MinioBucket)
	utils.InitSnowflake(1)
	database.BackfillOrderNo()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.Static("/static", "./static")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json"), ginSwagger.DefaultModelsExpandDepth(-1)))

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	r.GET("/api/products", handlers.GetProducts)
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

		admin.GET("/banners", handlers.GetAdminBanners)
		admin.POST("/banners", handlers.CreateBanner)
		admin.PUT("/banners/:id", handlers.UpdateBanner)
		admin.DELETE("/banners/:id", handlers.DeleteBanner)

		admin.GET("/pages", handlers.GetAdminPages)
		admin.PUT("/pages/:slug", handlers.UpdatePage)

		admin.PUT("/config/:key", handlers.UpdateConfig)
		admin.PUT("/settings/:key", handlers.UpdateSetting)

		admin.GET("/orders", handlers.GetAdminOrders)
		admin.GET("/orders/:id", handlers.GetAdminOrder)
		admin.PUT("/orders/:id/status", handlers.UpdateOrderStatus)

		admin.GET("/dashboard/stats", handlers.GetDashboardStats)
		admin.GET("/dashboard/recent-orders", handlers.GetDashboardRecentOrders)
		admin.GET("/dashboard/low-stock", handlers.GetDashboardLowStock)
		admin.GET("/dashboard/top-products", handlers.GetDashboardTopProducts)

		admin.GET("/users", handlers.GetAdminUsers)
		admin.GET("/users/:id", handlers.GetAdminUser)
		admin.PUT("/users/:id/role", handlers.UpdateUserRole)
	}

	log.Printf("Server running at http://localhost:%s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
