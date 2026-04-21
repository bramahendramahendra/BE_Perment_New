package service

import (
	"mime/multipart"

	dto "permen_api/domain/realisasi_kpi/dto"
	repo "permen_api/domain/realisasi_kpi/repo"
)

type (
	RealisasiKpiServiceInterface interface {
		// ValidateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/validate.
		ValidateRealisasiKpi(
			req *dto.ValidateRealisasiKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidateRealisasiKpiResponse, err error)

		// CreateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/create.
		// Mengubah status dari draft (80) ke pending approval (3) dan menyimpan approval chain.
		CreateRealisasiKpi(
			req *dto.CreateRealisasiKpiRequest,
		) (data dto.CreateRealisasiKpiResponse, err error)

		// RevisionRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/revision.
		RevisionRealisasiKpi(
			req *dto.RevisionRealisasiKpiRequest,
			file *multipart.FileHeader,
		) (data dto.RevisionRealisasiKpiResponse, err error)

		// ApprovePenyusunanKpi digunakan oleh endpoint POST /realisasi-kpi/approve.
		ApproveRealisasiKpi(
			req *dto.ApproveRealisasiKpiRequest,
		) (data dto.ApproveRealisasiKpiResponse, err error)

		// RejectPenyusunanKpi digunakan oleh endpoint POST /realisasi-kpi/reject.
		RejectRealisasiKpi(
			req *dto.RejectRealisasiKpiRequest,
		) (data dto.RejectRealisasiKpiResponse, err error)

		// GetAllRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all.
		GetAllRealisasiKpi(
			req *dto.GetAllRealisasiKpiRequest,
		) (data []*dto.GetAllRealisasiKpiResponse, total int64, err error)

		// GetAllApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-approval.
		GetAllApprovalRealisasiKpi(
			req *dto.GetAllApprovalRealisasiKpiRequest,
		) (data []*dto.GetAllApprovalRealisasiKpiResponse, total int64, err error)

		// GetAllTolakanRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-tolakan.
		GetAllTolakanRealisasiKpi(
			req *dto.GetAllTolakanRealisasiKpiRequest,
		) (data []*dto.GetAllTolakanRealisasiKpiResponse, total int64, err error)

		// GetAllDaftarRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-realisasi.
		GetAllDaftarRealisasiKpi(
			req *dto.GetAllDaftarRealisasiKpiRequest,
		) (data []*dto.GetAllDaftarRealisasiKpiResponse, total int64, err error)

		// GetAllDaftarApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalRealisasiKpi(
			req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
		) (data []*dto.GetAllDaftarApprovalRealisasiKpiResponse, total int64, err error)

		// GetDetailRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-detail.
		// Mengembalikan detail lengkap satu pengajuan beserta sub KPI, context, dan process list.
		GetDetailRealisasiKpi(
			req *dto.GetDetailRealisasiKpiRequest,
		) (data *dto.GetDetailRealisasiKpiResponse, err error)
	}

	realisasiKpiService struct {
		repo repo.RealisasiKpiRepoInterface
	}
)

func NewRealisasiKpiService(repo repo.RealisasiKpiRepoInterface) *realisasiKpiService {
	return &realisasiKpiService{repo: repo}
}
