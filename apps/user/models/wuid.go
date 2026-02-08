package models

type Wuid struct {
	H int  `gorm:"column:h;type:int(10);autoIncrement;uniqueIndex" json:"h"`
	X int8 `gorm:"column:x;type:tinyint(4);primaryKey;default:0" json:"x"`
}

func (Wuid) TableName() string {
	return "wuid"
}
