package service

import (
	dto "permen_api/domain/template/dto"
	repo "permen_api/domain/template/repo"
)

type (
	TemplateServiceInterface interface {
		// GenerateFormatPenyusunanKpi digunakan oleh endpoint POST /template/format-penyusunan-kpi.
		// Menghasilkan file Excel template penyusunan KPI sesuai triwulan (tanpa isi data baris).
		// Sheet 1 — nama sheet mengikuti nilai triwulan dari request (TW1, TW2, TW3, TW4):
		//   Jika triwulan TW1/TW3 → kolom A–O (format base).
		//   Jika triwulan TW2/TW4 → kolom A–U (format extended).
		// Sheet 2 — nama sheet "KPI":
		//   Kolom A (KPI) dan B (Polarisasi) dari join mst_kpi dan mst_polarisasi.
		GenerateFormatPenyusunanKpi(req *dto.FormatPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)

		// GenerateTolakanPenyusunanKpi digunakan oleh endpoint POST /template/tolakan-penyusunan-kpi.
		// Menghasilkan file Excel yang sudah terisi data baris sub KPI berdasarkan id_pengajuan.
		// Format kolom mengikuti triwulan dari DB:
		//   TW1/TW3 → kolom A–O (format base).
		//   TW2/TW4 → kolom A–U (format extended, kolom P–U terisi data result/method/challenge).
		// Sheet 2 — nama sheet "KPI":
		//   Kolom A (KPI) dan B (Polarisasi) dari join mst_kpi dan mst_polarisasi.
		GenerateTolakanPenyusunanKpi(req *dto.TolakanPenyusunanKpiRequest) (fileBytes []byte, filename string, err error)
	}

	templateService struct {
		repo repo.TemplateRepoInterface
	}
)

func NewTemplateService(repo repo.TemplateRepoInterface) *templateService {
	return &templateService{repo: repo}
}
