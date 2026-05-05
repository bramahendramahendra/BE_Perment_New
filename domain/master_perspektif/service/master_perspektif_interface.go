package service

import (
	dto "permen_api/domain/master_perspektif/dto"
	repo "permen_api/domain/master_perspektif/repo"
)

type (
	MasterPerspektifServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterPerspektif digunakan oleh endpoint POST /master-perspektif/get-all.
		GetAllMasterPerspektif() (data []dto.MasterPerspektifResponse, err error)
	}

	masterPerspektifService struct {
		repo repo.MasterPerspektifRepoInterface
	}
)

func NewMasterPerspektifService(repo repo.MasterPerspektifRepoInterface) *masterPerspektifService {
	return &masterPerspektifService{repo: repo}
}
