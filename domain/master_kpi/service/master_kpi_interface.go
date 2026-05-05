package service

import (
	dto "permen_api/domain/master_kpi/dto"
	repo "permen_api/domain/master_kpi/repo"
)

type (
	MasterKpiServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterKpi digunakan oleh endpoint POST /master-kpi/get-all.
		GetAllMasterKpi() (data []dto.MasterKpiResponse, err error)
	}

	masterKpiService struct {
		repo repo.MasterKpiRepoInterface
	}
)

func NewMasterKpiService(repo repo.MasterKpiRepoInterface) *masterKpiService {
	return &masterKpiService{repo: repo}
}
