package repo

import (
	"permen_api/domain/master_triwulan/model"

	"gorm.io/gorm"
)

type (
	MasterTriwulanRepoInterface interface {
		GetAllTriwulan() ([]*model.MstTriwulan, error)
		GetDB() *gorm.DB
	}

	masterTriwulanRepo struct {
		db *gorm.DB
	}
)

func NewMasterTriwulanRepo(db *gorm.DB) *masterTriwulanRepo {
	return &masterTriwulanRepo{db: db}
}

func (r *masterTriwulanRepo) GetDB() *gorm.DB {
	return r.db
}
