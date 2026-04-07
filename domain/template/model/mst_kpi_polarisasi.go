package model

type (
	MstKpiPolarisasi struct {
		Kpi        string `gorm:"column:kpi"`
		Polarisasi string `gorm:"column:polarisasi"`
	}
)
