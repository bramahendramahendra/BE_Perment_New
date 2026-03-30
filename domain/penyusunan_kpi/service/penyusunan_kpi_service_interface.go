package service

import (
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidatePenyusunanKpiResponse, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) (data dto.CreatePenyusunanKpiResponse, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-draft.
		GetAllDraftPenyusunanKpi(
			req *dto.GetAllDraftPenyusunanKpiRequest,
		) (data []*dto.GetAllDraftPenyusunanKpiResponse, total int64, err error)
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
