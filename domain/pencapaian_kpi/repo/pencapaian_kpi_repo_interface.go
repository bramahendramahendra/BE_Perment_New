package repo

import (
	dto "permen_api/domain/pencapaian_kpi/dto"
	model "permen_api/domain/pencapaian_kpi/model"

	"gorm.io/gorm"
)

type (
	PencapaianKpiRepoInterface interface {
		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-all-pencapaian.
		GetAllPencapaianKpi(
			req *dto.GetAllPencapaianKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// =============================================================================
		// GET DETAIL
		// =============================================================================

		// GetDetailPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-detail.
		GetDetailPencapaianKpi(
			req *dto.GetDetailPencapaianKpiRequest,
		) (*model.DataKpi, error)

		// =============================================================================
		// GET INDIKATOR
		// =============================================================================

		// GetIndikatorPencapaian mengambil semua indikator warna dari tabel indikator_pencapaian.
		GetIndikatorPencapaian() ([]*model.IndikatorPencapaian, error)
	}

	pencapaianKpiRepo struct {
		db *gorm.DB
	}
)

func NewPencapaianKpiRepo(db *gorm.DB) *pencapaianKpiRepo {
	return &pencapaianKpiRepo{db: db}
}
