package model

type DataKpiSubDetail struct {
	IdPengajuan               string  `gorm:"column:id_pengajuan"`
	IdDetail                  string  `gorm:"column:id_detail"`
	IdSubDetail               string  `gorm:"column:id_sub_detail"`
	Tahun                     string  `gorm:"column:tahun"`
	Triwulan                  string  `gorm:"column:triwulan"`
	IdKpi                     string  `gorm:"column:id_kpi"`
	Kpi                       string  `gorm:"column:kpi"`
	Rumus                     string  `gorm:"column:rumus"`
	Otomatis                  string  `gorm:"column:otomatis"`
	Bobot                     float64 `gorm:"column:bobot"`
	Capping                   string  `gorm:"column:capping"`
	TargetTriwulan            string  `gorm:"column:target_triwulan"`
	TargetKuantitatifTriwulan float64 `gorm:"column:target_kuantitatif_triwulan"`
	TargetTahunan             string  `gorm:"column:target_tahunan"`
	TargetKuantitatifTahunan  float64 `gorm:"column:target_kuantitatif_tahunan"`
	DeskripsiGlossary         string  `gorm:"column:deskripsi_glossary"`
	ItemQualifier             string  `gorm:"column:item_qualifier"`
	DeskripsiQualifier        string  `gorm:"column:deskripsi_qualifier"`
	TargetQualifier           string  `gorm:"column:target_qualifier"`
	IdKeteranganProject       string  `gorm:"column:id_keterangan_project"`
	IdQualifier               string  `gorm:"column:id_qualifier"`
}
