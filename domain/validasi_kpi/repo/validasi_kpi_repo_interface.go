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

		// CheckApprovalValidasiExists memvalidasi bahwa user adalah approver aktif untuk pengajuan (status=6, approval_posisi=user).
		CheckApprovalValidasiExists(user, idPengajuan string) (bool, error)

		// GetApprovalListValidasiJSON mengambil approval_list_validasi dalam format JSON string.
		GetApprovalListValidasiJSON(idPengajuan, userID string) (string, error)

		// =============================================================================
		// GET EXIST
		// =============================================================================

		// GetExistDataValidasi mengambil data header KPI berdasarkan id_pengajuan untuk validasi status.
		GetExistDataValidasi(idPengajuan string) (*model.DataKpiExist, error)
	}

	validasiKpiRepo struct {
		db *gorm.DB
	}
)

func NewValidasiKpiRepo(db *gorm.DB) *validasiKpiRepo {
	return &validasiKpiRepo{db: db}
}
