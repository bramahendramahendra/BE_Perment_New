package service

import (
	dto "permen_api/domain/sample/dto"
	repo "permen_api/domain/sample/repo"
)

type (
	UserIntegrationServiceInterface interface {
		GetAllUserIntegrations() (data []dto.UserIntegrationResponse, err error)
	}

	userIntegrationService struct {
		repo repo.UserIntegrationRepoInterface
	}
)

func NewUserIntegrationService(repo repo.UserIntegrationRepoInterface) *userIntegrationService {
	return &userIntegrationService{repo: repo}
}
