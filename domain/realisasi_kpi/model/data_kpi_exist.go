package model

type DataKpiExist struct {
	Tahun              string `gorm:"column:tahun"`
	Triwulan           string `gorm:"column:triwulan"`
	Kostl              string `gorm:"column:kostl"`
	KostlTx            string `gorm:"column:kostl_tx"`
	Status             int    `gorm:"column:status"`
	StatusDesc         string `gorm:"column:status_desc"`
	EntryUserRealisasi string `gorm:"column:entry_user_realisasi"`
	EntryNameRealisasi string `gorm:"column:entry_name_realisasi"`
}
