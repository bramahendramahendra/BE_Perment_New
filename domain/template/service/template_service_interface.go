package service

import (
	dto "permen_api/domain/template/dto"
)

type (
	TemplateServiceInterface interface {
		// Digunakan oleh endpoint GET /template/format-penyusunan-kpi.
		// Menghasilkan file Excel template penyusunan KPI sesuai triwulan.
		// Jika triwulan TW1/TW2/TW3 → sheet "Selain TW 4" (kolom A–O).
		// Jika triwulan TW4         → sheet "TW4"          (kolom A–U).
		GenerateFormatPenyusunanKpi(req *dto.FormatPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)
	}

	templateService struct{}
)

func NewTemplateService() *templateService {
	return &templateService{}
}
