package service

import (
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	repo "permen_api/domain/penyusunan_kpi/repo"
)

type (
	PenyusunanKpiServiceInterface interface {
		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidatePenyusunanKpiResponse, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) (data dto.CreatePenyusunanKpiResponse, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-draft.
		GetAllDraftPenyusunanKpi(
			req *dto.GetAllDraftPenyusunanKpiRequest,
		) (data []*dto.GetAllDraftPenyusunanKpiResponse, total int64, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
		// Mengembalikan 1 record lengkap (header + KpiDetail + ChallengeDetail + MethodDetail)
		// berdasarkan id_pengajuan.
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (data *dto.GetAllDraftPenyusunanKpiResponse, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-csv.
		// Menghasilkan file CSV berisi daftar sub KPI beserta bobot, target tahunan, dan capping.
		GetCsvPenyusunanKpi(
			req *dto.GetCsvPenyusunanKpiRequest,
		) (fileBytes []byte, filename string, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-pdf.
		// Menghasilkan file PDF berisi tabel KPI bergaya seperti tampilan aplikasi.
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
