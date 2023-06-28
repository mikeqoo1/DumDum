package main

import (
	"DumDum/lib/basic"
	nici "DumDum/lib/nici"
	"DumDum/lib/pvc"
	shuming "DumDum/lib/shuming"
	tidb "DumDum/lib/tidb"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	_ "DumDum/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var (
	//Googl雲端版本
	IsGoogle string
	conn     *gorm.DB

	userobj    []shuming.User
	orderobj   []shuming.Order
	productobj []shuming.Product

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
		url := "http://ip2c.org/" + clientIP
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

//	@Summary		測試
//	@Description	給書銘測試
//	@Tags			Test
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/example [get]
func hi腦包(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "腦包書銘兒, 你的網站想做啥阿?? 方便確認API方向",
	})
}

//	@Summary		取得User資料
//	@Description	回傳所有User的資料 跟 筆數
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [get]
func hiUser(c *gin.Context) {
	results := conn.Order("id desc").Find(&userobj)
	if results.Error != nil {
		Logger().Error("取得User資料錯誤", results.Error.Error())
	}
	Logger().Info("取得User資料", results)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   userobj,
		"msg":    "腦包書銘兒, 你好, 晚上峽谷見",
		//"errmsg": results.Error.Error(),
	})
}

//	@Summary		增加User
//	@Description	增加User的資料
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			name		body		string	true	"使用者名稱"
//	@Param			email		body		string	true	"電子信箱"
//	@Param			password	body		string	true	"密碼"
//	@Param			address		body		string	true	"住址"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [post]
func addUser(c *gin.Context) {
	username := c.PostForm("name")
	email := c.PostForm("email")
	pwd := c.PostForm("password")
	address := c.PostForm("address")
	Logger().Info("增加User資料:", username, email, pwd, address)
	var result shuming.User
	conn.First(&result, "username = ?", username)
	if result.Username == username {
		Logger().Error("User名稱重複了", result, username)
		c.JSON(http.StatusBadRequest, gin.H{
			"data": result,
			"msg":  "名稱重複了",
		})
		return
	} else {
		新腦包 := shuming.User{
			Username:     username,
			Email:        email,
			Password:     pwd,
			Address:      address,
			Payment_info: "付款資訊現在先用假的",
		}
		conn.Save(&新腦包)
		Logger().Info("新的User資料", 新腦包)
		c.JSON(http.StatusOK, gin.H{
			"data": 新腦包,
			"msg":  "增加客戶",
		})
	}
}

//	@Summary		更新User
//	@Description	更新User的資料
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id			body		string	true	"使用者ID"
//	@Param			name		body		string	true	"使用者名稱"
//	@Param			email		body		string	true	"電子信箱"
//	@Param			password	body		string	true	"密碼"
//	@Param			address		body		string	true	"住址"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [put]
func updateUser(c *gin.Context) {
	id := c.PostForm("id")
	uid, _ := strconv.ParseUint(id, 10, 64)
	username := c.PostForm("name")
	email := c.PostForm("email")
	pwd := c.PostForm("password")
	address := c.PostForm("address")
	Logger().Info("更新User資料:", uid, username, email, pwd, address)
	腦包 := shuming.User{
		ID:           uid,
		Username:     username,
		Email:        email,
		Password:     pwd,
		Address:      address,
		Payment_info: "付款資訊現在先用假的",
	}
	conn.Save(&腦包)
	c.JSON(http.StatusOK, gin.H{
		"data": 腦包,
		"msg":  "更新客戶資料",
	})
}

//	@Summary		刪掉User
//	@Description	刪掉User的資料
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"使用者ID"
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [delete]
func deleteUser(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	Logger().Info("刪掉User資料:", uid)
	conn.Delete(&shuming.User{}, uid)
	c.JSON(http.StatusOK, gin.H{
		"msg": "刪除客戶資料",
	})
}

//	@Summary		取得商品資料
//	@Description	回傳所有商品的資料 跟 筆數
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [get]
func hiProduct(c *gin.Context) {
	results := conn.Order("id desc").Find(&productobj)
	fmt.Println(results)
	if results.Error != nil {
		Logger().Error("取得商品資料錯誤", results.Error.Error())
	}
	fmt.Println(1)
	Logger().Info("取得商品資料", productobj)
	fmt.Println(2)
	fmt.Println(results.Error.Error())
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   productobj,
		"msg":    "商品列表列出來",
		//"errmsg": results.Error.Error(),
	})
}

//	@Summary		新增商品資料
//	@Description	新增商品資料
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			name		body		string	true	"商品名稱"
//	@Param			description	body		string	false	"描述"
//	@Param			price		body		string	false	"價格"
//	@Param			stock		body		string	false	"庫存"
//	@Param			sku			body		string	false	"庫存單位"
//	@Param			imageURL	body		string	false	"圖片"
//	@Param			category	body		string	false	"商品分類"
//	@Param			enabled		body		string	false	"商品啟用(0/1)"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [post]
func addProduct(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	stock := c.PostForm("stock")
	sku := c.PostForm("sku")
	url := c.PostForm("imageURL")
	category := c.PostForm("category")
	enabled := c.PostForm("enabled")
	//給初始值
	if price == "" {
		price = "1200"
	}
	if stock == "" {
		price = "999"
	}
	var result shuming.Product
	Logger().Info("新增商品資料:", name, description, price, stock, sku, url, category, enabled)
	conn.First(&result, "name = ?", name)
	if result.Name == name {
		Logger().Error("商品名稱重複了:", result, name)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "商品名稱重複了",
		})
		return
	} else {
		pricefff, err := strconv.ParseFloat(price, 64)
		if err != nil {
			Logger().Error("商品價格錯誤:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "商品價格錯誤",
			})
			return
		}
		stockiii, err := strconv.Atoi(stock)
		if err != nil {
			Logger().Error("商品庫存錯誤:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "商品庫存錯誤",
			})
			return
		}
		var 啟用 bool
		啟用 = false
		if enabled == "1" {
			啟用 = true
		} else if enabled == "0" {
			啟用 = false
		}
		腦包商品 := shuming.Product{
			Name:        name,
			Description: description,
			Price:       pricefff,
			Stock:       stockiii,
			SKU:         sku,
			ImageURL:    url,
			Category:    category,
			Is_enabled:  啟用,
		}
		Logger().Info("增加商品:", 腦包商品)
		conn.Save(&腦包商品)

		c.JSON(http.StatusOK, gin.H{
			"data": 腦包商品,
			"msg":  "增加商品",
		})
	}
}

//	@Summary		更新商品資料
//	@Description	更新商品資料
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id			body		string	true	"商品ID"
//	@Param			name		body		string	true	"商品名稱"
//	@Param			description	body		string	false	"描述"
//	@Param			price		body		string	false	"價格"
//	@Param			stock		body		string	false	"庫存"
//	@Param			sku			body		string	false	"庫存單位"
//	@Param			imageURL	body		string	false	"圖片"
//	@Param			category	body		string	false	"商品分類"
//	@Param			enabled		body		string	false	"商品啟用(0/1)"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [put]
func updateProduct(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	stock := c.PostForm("stock")
	sku := c.PostForm("sku")
	url := c.PostForm("imageURL")
	category := c.PostForm("category")
	enabled := c.PostForm("enabled")
	Logger().Info("更新商品資料:", name, description, price, stock, sku, url, category, enabled)
	pricefff, err := strconv.ParseFloat(price, 64)
	if err != nil {
		Logger().Error("商品價格錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "商品價格錯誤",
		})
		return
	}
	stockiii, err := strconv.Atoi(stock)
	if err != nil {
		Logger().Error("商品庫存錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "商品庫存錯誤",
		})
		return
	}
	var 啟用 bool
	啟用 = false
	if enabled == "1" {
		啟用 = true
	} else if enabled == "0" {
		啟用 = false
	} else {
		Logger().Error("商品狀態錯誤:", name, description, price, stock, sku, url, category, enabled, 啟用)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "商品狀態錯誤",
		})
		return
	}
	腦包商品 := shuming.Product{
		Name:        name,
		Description: description,
		Price:       pricefff,
		Stock:       stockiii,
		SKU:         sku,
		ImageURL:    url,
		Category:    category,
		Is_enabled:  啟用,
	}
	Logger().Info("更新商品:", 腦包商品)
	conn.Save(&腦包商品)
	c.JSON(http.StatusOK, gin.H{
		"data": 腦包商品,
		"msg":  "更新商品",
	})
}

//	@Summary		刪掉商品資料
//	@Description	刪掉商品資料
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"商品ID"
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [delete]
func deleteProduct(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	Logger().Info("刪掉商品資料:", uid)
	conn.Delete(&shuming.Product{}, uid)
	c.JSON(http.StatusOK, gin.H{
		"msg": "刪除商品資料",
	})
}

//	@Summary		取得訂單清單
//	@Description	回傳所有訂單的資料 跟 筆數
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [get]
func hiOrder(c *gin.Context) {
	results := conn.Order("id desc").Find(&orderobj)
	if results.Error != nil {
		Logger().Error("取得訂單錯誤", results.Error.Error())
	}
	Logger().Info("取得訂單清單", results)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   orderobj,
		"msg":    "訂單通通列出來",
		//"errmsg": results.Error.Error(),
	})
}

//	@Summary		新增訂單
//	@Description	新增訂單
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			user			body		string	true	"用戶名稱"
//	@Param			total_amount	body		string	true	"訂單總金額"
//	@Success		200				{object}	shuming.UserResponse
//	@Failure		400				{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [post]
func addOrder(c *gin.Context) {
	user := c.PostForm("user")
	total_amount := c.PostForm("total_amount")
	var result shuming.User
	conn.First(&result, "username = ?", user)
	if result.ID <= 0 {
		Logger().Error("查無此人:", result, user)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "查無此人" + user,
		})
	}

	total_amountffff, err := strconv.ParseFloat(total_amount, 64)
	if err != nil {
		Logger().Error("訂單總金額錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "訂單總金額錯誤",
		})
		return
	}
	腦包訂單 := shuming.Order{
		UserID:         result.ID,
		OrderDate:      time.Now(),
		PaymentStatus:  "Paid",
		ShippingStatus: "Shipped",
		TotalAmount:    total_amountffff,
	}
	Logger().Info("增加訂單:", 腦包訂單)
	conn.Create(&腦包訂單)
	msg := result.Username + "增加訂單"
	c.JSON(http.StatusOK, gin.H{
		"data": 腦包訂單,
		"msg":  msg,
	})
}

//	@Summary		更新訂單
//	@Description	更新訂單
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			user_id	body		string	true	"用戶資訊"
//	@Success		200		{object}	shuming.UserResponse
//	@Failure		400		{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [put]
func updateOrder(c *gin.Context) {
	user_id := c.PostForm("user_id")
	Logger().Info("更新訂單:", user_id)
	user_idiii, err := strconv.Atoi(user_id)
	if err != nil {
		Logger().Error("user_id錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user_id錯誤",
		})
	}

	conn.Where("user_id = ?", user_idiii).Find(&orderobj)
	c.JSON(http.StatusOK, gin.H{
		"data": orderobj,
		"msg":  "更新腦包訂單",
	})
}

//	@Summary		刪掉訂單
//	@Description	刪掉訂單
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"訂單ID"
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [delete]
func deleteOrder(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	Logger().Info("刪除訂單:", uid)
	conn.Delete(&shuming.Order{}, uid)
	c.JSON(http.StatusOK, gin.H{
		"msg": "刪除訂單資料",
	})
}

//	@Summary		取得報表
//	@Description	回傳統計值
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/report [get]
func hireport(c *gin.Context) {
	results := conn.Order("id desc").Find(&orderobj)
	results2 := conn.Order("id desc").Find(&productobj)
	results3 := conn.Order("id desc").Find(&userobj)
	Logger().Info("取得報表", orderobj, productobj, userobj)
	c.JSON(http.StatusOK, gin.H{
		"record":  results.RowsAffected,
		"data":    orderobj,
		"msg":     "訂單通通列出來",
		"errmsg":  results.Error.Error(),
		"errmsg2": results2.Error.Error(),
		"errmsg3": results3.Error.Error(),
	})
}

//	@title						書銘的API
//	@version					1.0
//	@description				This is a sample server celler server.
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						127.0.0.1:6620
//	@BasePath					/shumingyu
//	@securityDefinitions.basic	BasicAuth
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
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

		// 注入模板變量，並渲染模板
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   title,
			"message": message,
		})
	})

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
		otherRouter.GET("/", other)
		otherRouter.GET("/pig", otherPig)
		otherRouter.POST("/pig/translate", pigtranslate)
		otherRouter.GET("/love", otherLove)
	}

	shumingyuRouter := router.Group("/shumingyu")
	{
		shumingyuRouter.GET("/example", hi腦包)
		shumingyuRouter.GET("/user", hiUser)
		shumingyuRouter.POST("/user", addUser)
		shumingyuRouter.PUT("/user", updateUser)
		shumingyuRouter.DELETE("/user", deleteUser)

		shumingyuRouter.GET("/product", hiProduct)
		shumingyuRouter.POST("/product", addProduct)
		shumingyuRouter.PUT("/product", updateProduct)
		shumingyuRouter.DELETE("/product", deleteProduct)

		shumingyuRouter.GET("/order", hiOrder)
		shumingyuRouter.POST("/order", addOrder)
		shumingyuRouter.PUT("/order", updateOrder)
		shumingyuRouter.DELETE("/order", deleteOrder)

		shumingyuRouter.GET("/report", hireport)
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
