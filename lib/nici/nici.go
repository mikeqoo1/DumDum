package nici

import (
	tidb "DumDum/lib/tidb"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Nici struct {
	Id       int    `gorm:"primaryKey;type:int(11) NOT NULL auto_increment ;column:id"`
	Name     string `gorm:"primaryKey;type:varchar(50) NOT NULL;column:name"`
	Blood    string `gorm:"type:varchar(10) NOT NULL;column NOT NULL:blood"`
	Starsign string `gorm:"type:varchar(20) NOT NULL;column:starsign"`
	Series   string `gorm:"type:varchar(50) NOT NULL;column:series"`
	Img      string `gorm:"type:varchar(50) NOT NULL;column:img"`
}

func (Nici) TableName() string {
	return "nici"
}

var niciobj []Nici

// getFileName 統計不同系列的最大名稱
func getFileName(series string) string {
	files, err := os.ReadDir("./static/img")
	if err != nil {
		log.Fatal(err)
	}
	var friendsort []int
	var dragonsort []int
	var bearsort []int
	var unicornsort []int
	//先整理排序
	for _, file := range files {
		if strings.Contains(file.Name(), "friend") {
			str := strings.Split(file.Name(), "friend")
			num := strings.Split(str[1], ".jpg")
			i, _ := strconv.Atoi(num[0])
			friendsort = append(friendsort, i)
		} else if strings.Contains(file.Name(), "dragon") {
			str := strings.Split(file.Name(), "dragon")
			num := strings.Split(str[1], ".jpg")
			i, _ := strconv.Atoi(num[0])
			dragonsort = append(dragonsort, i)
		} else if strings.Contains(file.Name(), "bear") {
			str := strings.Split(file.Name(), "bear")
			num := strings.Split(str[1], ".jpg")
			i, _ := strconv.Atoi(num[0])
			bearsort = append(bearsort, i)
		} else if strings.Contains(file.Name(), "unicorn") {
			str := strings.Split(file.Name(), "unicorn")
			num := strings.Split(str[1], ".jpg")
			i, _ := strconv.Atoi(num[0])
			unicornsort = append(unicornsort, i)
		}
	}

	fileName := ""
	if series == "動物好夥伴" {
		sort.Ints(friendsort)
		max := friendsort[len(friendsort)-1]
		s := strconv.Itoa(max + 1)
		fileName = "friend" + s + ".jpg"
	} else if series == "恐龍時代" {
		sort.Ints(dragonsort)
		max := dragonsort[len(dragonsort)-1]
		s := strconv.Itoa(max + 1)
		fileName = "dragon" + s + ".jpg"
	} else if series == "熊熊家族" {
		sort.Ints(bearsort)
		max := bearsort[len(bearsort)-1]
		s := strconv.Itoa(max + 1)
		fileName = "bear" + s + ".jpg"
	} else if series == "獨角精靈" {
		sort.Ints(unicornsort)
		max := unicornsort[len(unicornsort)-1]
		s := strconv.Itoa(max + 1)
		fileName = "unicorn" + s + ".jpg"
	}
	return fileName
}

// checkname 更新小夥伴防呆
func checkname(name string) int64 {
	var sqlstr string
	var results *gorm.DB
	sqlstr = "name = ?"
	results = tidb.Globalconn.Where(sqlstr, name).Find(&niciobj)
	return results.RowsAffected
}

func Love(c *gin.Context) {
	//Find=> SELECT * FROM `nici ORDER BY series desc`
	results := tidb.Globalconn.Order("series desc").Find(&niciobj)
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "nici.html", gin.H{
		"title":  title,
		"record": results.RowsAffected,
		"data":   niciobj,
	})
}

func Destiny(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "destiny.html", gin.H{
		"title": title,
	})
}

func Conform(c *gin.Context) {
	title := "Nici好夥伴"
	blood := c.PostForm("blood")
	star := c.PostForm("star")
	var sqlstr string
	var results *gorm.DB
	if blood == "" || len(blood) == 0 {
		sqlstr = "starsign = ?"
		results = tidb.Globalconn.Where(sqlstr, star).Find(&niciobj)
	} else if star == "" || len(star) == 0 {
		sqlstr = "blood = ?"
		results = tidb.Globalconn.Where(sqlstr, blood).Find(&niciobj)
	} else {
		sqlstr = "starsign = ? AND blood = ?"
		results = tidb.Globalconn.Where(sqlstr, star, blood).Find(&niciobj)
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"blood":   blood,
	// 	"star":    star,
	// 	"record":  results.RowsAffected,
	// 	"results": niciobj,
	// })

	if results.RowsAffected == 0 {
		notfound := Nici{
			Name:     "查無資料",
			Blood:    blood,
			Starsign: star,
			Series:   "未知",
			Img:      "NotFound.jpg",
		}
		niciobj = append(niciobj, notfound)
		c.HTML(http.StatusOK, "only.html", gin.H{
			"title":  title,
			"record": results.RowsAffected,
			"data":   niciobj,
		})
	} else {
		c.HTML(http.StatusOK, "only.html", gin.H{
			"title":  title,
			"record": results.RowsAffected,
			"data":   niciobj,
		})
	}
}

func Newfriend(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "newfriend.html", gin.H{
		"title": title,
	})
}

func Update(c *gin.Context) {
	name := c.PostForm("name")
	blood := c.PostForm("blood")
	star := c.PostForm("star")
	series := c.PostForm("series")
	file, _ := c.FormFile("file0") // get file from form input name 'file0'
	yesorno := checkname(name)
	if yesorno > 0 {
		c.JSON(http.StatusOK, gin.H{
			"狀態": "新增失敗",
			"原因": "夥伴已經存在了",
		})
	} else {
		filename := getFileName(series)
		file.Filename = filename
		c.SaveUploadedFile(file, "static/img/"+file.Filename) // save file to tmp folder in current directory
		new := Nici{
			Name:     name,
			Blood:    blood,
			Starsign: star,
			Series:   series,
			Img:      file.Filename,
		}
		tidb.Globalconn.Save(&new)

		title := "Nici家族"
		message := "歡迎來到Nici家族"
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   title,
			"message": message,
		})
	}
}

func Price(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "price.html", gin.H{
		"title": title,
	})
}

func Priceee(c *gin.Context) {
	//Hevisaurus
	//Wind Rose
	//Dragon Force
	typeee := c.PostForm("typeee")
	size := c.PostForm("size")
	fmt.Println(typeee)
	fmt.Println(size)
	price_str := ""
	newprice_str := ""
	price1 := 0.0
	price2 := 0.0
	if typeee == "鑰匙圈" {
		price_str = "340~590"
		price1 = 340 * 0.9 * 0.95
		price2 = 590 * 0.9 * 0.95
	} else if typeee == "杯套" {
		price_str = "420~650"
		price1 = 420 * 0.9 * 0.95
		price2 = 650 * 0.9 * 0.95
	} else if typeee == "玩偶" {
		if size == "29" {
			price_str = "840~1090"
			price1 = 840 * 0.9 * 0.95
			price2 = 1090 * 0.9 * 0.95
		} else if size == "39" {
			price_str = "840~1600"
			price1 = 840 * 0.9 * 0.95
			price2 = 1600 * 0.9 * 0.95
		} else if size == "49" {
			price_str = "1290~1890"
			price1 = 1290 * 0.9 * 0.95
			price2 = 1890 * 0.9 * 0.95
		} else if size == "59" {
			price_str = "1470~2100"
			price1 = 1470 * 0.9 * 0.95
			price2 = 2100 * 0.9 * 0.95
		} else if size == "79" {
			price_str = "3570起"
			price1 = 3570 * 0.9 * 0.95
			price2 = 0
		} else if size == "100" {
			price_str = "5840起(誠品沒販售, 不適用估價)"
			price1 = 5840 * 0.9 * 0.95
			price2 = 0
		}
	}
	s1 := fmt.Sprintf("%f", price1)
	s2 := fmt.Sprintf("%f", price2)
	if s2 == "0" {
		s2 = "很多的錢錢"
	}
	newprice_str = s1 + "~" + s2
	c.HTML(http.StatusOK, "price.html", gin.H{
		"OGprice":  price_str,
		"Newprice": newprice_str,
	})
}
