package repo

import (
	dto "permen_api/domain/master_method/dto"
	model "permen_api/domain/master_method/model"

	"gorm.io/gorm"
)

type (
	MasterMethodRepoInterface interface {
		GetAllMasterMethod(req *dto.GetAllMasterMethodRequest) ([]*model.MstMethod, error)
		CheckTriwulanExists(idTriwulan string) (bool, error)
		GetDB() *gorm.DB
	}

	masterMethodRepo struct {
		db *gorm.DB
	}
)

func NewMasterMethodRepo(db *gorm.DB) *masterMethodRepo {
	return &masterMethodRepo{db: db}
}

func (r *masterMethodRepo) GetDB() *gorm.DB {
	return r.db
}
