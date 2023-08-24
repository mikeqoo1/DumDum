/*!40014 SET FOREIGN_KEY_CHECKS=0*/;
/*!40101 SET NAMES binary*/;
CREATE TABLE `orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用戶資訊',
  `order_date` date NOT NULL COMMENT '訂單日期',
  `payment_status` enum('Pending','Paid','Cancelled') NOT NULL COMMENT '付款狀態',
  `shipping_status` enum('Pending','Shipped','Delivered') NOT NULL COMMENT '配送狀態',
  `total_amount` decimal(10,2) NOT NULL COMMENT '訂單總金額',
  PRIMARY KEY (`id`) /*T![clustered_index] CLUSTERED */,
  KEY `fk_1` (`user_id`),
  CONSTRAINT `fk_1` FOREIGN KEY (`user_id`) REFERENCES `sea`.`users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='訂單資料表';
