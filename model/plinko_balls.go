package model

type PlinkoBalls struct {
	Id    int    `gorm:"type:int;primaryKey" json:"id"`
	Color string `gorm:"comment:顏色" json:"color"`
	Type  int    `gorm:"type:int;index:idx_name,unique;comment:球種" json:"type"`
	Row   int    `gorm:"type:int;index:idx_name,unique;comment:行數" json:"row"`
}
