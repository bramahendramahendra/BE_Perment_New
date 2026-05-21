package service

import (
	dto "permen_api/domain/edm/dto"
	edm "permen_api/pkg/external/edm"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type (
	EdmServiceInterface interface {
		GetKpi(req *dto.GetKpiRequest) (data interface{}, err error)
	}

	edmService struct {
		edm   edm.EdmClient
		redis *redis.Client
	}
)

func NewEdmService(db *gorm.DB, redisClient *redis.Client) *edmService {
	return &edmService{
		edm:   edm.New(db, false),
		redis: redisClient,
	}
}
