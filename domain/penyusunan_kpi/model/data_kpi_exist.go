package model

type DataKpiExist struct {
	Tahun      string `gorm:"column:tahun"`
	Triwulan   string `gorm:"column:triwulan"`
	Kostl      string `gorm:"column:kostl"`
	KostlTx    string `gorm:"column:kostl_tx"`
	Status     int    `gorm:"column:status"`
	StatusDesc string `gorm:"column:status_desc"`
	EntryUser  string `gorm:"column:entry_user"`
	EntryName  string `gorm:"column:entry_name"`
}
