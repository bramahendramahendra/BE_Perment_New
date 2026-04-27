package dto

type (
	GetRealisasiRequest struct {
		Tahun    string `json:"tahun" validate:"required"`
		Triwulan string `json:"triwulan" validate:"required"`
		IdKpi    string `json:"id_kpi" validate:"required"`
	}

	GetRealisasiResponse struct {
		Data interface{} `json:"data"`
	}
)
