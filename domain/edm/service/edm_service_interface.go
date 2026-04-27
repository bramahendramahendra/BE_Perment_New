package service

import (
	dto "permen_api/domain/edm/dto"

	"gorm.io/gorm"
)

type (
	EdmServiceInterface interface {
		GetRealisasi(req *dto.GetRealisasiRequest) (data interface{}, err error)
	}

	edmService struct {
		db *gorm.DB
	}
)

func NewEdmService(db *gorm.DB) *edmService {
	return &edmService{db: db}
}
