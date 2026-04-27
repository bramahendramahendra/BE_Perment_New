package service

import (
	dto "permen_api/domain/edm/dto"
	edm "permen_api/pkg/external/edm"

	"gorm.io/gorm"
)

type (
	EdmServiceInterface interface {
		GetRealisasi(req *dto.GetRealisasiRequest) (data interface{}, err error)
	}

	edmService struct {
		edm edm.EdmClient
	}
)

func NewEdmService(db *gorm.DB) *edmService {
	return &edmService{
		edm: edm.New(db, false),
	}
}
