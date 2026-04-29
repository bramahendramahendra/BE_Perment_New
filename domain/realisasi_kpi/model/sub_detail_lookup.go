package model

type SubDetailLookup struct {
	IdSubDetail               string  `gorm:"column:id_sub_detail"`
	IdDetail                  string  `gorm:"column:id_detail"`
	Rumus                     string  `gorm:"column:rumus"`
	Otomatis                  string  `gorm:"column:otomatis"`
	Glossary                  string  `gorm:"column:glossary"`
	TargetTriwulan            string  `gorm:"column:target_triwulan"`
	TargetKuantitatifTriwulan float64 `gorm:"column:target_kuantitatif_triwulan"`
	TargetTahunan             string  `gorm:"column:target_tahunan"`
	TargetKuantitatifTahunan  float64 `gorm:"column:target_kuantitatif_tahunan"`
	IdQualifier               string  `gorm:"column:id_qualifier"`
	Qualifier                 string  `gorm:"column:qualifier"`
	DeskripsiQualifier        string  `gorm:"column:deskripsi_qualifier"`
}
