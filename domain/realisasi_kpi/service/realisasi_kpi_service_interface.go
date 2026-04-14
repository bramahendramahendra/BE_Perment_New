package service

import (
	"mime/multipart"

	dto "permen_api/domain/realisasi_kpi/dto"
	repo "permen_api/domain/realisasi_kpi/repo"
)

type (
	RealisasiKpiServiceInterface interface {
		// ValidateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/validate.
		// Menerima file Excel realisasi, mem-parse, menghitung Pencapaian/Skor,
		// dan menyimpan ke DB sebagai draft realisasi (status 80).
		ValidateRealisasiKpi(
			req *dto.ValidateRealisasiKpiRequest,
			file *multipart.FileHeader,
		) (data dto.ValidateRealisasiKpiResponse, err error)

		// RevisionRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/revision.
		// Menerima file Excel realisasi baru, menghitung ulang, dan meng-update DB.
		RevisionRealisasiKpi(
			req *dto.RevisionRealisasiKpiRequest,
			file *multipart.FileHeader,
		) (data dto.RevisionRealisasiKpiResponse, err error)

		// CreateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/create.
		// Mengubah status dari draft (80) ke pending approval (3) dan menyimpan approval chain.
		CreateRealisasiKpi(
			req *dto.CreateRealisasiKpiRequest,
		) (data dto.CreateRealisasiKpiResponse, err error)

		// ApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/approval.
		// Memproses approve (chain atau final) atau reject realisasi.
		ApprovalRealisasiKpi(
			req *dto.ApprovalRealisasiKpiRequest,
		) (data dto.ApprovalRealisasiKpiResponse, err error)

		// GetAllApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-approval.
		// Mengembalikan list pengajuan berstatus 3 yang menunggu approval dari user tertentu.
		GetAllApprovalRealisasiKpi(
			req *dto.GetAllApprovalRealisasiKpiRequest,
		) (data []*dto.GetAllApprovalRealisasiKpiResponse, total int64, err error)

		// GetAllTolakanRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-tolakan.
		// Mengembalikan list pengajuan berstatus 4 (realisasi ditolak) milik user tertentu.
		GetAllTolakanRealisasiKpi(
			req *dto.GetAllTolakanRealisasiKpiRequest,
		) (data []*dto.GetAllTolakanRealisasiKpiResponse, total int64, err error)

		// GetAllDaftarRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-realisasi.
		// Mengembalikan semua pengajuan dalam konteks realisasi dengan filter opsional.
		GetAllDaftarRealisasiKpi(
			req *dto.GetAllDaftarRealisasiKpiRequest,
		) (data []*dto.GetAllDaftarRealisasiKpiResponse, total int64, err error)

		// GetAllDaftarApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-approval.
		// Mengembalikan semua pengajuan realisasi yang sudah masuk proses approval.
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
