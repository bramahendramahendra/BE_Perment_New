package repo

import (
	model "permen_api/domain/master_status/model"

	"gorm.io/gorm"
)

type (
	MasterStatusRepoInterface interface {
		GetAllMasterStatus() ([]*model.MstStatus, error)
		GetDraftMasterStatus() ([]*model.MstStatus, error)
		GetDB() *gorm.DB
	}

	masterStatusRepo struct {
		db *gorm.DB
	}
)

func NewMasterStatusRepo(db *gorm.DB) *masterStatusRepo {
	return &masterStatusRepo{db: db}
}

func (r *masterStatusRepo) GetDB() *gorm.DB {
	return r.db
}
