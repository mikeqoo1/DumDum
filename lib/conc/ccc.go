package concord

import (
	"DumDum/lib/basic"
	"DumDum/lib/tidb"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func (Family) TableName() string {
	return "family"
}

var familyobj []Family

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
