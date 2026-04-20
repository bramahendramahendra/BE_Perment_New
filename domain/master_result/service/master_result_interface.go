package service

import (
	dto "permen_api/domain/master_result/dto"
	repo "permen_api/domain/master_result/repo"
)

type (
	MasterResultServiceInterface interface {
		GetAllMasterResult(req *dto.GetAllMasterResultRequest) (data []dto.MasterResultResponse, err error)
	}

	masterResultService struct {
		repo repo.MasterResultRepoInterface
	}
)

func NewMasterResultService(repo repo.MasterResultRepoInterface) *masterResultService {
	return &masterResultService{repo: repo}
}
