package main

import (
	"DumDum/lib/basic"
	nici "DumDum/lib/nici"
	"DumDum/lib/pvc"
	tidb "DumDum/lib/tidb"
	"fmt"
	"io"
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
	//Googl雲端版本
	IsGoogle string
	conn     *gorm.DB
	niciobj  []nici.Nici
	// 產生上市櫃客戶端物件
	client = basic.TCPClient{
		SendCh:    make(chan string, 1024),
		ReceiveCh: make(chan string, 1024),
	}
	// 產生興櫃客戶端物件
	clientEM = basic.TCPClient{
		SendCh:    make(chan string, 1024),
		ReceiveCh: make(chan string, 1024),
	}
	p   = pvc.NewPVC()
	pem = pvc.NewPVC()
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
	results = conn.Where(sqlstr, name).Find(&niciobj)
	return results.RowsAffected
}

// StrPad
// input string 原字串
// padLength int 規定補完後的字串長度
// padString string 自定義填充字串
// padType string 填充類型:LEFT(向左填充,自動補齊位數), 默認右側
func StrPad(input string, padLength int, padString string, padType string) string {

	output := ""
	inputLen := len(input)

	if inputLen >= padLength {
		return input
	}

	padStringLen := len(padString)
	needFillLen := padLength - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
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

// checkIP 解析IP
func checkIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		url := "https://ip2c.org/" + clientIP
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("解析IP失敗", err.Error())
			Logger().Error("解析IP失敗", err.Error())
		}
		body, _ := io.ReadAll(res.Body)
		bodystr := string(body)
		bbb := strings.Split(bodystr, ";")
		if bbb[0] == "0" {
			fmt.Println("解析IP的API失敗", url)
			Logger().Error("解析IP的API失敗", url)
		} else if bbb[0] == "1" {
			fmt.Println("Two-letter: " + bbb[1])
			fmt.Println("Three-letter: " + bbb[2])
			fmt.Println("Full name: " + bbb[3])
			Logger().Info("國別碼2碼: ", bbb[1])
			Logger().Info("國別碼3碼: ", bbb[2])
			Logger().Info("國家全名: ", bbb[3])
		} else if bbb[0] == "2" {
			fmt.Println("解析IP的API失敗(Not found in database)")
			Logger().Error("解析IP的API失敗(Not found in database)")
		}
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
}

func price(c *gin.Context) {
	title := "Nici好夥伴"
	c.HTML(http.StatusOK, "price.html", gin.H{
		"title": title,
	})
}

func priceee(c *gin.Context) {
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

/*Nici API區*/

func getallnici(c *gin.Context) {
	results := conn.Order("series desc").Find(&niciobj)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   niciobj,
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
	newpig = strings.Replace(newpig, "牙", "阿", -1)
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
	// cid := c.PostForm("cid")
	// cid = StrPad(cid, 12, "0", "LEFT")
	cid := "Q00000000001"
	orderno := c.PostForm("orderno")
	orderno = StrPad(orderno, 5, "0", "RIGHT")
	stock := c.PostForm("stock")
	stock = StrPad(stock, 6, " ", "RIGHT")
	bs := c.PostForm("bscode")
	oederflag := c.PostForm("orderflag")
	excode := c.PostForm("excode")
	bhno := c.PostForm("bhno")
	delimiter := "\x01"
	msg := "11=" + cid + delimiter + "37=" + orderno + delimiter + "55=" + stock + delimiter + "54=" + bs + delimiter + "10000=" + oederflag + delimiter + "10002=" + excode + delimiter
	p.SetbrokId(bhno)
	p.SetwtmpId(cid)
	msg = p.CreateSearchMessages(msg)
	client.SendCh <- msg
	Logger().Info("上市櫃查詢電文", msg)
	myreport := <-p.FixReportCh
	reportmsg := myreport.Account + myreport.OrderID
	c.JSON(http.StatusOK, gin.H{
		"OrderMsg":  msg,
		"ReportMsg": reportmsg,
	})
}

func concordsEM(c *gin.Context) {
	c.HTML(http.StatusOK, "coo2.html", gin.H{})
}

func searchconcordsEM(c *gin.Context) {
	ct := time.Now()
	HHMMSS := ct.Format("150405")
	stock := c.PostForm("stock")
	stock = StrPad(stock, 6, " ", "RIGHT")
	bhno := c.PostForm("bhno")
	delimiter := "\x01"
	msg := "80001=03" + delimiter + "80002=05" + delimiter + "80003=03" + delimiter + "80004=0000" + delimiter + "81005=000" + delimiter + "55=" + stock + delimiter + "80024=" + HHMMSS + delimiter + "80014=Q0000001" + delimiter
	pem.SetbrokId(bhno)
	pem.SetwtmpId("Q0000001")
	clientEM.SendCh <- msg
	msg = p.CreateSearchMessagesEM(msg)
	Logger().Info("興櫃查詢電文", msg)
	myreport := <-pem.FixReportCh
	reportmsg := myreport.Account + myreport.OrderID
	c.JSON(http.StatusOK, gin.H{
		"OrderMsg":  msg,
		"ReportMsg": reportmsg,
	})
}

/*腦包書銘區*/

func hi腦包(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "腦包書銘兒, 你好, 晚上峽谷見",
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
	router.Use(checkIP())

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
		niciRouter.GET("/newfriend", newfriend)
		niciRouter.GET("/price", price)
		niciRouter.POST("/conform", conform)
		niciRouter.POST("/update", update)
		niciRouter.POST("/priceee", priceee)
		niciRouter.GET("/api/all", getallnici)
	}

	otherRouter := router.Group("/other")
	{
		otherRouter.GET("/", other)
		otherRouter.GET("/pig", otherPig)
		otherRouter.POST("/pig/translate", pigtranslate)
		otherRouter.GET("/love", otherLove)
	}

	shumingyuRouter := router.Group("/yu")
	{
		shumingyuRouter.GET("/", hi腦包)
	}

	if IsGoogle == "NO" {
		concordsRouter := router.Group("/concords")
		{
			concordsRouter.GET("/", concords)
			concordsRouter.GET("/EM", concordsEM)
			concordsRouter.POST("/search", searchconcords)
			concordsRouter.POST("/searchEM", searchconcordsEM)
		}

		if err := client.Connect("192.168.199.185:7052"); err != nil {
			fmt.Println("Error connecting TWSE/OTC:", err)
			return
		}
		defer client.Close()

		if err := clientEM.Connect("192.168.199.250:7080"); err != nil {
			fmt.Println("Error connecting EM:", err)
			return
		}
		defer clientEM.Close()

		message := p.CreateRegisterMsg()
		client.SendCh <- message
		go client.SendMessages()
		go client.ReceiveMessages()
		go p.ParseMessages(client)

		emmsg := pem.CreateRegisterMsg()
		clientEM.SendCh <- emmsg
		go clientEM.SendMessages()
		go clientEM.ReceiveMessages()
		go pem.ParseMessages(clientEM)
	}

	err = router.Run(addr)
	if err != nil {
		fmt.Println("Nici網頁啟動失敗" + err.Error())
	}
}
