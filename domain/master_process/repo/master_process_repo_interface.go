package repo

import (
	dto "permen_api/domain/master_process/dto"
	model "permen_api/domain/master_process/model"

	"gorm.io/gorm"
)

type (
	MasterProcessRepoInterface interface {
		GetAllMasterProcess(req *dto.GetAllMasterProcessRequest) ([]*model.MstMethod, error)
		CheckTriwulanExists(idTriwulan string) (bool, error)
		GetDB() *gorm.DB
	}

	masterProcessRepo struct {
		db *gorm.DB
	}
)

func NewMasterProcessRepo(db *gorm.DB) *masterProcessRepo {
	return &masterProcessRepo{db: db}
}

func (r *masterProcessRepo) GetDB() *gorm.DB {
	return r.db
}
