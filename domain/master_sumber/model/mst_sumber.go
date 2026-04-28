package model

type (
	MstSumber struct {
		IdSumber int    `gorm:"column:id_sumber"`
		Sumber   string `gorm:"column:sumber"`
	}
)
