-- Seed data for bycigar shop
-- Run: docker exec -i bycigar-mysql mysql -u root -proot123 bycigar < seed_products.sql

-- Clear existing data
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE products;
TRUNCATE TABLE categories;
SET FOREIGN_KEY_CHECKS = 1;

-- Insert categories
INSERT INTO categories (id, name, slug, parent_id, created_at, updated_at) VALUES
(1, 'Cuban Cigars', 'cuban', NULL, NOW(), NOW()),
(2, 'Nicaraguan Cigars', 'nicaraguan', NULL, NOW(), NOW()),
(3, 'Dominican Cigars', 'dominican', NULL, NOW(), NOW()),
(4, 'Accessories', 'accessories', NULL, NOW(), NOW()),
(5, 'Humidors', 'humidors', NULL, NOW(), NOW()),
(6, 'Gift Sets', 'gift-sets', NULL, NOW(), NOW());

-- Insert products (55 items)
INSERT INTO products (name, slug, description, price, image, images, category_id, stock, is_active, is_featured, created_at, updated_at) VALUES
-- Cuban Cigars (10)
('Cohiba Siglo I', 'cohiba-siglo-i', 'A medium-bodied cigar with smooth, creamy flavors and a hint of spice. Perfect for a quick smoke.', 28.00, 'https://picsum.photos/seed/cigar1/400/300', 'https://picsum.photos/seed/cigar1a/800/600,https://picsum.photos/seed/cigar1b/800/600', 1, 85, true, true, NOW(), NOW()),
('Cohiba Siglo II', 'cohiba-siglo-ii', 'Rich and aromatic with notes of cedar and leather. A classic Cuban experience.', 32.00, 'https://picsum.photos/seed/cigar2/400/300', 'https://picsum.photos/seed/cigar2a/800/600', 1, 72, true, false, NOW(), NOW()),
('Cohiba Siglo III', 'cohiba-siglo-iii', 'Complex flavors of coffee, cocoa, and subtle spices. Excellent construction.', 38.00, 'https://picsum.photos/seed/cigar3/400/300', 'https://picsum.photos/seed/cigar3a/800/600', 1, 60, true, true, NOW(), NOW()),
('Montecristo No.2', 'montecristo-no2', 'The legendary torpedo with rich, earthy flavors and a smooth finish.', 26.00, 'https://picsum.photos/seed/cigar4/400/300', 'https://picsum.photos/seed/cigar4a/800/600', 1, 120, true, true, NOW(), NOW()),
('Montecristo No.4', 'montecristo-no4', 'A bestseller with balanced flavors of wood, leather, and subtle sweetness.', 22.00, 'https://picsum.photos/seed/cigar5/400/300', 'https://picsum.photos/seed/cigar5a/800/600', 1, 150, true, false, NOW(), NOW()),
('Romeo y Julieta Churchill', 'romeo-julieta-churchill', 'Classic Churchill size with elegant, medium-bodied flavors.', 30.00, 'https://picsum.photos/seed/cigar6/400/300', 'https://picsum.photos/seed/cigar6a/800/600', 1, 45, true, false, NOW(), NOW()),
('Partagas Serie D No.4', 'partagas-serie-d-no4', 'Full-bodied with rich, earthy notes and a hint of pepper.', 25.00, 'https://picsum.photos/seed/cigar7/400/300', 'https://picsum.photos/seed/cigar7a/800/600', 1, 88, true, true, NOW(), NOW()),
('H. Upmann No.2', 'h-upmann-no2', 'Smooth and mellow with notes of honey and toasted nuts.', 24.00, 'https://picsum.photos/seed/cigar8/400/300', 'https://picsum.photos/seed/cigar8a/800/600', 1, 55, true, false, NOW(), NOW()),
('Bolivar Royal Corona', 'bolivar-royal-corona', 'Powerful and robust with intense flavors of leather and spice.', 28.00, 'https://picsum.photos/seed/cigar9/400/300', 'https://picsum.photos/seed/cigar9a/800/600', 1, 40, true, false, NOW(), NOW()),
('Punch Punch', 'punch-punch', 'Medium to full-bodied with woody and nutty characteristics.', 20.00, 'https://picsum.photos/seed/cigar10/400/300', 'https://picsum.photos/seed/cigar10a/800/600', 1, 95, true, false, NOW(), NOW()),

-- Nicaraguan Cigars (9)
('Padron 1964 Anniversary', 'padron-1964-anniversary', 'Award-winning cigar with rich, complex flavors of coffee and cocoa.', 35.00, 'https://picsum.photos/seed/cigar11/400/300', 'https://picsum.photos/seed/cigar11a/800/600', 2, 65, true, true, NOW(), NOW()),
('Padron 1926 Serie', 'padron-1926-serie', 'Luxurious smoke with notes of dark chocolate and espresso.', 42.00, 'https://picsum.photos/seed/cigar12/400/300', 'https://picsum.photos/seed/cigar12a/800/600', 2, 38, true, false, NOW(), NOW()),
('Oliva Serie V', 'oliva-serie-v', 'Full-bodied with rich, spicy flavors and a long, smooth finish.', 18.00, 'https://picsum.photos/seed/cigar13/400/300', 'https://picsum.photos/seed/cigar13a/800/600', 2, 110, true, true, NOW(), NOW()),
('My Father Le Bijou', 'my-father-le-bijou', 'Complex and refined with notes of black pepper and cedar.', 22.00, 'https://picsum.photos/seed/cigar14/400/300', 'https://picsum.photos/seed/cigar14a/800/600', 2, 78, true, false, NOW(), NOW()),
('Tatuaje Fausto', 'tatuaje-fausto', 'Bold and peppery with a rich, earthy core.', 16.00, 'https://picsum.photos/seed/cigar15/400/300', 'https://picsum.photos/seed/cigar15a/800/600', 2, 92, true, false, NOW(), NOW()),
('Drew Estate Liga Privada', 'drew-estate-liga-privada', 'Unique blend with notes of coffee, cocoa, and subtle sweetness.', 28.00, 'https://picsum.photos/seed/cigar16/400/300', 'https://picsum.photos/seed/cigar16a/800/600', 2, 52, true, true, NOW(), NOW()),
('AJ Fernandez New World', 'aj-fernandez-new-world', 'Rich and complex with flavors of leather, wood, and spice.', 14.00, 'https://picsum.photos/seed/cigar17/400/300', 'https://picsum.photos/seed/cigar17a/800/600', 2, 130, true, false, NOW(), NOW()),
('Puronicago', 'puronicago', 'Exceptional craftsmanship with balanced, nuanced flavors.', 25.00, 'https://picsum.photos/seed/cigar18/400/300', 'https://picsum.photos/seed/cigar18a/800/600', 2, 45, true, false, NOW(), NOW()),
('Plasencia Alma Fuerte', 'plasencia-alma-fuerte', 'Bold and complex with notes of dark fruit and spices.', 32.00, 'https://picsum.photos/seed/cigar19/400/300', 'https://picsum.photos/seed/cigar19a/800/600', 2, 58, true, false, NOW(), NOW()),

-- Dominican Cigars (9)
('Arturo Fuente OpusX', 'arturo-fuente-opusx', 'Legendary cigar with rich, full-bodied flavors and perfect construction.', 65.00, 'https://picsum.photos/seed/cigar20/400/300', 'https://picsum.photos/seed/cigar20a/800/600', 3, 25, true, true, NOW(), NOW()),
('Arturo Fuente Hemingway', 'arturo-fuente-hemingway', 'Elegant and smooth with notes of cedar and cream.', 22.00, 'https://picsum.photos/seed/cigar21/400/300', 'https://picsum.photos/seed/cigar21a/800/600', 3, 88, true, false, NOW(), NOW()),
('Davidoff Millennium', 'davidoff-millennium', 'Sophisticated blend with rich, complex flavors.', 48.00, 'https://picsum.photos/seed/cigar22/400/300', 'https://picsum.photos/seed/cigar22a/800/600', 3, 35, true, true, NOW(), NOW()),
('Davidoff Winston Churchill', 'davidoff-winston-churchill', 'Inspired by the iconic leader, refined and elegant.', 52.00, 'https://picsum.photos/seed/cigar23/400/300', 'https://picsum.photos/seed/cigar23a/800/600', 3, 42, true, false, NOW(), NOW()),
('Avo Classic', 'avo-classic', 'Smooth and mellow with a perfect balance of flavors.', 18.00, 'https://picsum.photos/seed/cigar24/400/300', 'https://picsum.photos/seed/cigar24a/800/600', 3, 95, true, false, NOW(), NOW()),
('Ashton Classic', 'ashton-classic', 'Mild and creamy with excellent construction.', 16.00, 'https://picsum.photos/seed/cigar25/400/300', 'https://picsum.photos/seed/cigar25a/800/600', 3, 105, true, false, NOW(), NOW()),
('Ashton VSG', 'ashton-vsg', 'Full-bodied with rich, complex flavors of earth and spice.', 24.00, 'https://picsum.photos/seed/cigar26/400/300', 'https://picsum.photos/seed/cigar26a/800/600', 3, 68, true, true, NOW(), NOW()),
('La Aurora 1903', 'la-aurora-1903', 'A century of tradition in every smoke. Rich and balanced.', 20.00, 'https://picsum.photos/seed/cigar27/400/300', 'https://picsum.photos/seed/cigar27a/800/600', 3, 75, true, false, NOW(), NOW()),
('Fuente Don Carlos', 'fuente-don-carlos', 'Premium blend with exceptional smoothness and depth.', 38.00, 'https://picsum.photos/seed/cigar28/400/300', 'https://picsum.photos/seed/cigar28a/800/600', 3, 48, true, false, NOW(), NOW()),

-- Accessories (10)
('Xikar Xi2 Cutter', 'xikar-xi2-cutter', 'Premium double-blade cutter for a perfect cut every time.', 45.00, 'https://picsum.photos/seed/acc1/400/300', 'https://picsum.photos/seed/acc1a/800/600', 4, 150, true, true, NOW(), NOW()),
('Colibri V-Cut Cutter', 'colibri-v-cut-cutter', 'Precision V-cut for enhanced flavor and draw.', 35.00, 'https://picsum.photos/seed/acc2/400/300', 'https://picsum.photos/seed/acc2a/800/600', 4, 120, true, false, NOW(), NOW()),
('S.T. Dupont Lighter', 'st-dupont-lighter', 'Luxury butane lighter with iconic "cling" sound.', 180.00, 'https://picsum.photos/seed/acc3/400/300', 'https://picsum.photos/seed/acc3a/800/600', 4, 45, true, true, NOW(), NOW()),
('Xikar HP Lighter', 'xikar-hp-lighter', 'Reliable soft flame lighter with elegant design.', 55.00, 'https://picsum.photos/seed/acc4/400/300', 'https://picsum.photos/seed/acc4a/800/600', 4, 85, true, false, NOW(), NOW()),
('Cigar Travel Case', 'cigar-travel-case', 'Leather travel case holds 5 cigars. Perfect for on-the-go.', 65.00, 'https://picsum.photos/seed/acc5/400/300', 'https://picsum.photos/seed/acc5a/800/600', 4, 70, true, false, NOW(), NOW()),
('Cedar Spill Holder', 'cedar-spill-holder', 'Traditional cedar spills for lighting cigars the proper way.', 28.00, 'https://picsum.photos/seed/acc6/400/300', 'https://picsum.photos/seed/acc6a/800/600', 4, 95, true, false, NOW(), NOW()),
('Cigar Ashtray Crystal', 'cigar-ashtray-crystal', 'Elegant crystal ashtray with 4 rests. A statement piece.', 120.00, 'https://picsum.photos/seed/acc7/400/300', 'https://picsum.photos/seed/acc7a/800/600', 4, 30, true, false, NOW(), NOW()),
('Leather Cigar Wallet', 'leather-cigar-wallet', 'Genuine leather wallet holds 3 cigars. Sleek and portable.', 48.00, 'https://picsum.photos/seed/acc8/400/300', 'https://picsum.photos/seed/acc8a/800/600', 4, 60, true, false, NOW(), NOW()),
('Cigar Punch Keyring', 'cigar-punch-keyring', 'Compact punch cutter on a keyring. Always ready.', 22.00, 'https://picsum.photos/seed/acc9/400/300', 'https://picsum.photos/seed/acc9a/800/600', 4, 180, true, false, NOW(), NOW()),
('Boveda Humidity Pack', 'boveda-humidity-pack', '62% humidity packs for optimal cigar storage. Pack of 10.', 25.00, 'https://picsum.photos/seed/acc10/400/300', 'https://picsum.photos/seed/acc10a/800/600', 4, 200, true, false, NOW(), NOW()),

-- Humidors (9)
('Desktop Humidor 25', 'desktop-humidor-25', 'Spanish cedar lined humidor holds 25 cigars. Classic design.', 85.00, 'https://picsum.photos/seed/hum1/400/300', 'https://picsum.photos/seed/hum1a/800/600', 5, 40, true, true, NOW(), NOW()),
('Desktop Humidor 50', 'desktop-humidor-50', 'Larger capacity humidor with digital hygrometer.', 145.00, 'https://picsum.photos/seed/hum2/400/300', 'https://picsum.photos/seed/hum2a/800/600', 5, 35, true, false, NOW(), NOW()),
('Glass Top Humidor', 'glass-top-humidor', 'Elegant glass top design. Holds 75 cigars.', 195.00, 'https://picsum.photos/seed/hum3/400/300', 'https://picsum.photos/seed/hum3a/800/600', 5, 22, true, true, NOW(), NOW()),
('Cabinet Humidor 200', 'cabinet-humidor-200', 'Premium cabinet humidor. Holds up to 200 cigars.', 450.00, 'https://picsum.photos/seed/hum4/400/300', 'https://picsum.photos/seed/hum4a/800/600', 5, 12, true, false, NOW(), NOW()),
('Travel Humidor', 'travel-humidor', 'Durable travel case holds 10 cigars. Waterproof.', 55.00, 'https://picsum.photos/seed/hum5/400/300', 'https://picsum.photos/seed/hum5a/800/600', 5, 75, true, false, NOW(), NOW()),
('Leather Humidor', 'leather-humidor', 'Genuine leather exterior with cedar interior. Holds 20.', 165.00, 'https://picsum.photos/seed/hum6/400/300', 'https://picsum.photos/seed/hum6a/800/600', 5, 28, true, false, NOW(), NOW()),
('Electric Humidor', 'electric-humidor', 'Temperature and humidity controlled. Holds 100 cigars.', 380.00, 'https://picsum.photos/seed/hum7/400/300', 'https://picsum.photos/seed/hum7a/800/600', 5, 18, true, true, NOW(), NOW()),
('Countertop Humidor', 'countertop-humidor', 'Display humidor for retail or home. Holds 150 cigars.', 320.00, 'https://picsum.photos/seed/hum8/400/300', 'https://picsum.photos/seed/hum8a/800/600', 5, 15, true, false, NOW(), NOW()),
('Mini Humidor 5', 'mini-humidor-5', 'Compact personal humidor. Perfect for beginners.', 45.00, 'https://picsum.photos/seed/hum9/400/300', 'https://picsum.photos/seed/hum9a/800/600', 5, 90, true, false, NOW(), NOW()),

-- Gift Sets (8)
('Cuban Sampler', 'cuban-sampler', '5 premium Cuban cigars in a gift box. Perfect introduction.', 125.00, 'https://picsum.photos/seed/gift1/400/300', 'https://picsum.photos/seed/gift1a/800/600', 6, 55, true, true, NOW(), NOW()),
('World Tour Sampler', 'world-tour-sampler', '8 cigars from Cuba, Nicaragua, and Dominican Republic.', 145.00, 'https://picsum.photos/seed/gift2/400/300', 'https://picsum.photos/seed/gift2a/800/600', 6, 42, true, true, NOW(), NOW()),
('Premium Starter Kit', 'premium-starter-kit', 'Includes humidor, cutter, lighter, and 10 cigars.', 220.00, 'https://picsum.photos/seed/gift3/400/300', 'https://picsum.photos/seed/gift3a/800/600', 6, 30, true, false, NOW(), NOW()),
('Executive Gift Set', 'executive-gift-set', 'Luxury set with Dupont lighter, Xikar cutter, and 5 cigars.', 350.00, 'https://picsum.photos/seed/gift4/400/300', 'https://picsum.photos/seed/gift4a/800/600', 6, 20, true, false, NOW(), NOW()),
('Connoisseur Collection', 'connoisseur-collection', '12 rare and aged cigars in a premium humidor.', 480.00, 'https://picsum.photos/seed/gift5/400/300', 'https://picsum.photos/seed/gift5a/800/600', 6, 15, true, true, NOW(), NOW()),
('Birthday Box', 'birthday-box', 'Customizable gift box with 6 cigars and accessories.', 95.00, 'https://picsum.photos/seed/gift6/400/300', 'https://picsum.photos/seed/gift6a/800/600', 6, 65, true, false, NOW(), NOW()),
('Anniversary Set', 'anniversary-set', 'Premium cigars paired with aged rum. A special gift.', 185.00, 'https://picsum.photos/seed/gift7/400/300', 'https://picsum.photos/seed/gift7a/800/600', 6, 38, true, false, NOW(), NOW()),
('Accessory Gift Box', 'accessory-gift-box', 'Complete accessory set: cutter, lighter, case, ashtray.', 165.00, 'https://picsum.photos/seed/gift8/400/300', 'https://picsum.photos/seed/gift8a/800/600', 6, 48, true, false, NOW(), NOW());

SELECT 'Seed completed!' AS status;
SELECT COUNT(*) AS total_products FROM products;
SELECT c.name, COUNT(p.id) AS product_count FROM categories c LEFT JOIN products p ON c.id = p.category_id GROUP BY c.id;
