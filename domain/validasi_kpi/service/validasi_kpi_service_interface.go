package service

import (
	dto "permen_api/domain/validasi_kpi/dto"
	repo "permen_api/domain/validasi_kpi/repo"
)

type (
	ValidasiKpiServiceInterface interface {
		// =============================================================================
		// INPUT
		// =============================================================================

		// InputValidasiKpi digunakan oleh endpoint POST /validasi-kpi/input.
		InputValidasiKpi(
			req *dto.InputValidasiKpiRequest,
		) (data dto.InputValidasiKpiResponse, err error)

		// =============================================================================
		// APPROVAL
		// =============================================================================

		// ApproveValidasiKpi digunakan oleh endpoint POST /validasi-kpi/approve.
		ApproveValidasiKpi(
			req *dto.ApproveValidasiKpiRequest,
		) (data dto.ApproveValidasiKpiResponse, err error)

		// RejectValidasiKpi digunakan oleh endpoint POST /validasi-kpi/reject.
		RejectValidasiKpi(
			req *dto.RejectValidasiKpiRequest,
		) (data dto.RejectValidasiKpiResponse, err error)

		// =============================================================================
		// GET ALL
		// =============================================================================

		// GetAllValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all.
		GetAllValidasiKpi(
			req *dto.GetAllValidasiKpiRequest,
		) (data []*dto.GetAllValidasiKpiResponse, total int64, err error)

		// GetAllApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
		GetAllApprovalValidasiKpi(
			req *dto.GetAllApprovalValidasiKpiRequest,
		) (data []*dto.GetAllApprovalValidasiKpiResponse, total int64, err error)

		// GetAllTolakanValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
		GetAllTolakanValidasiKpi(
			req *dto.GetAllTolakanValidasiKpiRequest,
		) (data []*dto.GetAllTolakanValidasiKpiResponse, total int64, err error)

		// GetAllDaftarValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-validasi.
		GetAllDaftarValidasiKpi(
			req *dto.GetAllDaftarPValidasiKpiRequest,
		) (data []*dto.GetAllDaftarValidasiKpiResponse, total int64, err error)

		// GetAllDaftarApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalValidasiKpi(
			req *dto.GetAllDaftarApprovalValidasiKpiRequest,
		) (data []*dto.GetAllDaftarApprovalValidasiKpiResponse, total int64, err error)

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
