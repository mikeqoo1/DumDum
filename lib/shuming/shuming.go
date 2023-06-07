package shuming

type Shuming struct {
	Id       int    `gorm:"primaryKey;type:int(11) NOT NULL auto_increment ;column:id"`
	Account  string `gorm:"primaryKey;type:varchar(50) NOT NULL;column:account"`
	Username string `gorm:"type:varchar(50) NOT NULL;column NOT NULL:username"`
	Status   int    `gorm:"type:int(2) NOT NULL DEFAULT 1;column:status"`
}

type UserResponse struct {
	// 定义你的字段...
}

type ErrorResponse struct {
	// 定义你的字段...
}

func (Shuming) TableName() string {
	return "userlist"
}
