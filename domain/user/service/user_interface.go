package service

import (
	dto "permen_api/domain/user/dto"
	repo "permen_api/domain/user/repo"
)

type (
	UserServiceInterface interface {
		GetAllUser(req *dto.GetAllUserRequest) (data []dto.UserResponse, err error)
	}

	userService struct {
		repo repo.UserRepoInterface
	}
)

func NewUserService(repo repo.UserRepoInterface) *userService {
	return &userService{repo: repo}
}
