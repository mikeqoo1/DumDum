package main

import (
	"DumDum/lib/basic"
	concord "DumDum/lib/conc"
	nici "DumDum/lib/nici"
	"DumDum/lib/pvc"
	shuming "DumDum/lib/shuming"
	"DumDum/lib/tidb"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "DumDum/docs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	//Googl雲端版本
	IsGoogle string
	//要不要開PVC
	IsPvc string

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

func init() {
	var errdb error
	mydb := tidb.NewTiDB()
	_, errdb = mydb.GetDB()
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
		url := "http://ip2c.org/" + clientIP
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("解析IP失敗", err.Error())
			basic.Logger().Error("解析IP失敗", err.Error())
		}
		body, _ := io.ReadAll(res.Body)
		bodystr := string(body)
		bbb := strings.Split(bodystr, ";")
		if bbb[0] == "0" {
			fmt.Println("解析IP的API失敗", url)
			basic.Logger().Error("解析IP的API失敗", url)
		} else if bbb[0] == "1" {
			fmt.Println("Two-letter: " + bbb[1])
			fmt.Println("Three-letter: " + bbb[2])
			fmt.Println("Full name: " + bbb[3])
			basic.Logger().Info("國別碼2碼: ", bbb[1])
			basic.Logger().Info("國別碼3碼: ", bbb[2])
			basic.Logger().Info("國家全名: ", bbb[3])
		} else if bbb[0] == "2" {
			fmt.Println("解析IP的API失敗(Not found in database)")
			basic.Logger().Error("解析IP的API失敗(Not found in database)")
		}
	}
}

func Searchconcords(c *gin.Context) {
	// cid := c.PostForm("cid")
	// cid = StrPad(cid, 12, "0", "LEFT")
	cid := "Q00000000001"
	orderno := c.PostForm("orderno")
	orderno = concord.StrPad(orderno, 5, "0", "RIGHT")
	stock := c.PostForm("stock")
	stock = concord.StrPad(stock, 6, " ", "RIGHT")
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
	basic.Logger().Info("上市櫃查詢電文", msg)
	myreport := <-p.FixReportCh
	reportmsg := myreport.Account + myreport.OrderID
	c.JSON(http.StatusOK, gin.H{
		"OrderMsg":  msg,
		"ReportMsg": reportmsg,
	})
}

func SearchconcordsEM(c *gin.Context) {
	ct := time.Now()
	HHMMSS := ct.Format("150405")
	stock := c.PostForm("stock")
	stock = concord.StrPad(stock, 6, " ", "RIGHT")
	bhno := c.PostForm("bhno")
	delimiter := "\x01"
	msg := "80001=03" + delimiter + "80002=05" + delimiter + "80003=03" + delimiter + "80004=0000" + delimiter + "81005=000" + delimiter + "55=" + stock + delimiter + "80024=" + HHMMSS + delimiter + "80014=Q0000001" + delimiter
	pem.SetbrokId(bhno)
	pem.SetwtmpId("Q0000001")
	clientEM.SendCh <- msg
	msg = p.CreateSearchMessagesEM(msg)
	basic.Logger().Info("興櫃查詢電文", msg)
	myreport := <-pem.FixReportCh
	reportmsg := myreport.Account + myreport.OrderID
	c.JSON(http.StatusOK, gin.H{
		"OrderMsg":  msg,
		"ReportMsg": reportmsg,
	})
}

// @title						書銘的API
// @version					1.0
// @description				This is a sample server celler server.
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.url				http://www.swagger.io/support
// @contact.email				support@swagger.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						127.0.0.1:6620
// @BasePath					/shumingyu
// @securityDefinitions.basic	BasicAuth
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
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
	router.Use(basic.LoggerToFile())
	router.Use(checkIP())

	// 設置模板路徑
	router.LoadHTMLGlob("templates/*.html")
	// 載入靜態文件
	router.Static("/static", "./static")
	router.Static("/docs", "./docs")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")

	//載入靜態資源 一般是上傳的資源 例如上傳的圖檔還是文件
	router.StaticFS("/upload", http.Dir("upload"))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 首頁
	router.GET("/", func(c *gin.Context) {
		// 定義模板變量
		title := "Nici家族"
		message := "歡迎來到Nici家族"
		donate := "網站經營不意 跪求贊助或是我幫放廣告 請用GitHub聯絡我"

		// 注入模板變量，並渲染模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   title,
			"message": message,
			"donate":  donate,
		})
	})

	router.GET("/alice", concord.Alice)
	router.POST("/alice/love", concord.AliceLove)
	router.GET("/alice/boys", concord.AliceBoys)
	router.GET("/alice/newboy", concord.AliceBoy)
	router.POST("/alice/addboy", concord.SadBoy)
	niciRouter := router.Group("/nici")
	{
		niciRouter.GET("/", nici.Love)           //列出所有
		niciRouter.GET("/destiny", nici.Destiny) //顯示輸入畫面
		niciRouter.GET("/newfriend", nici.Newfriend)
		niciRouter.GET("/price", nici.Price)
		niciRouter.POST("/conform", nici.Conform)
		niciRouter.POST("/update", nici.Update)
		niciRouter.POST("/priceee", nici.Priceee)
	}

	otherRouter := router.Group("/other")
	{
		otherRouter.GET("/", concord.Other)
		otherRouter.GET("/pig", concord.OtherPig)
		otherRouter.POST("/pig/translate", concord.Pigtranslate)
		otherRouter.GET("/love", concord.OtherLove)
	}

	shumingyuRouter := router.Group("/shumingyu")
	{
		shumingyuRouter.POST("/login", shuming.Login)
		shumingyuRouter.GET("/user", shuming.HiUser)
		shumingyuRouter.POST("/user", shuming.AddUser)
		shumingyuRouter.PUT("/user", shuming.UpdateUser)
		shumingyuRouter.DELETE("/user", shuming.DeleteUser)

		shumingyuRouter.GET("/product", shuming.HiProduct)
		shumingyuRouter.POST("/product", shuming.AddProduct)
		shumingyuRouter.PUT("/product", shuming.UpdateProduct)
		shumingyuRouter.DELETE("/product", shuming.DeleteProduct)
		shumingyuRouter.GET("/productcategory", shuming.GetProductCategory)
		shumingyuRouter.POST("/getoneproduct", shuming.GetOneProduct)

		shumingyuRouter.GET("/order", shuming.HiOrder)
		shumingyuRouter.POST("/order", shuming.AddOrder)
		shumingyuRouter.PUT("/order", shuming.UpdateOrder)
		shumingyuRouter.DELETE("/order", shuming.DeleteOrder)

		shumingyuRouter.GET("/report", shuming.Hireport)
	}

	//if IsGoogle == "NO" {
	concordsRouter := router.Group("/concords")
	{
		concordsRouter.GET("/", concord.Concords)
		concordsRouter.GET("/EM", concord.ConcordsEM)
		concordsRouter.GET("/family", concord.ConcordsFamily)
		concordsRouter.POST("/search", Searchconcords)
		concordsRouter.POST("/searchEM", SearchconcordsEM)
		concordsRouter.GET("/societies", concord.GetSocietiesAll)
		concordsRouter.GET("/moneylist", concord.GetSocietiesMoney)
		concordsRouter.GET("/adduser", concord.GetUserPage)
		concordsRouter.GET("/addevent", concord.GetEventPage)
		concordsRouter.POST("/newsocietiesuser", concord.AddSocietiesUser)
		concordsRouter.POST("/newsocietiesevent", concord.AddSocietiesEvent)

		concordsRouter.GET("/orderhome", concord.OrderHome)
		concordsRouter.GET("/neworder", concord.GetOrderPage)
		concordsRouter.POST("/order", concord.Order)
	}

	// if IsPvc == "YES" {
	// 	if err := client.Connect("192.168.199.185:7052"); err != nil {
	// 		fmt.Println("Error connecting TWSE/OTC:", err)
	// 		return
	// 	}
	// 	defer client.Close()

	// 	if err := clientEM.Connect("192.168.199.250:7080"); err != nil {
	// 		fmt.Println("Error connecting EM:", err)
	// 		return
	// 	}
	// 	defer clientEM.Close()

	// 	message := p.CreateRegisterMsg()
	// 	client.SendCh <- message
	// 	go client.SendMessages()
	// 	go client.ReceiveMessages()
	// 	go p.ParseMessages(client)

	// 	emmsg := pem.CreateRegisterMsg()
	// 	clientEM.SendCh <- emmsg
	// 	go clientEM.SendMessages()
	// 	go clientEM.ReceiveMessages()
	// 	go pem.ParseMessages(clientEM)
	// }
	//}

	go router.RunTLS(":443", "./certs/server.crt", "./certs/server.key")
	err = router.Run(addr)
	if err != nil {
		fmt.Println("Nici網頁啟動失敗" + err.Error())
	}
}
