package service

import (
	dto "permen_api/domain/master_tahun/dto"
)

type (
	MasterTahunServiceInterface interface {
		GetAllMasterTahun() (data []dto.MasterTahunResponse, err error)
	}

	masterTahunService struct{}
)

func NewMasterTahunService() *masterTahunService {
	return &masterTahunService{}
}
