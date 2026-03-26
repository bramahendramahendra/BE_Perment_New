package repo

import (
	model "permen_api/domain/sample/model"

	"gorm.io/gorm"
)

type (
	UserIntegrationRepoInterface interface {
		GetAllUserIntegrations() ([]*model.UserIntegration, error)
		GetDB() *gorm.DB
	}

	userIntegrationRepo struct {
		db *gorm.DB
	}
)

func NewUserIntegrationRepo(db *gorm.DB) *userIntegrationRepo {
	return &userIntegrationRepo{db: db}
}

func (r *userIntegrationRepo) GetDB() *gorm.DB {
	return r.db
}
