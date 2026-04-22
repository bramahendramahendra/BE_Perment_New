package model

type DataResultDetail struct {
	IdPengajuan     string `gorm:"column:id_pengajuan"`
	IdDetailResult  string `gorm:"column:id_detail_result"`
	Tahun           string `gorm:"column:tahun"`
	Triwulan        string `gorm:"column:triwulan"`
	NamaResult      string `gorm:"column:nama_result"`
	DeskripsiResult string `gorm:"column:deskripsi_result"`
}
