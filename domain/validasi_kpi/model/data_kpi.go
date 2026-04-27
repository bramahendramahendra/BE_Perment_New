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
	Status         int    `gorm:"column:status"`
	StatusDesc     string `gorm:"column:status_desc"`

	EntryUserRealisasi       string `gorm:"column:entry_user_realisasi"`
	EntryNameRealisasi       string `gorm:"column:entry_name_realisasi"`
	EntryTimeRealisasi       string `gorm:"column:entry_time_realisasi"`
	ApprovalListRealisasi    string `gorm:"column:approval_list_realisasi"`
	CatatanTolakan           string `gorm:"column:catatan_tolakan"`
	TotalBobot               string `gorm:"column:total_bobot"`
	TotalPencapaian          string `gorm:"column:total_pencapaian"`
	TotalBobotPengurang      string `gorm:"column:total_bobot_pengurang"`
	TotalPencapaianPost      string `gorm:"column:total_pencapaian_post"`
	EntryUserValidasi        string `gorm:"column:entry_user_validasi"`
	EntryNameValidasi        string `gorm:"column:entry_name_validasi"`
	EntryTimeValidasi        string `gorm:"column:entry_time_validasi"`
	ApprovalListValidasi     string `gorm:"column:approval_list_validasi"`
	LampiranValidasi         string `gorm:"column:lampiran_validasi"`
	QualifierOverallValidasi string `gorm:"column:qualifier_overall_validasi"`
}
