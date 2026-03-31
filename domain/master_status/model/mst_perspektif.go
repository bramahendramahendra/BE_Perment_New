package model

type (
	MstStatus struct {
		IdStatus   int    `gorm:"column:id_status"`
		StatusDesc string `gorm:"column:status_desc"`
	}
)
