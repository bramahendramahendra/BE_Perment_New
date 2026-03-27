package service

import (
	dto "permen_api/domain/master_method/dto"
	repo "permen_api/domain/master_method/repo"
)

type (
	MasterMethodServiceInterface interface {
		GetAllMasterMethod(req *dto.GetAllMasterMethodRequest) (data []dto.MasterMethodResponse, err error)
	}

	masterMethodService struct {
		repo repo.MasterMethodRepoInterface
	}
)

func NewMasterMethodService(repo repo.MasterMethodRepoInterface) *masterMethodService {
	return &masterMethodService{repo: repo}
}
