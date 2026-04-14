package model

type DataChallengeDetail struct {
	IdPengajuan        string `gorm:"column:id_pengajuan"`
	IdDetailChallenge  string `gorm:"column:id_detail_challenge"`
	Tahun              string `gorm:"column:tahun"`
	Triwulan           string `gorm:"column:triwulan"`
	NamaChallenge      string `gorm:"column:nama_challenge"`
	DeskripsiChallenge string `gorm:"column:deskripsi_challenge"`
}
