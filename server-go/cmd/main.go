package main

import (
	"log"

	"bycigar-server/internal/config"
	"bycigar-server/internal/database"
	"bycigar-server/internal/handlers"
	"bycigar-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()
	database.Migrate()
	database.Seed()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.Static("/static", "./static")

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	r.GET("/api/products", handlers.GetProducts)
	r.GET("/api/products/:id", handlers.GetProduct)
	r.GET("/api/categories", handlers.GetCategories)
	r.GET("/api/config", handlers.GetConfig)
	r.GET("/api/banners", handlers.GetBanners)

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

	r.PUT("/api/admin/config/:key", handlers.UpdateConfig)

	admin := r.Group("/api/admin")
	admin.Use(middleware.AdminOnly)
	{
		admin.POST("/upload", handlers.UploadImage)

		admin.GET("/products", handlers.GetAdminProducts)
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)

		admin.GET("/categories", handlers.GetAdminCategories)
		admin.POST("/categories", handlers.CreateCategory)
		admin.PUT("/categories/:id", handlers.UpdateCategory)
		admin.DELETE("/categories/:id", handlers.DeleteCategory)

		admin.GET("/banners", handlers.GetBanners)
		admin.POST("/banners", handlers.CreateBanner)
		admin.PUT("/banners/:id", handlers.UpdateBanner)
		admin.DELETE("/banners/:id", handlers.DeleteBanner)
	}

	log.Printf("Server running at http://localhost:%s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
