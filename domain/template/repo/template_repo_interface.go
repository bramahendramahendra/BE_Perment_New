package repo

import (
	model "permen_api/domain/template/model"

	"gorm.io/gorm"
)

type (
	TemplateRepoInterface interface {
		// GetKpiWithPolarisasi mengambil semua data mst_kpi beserta polarisasi-nya.
		// Jika rumus pada mst_kpi tidak ditemukan di mst_polarisasi, kolom polarisasi dikosongkan.
		GetKpiWithPolarisasi() ([]*model.MstKpiPolarisasi, error)
		GetDB() *gorm.DB
	}

	templateRepo struct {
		db *gorm.DB
	}
)

func NewTemplateRepo(db *gorm.DB) *templateRepo {
	return &templateRepo{db: db}
}

func (r *templateRepo) GetDB() *gorm.DB {
	return r.db
}
