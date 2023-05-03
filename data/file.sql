CREATE TABLE IF NOT EXISTS `nici` (
    id int(11) AUTO_INCREMENT NOT NULL COMMENT '流水編號',
    name varchar(50) NOT NULL COMMENT '名稱',
    blood varchar(10) NOT NULL COMMENT '血型',
    starsign varchar(20) NOT NULL COMMENT '星座',
    series varchar(50) NOT NULL COMMENT '系列',
    img varchar(50) NOT NULL COMMENT '圖檔名稱',
    PRIMARY KEY(`id`, `name`)
) COMMENT='Nici的身份詳細表';

