package model

type DataKpiSubDetailValidasi struct {
	IdSubDetail                      string
	IdKpi                            string
	Kpi                              string
	Bobot                            float64
	TargetTriwulan                   string
	TargetKuantitatifTriwulan        float64
	TargetQualifier                  string
	RealisasiValidated               string
	RealisasiKuantitatifValidated    float64
	ValidasiKeterangan               string
	Pencapaian                       float64
	Skor                             float64
	PencapaianQualifierValidated     float64
	PencapaianPostQualifierValidated float64
}
