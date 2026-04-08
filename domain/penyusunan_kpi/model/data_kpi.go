package model

type DataKpi struct {
	IdPengajuan    string `gorm:"column:id_pengajuan"`
	Tahun          string `gorm:"column:tahun"`
	Triwulan       string `gorm:"column:triwulan"`
	Kostl          string `gorm:"column:kostl"`
	KostlTx        string `gorm:"column:kostl_tx"`
	Orgeh          string `gorm:"column:orgeh"`
	OrgehTx        string `gorm:"column:orgeh_tx"`
	EntryUser      string `gorm:"column:entry_user"`
	EntryName      string `gorm:"column:entry_name"`
	EntryTime      string `gorm:"column:entry_time"`
	ApprovalPosisi string `gorm:"column:approval_posisi"`
	ApprovalList   string `gorm:"column:approval_list"`
	Status         string `gorm:"column:status"`
	StatusDesc     string `gorm:"column:status_desc"`
}
