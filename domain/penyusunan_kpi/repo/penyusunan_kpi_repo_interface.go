package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"

	"gorm.io/gorm"
)

type (
	PenyusunanKpiRepoInterface interface {
		LookupSubKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error)

		LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
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
		// Mengembalikan 1 record lengkap (header + KpiDetail + ChallengeDetail + MethodDetail)
		// berdasarkan id_pengajuan, tanpa filter status maupun entry_user.
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (*dto.GetAllDraftPenyusunanKpiResponse, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-csv dan /get-pdf.
		// Mengambil data header (nama divisi, tahun, triwulan) dan baris sub KPI
		// dari data_kpi_subdetail untuk keperluan ekspor dokumen.
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
