package service

import (
	dto "permen_api/domain/master_tahun/dto"
	repo "permen_api/domain/master_tahun/repo"
)

type (
	MasterTahunServiceInterface interface {
		GetAllMasterTahun() (data []dto.MasterTahunResponse, err error)
	}

	masterTahunService struct {
		repo repo.MasterTahunRepoInterface
	}
)

func NewMasterTahunService(repo repo.MasterTahunRepoInterface) *masterTahunService {
	return &masterTahunService{repo: repo}
}
