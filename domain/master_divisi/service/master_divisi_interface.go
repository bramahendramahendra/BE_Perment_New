package service

import (
	dto "permen_api/domain/master_divisi/dto"
	repo "permen_api/domain/master_divisi/repo"
)

type (
	MasterDivisiServiceInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterDivisi digunakan oleh endpoint POST /master-divisi/get-all.
		GetAllMasterDivisi() (data []dto.MasterDivisiResponse, err error)
	}

	masterDivisiService struct {
		repo repo.MasterDivisiRepoInterface
	}
)

func NewMasterDivisiService(repo repo.MasterDivisiRepoInterface) *masterDivisiService {
	return &masterDivisiService{repo: repo}
}
