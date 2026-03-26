package model

type MstTahun struct {
	Id         int `gorm:"column:id"`
	BatasAtas  int `gorm:"column:batas_atas"`
	BatasBawah int `gorm:"column:batas_bawah"`
}
