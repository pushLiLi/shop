package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

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
	DB, err = gorm.Open(mysqlConfig, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var sqlDB *sql.DB
	sqlDB, err = DB.DB()
	if err == nil {
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(4 * time.Minute)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(25)
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
		&models.Conversation{},
		&models.Message{},
		&models.QuickReply{},
		&models.Rating{},
		&models.PaymentMethod{},
		&models.PaymentProof{},
		&models.ContactMethod{},
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
	SeedBulkOrders()
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
		{ConfigKey: "home_promo_left_image", ConfigValue: ""},
		{ConfigKey: "home_promo_left_link", ConfigValue: ""},
		{ConfigKey: "home_promo_right_image", ConfigValue: ""},
		{ConfigKey: "home_promo_right_link", ConfigValue: ""},
		{ConfigKey: "home_featured_title", ConfigValue: "特别推荐"},
		{ConfigKey: "home_new_title", ConfigValue: "新品上架"},
		{ConfigKey: "home_topselling_title", ConfigValue: "热销排行"},
		{ConfigKey: "site_title", ConfigValue: "BYCIGAR | 权威正品雪茄在线购买商城"},
		{ConfigKey: "site_meta_description", ConfigValue: "BYCIGAR是中国领先的雪茄文化与在线购物平台。我们提供最新、最专业的雪茄测评、品牌新闻与养护知识，并为您甄选全球优质雪茄及配件，支持便捷在线购买。加入我们的雪茄社区，探索醇香世界。"},
		{ConfigKey: "favicon_url", ConfigValue: "/favicon.png"},
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

func SeedBulkOrders() {
	var products []models.Product
	DB.Where("is_active = ?", true).Find(&products)
	if len(products) == 0 {
		log.Println("SeedBulkOrders: 没有活跃商品，跳过订单生成")
		return
	}

	var users []models.User
	DB.Where("role = ?", "customer").Find(&users)
	for len(users) < 2 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
		newUser := models.User{
			Email:    fmt.Sprintf("testuser%d@test.com", len(users)+1),
			Password: string(hashedPassword),
			Name:     fmt.Sprintf("测试用户%d", len(users)+1),
			Role:     "customer",
		}
		DB.Create(&newUser)
		users = append(users, newUser)
	}

	var addresses []models.Address
	DB.Find(&addresses)
	for _, user := range users {
		hasAddress := false
		for _, a := range addresses {
			if a.UserID == user.ID {
				hasAddress = true
				break
			}
		}
		if !hasAddress {
			addr := models.Address{
				UserID:       user.ID,
				FullName:     user.Name,
				AddressLine1: fmt.Sprintf("测试地址 %d", user.ID),
				City:         "北京",
				State:        "北京",
				ZipCode:      "100000",
				Phone:        "13800138000",
				IsDefault:    true,
			}
			DB.Create(&addr)
			addresses = append(addresses, addr)
		}
	}

	var count int64
	DB.Model(&models.Order{}).Count(&count)
	const targetCount = 3000
	if count >= targetCount {
		return
	}

	statuses := []string{"pending", "processing", "shipped", "completed", "cancelled"}
	now := time.Now()
	r := rand.New(rand.NewSource(now.UnixNano()))

	needToCreate := targetCount - int(count)
	for i := 0; i < needToCreate; i++ {
		status := statuses[r.Intn(len(statuses))]
		user := users[r.Intn(len(users))]

		var userAddresses []models.Address
		for _, a := range addresses {
			if a.UserID == user.ID {
				userAddresses = append(userAddresses, a)
			}
		}
		if len(userAddresses) == 0 {
			userAddresses = addresses
		}
		address := userAddresses[r.Intn(len(userAddresses))]

		itemCount := r.Intn(5) + 1
		var total float64
		items := make([]models.OrderItem, itemCount)

		for j := 0; j < itemCount; j++ {
			product := products[r.Intn(len(products))]
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
			UserID:    user.ID,
			AddressID: address.ID,
			Total:     total,
			Status:    status,
			Remark:    fmt.Sprintf("测试订单 %d", count+int64(i)+1),
			Items:     items,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		}
		DB.Create(&order)

		if (i+1)%500 == 0 {
			log.Printf("已生成 %d / %d 订单", i+1, needToCreate)
		}
	}
	log.Printf("测试订单数据生成完成，共 %d 笔", targetCount)
}
