package repo

import (
	dto "permen_api/domain/master_context/dto"
	model "permen_api/domain/master_context/model"

	"gorm.io/gorm"
)

type (
	MasterContextRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterContext digunakan oleh endpoint POST /master-context/get-all.
		GetAllMasterContext(req *dto.GetAllMasterContextRequest) ([]*model.MstChallenge, error)

		// =============================================================================
		// CHECK EXIST
		// =============================================================================

		// CheckTriwulanExists digunakan oleh service untuk mengecek keberadaan data Triwulan.
		CheckTriwulanExists(idTriwulan string) (bool, error)

		GetDB() *gorm.DB
	}

	masterContextRepo struct {
		db *gorm.DB
	}
)

func NewMasterContextRepo(db *gorm.DB) *masterContextRepo {
	return &masterContextRepo{db: db}
}

func (r *masterContextRepo) GetDB() *gorm.DB {
	return r.db
}
