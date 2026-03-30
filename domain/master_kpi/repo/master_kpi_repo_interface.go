package repo

import (
	model "permen_api/domain/master_kpi/model"

	"gorm.io/gorm"
)

type (
	MasterKpiRepoInterface interface {
		GetAllMasterKpi() ([]*model.MstKpi, error)
		GetDB() *gorm.DB
	}

	masterKpiRepo struct {
		db *gorm.DB
	}
)

func NewMasterKpiRepo(db *gorm.DB) *masterKpiRepo {
	return &masterKpiRepo{db: db}
}

func (r *masterKpiRepo) GetDB() *gorm.DB {
	return r.db
}
