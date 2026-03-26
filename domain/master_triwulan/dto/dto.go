package dto

type TriwulanResponse struct {
	IdTriwulan string `gorm:"column:id_triwulan"`
	Triwulan   string `gorm:"column:triwulan"`
}
