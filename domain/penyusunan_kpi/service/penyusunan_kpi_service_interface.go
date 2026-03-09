package service

import (
	"mime/multipart"
	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			files []*multipart.FileHeader,
		) (*dto.InsertPenyusunanKpiResult, error)
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
