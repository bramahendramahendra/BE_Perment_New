package repo

import (
	dto "permen_api/domain/realisasi_kpi/dto"

	"gorm.io/gorm"
)

type (
	RealisasiKpiRepoInterface interface {
		// LookupSubDetailByKpiSubKpi mencari id_sub_detail, id_detail, target_kuantitatif_triwulan,
		// dan rumus (id_polarisasi) berdasarkan id_pengajuan + kpi_name + sub_kpi_name dari Excel.
		LookupSubDetailByKpiSubKpi(
			idPengajuan, kpiName, subKpiName string,
		) (idSubDetail, idDetail, rumus string, targetKuantitatifTriwulan float64, err error)

		// ValidateRealisasiKpi menyimpan data realisasi ke data_kpi_subdetail (status 80 = draft realisasi).
		// Juga meng-update data_challenge_detail dan data_method_detail jika ada extended data.
		ValidateRealisasiKpi(
			req *dto.ValidateRealisasiKpiRequest,
			rows []dto.RealisasiKpiRow,
		) error

		// RevisionRealisasiKpi meng-update ulang data realisasi di DB.
		// Mengizinkan update dari status 80 (draft) atau 4 (ditolak).
		RevisionRealisasiKpi(
			req *dto.RevisionRealisasiKpiRequest,
			rows []dto.RealisasiKpiRow,
		) error

		// CreateRealisasiKpi mengubah status dari 80 → 3 (pending approval realisasi)
		// dan menyimpan approval chain.
		CreateRealisasiKpi(
			req *dto.CreateRealisasiKpiRequest,
		) error

		// ApprovalRealisasiKpi memproses approve/reject realisasi.
		ApprovalRealisasiKpi(
			req *dto.ApprovalRealisasiKpiRequest,
		) error

		// GetAllApprovalRealisasiKpi mengembalikan list pengajuan berstatus 3 (pending approval realisasi)
		// yang approval_posisi-nya adalah user yang sedang login.
		GetAllApprovalRealisasiKpi(
			req *dto.GetAllApprovalRealisasiKpiRequest,
		) ([]*dto.DataKpiRealisasi, int64, error)

		// GetAllTolakanRealisasiKpi mengembalikan list pengajuan berstatus 4 (realisasi ditolak)
		// milik entry_user_realisasi tertentu.
		GetAllTolakanRealisasiKpi(
			req *dto.GetAllTolakanRealisasiKpiRequest,
		) ([]*dto.DataKpiRealisasi, int64, error)

		// GetAllDaftarRealisasiKpi mengembalikan semua pengajuan dalam konteks realisasi
		// (status 2, 3, 4, 5, 80) dengan filter opsional.
		GetAllDaftarRealisasiKpi(
			req *dto.GetAllDaftarRealisasiKpiRequest,
		) ([]*dto.DataKpiRealisasi, int64, error)

		// GetAllDaftarApprovalRealisasiKpi mengembalikan semua pengajuan realisasi
		// yang pernah melewati approval (status 3 atau 5) dengan filter opsional.
		GetAllDaftarApprovalRealisasiKpi(
			req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
		) ([]*dto.DataKpiRealisasi, int64, error)

		// GetDetailRealisasiKpi mengembalikan detail lengkap satu pengajuan beserta
		// sub KPI, context list, dan process list.
		GetDetailRealisasiKpi(
			req *dto.GetDetailRealisasiKpiRequest,
		) (*dto.GetDetailRealisasiKpiResponse, error)

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
