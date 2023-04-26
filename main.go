package main

import (
	tidb "DumDum/lib/tidb"
	"fmt"
)

type Nici struct {
	Id       int    `gorm:"primaryKey;type:int(11) NOT NULL auto_increment ;column:id"`
	Name     string `gorm:"primaryKey;type:varchar(50) NOT NULL;column:name"`
	Blood    string `gorm:"type:varchar(10) NOT NULL;column NOT NULL:blood"`
	Starsign string `gorm:"type:varchar(20) NOT NULL;column:starsign"`
	Series   string `gorm:"type:varchar(50) NOT NULL;column:series"`
}

func (Nici) TableName() string {
	return "nici"
}

func main() {
	mydb := tidb.NewTiDB("192.168.199.235")
	mydb.Database = "sea"
	mydb.User = "mike"
	mydb.Passwd = "110084"
	mydb.Ip = "192.168.199.235"
	conn, err := mydb.GetDB()
	if err != nil {
		fmt.Println("DB連線失敗->" + err.Error())
	}

	var niciobj []Nici

	//First=> SELECT * FROM `nici` ORDER BY `nici`.`id` LIMIT 1
	result := conn.First(&niciobj)
	fmt.Println(result.RowsAffected) // 找到的筆數
	fmt.Println(niciobj[0].Name)

	//Find=> SELECT * FROM `nici`
	results := conn.Find(&niciobj)
	fmt.Println(results.RowsAffected)
	fmt.Println(niciobj[10].Name)
}
