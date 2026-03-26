package repo

import (
	model "permen_api/domain/master_divisi/model"

	"gorm.io/gorm"
)

type (
	MasterDivisiRepoInterface interface {
		GetAllMasterDivisi() ([]*model.MstDivisi, error)
		GetDB() *gorm.DB
	}

	masterDivisiRepo struct {
		db *gorm.DB
	}
)

func NewMasterDivisiRepo(db *gorm.DB) *masterDivisiRepo {
	return &masterDivisiRepo{db: db}
}

func (r *masterDivisiRepo) GetDB() *gorm.DB {
	return r.db
}
