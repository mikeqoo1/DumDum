package tidb

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	viper "github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MaxLifetime  int = 10
	MaxOpenConns int = 10
	MaxIdleConns int = 10
)

type TiDB struct {
	Ip       string
	Port     int
	Database string
	User     string
	Passwd   string
}

// NewTiDB 產生一個DB實例
func NewTiDB() *TiDB {
	viper.SetConfigName("config") // 指定文件的名稱
	viper.AddConfigPath("config") // 配置文件和執行檔目錄
	err := viper.ReadInConfig()   // 根據以上定讀取文件
	if err != nil {
		fmt.Println("Fatal error config file" + err.Error())
		os.Exit(0)
	}
	host := viper.GetString("DB.host")
	port := viper.GetInt("DB.port")
	user := viper.GetString("DB.user")
	pw := viper.GetString("DB.password")
	db := viper.GetString("DB.database")
	//初始化
	tidb := &TiDB{
		Ip:       host,
		Port:     port,
		Database: db,
		User:     user,
		Passwd:   pw,
	}
	return tidb
}

func (tidb *TiDB) GetDB() (*gorm.DB, error) {
	var err error
	var conn *gorm.DB
	var db *sql.DB
	//組合sql連線字串
	addr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", tidb.User, tidb.Passwd, tidb.Ip, tidb.Port, tidb.Database)

	// 先開啟
	conn, err = gorm.Open(mysql.Open(addr), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
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
