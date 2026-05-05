package service

import (
	dto "permen_api/domain/master_context/dto"
	repo "permen_api/domain/master_context/repo"
)

type (
	MasterContextServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterContext digunakan oleh endpoint POST /master-context/get-all.
		GetAllMasterContext(req *dto.GetAllMasterContextRequest) (data []dto.MasterContextResponse, err error)
	}

	masterContextService struct {
		repo repo.MasterContextRepoInterface
	}
)

func NewMasterContextService(repo repo.MasterContextRepoInterface) *masterContextService {
	return &masterContextService{repo: repo}
}
