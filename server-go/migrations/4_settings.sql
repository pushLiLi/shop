-- BYCIGAR 配置数据
-- 说明: 网站设置、页面内容、SEO配置等

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 网站设置
-- ----------------------------
INSERT INTO `settings` (`key`, `value`, `created_at`, `updated_at`) VALUES
('site_name', 'BYCIGAR 雪茄旗舰店', NOW(), NOW()),
('chat_greeting', '您好！欢迎来到 HUAUHE，有什么可以帮助您的吗？', NOW(), NOW()),
('contact_phone', '400-888-9999', NOW(), NOW()),
('contact_email', 'support@bycigar.com', NOW(), NOW()),
('footer_content', '© 2024 BYCIGAR 版权所有', NOW(), NOW()),
('footer_description', 'BYCIGAR是中国领先的雪茄文化与在线购物平台。我们提供最新、最专业的雪茄测评、品牌新闻与养护知识，并为您甄选全球优质雪茄及配件，支持便捷在线购买。加入我们的雪茄社区，探索醇香世界。', NOW(), NOW()),
('footer_service_time', '客服在线时间每周一至周六 9:00到18:00', NOW(), NOW()),
('order_auto_close_days', '7', NOW(), NOW()),
('low_stock_threshold', '10', NOW(), NOW()),
('maintenance_mode', 'false', NOW(), NOW())
ON DUPLICATE KEY UPDATE `value` = VALUES(`value`);

-- ----------------------------
-- 网站配置 (SiteConfig)
-- ----------------------------
INSERT INTO `site_configs` (`config_key`, `config_value`, `created_at`, `updated_at`) VALUES
('site_title', 'BYCIGAR 雪茄商城', NOW(), NOW()),
('seo_description', '专业古巴雪茄及高档雪茄配件销售平台，提供高希霸、蒙特、大卫杜夫等品牌雪茄，品类齐全，品质保证。', NOW(), NOW()),
('customer_service_hours', '9:00-22:00', NOW(), NOW()),
('max_cart_items', '99', NOW(), NOW()),
('order_max_quantity', '50', NOW(), NOW()),
('exchange_rate_usd', '7.25', NOW(), NOW()),
('home_banner_1', '/media/bycigar/banner-1.png', NOW(), NOW()),
('home_banner_2', '/media/bycigar/banner-2.png', NOW(), NOW()),
('home_banner_3', '/media/bycigar/banner-3.png', NOW(), NOW()),
('home_promo_left_image', '', NOW(), NOW()),
('home_promo_left_link', '', NOW(), NOW()),
('home_promo_right_image', '', NOW(), NOW()),
('home_promo_right_link', '', NOW(), NOW()),
('home_featured_title', '特别推荐', NOW(), NOW()),
('home_new_title', '新品上架', NOW(), NOW()),
('home_topselling_title', '热销排行', NOW(), NOW())
ON DUPLICATE KEY UPDATE `config_value` = VALUES(`config_value`);

-- ----------------------------
-- 页面内容
-- ----------------------------
INSERT INTO `pages` (`slug`, `title`, `content`, `updated_at`) VALUES
('about', '关于我们', CONCAT('# 关于我们\n\nBYCIGAR 成立于2010年，是国内领先的高端雪茄零售商。我们致力于为雪茄爱好者提供最优质的古巴及世界顶级雪茄，产品涵盖高希霸、蒙特、帕塔加斯、大卫杜夫等经典品牌。公司总部位于上海，拥有专业化的恒温仓储物流体系，所有雪茄均经过严格品控，确保品质如一。\n\n我们的团队由一群热爱雪茄文化的专业人士组成，每位顾问都经过严格培训，能够为客户提供专业的品鉴建议和搭配推荐。BYCIGAR 还定期举办雪茄品鉴会，邀请国内外知名雪茄大师与会员互动交流。\n\n未来，我们将继续深耕雪茄文化，引入更多优质产品，为中国雪茄市场的发展贡献力量。'), NOW()),

('services', '服务条款', CONCAT('# 服务条款\n\n一、服务说明\n\nBYCIGAR 为您提供在线雪茄及配件销售服务，包括商品浏览、在线购买、支付结算、物流配送等完整流程。\n\n二、购买须知\n\n1. 雪茄产品需年满18周岁方可购买。\n2. 请如实填写收货信息，确保商品顺利送达。\n3. 雪茄属于特殊商品一经拆封不支持无理由退换。\n4. 如收到商品有质量问题，请在48小时内联系客服。\n\n三、物流配送\n\n我们合作的物流伙伴包括顺丰、京东等优质快递，全程冷链运输，确保雪茄在最佳状态下送达。部分地区可能需要身份证验证。\n\n四、售后服务\n\n客服热线：400-888-9999（工作日9:00-22:00）\n邮箱：support@bycigar.com'), NOW()),

('privacy-policy', '隐私政策', CONCAT('# 隐私政策\n\nBYCIGAR 非常重视您的个人信息的保护。本隐私政策说明了我们在您使用我们的服务时如何收集、使用、存储和保护您的个人信息。\n\n一、信息收集\n\n当您注册账户时，我们会收集您的姓名、邮箱、手机号码等基本信息。当您下单时，我们会收集收货地址信息以完成配送。\n\n二、信息使用\n\n您的个人信息将用于：处理订单、提供客户服务、发送订单状态通知、推荐适合的商品、改进我们的服务。\n\n三、信息保护\n\n我们采用行业标准的加密技术保护您的数据安全，未经您的授权，我们不会将个人信息提供给任何第三方。\n\n四、联系我们\n\n如对隐私政策有任何疑问，请联系：privacy@bycigar.com'), NOW()),

('statement', '免责说明', CONCAT('# 免责说明\n\n一、产品声明\n\nBYCIGAR 销售的所有雪茄产品均为正品。我们合作的品牌均经过官方授权，或从正规渠道采购。\n\n二、健康提示\n\n吸烟有害健康。雪茄含有焦油和一氧化碳等有害物质，未成年人、孕妇、哺乳期妇女及心血管疾病患者不应使用雪茄。\n\n三、使用风险\n\n雪茄的保存需要特定温度（18-21°C）和湿度（68-72%）。因客户保存不当导致的雪茄品质问题，不在退换货范围内。\n\n四、配送风险\n\n商品在运输过程中可能因不可抗力因素导致延误或损坏，我们将配合物流公司积极处理，但不承担额外赔偿责任。'), NOW()),

('shipping', '配送说明', CONCAT('# 配送说明\n\n一、配送范围\n\n我们支持全国大部分地区的配送，包括港澳台地区。部分偏远地区可能需要更长的配送时间。\n\n二、配送方式\n\n1. 标准配送：3-5个工作日送达\n2. 加急配送：1-2个工作日送达（需额外支付费用）\n3. 定时配送：您可选择具体的收货时间段\n\n三、配送费用\n\n订单满500元免运费，不满500元收取20元运费。新会员首单免运费。\n\n四、配送包装\n\n所有雪茄均采用专业保湿包装，配合保温箱确保品质。礼盒订单额外使用精美纸箱包装。\n\n五、签收提示\n\n请在签收前检查包装是否完好，如有问题请当场拒收并联系客服。'), NOW())
ON DUPLICATE KEY UPDATE `title` = VALUES(`title`), `content` = VALUES(`content`);

SET FOREIGN_KEY_CHECKS = 1;
