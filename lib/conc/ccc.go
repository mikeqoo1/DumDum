package concord

import (
	"DumDum/lib/basic"
	"DumDum/lib/tidb"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Family struct {
	Id            int
	Name          string
	Nickname      string
	Birthday      string
	Age           int
	Chinesezodiac string
	Zodiacsign    string
	Occupation    string
	Extension     string
	Profileimage  string
}

type Boy struct {
	Id         int
	Name       string
	District   string
	Occupation string
}

type Societies struct {
	Name  string
	Money string
}
type SocietiesUser struct {
	User          string
	Societiesname string
	Identity      string
}
type SocietiesEvent struct {
	Data_date string
	Eventname string
	Money     string
	Person    string
	People    string
}

func (Family) TableName() string {
	return "family"
}

func (Boy) TableName() string {
	return "boy"
}

func (Societies) TableName() string {
	return "societies"
}
func (SocietiesUser) TableName() string {
	return "societies_user"
}
func (SocietiesEvent) TableName() string {
	return "societies_event"
}

var familyobj []Family

var boysobj []Boy

var Societiesobj []Societies
var Societiesuserobj []SocietiesUser
var Societieseventobj []SocietiesEvent

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

func Other(c *gin.Context) {
	c.HTML(http.StatusOK, "other.html", gin.H{})
}

func OtherPig(c *gin.Context) {
	c.HTML(http.StatusOK, "otherpig.html", gin.H{})
}

func Pigtranslate(c *gin.Context) {
	pig := c.PostForm("piggg")
	basic.Logger().Info("海豬原文:", pig)
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
	newpig = strings.Replace(newpig, "ㄞ", "挨", -1)
	newpig = strings.Replace(newpig, "仍", "來", -1)
	newpig = strings.Replace(newpig, "度", "對", -1)
	newpig = strings.Replace(newpig, "迷", "沒", -1)
	newpig = strings.Replace(newpig, "抗", "看", -1)
	newpig = strings.Replace(newpig, "奏", "揍", -1)
	newpig = strings.Replace(newpig, "牙", "阿", -1)
	basic.Logger().Info("翻譯後:", newpig)
	c.HTML(http.StatusOK, "otherpig.html", gin.H{
		"pig":    pig,
		"newpig": newpig,
	})
}

func OtherLove(c *gin.Context) {
	c.HTML(http.StatusOK, "otherlove.html", gin.H{})
}

/*康x區*/

func Concords(c *gin.Context) {
	c.HTML(http.StatusOK, "coo.html", gin.H{})
}

func ConcordsEM(c *gin.Context) {
	c.HTML(http.StatusOK, "coo2.html", gin.H{})
}

func ConcordsFamily(c *gin.Context) {
	results := tidb.Globalconn.Find(&familyobj)
	title := "加菲教夥伴"
	c.HTML(http.StatusOK, "family.html", gin.H{
		"title":  title,
		"record": results.RowsAffected,
		"data":   familyobj,
	})
}

func GetSocietiesAll(c *gin.Context) {
	tidb.Globalconn.Find(&Societiesobj)
	tidb.Globalconn.Find(&Societiesuserobj)
	title := "社團資料"
	c.HTML(http.StatusOK, "societies.html", gin.H{
		"title": title,
		"data":  Societiesobj,
		"data2": Societiesuserobj,
	})
}

func GetSocietiesMoney(c *gin.Context) {
	tidb.Globalconn.Find(&Societieseventobj)
	title := "社團紀錄"
	c.HTML(http.StatusOK, "societiesmoney.html", gin.H{
		"title": title,
		"data":  Societieseventobj,
	})
}

func GetUserPage(c *gin.Context) {
	title := "新增社團成員"
	c.HTML(http.StatusOK, "societiesadd.html", gin.H{
		"title": title,
	})
}

func GetEventPage(c *gin.Context) {
	title := "新增社團活動"
	c.HTML(http.StatusOK, "societiesadd2.html", gin.H{
		"title": title,
	})
}

func AddSocietiesUser(c *gin.Context) {
	uesr := c.PostForm("user")
	societiesname := c.PostForm("societiesname")
	identity := c.PostForm("identity")
	password := c.PostForm("password")
	if password != "dj/3xj/6yjo41;4!" { //dj/3xj/6yjo41;4! 恐龍最棒!
		c.JSON(http.StatusBadRequest, gin.H{
			"狀態": "新增失敗",
			"原因": "代碼失敗",
		})
		return
	}
	新社員 := SocietiesUser{
		User:          uesr,
		Societiesname: societiesname,
		Identity:      identity,
	}
	//fmt.Println(新社員)

	if uesr == "" || societiesname == "" || identity == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"狀態": "新增失敗",
			"原因": "欄位不得為空",
		})
		return
	} else {
		var sqlstr string
		var results *gorm.DB
		sqlstr = "user = ?"
		results = tidb.Globalconn.Where(sqlstr, uesr).Find(&Societiesuserobj)
		if results.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"狀態": "新增失敗",
				"原因": "名稱已經存在了",
			})
			return
		}
	}
	tidb.Globalconn.Create(&新社員)
	c.JSON(http.StatusOK, gin.H{
		"狀態": "新增社員成功",
	})
}

func AddSocietiesEvent(c *gin.Context) {
	date := c.PostForm("date")
	eventname := c.PostForm("eventname")
	money := c.PostForm("money")
	person := c.PostForm("person")
	people := c.PostForm("people")
	password := c.PostForm("password")
	if password != "dj/3xj/6yjo41;4!" { //dj/3xj/6yjo41;4! 恐龍最棒!
		c.JSON(http.StatusBadRequest, gin.H{
			"狀態": "新增失敗",
			"原因": "代碼失敗",
		})
		return
	}
	新活動 := SocietiesEvent{
		Data_date: date,
		Eventname: eventname,
		Money:     money,
		Person:    person,
		People:    people,
	}

	if date == "" || eventname == "" || money == "" || person == "" || people == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"狀態": "新增失敗",
			"原因": "欄位不得為空",
		})
		return
	}
	tidb.Globalconn.Create(&新活動)
	c.JSON(http.StatusOK, gin.H{
		"狀態": "新增活動成功",
	})
}

func OrderHome(c *gin.Context) {
	c.HTML(http.StatusOK, "orderhome.html", gin.H{})
}

func GetOrderPage(c *gin.Context) {
	c.HTML(http.StatusOK, "orderpage.html", gin.H{})
}

func Order(c *gin.Context) {}

/*海豬區*/

func Alice(c *gin.Context) {
	title := "海豬問券"
	c.HTML(http.StatusOK, "alice.html", gin.H{
		"title": title,
	})
}

func AliceLove(c *gin.Context) {
	name := c.PostForm("name")
	age := c.PostForm("age")
	profession := c.PostForm("profession")
	yearmoney := c.PostForm("yearmoney")
	radio1 := c.PostForm("radio1")
	radio2 := c.PostForm("radio2")
	radio3 := c.PostForm("radio3")
	radio4 := c.PostForm("radio4")
	radio5 := c.PostForm("radio5")
	iage, err1 := strconv.Atoi(age)
	iyearmoney, err2 := strconv.Atoi(yearmoney)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "年紀或是年收入錯誤",
		})
		return
	}

	fmt.Println("姓名:", name)
	basic.Logger().Info("姓名:", name)
	fmt.Println("年紀:", age)
	basic.Logger().Info("年紀:", age)
	fmt.Println("職業:", profession)
	basic.Logger().Info("職業:", profession)
	fmt.Println("收入:", yearmoney)
	basic.Logger().Info("收入:", yearmoney)
	fmt.Println("1.答案:", radio1)
	basic.Logger().Info("1.答案:", radio1)
	fmt.Println("2.答案:", radio2)
	basic.Logger().Info("2.答案:", radio2)
	fmt.Println("3.答案:", radio3)
	basic.Logger().Info("3.答案:", radio3)
	fmt.Println("4.答案:", radio4)
	basic.Logger().Info("4.答案:", radio4)
	fmt.Println("5.答案:", radio5)
	basic.Logger().Info("5.答案:", radio5)

	if iage < 18 {
		basic.Logger().Info("年紀太小不符合資格")
		c.JSON(http.StatusOK, gin.H{
			"msg": "年紀太小不符合資格",
		})
		return
	}
	if iyearmoney < 100 {
		basic.Logger().Info("沒有百萬年收海豬看不上 在磨練磨練八")
		c.JSON(http.StatusOK, gin.H{
			"msg": "沒有百萬年收海豬看不上 在磨練磨練八",
		})
		return
	}
	y := 0
	if radio1 == "Yes" {
		y++
	}
	if radio2 == "Yes" {
		y++
	}
	if radio3 == "Yes" {
		y++
	}
	if radio4 == "Yes" {
		y++
	}
	if radio5 == "Yes" {
		y++
	}

	if y >= 5 {
		basic.Logger().Info("天選之人 恭喜有機會成為第N號男")
		c.JSON(http.StatusOK, gin.H{
			"msg": "天選之人 恭喜有機會成為第N號男",
		})
		return
	} else {
		basic.Logger().Info("你不適合海豬 滾八")
		c.JSON(http.StatusOK, gin.H{
			"msg": "你不適合海豬 滾八",
		})
	}
}

func AliceBoys(c *gin.Context) {
	title := "海豬男的公式書"
	results := tidb.Globalconn.Find(&boysobj)
	c.HTML(http.StatusOK, "aliceboys.html", gin.H{
		"title":  title,
		"record": results.RowsAffected,
		"data":   boysobj,
	})
}

func AliceBoy(c *gin.Context) {
	title := "增加海豬男的公式書"
	c.HTML(http.StatusOK, "alicenewboy.html", gin.H{
		"title": title,
	})
}

func SadBoy(c *gin.Context) {
	name := c.PostForm("name")
	district := c.PostForm("district")
	occupation := c.PostForm("occupation")
	sadboy := Boy{
		Name:       name,
		District:   district,
		Occupation: occupation,
	}

	if strings.Contains(name, "林") || strings.Contains(name, "奎") || strings.Contains(name, "翰") ||
		strings.Contains(name, "龍") || strings.Contains(name, "恐") || strings.Contains(name, "揆") ||
		strings.Contains(name, "和") || strings.Contains(name, "汗") || strings.Contains(name, "漢") ||
		strings.Contains(name, "魁") || strings.Contains(name, "淋") || strings.Contains(name, "籠") {
		c.JSON(http.StatusBadRequest, gin.H{
			"原因": "臭海豬 還想搞阿 下去拉",
		})
		return
	}

	if name == "" || district == "" || occupation == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"狀態": "新增失敗",
			"原因": "欄位不得為空",
		})
		return
	} else {
		var sqlstr string
		var results *gorm.DB
		sqlstr = "name = ?"
		results = tidb.Globalconn.Where(sqlstr, name).Find(&boysobj)
		if results.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"狀態": "新增失敗",
				"原因": "名稱已經存在了",
			})
			return
		}
	}
	tidb.Globalconn.Save(&sadboy)
	c.JSON(http.StatusOK, gin.H{
		"狀態": "新增成功",
	})
}
