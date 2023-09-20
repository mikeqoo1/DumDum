CREATE TABLE IF NOT EXISTS `nici` (
    id int(11) AUTO_INCREMENT NOT NULL COMMENT '流水編號',
    name varchar(50) NOT NULL COMMENT '名稱',
    blood varchar(10) NOT NULL COMMENT '血型',
    starsign varchar(20) NOT NULL COMMENT '星座',
    series varchar(50) NOT NULL COMMENT '系列',
    img varchar(50) NOT NULL COMMENT '圖檔名稱',
    PRIMARY KEY(`id`, `name`)
) COMMENT='Nici的身份詳細表';

-- 建立用戶資料表
CREATE TABLE IF NOT EXISTS `users` (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL COMMENT '用戶名稱',
  email VARCHAR(100) NOT NULL COMMENT '電子郵件',
  password VARCHAR(100) NOT NULL COMMENT '密碼',
  address TEXT COMMENT '配送地址',
  payment_info VARCHAR(200) COMMENT '付款資訊',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) COMMENT='用戶資料表';

-- 建立訂單資料表
CREATE TABLE IF NOT EXISTS `orders` (
  id INT AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL COMMENT '用戶資訊',
  order_date DATE NOT NULL COMMENT '訂單日期',
  payment_status ENUM('Pending', 'Paid', 'Cancelled') NOT NULL COMMENT '付款狀態',
  shipping_status ENUM('Pending', 'Shipped', 'Delivered') NOT NULL COMMENT '配送狀態',
  total_amount DECIMAL(10, 2) NOT NULL COMMENT '訂單總金額',
  FOREIGN KEY (user_id) REFERENCES users(id)
) COMMENT='訂單資料表';

-- 建立產品資料表
CREATE TABLE IF NOT EXISTS `products` (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL COMMENT '產品名稱',
  description TEXT COMMENT '描述',
  price DECIMAL(10, 2) NOT NULL COMMENT '價格',
  discount INT NOT NULL DEFAULT 0 COMMENT '折扣',
  stock INT NOT NULL COMMENT '庫存',
  sku VARCHAR(50) NOT NULL COMMENT 'SKU(庫存單位)',
  image_url VARCHAR(200) COMMENT '圖片',
  category VARCHAR(100) NOT NULL COMMENT '產品分類',
  enabled INT NOT NULL DEFAULT 1 COMMENT '啟用1 關閉0'
) COMMENT='產品資料表';

-- 建立家族資料表
CREATE TABLE IF NOT EXISTS `family` (
  id INT(10) PRIMARY KEY NOT NULL COMMENT '編號',
  name VARCHAR(100) NOT NULL COMMENT '名稱',
  nickname VARCHAR(50) NOT NULL COMMENT '暱稱',
  birthday VARCHAR(50) NOT NULL COMMENT '生日',
  age INT NOT NULL COMMENT '年齡',
  chinesezodiac VARCHAR(20) NOT NULL COMMENT '生肖',
  zodiacsign VARCHAR(20) NOT NULL COMMENT '星座',
  occupation VARCHAR(100) NOT NULL COMMENT '職業',
  extension VARCHAR(10) NOT NULL COMMENT '分機',
  profileimage VARCHAR(255) NOT NULL COMMENT '大頭貼'
) COMMENT='家族資料';

-- 建立獵物資料表
CREATE TABLE IF NOT EXISTS `boy` (
  id INT AUTO_INCREMENT PRIMARY KEY NOT NULL COMMENT '編號',
  name VARCHAR(100) NOT NULL COMMENT '名稱',
  district VARCHAR(100) NOT NULL COMMENT '地區',
  occupation VARCHAR(100) NOT NULL COMMENT '職業'
) COMMENT='海豬獵物公式書';

-- 建立社團資料表
CREATE TABLE IF NOT EXISTS `societies` (
  name VARCHAR(100) NOT NULL COMMENT '社團名稱',
  money VARCHAR(100) NOT NULL COMMENT '社團經費',
  PRIMARY KEY(`name`)
) COMMENT='社團資料';

-- 建立社團成員資料表
CREATE TABLE IF NOT EXISTS `societies_user` (
  user VARCHAR(100) NOT NULL COMMENT '姓名',
  societiesname VARCHAR(100) NOT NULL COMMENT '參加的社團',
  identity VARCHAR(100) NOT NULL COMMENT '社團身份',
  PRIMARY KEY(`user`)
) COMMENT='團員資料';