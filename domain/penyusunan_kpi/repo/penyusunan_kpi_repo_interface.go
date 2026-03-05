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
		// ID di-generate di backend mengikuti pola frontend lama:
		//   - IDPengajuan  = Kostl + Tahun + Triwulan + timestamp(ymdhis)
		//   - id_detail    = IDPengajuan + "P" + index KPI 3 digit (P001, P002, ...)
		//   - id_sub_detail = IDPengajuan + "C" + index SubKPI 3 digit (C001, C002, ...)
		//                     ⚠️ index reset setiap KPI baru
		//
		// Return: idPengajuan yang di-generate, error
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
		) (string, error)

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
