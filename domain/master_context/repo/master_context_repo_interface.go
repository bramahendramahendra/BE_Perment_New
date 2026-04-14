package repo

import (
	dto "permen_api/domain/master_context/dto"
	model "permen_api/domain/master_context/model"

	"gorm.io/gorm"
)

type (
	MasterContextRepoInterface interface {
		GetAllMasterContext(req *dto.GetAllMasterContextRequest) ([]*model.MstChallenge, error)
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
