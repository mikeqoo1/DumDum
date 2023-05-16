package pvc

import "fmt"

type pvc struct {
	start string // 封包起始字元  0xAA

	bodyId string // 無須填此欄位填 0XFF, FIX 必填 新單 'D' 改量 'G' 刪單 'F' 查詢 'H' 回報 '8', 回補'2'
	// TMP 檔送'S'(檔案由上端自行傳送 pvc bypass),  檔送's'(pvc 自行於ftp 下載)

	exCode string // TMP 普通 '4' 零股 '5' 標借 '6' 拍賣議價 '7' 標構 '8' 定價 '9' 借貸 'B' 檔收送'1' 盤中零股 'C'
	// 期貨TMP T盤 '1' T + 1盤 '2', 無須填此欄位填 0XFF,
	// fix '0' 普通, '2' 零股, '7' 定價 'C' 盤中零股

	msgTy string // 'T'(集中, 期貨) 'O'(櫃買, 選擇) 'R' (傳送) 'r'(註冊,僅連線時使用), 其他則不理

	typeId string // fix '0', twse tmp '1' tafiex '2',.... 註冊時會判別是否符合,不符會切斷

	connId string // 註冊碼'A'~'Z','a'~'z','1' ~ '9', 會區分大小寫, 同一註冊碼只能有一連線,
	// 重送(msgTy == 'R')如填'A'只送註冊碼為'A'的訊息,填'0'則送當時保存在queue 所有資料

	pvcId string // 無須填此欄位填 0XFF

	rtnState string // reserve 無須填此欄位填 0XFF 如PVC未準備好會填'V'送回, 資料已到PVC但未與證交連線則回'F', bodydata 不變
	// TMP 'Y' success 'E' error 'T' timeout 'N' 沒有檔案, 成交回報 'A'
	// 期貨TMP 'Y' 為回補 'B' 為盤後刪單

	bodyLen string // 4個文數字不足前面補 '0', 大小不包括最後'\n'

	brokId string // 委託證券商代號,目前4位多的後面補 0XFF

	wtmpId string // twse tmp Id, 無須填此欄位填 0XFF

	bodydata string // 字串最後\n,但bodyLen內容大小不包括最後'\n'
	// 註冊'r'時請將IP填入以做為是否為認證判斷
	// 註冊通過認證後 orderGW 會回送PVC數量資料
	// 所送之資料為TN及ON (T,O為市場別 N 1 byte數值為PVC連線的數量)
	// 如未通過認證則回傳 '1'typeId(系統別)錯誤 '2'connId錯誤 '3'IP錯誤
}

// NewPVC 產生PVC物件
func NewPVC() *pvc {
	pvc := &pvc{
		start:    "\xAA",
		bodyId:   "\xFF",
		exCode:   "\xFF",
		msgTy:    "r",
		typeId:   "0",
		connId:   "Q",
		pvcId:    "\xFF",
		rtnState: "\xFF",
		bodyLen:  "0000",
		brokId:   "0000\xFF\xFF\xFF\xFF",
		wtmpId:   "\xFF",
		bodydata: "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF",
	}
	return pvc
}

func (pvc *pvc) GetbodyId() string {
	return pvc.bodyId
}

func (pvc *pvc) GetexCode() string {
	return pvc.exCode
}

func (pvc *pvc) GetmsgTy() string {
	return pvc.msgTy
}

func (pvc *pvc) GettypeId() string {
	return pvc.typeId
}

func (pvc *pvc) GetconnId() string {
	return pvc.connId
}

func (pvc *pvc) GetpvcId() string {
	return pvc.pvcId
}

func (pvc *pvc) GetrtnState() string {
	return pvc.pvcId
}

func (pvc *pvc) GetbodyLen() string {
	return pvc.bodyLen
}

func (pvc *pvc) GetbrokId() string {
	return pvc.brokId
}

func (pvc *pvc) GetwtmpId() string {
	return pvc.wtmpId
}

func (pvc *pvc) Getbodydata() string {
	return pvc.bodydata
}

func (pvc *pvc) CreateRegisterMsg() string {
	msg := pvc.start + pvc.bodyId + pvc.exCode + pvc.msgTy + pvc.typeId + pvc.connId + pvc.pvcId + pvc.rtnState + pvc.bodyLen + pvc.brokId + pvc.wtmpId + pvc.bodydata
	fmt.Println("註冊電文:[" + msg + "]")
	return msg
}

func (pvc *pvc) ParseMessages(recvmsg string) {
	fmt.Println("收到的電文:[" + recvmsg + "]")
	pvc.msgTy = recvmsg[3:4]
	pvc.bodydata = recvmsg[32:]
	if pvc.msgTy == "r" {
		fmt.Println("註冊成功")
		fmt.Println("bodydata=" + pvc.bodydata)
	}
}
