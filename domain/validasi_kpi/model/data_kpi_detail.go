package model

type DataKpiDetailValidasi struct {
	IdDetail     string
	IdKpi        string
	Kpi          string
	Rumus        string
	TotalSubKpi  int
	KpiSubDetail []DataKpiSubDetailValidasi
}
