package repo

import (
	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"

	"gorm.io/gorm"
)

type (
	ValidasiKpiRepoInterface interface {
		// =============================================================================
		// INPUT
		// =============================================================================

		// InputValidasiKpi digunakan oleh endpoint POST /validasi-kpi/input.
		InputValidasiKpi(
			req *dto.InputValidasiKpiRequest,
		) error

		// DraftValidasiKpi digunakan oleh endpoint POST /validasi-kpi/draft.
		DraftValidasiKpi(
			req *dto.DraftValidasiKpiRequest,
		) error

		// =============================================================================
		// APPROVAL
		// =============================================================================

		// ApproveValidasiKpi digunakan oleh endpoint POST /validasi-kpi/approve.
		ApproveValidasiKpi(
			idPengajuan, approvalList, approvalPosisi, user string,
		) error

		// RejectValidasiKpi digunakan oleh endpoint POST /validasi-kpi/reject.
		RejectValidasiKpi(
			idPengajuan, approvalList, catatan, user string,
		) error

		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all.
		GetAllValidasiKpi(
			req *dto.GetAllValidasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
		GetAllApprovalValidasiKpi(
			req *dto.GetAllApprovalValidasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllTolakanValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
		GetAllTolakanValidasiKpi(
			req *dto.GetAllTolakanValidasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-validasi.
		GetAllDaftarValidasiKpi(
			req *dto.GetAllDaftarPValidasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// GetAllDaftarApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalValidasiKpi(
			req *dto.GetAllDaftarApprovalValidasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// =============================================================================
		// GET DETAIL
		// =============================================================================

		// GetDetailRealisasiKpi digunakan oleh endpoint POST /validasi-kpi/get-detail.
		GetDetailValidasiKpi(
			req *dto.GetDetailValidasiKpiRequest,
		) (*model.DataKpi, error)

		// =============================================================================
		// APPROVAL HELPER
		// =============================================================================

		// GetApprovalListJSON digunakan oleh service ApproveValidasiKpi dan RejectValidasiKpi untuk mengambil daftar approval dalam format JSON.
		GetApprovalListJSON(idPengajuan, userID string) (string, error)

		// GetCatatanTolakan digunakan oleh service RejectValidasiKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
		GetCatatanTolakan(idPengajuan string) (string, error)

		// CheckApprovalValidasiExists digunakan oleh service ApproveValidasiKpi dan RejectValidasiKpi untuk memvalidasi keberadaan approval.
		CheckApprovalValidasiExists(user, idPengajuan string) (bool, error)

		// =============================================================================
		// GET EXIST
		// =============================================================================

		// GetExistDataKpi digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
		GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error)

		// GetIndikatorPencapaian mengambil semua indikator warna dari tabel indikator_pencapaian.
		GetIndikatorPencapaian() ([]*model.IndikatorPencapaian, error)
	}

	validasiKpiRepo struct {
		db *gorm.DB
	}
)

func NewValidasiKpiRepo(db *gorm.DB) *validasiKpiRepo {
	return &validasiKpiRepo{db: db}
}
