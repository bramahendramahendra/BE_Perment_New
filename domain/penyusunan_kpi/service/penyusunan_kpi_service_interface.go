package service

import (
	"mime/multipart"
	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

// =============================================
// INTERFACE & STRUCT DEFINITION
// =============================================

type (
	PenyusunanKpiServiceInterface interface {
		// InsertPenyusunanKpi memproses request insert KPI:
		//  1. Validasi jumlah file sesuai jumlah KPI
		//  2. Parse & validasi semua file Excel (semua harus valid sebelum insert DB)
		//  3. Panggil repo untuk insert dalam 1 transaksi
		//
		// Return: idPengajuan yang di-generate backend, error
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			files []*multipart.FileHeader,
		) (string, error)
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

// NewPenyusunanKpiService membuat instance baru penyusunanKpiService
func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
