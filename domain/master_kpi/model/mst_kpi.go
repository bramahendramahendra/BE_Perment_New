package model

type (
	MstKpi struct {
		IdKpi int    `gorm:"column:id_kpi"`
		Kpi   string `gorm:"column:kpi"`
		Rumus string `gorm:"column:rumus"`
	}
)
