package repo

import (
	dto "permen_api/domain/master_result/dto"
	model "permen_api/domain/master_result/model"

	"gorm.io/gorm"
)

type (
	MasterResultRepoInterface interface {
		GetAllMasterResult(req *dto.GetAllMasterResultRequest) ([]*model.MstResult, error)
		CheckTriwulanExists(idTriwulan string) (bool, error)
		GetDB() *gorm.DB
	}

	masterResultRepo struct {
		db *gorm.DB
	}
)

func NewMasterResultRepo(db *gorm.DB) *masterResultRepo {
	return &masterResultRepo{db: db}
}

func (r *masterResultRepo) GetDB() *gorm.DB {
	return r.db
}
