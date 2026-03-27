package model

type MstChallenge struct {
	IdChallenge   int    `gorm:"column:id_challenge"`
	NamaChallenge string `gorm:"column:nama_challenge"`
	DescChallenge string `gorm:"column:desc_challenge"`
	Tahun         string `gorm:"column:tahun"`
	Triwulan      string `gorm:"column:triwulan"`
	EntryUser     string `gorm:"column:entry_user"`
	EntryName     string `gorm:"column:entry_name"`
	EntryTime     string `gorm:"column:entry_time"`
}
