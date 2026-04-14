package service

import (
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		// ValidatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidatePenyusunanKpiResponse, err error)

		// RevisionPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/revision.
		RevisionPenyusunanKpi(
			req *dto.RevisionPenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.RevisionPenyusunanKpiResponse, err error)

		// CreatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) (data dto.CreatePenyusunanKpiResponse, err error)

		// ApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/approval.
		ApprovalPenyusunanKpi(
			req *dto.ApprovalPenyusunanKpiRequest,
		) (data dto.ApprovalPenyusunanKpiResponse, err error)

		// GetAllApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-approval.
		GetAllApprovalPenyusunanKpi(
			req *dto.GetAllApprovalPenyusunanKpiRequest,
		) (data []*dto.GetAllApprovalPenyusunanKpiResponse, total int64, err error)

		// GetAllTolakanPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-tolakan.
		GetAllTolakanPenyusunanKpi(
			req *dto.GetAllTolakanPenyusunanKpiRequest,
		) (data []*dto.GetAllTolakanPenyusunanKpiResponse, total int64, err error)

		// GetAllDaftarPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
		GetAllDaftarPenyusunanKpi(
			req *dto.GetAllDaftarPenyusunanKpiRequest,
		) (data []*dto.GetAllDaftarPenyusunanKpiResponse, total int64, err error)

		// GetAllDaftarApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalPenyusunanKpi(
			req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
		) (data []*dto.GetAllDaftarApprovalPenyusunanKpiResponse, total int64, err error)

		// GetDetailPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
		// Mengembalikan GetDetailPenyusunanKpiResponse dengan struktur response baru
		// (nested divisi, entry, approvalList sebagai array, totalKpi, totalResult, dst).
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (data *dto.GetDetailPenyusunanKpiResponse, err error)

		// GetExcelPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-excel.
		GetExcelPenyusunanKpi(
			req *dto.GetExcelPenyusunanKpiRequest,
		) (fileBytes []byte, filename string, err error)

		// GetPdfPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-pdf.
		GetPdfPenyusunanKpi(
			req *dto.GetPdfPenyusunanKpiRequest,
		) (fileBytes []byte, filename string, err error)
	}

	penyusunanKpiService struct {
		repo repo.PenyusunanKpiRepoInterface
	}
)

func NewPenyusunanKpiService(repo repo.PenyusunanKpiRepoInterface) *penyusunanKpiService {
	return &penyusunanKpiService{repo: repo}
}
