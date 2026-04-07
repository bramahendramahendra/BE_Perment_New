package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"

	"gorm.io/gorm"
)

type (
	PenyusunanKpiRepoInterface interface {
		// LookupKpiMaster mencari id_kpi, kpi, dan rumus dari mst_kpi.
		LookupKpiMaster(kpiText string) (idKpi, kpiFromDB, rumus string, err error)

		// LookupPolarisasi mencari id_polarisasi dari mst_polarisasi.
		LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error)
		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.PenyusunanResult,
			methodList []dto.PenyusunanMethod,
			challengeList []dto.PenyusunanChallenge,
		) (string, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-draft.
		GetAllDraftPenyusunanKpi(
			req *dto.GetAllDraftPenyusunanKpiRequest,
		) ([]*dto.GetAllDraftPenyusunanKpiResponse, int64, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (*dto.GetAllDraftPenyusunanKpiResponse, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-csv dan /get-pdf.
		GetKpiExportData(idPengajuan string) (*dto.KpiExportData, error)

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
