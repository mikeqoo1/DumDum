package shuming

import (
	"DumDum/lib/basic"
	"DumDum/lib/tidb"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 定義用戶模型結構
type User struct {
	ID           uint64
	Username     string
	Email        string
	Password     string
	Address      string
	Payment_info string
	CreatedAt    time.Time
}

// 定義訂單模型結構
type Order struct {
	ID             uint64
	UserID         uint64
	User           User `gorm:"foreignKey:UserID"`
	OrderDate      time.Time
	PaymentStatus  string
	ShippingStatus string
	TotalAmount    float64
}

// 定義產品模型結構
type Product struct {
	ID          uint64
	Name        string
	Description string
	Price       float64
	Discount    int
	Stock       int
	SKU         string
	ImageURL    string
	Category    string
	Enabled     bool
}

// 定義報表結構
type Report struct {
	ID          uint64
	Name        string
	Description string
	Price       float64
	Stock       int
	SKU         string
	ImageURL    string
	Category    string
	Enabled     bool
}

type UserResponse struct {
	Data     string `json:"data"`
	Msg      string `json:"msg"`
	Record   int    `json:"record"`
	ErrorMag string `json:"errmsg"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type product_Category struct {
	All              string `json:"All"`
	ExperienceClass  string `json:"ExperienceClass"`
	Aerobics         string `json:"Aerobics"`
	Yoga             string `json:"Yoga"`
	StrengthTraining string `json:"StrengthTraining"`
}

func (User) TableName() string {
	return "users"
}

func (Order) TableName() string {
	return "orders"
}

func (Product) TableName() string {
	return "products"
}

var userobj []User
var orderobj []Order
var productobj []Product

//	@Summary		登入功能
//	@Description	登入功能
//	@Tags			Login
//	@Accept			json
//	@Produce		json
//	@Param			Name		body		string	true	"使用者名稱"
//	@Param			Password	body		string	true	"密碼"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/login [post]
func Login(c *gin.Context) {
	username := c.PostForm("Name")
	pwd := c.PostForm("Password")
	basic.Logger().Info("登入紀錄:", username, pwd)
	var u User
	results := tidb.Globalconn.First(&u, "username = ?", username)
	if u.Username == username && u.Password == pwd {
		c.JSON(http.StatusOK, gin.H{
			"msg": "登入成功",
		})
		return
	} else {
		if u.Username == username && u.Password != pwd {
			basic.Logger().Error("User密碼錯誤", u, username)
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": "密碼錯誤!!!",
			})
			return
		}
		if results.RowsAffected == 0 {
			basic.Logger().Error("查無User", u, username)
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": "帳號錯誤!!!",
			})
			return
		}
		if results.Error != nil {
			basic.Logger().Error("DB錯誤", results.Error.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": results.Error.Error(),
			})
			return
		}
	}
}

//	@Summary		取得User資料
//	@Description	回傳所有User的資料 跟 筆數
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [get]
func HiUser(c *gin.Context) {
	results := tidb.Globalconn.Order("id desc").Find(&userobj)
	if results.Error != nil {
		basic.Logger().Error("取得User資料錯誤", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	basic.Logger().Info("取得User資料", results)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   userobj,
		"msg":    "取得User資料",
	})
}

//	@Summary		增加User
//	@Description	增加User的資料
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			Name		body		string	true	"使用者名稱"
//	@Param			Email		body		string	true	"電子信箱"
//	@Param			Password	body		string	true	"密碼"
//	@Param			Address		body		string	true	"住址"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [post]
func AddUser(c *gin.Context) {
	username := c.PostForm("Name")
	email := c.PostForm("Email")
	pwd := c.PostForm("Password")
	address := c.PostForm("Address")
	basic.Logger().Info("增加User資料:", username, email, pwd, address)
	var result User
	tidb.Globalconn.First(&result, "username = ?", username)
	if result.Username == username {
		basic.Logger().Error("User名稱重複了", result, username)
		c.JSON(http.StatusBadRequest, gin.H{
			"data": result,
			"msg":  "名稱重複了",
		})
		return
	} else {
		新腦包 := User{
			Username:     username,
			Email:        email,
			Password:     pwd,
			Address:      address,
			Payment_info: "付款資訊現在先用假的",
		}
		tidb.Globalconn.Save(&新腦包)
		basic.Logger().Info("新的User資料", 新腦包)
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
//	@Param			Id			body		string	true	"使用者ID"
//	@Param			Name		body		string	true	"使用者名稱"
//	@Param			Email		body		string	true	"電子信箱"
//	@Param			Password	body		string	true	"密碼"
//	@Param			Address		body		string	true	"住址"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/user [put]
func UpdateUser(c *gin.Context) {
	id := c.PostForm("Id")
	uid, _ := strconv.ParseUint(id, 10, 64)
	username := c.PostForm("Name")
	email := c.PostForm("Email")
	pwd := c.PostForm("Password")
	address := c.PostForm("Address")
	basic.Logger().Info("更新User資料:", uid, username, email, pwd, address)
	腦包 := User{
		ID:           uid,
		Username:     username,
		Email:        email,
		Password:     pwd,
		Address:      address,
		Payment_info: "付款資訊現在先用假的",
	}
	tidb.Globalconn.Save(&腦包)
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
func DeleteUser(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	basic.Logger().Info("刪掉User資料:", uid)
	var u User
	results := tidb.Globalconn.First(&u, "id = ?", uid)
	if results.RowsAffected == 0 {
		basic.Logger().Error("找不到User資訊 id=", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "找不到User資訊:" + c.Query("id"),
		})
		return
	}
	results = tidb.Globalconn.Delete(&User{}, uid)
	if results.Error != nil {
		basic.Logger().Error("刪掉User資料錯誤:", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "刪除客戶資料",
		"data": u,
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
func HiProduct(c *gin.Context) {
	results := tidb.Globalconn.Order("id desc").Find(&productobj)
	if results.Error != nil {
		basic.Logger().Error("取得商品資料錯誤", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	basic.Logger().Info("取得商品資料", productobj)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   productobj,
		"msg":    "商品列表列出來",
	})
}

//	@Summary		新增商品資料
//	@Description	新增商品資料
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			Name		body		string	true	"商品名稱"
//	@Param			Description	body		string	false	"描述"
//	@Param			Price		body		string	false	"價格"
//	@Param			Discount	body		string	false	"折扣 例如10代表打9折"
//	@Param			Stock		body		string	false	"庫存"
//	@Param			SKU			body		string	false	"庫存單位"
//	@Param			ImageURL	body		string	false	"圖片"
//	@Param			Category	body		string	false	"商品分類"
//	@Param			Enabled		body		string	false	"商品啟用(0/1)"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [post]
func AddProduct(c *gin.Context) {
	name := c.PostForm("Name")
	description := c.PostForm("Description")
	price := c.PostForm("Price")
	discount := c.PostForm("Discount")
	stock := c.PostForm("Stock")
	sku := c.PostForm("SKU")
	url := c.PostForm("ImageURL")
	category := c.PostForm("Category")
	enabled := c.PostForm("Enabled")
	//給初始值
	if price == "" {
		price = "1200"
	}
	if stock == "" {
		stock = "999"
	}
	if discount == "" {
		discount = "0"
	}
	var result Product
	basic.Logger().Info("新增商品資料:",
		name, description, price, discount, stock, sku, url, category, enabled)
	tidb.Globalconn.First(&result, "name = ?", name)
	if result.Name == name {
		basic.Logger().Error("商品名稱重複了:", result, name)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "商品名稱重複了",
		})
		return
	} else {
		pricefff, err := strconv.ParseFloat(price, 64)
		if err != nil {
			basic.Logger().Error("商品價格錯誤:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": "商品價格錯誤",
			})
			return
		}
		stockiii, err := strconv.Atoi(stock)
		if err != nil {
			basic.Logger().Error("商品庫存錯誤:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": "商品庫存錯誤",
			})
			return
		}
		discountiii, err := strconv.Atoi(discount)
		if err != nil {
			basic.Logger().Error("折扣錯誤:", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{
				"errmsg": "商品折扣錯誤",
			})
			return
		}
		var 啟用 bool
		啟用 = false
		if enabled == "true" {
			啟用 = true
		} else if enabled == "false" {
			啟用 = false
		}
		腦包商品 := Product{
			Name:        name,
			Description: description,
			Price:       pricefff,
			Discount:    discountiii,
			Stock:       stockiii,
			SKU:         sku,
			ImageURL:    url,
			Category:    category,
			Enabled:     啟用,
		}
		basic.Logger().Info("增加商品:", 腦包商品)
		tidb.Globalconn.Save(&腦包商品)

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
//	@Param			Id			body		string	true	"商品ID"
//	@Param			Name		body		string	true	"商品名稱"
//	@Param			Description	body		string	false	"描述"
//	@Param			Price		body		string	false	"價格"
//	@Param			Discount	body		string	false	"折扣 例如10代表打9折"
//	@Param			Stock		body		string	false	"庫存"
//	@Param			SKU			body		string	false	"庫存單位"
//	@Param			ImageURL	body		string	false	"圖片"
//	@Param			Category	body		string	false	"商品分類"
//	@Param			Enabled		body		string	false	"商品啟用(0/1)"
//	@Success		200			{object}	shuming.UserResponse
//	@Failure		400			{object}	shuming.ErrorResponse
//	@Router			/shumingyu/product [put]
func UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.PostForm("ID"), 10, 64)
	name := c.PostForm("Name")
	description := c.PostForm("Description")
	price := c.PostForm("Price")
	discount := c.PostForm("Discount")
	stock := c.PostForm("Stock")
	sku := c.PostForm("SKU")
	url := c.PostForm("ImageURL")
	category := c.PostForm("Category")
	enabled := c.PostForm("Enabled")
	basic.Logger().Info("更新商品ID:", id)
	basic.Logger().Info("更新商品內容:",
		name, ",",
		description, ",", price, ",", stock, ",",
		sku, ",", url, ",", category, ",", enabled)
	pricefff, err := strconv.ParseFloat(price, 64)
	if err != nil {
		basic.Logger().Error("商品價格錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "商品價格錯誤",
		})
		return
	}
	stockiii, err := strconv.Atoi(stock)
	if err != nil {
		basic.Logger().Error("商品庫存錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "商品庫存錯誤",
		})
		return
	}
	discountiii, err := strconv.Atoi(discount)
	if err != nil {
		basic.Logger().Error("折扣錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "商品折扣錯誤",
		})
		return
	}
	var 啟用 bool
	啟用 = false
	if enabled == "true" {
		啟用 = true
	} else if enabled == "false" {
		啟用 = false
	} else {
		basic.Logger().Error("商品狀態錯誤:", enabled)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "商品狀態錯誤",
		})
		return
	}

	var p Product
	results := tidb.Globalconn.First(&p, "id = ?", id)
	if results.RowsAffected == 0 {
		basic.Logger().Warn("找不到商品資訊 id=", id, "變新增一筆")
		pricefff = 1200
		stockiii = 999
		discountiii = 0
	}

	腦包商品 := Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       pricefff,
		Discount:    discountiii,
		Stock:       stockiii,
		SKU:         sku,
		ImageURL:    url,
		Category:    category,
		Enabled:     啟用,
	}
	basic.Logger().Info("更新商品:", 腦包商品)
	tidb.Globalconn.Save(&腦包商品)
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
func DeleteProduct(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	basic.Logger().Info("刪掉商品資料:", uid)
	var p Product
	results := tidb.Globalconn.First(&p, "id = ?", uid)
	if results.RowsAffected == 0 {
		basic.Logger().Error("找不到商品資訊 id=", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "找不到商品資訊:" + c.Query("id"),
		})
		return
	}
	results = tidb.Globalconn.Delete(&Product{}, uid)
	if results.Error != nil {
		basic.Logger().Error("刪掉商品資料錯誤:", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "刪除商品資料",
		"data": p,
	})
}

//	@Summary		取得單一商品
//	@Description	取得單一商品
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string	true	"商品ID"
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/getoneproduct [post]
func GetOneProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.PostForm("ID"), 10, 64)
	basic.Logger().Info("取得單一商品ID:", id)
	var p Product
	results := tidb.Globalconn.First(&p, "id = ?", id)
	if results.RowsAffected == 0 {
		basic.Logger().Warn("找不到商品資訊 id=", id)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "找不到商品資訊",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": p,
			"msg":  "取得商品成功",
		})
	}
}

//	@Summary		取得商品種類
//	@Description	回傳商品的所有種類
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	shuming.UserResponse
//	@Failure		400	{object}	shuming.ErrorResponse
//	@Router			/shumingyu/productcategory [get]
func GetProductCategory(c *gin.Context) {
	results := tidb.Globalconn.Select("category").Group("category").Find(&productobj)
	if results.Error != nil {
		basic.Logger().Error("取得商品資料錯誤", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	var pc product_Category
	pc.All = "全部商品"
	for i := 0; i < len(productobj); i++ {
		if productobj[i].Category == "瑜珈" {
			pc.Yoga = productobj[i].Category
		} else if productobj[i].Category == "有氧" {
			pc.Aerobics = productobj[i].Category
		} else if productobj[i].Category == "體驗" {
			pc.ExperienceClass = productobj[i].Category
		} else if productobj[i].Category == "肌力訓練" {
			pc.StrengthTraining = productobj[i].Category
		}
	}
	basic.Logger().Info("取得商品種類", pc)
	c.JSON(http.StatusOK, gin.H{
		"data": pc,
		"msg":  "全部商品種類",
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
func HiOrder(c *gin.Context) {
	results := tidb.Globalconn.Order("id desc").Find(&orderobj)
	if results.Error != nil {
		basic.Logger().Error("取得訂單錯誤", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	basic.Logger().Info("取得訂單清單", results)
	c.JSON(http.StatusOK, gin.H{
		"record": results.RowsAffected,
		"data":   orderobj,
		"msg":    "訂單通通列出來",
	})
}

//	@Summary		新增訂單
//	@Description	新增訂單
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			User			body		string	true	"用戶名稱"
//	@Param			Total_amount	body		string	true	"訂單總金額"
//	@Success		200				{object}	shuming.UserResponse
//	@Failure		400				{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [post]
func AddOrder(c *gin.Context) {
	user := c.PostForm("User")
	total_amount := c.PostForm("Total_amount")
	var result User
	tidb.Globalconn.First(&result, "username = ?", user)
	if result.ID <= 0 {
		basic.Logger().Error("查無此人:", result, user)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "查無此人" + user,
		})
	}

	total_amountffff, err := strconv.ParseFloat(total_amount, 64)
	if err != nil {
		basic.Logger().Error("訂單總金額錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "訂單總金額錯誤",
		})
		return
	}
	腦包訂單 := Order{
		UserID:         result.ID,
		OrderDate:      time.Now(),
		PaymentStatus:  "Paid",
		ShippingStatus: "Shipped",
		TotalAmount:    total_amountffff,
	}
	basic.Logger().Info("增加訂單:", 腦包訂單)
	tidb.Globalconn.Create(&腦包訂單)
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
//	@Param			User_id	body		string	true	"用戶資訊"
//	@Success		200		{object}	shuming.UserResponse
//	@Failure		400		{object}	shuming.ErrorResponse
//	@Router			/shumingyu/order [put]
func UpdateOrder(c *gin.Context) {
	user_id := c.PostForm("User_id")
	basic.Logger().Info("更新訂單:", user_id)
	user_idiii, err := strconv.Atoi(user_id)
	if err != nil {
		basic.Logger().Error("user_id錯誤:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "user_id錯誤",
		})
	}

	tidb.Globalconn.Where("user_id = ?", user_idiii).Find(&orderobj)
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
func DeleteOrder(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	basic.Logger().Info("刪除訂單:", uid)
	var o Order
	results := tidb.Globalconn.First(&o, "id = ?", uid)
	if results.RowsAffected == 0 {
		basic.Logger().Error("找不到訂單資訊 id=", uid)
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": "找不到訂單資訊:" + c.Query("id"),
		})
		return
	}
	results = tidb.Globalconn.Delete(&Order{}, uid)
	if results.Error != nil {
		basic.Logger().Error("刪掉訂單錯誤:", results.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"errmsg": results.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "刪除訂單資料",
		"data": o,
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
func Hireport(c *gin.Context) {
	results := tidb.Globalconn.Order("id desc").Find(&orderobj)
	results2 := tidb.Globalconn.Order("id desc").Find(&productobj)
	results3 := tidb.Globalconn.Order("id desc").Find(&userobj)
	basic.Logger().Info("取得報表", orderobj, productobj, userobj)
	c.JSON(http.StatusOK, gin.H{
		"record":  results.RowsAffected,
		"data":    orderobj,
		"msg":     "訂單通通列出來",
		"errmsg":  results.Error.Error(),
		"errmsg2": results2.Error.Error(),
		"errmsg3": results3.Error.Error(),
	})
}
