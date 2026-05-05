package repo

import (
	model "permen_api/domain/master_triwulan/model"

	"gorm.io/gorm"
)

type (
	MasterTriwulanRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterTriwulan digunakan oleh endpoint POST /master-triwulan/get-all.
		GetAllMasterTriwulan() ([]*model.MstTriwulan, error)
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
