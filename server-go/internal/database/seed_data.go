package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

var seedUsers []models.User
var seedCategories []models.Category
var seedProducts []models.Product
var seedAddresses []models.Address
var seedPaymentMethods []models.PaymentMethod
var seedContactMethods []models.ContactMethod

func SeedTestData() {
	log.Println("Clearing existing data for re-seeding...")
	DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	DB.Exec("TRUNCATE TABLE order_items")
	DB.Exec("TRUNCATE TABLE orders")
	DB.Exec("TRUNCATE TABLE favorites")
	DB.Exec("TRUNCATE TABLE cart_items")
	DB.Exec("TRUNCATE TABLE addresses")
	DB.Exec("TRUNCATE TABLE products")
	DB.Exec("TRUNCATE TABLE categories")
	DB.Exec("TRUNCATE TABLE banners")
	DB.Exec("TRUNCATE TABLE conversations")
	DB.Exec("TRUNCATE TABLE messages")
	DB.Exec("TRUNCATE TABLE notifications")
	DB.Exec("TRUNCATE TABLE payment_proofs")
	DB.Exec("TRUNCATE TABLE payment_methods")
	DB.Exec("TRUNCATE TABLE contact_methods")
	DB.Exec("TRUNCATE TABLE settings")
	DB.Exec("TRUNCATE TABLE pages")
	DB.Exec("TRUNCATE TABLE configs")
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("Seeding test data...")
	utils.InitSnowflake(1)

	seedUsersData()
	seedCategoriesData()
	seedProductsData()
	seedBannersData()
	seedPaymentMethodsData()
	seedSettingsData()
	seedPagesData()
	seedConfigsData()
	seedAddressesData()
	seedOrdersData()
	seedPaymentProofsData()
	seedCartAndFavorites()
	seedConversationsAndMessages()
	seedNotificationsData()
	seedContactMethodsData()

	log.Println("Test data seeded successfully")
}

func seedUsersData() {
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	admins := []struct {
		email string
		name  string
		role  string
	}{
		{"admin@bycigar.com", "系统管理员", "admin"},
		{"service1@bycigar.com", "客服小王", "service"},
		{"service2@bycigar.com", "客服小李", "service"},
	}

	for _, u := range admins {
		var existing models.User
		if err := DB.Where("email = ?", u.email).First(&existing).Error; err != nil {
			user := models.User{Email: u.email, Password: string(password), Name: u.name, Role: u.role}
			DB.Create(&user)
			seedUsers = append(seedUsers, user)
		} else {
			seedUsers = append(seedUsers, existing)
		}
	}

	customerNames := []string{
		"张伟", "李娜", "王芳", "刘洋", "陈杰", "杨秀英", "黄明", "赵磊", "周静", "吴强",
		"徐丽", "孙浩", "马超", "胡蝶", "朱琳", "郭锐", "林涛", "何雪", "高建", "罗岚",
		"梁志明", "宋雨晴", "郑浩", "谢辉", "田亮",
	}

	daysAgo := 90
	for i, name := range customerNames {
		email := fmt.Sprintf("user%02d@test.com", i+1)
		var existing models.User
		if err := DB.Where("email = ?", email).First(&existing).Error; err == nil {
			seedUsers = append(seedUsers, existing)
			continue
		}
		user := models.User{
			Email:    email,
			Password: string(password),
			Name:     name,
			Role:     "customer",
		}
		DB.Create(&user)
		seedUsers = append(seedUsers, user)

		createdAt := time.Now().AddDate(0, 0, -rand.Intn(daysAgo+1))
		DB.Model(&user).Update("created_at", createdAt)
	}

	log.Printf("Users seeded: %d (admins=%d, customers=%d)", len(seedUsers), 3, len(customerNames))
}

func seedCategoriesData() {
	topCategories := []struct {
		name string
		slug string
	}{
		{"精品雪茄", "premium-cigars"},
		{"雪茄配件", "accessories"},
		{"生活方式", "lifestyle"},
		{"礼盒套装", "gift-sets"},
	}

	var createdTop []models.Category
	for _, c := range topCategories {
		var cat models.Category
		if err := DB.Where("slug = ?", c.slug).First(&cat).Error; err != nil {
			cat = models.Category{Name: c.name, Slug: c.slug}
			DB.Create(&cat)
		}
		createdTop = append(createdTop, cat)
	}

	premiumID := createdTop[0].ID
	accessoriesID := createdTop[1].ID
	lifestyleID := createdTop[2].ID
	giftSetsID := createdTop[3].ID

	subCategories := []struct {
		name     string
		slug     string
		parentID uint
	}{
		{"古巴经典", "cuba-classic", premiumID},
		{"多米尼加", "dominican", premiumID},
		{"尼加拉瓜", "nicaragua", premiumID},
		{"切割工具", "cutters", accessoriesID},
		{"保湿存储", "humidors", accessoriesID},
		{"点火设备", "lighters", accessoriesID},
		{"酒水搭配", "spirits-pairing", lifestyleID},
		{"入门礼盒", "starter-kits", giftSetsID},
	}

	for _, sc := range subCategories {
		var cat models.Category
		if err := DB.Where("slug = ?", sc.slug).First(&cat).Error; err != nil {
			cat = models.Category{Name: sc.name, Slug: sc.slug, ParentID: &sc.parentID}
			DB.Create(&cat)
		}
		seedCategories = append(seedCategories, cat)
	}
	seedCategories = append(seedCategories, createdTop...)

	log.Printf("Categories seeded: %d", len(seedCategories))
}

func catIDBySlug(slug string) uint {
	var cat models.Category
	DB.Where("slug = ?", slug).First(&cat)
	return cat.ID
}

func seedProductsData() {
	cuba := catIDBySlug("cuba-classic")
	dominican := catIDBySlug("dominican")
	nicaragua := catIDBySlug("nicaragua")
	cutters := catIDBySlug("cutters")
	humidors := catIDBySlug("humidors")
	lighters := catIDBySlug("lighters")
	spirits := catIDBySlug("spirits-pairing")
	starter := catIDBySlug("starter-kits")

	products := []models.Product{
		{Name: "高希霸世纪一号", Slug: "cohiba-siglo-i", Price: 128, CategoryID: cuba, Stock: 50, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba1/400/400", Description: "高希霸世纪系列入门款，口感温和细腻。"},
		{Name: "高希霸世纪二号", Slug: "cohiba-siglo-ii", Price: 158, CategoryID: cuba, Stock: 35, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohiba2/400/400", Description: "中等浓郁的世纪二号，带有奶油和咖啡的香气。"},
		{Name: "高希霸世纪三号", Slug: "cohiba-siglo-iii", Price: 198, CategoryID: cuba, Stock: 0, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohiba3/400/400", Description: "浓郁度更高，目前缺货。"},
		{Name: "高希霸世纪四号", Slug: "cohiba-siglo-iv", Price: 258, CategoryID: cuba, Stock: 22, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba4/400/400", Description: "世纪系列中最受欢迎的型号，口感丰富饱满。"},
		{Name: "高希霸世纪五号", Slug: "cohiba-siglo-v", Price: 328, CategoryID: cuba, Stock: 15, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba5/400/400", Description: "世纪系列旗舰款，复杂多变的风味层次。"},
		{Name: "高希霸短号", Slug: "cohiba-short", Price: 68, CategoryID: cuba, Stock: 120, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohibas/400/400", Description: "短小精悍的日常雪茄，适合短暂休憩。"},
		{Name: "高希霸马杜罗5号魔术师", Slug: "cohiba-maduro-5", Price: 388, CategoryID: cuba, Stock: 12, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohibam5/400/400", Description: "深色马杜罗茄衣，浓郁甜蜜。"},
		{Name: "高希霸贝伊可52号", Slug: "cohiba-behike-52", Price: 680, CategoryID: cuba, Stock: 3, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/bk52/400/400", Description: "贝伊可系列，极其稀有，低库存。"},
		{Name: "高希霸长矛", Slug: "cohiba-lancero", Price: 218, CategoryID: cuba, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/cohibal/400/400", Description: "已下架的经典长矛款。"},
		{Name: "科伊巴鱼雷限量版", Slug: "cohiba-torpedo-limited", Price: 888, CategoryID: cuba, Stock: 0, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohibator/400/400", Description: "限量版鱼雷，极高品质。"},
		{Name: "蒙特2号", Slug: "montecristo-no2", Price: 108, CategoryID: cuba, Stock: 45, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/monte2/400/400", Description: "蒙特最经典的鱼雷型号，全球畅销。"},
		{Name: "蒙特4号", Slug: "montecristo-no4", Price: 78, CategoryID: cuba, Stock: 80, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/monte4/400/400", Description: "世界上最畅销的雪茄之一。"},
		{Name: "蒙特埃德蒙多", Slug: "montecristo-edmundo", Price: 138, CategoryID: cuba, Stock: 30, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/monteed/400/400", Description: "丰富的层次感，中等偏浓郁。"},
		{Name: "蒙特双埃德蒙多", Slug: "montecristo-double", Price: 168, CategoryID: cuba, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/montede/400/400", Description: "加粗版埃德蒙多，更长品吸时间。"},
		{Name: "蒙特OPEN初级", Slug: "montecristo-open-junior", Price: 58, CategoryID: cuba, Stock: 60, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/montejr/400/400", Description: "入门级蒙特，温和易入口。"},
		{Name: "帕塔加斯D4号", Slug: "partagas-d4", Price: 98, CategoryID: cuba, Stock: 40, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/partagas4/400/400", Description: "帕塔加斯最经典的罗布斯托。"},
		{Name: "帕塔加斯D6号", Slug: "partagas-d6", Price: 128, CategoryID: cuba, Stock: 25, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/partagas6/400/400", Description: "浓郁的泥土和香料风味。"},
		{Name: "帕塔加斯卢西塔尼亚", Slug: "partagas-lusitanias", Price: 358, CategoryID: cuba, Stock: 5, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/partagasl/400/400", Description: "大尺寸双皇冠，低库存限量。"},
		{Name: "帕塔加斯超级皇冠", Slug: "partagas-super-corona", Price: 198, CategoryID: cuba, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/partagassc/400/400", Description: "已下架。"},
		{Name: "罗密欧2号", Slug: "romeo-no2", Price: 88, CategoryID: cuba, Stock: 55, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/romeo2/400/400", Description: "经典鱼雷款，浪漫之名。"},
		{Name: "罗密欧宽丘吉尔", Slug: "romeo-wide-churchill", Price: 178, CategoryID: cuba, Stock: 18, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeowc/400/400", Description: "宽环规丘吉尔，品吸时间充裕。"},
		{Name: "罗密欧短丘吉尔", Slug: "romeo-short-churchill", Price: 108, CategoryID: cuba, Stock: 35, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeosc/400/400", Description: "短丘吉尔，适合午间休息。"},
		{Name: "罗密欧俱乐部", Slug: "romeo-club", Price: 48, CategoryID: cuba, Stock: 0, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeoclub/400/400", Description: "缺货的小俱乐部款。"},
		{Name: "大卫杜夫2000", Slug: "davidoff-2000", Price: 198, CategoryID: dominican, Stock: 30, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/dav2000/400/400", Description: "瑞士精工品质，细腻顺滑。"},
		{Name: "大卫杜夫千年系列", Slug: "davidoff-millennium", Price: 298, CategoryID: dominican, Stock: 12, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/davmil/400/400", Description: "千年系列，浓郁的香料和咖啡。"},
		{Name: "大卫杜夫温斯顿丘吉尔", Slug: "davidoff-wsc", Price: 458, CategoryID: dominican, Stock: 8, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/davwc/400/400", Description: "致敬伟人的顶级系列。"},
		{Name: "大卫杜夫埃斯库里奥", Slug: "davidoff-escurio", Price: 168, CategoryID: dominican, Stock: 25, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/davesc/400/400", Description: "巴西茄叶，甜美辛辣。"},
		{Name: "大卫杜夫格兰德", Slug: "davidoff-grande", Price: 388, CategoryID: dominican, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/davgr/400/400", Description: "已下架。"},
		{Name: "富恩特唐卡洛斯", Slug: "fuente-don-carlos", Price: 228, CategoryID: dominican, Stock: 15, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/fuentedc/400/400", Description: "多米尼加之光，手工精选。"},
		{Name: "富恩特OpusX", Slug: "fuente-opus-x", Price: 688, CategoryID: dominican, Stock: 4, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/fuentox/400/400", Description: "传说中的OpusX，极低库存。"},
		{Name: "富恩特海明威经典", Slug: "fuente-hemingway", Price: 188, CategoryID: dominican, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/fuentehw/400/400", Description: "完美造型，大师级卷制。"},
		{Name: "盛赛迪亚马杜罗", Slug: "盛赛迪亚-maduro", Price: 258, CategoryID: dominican, Stock: 12, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed盛赛迪亚/400/400", Description: "深色马杜罗风格。"},
		{Name: "帕德龙1964周年", Slug: "padron-1964", Price: 278, CategoryID: nicaragua, Stock: 18, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/padron64/400/400", Description: "尼加拉瓜经典，周年纪念系列。"},
		{Name: "帕德龙1926系列80年", Slug: "padron-1926-80", Price: 488, CategoryID: nicaragua, Stock: 6, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/padron26/400/400", Description: "80周年纪念款，极高品质。"},
		{Name: "帕德龙家族Reserve", Slug: "padron-family-reserve", Price: 568, CategoryID: nicaragua, Stock: 2, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/padronfr/400/400", Description: "家族珍藏，极低库存。"},
		{Name: "帕德龙大师系列", Slug: "padron-master", Price: 328, CategoryID: nicaragua, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/padronm/400/400", Description: "已下架。"},
		{Name: "奥利瓦V系列", Slug: "oliva-v-series", Price: 198, CategoryID: nicaragua, Stock: 22, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/olivav/400/400", Description: "尼加拉瓜名庄奥利瓦代表作。"},
		{Name: "AJ费尔南德斯新世界", Slug: "aj-fernandez-nw", Price: 168, CategoryID: nicaragua, Stock: 15, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/ajfw/400/400", Description: "新世界风格，浓郁饱满。"},
		{Name: "雪茄剪双刃不锈钢", Slug: "cutter-stainless", Price: 128, CategoryID: cutters, Stock: 100, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cutter1/400/400", Description: "高品质不锈钢双刃雪茄剪。"},
		{Name: "雪茄剪V口切割器", Slug: "cutter-v-cut", Price: 88, CategoryID: cutters, Stock: 80, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cutter2/400/400", Description: "V型切口，完美品吸体验。"},
		{Name: "雪茄钻孔器", Slug: "cutter-punch", Price: 58, CategoryID: cutters, Stock: 150, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/punch1/400/400", Description: "便捷雪茄钻孔器，随身携带。"},
		{Name: "专业雪茄剪套装", Slug: "cutter-pro-set", Price: 298, CategoryID: cutters, Stock: 10, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cutterpro/400/400", Description: "专业级套装，含剪刀和V口器。"},
		{Name: "电动雪茄剪", Slug: "cutter-electric", Price: 388, CategoryID: cutters, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/cuttere/400/400", Description: "已下架。"},
		{Name: "桌面保湿盒50支", Slug: "humidor-50", Price: 388, CategoryID: humidors, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/humidor50/400/400", Description: "经典桃花心木桌面保湿盒。"},
		{Name: "旅行保湿盒5支", Slug: "humidor-travel-5", Price: 168, CategoryID: humidors, Stock: 35, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/humidort5/400/400", Description: "便携旅行装，密封防干。"},
		{Name: "豪华保湿柜200支", Slug: "humidor-cabinet-200", Price: 2800, CategoryID: humidors, Stock: 3, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/humidor200/400/400", Description: "顶级豪华展示柜，电子恒温恒湿。"},
		{Name: "电子湿度计", Slug: "hygrometer-digital", Price: 88, CategoryID: humidors, Stock: 50, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/hygrometer/400/400", Description: "精准电子湿度温度计。"},
		{Name: "保湿包套装", Slug: "humidor-pack-set", Price: 68, CategoryID: humidors, Stock: 0, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/humidorpack/400/400", Description: "缺货中。"},
		{Name: "雪松木保湿内衬", Slug: "humidor-cedar-liner", Price: 128, CategoryID: humidors, Stock: 40, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cedarliner/400/400", Description: "雪松木内衬，增强保湿效果。"},
		{Name: "气点火器", Slug: "lighter-butane", Price: 158, CategoryID: lighters, Stock: 60, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/lightbut/400/400", Description: "丁烷气点火器，可调节火焰。"},
		{Name: "松木火柴", Slug: "lighter-wood-matches", Price: 28, CategoryID: lighters, Stock: 200, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/matchwood/400/400", Description: "天然松木火柴，品茄专用。"},
		{Name: "丁烷火枪", Slug: "lighter-torch", Price: 198, CategoryID: lighters, Stock: 45, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/torch/400/400", Description: "双焰火枪，点燃方便。"},
		{Name: "Xikar点火器", Slug: "lighter-xikar", Price: 288, CategoryID: lighters, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/xikar/400/400", Description: "已下架。"},
		{Name: "威士忌杯套装", Slug: "whisky-glass-set", Price: 168, CategoryID: spirits, Stock: 30, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/whiskyglass/400/400", Description: "手工切割水晶威士忌杯套装。"},
		{Name: "白兰地杯", Slug: "brandy-glass", Price: 128, CategoryID: spirits, Stock: 25, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/brandyglass/400/400", Description: "经典白兰地杯，品茄佳配。"},
		{Name: "朗姆酒精选", Slug: "rum-premium", Price: 298, CategoryID: spirits, Stock: 15, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/rumprem/400/400", Description: "陈年朗姆酒，搭配雪茄绝配。"},
		{Name: "雪茄核桃木托盘", Slug: "cigar-walnut-tray", Price: 388, CategoryID: spirits, Stock: 8, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/walnuttray/400/400", Description: "胡桃木托盘，兼具美观与实用。"},
		{Name: "入门体验装5支", Slug: "starter-pack-5", Price: 198, CategoryID: starter, Stock: 25, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/beginpack/400/400", Description: "精选5支入门雪茄套装，附品鉴指南。"},
		{Name: "高希霸品鉴礼盒", Slug: "cohiba-tasting-box", Price: 666, CategoryID: starter, Stock: 10, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohibabox/400/400", Description: "含世纪1-5号各一支的豪华品鉴礼盒。"},
		{Name: "旅行雪茄套装", Slug: "travel-cigar-set", Price: 358, CategoryID: starter, Stock: 18, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/travelset/400/400", Description: "含便携保湿盒和5支精选雪茄。"},
		{Name: "送礼豪华礼盒", Slug: "gift-luxury-box", Price: 1288, CategoryID: starter, Stock: 5, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/giftlux/400/400", Description: "高端送礼首选，附精美包装。"},
		{Name: "限量版年度礼盒", Slug: "limited-annual-box", Price: 2588, CategoryID: starter, Stock: 2, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/annualbox/400/400", Description: "年度限量版，极具收藏价值。"},
		{Name: "节日特别版套装", Slug: "festival-special-set", Price: 988, CategoryID: starter, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/festset/400/400", Description: "已下架。"},
		{Name: "已下架测试商品", Slug: "discontinued-test", Price: 99, CategoryID: cuba, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/discont/400/400", Description: "已下架。"},
	}

	for i := range products {
		var existing models.Product
		if err := DB.Where("slug = ?", products[i].Slug).First(&existing).Error; err != nil {
			DB.Create(&products[i])
		} else {
			products[i].ID = existing.ID
		}
	}

	seedProducts = products
	log.Printf("Products seeded: %d", len(seedProducts))
}

func seedBannersData() {
	banners := []models.Banner{
		{Title: "古巴经典 · 传承百年", Image: "https://picsum.photos/seed/banner-cuba/1200/400", Link: "/category/cuba-classic", SortOrder: 1, IsActive: true},
		{Title: "多米尼加风情 · 细腻优雅", Image: "https://picsum.photos/seed/banner-dom/1200/400", Link: "/category/dominican", SortOrder: 2, IsActive: true},
		{Title: "尼加拉瓜激情 · 浓郁澎湃", Image: "https://picsum.photos/seed/banner-nic/1200/400", Link: "/category/nicaragua", SortOrder: 3, IsActive: true},
		{Title: "配件专区 · 点亮品茄时刻", Image: "https://picsum.photos/seed/banner-acc/1200/400", Link: "/category/cutters", SortOrder: 4, IsActive: true},
		{Title: "生活方式 · 雪茄与美酒", Image: "https://picsum.photos/seed/banner-life/1200/400", Link: "/category/spirits-pairing", SortOrder: 5, IsActive: true},
		{Title: "送礼佳选 · 礼盒套装", Image: "https://picsum.photos/seed/banner-gift/1200/400", Link: "/category/starter-kits", SortOrder: 6, IsActive: true},
	}

	for i := range banners {
		var existing models.Banner
		if err := DB.Where("title = ?", banners[i].Title).First(&existing).Error; err != nil {
			DB.Create(&banners[i])
		}
	}
	log.Printf("Banners seeded: %d", len(banners))
}

func seedPaymentMethodsData() {
	methods := []models.PaymentMethod{
		{Name: "微信支付", QRCodeUrl: "https://picsum.photos/seed/wxpay/200/200", Instructions: "请扫描二维码支付", IsActive: true, SortOrder: 1},
		{Name: "支付宝", QRCodeUrl: "https://picsum.photos/seed/alipay/200/200", Instructions: "请扫描二维码支付", IsActive: true, SortOrder: 2},
		{Name: "银行转账", QRCodeUrl: "", Instructions: "转账至指定银行账户", IsActive: true, SortOrder: 3},
		{Name: "微信支付(旧版)", QRCodeUrl: "https://picsum.photos/seed/wxpay-old/200/200", Instructions: "已停用，请使用新版", IsActive: false, SortOrder: 4},
		{Name: "支付宝(旧版)", QRCodeUrl: "https://picsum.photos/seed/alipay-old/200/200", Instructions: "已停用，请使用新版", IsActive: false, SortOrder: 5},
		{Name: "货到付款", QRCodeUrl: "", Instructions: "收货时付款", IsActive: true, SortOrder: 6},
		{Name: "PayPal", QRCodeUrl: "https://picsum.photos/seed/paypal/200/200", Instructions: "支持国际信用卡支付", IsActive: true, SortOrder: 7},
		{Name: "企业转账", QRCodeUrl: "", Instructions: "仅限企业客户", IsActive: false, SortOrder: 8},
	}

	for i := range methods {
		var existing models.PaymentMethod
		if err := DB.Where("name = ?", methods[i].Name).First(&existing).Error; err != nil {
			DB.Create(&methods[i])
		} else {
			methods[i].ID = existing.ID
		}
	}
	seedPaymentMethods = methods
	log.Printf("PaymentMethods seeded: %d", len(seedPaymentMethods))
}

func seedSettingsData() {
	settings := []struct {
		key   string
		value string
	}{
		{"site_name", "BYCIGAR 雪茄旗舰店"},
		{"chat_greeting", "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？"},
		{"contact_phone", "400-888-9999"},
		{"contact_email", "support@bycigar.com"},
		{"footer_content", "© 2024 BYCIGAR 版权所有"},
		{"order_auto_close_days", "7"},
		{"low_stock_threshold", "10"},
		{"maintenance_mode", "false"},
	}

	for _, s := range settings {
		var existing models.Setting
		if err := DB.Where("`key` = ?", s.key).First(&existing).Error; err != nil {
			DB.Create(&models.Setting{Key: s.key, Value: s.value})
		}
	}
	log.Printf("Settings seeded: %d", len(settings))
}

func seedPagesData() {
	pages := []struct {
		slug    string
		title   string
		content string
	}{
		{"about", "关于我们", "BYCIGAR 成立于2010年，是国内领先的高端雪茄零售商。我们致力于为雪茄爱好者提供最优质的古巴及世界顶级雪茄，产品涵盖高希霸、蒙特、帕塔加斯、大卫杜夫等经典品牌。公司总部位于上海，拥有专业化的恒温仓储物流体系，所有雪茄均经过严格品控，确保品质如一。\n\n我们的团队由一群热爱雪茄文化的专业人士组成，每位顾问都经过严格培训，能够为客户提供专业的品鉴建议和搭配推荐。BYCIGAR 还定期举办雪茄品鉴会，邀请国内外知名雪茄大师与会员互动交流。\n\n未来，我们将继续深耕雪茄文化，引入更多优质产品，为中国雪茄市场的发展贡献力量。"},
		{"services", "服务条款", "一、服务说明\n\nBYCIGAR 为您提供在线雪茄及配件销售服务，包括商品浏览、在线购买、支付结算、物流配送等完整流程。\n\n二、购买须知\n\n1. 雪茄产品需年满18周岁方可购买。\n2. 请如实填写收货信息，确保商品顺利送达。\n3. 雪茄属于特殊商品一经拆封不支持无理由退换。\n4. 如收到商品有质量问题，请在48小时内联系客服。\n\n三、物流配送\n\n我们合作的物流伙伴包括顺丰、京东等优质快递，全程冷链运输，确保雪茄在最佳状态下送达。部分地区可能需要身份证验证。\n\n四、售后服务\n\n客服热线：400-888-9999（工作日9:00-22:00）\n邮箱：support@bycigar.com"},
		{"privacy-policy", "隐私政策", "BYCIGAR 非常重视您的个人信息的保护。本隐私政策说明了我们在您使用我们的服务时如何收集、使用、存储和保护您的个人信息。\n\n一、信息收集\n\n当您注册账户时，我们会收集您的姓名、邮箱、手机号码等基本信息。当您下单时，我们会收集收货地址信息以完成配送。\n\n二、信息使用\n\n您的个人信息将用于：处理订单、提供客户服务、发送订单状态通知、推荐适合的商品、改进我们的服务。\n\n三、信息保护\n\n我们采用行业标准的加密技术保护您的数据安全，未经您的授权，我们不会将个人信息提供给任何第三方。\n\n四、联系我们\n\n如对隐私政策有任何疑问，请联系：privacy@bycigar.com"},
		{"statement", "免责说明", "一、产品声明\n\nBYCIGAR 销售的所有雪茄产品均为正品。我们合作的品牌均经过官方授权，或从正规渠道采购。\n\n二、健康提示\n\n吸烟有害健康。雪茄含有焦油和一氧化碳等有害物质，未成年人、孕妇、哺乳期妇女及心血管疾病患者不应使用雪茄。\n\n三、使用风险\n\n雪茄的保存需要特定温度（18-21°C）和湿度（68-72%）。因客户保存不当导致的雪茄品质问题，不在退换货范围内。\n\n四、配送风险\n\n商品在运输过程中可能因不可抗力因素导致延误或损坏，我们将配合物流公司积极处理，但不承担额外赔偿责任。"},
		{"shipping", "配送说明", "一、配送范围\n\n我们支持全国大部分地区的配送，包括港澳台地区。部分偏远地区可能需要更长的配送时间。\n\n二、配送方式\n\n1. 标准配送：3-5个工作日送达\n2. 加急配送：1-2个工作日送达（需额外支付费用）\n3. 定时配送：您可选择具体的收货时间段\n\n三、配送费用\n\n订单满500元免运费，不满500元收取20元运费。新会员首单免运费。\n\n四、配送包装\n\n所有雪茄均采用专业保湿包装，配合保温箱确保品质。礼盒订单额外使用精美纸箱包装。\n\n五、签收提示\n\n请在签收前检查包装是否完好，如有问题请当场拒收并联系客服。"},
	}

	for _, p := range pages {
		var existing models.Page
		if err := DB.Where("slug = ?", p.slug).First(&existing).Error; err != nil {
			DB.Create(&models.Page{Slug: p.slug, Title: p.title, Content: p.content})
		}
	}
	log.Printf("Pages seeded: %d", len(pages))
}

func seedConfigsData() {
	configs := []struct {
		key   string
		value string
	}{
		{"site_title", "BYCIGAR 雪茄商城"},
		{"seo_description", "专业古巴雪茄及高档雪茄配件销售平台，提供高希霸、蒙特、大卫杜夫等品牌雪茄，品类齐全，品质保证。"},
		{"customer_service_hours", "9:00-22:00"},
		{"max_cart_items", "99"},
		{"order_max_quantity", "50"},
		{"exchange_rate_usd", "7.25"},
	}

	for _, c := range configs {
		var existing models.SiteConfig
		if err := DB.Where("config_key = ?", c.key).First(&existing).Error; err != nil {
			DB.Create(&models.SiteConfig{ConfigKey: c.key, ConfigValue: c.value})
		}
	}
	log.Printf("Configs seeded: %d", len(configs))
}

func seedAddressesData() {
	if len(seedUsers) < 25 {
		log.Println("Not enough users for address seeding")
		return
	}

	cities := []struct {
		city  string
		state string
	}{
		{"北京", "北京"}, {"上海", "上海"}, {"广州", "广东"}, {"深圳", "广东"},
		{"杭州", "浙江"}, {"成都", "四川"}, {"武汉", "湖北"}, {"天津", "天津"},
		{"南京", "江苏"}, {"长沙", "湖南"}, {"西安", "陕西"}, {"厦门", "福建"},
		{"重庆", "重庆"}, {"苏州", "江苏"}, {"郑州", "河南"}, {"青岛", "山东"},
		{"沈阳", "辽宁"}, {"大连", "辽宁"}, {"济南", "山东"}, {"福州", "福建"},
		{"昆明", "云南"}, {"哈尔滨", "黑龙江"}, {"长春", "吉林"}, {"石家庄", "河北"},
		{"南昌", "江西"},
	}

	streets := []string{
		"建国路88号SOHO现代城", "南京东路100号", "天河路385号太古汇",
		"科技园南区深南大道9966号", "文三路553号", "红星路三段1号",
		"中北路109号", "南京路189号", "中央路331号",
		"麓谷大道627号", "长安南路86号", "厦禾路888号",
		"解放碑步行街88号", "观前街168号", "二七路72号",
		"香港中路100号", "五四路200号", "长江道1号",
		"东风路388号", "五一路399号", "人民路555号",
		"建设路188号", "中山路299号", "胜利北路66号",
		"滨海大道88号中心大厦",
	}

	phones := []string{
		"13800138001", "13900139002", "13700137003", "13600136004",
		"13500135005", "13400134006", "13300133007", "13200132008",
		"13100131009", "13000130010", "15000150011", "15100151012",
		"15200152013", "15300153014", "15500155015", "15600156016",
		"15700157017", "15800158018", "15900159019", "18000180020",
		"18100181021", "18200182022", "18300183023", "18400184024",
		"18500185025",
	}

	zips := []string{
		"100022", "200001", "510620", "518057", "310012", "610021",
		"430071", "300051", "210008", "410205", "710061", "361003",
		"400010", "215000", "450000", "266000", "110001", "116001",
		"250001", "350001", "650000", "150001", "130000", "050000",
		"330000",
	}

	var addresses []models.Address
	for i := 0; i < 25; i++ {
		userID := seedUsers[i+3].ID
		name := seedUsers[i+3].Name
		c := cities[i]
		street := streets[i]
		phone := phones[i]
		zip := zips[i]

		addr1 := models.Address{
			UserID: userID, FullName: name, AddressLine1: street,
			City: c.city, State: c.state, ZipCode: zip, Phone: phone, IsDefault: true,
		}
		DB.Create(&addr1)
		addresses = append(addresses, addr1)

		if i%3 == 0 {
			addr2 := models.Address{
				UserID: userID, FullName: name, AddressLine1: street + "B座",
				City: c.city, State: c.state, ZipCode: zip, Phone: phone, IsDefault: false,
			}
			DB.Create(&addr2)
			addresses = append(addresses, addr2)
		}
		if i%5 == 0 {
			addr3 := models.Address{
				UserID: userID, FullName: name, AddressLine1: street + "（仓库）",
				City: c.city, State: c.state, ZipCode: zip, Phone: phone, IsDefault: false,
			}
			DB.Create(&addr3)
			addresses = append(addresses, addr3)
		}
	}

	seedAddresses = addresses
	log.Printf("Addresses seeded: %d", len(addresses))
}

func seedOrdersData() {
	if len(seedUsers) < 25 || len(seedProducts) < 50 || len(seedAddresses) < 30 {
		log.Println("Not enough data for order seeding")
		return
	}

	type orderDef struct {
		userIdx   int
		addrIdx   int
		status    string
		remark    string
		itemIdxes []int
		qtys      []int
		daysAgo   int
		tracking  string
		trackingN string
	}

	defs := []orderDef{
		{0, 0, "pending", "", []int{0, 11}, []int{2, 1}, 0, "", ""},
		{1, 1, "pending", "请尽快发货", []int{4}, []int{1}, 0, "", ""},
		{2, 2, "pending", "", []int{20}, []int{3}, 1, "", ""},
		{3, 3, "pending", "需要发票", []int{32, 33}, []int{1, 1}, 1, "", ""},
		{4, 4, "pending", "", []int{42}, []int{2}, 2, "", ""},
		{5, 5, "pending", "小心轻放", []int{15}, []int{4}, 0, "", ""},
		{6, 6, "pending", "", []int{22, 12}, []int{1, 2}, 1, "", ""},
		{7, 7, "pending", "生日礼物", []int{28}, []int{1}, 0, "", ""},

		{8, 8, "paid", "已转账", []int{43}, []int{1}, 0, "", ""},
		{9, 9, "paid", "", []int{35, 36}, []int{1, 2}, 1, "", ""},
		{10, 10, "paid", "", []int{0, 5}, []int{1, 1}, 0, "", ""},
		{11, 11, "paid", "谢谢", []int{22}, []int{2}, 1, "", ""},
		{12, 12, "paid", "", []int{52}, []int{1}, 0, "", ""},

		{13, 13, "processing", "", []int{1, 10}, []int{1, 3}, 2, "", ""},
		{14, 14, "processing", "请仔细包装", []int{5}, []int{2}, 3, "顺丰速运", "SF1234567890"},
		{15, 15, "processing", "", []int{8, 14}, []int{1, 1}, 3, "京东快递", "JD0001234567"},
		{16, 16, "processing", "", []int{20, 21}, []int{2, 1}, 4, "顺丰速运", "SF9876543210"},
		{17, 17, "processing", "节假日不送", []int{28}, []int{1}, 4, "EMS", "EM000123456"},
		{18, 18, "processing", "", []int{40, 41}, []int{1, 2}, 5, "顺丰速运", "SF5555666677"},
		{19, 19, "processing", "办公室地址", []int{15, 16, 17}, []int{1, 1, 1}, 5, "京东快递", "JD888999000"},
		{20, 20, "processing", "", []int{22, 23}, []int{1, 1}, 6, "顺丰速运", "SF1111222233"},
		{21, 21, "processing", "尽快", []int{50}, []int{1}, 6, "德邦物流", "DB444455556"},
		{22, 22, "processing", "", []int{3, 4}, []int{1, 1}, 7, "顺丰速运", "SF7777888899"},

		{0, 0, "shipped", "", []int{6}, []int{2}, 8, "顺丰速运", "SF2222333344"},
		{1, 1, "shipped", "快递柜自提", []int{10, 14}, []int{1, 2}, 8, "丰巢快递", "FC123123123"},
		{2, 2, "shipped", "", []int{8}, []int{1}, 9, "京东快递", "JD5555666677"},
		{3, 3, "shipped", "", []int{0}, []int{5}, 9, "顺丰速运", "SF8888999900"},
		{4, 4, "shipped", "家人代收", []int{35, 36}, []int{1, 2}, 10, "EMS", "EM888999000"},
		{5, 5, "shipped", "", []int{43}, []int{1}, 10, "顺丰速运", "SF1111333355"},
		{6, 6, "shipped", "白天不在", []int{50, 51}, []int{1, 1}, 11, "京东快递", "JD2222333344"},
		{7, 7, "shipped", "", []int{52}, []int{2}, 11, "顺丰速运", "SF5555667788"},
		{8, 8, "shipped", "请敲门", []int{40}, []int{1}, 12, "德邦物流", "DB6666777788"},
		{9, 9, "shipped", "", []int{20, 21, 22}, []int{1, 1, 1}, 12, "顺丰速运", "SF9999000011"},
		{10, 10, "shipped", "", []int{1}, []int{3}, 13, "EMS", "EM111222333"},
		{11, 11, "shipped", "周末在家", []int{4, 5}, []int{1, 1}, 13, "京东快递", "JD444555666"},

		{12, 12, "completed", "非常满意", []int{0}, []int{2}, 20, "", ""},
		{13, 13, "completed", "", []int{10, 11}, []int{1, 3}, 25, "", ""},
		{14, 14, "completed", "", []int{4, 22}, []int{1, 1}, 28, "", ""},
		{15, 15, "completed", "好评", []int{25}, []int{2}, 30, "", ""},
		{16, 16, "completed", "", []int{20, 21}, []int{2, 1}, 35, "", ""},
		{17, 17, "completed", "", []int{32}, []int{3}, 40, "", ""},
		{18, 18, "completed", "谢谢", []int{42, 43}, []int{1, 1}, 45, "", ""},
		{19, 19, "completed", "", []int{15, 16, 17}, []int{1, 1, 1}, 50, "", ""},
		{20, 20, "completed", "很满意", []int{28}, []int{2}, 55, "", ""},
		{21, 21, "completed", "", []int{40}, []int{1}, 60, "", ""},
		{22, 22, "completed", "", []int{35}, []int{2}, 65, "", ""},
		{23, 23, "completed", "还会再来", []int{50, 51}, []int{1, 1}, 70, "", ""},
		{24, 24, "completed", "", []int{1, 2, 3}, []int{1, 1, 1}, 75, "", ""},
		{0, 0, "completed", "", []int{52}, []int{1}, 80, "", ""},
		{1, 1, "completed", "品质一流", []int{4}, []int{2}, 85, "", ""},

		{2, 2, "cancelled", "不想要了", []int{1}, []int{1}, 3, "", ""},
		{3, 3, "cancelled", "", []int{23}, []int{1}, 5, "", ""},
		{4, 4, "cancelled", "价格太高", []int{0, 20}, []int{1, 1}, 8, "", ""},
		{5, 5, "cancelled", "", []int{32}, []int{2}, 12, "", ""},
		{6, 6, "cancelled", "地址填错", []int{42}, []int{1}, 6, "", ""},
	}

	var addrCount int64
	DB.Model(&models.Address{}).Count(&addrCount)

	for _, d := range defs {
		userIdx := d.userIdx + 3
		if userIdx >= len(seedUsers) {
			continue
		}

		addrIdx := d.addrIdx
		if int64(addrIdx) >= addrCount {
			addrIdx = 0
		}

		var addr models.Address
		DB.Offset(addrIdx).First(&addr)

		var total float64
		var items []models.OrderItem
		for j, pIdx := range d.itemIdxes {
			if pIdx >= len(seedProducts) {
				continue
			}
			qty := d.qtys[j]
			price := seedProducts[pIdx].Price
			total += price * float64(qty)
			items = append(items, models.OrderItem{
				ProductID: seedProducts[pIdx].ID,
				Quantity:  qty,
				Price:     price,
			})
		}

		order := models.Order{
			OrderNo:         utils.GenerateOrderNo(),
			UserID:          seedUsers[userIdx].ID,
			AddressID:       addr.ID,
			Total:           total,
			Status:          d.status,
			Remark:          d.remark,
			TrackingCompany: d.tracking,
			TrackingNumber:  d.trackingN,
			Items:           items,
		}

		DB.Create(&order)

		if d.daysAgo > 0 {
			createdAt := time.Now().AddDate(0, 0, -d.daysAgo)
			DB.Model(&models.Order{}).Where("id = ?", order.ID).Update("created_at", createdAt)
		}
	}

	log.Printf("Orders seeded: %d", len(defs))
}

func seedPaymentProofsData() {
	if len(seedUsers) < 25 || len(seedPaymentMethods) < 3 {
		log.Println("Not enough data for payment proof seeding")
		return
	}

	var pendingOrders []models.Order
	DB.Where("status = ? OR status = ?", "pending", "paid").Limit(12).Find(&pendingOrders)

	var processingOrders []models.Order
	DB.Where("status = ?", "processing").Limit(7).Find(&processingOrders)

	var proofIdx int
	for _, order := range pendingOrders {
		if proofIdx >= 12 {
			break
		}
		userIdx := 0
		for j, u := range seedUsers {
			if u.ID == order.UserID {
				userIdx = j
				break
			}
		}

		status := "pending"
		rejectReason := ""
		if proofIdx >= 9 {
			status = "rejected"
			reviewerID := seedUsers[1].ID
			t := time.Now()
			reasons := []string{"图片不清晰，请重新上传", "付款金额与订单不符", "未填写订单备注"}
			rejectReason = reasons[proofIdx-9]
			pm := seedPaymentMethods[proofIdx%3]
			proof := models.PaymentProof{
				OrderID:         order.ID,
				UserID:          seedUsers[userIdx].ID,
				PaymentMethodID: pm.ID,
				ImageUrl:        fmt.Sprintf("https://picsum.photos/seed/proof%02d/400/300", proofIdx),
				Status:          status,
				RejectReason:    rejectReason,
				ReviewerID:      &reviewerID,
				ReviewedAt:      &t,
			}
			DB.Create(&proof)
			proofIdx++
			continue
		}

		pm := seedPaymentMethods[proofIdx%3]
		proof := models.PaymentProof{
			OrderID:         order.ID,
			UserID:          seedUsers[userIdx].ID,
			PaymentMethodID: pm.ID,
			ImageUrl:        fmt.Sprintf("https://picsum.photos/seed/proof%02d/400/300", proofIdx),
			Status:          status,
			RejectReason:    rejectReason,
		}
		DB.Create(&proof)
		proofIdx++
	}

	for i, order := range processingOrders {
		userIdx := 0
		for j, u := range seedUsers {
			if u.ID == order.UserID {
				userIdx = j
				break
			}
		}

		pm := seedPaymentMethods[i%3]
		proof := models.PaymentProof{
			OrderID:         order.ID,
			UserID:          seedUsers[userIdx].ID,
			PaymentMethodID: pm.ID,
			ImageUrl:        fmt.Sprintf("https://picsum.photos/seed/proofapproved%02d/400/300", i),
			Status:          "approved",
			ReviewerID:      &seedUsers[1].ID,
		}
		t := time.Now()
		proof.ReviewedAt = &t
		DB.Create(&proof)
	}

	log.Printf("PaymentProofs seeded")
}

func seedCartAndFavorites() {
	if len(seedUsers) < 25 || len(seedProducts) < 50 {
		log.Println("Not enough data for cart/favorites seeding")
		return
	}

	var cartItems []models.CartItem
	for i := 0; i < 25; i++ {
		userID := seedUsers[i+3].ID
		productIndices := []int{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48, 50, 52}
		count := 1 + (i % 3)
		for j := 0; j < count && j < len(productIndices); j++ {
			pIdx := productIndices[(i+j)%len(productIndices)]
			qty := 1 + (i % 2)
			item := models.CartItem{
				UserID:    userID,
				ProductID: seedProducts[pIdx].ID,
				Quantity:  qty,
			}
			var existing models.CartItem
			if err := DB.Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).First(&existing).Error; err != nil {
				DB.Create(&item)
				cartItems = append(cartItems, item)
			}
		}
	}
	log.Printf("Cart items seeded: %d", len(cartItems))

	var favorites []models.Favorite
	for i := 0; i < 25; i++ {
		userID := seedUsers[i+3].ID
		productIndices := []int{1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49, 51, 53}
		count := 2 + (i % 3)
		for j := 0; j < count && j < len(productIndices); j++ {
			pIdx := productIndices[(i+j)%len(productIndices)]
			fav := models.Favorite{UserID: userID, ProductID: seedProducts[pIdx].ID}
			var existing models.Favorite
			if err := DB.Where("user_id = ? AND product_id = ?", fav.UserID, fav.ProductID).First(&existing).Error; err != nil {
				DB.Create(&fav)
				favorites = append(favorites, fav)
			}
		}
	}
	log.Printf("Favorites seeded: %d", len(favorites))
}

func seedConversationsAndMessages() {
	if len(seedUsers) < 25 {
		log.Println("Not enough users for conversation seeding")
		return
	}

	customerUsers := seedUsers[3:28]

	serviceIDs := []uint{seedUsers[1].ID, seedUsers[2].ID}

	openConvDefs := []struct {
		customerIdx int
		msgPairs    []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}
	}{
		{0, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "你好，我想咨询一下高希霸世纪一号的口感", "text"},
			{"service", 0, "高希霸世纪一号是我们最受欢迎的入门款，口感温和细腻，带有淡淡的奶油和咖啡香气。环规40，长度5英寸，非常适合初次品鉴古巴雪茄的朋友。", "text"},
			{"customer", 0, "那库存现在有货吗？", "text"},
			{"service", 0, "目前库存充足，我们有50支现货，可以当天发货。", "text"},
			{"customer", 0, "好的，我下单了", "text"},
		}},
		{1, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "请问有适合送礼的礼盒吗？", "text"},
			{"service", 0, "当然，我们有几款非常适合送礼的礼盒。推荐您看看我们的「送礼豪华礼盒」，内含多支精选雪茄，配有精美包装，是送礼的首选。另外还有「高希霸品鉴礼盒」，包含世纪1-5号各一支，非常适合送给你尊重的长辈或领导。", "text"},
			{"customer", 0, "送礼豪华礼盒的价格是多少？", "text"},
			{"service", 0, "送礼豪华礼盒定价1288元，目前还有5套现货。", "text"},
			{"customer", 0, "好的，我去看看", "text"},
			{"service", 0, "好的，有任何问题随时联系我。祝您购物愉快！", "text"},
		}},
		{2, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "我的订单怎么还没发货？", "text"},
			{"service", 0, "您好，让我帮您查询一下订单状态。请问您的订单号是多少？", "text"},
			{"customer", 0, "好的稍等我看一下", "text"},
		}},
		{3, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "请问雪茄剪双刃不锈钢和V口哪个好？", "text"},
			{"service", 0, "各有优势。双刃剪适合一次性整齐切割，适合大多数雪茄；V口剪则切割深度较浅，适合环规较粗的雪茄。如果您经常抽粗环规的雪茄，推荐V口剪；如果追求整齐的切口，双刃剪更合适。", "text"},
			{"customer", 0, "明白了，谢谢", "text"},
			{"service", 0, "不客气！欢迎下次光临。", "text"},
		}},
		{4, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "你们支持货到付款吗？", "text"},
			{"service", 0, "支持的，我们支持货到付款、微信支付、支付宝、银行转账以及PayPal等多种支付方式。", "text"},
			{"customer", 0, "好的，那发货用什么快递？", "text"},
			{"service", 0, "我们主要使用顺丰速运和京东快递，大部分城市支持隔日达。偏远地区可能需要3-5天。", "text"},
			{"customer", 0, "明白了，已下单", "text"},
		}},
		{5, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "保湿盒50支装的是什么木材？", "text"},
			{"service", 0, "我们50支装的桌面保湿盒采用的是西班牙雪松木，具有天然的保湿和驱虫效果，是保存雪茄的理想木材。内部配有电子湿度计，可以实时监控温湿度。", "text"},
			{"customer", 0, "价格呢？", "text"},
			{"service", 0, "桌面保湿盒50支装定价388元，含电子湿度计。顺丰包邮。", "text"},
			{"customer", 0, "有保修吗？", "text"},
			{"service", 0, "我们提供一年质保服务，非人为损坏可以免费维修或更换。", "text"},
		}},
		{6, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "想问问帕德龙1964周年口感如何", "text"},
			{"service", 0, "帕德龙1964是尼加拉瓜最经典的雪茄之一，采用该国最好的茄叶，带有浓郁的咖啡、巧克力和香料风味，燃烧均匀，是资深雪茄客的心头好。", "text"},
			{"customer", 0, "有陈年潜力吗？", "text"},
			{"service", 0, "非常有。帕德龙1964建议陈年5年以上，风味会更加圆润醇厚。我们库存18支，生产日期是半年前，正是入手的好时机。", "text"},
		}},
		{7, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "请问有图片吗想看看", "text"},
			{"service", 0, "https://picsum.photos/seed/chatimg01/400/300", "image"},
			{"service", 0, "这是商品实拍图，您可以看到细节。", "text"},
			{"customer", 0, "收到，看起来不错", "text"},
		}},
		{8, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "想退换货怎么办理", "text"},
			{"service", 0, "您好，雪茄属于特殊商品，若因质量问题（如发霉、干燥等）可在收货48小时内申请退换。请提供照片凭证，我们的客服会第一时间处理。", "text"},
			{"customer", 0, "好的明白了", "text"},
		}},
		{9, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "我是老客户了有优惠吗", "text"},
			{"service", 0, "感谢您一直以来的支持！老客户我们有专属折扣码，可以享受9折优惠。请加客服微信获取今日特惠码。", "text"},
			{"customer", 0, "好的谢谢", "text"},
		}},
		{10, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "入门体验装适合新手吗", "text"},
			{"service", 0, "非常适合！入门体验装精选5支不同品牌的经典雪茄，口感从温和到浓郁都有覆盖，还附赠专业品鉴指南，是新手了解雪茄世界的绝佳选择。", "text"},
			{"customer", 0, "价格多少", "text"},
			{"service", 0, "定价198元，新客户首单还有额外9折。", "text"},
			{"customer", 0, "已下单", "text"},
		}},
		{11, []struct {
			senderType string
			senderID   uint
			content    string
			msgType    string
		}{
			{"service", 0, "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？", "text"},
			{"customer", 0, "能开发票吗", "text"},
			{"service", 0, "可以开普通发票或增值税专用发票，请在下单时备注发票抬头和税号。", "text"},
			{"customer", 0, "好的", "text"},
		}},
	}

	for _, def := range openConvDefs {
		customer := customerUsers[def.customerIdx]
		serviceID := serviceIDs[def.customerIdx%2]

		conv := models.Conversation{
			UserID:        customer.ID,
			Status:        "open",
			LastMessageAt: time.Now(),
		}
		DB.Create(&conv)

		for _, msg := range def.msgPairs {
			senderID := customer.ID
			if msg.senderType == "service" {
				senderID = serviceID
			}
			message := models.Message{
				ConversationID: conv.ID,
				SenderType:    msg.senderType,
				SenderID:      senderID,
				MessageType:   msg.msgType,
				Content:       msg.content,
				IsRead:        msg.senderType == "service",
			}
			DB.Create(&message)
		}

		DB.Model(&conv).Update("last_message_at", time.Now())
	}

	closedConvDefs := []struct {
		customerIdx int
		closeBy    string
	}{
		{0, "customer"},
		{1, "service"},
		{2, "customer"},
		{3, "service"},
		{4, "customer"},
		{5, "service"},
		{6, "customer"},
		{7, "service"},
	}

	for _, def := range closedConvDefs {
		customer := customerUsers[def.customerIdx+12]
		serviceID := serviceIDs[def.customerIdx%2]

		conv := models.Conversation{
			UserID:        customer.ID,
			Status:        "closed",
			LastMessageAt: time.Now().AddDate(0, 0, -1),
		}
		DB.Create(&conv)

		msgs := []struct {
			senderType string
			content    string
		}{
			{"service", "您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？"},
			{"customer", "你好"},
			{"service", "您好，请问有什么可以帮您？"},
			{"customer", "我想问一下物流时效"},
			{"service", "顺丰快递一般1-3天送达，偏远地区3-5天。"},
			{"customer", "好的谢谢"},
			{"service", "不客气！祝您生活愉快！"},
		}
		if def.closeBy == "customer" {
			msgs = append(msgs, struct {
				senderType string
				content    string
			}{"system", "客户已结束对话"})
		} else {
			msgs = append(msgs, struct {
				senderType string
				content    string
			}{"system", "客服已结束对话"})
		}

		for _, m := range msgs {
			senderID := customer.ID
			if m.senderType == "service" {
				senderID = serviceID
			} else if m.senderType == "system" {
				senderID = 0
			}
			message := models.Message{
				ConversationID: conv.ID,
				SenderType:     m.senderType,
				SenderID:       senderID,
				MessageType:    "text",
				Content:        m.content,
				IsRead:         true,
			}
			DB.Create(&message)
		}

		DB.Model(&conv).Update("last_message_at", time.Now().AddDate(0, 0, -1))
	}

	log.Printf("Conversations and messages seeded: %d open, %d closed", len(openConvDefs), len(closedConvDefs))
}

func seedNotificationsData() {
	if len(seedUsers) < 25 || len(seedProducts) < 50 {
		log.Println("Not enough data for notifications seeding")
		return
	}

	var notifications []models.Notification

	for i := 0; i < 25; i++ {
		userID := seedUsers[i+3].ID

		notifCount := 4 + (i % 3)

		for j := 0; j < notifCount; j++ {
			notifType := "order_status"
			title := "订单已发货"
			content := "您的订单已发货，快递正在配送中，请注意查收。"
			link := "/orders"
			var productID *uint
			var orderID *uint

			cat := j % 3
			if cat == 0 {
				notifType = "order_status"
				if i%2 == 0 {
					title = "订单已发货"
					content = "您的订单已由顺丰速运发出，请注意查收。"
				} else {
					title = "订单已完成"
					content = "感谢您的购买，期待下次光临！"
				}
				orderPtr := uint(i*3 + j + 100)
				orderID = &orderPtr
			} else if cat == 1 {
				notifType = "back_in_stock"
				title = "商品到货通知"
				productIdx := (i + j) % len(seedProducts)
				content = fmt.Sprintf("您关注的「%s」已到货，欢迎购买！", seedProducts[productIdx].Name)
				productID = &seedProducts[productIdx].ID
				link = fmt.Sprintf("/products/%s", seedProducts[productIdx].Slug)
			} else {
				notifType = "price_drop"
				title = "价格下调通知"
				productIdx := (i + j + 5) % len(seedProducts)
				oldPrice := seedProducts[productIdx].Price * 1.2
				content = fmt.Sprintf("您关注的「%s」价格下调，由¥%.0f调整为¥%.0f，机会不容错过！", seedProducts[productIdx].Name, oldPrice, seedProducts[productIdx].Price)
				productID = &seedProducts[productIdx].ID
				link = fmt.Sprintf("/products/%s", seedProducts[productIdx].Slug)
			}

			isRead := (i+j)%5 != 0

			notif := models.Notification{
				UserID:    userID,
				Type:      notifType,
				Title:     title,
				Content:   content,
				IsRead:    isRead,
				Link:      link,
				ProductID: productID,
				OrderID:   orderID,
			}
			notifications = append(notifications, notif)
		}
	}

	for i := range notifications {
		DB.Create(&notifications[i])
	}

	log.Printf("Notifications seeded: %d", len(notifications))
}

func seedContactMethodsData() {
	methods := []models.ContactMethod{
		{Type: "phone", Label: "客服热线", Value: "400-888-9999", IsActive: true, SortOrder: 1},
		{Type: "email", Label: "邮箱支持", Value: "support@bycigar.com", IsActive: true, SortOrder: 2},
		{Type: "wechat", Label: "微信客服", Value: "BYCIGAR_CS", QRCodeUrl: "https://picsum.photos/seed/wechat-qr/200/200", IsActive: true, SortOrder: 3},
		{Type: "whatsapp", Label: "WhatsApp", Value: "8613800138000", IsActive: true, SortOrder: 4},
		{Type: "telegram", Label: "Telegram", Value: "bycigar_support", IsActive: false, SortOrder: 5},
		{Type: "qq", Label: "QQ客服", Value: "88889999", QRCodeUrl: "https://picsum.photos/seed/qq-qr/200/200", IsActive: false, SortOrder: 6},
	}

	for i := range methods {
		var existing models.ContactMethod
		if err := DB.Where("type = ? AND label = ?", methods[i].Type, methods[i].Label).First(&existing).Error; err != nil {
			DB.Create(&methods[i])
		} else {
			methods[i].ID = existing.ID
		}
	}
	seedContactMethods = methods
	log.Printf("ContactMethods seeded: %d", len(seedContactMethods))
}
