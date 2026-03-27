package repo

import (
	dto "permen_api/domain/user/dto"
	model "permen_api/domain/user/model"

	"gorm.io/gorm"
)

type (
	UserRepoInterface interface {
		GetAllUser(req *dto.GetAllUserRequest) ([]*model.User, error)
		GetDB() *gorm.DB
	}

	userRepo struct {
		db *gorm.DB
	}
)

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetDB() *gorm.DB {
	return r.db
}
