package service

import (
	dto "permen_api/domain/pencapaian_kpi/dto"
	repo "permen_api/domain/pencapaian_kpi/repo"
)

type (
	PencapaianKpiServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-all-pencapaian.
		GetAllPencapaianKpi(
			req *dto.GetAllPencapaianKpiRequest,
		) (data []*dto.GetAllPencapaianKpiResponse, total int64, err error)

		// =============================================================================
		// GET DETAIL
		// =============================================================================

		// GetDetailPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-detail.
		GetDetailPencapaianKpi(
			req *dto.GetDetailPencapaianKpiRequest,
		) (data *dto.GetDetailPencapaianKpiResponse, err error)

		// =============================================================================
		// DOWNLOAD
		// =============================================================================

		// GetExcelPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-excel.
		GetExcelPencapaianKpi(req *dto.GetExcelPencapaianKpiRequest) ([]byte, string, error)

		// GetPdfPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-pdf.
		GetPdfPencapaianKpi(req *dto.GetPdfPencapaianKpiRequest) ([]byte, string, error)
	}

	pencapaianKpiService struct {
		repo repo.PencapaianKpiRepoInterface
	}
)

func NewPencapaianKpiService(repo repo.PencapaianKpiRepoInterface) *pencapaianKpiService {
	return &pencapaianKpiService{repo: repo}
}
