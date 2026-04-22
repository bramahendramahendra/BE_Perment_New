package model

type DataKpiDetail struct {
	IdPengajuan         string `gorm:"column:id_pengajuan"`
	IdDetail            string `gorm:"column:id_detail"`
	Tahun               string `gorm:"column:tahun"`
	Triwulan            string `gorm:"column:triwulan"`
	IdKpi               string `gorm:"column:id_kpi"`
	Kpi                 string `gorm:"column:kpi"`
	Rumus               string `gorm:"column:rumus"`
	IdPersfektif        string `gorm:"column:id_perspektif"`
	Perspektif          string `gorm:"column:perspektif"`
	IdKeteranganProject string `gorm:"column:id_keterangan_project"`
	KeteranganProject   string `gorm:"column:keterangan_project"`

	TotalSubKpi  int              `gorm:"-"`
	KpiSubDetail []DataKpiSubDetail `gorm:"-"`
}
