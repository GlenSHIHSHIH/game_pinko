package model

type PlinkoOdds struct {
	Id            int         `gorm:"primaryKey;type:int" json:"id"`
	PlinkoBallsId int         `gorm:"type:int" json:"plink_balls_id"`
	Odd           float32     `gorm:"type:float" json:"odd"`
	PlinkoBalls   PlinkoBalls `gorm:"foreignKey:PlinkoBallsId" `
}
