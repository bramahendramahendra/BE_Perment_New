package service

import (
	dto "permen_api/domain/master_process/dto"
	repo "permen_api/domain/master_process/repo"
)

type (
	MasterProcessServiceInterface interface {
		GetAllMasterProcess(req *dto.GetAllMasterProcessRequest) (data []dto.MasterProcessResponse, err error)
	}

	masterProcessService struct {
		repo repo.MasterProcessRepoInterface
	}
)

func NewMasterProcessService(repo repo.MasterProcessRepoInterface) *masterProcessService {
	return &masterProcessService{repo: repo}
}
