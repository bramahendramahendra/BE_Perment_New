package repo

import (
	model "permen_api/domain/master_sumber/model"

	"gorm.io/gorm"
)

type (
	MasterSumberRepoInterface interface {
		GetAllMasterSumber() ([]*model.MstSumber, error)
		GetDB() *gorm.DB
	}

	masterSumberRepo struct {
		db *gorm.DB
	}
)

func NewMasterSumberRepo(db *gorm.DB) *masterSumberRepo {
	return &masterSumberRepo{db: db}
}

func (r *masterSumberRepo) GetDB() *gorm.DB {
	return r.db
}
