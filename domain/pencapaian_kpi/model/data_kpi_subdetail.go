package model

type DataKpiSubDetail struct {
	IdPengajuan                      string  `gorm:"column:id_pengajuan"`
	IdDetail                         string  `gorm:"column:id_detail"`
	IdSubDetail                      string  `gorm:"column:id_sub_detail"`
	Tahun                            string  `gorm:"column:tahun"`
	Triwulan                         string  `gorm:"column:triwulan"`
	IdKpi                            string  `gorm:"column:id_kpi"`
	Kpi                              string  `gorm:"column:kpi"`
	Rumus                            string  `gorm:"column:rumus"`
	IdPolarisasi                     string  `gorm:"column:id_polarisasi"`
	Polarisasi                       string  `gorm:"column:polarisasi"`
	Otomatis                         string  `gorm:"column:otomatis"`
	Bobot                            float64 `gorm:"column:bobot"`
	Capping                          string  `gorm:"column:capping"`
	TargetTriwulan                   string  `gorm:"column:target_triwulan"`
	TargetKuantitatifTriwulan        float64 `gorm:"column:target_kuantitatif_triwulan"`
	TargetTahunan                    string  `gorm:"column:target_tahunan"`
	TargetKuantitatifTahunan         float64 `gorm:"column:target_kuantitatif_tahunan"`
	Realisasi                        string  `gorm:"column:realisasi"`
	RealisasiKuantitatif             float64 `gorm:"column:realisasi_kuantitatif"`
	RealisasiKeterangan              string  `gorm:"column:realisasi_keterangan"`
	RealisasiValidated               string  `gorm:"column:realisasi_validated"`
	RealisasiKuantitatifValidated    float64 `gorm:"column:realisasi_kuantitatif_validated"`
	IdSumber                         string  `gorm:"column:id_sumber"`
	Sumber                           string  `gorm:"column:sumber"`
	ValidasiKeterangan               string  `gorm:"column:validasi_keterangan"`
	Pencapaian                       float64 `gorm:"column:pencapaian"`
	Skor                             float64 `gorm:"column:skor"`
	DeskripsiGlossary                string  `gorm:"column:deskripsi_glossary"`
	DeskripsiQualifier               string  `gorm:"column:deskripsi_qualifier"`
	TargetQualifier                  string  `gorm:"column:target_qualifier"`
	IdKeteranganProject              string  `gorm:"column:id_keterangan_project"`
	KeteranganProject                string  `gorm:"column:keterangan_project"`
	IdQualifier                      string  `gorm:"column:id_qualifier"`
	RealisasiQualifier               string  `gorm:"column:realisasi_qualifier"`
	RealisasiKuantitatifQualifier    string  `gorm:"column:realisasi_kuantitatif_qualifier"`
	PencapaianQualifierValidated     float64 `gorm:"column:pencapaian_qualifier_validated"`
	PencapaianPostQualifierValidated float64 `gorm:"column:pencapaian_post_qualifier_validated"`
	ItemQualifier                    string  `gorm:"column:item_qualifier"`
}
