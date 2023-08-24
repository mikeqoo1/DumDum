/*!40014 SET FOREIGN_KEY_CHECKS=0*/;
/*!40101 SET NAMES binary*/;
CREATE TABLE `products` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '產品名稱',
  `description` text DEFAULT NULL COMMENT '描述',
  `price` decimal(10,2) NOT NULL COMMENT '價格',
  `stock` int(11) NOT NULL COMMENT '庫存',
  `sku` varchar(50) NOT NULL COMMENT 'SKU(庫存單位)',
  `image_url` varchar(200) DEFAULT NULL COMMENT '圖片',
  `category` varchar(100) NOT NULL COMMENT '產品分類',
  `enabled` int(11) DEFAULT NULL,
  `discount` int(11) NOT NULL DEFAULT '0' COMMENT '折扣',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=30001 COMMENT='產品資料表';
