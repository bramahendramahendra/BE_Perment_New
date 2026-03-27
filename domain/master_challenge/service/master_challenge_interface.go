package service

import (
	dto "permen_api/domain/sample/dto"
	repo "permen_api/domain/sample/repo"
	transport "permen_api/pkg/transport"
)

type (
	UserIntegrationServiceInterface interface {
		GetUserIntegrationByUsername(username string) (data dto.UserIntegrationResponse, err error)
		GetAllUserIntegrations() (data []dto.UserIntegrationResponse, err error)
	}

	userIntegrationService struct {
		repo repo.UserIntegrationRepoInterface
		esb  *transport.RestClient
	}
)

func NewUserIntegrationService(repo repo.UserIntegrationRepoInterface, esb *transport.RestClient) *userIntegrationService {
	return &userIntegrationService{repo: repo, esb: esb}
}
