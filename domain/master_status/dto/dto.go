package dto

type (
	MasterStatusResponse struct {
		IdStatus   int    `gorm:"column:id_status"`
		StatusDesc string `gorm:"column:status_desc"`
	}
)
