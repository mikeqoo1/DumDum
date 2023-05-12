package main

import (
	nici "DumDum/lib/nici"
	tidb "DumDum/lib/tidb"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	conn    *gorm.DB
	niciobj []nici.Nici
)

func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
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

	logger := logrus.New()
	logger.Out = src
	//設log等級
	logger.SetLevel(logrus.DebugLevel)
	//log格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
	return func(c *gin.Context) {
		startTime := time.Now()
		// 處理請求
		c.Next()
		endTime := time.Now()
		// 執行時間
		latencyTime := endTime.Sub(startTime)
		// 請求方式
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		//log格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// getFileName 統計不同系列的最大名稱
func getFileName(series string) string {
	files, err := ioutil.ReadDir("./static/img")
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

func init() {
	var errdb error
	mydb := tidb.NewTiDB()
	conn, errdb = mydb.GetDB()
	if errdb != nil {
		fmt.Println("DB連線失敗->" + errdb.Error())
		os.Exit(0)
	}
}

// crosHandler 處理跨域問題
func crosHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //請求頭部
		if origin != "" {
			//接收客戶端傳送的origin (重要)
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//伺服器支援的所有跨域請求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允許跨域設定可以返回其他子段，可以自定義欄位
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, "+
				"X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, "+
				"X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, "+
				"Content-Type, Pragma, token, openid, opentoken")
			//允許瀏覽器(客戶端)可以解析的頭部 (重要)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, "+
				"Expires, Last-Modified, Pragma, FooBar")
			//設定快取時間
			c.Header("Access-Control-Max-Age", "172800")
			//允許客戶端傳遞校驗資訊比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允許型別校驗
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		c.Next()
	}
}

/*Nici區*/

func love(c *gin.Context) {
	//Find=> SELECT * FROM `nici ORDER BY series desc`
	results := conn.Order("series desc").Find(&niciobj)
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "nici.html", gin.H{
		"title":  title,
		"record": results.RowsAffected,
		"data":   niciobj,
	})
}

func destiny(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "destiny.html", gin.H{
		"title": title,
	})
}

func conform(c *gin.Context) {
	title := "Nici好夥伴"
	blood := c.PostForm("blood")
	star := c.PostForm("star")
	var sqlstr string
	var results *gorm.DB
	if blood == "" || len(blood) == 0 {
		sqlstr = "starsign = ?"
		results = conn.Where(sqlstr, star).Find(&niciobj)
	} else if star == "" || len(star) == 0 {
		sqlstr = "blood = ?"
		results = conn.Where(sqlstr, blood).Find(&niciobj)
	} else {
		sqlstr = "starsign = ? AND blood = ?"
		results = conn.Where(sqlstr, star, blood).Find(&niciobj)
	}
	Logger().Info("輸入的條件", star, blood)
	// c.JSON(http.StatusOK, gin.H{
	// 	"blood":   blood,
	// 	"star":    star,
	// 	"record":  results.RowsAffected,
	// 	"results": niciobj,
	// })

	if results.RowsAffected == 0 {
		notfound := nici.Nici{
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

func newfriend(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "newfriend.html", gin.H{
		"title": title,
	})
}

func update(c *gin.Context) {
	name := c.PostForm("name")
	blood := c.PostForm("blood")
	star := c.PostForm("star")
	series := c.PostForm("series")
	file, _ := c.FormFile("file0") // get file from form input name 'file0'
	filename := getFileName(series)
	file.Filename = filename
	c.SaveUploadedFile(file, "static/img/"+file.Filename) // save file to tmp folder in current directory
	new := nici.Nici{
		Name:     name,
		Blood:    blood,
		Starsign: star,
		Series:   series,
		Img:      file.Filename,
	}
	conn.Save(&new)

	title := "Nici家族"
	message := "歡迎來到Nici家族"
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   title,
		"message": message,
	})
}

/*其他區*/

func other(c *gin.Context) {
	c.HTML(http.StatusOK, "other.html", gin.H{})
}

func otherPig(c *gin.Context) {
	c.HTML(http.StatusOK, "otherpig.html", gin.H{})
}

func pigtranslate(c *gin.Context) {
	pig := c.PostForm("piggg")
	Logger().Info("海豬原文:", pig)
	newpig := strings.Replace(pig, "窩咬", "我要", -1)
	newpig = strings.Replace(newpig, "女森", "女生", -1)
	newpig = strings.Replace(newpig, "蛇摸", "什麼", -1)
	newpig = strings.Replace(newpig, "笑史", "笑死", -1)
	newpig = strings.Replace(newpig, "把拔", "爸爸", -1)
	newpig = strings.Replace(newpig, "縮", "說", -1)
	newpig = strings.Replace(newpig, "窩", "我", -1)
	newpig = strings.Replace(newpig, "尼", "你", -1)
	newpig = strings.Replace(newpig, "惹", "了", -1)
	newpig = strings.Replace(newpig, "倫", "人", -1)
	newpig = strings.Replace(newpig, "ㄅ", "不", -1)
	newpig = strings.Replace(newpig, "ㄇ", "嗎", -1)
	newpig = strings.Replace(newpig, "ㄍ", "個", -1)
	newpig = strings.Replace(newpig, "仍", "來", -1)
	newpig = strings.Replace(newpig, "度", "對", -1)
	newpig = strings.Replace(newpig, "迷", "沒", -1)
	newpig = strings.Replace(newpig, "抗", "看", -1)
	newpig = strings.Replace(newpig, "奏", "揍", -1)
	Logger().Info("翻譯後長:", newpig)
	c.HTML(http.StatusOK, "otherpig.html", gin.H{
		"pig":    pig,
		"newpig": newpig,
	})
}

func otherLove(c *gin.Context) {
	c.HTML(http.StatusOK, "otherlove.html", gin.H{})
}

/*康x區*/

func concords(c *gin.Context) {
	c.HTML(http.StatusOK, "coo.html", gin.H{})
}

func searchconcords(c *gin.Context) {
	cid := c.PostForm("cid")
	orderno := c.PostForm("orderno")
	stock := c.PostForm("stock")
	bs := c.PostForm("bscode")
	oederflag := c.PostForm("orderflag")
	excode := c.PostForm("excode")

	msg := "11=" + cid + "37=" + orderno + "55=" + stock + "54=" + bs + "10000=" + oederflag + "10002=" + excode

	Logger().Info("查詢電文", msg)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func main() {
	viper.SetConfigName("config") // 指定文件的名稱
	viper.AddConfigPath("config") // 配置文件和執行檔目錄
	err := viper.ReadInConfig()   // 根據以上定讀取文件
	if err != nil {
		fmt.Println("Fatal error config file" + err.Error())
		os.Exit(0)
	}
	host := viper.GetString("Server.ip")
	port := viper.GetInt("Server.port")
	addr := fmt.Sprintf("%s:%d", host, port)
	router := gin.Default()
	router.Use(crosHandler())
	router.Use(LoggerToFile())

	// 設置模板路徑
	router.LoadHTMLGlob("templates/*.html")
	// 載入靜態文件
	router.Static("/static", "./static")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	//載入靜態資源 一般是上傳的資源 例如上傳的圖檔還是文件
	router.StaticFS("/upload", http.Dir("upload"))

	// 首頁
	router.GET("/", func(c *gin.Context) {
		// 定義模板變量
		title := "Nici家族"
		message := "歡迎來到Nici家族"

		// 注入模板變量，並渲染模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   title,
			"message": message,
		})
	})

	niciRouter := router.Group("/nici")
	{
		niciRouter.GET("/", love)           //列出所有
		niciRouter.GET("/destiny", destiny) //顯示輸入畫面
		niciRouter.POST("/conform", conform)
		niciRouter.GET("/newfriend", newfriend)
		niciRouter.POST("/update", update)
	}

	otherRouter := router.Group("/other")
	{
		otherRouter.GET("/", other)
		otherRouter.GET("/pig", otherPig)
		otherRouter.POST("/pig/translate", pigtranslate)
		otherRouter.GET("/love", otherLove)
	}

	// concordsRouter := router.Group("/concords")
	// {
	// 	concordsRouter.GET("/", concords)
	// 	concordsRouter.POST("/search", searchconcords)
	// }

	err = router.Run(addr)
	if err != nil {
		fmt.Println("Nici網頁啟動失敗" + err.Error())
	}
}
