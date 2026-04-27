package dto

type (
	GetRealisasiRequest struct {
		Tahun    string `json:"Tahun" validate:"required"`
		Triwulan string `json:"Triwulan" validate:"required"`
		IdKpi    string `json:"Id_kpi" validate:"required"`
	}

	GetRealisasiResponse struct {
		Data interface{} `json:"data"`
	}
)
