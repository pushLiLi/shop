-- BYCIGAR 演示数据
-- 说明: 用户、地址、订单、聊天记录、通知等演示数据
-- 用途: 开发测试、展示演示
-- 警告: 此脚本会创建大量数据，首次运行请确保数据库已初始化

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 演示用户 (25个客户账号)
-- 密码: 123456
-- ----------------------------
INSERT INTO `users` (`id`, `email`, `password`, `name`, `role`, `created_at`, `updated_at`) VALUES
(4, 'user01@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '张伟', 'customer', DATE_SUB(NOW(), INTERVAL 90 DAY), NOW()),
(5, 'user02@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '李娜', 'customer', DATE_SUB(NOW(), INTERVAL 85 DAY), NOW()),
(6, 'user03@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '王芳', 'customer', DATE_SUB(NOW(), INTERVAL 80 DAY), NOW()),
(7, 'user04@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '刘洋', 'customer', DATE_SUB(NOW(), INTERVAL 75 DAY), NOW()),
(8, 'user05@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '陈杰', 'customer', DATE_SUB(NOW(), INTERVAL 70 DAY), NOW()),
(9, 'user06@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '杨秀英', 'customer', DATE_SUB(NOW(), INTERVAL 65 DAY), NOW()),
(10, 'user07@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '黄明', 'customer', DATE_SUB(NOW(), INTERVAL 60 DAY), NOW()),
(11, 'user08@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '赵磊', 'customer', DATE_SUB(NOW(), INTERVAL 55 DAY), NOW()),
(12, 'user09@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '周静', 'customer', DATE_SUB(NOW(), INTERVAL 50 DAY), NOW()),
(13, 'user10@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '吴强', 'customer', DATE_SUB(NOW(), INTERVAL 45 DAY), NOW()),
(14, 'user11@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '徐丽', 'customer', DATE_SUB(NOW(), INTERVAL 40 DAY), NOW()),
(15, 'user12@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '孙浩', 'customer', DATE_SUB(NOW(), INTERVAL 35 DAY), NOW()),
(16, 'user13@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '马超', 'customer', DATE_SUB(NOW(), INTERVAL 30 DAY), NOW()),
(17, 'user14@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '胡蝶', 'customer', DATE_SUB(NOW(), INTERVAL 25 DAY), NOW()),
(18, 'user15@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '朱琳', 'customer', DATE_SUB(NOW(), INTERVAL 20 DAY), NOW()),
(19, 'user16@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '郭锐', 'customer', DATE_SUB(NOW(), INTERVAL 18 DAY), NOW()),
(20, 'user17@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '林涛', 'customer', DATE_SUB(NOW(), INTERVAL 15 DAY), NOW()),
(21, 'user18@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '何雪', 'customer', DATE_SUB(NOW(), INTERVAL 12 DAY), NOW()),
(22, 'user19@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '高建', 'customer', DATE_SUB(NOW(), INTERVAL 10 DAY), NOW()),
(23, 'user20@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '罗岚', 'customer', DATE_SUB(NOW(), INTERVAL 8 DAY), NOW()),
(24, 'user21@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '梁志明', 'customer', DATE_SUB(NOW(), INTERVAL 6 DAY), NOW()),
(25, 'user22@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '宋雨晴', 'customer', DATE_SUB(NOW(), INTERVAL 5 DAY), NOW()),
(26, 'user23@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '郑浩', 'customer', DATE_SUB(NOW(), INTERVAL 4 DAY), NOW()),
(27, 'user24@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '谢辉', 'customer', DATE_SUB(NOW(), INTERVAL 3 DAY), NOW()),
(28, 'user25@test.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '田亮', 'customer', DATE_SUB(NOW(), INTERVAL 2 DAY), NOW())
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- ----------------------------
-- 收货地址
-- ----------------------------
INSERT INTO `addresses` (`user_id`, `full_name`, `address_line1`, `city`, `state`, `zip_code`, `phone`, `is_default`, `created_at`, `updated_at`) VALUES
-- user01 (id=4)
(4, '张伟', '建国路88号SOHO现代城', '北京', '北京', '100022', '13800138001', 1, NOW(), NOW()),
(4, '张伟', '建国路88号SOHO现代城B座', '北京', '北京', '100022', '13800138001', 0, NOW(), NOW()),
-- user02 (id=5)
(5, '李娜', '南京东路100号', '上海', '上海', '200001', '13900139002', 1, NOW(), NOW()),
-- user03 (id=6)
(6, '王芳', '天河路385号太古汇', '广州', '广东', '510620', '13700137003', 1, NOW(), NOW()),
(6, '王芳', '天河路385号太古汇（仓库）', '广州', '广东', '510620', '13700137003', 0, NOW(), NOW()),
-- user04 (id=7)
(7, '刘洋', '科技园南区深南大道9966号', '深圳', '广东', '518057', '13600136004', 1, NOW(), NOW()),
-- user05 (id=8)
(8, '陈杰', '文三路553号', '杭州', '浙江', '310012', '13500135005', 1, NOW(), NOW()),
-- user06 (id=9)
(9, '杨秀英', '红星路三段1号', '成都', '四川', '610021', '13400134006', 1, NOW(), NOW()),
(9, '杨秀英', '红星路三段1号B座', '成都', '四川', '610021', '13400134006', 0, NOW(), NOW()),
-- user07 (id=10)
(10, '黄明', '中北路109号', '武汉', '湖北', '430071', '13300133007', 1, NOW(), NOW()),
-- user08 (id=11)
(11, '赵磊', '南京路189号', '天津', '天津', '300051', '13200132008', 1, NOW(), NOW()),
-- user09 (id=12)
(12, '周静', '中央路331号', '南京', '江苏', '210008', '13100131009', 1, NOW(), NOW()),
-- user10 (id=13)
(13, '吴强', '麓谷大道627号', '长沙', '湖南', '410205', '13000130010', 1, NOW(), NOW()),
(13, '吴强', '麓谷大道627号（仓库）', '长沙', '湖南', '410205', '13000130010', 0, NOW(), NOW()),
-- user11 (id=14)
(14, '徐丽', '长安南路86号', '西安', '陕西', '710061', '15000150011', 1, NOW(), NOW()),
-- user12 (id=15)
(15, '孙浩', '厦禾路888号', '厦门', '福建', '361003', '15100151012', 1, NOW(), NOW()),
-- user13 (id=16)
(16, '马超', '解放碑步行街88号', '重庆', '重庆', '400010', '15200152013', 1, NOW(), NOW()),
-- user14 (id=17)
(17, '胡蝶', '观前街168号', '苏州', '江苏', '215000', '15300153014', 1, NOW(), NOW()),
-- user15 (id=18)
(18, '朱琳', '二七路72号', '郑州', '河南', '450000', '15500155015', 1, NOW(), NOW()),
-- user16 (id=19)
(19, '郭锐', '香港中路100号', '青岛', '山东', '266000', '15600156016', 1, NOW(), NOW()),
-- user17 (id=20)
(20, '林涛', '五四路200号', '沈阳', '辽宁', '110001', '15700157017', 1, NOW(), NOW()),
-- user18 (id=21)
(21, '何雪', '长江道1号', '大连', '辽宁', '116001', '15800158018', 1, NOW(), NOW()),
-- user19 (id=22)
(22, '高建', '东风路388号', '济南', '山东', '250001', '15900159019', 1, NOW(), NOW()),
-- user20 (id=23)
(23, '罗岚', '五一路399号', '福州', '福建', '350001', '18000180020', 1, NOW(), NOW()),
-- user21 (id=24)
(24, '梁志明', '建设路188号', '昆明', '云南', '650000', '18100181021', 1, NOW(), NOW()),
-- user22 (id=25)
(25, '宋雨晴', '中山路299号', '哈尔滨', '黑龙江', '150001', '18200182022', 1, NOW(), NOW()),
-- user23 (id=26)
(26, '郑浩', '胜利北路66号', '长春', '吉林', '130000', '18300183023', 1, NOW(), NOW()),
-- user24 (id=27)
(27, '谢辉', '人民路555号', '石家庄', '河北', '050000', '18400184024', 1, NOW(), NOW()),
-- user25 (id=28)
(28, '田亮', '滨海大道88号中心大厦', '南昌', '江西', '330000', '18500185025', 1, NOW(), NOW());

-- ----------------------------
-- 快捷回复
-- ----------------------------
INSERT INTO `quick_replies` (`title`, `content`, `created_by`, `sort_order`, `created_at`, `updated_at`) VALUES
('欢迎语', '您好！欢迎来到 HUAUHE，有什么可以帮助您的吗？', 2, 1, NOW(), NOW()),
('产品咨询', '感谢您的咨询，我们的产品均为正品，如有任何问题随时联系客服。', 2, 2, NOW(), NOW()),
('发货说明', '您的订单已发货，快递正在配送中，请注意查收。', 2, 3, NOW(), NOW()),
('感谢语', '感谢您的支持，祝您生活愉快！', 2, 4, NOW(), NOW()),
('缺货处理', '非常抱歉，您咨询的商品目前缺货中，可以留下您的联系方式，到货后我们会第一时间通知您。', 2, 5, NOW(), NOW()),
('退换货说明', '雪茄属于特殊商品，若因质量问题（如发霉、干燥等）可在收货48小时内申请退换。', 3, 1, NOW(), NOW()),
('支付帮助', '我们支持微信支付、支付宝、银行转账、PayPal以及货到付款等多种支付方式。', 3, 2, NOW(), NOW()),
('物流时效', '顺丰快递一般1-3天送达，偏远地区3-5天。', 3, 3, NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
