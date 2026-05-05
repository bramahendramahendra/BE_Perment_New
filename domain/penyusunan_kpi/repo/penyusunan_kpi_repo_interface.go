package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"
	model "permen_api/domain/penyusunan_kpi/model"

	"gorm.io/gorm"
)

type (
	PenyusunanKpiRepoInterface interface {

		// ValidatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
			idLama string,
		) (string, error)

		// CreatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) error

		// RevisionPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/revision.
		RevisionPenyusunanKpi(
			req *dto.RevisionPenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
		) error

		// ApprovePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/approve.
		ApprovePenyusunanKpi(
			idPengajuan, approvalList, approvalPosisi, user string,
		) error

		// RejectPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/reject.
		RejectPenyusunanKpi(
			idPengajuan, approvalList, catatan, user string,
		) error

		// GetAllApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-approval.
		GetAllApprovalPenyusunanKpi(
			req *dto.GetAllApprovalPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllTolakanPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-tolakan.
		GetAllTolakanPenyusunanKpi(
			req *dto.GetAllTolakanPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
		GetAllDaftarPenyusunanKpi(
			req *dto.GetAllDaftarPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalPenyusunanKpi(
			req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetDetailPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (*model.DataKpi, error)

		// GetKpiExportData digunakan oleh endpoint POST /penyusunan-kpi/get-excel dan /penyusunan-kpi/get-pdf.
		GetKpiExportData(idPengajuan, kostl, tahun, triwulan string) (*dto.KpiExportData, error)

		// =============================================================================
		// Digunakan oleh service ApprovePenyusunanKpi dan RejectPenyusunanKpi
		// =============================================================================

		// GetApprovalListJSON digunakan oleh service ApprovePenyusunanKpi dan RejectPenyusunanKpi untuk mengambil daftar approval dalam format JSON.
		GetApprovalListJSON(idPengajuan, userID string) (string, error)

		// GetCatatanTolakan digunakan oleh service RejectPenyusunanKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
		GetCatatanTolakan(idPengajuan string) (string, error)

		// CheckApprovalPenyusunanExists digunakan oleh service ApprovePenyusunanKpi dan RejectPenyusunanKpi untuk memvalidasi keberadaan approval.
		CheckApprovalPenyusunanExists(user, idPengajuan string) (bool, error)

		// =============================================================================
		// Get Exist
		// =============================================================================

		// GetExistDataKpi digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
		GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error)

		// GetExistDataKpiStatus digunakan oleh service untuk mengecek keberadaan data KPI dan mengembalikan id_pengajuan beserta statusnya.
		GetExistDataKpiStatus(tahun, triwulan, kostl string) (idPengajuan string, status int, found bool, err error)

		// =============================================================================
		// Helpers Service
		// =============================================================================

		// LookupKpiMaster digunakan oleh service ValidatePenyusunanKpi untuk mencari id_kpi, kpi, dan rumus dari mst_kpi.
		LookupKpiMaster(kpiText string) (idKpi, kpiFromDB, rumus string, err error)

		// LookupPolarisasi digunakan oleh service ValidatePenyusunanKpi untuk mencari id_polarisasi dari mst_polarisasi.
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
