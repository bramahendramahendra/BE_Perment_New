package model

type DataMethodDetail struct {
	IdPengajuan     string `gorm:"column:id_pengajuan"`
	IdDetailMethod  string `gorm:"column:id_detail_method"`
	Tahun           string `gorm:"column:tahun"`
	Triwulan        string `gorm:"column:triwulan"`
	NamaMethod      string `gorm:"column:nama_method"`
	DeskripsiMethod string `gorm:"column:deskripsi_method"`
	RealisasiMethod string `gorm:"column:realisasi_method"`
	LampiranEvidence string `gorm:"column:lampiran_evidence"`
}
