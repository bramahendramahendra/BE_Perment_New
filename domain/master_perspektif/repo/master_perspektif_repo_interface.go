package repo

import (
	model "permen_api/domain/master_perspektif/model"

	"gorm.io/gorm"
)

type (
	MasterPerspektifRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterPerspektif digunakan oleh endpoint POST /master-perspektif/get-all.
		GetAllMasterPerspektif() ([]*model.MstPerspektif, error)
		GetDB() *gorm.DB
	}

	masterPerspektifRepo struct {
		db *gorm.DB
	}
)

func NewMasterPerspektifRepo(db *gorm.DB) *masterPerspektifRepo {
	return &masterPerspektifRepo{db: db}
}

func (r *masterPerspektifRepo) GetDB() *gorm.DB {
	return r.db
}
