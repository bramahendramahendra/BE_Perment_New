package model

type MstResult struct {
	IdResult   int    `gorm:"column:id_result"`
	NamaResult string `gorm:"column:nama_result"`
	DescResult string `gorm:"column:desc_result"`
	Tahun      string `gorm:"column:tahun"`
	Triwulan   string `gorm:"column:triwulan"`
	EntryUser  string `gorm:"column:entry_user"`
	EntryName  string `gorm:"column:entry_name"`
	EntryTime  string `gorm:"column:entry_time"`
}
