package service

import (
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		// ValidatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/validate.
		// REQUEST hanya memerlukan Divisi, Tahun, dan Triwulan.
		// KPI, ChallengeList, dan MethodList diambil dari file Excel dan tabel mst_kpi.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidatePenyusunanKpiResponse, err error)

		// CreatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) (data dto.CreatePenyusunanKpiResponse, err error)

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
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (data *dto.GetAllDataPenyusunanKpiResponse, err error)

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
