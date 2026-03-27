package model

type MstMethod struct {
	IdMethodUse int    `gorm:"column:id_method_use"`
	NamaMethod  string `gorm:"column:nama_method"`
	DescMethod  string `gorm:"column:desc_method"`
	Tahun       string `gorm:"column:tahun"`
	Triwulan    string `gorm:"column:triwulan"`
	EntryUser   string `gorm:"column:entry_user"`
	EntryName   string `gorm:"column:entry_name"`
	EntryTime   string `gorm:"column:entry_time"`
}
