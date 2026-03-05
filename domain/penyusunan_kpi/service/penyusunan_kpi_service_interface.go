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
		//  1. Validasi mandatory field
		//  2. Validasi jumlah file sesuai jumlah KPI
		//  3. Parse & validasi semua file Excel (semua harus valid sebelum insert DB)
		//  4. Panggil repo untuk insert dalam 1 transaksi
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			files []*multipart.FileHeader,
		) error
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

// NewPenyusunanKpiService membuat instance baru penyusunanKpiService
func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
