/*!40014 SET FOREIGN_KEY_CHECKS=0*/;
/*!40101 SET NAMES binary*/;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用戶名稱',
  `email` varchar(100) NOT NULL COMMENT '電子郵件',
  `password` varchar(100) NOT NULL COMMENT '密碼',
  `address` text DEFAULT NULL COMMENT '配送地址',
  `payment_info` varchar(200) DEFAULT NULL COMMENT '付款資訊',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=30001 COMMENT='用戶資料表';
