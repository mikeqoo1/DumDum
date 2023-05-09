package tidb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	viper "github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	//產生sql log
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + "_sql.log"
	//log文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//寫入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	newLogger := logger.New(
		log.New(src, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	// 先開啟
	conn, err = gorm.Open(mysql.Open(addr), &gorm.Config{Logger: newLogger})
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
