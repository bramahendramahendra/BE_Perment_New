package service

import (
	dto "permen_api/domain/master_divisi/dto"
	repo "permen_api/domain/master_divisi/repo"
)

type (
	MasterDivisiServiceInterface interface {
		GetAllMasterDivisi() (data []dto.MasterDivisiResponse, err error)
	}

	masterDivisiService struct {
		repo repo.MasterDivisiRepoInterface
	}
)

func NewMasterDivisiService(repo repo.MasterDivisiRepoInterface) *masterDivisiService {
	return &masterDivisiService{repo: repo}
}
