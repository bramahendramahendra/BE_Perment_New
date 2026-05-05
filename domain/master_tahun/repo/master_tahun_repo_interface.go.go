package repo

import (
	model "permen_api/domain/master_tahun/model"

	"gorm.io/gorm"
)

type (
	MasterTahunRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterTahun digunakan oleh endpoint POST /master-tahun/get-all.
		GetAllMasterTahun() (*model.MstTahun, error)
		GetDB() *gorm.DB
	}

	masterTahunRepo struct {
		db *gorm.DB
	}
)

func NewMasterTahunRepo(db *gorm.DB) *masterTahunRepo {
	return &masterTahunRepo{db: db}
}

func (r *masterTahunRepo) GetDB() *gorm.DB {
	return r.db
}
