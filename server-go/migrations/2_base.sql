-- BYCIGAR 基础种子数据
-- 说明: 管理员账号、分类、商品、banner、支付方式等基础数据
-- 密码: 123456 (bcrypt hash)

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 管理员和客服账号
-- 注意: bcrypt("123456") = $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
-- ----------------------------
INSERT INTO `users` (`id`, `email`, `password`, `name`, `role`, `created_at`, `updated_at`) VALUES
(1, 'admin@bycigar.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '系统管理员', 'admin', NOW(), NOW()),
(2, 'service1@bycigar.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '客服小王', 'service', NOW(), NOW()),
(3, 'service2@bycigar.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '客服小李', 'service', NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `role` = VALUES(`role`);

-- ----------------------------
-- 一级分类
-- ----------------------------
INSERT INTO `categories` (`id`, `name`, `slug`, `parent_id`, `created_at`, `updated_at`) VALUES
(1, '精品雪茄', 'premium-cigars', NULL, NOW(), NOW()),
(2, '雪茄配件', 'accessories', NULL, NOW(), NOW()),
(3, '生活方式', 'lifestyle', NULL, NOW(), NOW()),
(4, '礼盒套装', 'gift-sets', NULL, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ----------------------------
-- 二级分类
-- ----------------------------
INSERT INTO `categories` (`id`, `name`, `slug`, `parent_id`, `created_at`, `updated_at`) VALUES
(5, '古巴经典', 'cuba-classic', 1, NOW(), NOW()),
(6, '多米尼加', 'dominican', 1, NOW(), NOW()),
(7, '尼加拉瓜', 'nicaragua', 1, NOW(), NOW()),
(8, '切割工具', 'cutters', 2, NOW(), NOW()),
(9, '保湿存储', 'humidors', 2, NOW(), NOW()),
(10, '点火设备', 'lighters', 2, NOW(), NOW()),
(11, '酒水搭配', 'spirits-pairing', 3, NOW(), NOW()),
(12, '入门礼盒', 'starter-kits', 4, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `parent_id` = VALUES(`parent_id`);

-- ----------------------------
-- 商品 - 古巴经典
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(1, '高希霸世纪一号', 'cohiba-siglo-i', '高希霸世纪系列入门款，口感温和细腻。', 128, 'https://picsum.photos/seed/cohiba1/400/400', 5, 50, 1, 1, NOW(), NOW()),
(2, '高希霸世纪二号', 'cohiba-siglo-ii', '中等浓郁的世纪二号，带有奶油和咖啡的香气。', 158, 'https://picsum.photos/seed/cohiba2/400/400', 5, 35, 1, 0, NOW(), NOW()),
(3, '高希霸世纪三号', 'cohiba-siglo-iii', '浓郁度更高，目前缺货。', 198, 'https://picsum.photos/seed/cohiba3/400/400', 5, 0, 1, 0, NOW(), NOW()),
(4, '高希霸世纪四号', 'cohiba-siglo-iv', '世纪系列中最受欢迎的型号，口感丰富饱满。', 258, 'https://picsum.photos/seed/cohiba4/400/400', 5, 22, 1, 1, NOW(), NOW()),
(5, '高希霸世纪五号', 'cohiba-siglo-v', '世纪系列旗舰款，复杂多变的风味层次。', 328, 'https://picsum.photos/seed/cohiba5/400/400', 5, 15, 1, 1, NOW(), NOW()),
(6, '高希霸短号', 'cohiba-short', '短小精悍的日常雪茄，适合短暂休憩。', 68, 'https://picsum.photos/seed/cohibas/400/400', 5, 120, 1, 0, NOW(), NOW()),
(7, '高希霸马杜罗5号魔术师', 'cohiba-maduro-5', '深色马杜罗茄衣，浓郁甜蜜。', 388, 'https://picsum.photos/seed/cohibam5/400/400', 5, 12, 1, 1, NOW(), NOW()),
(8, '高希霸贝伊可52号', 'cohiba-behike-52', '贝伊可系列，极其稀有，低库存。', 680, 'https://picsum.photos/seed/bk52/400/400', 5, 3, 1, 0, NOW(), NOW()),
(9, '高希霸长矛', 'cohiba-lancero', '已下架的经典长矛款。', 218, 'https://picsum.photos/seed/cohibal/400/400', 5, 0, 0, 0, NOW(), NOW()),
(10, '科伊巴鱼雷限量版', 'cohiba-torpedo-limited', '限量版鱼雷，极高品质。', 888, 'https://picsum.photos/seed/cohibator/400/400', 5, 0, 1, 1, NOW(), NOW()),
(11, '蒙特2号', 'montecristo-no2', '蒙特最经典的鱼雷型号，全球畅销。', 108, 'https://picsum.photos/seed/monte2/400/400', 5, 45, 1, 0, NOW(), NOW()),
(12, '蒙特4号', 'montecristo-no4', '世界上最畅销的雪茄之一。', 78, 'https://picsum.photos/seed/monte4/400/400', 5, 80, 1, 1, NOW(), NOW()),
(13, '蒙特埃德蒙多', 'montecristo-edmundo', '丰富的层次感，中等偏浓郁。', 138, 'https://picsum.photos/seed/monteed/400/400', 5, 30, 1, 0, NOW(), NOW()),
(14, '蒙特双埃德蒙多', 'montecristo-double', '加粗版埃德蒙多，更长品吸时间。', 168, 'https://picsum.photos/seed/montede/400/400', 5, 20, 1, 0, NOW(), NOW()),
(15, '蒙特OPEN初级', 'montecristo-open-junior', '入门级蒙特，温和易入口。', 58, 'https://picsum.photos/seed/montejr/400/400', 5, 60, 1, 0, NOW(), NOW()),
(16, '帕塔加斯D4号', 'partagas-d4', '帕塔加斯最经典的罗布斯托。', 98, 'https://picsum.photos/seed/partagas4/400/400', 5, 40, 1, 0, NOW(), NOW()),
(17, '帕塔加斯D6号', 'partagas-d6', '浓郁的泥土和香料风味。', 128, 'https://picsum.photos/seed/partagas6/400/400', 5, 25, 1, 1, NOW(), NOW()),
(18, '帕塔加斯卢西塔尼亚', 'partagas-lusitanias', '大尺寸双皇冠，低库存限量。', 358, 'https://picsum.photos/seed/partagasl/400/400', 5, 5, 1, 0, NOW(), NOW()),
(19, '帕塔加斯超级皇冠', 'partagas-super-corona', '已下架。', 198, 'https://picsum.photos/seed/partagassc/400/400', 5, 0, 0, 0, NOW(), NOW()),
(20, '罗密欧2号', 'romeo-no2', '经典鱼雷款，浪漫之名。', 88, 'https://picsum.photos/seed/romeo2/400/400', 5, 55, 1, 1, NOW(), NOW()),
(21, '罗密欧宽丘吉尔', 'romeo-wide-churchill', '宽环规丘吉尔，品吸时间充裕。', 178, 'https://picsum.photos/seed/romeowc/400/400', 5, 18, 1, 0, NOW(), NOW()),
(22, '罗密欧短丘吉尔', 'romeo-short-churchill', '短丘吉尔，适合午间休息。', 108, 'https://picsum.photos/seed/romeosc/400/400', 5, 35, 1, 0, NOW(), NOW()),
(23, '罗密欧俱乐部', 'romeo-club', '缺货的小俱乐部款。', 48, 'https://picsum.photos/seed/romeoclub/400/400', 5, 0, 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`), `is_featured` = VALUES(`is_featured`);

-- ----------------------------
-- 商品 - 多米尼加
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(24, '大卫杜夫2000', 'davidoff-2000', '瑞士精工品质，细腻顺滑。', 198, 'https://picsum.photos/seed/dav2000/400/400', 6, 30, 1, 1, NOW(), NOW()),
(25, '大卫杜夫千年系列', 'davidoff-millennium', '千年系列，浓郁的香料和咖啡。', 298, 'https://picsum.photos/seed/davmil/400/400', 6, 12, 1, 0, NOW(), NOW()),
(26, '大卫杜夫温斯顿丘吉尔', 'davidoff-wsc', '致敬伟人的顶级系列。', 458, 'https://picsum.photos/seed/davwc/400/400', 6, 8, 1, 1, NOW(), NOW()),
(27, '大卫杜夫埃斯库里奥', 'davidoff-escurio', '巴西茄叶，甜美辛辣。', 168, 'https://picsum.photos/seed/davesc/400/400', 6, 25, 1, 0, NOW(), NOW()),
(28, '大卫杜夫格兰德', 'davidoff-grande', '已下架。', 388, 'https://picsum.photos/seed/davgr/400/400', 6, 0, 0, 0, NOW(), NOW()),
(29, '富恩特唐卡洛斯', 'fuente-don-carlos', '多米尼加之光，手工精选。', 228, 'https://picsum.photos/seed/fuentedc/400/400', 6, 15, 1, 0, NOW(), NOW()),
(30, '富恩特OpusX', 'fuente-opus-x', '传说中的OpusX，极低库存。', 688, 'https://picsum.photos/seed/fuentox/400/400', 6, 4, 1, 1, NOW(), NOW()),
(31, '富恩特海明威经典', 'fuente-hemingway', '完美造型，大师级卷制。', 188, 'https://picsum.photos/seed/fuentehw/400/400', 6, 20, 1, 0, NOW(), NOW()),
(32, '盛赛迪尔兰多', '盛赛迪亚-maduro', '深色马杜罗风格。', 258, 'https://picsum.photos/seed盛赛迪亚/400/400', 6, 12, 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 尼加拉瓜
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(33, '帕德龙1964周年', 'padron-1964', '尼加拉瓜经典，周年纪念系列。', 278, 'https://picsum.photos/seed/padron64/400/400', 7, 18, 1, 0, NOW(), NOW()),
(34, '帕德龙1926系列80年', 'padron-1926-80', '80周年纪念款，极高品质。', 488, 'https://picsum.photos/seed/padron26/400/400', 7, 6, 1, 1, NOW(), NOW()),
(35, '帕德龙家族Reserve', 'padron-family-reserve', '家族珍藏，极低库存。', 568, 'https://picsum.photos/seed/padronfr/400/400', 7, 2, 1, 0, NOW(), NOW()),
(36, '帕德龙大师系列', 'padron-master', '已下架。', 328, 'https://picsum.photos/seed/padronm/400/400', 7, 0, 0, 0, NOW(), NOW()),
(37, '奥利瓦V系列', 'oliva-v-series', '尼加拉瓜名庄奥利瓦代表作。', 198, 'https://picsum.photos/seed/olivav/400/400', 7, 22, 1, 1, NOW(), NOW()),
(38, 'AJ费尔南德斯新世界', 'aj-fernandez-nw', '新世界风格，浓郁饱满。', 168, 'https://picsum.photos/seed/ajfw/400/400', 7, 15, 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 切割工具
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(39, '雪茄剪双刃不锈钢', 'cutter-stainless', '高品质不锈钢双刃雪茄剪。', 128, 'https://picsum.photos/seed/cutter1/400/400', 8, 100, 1, 0, NOW(), NOW()),
(40, '雪茄剪V口切割器', 'cutter-v-cut', 'V型切口，完美品吸体验。', 88, 'https://picsum.photos/seed/cutter2/400/400', 8, 80, 1, 1, NOW(), NOW()),
(41, '雪茄钻孔器', 'cutter-punch', '便捷雪茄钻孔器，随身携带。', 58, 'https://picsum.photos/seed/punch1/400/400', 8, 150, 1, 0, NOW(), NOW()),
(42, '专业雪茄剪套装', 'cutter-pro-set', '专业级套装，含剪刀和V口器。', 298, 'https://picsum.photos/seed/cutterpro/400/400', 8, 10, 1, 1, NOW(), NOW()),
(43, '电动雪茄剪', 'cutter-electric', '已下架。', 388, 'https://picsum.photos/seed/cuttere/400/400', 8, 0, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 保湿存储
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(44, '桌面保湿盒50支', 'humidor-50', '经典桃花心木桌面保湿盒。', 388, 'https://picsum.photos/seed/humidor50/400/400', 9, 20, 1, 0, NOW(), NOW()),
(45, '旅行保湿盒5支', 'humidor-travel-5', '便携旅行装，密封防干。', 168, 'https://picsum.photos/seed/humidort5/400/400', 9, 35, 1, 1, NOW(), NOW()),
(46, '豪华保湿柜200支', 'humidor-cabinet-200', '顶级豪华展示柜，电子恒温恒湿。', 2800, 'https://picsum.photos/seed/humidor200/400/400', 9, 3, 1, 1, NOW(), NOW()),
(47, '电子湿度计', 'hygrometer-digital', '精准电子湿度温度计。', 88, 'https://picsum.photos/seed/hygrometer/400/400', 9, 50, 1, 0, NOW(), NOW()),
(48, '保湿包套装', 'humidor-pack-set', '缺货中。', 68, 'https://picsum.photos/seed/humidorpack/400/400', 9, 0, 1, 0, NOW(), NOW()),
(49, '雪松木保湿内衬', 'humidor-cedar-liner', '雪松木内衬，增强保湿效果。', 128, 'https://picsum.photos/seed/cedarliner/400/400', 9, 40, 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 点火设备
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(50, '气点火器', 'lighter-butane', '丁烷气点火器，可调节火焰。', 158, 'https://picsum.photos/seed/lightbut/400/400', 10, 60, 1, 1, NOW(), NOW()),
(51, '松木火柴', 'lighter-wood-matches', '天然松木火柴，品茄专用。', 28, 'https://picsum.photos/seed/matchwood/400/400', 10, 200, 1, 0, NOW(), NOW()),
(52, '丁烷火枪', 'lighter-torch', '双焰火枪，点燃方便。', 198, 'https://picsum.photos/seed/torch/400/400', 10, 45, 1, 0, NOW(), NOW()),
(53, 'Xikar点火器', 'lighter-xikar', '已下架。', 288, 'https://picsum.photos/seed/xikar/400/400', 10, 0, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 酒水搭配
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(54, '威士忌杯套装', 'whisky-glass-set', '手工切割水晶威士忌杯套装。', 168, 'https://picsum.photos/seed/whiskyglass/400/400', 11, 30, 1, 1, NOW(), NOW()),
(55, '白兰地杯', 'brandy-glass', '经典白兰地杯，品茄佳配。', 128, 'https://picsum.photos/seed/brandyglass/400/400', 11, 25, 1, 0, NOW(), NOW()),
(56, '朗姆酒精选', 'rum-premium', '陈年朗姆酒，搭配雪茄绝配。', 298, 'https://picsum.photos/seed/rumprem/400/400', 11, 15, 1, 0, NOW(), NOW()),
(57, '雪茄核桃木托盘', 'cigar-walnut-tray', '胡桃木托盘，兼具美观与实用。', 388, 'https://picsum.photos/seed/walnuttray/400/400', 11, 8, 1, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- 商品 - 礼盒套装
-- ----------------------------
INSERT INTO `products` (`id`, `name`, `slug`, `description`, `price`, `image`, `category_id`, `stock`, `is_active`, `is_featured`, `created_at`, `updated_at`) VALUES
(58, '入门体验装5支', 'starter-pack-5', '精选5支入门雪茄套装，附品鉴指南。', 198, 'https://picsum.photos/seed/beginpack/400/400', 12, 25, 1, 1, NOW(), NOW()),
(59, '高希霸品鉴礼盒', 'cohiba-tasting-box', '含世纪1-5号各一支的豪华品鉴礼盒。', 666, 'https://picsum.photos/seed/cohibabox/400/400', 12, 10, 1, 1, NOW(), NOW()),
(60, '旅行雪茄套装', 'travel-cigar-set', '含便携保湿盒和5支精选雪茄。', 358, 'https://picsum.photos/seed/travelset/400/400', 12, 18, 1, 0, NOW(), NOW()),
(61, '送礼豪华礼盒', 'gift-luxury-box', '高端送礼首选，附精美包装。', 1288, 'https://picsum.photos/seed/giftlux/400/400', 12, 5, 1, 1, NOW(), NOW()),
(62, '限量版年度礼盒', 'limited-annual-box', '年度限量版，极具收藏价值。', 2588, 'https://picsum.photos/seed/annualbox/400/400', 12, 2, 1, 0, NOW(), NOW()),
(63, '节日特别版套装', 'festival-special-set', '已下架。', 988, 'https://picsum.photos/seed/festset/400/400', 12, 0, 0, 0, NOW(), NOW()),
(64, '已下架测试商品', 'discontinued-test', '已下架。', 99, 'https://picsum.photos/seed/discont/400/400', 5, 0, 0, 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `price` = VALUES(`price`), `stock` = VALUES(`stock`), `is_active` = VALUES(`is_active`);

-- ----------------------------
-- Banner 横幅
-- ----------------------------
INSERT INTO `banners` (`title`, `image`, `link`, `sort_order`, `is_active`, `created_at`, `updated_at`) VALUES
('古巴经典 · 传承百年', 'https://picsum.photos/seed/banner-cuba/1200/400', '/category/cuba-classic', 1, 1, NOW(), NOW()),
('多米尼加风情 · 细腻优雅', 'https://picsum.photos/seed/banner-dom/1200/400', '/category/dominican', 2, 1, NOW(), NOW()),
('尼加拉瓜激情 · 浓郁澎湃', 'https://picsum.photos/seed/banner-nic/1200/400', '/category/nicaragua', 3, 1, NOW(), NOW()),
('配件专区 · 点亮品茄时刻', 'https://picsum.photos/seed/banner-acc/1200/400', '/category/cutters', 4, 1, NOW(), NOW()),
('生活方式 · 雪茄与美酒', 'https://picsum.photos/seed/banner-life/1200/400', '/category/spirits-pairing', 5, 1, NOW(), NOW()),
('送礼佳选 · 礼盒套装', 'https://picsum.photos/seed/banner-gift/1200/400', '/category/starter-kits', 6, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `image` = VALUES(`image`), `link` = VALUES(`link`), `sort_order` = VALUES(`sort_order`);

-- ----------------------------
-- 支付方式
-- ----------------------------
INSERT INTO `payment_methods` (`name`, `qr_code_url`, `instructions`, `is_active`, `sort_order`, `created_at`, `updated_at`) VALUES
('微信支付', 'https://picsum.photos/seed/wxpay/200/200', '请扫描二维码支付', 1, 1, NOW(), NOW()),
('支付宝', 'https://picsum.photos/seed/alipay/200/200', '请扫描二维码支付', 1, 2, NOW(), NOW()),
('银行转账', '', '转账至指定银行账户', 1, 3, NOW(), NOW()),
('微信支付(旧版)', 'https://picsum.photos/seed/wxpay-old/200/200', '已停用，请使用新版', 0, 4, NOW(), NOW()),
('支付宝(旧版)', 'https://picsum.photos/seed/alipay-old/200/200', '已停用，请使用新版', 0, 5, NOW(), NOW()),
('货到付款', '', '收货时付款', 1, 6, NOW(), NOW()),
('PayPal', '', 'https://paypal.me/yourname', '点击上方按钮前往 PayPal 完成付款，付款后请上传截图', 1, 7, NOW(), NOW()),
('企业转账', '', '仅限企业客户', 0, 8, NOW(), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `qr_code_url` = VALUES(`qr_code_url`), `instructions` = VALUES(`instructions`), `is_active` = VALUES(`is_active`), `sort_order` = VALUES(`sort_order`);

-- ----------------------------
-- 联系人方式
-- ----------------------------
INSERT INTO `contact_methods` (`type`, `label`, `value`, `qr_code_url`, `is_active`, `sort_order`, `created_at`, `updated_at`) VALUES
('phone', '客服热线', '400-888-9999', '', 1, 1, NOW(), NOW()),
('email', '邮箱支持', 'support@bycigar.com', '', 1, 2, NOW(), NOW()),
('wechat', '微信客服', 'BYCIGAR_CS', 'https://picsum.photos/seed/wechat-qr/200/200', 1, 3, NOW(), NOW()),
('whatsapp', 'WhatsApp', '8613800138000', '', 1, 4, NOW(), NOW()),
('telegram', 'Telegram', 'bycigar_support', '', 0, 5, NOW(), NOW()),
('qq', 'QQ客服', '88889999', 'https://picsum.photos/seed/qq-qr/200/200', 0, 6, NOW(), NOW())
ON DUPLICATE KEY UPDATE `label` = VALUES(`label`), `value` = VALUES(`value`), `is_active` = VALUES(`is_active`);

SET FOREIGN_KEY_CHECKS = 1;
