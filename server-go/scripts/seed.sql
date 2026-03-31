-- BYCIGAR Shop SQL Seed Data
-- Run: docker exec -i bycigar-mysql mysql -u root -proot123 bycigar < seed.sql
-- Password for all users: 123456 (bcrypt hash: $2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq)

SET FOREIGN_KEY_CHECKS = 0;

-- Truncate all tables in reverse dependency order
TRUNCATE TABLE order_items;
TRUNCATE TABLE orders;
TRUNCATE TABLE favorites;
TRUNCATE TABLE cart_items;
TRUNCATE TABLE addresses;
TRUNCATE TABLE notifications;
TRUNCATE TABLE conversations;
TRUNCATE TABLE messages;
TRUNCATE TABLE payment_proofs;
TRUNCATE TABLE banners;
TRUNCATE TABLE products;
TRUNCATE TABLE categories;
TRUNCATE TABLE payment_methods;
TRUNCATE TABLE contact_methods;
TRUNCATE TABLE settings;
TRUNCATE TABLE pages;
TRUNCATE TABLE configs;
TRUNCATE TABLE users;

SET FOREIGN_KEY_CHECKS = 1;

-- =============================================
-- CATEGORIES (12 rows: 4 top-level + 8 subcategories)
-- =============================================
INSERT INTO categories (id, name, slug, parent_id, created_at, updated_at) VALUES
(1, '精品雪茄', 'premium-cigars', NULL, NOW(), NOW()),
(2, '雪茄配件', 'accessories', NULL, NOW(), NOW()),
(3, '生活方式', 'lifestyle', NULL, NOW(), NOW()),
(4, '礼盒套装', 'gift-sets', NULL, NOW(), NOW()),
(5, '古巴经典', 'cuba-classic', 1, NOW(), NOW()),
(6, '多米尼加', 'dominican', 1, NOW(), NOW()),
(7, '尼加拉瓜', 'nicaragua', 1, NOW(), NOW()),
(8, '切割工具', 'cutters', 2, NOW(), NOW()),
(9, '保湿存储', 'humidors', 2, NOW(), NOW()),
(10, '点火设备', 'lighters', 2, NOW(), NOW()),
(11, '酒水搭配', 'spirits-pairing', 3, NOW(), NOW()),
(12, '入门礼盒', 'starter-kits', 4, NOW(), NOW());

-- =============================================
-- USERS (28 rows: 3 admin/service + 25 customers)
-- =============================================
INSERT INTO users (id, email, password, name, role, created_at, updated_at) VALUES
(1, 'admin@bycigar.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '系统管理员', 'admin', NOW(), NOW()),
(2, 'service1@bycigar.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '客服小王', 'service', NOW(), NOW()),
(3, 'service2@bycigar.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '客服小李', 'service', NOW(), NOW()),
(4, 'user01@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '张伟', 'customer', NOW() - INTERVAL 90 DAY, NOW()),
(5, 'user02@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '李娜', 'customer', NOW() - INTERVAL 85 DAY, NOW()),
(6, 'user03@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '王芳', 'customer', NOW() - INTERVAL 80 DAY, NOW()),
(7, 'user04@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '刘洋', 'customer', NOW() - INTERVAL 75 DAY, NOW()),
(8, 'user05@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '陈杰', 'customer', NOW() - INTERVAL 70 DAY, NOW()),
(9, 'user06@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '杨秀英', 'customer', NOW() - INTERVAL 65 DAY, NOW()),
(10, 'user07@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '黄明', 'customer', NOW() - INTERVAL 60 DAY, NOW()),
(11, 'user08@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '赵磊', 'customer', NOW() - INTERVAL 55 DAY, NOW()),
(12, 'user09@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '周静', 'customer', NOW() - INTERVAL 50 DAY, NOW()),
(13, 'user10@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '吴强', 'customer', NOW() - INTERVAL 45 DAY, NOW()),
(14, 'user11@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '徐丽', 'customer', NOW() - INTERVAL 40 DAY, NOW()),
(15, 'user12@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '孙浩', 'customer', NOW() - INTERVAL 35 DAY, NOW()),
(16, 'user13@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '马超', 'customer', NOW() - INTERVAL 30 DAY, NOW()),
(17, 'user14@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '胡蝶', 'customer', NOW() - INTERVAL 25 DAY, NOW()),
(18, 'user15@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '朱琳', 'customer', NOW() - INTERVAL 20 DAY, NOW()),
(19, 'user16@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '郭锐', 'customer', NOW() - INTERVAL 18 DAY, NOW()),
(20, 'user17@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '林涛', 'customer', NOW() - INTERVAL 15 DAY, NOW()),
(21, 'user18@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '何雪', 'customer', NOW() - INTERVAL 12 DAY, NOW()),
(22, 'user19@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '高建', 'customer', NOW() - INTERVAL 10 DAY, NOW()),
(23, 'user20@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '罗岚', 'customer', NOW() - INTERVAL 8 DAY, NOW()),
(24, 'user21@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '梁志明', 'customer', NOW() - INTERVAL 6 DAY, NOW()),
(25, 'user22@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '宋雨晴', 'customer', NOW() - INTERVAL 5 DAY, NOW()),
(26, 'user23@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '郑浩', 'customer', NOW() - INTERVAL 4 DAY, NOW()),
(27, 'user24@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '谢辉', 'customer', NOW() - INTERVAL 3 DAY, NOW()),
(28, 'user25@test.com', '$2a$10$rqeq3t0io9etMQzyBfaWoexzFVxTRONUuvya5QhgLjEoyXy57ZWSq', '田亮', 'customer', NOW() - INTERVAL 2 DAY, NOW());

-- =============================================
-- PRODUCTS (54 rows)
-- =============================================
INSERT INTO products (id, name, slug, description, price, image, images, thumbnail_image, category_id, stock, is_active, is_featured, created_at, updated_at) VALUES
-- 古巴经典 (10 products, category_id=5)
(1, '高希霸世纪一号', 'cohiba-siglo-i', '高希霸世纪系列入门款，口感温和细腻。', 128.00, 'https://picsum.photos/seed/cohiba1/400/400', 'https://picsum.photos/seed/cohiba1a/800/600,https://picsum.photos/seed/cohiba1b/800/600', 'https://picsum.photos/seed/cohiba1/200/200', 5, 50, true, true, NOW(), NOW()),
(2, '高希霸世纪二号', 'cohiba-siglo-ii', '中等浓郁的世纪二号，带有奶油和咖啡的香气。', 158.00, 'https://picsum.photos/seed/cohiba2/400/400', 'https://picsum.photos/seed/cohiba2a/800/600', 'https://picsum.photos/seed/cohiba2/200/200', 5, 35, true, false, NOW(), NOW()),
(3, '高希霸世纪三号', 'cohiba-siglo-iii', '浓郁度更高，目前缺货。', 198.00, 'https://picsum.photos/seed/cohiba3/400/400', 'https://picsum.photos/seed/cohiba3a/800/600', 'https://picsum.photos/seed/cohiba3/200/200', 5, 0, true, false, NOW(), NOW()),
(4, '高希霸世纪四号', 'cohiba-siglo-iv', '世纪系列中最受欢迎的型号，口感丰富饱满。', 258.00, 'https://picsum.photos/seed/cohiba4/400/400', 'https://picsum.photos/seed/cohiba4a/800/600', 'https://picsum.photos/seed/cohiba4/200/200', 5, 22, true, true, NOW(), NOW()),
(5, '高希霸世纪五号', 'cohiba-siglo-v', '世纪系列旗舰款，复杂多变的风味层次。', 328.00, 'https://picsum.photos/seed/cohiba5/400/400', 'https://picsum.photos/seed/cohiba5a/800/600', 'https://picsum.photos/seed/cohiba5/200/200', 5, 15, true, true, NOW(), NOW()),
(6, '高希霸短号', 'cohiba-short', '短小精悍的日常雪茄，适合短暂休憩。', 68.00, 'https://picsum.photos/seed/cohibas/400/400', 'https://picsum.photos/seed/cohibasa/800/600', 'https://picsum.photos/seed/cohibas/200/200', 5, 120, true, false, NOW(), NOW()),
(7, '高希霸马杜罗5号魔术师', 'cohiba-maduro-5', '深色马杜罗茄衣，浓郁甜蜜。', 388.00, 'https://picsum.photos/seed/cohibam5/400/400', 'https://picsum.photos/seed/cohibam5a/800/600', 'https://picsum.photos/seed/cohibam5/200/200', 5, 12, true, true, NOW(), NOW()),
(8, '高希霸贝伊可52号', 'cohiba-behike-52', '贝伊可系列，极其稀有，低库存。', 680.00, 'https://picsum.photos/seed/bk52/400/400', 'https://picsum.photos/seed/bk52a/800/600', 'https://picsum.photos/seed/bk52/200/200', 5, 3, true, false, NOW(), NOW()),
(9, '高希霸长矛', 'cohiba-lancero', '已下架的经典长矛款。', 218.00, 'https://picsum.photos/seed/cohibal/400/400', '', '', 5, 0, false, false, NOW(), NOW()),
(10, '科伊巴鱼雷限量版', 'cohiba-torpedo-limited', '限量版鱼雷，极高品质。', 888.00, 'https://picsum.photos/seed/cohibator/400/400', 'https://picsum.photos/seed/cohibatora/800/600', 'https://picsum.photos/seed/cohibator/200/200', 5, 0, true, true, NOW(), NOW()),
(11, '蒙特2号', 'montecristo-no2', '蒙特最经典的鱼雷型号，全球畅销。', 108.00, 'https://picsum.photos/seed/monte2/400/400', 'https://picsum.photos/seed/monte2a/800/600', 'https://picsum.photos/seed/monte2/200/200', 5, 45, true, false, NOW(), NOW()),
(12, '蒙特4号', 'montecristo-no4', '世界上最畅销的雪茄之一。', 78.00, 'https://picsum.photos/seed/monte4/400/400', 'https://picsum.photos/seed/monte4a/800/600', 'https://picsum.photos/seed/monte4/200/200', 5, 80, true, true, NOW(), NOW()),
(13, '蒙特埃德蒙多', 'montecristo-edmundo', '丰富的层次感，中等偏浓郁。', 138.00, 'https://picsum.photos/seed/monteed/400/400', 'https://picsum.photos/seed/monteeda/800/600', 'https://picsum.photos/seed/monteed/200/200', 5, 30, true, false, NOW(), NOW()),
(14, '蒙特双埃德蒙多', 'montecristo-double', '加粗版埃德蒙多，更长品吸时间。', 168.00, 'https://picsum.photos/seed/montede/400/400', 'https://picsum.photos/seed/montedea/800/600', 'https://picsum.photos/seed/montede/200/200', 5, 20, true, false, NOW(), NOW()),
(15, '蒙特OPEN初级', 'montecristo-open-junior', '入门级蒙特，温和易入口。', 58.00, 'https://picsum.photos/seed/montejr/400/400', 'https://picsum.photos/seed/montejra/800/600', 'https://picsum.photos/seed/montejr/200/200', 5, 60, true, false, NOW(), NOW()),
(16, '帕塔加斯D4号', 'partagas-d4', '帕塔加斯最经典的罗布斯托。', 98.00, 'https://picsum.photos/seed/partagas4/400/400', 'https://picsum.photos/seed/partagas4a/800/600', 'https://picsum.photos/seed/partagas4/200/200', 5, 40, true, false, NOW(), NOW()),
(17, '帕塔加斯D6号', 'partagas-d6', '浓郁的泥土和香料风味。', 128.00, 'https://picsum.photos/seed/partagas6/400/400', 'https://picsum.photos/seed/partagas6a/800/600', 'https://picsum.photos/seed/partagas6/200/200', 5, 25, true, true, NOW(), NOW()),
(18, '帕塔加斯卢西塔尼亚', 'partagas-lusitanias', '大尺寸双皇冠，低库存限量。', 358.00, 'https://picsum.photos/seed/partagasl/400/400', 'https://picsum.photos/seed/partagasla/800/600', 'https://picsum.photos/seed/partagasl/200/200', 5, 5, true, false, NOW(), NOW()),
(19, '帕塔加斯超级皇冠', 'partagas-super-corona', '已下架。', 198.00, 'https://picsum.photos/seed/partagassc/400/400', '', '', 5, 0, false, false, NOW(), NOW()),
(20, '罗密欧2号', 'romeo-no2', '经典鱼雷款，浪漫之名。', 88.00, 'https://picsum.photos/seed/romeo2/400/400', 'https://picsum.photos/seed/romeo2a/800/600', 'https://picsum.photos/seed/romeo2/200/200', 5, 55, true, true, NOW(), NOW()),
(21, '罗密欧宽丘吉尔', 'romeo-wide-churchill', '宽环规丘吉尔，品吸时间充裕。', 178.00, 'https://picsum.photos/seed/romeowc/400/400', 'https://picsum.photos/seed/romeowca/800/600', 'https://picsum.photos/seed/romeowc/200/200', 5, 18, true, false, NOW(), NOW()),
(22, '罗密欧短丘吉尔', 'romeo-short-churchill', '短丘吉尔，适合午间休息。', 108.00, 'https://picsum.photos/seed/romeosc/400/400', 'https://picsum.photos/seed/romeosca/800/600', 'https://picsum.photos/seed/romeosc/200/200', 5, 35, true, false, NOW(), NOW()),
(23, '罗密欧俱乐部', 'romeo-club', '缺货的小俱乐部款。', 48.00, 'https://picsum.photos/seed/romeoclub/400/400', '', '', 5, 0, true, false, NOW(), NOW()),

-- 多米尼加 (8 products, category_id=6)
(24, '大卫杜夫2000', 'davidoff-2000', '瑞士精工品质，细腻顺滑。', 198.00, 'https://picsum.photos/seed/dav2000/400/400', 'https://picsum.photos/seed/dav2000a/800/600', 'https://picsum.photos/seed/dav2000/200/200', 6, 30, true, true, NOW(), NOW()),
(25, '大卫杜夫千年系列', 'davidoff-millennium', '千年系列，浓郁的香料和咖啡。', 298.00, 'https://picsum.photos/seed/davmil/400/400', 'https://picsum.photos/seed/davmila/800/600', 'https://picsum.photos/seed/davmil/200/200', 6, 12, true, false, NOW(), NOW()),
(26, '大卫杜夫温斯顿丘吉尔', 'davidoff-wsc', '致敬伟人的顶级系列。', 458.00, 'https://picsum.photos/seed/davwc/400/400', 'https://picsum.photos/seed/davwca/800/600', 'https://picsum.photos/seed/davwc/200/200', 6, 8, true, true, NOW(), NOW()),
(27, '大卫杜夫埃斯库里奥', 'davidoff-escurio', '巴西茄叶，甜美辛辣。', 168.00, 'https://picsum.photos/seed/davesc/400/400', 'https://picsum.photos/seed/davesca/800/600', 'https://picsum.photos/seed/davesc/200/200', 6, 25, true, false, NOW(), NOW()),
(28, '大卫杜夫格兰德', 'davidoff-grande', '已下架。', 388.00, 'https://picsum.photos/seed/davgr/400/400', '', '', 6, 0, false, false, NOW(), NOW()),
(29, '富恩特唐卡洛斯', 'fuente-don-carlos', '多米尼加之光，手工精选。', 228.00, 'https://picsum.photos/seed/fuentedc/400/400', 'https://picsum.photos/seed/fuentedca/800/600', 'https://picsum.photos/seed/fuentedc/200/200', 6, 15, true, false, NOW(), NOW()),
(30, '富恩特OpusX', 'fuente-opus-x', '传说中的OpusX，极低库存。', 688.00, 'https://picsum.photos/seed/fuentox/400/400', 'https://picsum.photos/seed/fuentoxa/800/600', 'https://picsum.photos/seed/fuentox/200/200', 6, 4, true, true, NOW(), NOW()),
(31, '富恩特海明威经典', 'fuente-hemingway', '完美造型，大师级卷制。', 188.00, 'https://picsum.photos/seed/fuentehw/400/400', 'https://picsum.photos/seed/fuentehwa/800/600', 'https://picsum.photos/seed/fuentehw/200/200', 6, 20, true, false, NOW(), NOW()),
(32, '盛赛迪尔兰杜罗', 'shengsaiyadi-maduro', '深色马杜罗风格。', 258.00, 'https://picsum.photos/seed/shengsaiyadi/400/400', 'https://picsum.photos/seed/shengsaiyadia/800/600', 'https://picsum.photos/seed/shengsaiyadi/200/200', 6, 12, true, false, NOW(), NOW()),

-- 尼加拉瓜 (6 products, category_id=7)
(33, '帕德龙1964周年', 'padron-1964', '尼加拉瓜经典，周年纪念系列。', 278.00, 'https://picsum.photos/seed/padron64/400/400', 'https://picsum.photos/seed/padron64a/800/600', 'https://picsum.photos/seed/padron64/200/200', 7, 18, true, false, NOW(), NOW()),
(34, '帕德龙1926系列80年', 'padron-1926-80', '80周年纪念款，极高品质。', 488.00, 'https://picsum.photos/seed/padron26/400/400', 'https://picsum.photos/seed/padron26a/800/600', 'https://picsum.photos/seed/padron26/200/200', 7, 6, true, true, NOW(), NOW()),
(35, '帕德龙家族Reserve', 'padron-family-reserve', '家族珍藏，极低库存。', 568.00, 'https://picsum.photos/seed/padronfr/400/400', 'https://picsum.photos/seed/padronfra/800/600', 'https://picsum.photos/seed/padronfr/200/200', 7, 2, true, false, NOW(), NOW()),
(36, '帕德龙大师系列', 'padron-master', '已下架。', 328.00, 'https://picsum.photos/seed/padronm/400/400', '', '', 7, 0, false, false, NOW(), NOW()),
(37, '奥利瓦V系列', 'oliva-v-series', '尼加拉瓜名庄奥利瓦代表作。', 198.00, 'https://picsum.photos/seed/olivav/400/400', 'https://picsum.photos/seed/olivava/800/600', 'https://picsum.photos/seed/olivav/200/200', 7, 22, true, true, NOW(), NOW()),
(38, 'AJ费尔南德斯新世界', 'aj-fernandez-nw', '新世界风格，浓郁饱满。', 168.00, 'https://picsum.photos/seed/ajfw/400/400', 'https://picsum.photos/seed/ajfwa/800/600', 'https://picsum.photos/seed/ajfw/200/200', 7, 15, true, false, NOW(), NOW()),

-- 切割工具 (5 products, category_id=8)
(39, '雪茄剪双刃不锈钢', 'cutter-stainless', '高品质不锈钢双刃雪茄剪。', 128.00, 'https://picsum.photos/seed/cutter1/400/400', 'https://picsum.photos/seed/cutter1a/800/600', 'https://picsum.photos/seed/cutter1/200/200', 8, 100, true, false, NOW(), NOW()),
(40, '雪茄剪V口切割器', 'cutter-v-cut', 'V型切口，完美品吸体验。', 88.00, 'https://picsum.photos/seed/cutter2/400/400', 'https://picsum.photos/seed/cutter2a/800/600', 'https://picsum.photos/seed/cutter2/200/200', 8, 80, true, true, NOW(), NOW()),
(41, '雪茄钻孔器', 'cutter-punch', '便捷雪茄钻孔器，随身携带。', 58.00, 'https://picsum.photos/seed/punch1/400/400', 'https://picsum.photos/seed/punch1a/800/600', 'https://picsum.photos/seed/punch1/200/200', 8, 150, true, false, NOW(), NOW()),
(42, '专业雪茄剪套装', 'cutter-pro-set', '专业级套装，含剪刀和V口器。', 298.00, 'https://picsum.photos/seed/cutterpro/400/400', 'https://picsum.photos/seed/cutterproa/800/600', 'https://picsum.photos/seed/cutterpro/200/200', 8, 10, true, true, NOW(), NOW()),
(43, '电动雪茄剪', 'cutter-electric', '已下架。', 388.00, 'https://picsum.photos/seed/cuttere/400/400', '', '', 8, 0, false, false, NOW(), NOW()),

-- 保湿存储 (6 products, category_id=9)
(44, '桌面保湿盒50支', 'humidor-50', '经典桃花心木桌面保湿盒。', 388.00, 'https://picsum.photos/seed/humidor50/400/400', 'https://picsum.photos/seed/humidor50a/800/600', 'https://picsum.photos/seed/humidor50/200/200', 9, 20, true, false, NOW(), NOW()),
(45, '旅行保湿盒5支', 'humidor-travel-5', '便携旅行装，密封防干。', 168.00, 'https://picsum.photos/seed/humidort5/400/400', 'https://picsum.photos/seed/humidort5a/800/600', 'https://picsum.photos/seed/humidort5/200/200', 9, 35, true, true, NOW(), NOW()),
(46, '豪华保湿柜200支', 'humidor-cabinet-200', '顶级豪华展示柜，电子恒温恒湿。', 2800.00, 'https://picsum.photos/seed/humidor200/400/400', 'https://picsum.photos/seed/humidor200a/800/600', 'https://picsum.photos/seed/humidor200/200/200', 9, 3, true, true, NOW(), NOW()),
(47, '电子湿度计', 'hygrometer-digital', '精准电子湿度温度计。', 88.00, 'https://picsum.photos/seed/hygrometer/400/400', 'https://picsum.photos/seed/hygrometera/800/600', 'https://picsum.photos/seed/hygrometer/200/200', 9, 50, true, false, NOW(), NOW()),
(48, '保湿包套装', 'humidor-pack-set', '缺货中。', 68.00, 'https://picsum.photos/seed/humidorpack/400/400', '', '', 9, 0, true, false, NOW(), NOW()),
(49, '雪松木保湿内衬', 'humidor-cedar-liner', '雪松木内衬，增强保湿效果。', 128.00, 'https://picsum.photos/seed/cedarliner/400/400', 'https://picsum.photos/seed/cedarlinera/800/600', 'https://picsum.photos/seed/cedarliner/200/200', 9, 40, true, false, NOW(), NOW()),

-- 点火设备 (4 products, category_id=10)
(50, '丁烷气点火器', 'lighter-butane', '丁烷气点火器，可调节火焰。', 158.00, 'https://picsum.photos/seed/lightbut/400/400', 'https://picsum.photos/seed/lightbuta/800/600', 'https://picsum.photos/seed/lightbut/200/200', 10, 60, true, true, NOW(), NOW()),
(51, '松木火柴', 'lighter-wood-matches', '天然松木火柴，品茄专用。', 28.00, 'https://picsum.photos/seed/matchwood/400/400', 'https://picsum.photos/seed/matchwooda/800/600', 'https://picsum.photos/seed/matchwood/200/200', 10, 200, true, false, NOW(), NOW()),
(52, '丁烷火枪', 'lighter-torch', '双焰火枪，点燃方便。', 198.00, 'https://picsum.photos/seed/torch/400/400', 'https://picsum.photos/seed/torcha/800/600', 'https://picsum.photos/seed/torch/200/200', 10, 45, true, false, NOW(), NOW()),
(53, 'Xikar点火器', 'lighter-xikar', '已下架。', 288.00, 'https://picsum.photos/seed/xikar/400/400', '', '', 10, 0, false, false, NOW(), NOW()),

-- 酒水搭配 (4 products, category_id=11)
(54, '威士忌杯套装', 'whisky-glass-set', '手工切割水晶威士忌杯套装。', 168.00, 'https://picsum.photos/seed/whiskyglass/400/400', 'https://picsum.photos/seed/whiskyglassa/800/600', 'https://picsum.photos/seed/whiskyglass/200/200', 11, 30, true, true, NOW(), NOW()),
(55, '白兰地杯', 'brandy-glass', '经典白兰地杯，品茄佳配。', 128.00, 'https://picsum.photos/seed/brandyglass/400/400', 'https://picsum.photos/seed/brandyglassa/800/600', 'https://picsum.photos/seed/brandyglass/200/200', 11, 25, true, false, NOW(), NOW()),
(56, '朗姆酒精选', 'rum-premium', '陈年朗姆酒，搭配雪茄绝配。', 298.00, 'https://picsum.photos/seed/rumprem/400/400', 'https://picsum.photos/seed/rumprema/800/600', 'https://picsum.photos/seed/rumprem/200/200', 11, 15, true, false, NOW(), NOW()),
(57, '雪茄核桃木托盘', 'cigar-walnut-tray', '胡桃木托盘，兼具美观与实用。', 388.00, 'https://picsum.photos/seed/walnuttray/400/400', 'https://picsum.photos/seed/walnuttraya/800/600', 'https://picsum.photos/seed/walnuttray/200/200', 11, 8, true, false, NOW(), NOW()),

-- 入门礼盒 (7 products, category_id=12)
(58, '入门体验装5支', 'starter-pack-5', '精选5支入门雪茄套装，附品鉴指南。', 198.00, 'https://picsum.photos/seed/beginpack/400/400', 'https://picsum.photos/seed/beginpacka/800/600', 'https://picsum.photos/seed/beginpack/200/200', 12, 25, true, true, NOW(), NOW()),
(59, '高希霸品鉴礼盒', 'cohiba-tasting-box', '含世纪1-5号各一支的豪华品鉴礼盒。', 666.00, 'https://picsum.photos/seed/cohibabox/400/400', 'https://picsum.photos/seed/cohibaboxa/800/600', 'https://picsum.photos/seed/cohibabox/200/200', 12, 10, true, true, NOW(), NOW()),
(60, '旅行雪茄套装', 'travel-cigar-set', '含便携保湿盒和5支精选雪茄。', 358.00, 'https://picsum.photos/seed/travelset/400/400', 'https://picsum.photos/seed/travelseta/800/600', 'https://picsum.photos/seed/travelset/200/200', 12, 18, true, false, NOW(), NOW()),
(61, '送礼豪华礼盒', 'gift-luxury-box', '高端送礼首选，附精美包装。', 1288.00, 'https://picsum.photos/seed/giftlux/400/400', 'https://picsum.photos/seed/giftluxa/800/600', 'https://picsum.photos/seed/giftlux/200/200', 12, 5, true, true, NOW(), NOW()),
(62, '限量版年度礼盒', 'limited-annual-box', '年度限量版，极具收藏价值。', 2588.00, 'https://picsum.photos/seed/annualbox/400/400', 'https://picsum.photos/seed/annualboxa/800/600', 'https://picsum.photos/seed/annualbox/200/200', 12, 2, true, false, NOW(), NOW()),
(63, '节日特别版套装', 'festival-special-set', '已下架。', 988.00, 'https://picsum.photos/seed/festset/400/400', '', '', 12, 0, false, false, NOW(), NOW()),
(64, '已下架测试商品', 'discontinued-test', '已下架。', 99.00, 'https://picsum.photos/seed/discont/400/400', '', '', 12, 0, false, false, NOW(), NOW());

-- =============================================
-- BANNERS (6 rows)
-- =============================================
INSERT INTO banners (id, title, image, link, sort_order, is_active, created_at, updated_at) VALUES
(1, '古巴经典 · 传承百年', 'https://picsum.photos/seed/banner-cuba/1200/400', '/category/cuba-classic', 1, true, NOW(), NOW()),
(2, '多米尼加风情 · 细腻优雅', 'https://picsum.photos/seed/banner-dom/1200/400', '/category/dominican', 2, true, NOW(), NOW()),
(3, '尼加拉瓜激情 · 浓郁澎湃', 'https://picsum.photos/seed/banner-nic/1200/400', '/category/nicaragua', 3, true, NOW(), NOW()),
(4, '配件专区 · 点亮品茄时刻', 'https://picsum.photos/seed/banner-acc/1200/400', '/category/cutters', 4, true, NOW(), NOW()),
(5, '生活方式 · 雪茄与美酒', 'https://picsum.photos/seed/banner-life/1200/400', '/category/spirits-pairing', 5, true, NOW(), NOW()),
(6, '送礼佳选 · 礼盒套装', 'https://picsum.photos/seed/banner-gift/1200/400', '/category/starter-kits', 6, true, NOW(), NOW());

-- =============================================
-- PAYMENT_METHODS (8 rows)
-- =============================================
INSERT INTO payment_methods (id, name, qr_code_url, payment_url, instructions, is_active, sort_order, created_at, updated_at) VALUES
(1, '微信支付', 'https://picsum.photos/seed/wxpay/200/200', '', '请扫描二维码支付', true, 1, NOW(), NOW()),
(2, '支付宝', 'https://picsum.photos/seed/alipay/200/200', '', '请扫描二维码支付', true, 2, NOW(), NOW()),
(3, '银行转账', '', '', '转账至指定银行账户', true, 3, NOW(), NOW()),
(4, '微信支付(旧版)', 'https://picsum.photos/seed/wxpay-old/200/200', '', '已停用，请使用新版', false, 4, NOW(), NOW()),
(5, '支付宝(旧版)', 'https://picsum.photos/seed/alipay-old/200/200', '', '已停用，请使用新版', false, 5, NOW(), NOW()),
(6, '货到付款', '', '', '收货时付款', true, 6, NOW(), NOW()),
(7, 'PayPal', '', 'https://paypal.me/yourname', '点击上方按钮前往 PayPal 完成付款，付款后请上传截图', true, 7, NOW(), NOW()),
(8, '企业转账', '', '', '仅限企业客户', false, 8, NOW(), NOW());

-- =============================================
-- SETTINGS (8 rows)
-- =============================================
INSERT INTO settings (id, `key`, value, created_at, updated_at) VALUES
(1, 'site_name', 'BYCIGAR 雪茄旗舰店', NOW(), NOW()),
(2, 'chat_greeting', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', NOW(), NOW()),
(3, 'contact_phone', '400-888-9999', NOW(), NOW()),
(4, 'contact_email', 'support@bycigar.com', NOW(), NOW()),
(5, 'footer_content', '© 2024 BYCIGAR 版权所有', NOW(), NOW()),
(6, 'order_auto_close_days', '7', NOW(), NOW()),
(7, 'low_stock_threshold', '10', NOW(), NOW()),
(8, 'maintenance_mode', 'false', NOW(), NOW());

-- =============================================
-- PAGES (5 rows)
-- =============================================
INSERT INTO pages (id, slug, title, content, updated_at) VALUES
(1, 'about', '关于我们', 'BYCIGAR 成立于2010年，是国内领先的高端雪茄零售商。我们致力于为雪茄爱好者提供最优质的古巴及世界顶级雪茄，产品涵盖高希霸、蒙特、帕塔加斯、大卫杜夫等经典品牌。公司总部位于上海，拥有专业化的恒温仓储物流体系，所有雪茄均经过严格品控，确保品质如一。\n\n我们的团队由一群热爱雪茄文化的专业人士组成，每位顾问都经过严格培训，能够为客户提供专业的品鉴建议和搭配推荐。BYCIGAR 还定期举办雪茄品鉴会，邀请国内外知名雪茄大师与会员互动交流。\n\n未来，我们将继续深耕雪茄文化，引入更多优质产品，为中国雪茄市场的发展贡献力量。', NOW()),
(2, 'services', '服务条款', '一、服务说明\n\nBYCIGAR 为您提供在线雪茄及配件销售服务，包括商品浏览、在线购买、支付结算、物流配送等完整流程。\n\n二、购买须知\n\n1. 雪茄产品需年满18周岁方可购买。\n2. 请如实填写收货信息，确保商品顺利送达。\n3. 雪茄属于特殊商品一经拆封不支持无理由退换。\n4. 如收到商品有质量问题，请在48小时内联系客服。\n\n三、物流配送\n\n我们合作的物流伙伴包括顺丰、京东等优质快递，全程冷链运输，确保雪茄在最佳状态下送达。部分地区可能需要身份证验证。\n\n四、售后服务\n\n客服热线：400-888-9999（工作日9:00-22:00）\n邮箱：support@bycigar.com', NOW()),
(3, 'privacy-policy', '隐私政策', 'BYCIGAR 非常重视您的个人信息的保护。本隐私政策说明了我们在您使用我们的服务时如何收集、使用、存储和保护您的个人信息。\n\n一、信息收集\n\n当您注册账户时，我们会收集您的姓名、邮箱、手机号码等基本信息。当您下单时，我们会收集收货地址信息以完成配送。\n\n二、信息使用\n\n您的个人信息将用于：处理订单、提供客户服务、发送订单状态通知、推荐适合的商品、改进我们的服务。\n\n三、信息保护\n\n我们采用行业标准的加密技术保护您的数据安全，未经您的授权，我们不会将个人信息提供给任何第三方。\n\n四、联系我们\n\n如对隐私政策有任何疑问，请联系：privacy@bycigar.com', NOW()),
(4, 'statement', '免责说明', '一、产品声明\n\nBYCIGAR 销售的所有雪茄产品均为正品。我们合作的品牌均经过官方授权，或从正规渠道采购。\n\n二、健康提示\n\n吸烟有害健康。雪茄含有焦油和一氧化碳等有害物质，未成年人、孕妇、哺乳期妇女及心血管疾病患者不应使用雪茄。\n\n三、使用风险\n\n雪茄的保存需要特定温度（18-21°C）和湿度（68-72%）。因客户保存不当导致的雪茄品质问题，不在退换货范围内。\n\n四、配送风险\n\n商品在运输过程中可能因不可抗力因素导致延误或损坏，我们将配合物流公司积极处理，但不承担额外赔偿责任。', NOW()),
(5, 'shipping', '配送说明', '一、配送范围\n\n我们支持全国大部分地区的配送，包括港澳台地区。部分偏远地区可能需要更长的配送时间。\n\n二、配送方式\n\n1. 标准配送：3-5个工作日送达\n2. 加急配送：1-2个工作日送达（需额外支付费用）\n3. 定时配送：您可选择具体的收货时间段\n\n三、配送费用\n\n订单满500元免运费，不满500元收取20元运费。新会员首单免运费。\n\n四、配送包装\n\n所有雪茄均采用专业保湿包装，配合保温箱确保品质。礼盒订单额外使用精美纸箱包装。\n\n五、签收提示\n\n请在签收前检查包装是否完好，如有问题请当场拒收并联系客服。', NOW());

-- =============================================
-- CONFIGS (6 rows)
-- =============================================
INSERT INTO configs (id, config_key, config_value, created_at, updated_at) VALUES
(1, 'site_title', 'BYCIGAR 雪茄商城', NOW(), NOW()),
(2, 'seo_description', '专业古巴雪茄及高档雪茄配件销售平台，提供高希霸、蒙特、大卫杜夫等品牌雪茄，品类齐全，品质保证。', NOW(), NOW()),
(3, 'customer_service_hours', '9:00-22:00', NOW(), NOW()),
(4, 'max_cart_items', '99', NOW(), NOW()),
(5, 'order_max_quantity', '50', NOW(), NOW()),
(6, 'exchange_rate_usd', '7.25', NOW(), NOW());

-- =============================================
-- CONTACT_METHODS (6 rows)
-- =============================================
INSERT INTO contact_methods (id, type, label, value, qr_code_url, is_active, sort_order, created_at, updated_at) VALUES
(1, 'phone', '客服热线', '400-888-9999', '', true, 1, NOW(), NOW()),
(2, 'email', '邮箱支持', 'support@bycigar.com', '', true, 2, NOW(), NOW()),
(3, 'wechat', '微信客服', 'BYCIGAR_CS', 'https://picsum.photos/seed/wechat-qr/200/200', true, 3, NOW(), NOW()),
(4, 'whatsapp', 'WhatsApp', '8613800138000', '', true, 4, NOW(), NOW()),
(5, 'telegram', 'Telegram', 'bycigar_support', '', false, 5, NOW(), NOW()),
(6, 'qq', 'QQ客服', '88889999', 'https://picsum.photos/seed/qq-qr/200/200', false, 6, NOW(), NOW());

-- =============================================
-- ADDRESSES (25 rows, one per customer)
-- =============================================
INSERT INTO addresses (id, user_id, full_name, address_line1, address_line2, city, state, zip_code, phone, is_default, created_at, updated_at) VALUES
(1, 4, '张伟', '建国路88号SOHO现代城', '', '北京', '北京', '100022', '13800138001', true, NOW(), NOW()),
(2, 5, '李娜', '南京东路100号', '', '上海', '上海', '200001', '13900139002', true, NOW(), NOW()),
(3, 6, '王芳', '天河路385号太古汇', '', '广州', '广东', '510620', '13700137003', true, NOW(), NOW()),
(4, 7, '刘洋', '科技园南区深南大道9966号', '', '深圳', '广东', '518057', '13600136004', true, NOW(), NOW()),
(5, 8, '陈杰', '文三路553号', '', '杭州', '浙江', '310012', '13500135005', true, NOW(), NOW()),
(6, 9, '杨秀英', '红星路三段1号', '', '成都', '四川', '610021', '13400134006', true, NOW(), NOW()),
(7, 10, '黄明', '中北路109号', '', '武汉', '湖北', '430071', '13300133007', true, NOW(), NOW()),
(8, 11, '赵磊', '南京路189号', '', '天津', '天津', '300051', '13200132008', true, NOW(), NOW()),
(9, 12, '周静', '中央路331号', '', '南京', '江苏', '210008', '13100131009', true, NOW(), NOW()),
(10, 13, '吴强', '麓谷大道627号', '', '长沙', '湖南', '410205', '13000130010', true, NOW(), NOW()),
(11, 14, '徐丽', '长安南路86号', '', '西安', '陕西', '710061', '15000150011', true, NOW(), NOW()),
(12, 15, '孙浩', '厦禾路888号', '', '厦门', '福建', '361003', '15100151012', true, NOW(), NOW()),
(13, 16, '马超', '解放碑步行街88号', '', '重庆', '重庆', '400010', '15200152013', true, NOW(), NOW()),
(14, 17, '胡蝶', '观前街168号', '', '苏州', '江苏', '215000', '15300153014', true, NOW(), NOW()),
(15, 18, '朱琳', '二七路72号', '', '郑州', '河南', '450000', '15500155015', true, NOW(), NOW()),
(16, 19, '郭锐', '香港中路100号', '', '青岛', '山东', '266000', '15600156016', true, NOW(), NOW()),
(17, 20, '林涛', '五四路200号', '', '沈阳', '辽宁', '110001', '15700157017', true, NOW(), NOW()),
(18, 21, '何雪', '长江道1号', '', '大连', '辽宁', '116001', '15800158018', true, NOW(), NOW()),
(19, 22, '高建', '东风路388号', '', '济南', '山东', '250001', '15900159019', true, NOW(), NOW()),
(20, 23, '罗岚', '五一路399号', '', '福州', '福建', '350001', '18000180020', true, NOW(), NOW()),
(21, 24, '梁志明', '人民路555号', '', '昆明', '云南', '650000', '18100181021', true, NOW(), NOW()),
(22, 25, '宋雨晴', '建设路188号', '', '哈尔滨', '黑龙江', '150001', '18200182022', true, NOW(), NOW()),
(23, 26, '郑浩', '中山路299号', '', '长春', '吉林', '130000', '18300183023', true, NOW(), NOW()),
(24, 27, '谢辉', '胜利北路66号', '', '石家庄', '河北', '050000', '18400184024', true, NOW(), NOW()),
(25, 28, '田亮', '滨海大道88号中心大厦', '', '南昌', '江西', '330000', '18500185025', true, NOW(), NOW());

-- Additional addresses for some users (users 4, 7, 10, 13, 16, 19, 22, 25 get a second address)
INSERT INTO addresses (id, user_id, full_name, address_line1, address_line2, city, state, zip_code, phone, is_default, created_at, updated_at) VALUES
(26, 4, '张伟', '建国路88号SOHO现代城B座', '', '北京', '北京', '100022', '13800138001', false, NOW(), NOW()),
(27, 7, '刘洋', '科技园南区深南大道9966号B座', '', '深圳', '广东', '518057', '13600136004', false, NOW(), NOW()),
(28, 10, '黄明', '中北路109号B座', '', '武汉', '湖北', '430071', '13300133007', false, NOW(), NOW()),
(29, 13, '吴强', '麓谷大道627号B座', '', '长沙', '湖南', '410205', '13000130010', false, NOW(), NOW()),
(30, 16, '马超', '解放碑步行街88号B座', '', '重庆', '重庆', '400010', '15200152013', false, NOW(), NOW()),
(31, 19, '郭锐', '香港中路100号（仓库）', '', '青岛', '山东', '266000', '15600156016', false, NOW(), NOW()),
(32, 22, '宋雨晴', '人民路555号（仓库）', '', '昆明', '云南', '650000', '18200182022', false, NOW(), NOW()),
(33, 25, '田亮', '滨海大道88号中心大厦B座', '', '南昌', '江西', '330000', '18500185025', false, NOW(), NOW());

-- =============================================
-- ORDERS (40 rows)
-- =============================================
INSERT INTO orders (id, order_no, user_id, address_id, total, status, remark, tracking_company, tracking_number, created_at, updated_at) VALUES
-- Pending orders (8)
(1, 'BYC2026033100001001', 4, 1, 514.00, 'pending', '', '', '', NOW() - INTERVAL 0 DAY, NOW()),
(2, 'BYC2026033100001002', 5, 2, 258.00, 'pending', '请尽快发货', '', '', NOW() - INTERVAL 0 DAY, NOW()),
(3, 'BYC2026033000001003', 6, 3, 168.00, 'pending', '', '', '', NOW() - INTERVAL 1 DAY, NOW()),
(4, 'BYC2026033000001004', 7, 4, 216.00, 'pending', '需要发票', '', '', NOW() - INTERVAL 1 DAY, NOW()),
(5, 'BYC2026032900001005', 8, 5, 28.00, 'pending', '', '', '', NOW() - INTERVAL 2 DAY, NOW()),
(6, 'BYC2026033100001006', 9, 6, 198.00, 'pending', '小心轻放', '', '', NOW() - INTERVAL 0 DAY, NOW()),
(7, 'BYC2026033000001007', 10, 7, 456.00, 'pending', '', '', '', NOW() - INTERVAL 1 DAY, NOW()),
(8, 'BYC2026033100001008', 11, 8, 688.00, 'pending', '生日礼物', '', '', NOW() - INTERVAL 0 DAY, NOW()),

-- Paid orders (5)
(9, 'BYC2026033100001009', 12, 9, 298.00, 'paid', '已转账', '', '', NOW() - INTERVAL 0 DAY, NOW()),
(10, 'BYC2026033000001010', 13, 10, 644.00, 'paid', '', '', '', NOW() - INTERVAL 1 DAY, NOW()),
(11, 'BYC2026033100001011', 14, 11, 426.00, 'paid', '', '', '', NOW() - INTERVAL 0 DAY, NOW()),
(12, 'BYC2026033000001012', 15, 12, 258.00, 'paid', '谢谢', '', '', NOW() - INTERVAL 1 DAY, NOW()),
(13, 'BYC2026033100001013', 16, 13, 666.00, 'paid', '', '', '', NOW() - INTERVAL 0 DAY, NOW()),

-- Processing orders (10)
(14, 'BYC2026032900001014', 17, 14, 356.00, 'processing', '', '顺丰速运', 'SF1234567890', NOW() - INTERVAL 2 DAY, NOW()),
(15, 'BYC2026032800001015', 18, 15, 116.00, 'processing', '请仔细包装', '京东快递', 'JD0001234567', NOW() - INTERVAL 3 DAY, NOW()),
(16, 'BYC2026032800001016', 19, 16, 266.00, 'processing', '', '顺丰速运', 'SF9876543210', NOW() - INTERVAL 4 DAY, NOW()),
(17, 'BYC2026032700001017', 20, 17, 666.00, 'processing', '节假日不送', 'EMS', 'EM000123456', NOW() - INTERVAL 4 DAY, NOW()),
(18, 'BYC2026032700001018', 21, 18, 226.00, 'processing', '', '顺丰速运', 'SF5555666677', NOW() - INTERVAL 5 DAY, NOW()),
(19, 'BYC2026032700001019', 22, 19, 336.00, 'processing', '办公室地址', '京东快递', 'JD888999000', NOW() - INTERVAL 5 DAY, NOW()),
(20, 'BYC2026032600001020', 23, 20, 486.00, 'processing', '', '顺丰速运', 'SF1111222233', NOW() - INTERVAL 6 DAY, NOW()),
(21, 'BYC2026032600001021', 24, 21, 328.00, 'processing', '尽快', '德邦物流', 'DB444455556', NOW() - INTERVAL 6 DAY, NOW()),
(22, 'BYC2026032500001022', 25, 22, 564.00, 'processing', '', '顺丰速运', 'SF7777888899', NOW() - INTERVAL 7 DAY, NOW()),

-- Shipped orders (12)
(23, 'BYC2026032400001023', 4, 1, 328.00, 'shipped', '', '顺丰速运', 'SF2222333344', NOW() - INTERVAL 8 DAY, NOW()),
(24, 'BYC2026032400001024', 5, 2, 194.00, 'shipped', '快递柜自提', '丰巢快递', 'FC123123123', NOW() - INTERVAL 8 DAY, NOW()),
(25, 'BYC2026032300001025', 6, 3, 98.00, 'shipped', '', '京东快递', 'JD5555666677', NOW() - INTERVAL 9 DAY, NOW()),
(26, 'BYC2026032300001026', 7, 4, 640.00, 'shipped', '', '顺丰速运', 'SF8888999900', NOW() - INTERVAL 9 DAY, NOW()),
(27, 'BYC2026032200001027', 8, 5, 644.00, 'shipped', '家人代收', 'EMS', 'EM888999000', NOW() - INTERVAL 10 DAY, NOW()),
(28, 'BYC2026032200001028', 9, 6, 1288.00, 'shipped', '', '顺丰速运', 'SF1111333355', NOW() - INTERVAL 10 DAY, NOW()),
(29, 'BYC2026032100001029', 10, 7, 584.00, 'shipped', '白天不在', '京东快递', 'JD2222333344', NOW() - INTERVAL 11 DAY, NOW()),
(30, 'BYC2026032100001030', 11, 8, 666.00, 'shipped', '', '顺丰速运', 'SF5555667788', NOW() - INTERVAL 11 DAY, NOW()),
(31, 'BYC2026032000001031', 12, 9, 158.00, 'shipped', '请敲门', '德邦物流', 'DB6666777788', NOW() - INTERVAL 12 DAY, NOW()),
(32, 'BYC2026032000001032', 13, 10, 594.00, 'shipped', '', '顺丰速运', 'SF9999000011', NOW() - INTERVAL 12 DAY, NOW()),
(33, 'BYC2026031900001033', 14, 11, 158.00, 'shipped', '', 'EMS', 'EM111222333', NOW() - INTERVAL 13 DAY, NOW()),
(34, 'BYC2026031900001034', 15, 12, 336.00, 'shipped', '周末在家', '京东快递', 'JD444555666', NOW() - INTERVAL 13 DAY, NOW()),

-- Completed orders (15)
(35, 'BYC2026031500001035', 16, 13, 256.00, 'completed', '非常满意', '', '', NOW() - INTERVAL 16 DAY, NOW()),
(36, 'BYC2026031200001036', 17, 14, 354.00, 'completed', '', '', '', NOW() - INTERVAL 19 DAY, NOW()),
(37, 'BYC2026031000001037', 18, 15, 176.00, 'completed', '', '', '', NOW() - INTERVAL 21 DAY, NOW()),
(38, 'BYC2026030700001038', 19, 16, 36.00, 'completed', '好评', '', '', NOW() - INTERVAL 24 DAY, NOW()),
(39, 'BYC2026030500001039', 20, 17, 366.00, 'completed', '', '', '', NOW() - INTERVAL 26 DAY, NOW()),
(40, 'BYC2026030200001040', 21, 18, 168.00, 'completed', '谢谢', '', '', NOW() - INTERVAL 29 DAY, NOW()),
(41, 'BYC2026022800001041', 22, 19, 336.00, 'completed', '', '', '', NOW() - INTERVAL 33 DAY, NOW()),
(42, 'BYC2026022500001042', 23, 20, 1288.00, 'completed', '还会再来', '', '', NOW() - INTERVAL 36 DAY, NOW()),
(43, 'BYC2026022000001043', 24, 21, 496.00, 'completed', '', '', '', NOW() - INTERVAL 41 DAY, NOW()),
(44, 'BYC2026021500001044', 25, 22, 584.00, 'completed', '', '', '', NOW() - INTERVAL 46 DAY, NOW()),
(45, 'BYC2026021000001045', 4, 26, 426.00, 'completed', '', '', '', NOW() - INTERVAL 51 DAY, NOW()),
(46, 'BYC2026020500001046', 5, 2, 516.00, 'completed', '品质一流', '', '', NOW() - INTERVAL 56 DAY, NOW()),
(47, 'BYC2026013100001047', 6, 3, 258.00, 'completed', '', '', '', NOW() - INTERVAL 61 DAY, NOW()),
(48, 'BYC2026012600001048', 7, 4, 386.00, 'completed', '', '', '', NOW() - INTERVAL 66 DAY, NOW()),
(49, 'BYC2026012000001049', 8, 5, 666.00, 'completed', '非常满意', '', '', NOW() - INTERVAL 72 DAY, NOW()),

-- Cancelled orders (6)
(50, 'BYC2026032800001050', 9, 6, 158.00, 'cancelled', '不想要了', '', '', NOW() - INTERVAL 3 DAY, NOW()),
(51, 'BYC2026032600001051', 10, 7, 258.00, 'cancelled', '', '', '', NOW() - INTERVAL 5 DAY, NOW()),
(52, 'BYC2026032300001052', 11, 8, 386.00, 'cancelled', '价格太高', '', '', NOW() - INTERVAL 8 DAY, NOW()),
(53, 'BYC2026031900001053', 12, 9, 216.00, 'cancelled', '', '', '', NOW() - INTERVAL 12 DAY, NOW()),
(54, 'BYC2026032500001054', 13, 10, 28.00, 'cancelled', '地址填错', '', '', NOW() - INTERVAL 6 DAY, NOW());

-- =============================================
-- ORDER_ITEMS (~100 rows)
-- =============================================
INSERT INTO order_items (order_id, product_id, quantity, price, created_at) VALUES
-- Order 1: (cohiba-siglo-i x2, cohiba-siglo-ii x1)
(1, 1, 2, 128.00, NOW()),
(1, 2, 1, 158.00, NOW()),
-- Order 2: (cohiba-siglo-iv x1)
(2, 4, 1, 258.00, NOW()),
-- Order 3: (davidoff-escurio x1)
(3, 27, 1, 168.00, NOW()),
-- Order 4: (cutter-v-cut x1, cutter-punch x1)
(4, 40, 1, 88.00, NOW()),
(4, 41, 1, 128.00, NOW()),
-- Order 5: (lighter-wood-matches x1)
(5, 51, 1, 28.00, NOW()),
-- Order 6: (davidoff-wsc x1)
(6, 26, 1, 458.00, NOW()),
-- Order 7: (cohiba-siglo-v x1, Oliva-v-series x1)
(7, 5, 1, 328.00, NOW()),
(7, 37, 1, 198.00, NOW()),
-- Order 8: (fuente-opus-x x1)
(8, 30, 1, 688.00, NOW()),

-- Order 9: (davidoff-millennium x1)
(9, 25, 1, 298.00, NOW()),
-- Order 10: (humidor-50 x1, humidor-cedar-liner x1)
(10, 44, 1, 388.00, NOW()),
(10, 49, 1, 128.00, NOW()),
(10, 47, 1, 88.00, NOW()),
-- Order 11: (cohiba-siglo-i x1, cohiba-short x1)
(11, 1, 1, 128.00, NOW()),
(11, 6, 1, 68.00, NOW()),
-- Order 12: (davidoff-escurio x1)
(12, 27, 1, 168.00, NOW()),
(12, 29, 1, 228.00, NOW()),
-- Order 13: (cohiba-tasting-box x1)
(13, 59, 1, 666.00, NOW()),

-- Order 14: (cohiba-siglo-ii x1, monte-no2 x3)
(14, 2, 1, 158.00, NOW()),
(14, 11, 3, 108.00, NOW()),
-- Order 15: (cohiba-short x2)
(15, 6, 2, 68.00, NOW()),
-- Order 16: (partagas-d4 x1, my-father x1)
(16, 16, 1, 98.00, NOW()),
(16, 17, 1, 128.00, NOW()),
-- Order 17: (romeo-no2 x1)
(17, 20, 1, 88.00, NOW()),
-- Order 18: (butane-lighter x1, wood-matches x2)
(18, 50, 1, 158.00, NOW()),
(18, 51, 2, 28.00, NOW()),
-- Order 19: (davidoff-wsc x1, Oliva-v-series x1, aj-fernandez x1)
(19, 26, 1, 458.00, NOW()),
(19, 37, 1, 198.00, NOW()),
(19, 38, 1, 168.00, NOW()),
-- Order 20: (davidoff-escurio x1, fuente-hemingway x1)
(20, 27, 1, 168.00, NOW()),
(20, 31, 1, 188.00, NOW()),
-- Order 21: (cohiba-tasting-box x1)
(21, 59, 1, 328.00, NOW()),
-- Order 22: (cohiba-siglo-i x1, cohiba-siglo-iv x1)
(22, 1, 1, 128.00, NOW()),
(22, 4, 1, 258.00, NOW()),

-- Order 23: (cohiba-short x2)
(23, 6, 2, 68.00, NOW()),
-- Order 24: (monte-no2 x1, partagas-d4 x2)
(24, 11, 1, 108.00, NOW()),
(24, 16, 2, 98.00, NOW()),
-- Order 25: (partagas-d4 x1)
(25, 16, 1, 98.00, NOW()),
-- Order 26: (cohiba-siglo-i x5)
(26, 1, 5, 128.00, NOW()),
-- Order 27: (humidor-50 x1, humidor-cedar-liner x2)
(27, 44, 1, 388.00, NOW()),
(27, 49, 2, 128.00, NOW()),
-- Order 28: (gift-luxury-box x1)
(28, 61, 1, 1288.00, NOW()),
-- Order 29: (humidor-cabinet-200 x1)
(29, 46, 1, 2800.00, NOW()),
-- Order 30: (cohiba-tasting-box x1)
(30, 59, 1, 666.00, NOW()),
-- Order 31: (butane-lighter x1)
(31, 50, 1, 158.00, NOW()),
-- Order 32: (davidoff-escurio x1, fuente-hemingway x1, fuente-don-carlos x1)
(32, 27, 1, 168.00, NOW()),
(32, 31, 1, 188.00, NOW()),
(32, 29, 1, 228.00, NOW()),
-- Order 33: (cohiba-siglo-ii x3)
(33, 2, 3, 158.00, NOW()),
-- Order 34: (cohiba-siglo-iv x1, cohiba-short x1)
(34, 4, 1, 258.00, NOW()),
(34, 6, 1, 68.00, NOW()),

-- Order 35: (cohiba-siglo-i x2)
(35, 1, 2, 128.00, NOW()),
-- Order 36: (monte-no2 x1, monte-edmundo x3)
(36, 11, 1, 108.00, NOW()),
(36, 13, 3, 138.00, NOW()),
-- Order 37: (cohiba-siglo-iv x1, davidoff-escurio x1)
(37, 4, 1, 258.00, NOW()),
(37, 27, 1, 168.00, NOW()),
-- Order 38: (cutter-v-cut x2)
(38, 40, 2, 88.00, NOW()),
-- Order 39: (davidoff-escurio x1, fuente-hemingway x1)
(39, 27, 1, 168.00, NOW()),
(39, 31, 1, 188.00, NOW()),
-- Order 40: (butane-lighter x1)
(40, 50, 1, 158.00, NOW()),
-- Order 41: (humidor-50 x1, humidor-cedar-liner x1)
(41, 44, 1, 388.00, NOW()),
(41, 49, 1, 128.00, NOW()),
-- Order 42: (gift-luxury-box x1)
(42, 61, 1, 1288.00, NOW()),
-- Order 43: (cohiba-tasting-box x1)
(43, 59, 1, 666.00, NOW()),
-- Order 44: (humidor-travel-5 x1, starter-pack-5 x2)
(44, 45, 1, 168.00, NOW()),
(44, 58, 2, 198.00, NOW()),
-- Order 45: (cohiba-siglo-i x1, cohiba-siglo-ii x1, cohiba-siglo-iii x1)
(45, 1, 1, 128.00, NOW()),
(45, 2, 1, 158.00, NOW()),
(45, 3, 1, 198.00, NOW()),
-- Order 46: (cohiba-siglo-iv x2)
(46, 4, 2, 258.00, NOW()),
-- Order 47: (davidoff-escurio x1)
(47, 27, 1, 168.00, NOW()),
(47, 31, 1, 188.00, NOW()),
-- Order 48: ( Oliva-v-series x1, aj-fernandez x1)
(48, 37, 1, 198.00, NOW()),
(48, 38, 1, 168.00, NOW()),
-- Order 49: (cohiba-tasting-box x1, travel-cigar-set x1)
(49, 59, 1, 666.00, NOW()),
(49, 60, 1, 358.00, NOW()),

-- Cancelled orders
-- Order 50: (davidoff-millennium x1)
(50, 25, 1, 298.00, NOW()),
-- Order 51: (shengsaiyadi x1)
(51, 32, 1, 258.00, NOW()),
-- Order 52: (cohiba-siglo-i x1, davidoff-escurio x1)
(52, 1, 1, 128.00, NOW()),
(52, 27, 1, 168.00, NOW()),
-- Order 53: (cutter-v-cut x1, cutter-punch x2)
(53, 40, 1, 88.00, NOW()),
(53, 41, 2, 58.00, NOW()),
-- Order 54: (lighter-wood-matches x1)
(54, 51, 1, 28.00, NOW());

-- =============================================
-- PAYMENT_PROOFS (19 rows)
-- =============================================
INSERT INTO payment_proofs (id, order_id, user_id, payment_method_id, image_url, status, reject_reason, reviewer_id, reviewed_at, created_at, updated_at) VALUES
-- Pending proofs for pending orders (6)
(1, 1, 4, 1, 'https://picsum.photos/seed/proof00/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(2, 2, 5, 2, 'https://picsum.photos/seed/proof01/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(3, 3, 6, 1, 'https://picsum.photos/seed/proof02/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(4, 4, 7, 2, 'https://picsum.photos/seed/proof03/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(5, 5, 8, 3, 'https://picsum.photos/seed/proof04/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(6, 6, 9, 1, 'https://picsum.photos/seed/proof05/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),

-- Pending proofs for paid orders (6)
(7, 7, 10, 2, 'https://picsum.photos/seed/proof06/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(8, 8, 11, 1, 'https://picsum.photos/seed/proof07/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(9, 9, 12, 2, 'https://picsum.photos/seed/proof08/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(10, 10, 13, 3, 'https://picsum.photos/seed/proof09/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(11, 11, 14, 1, 'https://picsum.photos/seed/proof10/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),
(12, 12, 15, 2, 'https://picsum.photos/seed/proof11/400/300', 'pending', NULL, NULL, NULL, NOW(), NOW()),

-- Rejected proofs (3)
(13, 13, 16, 1, 'https://picsum.photos/seed/proof12/400/300', 'rejected', '图片不清晰，请重新上传', 2, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY),
(14, 14, 17, 2, 'https://picsum.photos/seed/proof13/400/300', 'rejected', '付款金额与订单不符', 2, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY),
(15, 15, 18, 3, 'https://picsum.photos/seed/proof14/400/300', 'rejected', '未填写订单备注', 2, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY),

-- Approved proofs for processing orders (7)
(16, 16, 19, 1, 'https://picsum.photos/seed/proofapproved00/400/300', 'approved', NULL, 2, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 3 DAY),
(17, 17, 20, 2, 'https://picsum.photos/seed/proofapproved01/400/300', 'approved', NULL, 2, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 3 DAY),
(18, 18, 21, 1, 'https://picsum.photos/seed/proofapproved02/400/300', 'approved', NULL, 2, NOW() - INTERVAL 4 DAY, NOW() - INTERVAL 4 DAY),
(19, 19, 22, 2, 'https://picsum.photos/seed/proofapproved03/400/300', 'approved', NULL, 2, NOW() - INTERVAL 4 DAY, NOW() - INTERVAL 4 DAY),
(20, 20, 23, 3, 'https://picsum.photos/seed/proofapproved04/400/300', 'approved', NULL, 2, NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 5 DAY),
(21, 21, 24, 1, 'https://picsum.photos/seed/proofapproved05/400/300', 'approved', NULL, 2, NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 5 DAY),
(22, 22, 25, 2, 'https://picsum.photos/seed/proofapproved06/400/300', 'approved', NULL, 2, NOW() - INTERVAL 6 DAY, NOW() - INTERVAL 6 DAY);

-- =============================================
-- CART_ITEMS (~40 rows)
-- =============================================
INSERT INTO cart_items (user_id, product_id, quantity, created_at, updated_at) VALUES
-- User 4 (张伟)
(4, 1, 2, NOW(), NOW()),
(4, 12, 1, NOW(), NOW()),
(4, 24, 1, NOW(), NOW()),
-- User 5 (李娜)
(5, 4, 1, NOW(), NOW()),
(5, 25, 1, NOW(), NOW()),
-- User 6 (王芳)
(6, 6, 3, NOW(), NOW()),
-- User 7 (刘洋)
(7, 11, 1, NOW(), NOW()),
(7, 26, 1, NOW(), NOW()),
(7, 37, 1, NOW(), NOW()),
-- User 8 (陈杰)
(8, 16, 2, NOW(), NOW()),
-- User 9 (杨秀英)
(9, 20, 1, NOW(), NOW()),
(9, 30, 1, NOW(), NOW()),
(9, 40, 1, NOW(), NOW()),
-- User 10 (黄明)
(10, 1, 1, NOW(), NOW()),
(10, 44, 1, NOW(), NOW()),
-- User 11 (赵磊)
(11, 5, 1, NOW(), NOW()),
(11, 59, 1, NOW(), NOW()),
-- User 12 (周静)
(12, 40, 1, NOW(), NOW()),
-- User 13 (吴强)
(13, 1, 2, NOW(), NOW()),
(13, 2, 1, NOW(), NOW()),
(13, 3, 1, NOW(), NOW()),
-- User 14 (徐丽)
(14, 45, 1, NOW(), NOW()),
(14, 50, 1, NOW(), NOW()),
-- User 15 (孙浩)
(15, 24, 1, NOW(), NOW()),
-- User 16 (马超)
(16, 33, 1, NOW(), NOW()),
(16, 58, 1, NOW(), NOW()),
-- User 17 (胡蝶)
(17, 4, 1, NOW(), NOW()),
-- User 18 (朱琳)
(18, 11, 2, NOW(), NOW()),
(18, 27, 1, NOW(), NOW()),
-- User 19 (郭锐)
(19, 40, 1, NOW(), NOW()),
(19, 51, 2, NOW(), NOW()),
-- User 20 (林涛)
(20, 1, 1, NOW(), NOW()),
(20, 46, 1, NOW(), NOW()),
-- User 21 (何雪)
(21, 20, 1, NOW(), NOW()),
(21, 61, 1, NOW(), NOW()),
-- User 22 (高建)
(22, 5, 1, NOW(), NOW()),
(22, 59, 1, NOW(), NOW()),
-- User 23 (罗岚)
(23, 44, 1, NOW(), NOW()),
-- User 24 (梁志明)
(24, 58, 2, NOW(), NOW()),
-- User 25 (宋雨晴)
(25, 1, 1, NOW(), NOW()),
(25, 60, 1, NOW(), NOW());

-- =============================================
-- FAVORITES (~60 rows)
-- =============================================
INSERT INTO favorites (user_id, product_id, created_at) VALUES
-- User 4 (张伟)
(4, 2, NOW()), (4, 13, NOW()), (4, 25, NOW()), (4, 37, NOW()),
-- User 5 (李娜)
(5, 5, NOW()), (5, 14, NOW()), (5, 26, NOW()), (5, 59, NOW()),
-- User 6 (王芳)
(6, 7, NOW()), (6, 20, NOW()), (6, 30, NOW()),
-- User 7 (刘洋)
(7, 1, NOW()), (7, 11, NOW()), (7, 33, NOW()), (7, 44, NOW()), (7, 61, NOW()),
-- User 8 (陈杰)
(8, 4, NOW()), (8, 15, NOW()), (8, 27, NOW()),
-- User 9 (杨秀英)
(9, 6, NOW()), (9, 21, NOW()), (9, 31, NOW()), (9, 46, NOW()),
-- User 10 (黄明)
(10, 8, NOW()), (10, 22, NOW()), (10, 34, NOW()), (10, 58, NOW()),
-- User 11 (赵磊)
(11, 1, NOW()), (11, 12, NOW()), (11, 24, NOW()), (11, 37, NOW()), (11, 59, NOW()),
-- User 12 (周静)
(12, 40, NOW()), (12, 45, NOW()), (12, 50, NOW()),
-- User 13 (吴强)
(13, 5, NOW()), (13, 16, NOW()), (13, 26, NOW()), (13, 33, NOW()),
-- User 14 (徐丽)
(14, 1, NOW()), (14, 20, NOW()), (14, 38, NOW()), (14, 44, NOW()), (14, 60, NOW()),
-- User 15 (孙浩)
(15, 11, NOW()), (15, 27, NOW()), (15, 40, NOW()),
-- User 16 (马超)
(16, 2, NOW()), (16, 13, NOW()), (16, 30, NOW()), (16, 46, NOW()), (16, 61, NOW()),
-- User 17 (胡蝶)
(17, 6, NOW()), (17, 17, NOW()), (17, 33, NOW()),
-- User 18 (朱琳)
(18, 4, NOW()), (18, 14, NOW()), (18, 25, NOW()), (18, 37, NOW()), (18, 58, NOW()),
-- User 19 (郭锐)
(19, 7, NOW()), (19, 22, NOW()), (19, 34, NOW()),
-- User 20 (林涛)
(20, 1, NOW()), (20, 15, NOW()), (20, 26, NOW()), (20, 44, NOW()),
-- User 21 (何雪)
(21, 8, NOW()), (21, 21, NOW()), (21, 35, NOW()), (21, 59, NOW()),
-- User 22 (高建)
(22, 5, NOW()), (22, 16, NOW()), (22, 27, NOW()), (22, 38, NOW()),
-- User 23 (罗岚)
(23, 11, NOW()), (23, 30, NOW()), (23, 45, NOW()),
-- User 24 (梁志明)
(24, 1, NOW()), (24, 20, NOW()), (24, 33, NOW()), (24, 46, NOW()), (24, 61, NOW()),
-- User 25 (宋雨晴)
(25, 6, NOW()), (25, 17, NOW()), (25, 40, NOW());

-- =============================================
-- CONVERSATIONS (20 rows: 12 open + 8 closed)
-- =============================================
INSERT INTO conversations (id, user_id, status, assigned_to, last_message_at, created_at, updated_at) VALUES
-- Open conversations (12)
(1, 4, 'open', 2, NOW(), NOW() - INTERVAL 2 DAY, NOW()),
(2, 5, 'open', 3, NOW(), NOW() - INTERVAL 1 DAY, NOW()),
(3, 6, 'open', 2, NOW(), NOW() - INTERVAL 0 DAY, NOW()),
(4, 7, 'open', 3, NOW(), NOW() - INTERVAL 3 DAY, NOW()),
(5, 8, 'open', 2, NOW(), NOW() - INTERVAL 1 DAY, NOW()),
(6, 9, 'open', 3, NOW(), NOW() - INTERVAL 2 DAY, NOW()),
(7, 10, 'open', 2, NOW(), NOW() - INTERVAL 0 DAY, NOW()),
(8, 11, 'open', 3, NOW(), NOW() - INTERVAL 1 DAY, NOW()),
(9, 12, 'open', 2, NOW(), NOW() - INTERVAL 2 DAY, NOW()),
(10, 13, 'open', 3, NOW(), NOW() - INTERVAL 0 DAY, NOW()),
(11, 14, 'open', 2, NOW(), NOW() - INTERVAL 1 DAY, NOW()),
(12, 15, 'open', 3, NOW(), NOW() - INTERVAL 0 DAY, NOW()),

-- Closed conversations (8)
(13, 16, 'closed', 2, NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 7 DAY, NOW() - INTERVAL 5 DAY),
(14, 17, 'closed', 3, NOW() - INTERVAL 6 DAY, NOW() - INTERVAL 8 DAY, NOW() - INTERVAL 6 DAY),
(15, 18, 'closed', 2, NOW() - INTERVAL 4 DAY, NOW() - INTERVAL 6 DAY, NOW() - INTERVAL 4 DAY),
(16, 19, 'closed', 3, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 5 DAY, NOW() - INTERVAL 3 DAY),
(17, 20, 'closed', 2, NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 4 DAY, NOW() - INTERVAL 2 DAY),
(18, 21, 'closed', 3, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 1 DAY),
(19, 22, 'closed', 2, NOW() - INTERVAL 0 DAY, NOW() - INTERVAL 2 DAY, NOW()),
(20, 23, 'closed', 3, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY);

-- =============================================
-- MESSAGES (~80 rows)
-- =============================================
INSERT INTO messages (conversation_id, sender_type, sender_id, message_type, content, thumbnail_url, is_read, created_at) VALUES
-- Conversation 1 (user 4, service 2)
(1, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 2 DAY),
(1, 'customer', 4, 'text', '你好，我想咨询一下高希霸世纪一号的口感', '', false, NOW() - INTERVAL 2 DAY),
(1, 'service', 2, 'text', '高希霸世纪一号是我们最受欢迎的入门款，口感温和细腻，带有淡淡的奶油和咖啡香气。环规40，长度5英寸，非常适合初次品鉴古巴雪茄的朋友。', '', true, NOW() - INTERVAL 2 DAY),
(1, 'customer', 4, 'text', '那库存现在有货吗？', '', false, NOW() - INTERVAL 2 DAY),
(1, 'service', 2, 'text', '目前库存充足，我们有50支现货，可以当天发货。', '', true, NOW() - INTERVAL 2 DAY),
(1, 'customer', 4, 'text', '好的，我下单了', '', false, NOW() - INTERVAL 1 DAY),

-- Conversation 2 (user 5, service 3)
(2, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 1 DAY),
(2, 'customer', 5, 'text', '请问有适合送礼的礼盒吗？', '', false, NOW() - INTERVAL 1 DAY),
(2, 'service', 3, 'text', '当然，我们有几款非常适合送礼的礼盒。推荐您看看我们的「送礼豪华礼盒」，内含多支精选雪茄，配有精美包装，是送礼的首选。', '', true, NOW() - INTERVAL 1 DAY),
(2, 'customer', 5, 'text', '送礼豪华礼盒的价格是多少？', '', false, NOW() - INTERVAL 1 DAY),
(2, 'service', 3, 'text', '送礼豪华礼盒定价1288元，目前还有5套现货。', '', true, NOW() - INTERVAL 1 DAY),
(2, 'customer', 5, 'text', '好的，我去看看', '', false, NOW()),
(2, 'service', 3, 'text', '好的，有任何问题随时联系我。祝您购物愉快！', '', true, NOW()),

-- Conversation 3 (user 6, service 2)
(3, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW()),
(3, 'customer', 6, 'text', '我的订单怎么还没发货？', '', false, NOW()),
(3, 'service', 2, 'text', '您好，让我帮您查询一下订单状态。请问您的订单号是多少？', '', true, NOW()),
(3, 'customer', 6, 'text', '好的稍等我看一下', '', false, NOW()),

-- Conversation 4 (user 7, service 3)
(4, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 3 DAY),
(4, 'customer', 7, 'text', '请问雪茄剪双刃不锈钢和V口哪个好？', '', false, NOW() - INTERVAL 3 DAY),
(4, 'service', 3, 'text', '各有优势。双刃剪适合一次性整齐切割，适合大多数雪茄；V口剪则切割深度较浅，适合环规较粗的雪茄。', '', true, NOW() - INTERVAL 3 DAY),
(4, 'customer', 7, 'text', '明白了，谢谢', '', false, NOW() - INTERVAL 2 DAY),
(4, 'service', 3, 'text', '不客气！欢迎下次光临。', '', true, NOW() - INTERVAL 2 DAY),

-- Conversation 5 (user 8, service 2)
(5, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 1 DAY),
(5, 'customer', 8, 'text', '你们支持货到付款吗？', '', false, NOW() - INTERVAL 1 DAY),
(5, 'service', 2, 'text', '支持的，我们支持货到付款、微信支付、支付宝、银行转账以及PayPal等多种支付方式。', '', true, NOW() - INTERVAL 1 DAY),
(5, 'customer', 8, 'text', '好的，那发货用什么快递？', '', false, NOW() - INTERVAL 1 DAY),
(5, 'service', 2, 'text', '我们主要使用顺丰速运和京东快递，大部分城市支持隔日达。', '', true, NOW() - INTERVAL 1 DAY),
(5, 'customer', 8, 'text', '明白了，已下单', '', false, NOW()),

-- Conversation 6 (user 9, service 3)
(6, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 2 DAY),
(6, 'customer', 9, 'text', '保湿盒50支装的是什么木材？', '', false, NOW() - INTERVAL 2 DAY),
(6, 'service', 3, 'text', '我们50支装的桌面保湿盒采用的是西班牙雪松木，具有天然的保湿和驱虫效果。', '', true, NOW() - INTERVAL 2 DAY),
(6, 'customer', 9, 'text', '价格呢？', '', false, NOW() - INTERVAL 2 DAY),
(6, 'service', 3, 'text', '桌面保湿盒50支装定价388元，含电子湿度计。顺丰包邮。', '', true, NOW() - INTERVAL 2 DAY),
(6, 'customer', 9, 'text', '有保修吗？', '', false, NOW() - INTERVAL 1 DAY),
(6, 'service', 3, 'text', '我们提供一年质保服务，非人为损坏可以免费维修或更换。', '', true, NOW() - INTERVAL 1 DAY),

-- Conversation 7 (user 10, service 2)
(7, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW()),
(7, 'customer', 10, 'text', '想问问帕德龙1964周年口感如何', '', false, NOW()),
(7, 'service', 2, 'text', '帕德龙1964是尼加拉瓜最经典的雪茄之一，采用该国最好的茄叶，带有浓郁的咖啡、巧克力和香料风味。', '', true, NOW()),
(7, 'customer', 10, 'text', '有陈年潜力吗？', '', false, NOW()),
(7, 'service', 2, 'text', '非常有。帕德龙1964建议陈年5年以上，风味会更加圆润醇厚。', '', true, NOW()),

-- Conversation 8 (user 11, service 3)
(8, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 1 DAY),
(8, 'customer', 11, 'text', '请问有图片吗想看看', '', false, NOW() - INTERVAL 1 DAY),
(8, 'service', 3, 'image', 'https://picsum.photos/seed/chatimg01/400/300', 'https://picsum.photos/seed/chatimg01/200/200', true, NOW() - INTERVAL 1 DAY),
(8, 'service', 3, 'text', '这是商品实拍图，您可以看到细节。', '', true, NOW() - INTERVAL 1 DAY),
(8, 'customer', 11, 'text', '收到，看起来不错', '', false, NOW()),

-- Conversation 9 (user 12, service 2)
(9, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 2 DAY),
(9, 'customer', 12, 'text', '想退换货怎么办理', '', false, NOW() - INTERVAL 2 DAY),
(9, 'service', 2, 'text', '您好，雪茄属于特殊商品，若因质量问题可在收货48小时内申请退换。请提供照片凭证。', '', true, NOW() - INTERVAL 2 DAY),
(9, 'customer', 12, 'text', '好的明白了', '', false, NOW() - INTERVAL 1 DAY),

-- Conversation 10 (user 13, service 3)
(10, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW()),
(10, 'customer', 13, 'text', '我是老客户了有优惠吗', '', false, NOW()),
(10, 'service', 3, 'text', '感谢您一直以来的支持！老客户我们有专属折扣码，可以享受9折优惠。', '', true, NOW()),
(10, 'customer', 13, 'text', '好的谢谢', '', false, NOW()),

-- Conversation 11 (user 14, service 2)
(11, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 1 DAY),
(11, 'customer', 14, 'text', '入门体验装适合新手吗', '', false, NOW() - INTERVAL 1 DAY),
(11, 'service', 2, 'text', '非常适合！入门体验装精选5支不同品牌的经典雪茄，口感从温和到浓郁都有覆盖，还附赠专业品鉴指南。', '', true, NOW() - INTERVAL 1 DAY),
(11, 'customer', 14, 'text', '价格多少', '', false, NOW() - INTERVAL 1 DAY),
(11, 'service', 2, 'text', '定价198元，新客户首单还有额外9折。', '', true, NOW() - INTERVAL 1 DAY),
(11, 'customer', 14, 'text', '已下单', '', false, NOW()),

-- Conversation 12 (user 15, service 3)
(12, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW()),
(12, 'customer', 15, 'text', '能开发票吗', '', false, NOW()),
(12, 'service', 3, 'text', '可以开普通发票或增值税专用发票，请在下单时备注发票抬头和税号。', '', true, NOW()),
(12, 'customer', 15, 'text', '好的', '', false, NOW()),

-- Closed conversations (13-20) - simpler messages
(13, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 7 DAY),
(13, 'customer', 16, 'text', '你好', '', true, NOW() - INTERVAL 7 DAY),
(13, 'service', 2, 'text', '您好，请问有什么可以帮您？', '', true, NOW() - INTERVAL 7 DAY),
(13, 'customer', 16, 'text', '我想问一下物流时效', '', true, NOW() - INTERVAL 7 DAY),
(13, 'service', 2, 'text', '顺丰快递一般1-3天送达，偏远地区3-5天。', '', true, NOW() - INTERVAL 7 DAY),
(13, 'customer', 16, 'text', '好的谢谢', '', true, NOW() - INTERVAL 7 DAY),
(13, 'service', 2, 'text', '不客气！祝您生活愉快！', '', true, NOW() - INTERVAL 7 DAY),
(13, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 5 DAY),

(14, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 8 DAY),
(14, 'customer', 17, 'text', '请问有蒙特4号吗', '', true, NOW() - INTERVAL 8 DAY),
(14, 'service', 3, 'text', '有的，蒙特4号库存80支，正常发货。', '', true, NOW() - INTERVAL 8 DAY),
(14, 'customer', 17, 'text', '好的谢谢', '', true, NOW() - INTERVAL 8 DAY),
(14, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 6 DAY),

(15, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 6 DAY),
(15, 'customer', 18, 'text', '你们支持退换吗', '', true, NOW() - INTERVAL 6 DAY),
(15, 'service', 2, 'text', '质量问题可以48小时内申请退换。', '', true, NOW() - INTERVAL 6 DAY),
(15, 'customer', 18, 'text', '明白了', '', true, NOW() - INTERVAL 6 DAY),
(15, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 4 DAY),

(16, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 5 DAY),
(16, 'customer', 19, 'text', '发货地是哪里', '', true, NOW() - INTERVAL 5 DAY),
(16, 'service', 3, 'text', '我们从上海仓库发货，顺丰隔日达。', '', true, NOW() - INTERVAL 5 DAY),
(16, 'customer', 19, 'text', '好的', '', true, NOW() - INTERVAL 5 DAY),
(16, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 3 DAY),

(17, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 4 DAY),
(17, 'customer', 20, 'text', '大卫杜夫2000多少钱', '', true, NOW() - INTERVAL 4 DAY),
(17, 'service', 2, 'text', '大卫杜夫2000定价198元。', '', true, NOW() - INTERVAL 4 DAY),
(17, 'customer', 20, 'text', 'ok', '', true, NOW() - INTERVAL 4 DAY),
(17, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 2 DAY),

(18, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 3 DAY),
(18, 'customer', 21, 'text', '怎么联系你们', '', true, NOW() - INTERVAL 3 DAY),
(18, 'service', 3, 'text', '电话400-888-9999，微信BYCIGAR_CS。', '', true, NOW() - INTERVAL 3 DAY),
(18, 'customer', 21, 'text', '知道了', '', true, NOW() - INTERVAL 3 DAY),
(18, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 1 DAY),

(19, 'service', 2, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 2 DAY),
(19, 'customer', 22, 'text', '礼盒可以定制吗', '', true, NOW() - INTERVAL 2 DAY),
(19, 'service', 2, 'text', '部分礼盒支持定制，详情请咨询客服。', '', true, NOW() - INTERVAL 2 DAY),
(19, 'customer', 22, 'text', '好的谢谢', '', true, NOW() - INTERVAL 2 DAY),
(19, 'system', 0, 'text', '客服已结束对话', '', true, NOW()),

(20, 'service', 3, 'text', '您好！欢迎来到 BYCIGAR，有什么可以帮助您的吗？', '', true, NOW() - INTERVAL 1 DAY),
(20, 'customer', 23, 'text', '新会员有什么优惠', '', true, NOW() - INTERVAL 1 DAY),
(20, 'service', 3, 'text', '新会员首单免运费，另有专属折扣码。', '', true, NOW() - INTERVAL 1 DAY),
(20, 'customer', 23, 'text', '好的', '', true, NOW() - INTERVAL 1 DAY),
(20, 'system', 0, 'text', '客服已结束对话', '', true, NOW() - INTERVAL 1 DAY);

-- =============================================
-- NOTIFICATIONS (~100 rows)
-- =============================================
INSERT INTO notifications (user_id, type, title, content, is_read, link, product_id, order_id, created_at) VALUES
-- User 4 (张伟)
(4, 'order_status', '订单已发货', '您的订单已由顺丰速运发出，请注意查收。', false, '/orders', NULL, 23, NOW() - INTERVAL 8 DAY),
(4, 'back_in_stock', '商品到货通知', '您关注的「高希霸世纪二号」已到货，欢迎购买！', false, '/products/cohiba-siglo-ii', 2, NULL, NOW() - INTERVAL 5 DAY),
(4, 'price_drop', '价格下调通知', '您关注的「大卫杜夫2000」价格下调，由¥240调整为¥198，机会不容错过！', false, '/products/davidoff-2000', 24, NULL, NOW() - INTERVAL 3 DAY),
(4, 'order_status', '订单已完成', '感谢您的购买，期待下次光临！', true, '/orders', NULL, 35, NOW() - INTERVAL 16 DAY),

-- User 5 (李娜)
(5, 'order_status', '订单已发货', '您的订单已发货，快递正在配送中。', false, '/orders', NULL, 24, NOW() - INTERVAL 8 DAY),
(5, 'back_in_stock', '商品到货通知', '您关注的「大卫杜夫千年系列」已到货！', true, '/products/davidoff-millennium', 25, NULL, NOW() - INTERVAL 10 DAY),
(5, 'price_drop', '价格下调通知', '您关注的「帕德龙1964周年」价格下调，由¥330调整为¥278！', false, '/products/padron-1964', 33, NULL, NOW() - INTERVAL 2 DAY),
(5, 'order_status', '订单已完成', '感谢您的购买，期待下次光临！', true, '/orders', NULL, 36, NOW() - INTERVAL 19 DAY),

-- User 6 (王芳)
(6, 'order_status', '订单已发货', '您的订单已发货。', false, '/orders', NULL, 25, NOW() - INTERVAL 9 DAY),
(6, 'back_in_stock', '商品到货通知', '您关注的「大卫杜夫埃斯库里奥」已到货！', false, '/products/davidoff-escurio', 27, NULL, NOW() - INTERVAL 4 DAY),
(6, 'order_status', '订单已完成', '感谢您的购买，期待下次光临！', true, '/orders', NULL, 37, NOW() - INTERVAL 21 DAY),

-- User 7 (刘洋)
(7, 'order_status', '订单已发货', '您的订单已由EMS发出，请注意查收。', false, '/orders', NULL, 27, NOW() - INTERVAL 10 DAY),
(7, 'back_in_stock', '商品到货通知', '您关注的「大卫杜夫温斯顿丘吉尔」已到货！', false, '/products/davidoff-wsc', 26, NULL, NOW() - INTERVAL 6 DAY),
(7, 'price_drop', '价格下调通知', '您关注的「大卫杜夫格兰德」价格下调！', true, '/products/davidoff-grande', 28, NULL, NOW() - INTERVAL 15 DAY),
(7, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 48, NOW() - INTERVAL 66 DAY),

-- User 8 (陈杰)
(8, 'order_status', '订单已发货', '您的订单已发货。', false, '/orders', NULL, 28, NOW() - INTERVAL 10 DAY),
(8, 'back_in_stock', '商品到货通知', '您关注的「高希霸品鉴礼盒」已到货！', false, '/products/cohiba-tasting-box', 59, NULL, NOW() - INTERVAL 3 DAY),
(8, 'order_status', '订单已完成', '感谢您的购买，期待下次光临！', true, '/orders', NULL, 49, NOW() - INTERVAL 72 DAY),

-- User 9 (杨秀英)
(9, 'order_status', '订单已发货', '您的订单已发货。', false, '/orders', NULL, 29, NOW() - INTERVAL 11 DAY),
(9, 'back_in_stock', '商品到货通知', '您关注的「豪华保湿柜200支」已到货！', false, '/products/humidor-cabinet-200', 46, NULL, NOW() - INTERVAL 7 DAY),
(9, 'order_status', '订单已取消', '您的订单已取消。', true, '/orders', NULL, 50, NOW() - INTERVAL 3 DAY),

-- User 10 (黄明)
(10, 'order_status', '订单已发货', '您的订单已由京东快递发出。', false, '/orders', NULL, 30, NOW() - INTERVAL 11 DAY),
(10, 'back_in_stock', '商品到货通知', '您关注的「送礼豪华礼盒」已到货！', false, '/products/gift-luxury-box', 61, NULL, NOW() - INTERVAL 5 DAY),
(10, 'price_drop', '价格下调通知', '您关注的「高希霸贝伊可52号」价格下调！', false, '/products/cohiba-behike-52', 8, NULL, NOW() - INTERVAL 1 DAY),
(10, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 38, NOW() - INTERVAL 24 DAY),

-- User 11 (赵磊)
(11, 'order_status', '订单已发货', '您的订单已发货。', false, '/orders', NULL, 31, NOW() - INTERVAL 12 DAY),
(11, 'back_in_stock', '商品到货通知', '您关注的「雪茄剪V口切割器」已到货！', true, '/products/cutter-v-cut', 40, NULL, NOW() - INTERVAL 14 DAY),
(11, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 39, NOW() - INTERVAL 26 DAY),

-- User 12 (周静)
(12, 'order_status', '订单已发货', '您的订单已由德邦物流发出。', false, '/orders', NULL, 32, NOW() - INTERVAL 12 DAY),
(12, 'back_in_stock', '商品到货通知', '您关注的「帕德龙1926系列80年」已到货！', false, '/products/padron-1926-80', 34, NULL, NOW() - INTERVAL 8 DAY),
(12, 'price_drop', '价格下调通知', '您关注的「入门体验装5支」价格下调，由¥240调整为¥198！', true, '/products/starter-pack-5', 58, NULL, NOW() - INTERVAL 20 DAY),
(12, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 40, NOW() - INTERVAL 29 DAY),

-- User 13 (吴强)
(13, 'order_status', '订单已发货', '您的订单已发货。', false, '/orders', NULL, 33, NOW() - INTERVAL 13 DAY),
(13, 'back_in_stock', '商品到货通知', '您关注的「大卫杜夫2000」已到货！', false, '/products/davidoff-2000', 24, NULL, NOW() - INTERVAL 9 DAY),
(13, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 41, NOW() - INTERVAL 33 DAY),

-- User 14 (徐丽)
(14, 'order_status', '订单已发货', '您的订单已由EMS发出。', false, '/orders', NULL, 34, NOW() - INTERVAL 13 DAY),
(14, 'back_in_stock', '商品到货通知', '您关注的「旅行保湿盒5支」已到货！', false, '/products/humidor-travel-5', 45, NULL, NOW() - INTERVAL 6 DAY),
(14, 'price_drop', '价格下调通知', '您关注的「大卫杜夫温斯顿丘吉尔」价格下调！', true, '/products/davidoff-wsc', 26, NULL, NOW() - INTERVAL 18 DAY),
(14, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 42, NOW() - INTERVAL 36 DAY),

-- User 15 (孙浩)
(15, 'order_status', '订单已发货', '您的订单已由京东快递发出。', false, '/orders', NULL, 35, NOW() - INTERVAL 14 DAY),
(15, 'back_in_stock', '商品到货通知', '您关注的「丁烷气点火器」已到货！', false, '/products/lighter-butane', 50, NULL, NOW() - INTERVAL 10 DAY),
(15, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 43, NOW() - INTERVAL 41 DAY),

-- User 16 (马超)
(16, 'order_status', '订单已完成', '感谢您的购买，期待下次光临！', true, '/orders', NULL, 44, NOW() - INTERVAL 46 DAY),
(16, 'back_in_stock', '商品到货通知', '您关注的「高希霸世纪五号」已到货！', true, '/products/cohiba-siglo-v', 5, NULL, NOW() - INTERVAL 30 DAY),

-- User 17 (胡蝶)
(17, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 45, NOW() - INTERVAL 51 DAY),
(17, 'back_in_stock', '商品到货通知', '您关注的「奥利瓦V系列」已到货！', true, '/products/oliva-v-series', 37, NULL, NOW() - INTERVAL 35 DAY),

-- User 18 (朱琳)
(18, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 46, NOW() - INTERVAL 56 DAY),
(18, 'price_drop', '价格下调通知', '您关注的「AJ费尔南德斯新世界」价格下调！', true, '/products/aj-fernandez-nw', 38, NULL, NOW() - INTERVAL 40 DAY),

-- User 19 (郭锐)
(19, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 47, NOW() - INTERVAL 61 DAY),

-- User 20 (林涛)
(20, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 48, NOW() - INTERVAL 66 DAY),

-- User 21 (何雪)
(21, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 49, NOW() - INTERVAL 72 DAY),

-- User 22 (高建)
(22, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 44, NOW() - INTERVAL 75 DAY),

-- User 23 (罗岚)
(23, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 45, NOW() - INTERVAL 80 DAY),

-- User 24 (梁志明)
(24, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 46, NOW() - INTERVAL 85 DAY),

-- User 25 (宋雨晴)
(25, 'order_status', '订单已完成', '感谢您的购买！', true, '/orders', NULL, 47, NOW() - INTERVAL 90 DAY),
(25, 'back_in_stock', '商品到货通知', '您关注的「高希霸世纪一号」已到货！', true, '/products/cohiba-siglo-i', 1, NULL, NOW() - INTERVAL 60 DAY);

-- =============================================
-- Summary
-- =============================================
SELECT 'Seed completed!' AS status;
SELECT COUNT(*) AS total_users FROM users;
SELECT COUNT(*) AS total_categories FROM categories;
SELECT COUNT(*) AS total_products FROM products;
SELECT COUNT(*) AS total_addresses FROM addresses;
SELECT COUNT(*) AS total_orders FROM orders;
SELECT COUNT(*) AS total_order_items FROM order_items;
SELECT COUNT(*) AS total_notifications FROM notifications;
