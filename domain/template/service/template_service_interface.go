package service

import (
	dto "permen_api/domain/template/dto"
)

type (
	TemplateServiceInterface interface {
		// Digunakan oleh endpoint POST /template/format-penyusunan-kpi.
		// Menghasilkan file Excel template penyusunan KPI sesuai triwulan.
		// Nama sheet mengikuti nilai triwulan dari request (TW1, TW2, TW3, TW4).
		// Jika triwulan TW1/TW3 → kolom A–O (format base).
		// Jika triwulan TW2/TW4 → kolom A–U (format extended).
		GenerateFormatPenyusunanKpi(req *dto.FormatPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)
	}

	templateService struct{}
)

func NewTemplateService() *templateService {
	return &templateService{}
}
