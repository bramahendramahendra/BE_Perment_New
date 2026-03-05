package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"

	"gorm.io/gorm"
)

// =============================================
// INTERFACE & STRUCT DEFINITION
// =============================================

type (
	PenyusunanKpiRepoInterface interface {
		// InsertPenyusunanKpi melakukan insert ke semua tabel terkait dalam 1 transaksi DB:
		//   - data_kpi
		//   - data_kpi_detail
		//   - data_kpi_subdetail  (berasal dari hasil parse Excel, per index KPI)
		//   - data_challenge_detail
		//   - data_method_detail
		//
		// Parameter kpiSubDetails adalah map dengan key = index KPI (sesuai urutan req.Kpi),
		// value = slice baris hasil parse Excel untuk KPI tersebut.
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
		) error

		// GetDB mengembalikan instance *gorm.DB untuk keperluan lain (logging, dsb)
		GetDB() *gorm.DB
	}

	penyusunanKpiRepo struct {
		db *gorm.DB
	}
)

// NewPenyusunanKpiRepo membuat instance baru penyusunanKpiRepo
func NewPenyusunanKpiRepo(db *gorm.DB) *penyusunanKpiRepo {
	return &penyusunanKpiRepo{db: db}
}

// GetDB mengembalikan instance *gorm.DB
func (r *penyusunanKpiRepo) GetDB() *gorm.DB {
	return r.db
}
