package model

// =============================================
// DATABASE MODEL
// =============================================

// DataKpi merepresentasikan tabel data_kpi (header utama pengajuan KPI)
// type DataKpi struct {
// 	IdPengajuan    string `gorm:"column:id_pengajuan"`
// 	Tahun          string `gorm:"column:tahun"`
// 	Triwulan       string `gorm:"column:triwulan"`
// 	Kostl          string `gorm:"column:kostl"`
// 	KostlTx        string `gorm:"column:kostl_tx"`
// 	Orgeh          string `gorm:"column:orgeh"`
// 	OrgehTx        string `gorm:"column:orgeh_tx"`
// 	EntryUser      string `gorm:"column:entry_user"`
// 	EntryName      string `gorm:"column:entry_name"`
// 	EntryTime      string `gorm:"column:entry_time"`
// 	ApprovalPosisi string `gorm:"column:approval_posisi"`
// 	ApprovalList   string `gorm:"column:approval_list"`
// 	Status         *int   `gorm:"column:status"`
// }

// DataKpiDetail merepresentasikan tabel data_kpi_detail (satu baris per KPI)
// type DataKpiDetail struct {
// 	IdPengajuan         string `gorm:"column:id_pengajuan"`
// 	IdDetail            string `gorm:"column:id_detail"`
// 	Tahun               string `gorm:"column:tahun"`
// 	Triwulan            string `gorm:"column:triwulan"`
// 	IdKpi               string `gorm:"column:id_kpi"`
// 	Kpi                 string `gorm:"column:kpi"`
// 	Rumus               string `gorm:"column:rumus"`
// 	IdPersfektif        string `gorm:"column:id_perspektif"`
// 	IdKeteranganProject string `gorm:"column:id_keterangan_project"`
// }

// DataKpiSubDetail merepresentasikan tabel data_kpi_subdetail (satu baris per Sub KPI dari Excel)
// type DataKpiSubDetail struct {
// 	IdPengajuan               string  `gorm:"column:id_pengajuan"`
// 	IdDetail                  string  `gorm:"column:id_detail"`
// 	IdSubDetail               string  `gorm:"column:id_sub_detail"`
// 	Tahun                     string  `gorm:"column:tahun"`
// 	Triwulan                  string  `gorm:"column:triwulan"`
// 	IdKpi                     string  `gorm:"column:id_kpi"`
// 	Kpi                       string  `gorm:"column:kpi"`
// 	Rumus                     string  `gorm:"column:rumus"`
// 	Otomatis                  string  `gorm:"column:otomatis"`
// 	Bobot                     float64 `gorm:"column:bobot"`
// 	Capping                   string  `gorm:"column:capping"`
// 	TargetTriwulan            string  `gorm:"column:target_triwulan"`
// 	TargetKuantitatifTriwulan float64 `gorm:"column:target_kuantitatif_triwulan"`
// 	TargetTahunan             string  `gorm:"column:target_tahunan"`
// 	TargetKuantitatifTahunan  float64 `gorm:"column:target_kuantitatif_tahunan"`
// 	DeskripsiGlossary         string  `gorm:"column:deskripsi_glossary"`
// 	ItemQualifier             string  `gorm:"column:item_qualifier"`
// 	DeskripsiQualifier        string  `gorm:"column:deskripsi_qualifier"`
// 	TargetQualifier           string  `gorm:"column:target_qualifier"`
// 	IdKeteranganProject       string  `gorm:"column:id_keterangan_project"`
// 	IdQualifier               string  `gorm:"column:id_qualifier"`
// }

// DataChallengeDetail merepresentasikan tabel data_challenge_detail
// type DataChallengeDetail struct {
// 	IdPengajuan        string `gorm:"column:id_pengajuan"`
// 	IdDetailChallenge  string `gorm:"column:id_detail_challenge"`
// 	Tahun              string `gorm:"column:tahun"`
// 	Triwulan           string `gorm:"column:triwulan"`
// 	NamaChallenge      string `gorm:"column:nama_challenge"`
// 	DeskripsiChallenge string `gorm:"column:deskripsi_challenge"`
// }

// DataMethodDetail merepresentasikan tabel data_method_detail
// type DataMethodDetail struct {
// 	IdPengajuan     string `gorm:"column:id_pengajuan"`
// 	IdDetailMethod  string `gorm:"column:id_detail_method"`
// 	Tahun           string `gorm:"column:tahun"`
// 	Triwulan        string `gorm:"column:triwulan"`
// 	NamaMethod      string `gorm:"column:nama_method"`
// 	DeskripsiMethod string `gorm:"column:deskripsi_method"`
// }

// UserOrgeh digunakan untuk query orgeh dan orgeh_tx dari tabel user
type UserOrgeh struct {
	Orgeh   string `gorm:"column:orgeh"`
	OrgehTx string `gorm:"column:orgeh_tx"`
}
