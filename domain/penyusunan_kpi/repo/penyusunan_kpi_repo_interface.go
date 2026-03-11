package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"

	"gorm.io/gorm"
)

type (
	PenyusunanKpiRepoInterface interface {
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
		) (string, error)

		LookupSubKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error)

		LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error)

		GetDB() *gorm.DB
	}

	penyusunanKpiRepo struct {
		db *gorm.DB
	}
)

func NewPenyusunanKpiRepo(db *gorm.DB) *penyusunanKpiRepo {
	return &penyusunanKpiRepo{db: db}
}

func (r *penyusunanKpiRepo) GetDB() *gorm.DB {
	return r.db
}
