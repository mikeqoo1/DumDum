package nici

type Nici struct {
	Id       int    `gorm:"primaryKey;type:int(11) NOT NULL auto_increment ;column:id"`
	Name     string `gorm:"primaryKey;type:varchar(50) NOT NULL;column:name"`
	Blood    string `gorm:"type:varchar(10) NOT NULL;column NOT NULL:blood"`
	Starsign string `gorm:"type:varchar(20) NOT NULL;column:starsign"`
	Series   string `gorm:"type:varchar(50) NOT NULL;column:series"`
}

func (Nici) TableName() string {
	return "nici"
}

