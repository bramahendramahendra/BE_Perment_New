package model

type IndikatorPencapaian struct {
	IndikatorWarna string  `gorm:"column:indikator_warna"`
	IndikatorValue float64 `gorm:"column:indikator_value"`
}
