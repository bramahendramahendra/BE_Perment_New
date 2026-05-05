package repo

import (
	model "permen_api/domain/master_kpi/model"

	"gorm.io/gorm"
)

type (
	MasterKpiRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterKpi digunakan oleh endpoint POST /master-kpi/get-all.
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
