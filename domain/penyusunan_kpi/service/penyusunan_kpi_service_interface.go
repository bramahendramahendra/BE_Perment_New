package service

import (
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.CreatePenyusunanKpiResponse, err error)
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
