package database

import (
	"database/sql"
	"fmt"
	"log"

	"bycigar-server/internal/config"
	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	mysqlConfig := mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         255,
		SkipInitializeWithVersion: false,
	})
	DB, err = gorm.Open(mysqlConfig, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var sqlDB *sql.DB
	sqlDB, err = DB.DB()
	if err == nil {
		sqlDB.SetConnMaxLifetime(0)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}

	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.CartItem{},
		&models.Favorite{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.SiteConfig{},
		&models.Banner{},
		&models.Page{},
		&models.Setting{},
		&models.Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}

func Seed() {
	SeedAdminUser()
	SeedPages()
	SeedSettings()
	SeedSiteConfig()
	SeedTestData()
}

func SeedAdminUser() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := models.User{
		Email:    "admin@admin.com",
		Password: string(hashedPassword),
		Name:     "Admin",
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("Failed to create admin user: %v", err)
		return
	}

	log.Println("Admin user created successfully")
}

func SeedPages() {
	defaultPages := []models.Page{
		{Slug: "about", Title: "关于我们", Content: "# 关于我们\n\n欢迎使用 BYCIGAR，我们致力于为您提供最优质的雪茄产品和服务。"},
		{Slug: "services", Title: "服务条款", Content: "# 服务条款\n\n使用本网站即表示您同意以下服务条款..."},
		{Slug: "privacy-policy", Title: "隐私政策", Content: "# 隐私政策\n\n我们重视您的隐私，以下是我们的隐私政策..."},
		{Slug: "statement", Title: "严正声明", Content: "# 严正声明\n\n本网站所有商品均为正品，假一赔十..."},
	}

	for _, page := range defaultPages {
		var existing models.Page
		if err := DB.Where("slug = ?", page.Slug).First(&existing).Error; err != nil {
			DB.Create(&page)
		}
	}
}

func SeedSettings() {
	defaultSettings := []models.Setting{
		{Key: "footer_description", Value: "BYCIGAR是中国领先的雪茄文化与在线购物平台。我们提供最新、最专业的雪茄测评、品牌新闻与养护知识，并为您甄选全球优质雪茄及配件，支持便捷在线购买。加入我们的雪茄社区，探索醇香世界。"},
		{Key: "footer_service_time", Value: "客服在线时间每周一至周六 9:00到18:00"},
	}

	for _, setting := range defaultSettings {
		var existing models.Setting
		if err := DB.Where("key = ?", setting.Key).First(&existing).Error; err != nil {
			DB.Create(&setting)
		}
	}
}

func SeedSiteConfig() {
	defaultConfigs := []models.SiteConfig{
		{ConfigKey: "home_banner_1", ConfigValue: "/media/bycigar/banner-1.png"},
		{ConfigKey: "home_banner_2", ConfigValue: "/media/bycigar/banner-2.png"},
		{ConfigKey: "home_banner_3", ConfigValue: "/media/bycigar/banner-3.png"},
	}

	for _, cfg := range defaultConfigs {
		var existing models.SiteConfig
		if err := DB.Where("config_key = ?", cfg.ConfigKey).First(&existing).Error; err != nil {
			DB.Create(&cfg)
		}
	}
}

func BackfillOrderNo() {
	var orders []models.Order
	DB.Where("order_no = '' OR order_no IS NULL").Find(&orders)
	for _, order := range orders {
		DB.Model(&order).Update("order_no", utils.GenerateOrderNo())
	}
	if len(orders) > 0 {
		log.Printf("Backfilled order_no for %d orders", len(orders))
	}
}
