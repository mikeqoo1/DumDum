/*!40014 SET FOREIGN_KEY_CHECKS=0*/;
/*!40101 SET NAMES binary*/;
CREATE TABLE `nici` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '流水編號',
  `name` varchar(50) NOT NULL COMMENT '名稱',
  `blood` varchar(10) NOT NULL COMMENT '血型',
  `starsign` varchar(20) NOT NULL COMMENT '星座',
  `series` varchar(50) NOT NULL COMMENT '系列',
  `img` varchar(50) NOT NULL COMMENT '圖檔名稱',
  PRIMARY KEY (`id`,`name`) /*T![clustered_index] CLUSTERED */
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin AUTO_INCREMENT=30001 COMMENT='Nici的身份詳細表';
