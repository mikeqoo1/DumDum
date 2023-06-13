package shuming

import "time"

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
	Stock       int
	SKU         string
	ImageURL    string
	Category    string
	Is_enabled  bool
}

type UserResponse struct {
	Data   string `json:"data"`
	Msg    string `json:"msg"`
	Record int    `json:"record"`
}

type ErrorResponse struct {
	Message string `json:"message"`
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
