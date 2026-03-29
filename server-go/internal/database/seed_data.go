package database

import (
	"fmt"
	"log"

	"bycigar-server/internal/models"
	"bycigar-server/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

var seedUsers []models.User
var seedCategories []models.Category
var seedProducts []models.Product
var seedAddresses []models.Address

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
	DB.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("Seeding test data...")
	utils.InitSnowflake(1)

	seedUsersData()
	seedCategoriesData()
	seedProductsData()
	seedBannersData()
	seedAddressesData()
	seedCartAndFavorites()
	seedOrdersData()

	log.Println("Test data seeded successfully")
}

func seedUsersData() {
	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	users := []struct {
		email string
		name  string
		role  string
	}{
		{"user1@test.com", "张伟", "customer"},
		{"user2@test.com", "李娜", "customer"},
		{"user3@test.com", "王芳", "customer"},
		{"user4@test.com", "刘洋", "customer"},
		{"user5@test.com", "陈杰", "customer"},
		{"user6@test.com", "杨秀英", "customer"},
		{"user7@test.com", "黄明", "customer"},
		{"user8@test.com", "赵磊", "customer"},
		{"user9@test.com", "周静", "customer"},
		{"user10@test.com", "吴强", "customer"},
		{"user11@test.com", "徐丽", "customer"},
		{"user12@test.com", "孙浩", "customer"},
	}

	for _, u := range users {
		var existing models.User
		if err := DB.Where("email = ?", u.email).First(&existing).Error; err == nil {
			seedUsers = append(seedUsers, existing)
			continue
		}
		user := models.User{
			Email:    u.email,
			Password: string(password),
			Name:     u.name,
			Role:     u.role,
		}
		DB.Create(&user)
		seedUsers = append(seedUsers, user)
	}

	log.Printf("Users seeded: %d", len(seedUsers))
}

func seedCategoriesData() {
	topCategories := []struct {
		name string
		slug string
	}{
		{"古巴雪茄", "cuba"},
		{"多米尼加", "dominican"},
		{"尼加拉瓜", "nicaragua"},
		{"洪都拉斯", "honduras"},
		{"配件工具", "accessories"},
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

	cubaID := createdTop[0].ID
	dominicanID := createdTop[1].ID
	nicaraguaID := createdTop[2].ID
	accessoriesID := createdTop[4].ID

	subCategories := []struct {
		name     string
		slug     string
		parentID uint
	}{
		{"高希霸", "cohiba", cubaID},
		{"蒙特", "montecristo", cubaID},
		{"帕塔加斯", "partagas", cubaID},
		{"罗密欧与朱丽叶", "romeo-y-julieta", cubaID},
		{"大卫杜夫", "davidoff", dominicanID},
		{"富恩特", "arturo-fuente", dominicanID},
		{"帕德龙", "padron", nicaraguaID},
		{"雪茄剪", "cigar-cutter", accessoriesID},
		{"保湿盒", "humidor", accessoriesID},
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
	cohiba := catIDBySlug("cohiba")
	monte := catIDBySlug("montecristo")
	partagas := catIDBySlug("partagas")
	romeo := catIDBySlug("romeo-y-julieta")
	davidoff := catIDBySlug("davidoff")
	fuente := catIDBySlug("arturo-fuente")
	padron := catIDBySlug("padron")
	cutter := catIDBySlug("cigar-cutter")
	humidor := catIDBySlug("humidor")

	products := []models.Product{
		{Name: "高希霸世纪一号", Slug: "cohiba-siglo-i", Price: 128, CategoryID: cohiba, Stock: 50, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba1/400/400", Description: "高希霸世纪系列入门款，口感温和细腻，适合初学者品尝古巴雪茄的经典风味。环规40，长度5英寸。"},
		{Name: "高希霸世纪二号", Slug: "cohiba-siglo-ii", Price: 158, CategoryID: cohiba, Stock: 35, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohiba2/400/400", Description: "中等浓郁的世纪二号，带有奶油和咖啡的香气。"},
		{Name: "高希霸世纪三号", Slug: "cohiba-siglo-iii", Price: 198, CategoryID: cohiba, Stock: 0, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohiba3/400/400", Description: "目前缺货的世纪三号，浓郁度更高。"},
		{Name: "高希霸世纪四号", Slug: "cohiba-siglo-iv", Price: 258, CategoryID: cohiba, Stock: 22, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba4/400/400", Description: "世纪系列中最受欢迎的型号，口感丰富饱满。"},
		{Name: "高希霸世纪五号", Slug: "cohiba-siglo-v", Price: 328, CategoryID: cohiba, Stock: 15, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba5/400/400", Description: "世纪系列旗舰款，复杂多变的风味层次。"},
		{Name: "高希霸世纪六号", Slug: "cohiba-siglo-vi", Price: 468, CategoryID: cohiba, Stock: 8, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohiba6/400/400", Description: "世纪系列顶级款，极富层次感的品吸体验。"},
		{Name: "高希霸短号", Slug: "cohiba-short", Price: 68, CategoryID: cohiba, Stock: 120, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cohibas/400/400", Description: "短小精悍的日常雪茄，适合短暂休憩。"},
		{Name: "高希霸马杜罗5号魔术师", Slug: "cohiba-maduro-5-magicos", Price: 388, CategoryID: cohiba, Stock: 12, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohibam5/400/400", Description: "深色马杜罗茄衣，浓郁甜蜜。"},
		{Name: "高希霸贝伊可52号", Slug: "cohiba-behike-52", Price: 680, CategoryID: cohiba, Stock: 3, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/bk52/400/400", Description: "贝伊可系列，极其稀有，低库存。"},
		{Name: "高希霸长矛", Slug: "cohiba-lancero", Price: 218, CategoryID: cohiba, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/cohibal/400/400", Description: "已下架的经典长矛款。"},

		{Name: "蒙特2号", Slug: "montecristo-no2", Price: 108, CategoryID: monte, Stock: 45, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/monte2/400/400", Description: "蒙特最经典的鱼雷型号，全球畅销。"},
		{Name: "蒙特4号", Slug: "montecristo-no4", Price: 78, CategoryID: monte, Stock: 80, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/monte4/400/400", Description: "世界上最畅销的雪茄之一。"},
		{Name: "蒙特埃德蒙多", Slug: "montecristo-edmundo", Price: 138, CategoryID: monte, Stock: 30, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/monteed/400/400", Description: "丰富的层次感，中等偏浓郁。"},
		{Name: "蒙特双埃德蒙多", Slug: "montecristo-double-edmundo", Price: 168, CategoryID: monte, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/montede/400/400", Description: "加粗版的埃德蒙多，更长的品吸时间。"},
		{Name: "蒙特open初级", Slug: "montecristo-open-junior", Price: 58, CategoryID: monte, Stock: 60, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/montejr/400/400", Description: "入门级蒙特，温和易入口。"},

		{Name: "帕塔加斯D4号", Slug: "partagas-d4", Price: 98, CategoryID: partagas, Stock: 40, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/partagas4/400/400", Description: "帕塔加斯最经典的罗布斯托。"},
		{Name: "帕塔加斯D6号", Slug: "partagas-d6", Price: 128, CategoryID: partagas, Stock: 25, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/partagas6/400/400", Description: "浓郁的泥土和香料风味。"},
		{Name: "帕塔加斯Lusitanias", Slug: "partagas-lusitanias", Price: 358, CategoryID: partagas, Stock: 5, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/partagasl/400/400", Description: "大尺寸双皇冠，低库存限量。"},

		{Name: "罗密欧与朱丽叶2号", Slug: "romeo-no2", Price: 88, CategoryID: romeo, Stock: 55, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeo2/400/400", Description: "经典鱼雷款，浪漫之名。"},
		{Name: "罗密欧与朱丽叶宽丘吉尔", Slug: "romeo-wide-churchill", Price: 178, CategoryID: romeo, Stock: 18, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeowc/400/400", Description: "宽环规丘吉尔，品吸时间充裕。"},
		{Name: "罗密欧与朱丽叶短丘吉尔", Slug: "romeo-short-churchill", Price: 108, CategoryID: romeo, Stock: 35, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeosc/400/400", Description: "短丘吉尔，适合午间休息。"},
		{Name: "罗密欧与朱丽叶俱乐部", Slug: "romeo-club", Price: 48, CategoryID: romeo, Stock: 0, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/romeoclub/400/400", Description: "缺货的小俱乐部款。"},

		{Name: "大卫杜夫2000", Slug: "davidoff-2000", Price: 198, CategoryID: davidoff, Stock: 30, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/dav2000/400/400", Description: "瑞士精工品质，细腻顺滑。"},
		{Name: "大卫杜夫千年系列", Slug: "davidoff-millennium", Price: 298, CategoryID: davidoff, Stock: 12, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/davmil/400/400", Description: "千年系列，浓郁的香料和咖啡。"},
		{Name: "大卫杜夫温斯顿丘吉尔", Slug: "davidoff-winston-churchill", Price: 458, CategoryID: davidoff, Stock: 8, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/davwc/400/400", Description: "致敬伟人的顶级系列。"},
		{Name: "大卫杜夫埃斯库里奥", Slug: "davidoff-escurio", Price: 168, CategoryID: davidoff, Stock: 25, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/davesc/400/400", Description: "巴西巴西利亚茄叶，甜美辛辣。"},

		{Name: "富恩特唐卡洛斯", Slug: "fuente-don-carlos", Price: 228, CategoryID: fuente, Stock: 15, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/fuentedc/400/400", Description: "多米尼加之光，手工精选。"},
		{Name: "富恩特OpusX", Slug: "fuente-opus-x", Price: 688, CategoryID: fuente, Stock: 4, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/fuentox/400/400", Description: "传说中的OpusX，极低库存。"},
		{Name: "富恩特海明威经典", Slug: "fuente-hemingway", Price: 188, CategoryID: fuente, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/fuentehw/400/400", Description: "完美造型，大师级卷制。"},

		{Name: "帕德龙1964周年", Slug: "padron-1964", Price: 278, CategoryID: padron, Stock: 18, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/padron64/400/400", Description: "尼加拉瓜经典，周年纪念系列。"},
		{Name: "帕德龙1926系列80年", Slug: "padron-1926-80", Price: 488, CategoryID: padron, Stock: 6, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/padron26/400/400", Description: "80周年纪念款，极高品质。"},
		{Name: "帕德龙家族 Reserve", Slug: "padron-family-reserve", Price: 568, CategoryID: padron, Stock: 2, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/padronfr/400/400", Description: "家族珍藏，极低库存。"},

		{Name: "雪茄剪-双刃不锈钢", Slug: "cutter-stainless", Price: 128, CategoryID: cutter, Stock: 100, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cutter1/400/400", Description: "高品质不锈钢双刃雪茄剪。"},
		{Name: "雪茄剪-V口切割器", Slug: "cutter-v-cut", Price: 88, CategoryID: cutter, Stock: 80, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cutter2/400/400", Description: "V型切口，完美品吸体验。"},
		{Name: "雪茄钻孔器", Slug: "cutter-punch", Price: 58, CategoryID: cutter, Stock: 150, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/punch1/400/400", Description: "便捷雪茄钻孔器，随身携带。"},
		{Name: "专业雪茄剪套装", Slug: "cutter-pro-set", Price: 298, CategoryID: cutter, Stock: 10, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/cutterpro/400/400", Description: "专业级套装，含剪刀和V口器。"},

		{Name: "桌面保湿盒-50支装", Slug: "humidor-50", Price: 388, CategoryID: humidor, Stock: 20, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/humidor50/400/400", Description: "经典桃花心木桌面保湿盒。"},
		{Name: "旅行保湿盒-5支装", Slug: "humidor-travel-5", Price: 168, CategoryID: humidor, Stock: 35, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/humidort5/400/400", Description: "便携旅行装，密封防干。"},
		{Name: "豪华保湿柜-200支", Slug: "humidor-cabinet-200", Price: 2800, CategoryID: humidor, Stock: 3, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/humidor200/400/400", Description: "顶级豪华展示柜，电子恒温恒湿。"},
		{Name: "电子湿度计", Slug: "hygrometer-digital", Price: 88, CategoryID: humidor, Stock: 50, IsActive: true, IsFeatured: false, Image: "https://picsum.photos/seed/hygrometer/400/400", Description: "精准电子湿度温度计。"},

		{Name: "古巴特别纪念版", Slug: "cuba-special-edition", Price: 888, CategoryID: cohiba, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/cubaspec/400/400", Description: "已下架的特别纪念版。"},
		{Name: "蒙特大师精选", Slug: "montecristo-master", Price: 38, CategoryID: monte, Stock: 0, IsActive: false, IsFeatured: false, Image: "https://picsum.photos/seed/montemas/400/400", Description: "已下架，最低价产品。"},
		{Name: "入门体验装-5支", Slug: "beginner-pack-5", Price: 198, CategoryID: cohiba, Stock: 25, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/beginpack/400/400", Images: "https://picsum.photos/seed/beginpack2/400/400,https://picsum.photos/seed/beginpack3/400/400", Description: "精选5支入门雪茄套装，附品鉴指南。"},
		{Name: "高希霸品鉴礼盒", Slug: "cohiba-tasting-box", Price: 666, CategoryID: cohiba, Stock: 10, IsActive: true, IsFeatured: true, Image: "https://picsum.photos/seed/cohibabox/400/400", Images: "https://picsum.photos/seed/cohibabox2/400/400,https://picsum.photos/seed/cohibabox3/400/400,https://picsum.photos/seed/cohibabox4/400/400", Description: "含世纪1-5号各一支的豪华品鉴礼盒。"},
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
		{Title: "古巴精选 · 传承经典", Image: "https://picsum.photos/seed/banner-cuba/1200/400", Link: "/category/cuba", SortOrder: 1, IsActive: true},
		{Title: "新品上市 · 大卫杜夫温斯顿系列", Image: "https://picsum.photos/seed/banner-davidoff/1200/400", Link: "/category/davidoff", SortOrder: 2, IsActive: true},
		{Title: "配件专区 · 保湿盒限时特惠", Image: "https://picsum.photos/seed/banner-accessory/1200/400", Link: "/category/accessories", SortOrder: 3, IsActive: true},
		{Title: "品鉴礼盒 · 送礼首选", Image: "https://picsum.photos/seed/banner-gift/1200/400", Link: "/product/cohiba-tasting-box", SortOrder: 4, IsActive: true},
		{Title: "已过期活动 · 圣诞特辑", Image: "https://picsum.photos/seed/banner-xmas/1200/400", SortOrder: 5, IsActive: false},
	}

	for i := range banners {
		var existing models.Banner
		if err := DB.Where("title = ?", banners[i].Title).First(&existing).Error; err != nil {
			DB.Create(&banners[i])
		}
	}
	log.Printf("Banners seeded: %d", len(banners))
}

func seedAddressesData() {
	if len(seedUsers) < 12 {
		log.Println("Not enough users for address seeding")
		return
	}

	addresses := []models.Address{
		{UserID: seedUsers[0].ID, FullName: "张伟", AddressLine1: "朝阳区建国路88号SOHO现代城", City: "北京", State: "北京", ZipCode: "100022", Phone: "13800138001", IsDefault: true},
		{UserID: seedUsers[0].ID, FullName: "张伟", AddressLine1: "海淀区中关村大街1号", City: "北京", State: "北京", ZipCode: "100080", Phone: "13800138001", IsDefault: false},
		{UserID: seedUsers[1].ID, FullName: "李娜", AddressLine1: "黄浦区南京东路100号", City: "上海", State: "上海", ZipCode: "200001", Phone: "13900139002", IsDefault: true},
		{UserID: seedUsers[1].ID, FullName: "李娜", AddressLine1: "浦东新区陆家嘴环路500号", City: "上海", State: "上海", ZipCode: "200120", Phone: "13900139002", IsDefault: false},
		{UserID: seedUsers[2].ID, FullName: "王芳", AddressLine1: "天河区天河路385号太古汇", City: "广州", State: "广东", ZipCode: "510620", Phone: "13700137003", IsDefault: true},
		{UserID: seedUsers[3].ID, FullName: "刘洋", AddressLine1: "南山区科技园南区深南大道9966号", City: "深圳", State: "广东", ZipCode: "518057", Phone: "13600136004", IsDefault: true},
		{UserID: seedUsers[4].ID, FullName: "陈杰", AddressLine1: "西湖区文三路553号", City: "杭州", State: "浙江", ZipCode: "310012", Phone: "13500135005", IsDefault: true},
		{UserID: seedUsers[5].ID, FullName: "杨秀英", AddressLine1: "锦江区红星路三段1号国际金融中心", City: "成都", State: "四川", ZipCode: "610021", Phone: "13400134006", IsDefault: true},
		{UserID: seedUsers[6].ID, FullName: "黄明", AddressLine1: "武昌区中北路109号", City: "武汉", State: "湖北", ZipCode: "430071", Phone: "13300133007", IsDefault: true},
		{UserID: seedUsers[7].ID, FullName: "赵磊", AddressLine1: "和平区南京路189号", City: "天津", State: "天津", ZipCode: "300051", Phone: "13200132008", IsDefault: true},
		{UserID: seedUsers[8].ID, FullName: "周静", AddressLine1: "鼓楼区中央路331号", City: "南京", State: "江苏", ZipCode: "210008", Phone: "13100131009", IsDefault: true},
		{UserID: seedUsers[9].ID, FullName: "吴强", AddressLine1: "岳麓区麓谷大道627号", City: "长沙", State: "湖南", ZipCode: "410205", Phone: "13000130010", IsDefault: true},
		{UserID: seedUsers[10].ID, FullName: "徐丽", AddressLine1: "雁塔区长安南路86号", City: "西安", State: "陕西", ZipCode: "710061", Phone: "15000150011", IsDefault: true},
		{UserID: seedUsers[11].ID, FullName: "孙浩", AddressLine1: "思明区厦禾路888号", City: "厦门", State: "福建", ZipCode: "361003", Phone: "15100151012", IsDefault: true},
	}

	for i := range addresses {
		DB.Create(&addresses[i])
	}
	seedAddresses = addresses
	log.Printf("Addresses seeded: %d", len(seedAddresses))
}

func seedCartAndFavorites() {
	if len(seedUsers) < 12 || len(seedProducts) < 40 {
		log.Println("Not enough data for cart/favorites seeding")
		return
	}

	cartItems := []models.CartItem{
		{UserID: seedUsers[0].ID, ProductID: seedProducts[0].ID, Quantity: 2},
		{UserID: seedUsers[0].ID, ProductID: seedProducts[10].ID, Quantity: 1},
		{UserID: seedUsers[1].ID, ProductID: seedProducts[4].ID, Quantity: 1},
		{UserID: seedUsers[1].ID, ProductID: seedProducts[20].ID, Quantity: 3},
		{UserID: seedUsers[2].ID, ProductID: seedProducts[8].ID, Quantity: 1},
		{UserID: seedUsers[3].ID, ProductID: seedProducts[32].ID, Quantity: 2},
		{UserID: seedUsers[3].ID, ProductID: seedProducts[33].ID, Quantity: 1},
		{UserID: seedUsers[4].ID, ProductID: seedProducts[42].ID, Quantity: 1},
		{UserID: seedUsers[5].ID, ProductID: seedProducts[15].ID, Quantity: 4},
		{UserID: seedUsers[6].ID, ProductID: seedProducts[43].ID, Quantity: 1},
	}

	for i := range cartItems {
		var existing models.CartItem
		if err := DB.Where("user_id = ? AND product_id = ?", cartItems[i].UserID, cartItems[i].ProductID).First(&existing).Error; err != nil {
			DB.Create(&cartItems[i])
		}
	}
	log.Printf("Cart items seeded: %d", len(cartItems))

	favorites := []models.Favorite{
		{UserID: seedUsers[0].ID, ProductID: seedProducts[0].ID},
		{UserID: seedUsers[0].ID, ProductID: seedProducts[8].ID},
		{UserID: seedUsers[0].ID, ProductID: seedProducts[42].ID},
		{UserID: seedUsers[1].ID, ProductID: seedProducts[4].ID},
		{UserID: seedUsers[1].ID, ProductID: seedProducts[22].ID},
		{UserID: seedUsers[1].ID, ProductID: seedProducts[25].ID},
		{UserID: seedUsers[2].ID, ProductID: seedProducts[10].ID},
		{UserID: seedUsers[2].ID, ProductID: seedProducts[43].ID},
		{UserID: seedUsers[3].ID, ProductID: seedProducts[5].ID},
		{UserID: seedUsers[3].ID, ProductID: seedProducts[27].ID},
		{UserID: seedUsers[4].ID, ProductID: seedProducts[15].ID},
		{UserID: seedUsers[4].ID, ProductID: seedProducts[20].ID},
		{UserID: seedUsers[5].ID, ProductID: seedProducts[22].ID},
		{UserID: seedUsers[6].ID, ProductID: seedProducts[0].ID},
		{UserID: seedUsers[6].ID, ProductID: seedProducts[32].ID},
		{UserID: seedUsers[7].ID, ProductID: seedProducts[42].ID},
		{UserID: seedUsers[8].ID, ProductID: seedProducts[28].ID},
		{UserID: seedUsers[9].ID, ProductID: seedProducts[12].ID},
		{UserID: seedUsers[9].ID, ProductID: seedProducts[40].ID},
		{UserID: seedUsers[10].ID, ProductID: seedProducts[43].ID},
	}

	for i := range favorites {
		var existing models.Favorite
		if err := DB.Where("user_id = ? AND product_id = ?", favorites[i].UserID, favorites[i].ProductID).First(&existing).Error; err != nil {
			DB.Create(&favorites[i])
		}
	}
	log.Printf("Favorites seeded: %d", len(favorites))
}

func seedOrdersData() {
	if len(seedUsers) < 12 || len(seedProducts) < 40 || len(seedAddresses) < 14 {
		log.Println("Not enough data for order seeding")
		return
	}

	type orderDef struct {
		userIdx    int
		addrIdx    int
		status     string
		remark     string
		itemIdxes  []int
		quantities []int
		daysAgo    int
	}

	defs := []orderDef{
		{0, 0, "pending", "", []int{0, 10}, []int{2, 1}, 0},
		{1, 2, "pending", "请尽快发货", []int{4}, []int{1}, 0},
		{2, 4, "pending", "", []int{20}, []int{3}, 1},
		{3, 5, "pending", "需要发票", []int{32, 33}, []int{1, 1}, 1},
		{4, 6, "pending", "", []int{42}, []int{2}, 2},

		{0, 1, "processing", "", []int{5}, []int{1}, 2},
		{5, 7, "processing", "小心轻放", []int{15}, []int{3}, 3},
		{6, 8, "processing", "", []int{22, 12}, []int{1, 2}, 3},
		{7, 9, "processing", "", []int{28}, []int{1}, 4},

		{1, 3, "shipped", "", []int{10, 14}, []int{1, 2}, 5},
		{2, 4, "shipped", "快递柜自提", []int{8}, []int{1}, 5},
		{3, 5, "shipped", "", []int{0}, []int{5}, 6},
		{8, 10, "shipped", "", []int{35, 36}, []int{1, 2}, 7},
		{9, 11, "shipped", "", []int{43}, []int{1}, 7},

		{0, 0, "completed", "", []int{0}, []int{2}, 10},
		{0, 0, "completed", "", []int{10, 11}, []int{1, 3}, 14},
		{1, 2, "completed", "", []int{4, 22}, []int{1, 1}, 15},
		{1, 2, "completed", "好评", []int{25}, []int{2}, 20},
		{2, 4, "completed", "", []int{20, 21}, []int{2, 1}, 18},
		{3, 5, "completed", "", []int{32}, []int{3}, 22},
		{4, 6, "completed", "", []int{42, 43}, []int{1, 1}, 25},
		{5, 7, "completed", "", []int{15, 16, 17}, []int{1, 1, 1}, 30},
		{6, 8, "completed", "非常满意", []int{28}, []int{2}, 35},
		{7, 9, "completed", "", []int{40}, []int{1}, 40},
		{8, 10, "completed", "", []int{35}, []int{2}, 45},

		{0, 0, "cancelled", "不想要了", []int{1}, []int{1}, 3},
		{1, 2, "cancelled", "", []int{23}, []int{1}, 5},
		{2, 4, "cancelled", "价格不对", []int{0, 20}, []int{1, 1}, 8},
		{3, 5, "cancelled", "", []int{32}, []int{2}, 12},
		{9, 11, "cancelled", "地址填错了", []int{42}, []int{1}, 6},
		{10, 12, "cancelled", "", []int{12}, []int{1}, 15},
	}

	var count int64
	DB.Model(&models.Order{}).Count(&count)
	if count >= int64(len(defs)) {
		log.Println("Orders already seeded")
		return
	}

	for i, d := range defs {
		if d.userIdx >= len(seedUsers) || d.addrIdx >= len(seedAddresses) {
			continue
		}

		var total float64
		var items []models.OrderItem
		for j, pIdx := range d.itemIdxes {
			if pIdx >= len(seedProducts) {
				continue
			}
			qty := d.quantities[j]
			price := seedProducts[pIdx].Price
			total += price * float64(qty)
			items = append(items, models.OrderItem{
				ProductID: seedProducts[pIdx].ID,
				Quantity:  qty,
				Price:     price,
			})
		}

		order := models.Order{
			OrderNo:   fmt.Sprintf("ORD%s%d", utils.GenerateOrderNo(), i),
			UserID:    seedUsers[d.userIdx].ID,
			AddressID: seedAddresses[d.addrIdx].ID,
			Total:     total,
			Status:    d.status,
			Remark:    d.remark,
			Items:     items,
		}

		DB.Create(&order)

		if d.daysAgo > 0 {
			DB.Model(&models.Order{}).Where("id = ?", order.ID).
				Update("created_at", order.CreatedAt.AddDate(0, 0, -d.daysAgo))
		}
	}

	log.Printf("Orders seeded: %d", len(defs))
}
