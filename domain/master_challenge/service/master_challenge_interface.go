package service

import (
	dto "permen_api/domain/master_challenge/dto"
	repo "permen_api/domain/master_challenge/repo"
)

type (
	MasterChallengeServiceInterface interface {
		GetAllMasterChallenge(req *dto.GetAllMasterChallengeRequest) (data []dto.MasterChallengeResponse, err error)
	}

	masterChallengeService struct {
		repo repo.MasterChallengeRepoInterface
	}
)

func NewMasterChallengeService(repo repo.MasterChallengeRepoInterface) *masterChallengeService {
	return &masterChallengeService{repo: repo}
}
