package tidb

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	Port         int = 4000
	MaxLifetime  int = 10
	MaxOpenConns int = 10
	MaxIdleConns int = 10
)

type TiDB struct {
	Ip       string
	Database string
	User     string
	Passwd   string
}

//NewTiDB 產生一個DB實例
func NewTiDB(ip string) *TiDB {
	//初始化
	tidb := &TiDB{
		Ip:       ip,
		Database: " ",
		User:     " ",
		Passwd:   " ",
	}
	return tidb
}

func (tidb *TiDB) GetDB() (*gorm.DB, error) {
	var err error
	var conn *gorm.DB
	var db *sql.DB
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", tidb.User, tidb.Passwd, tidb.Ip, Port, tidb.Database)

	// 先開啟
	conn, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//設定ConnMaxLifetime/MaxIdleConns/MaxOpenConns
	db, err = conn.DB()
	if err != nil {
		fmt.Println("get db failed:", err)
		return nil, err
	}
	db.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)

	return conn, nil
}
