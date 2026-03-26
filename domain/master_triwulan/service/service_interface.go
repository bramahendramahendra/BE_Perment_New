package service

import (
	dto "permen_api/domain/master_triwulan/dto"
	repo "permen_api/domain/master_triwulan/repo"
)

type (
	MasterTriwulanServiceInterface interface {
		GetAllTriwulan() (data []dto.TriwulanResponse, err error)
	}

	masterTriwulanService struct {
		repo repo.MasterTriwulanRepoInterface
	}
)

func NewMasterTriwulanService(repo repo.MasterTriwulanRepoInterface) *masterTriwulanService {
	return &masterTriwulanService{repo: repo}
}
