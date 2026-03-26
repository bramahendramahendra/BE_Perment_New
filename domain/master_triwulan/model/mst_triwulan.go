package model

type MstTriwulan struct {
	IdTriwulan string `gorm:"column:id_triwulan"`
	Triwulan   string `gorm:"column:triwulan"`
}
