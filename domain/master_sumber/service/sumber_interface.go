package service

import (
	dto "permen_api/domain/master_sumber/dto"
	repo "permen_api/domain/master_sumber/repo"
)

type (
	MasterSumberServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterSumber digunakan oleh endpoint POST /master-sumber/get-all.
		GetAllMasterSumber() (data []dto.MasterSumberResponse, err error)
	}

	masterSumberService struct {
		repo repo.MasterSumberRepoInterface
	}
)

func NewMasterSumberService(repo repo.MasterSumberRepoInterface) *masterSumberService {
	return &masterSumberService{repo: repo}
}
