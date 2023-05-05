package main

import (
	nici "DumDum/lib/nici"
	tidb "DumDum/lib/tidb"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	conn    *gorm.DB
	niciobj []nici.Nici
)

func init() {
	var errdb error
	mydb := tidb.NewTiDB("127.0.0.1")
	mydb.Database = "sea"
	// mydb.User = "mike"
	// mydb.Passwd = "110084"
	mydb.User = "root"
	mydb.Passwd = ""
	mydb.Ip = "127.0.0.1"
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
		fmt.Println("進入星座", star)
		sqlstr = "starsign = ?"
		results = conn.Where(sqlstr, star).Find(&niciobj)
	} else if star == "" || len(star) == 0 {
		fmt.Println("進入血型", blood)
		sqlstr = "blood = ?"
		results = conn.Where(sqlstr, blood).Find(&niciobj)
	} else {
		fmt.Println("都進", star, blood)
		sqlstr = "starsign = ? AND blood = ?"
		results = conn.Where(sqlstr, star, blood).Find(&niciobj)
	}
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

func main() {
	// addr := fmt.Sprintf("%s:%d", "127.0.0.1", 6620)
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 80)
	router := gin.Default()
	router.Use(crosHandler())

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
	}

	err := router.Run(addr)
	if err != nil {
		fmt.Println("Nici網頁啟動失敗" + err.Error())
	}
}
