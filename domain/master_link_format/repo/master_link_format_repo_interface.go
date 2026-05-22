package repo

import (
	model "permen_api/domain/master_link_format/model"

	"gorm.io/gorm"
)

type (
	MasterLinkFormatRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllMasterLinkFormat digunakan oleh endpoint POST /master-link_format/get-all.
		GetAllMasterLinkFormat() ([]*model.MstLinkFormat, error)
		GetDB() *gorm.DB
	}

	masterLinkFormatRepo struct {
		db *gorm.DB
	}
)

func NewMasterLinkFormatRepo(db *gorm.DB) *masterLinkFormatRepo {
	return &masterLinkFormatRepo{db: db}
}

func (r *masterLinkFormatRepo) GetDB() *gorm.DB {
	return r.db
}
