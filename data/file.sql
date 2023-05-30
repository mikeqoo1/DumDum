CREATE TABLE IF NOT EXISTS `nici` (
    id int(11) AUTO_INCREMENT NOT NULL COMMENT '流水編號',
    name varchar(50) NOT NULL COMMENT '名稱',
    blood varchar(10) NOT NULL COMMENT '血型',
    starsign varchar(20) NOT NULL COMMENT '星座',
    series varchar(50) NOT NULL COMMENT '系列',
    img varchar(50) NOT NULL COMMENT '圖檔名稱',
    PRIMARY KEY(`id`, `name`)
) COMMENT='Nici的身份詳細表';

CREATE TABLE IF NOT EXISTS `userlist` (
    id int(11) AUTO_INCREMENT NOT NULL COMMENT '流水編號',
    account varchar(50) NOT NULL COMMENT '使用者帳號',
    username varchar(50) NOT NULL COMMENT '使用者名稱',
    status int(2) NOT NULL DEFAULT 1 COMMENT '狀態 0:關閉 1:啟用',
    PRIMARY KEY(`id`, `account`)
) COMMENT='使用者清單';

