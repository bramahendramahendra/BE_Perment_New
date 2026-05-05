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
		// Service
		// =============================================================================

		// Digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi
		GetApprovalListJSON(idPengajuan, userID string) (string, error)
		GetCatatanTolakan(idPengajuan string) (string, error)

		// CheckApprovalRealisasiExists memeriksa apakah user adalah approval_posisi aktif untuk id_pengajuan (status 3).
		CheckApprovalRealisasiExists(user, idPengajuan string) (bool, error)

		// =============================================================================
		// GET EXIST
		// =============================================================================
		// Digunakan oleh service RevisionPenyusunanKpi untuk mengambil header dari DB.
		GetKpiHeader(idPengajuan string) (tahun, triwulan, kostl, kostlTx, entryUser, entryName string, status int, statusDesc string, err error)

		// Digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
		GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error)

		// GetLinkFormats mengambil semua url_prefix yang aktif dari mst_link_format.
		// Digunakan untuk memvalidasi kolom "Link Dokumen Sumber" pada Excel upload.
		GetLinkFormats() ([]string, error)

		// CheckExistRealisasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan input realisasi (2, 4, 80, 81).
		CheckExistRealisasi(idPengajuan string) (bool, error)

		// CheckStatusCreateRealisasi memeriksa apakah id_pengajuan ada dengan status draft realisasi (80).
		CheckStatusCreateRealisasi(idPengajuan string) (bool, error)

		// CheckStatusRevisiRealisasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan revisi (4 atau 80).
		CheckStatusRevisiRealisasi(idPengajuan string) (bool, error)

		// GetTriwulanByIdPengajuan mengambil nilai triwulan dari data_kpi berdasarkan id_pengajuan.
		GetTriwulanByIdPengajuan(idPengajuan string) (string, error)

		// GetKpiHeaderByIdPengajuan mengambil field header (tahun, triwulan, kostl, kostl_tx)
		// dari data_kpi berdasarkan id_pengajuan. Digunakan untuk membangun response validate/revision.
		GetKpiHeaderByIdPengajuan(idPengajuan string) (tahun, triwulan, kostl, kostlTx string, err error)

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
