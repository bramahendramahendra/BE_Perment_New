package service

import (
	dto "permen_api/domain/validasi_kpi/dto"
	repo "permen_api/domain/validasi_kpi/repo"
)

type (
	ValidasiKpiServiceInterface interface {
		// =============================================================================
		// INPUT VALIDASI (validate + create + revision)
		// =============================================================================

		// InputValidasi digunakan oleh endpoint POST /validasi-kpi/input.
		// Menyimpan data validasi KPI, mengatur rantai approval, dan mengirim notifikasi (status → 6).
		// Berlaku untuk status 5 (baru), 7 (revisi), 90/91 (ulang setelah batal).
		InputValidasi(req *dto.InputValidasiRequest) (data dto.InputValidasiResponse, err error)

		// =============================================================================
		// APPROVAL
		// =============================================================================

		// ApproveValidasi digunakan oleh endpoint POST /validasi-kpi/approve.
		// Memproses approve dalam rantai approval; jika approver terakhir → status = 8 (final).
		ApproveValidasi(req *dto.ApproveValidasiRequest) (data dto.ApproveValidasiResponse, err error)

		// RejectValidasi digunakan oleh endpoint POST /validasi-kpi/reject.
		// Memproses penolakan validasi KPI (status → 7) dan mengirim notifikasi ke pengaju.
		RejectValidasi(req *dto.RejectValidasiRequest) (data dto.RejectValidasiResponse, err error)

		// =============================================================================
		// BATAL
		// =============================================================================

		// ValidasiBatal digunakan oleh endpoint POST /validasi-kpi/batal.
		// Membatalkan proses validasi (status → 91) dan menghapus notifikasi terkait.
		ValidasiBatal(req *dto.ValidasiBatalRequest) (data dto.ValidasiBatalResponse, err error)

		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllApprovalValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
		GetAllApprovalValidasi(req *dto.GetAllApprovalValidasiRequest) (data []*dto.GetAllValidasiResponse, total int64, err error)

		// GetAllTolakanValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
		GetAllTolakanValidasi(req *dto.GetAllTolakanValidasiRequest) (data []*dto.GetAllValidasiResponse, total int64, err error)

		// GetAllDaftarPenyusunanValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-penyusunan.
		GetAllDaftarPenyusunanValidasi(req *dto.GetAllDaftarPenyusunanValidasiRequest) (data []*dto.GetAllValidasiResponse, total int64, err error)

		// GetAllDaftarApprovalValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalValidasi(req *dto.GetAllDaftarApprovalValidasiRequest) (data []*dto.GetAllValidasiResponse, total int64, err error)

		// GetAllValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-validasi.
		GetAllValidasi(req *dto.GetAllValidasiRequest) (data []*dto.GetAllValidasiResponse, total int64, err error)

		// =============================================================================
		// GET DETAIL
		// =============================================================================

		// GetDetailValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-detail.
		GetDetailValidasiKpi(req *dto.GetDetailValidasiKpiRequest) (data *dto.GetDetailValidasiKpiResponse, err error)
	}

	validasiKpiService struct {
		repo repo.ValidasiKpiRepoInterface
	}
)

func NewValidasiKpiService(repo repo.ValidasiKpiRepoInterface) *validasiKpiService {
	return &validasiKpiService{repo: repo}
}
