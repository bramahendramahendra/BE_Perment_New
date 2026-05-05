package repo

import (
	dto "permen_api/domain/realisasi_kpi/dto"
	model "permen_api/domain/realisasi_kpi/model"

	"gorm.io/gorm"
)

type (
	RealisasiKpiRepoInterface interface {
		// ValidateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/validate.
		ValidateRealisasiKpi(
			req *dto.ValidateRealisasiKpiRequest,
			kpiRows []dto.RealisasiKpiRow,
			kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
		) error

		// CreateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/create.
		CreateRealisasiKpi(
			req *dto.CreateRealisasiKpiRequest,
		) error

		// RevisionRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/revision.
		RevisionRealisasiKpi(
			req *dto.RevisionRealisasiKpiRequest,
			kpiRows []dto.RealisasiKpiRow,
			kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
		) error

		// ApprovePenyusunanKpi digunakan oleh endpoint POST /realisasi-kpi/approve.
		ApproveRealisasiKpi(
			idPengajuan, approvalList, approvalPosisi, user string,
		) error

		// RejectPenyusunanKpi digunakan oleh endpoint POST /realisasi-kpi/reject.
		RejectRealisasiKpi(
			idPengajuan, approvalList, catatan, user string,
		) error

		// GetAllRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all.
		GetAllRealisasiKpi(
			req *dto.GetAllRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-approval.
		GetAllApprovalRealisasiKpi(
			req *dto.GetAllApprovalRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllTolakanRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-tolakan.
		GetAllTolakanRealisasiKpi(
			req *dto.GetAllTolakanRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-realisasi.
		GetAllDaftarRealisasiKpi(
			req *dto.GetAllDaftarRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalRealisasiKpi(
			req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetDetailRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-detail.
		GetDetailRealisasiKpi(
			req *dto.GetDetailRealisasiKpiRequest,
		) (*model.DataKpi, error)

		// =============================================================================
		// Approval Helper
		// =============================================================================

		// Digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi

		// GetApprovalListJSON digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi untuk mengambil daftar approval dalam format JSON.
		GetApprovalListJSON(idPengajuan, userID string) (string, error)

		// GetCatatanTolakan digunakan oleh service RejectRealisasiKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
		GetCatatanTolakan(idPengajuan string) (string, error)

		// CheckApprovalRealisasiExists digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi untuk memvalidasi keberadaan approval.
		CheckApprovalRealisasiExists(user, idPengajuan string) (bool, error)

		// =============================================================================
		// GET EXIST
		// =============================================================================

		// GetExistDataKpi digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
		GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error)

		// =============================================================================
		// Service Helpers
		// =============================================================================

		GetLinkFormats() ([]string, error)

		// CheckExistRealisasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan input realisasi (2, 4, 80, 81).
		CheckExistRealisasi(idPengajuan string) (bool, error)

		// LookupSubDetailByKpiSubKpi mencari data sub detail berdasarkan id_pengajuan + kpi_name + sub_kpi_name dari Excel.
		LookupSubDetailByKpiSubKpi(
			idPengajuan, kpiName, subKpiName string,
		) (*model.SubDetailLookup, error)

		GetDB() *gorm.DB
	}

	realisasiKpiRepo struct {
		db *gorm.DB
	}
)

func NewRealisasiKpiRepo(db *gorm.DB) *realisasiKpiRepo {
	return &realisasiKpiRepo{db: db}
}

func (r *realisasiKpiRepo) GetDB() *gorm.DB {
	return r.db
}
