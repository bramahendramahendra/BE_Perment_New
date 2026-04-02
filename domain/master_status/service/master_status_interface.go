package service

import (
	dto "permen_api/domain/master_status/dto"
	repo "permen_api/domain/master_status/repo"
)

type (
	MasterStatusServiceInterface interface {
		GetAllMasterStatus() (data []dto.MasterStatusResponse, err error)
		GetDraftMasterStatus() (data []dto.MasterStatusResponse, err error)
	}

	masterStatusService struct {
		repo repo.MasterStatusRepoInterface
	}
)

func NewMasterStatusService(repo repo.MasterStatusRepoInterface) *masterStatusService {
	return &masterStatusService{repo: repo}
}
