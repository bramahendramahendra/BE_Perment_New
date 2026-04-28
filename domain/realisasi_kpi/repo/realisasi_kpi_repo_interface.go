package repo

import (
	dto "permen_api/domain/realisasi_kpi/dto"
	model "permen_api/domain/realisasi_kpi/model"

	"gorm.io/gorm"
)

type (
	RealisasiKpiRepoInterface interface {
		// CheckExistRealisasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan input realisasi (2, 4, 80, 81).
		CheckExistRealisasi(idPengajuan string) (bool, error)

		// CheckStatusCreateRealisasi memeriksa apakah id_pengajuan ada dengan status draft realisasi (80).
		CheckStatusCreateRealisasi(idPengajuan string) (bool, error)

		// CheckStatusRevisiRealisasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan revisi (4 atau 80).
		CheckStatusRevisiRealisasi(idPengajuan string) (bool, error)

		// CheckApprovalRealisasiExists memeriksa apakah user adalah approval_posisi aktif untuk id_pengajuan (status 3).
		CheckApprovalRealisasiExists(user, idPengajuan string) (bool, error)

		// GetTriwulanByIdPengajuan mengambil nilai triwulan dari data_kpi berdasarkan id_pengajuan.
		GetTriwulanByIdPengajuan(idPengajuan string) (string, error)

		// GetKpiHeaderByIdPengajuan mengambil field header (tahun, triwulan, kostl, kostl_tx)
		// dari data_kpi berdasarkan id_pengajuan. Digunakan untuk membangun response validate/revision.
		GetKpiHeaderByIdPengajuan(idPengajuan string) (tahun, triwulan, kostl, kostlTx string, err error)

		// LookupSubDetailByKpiSubKpi mencari id_sub_detail, id_detail, target_kuantitatif_triwulan,
		// rumus (id_polarisasi), dan id_qualifier berdasarkan id_pengajuan + kpi_name + sub_kpi_name dari Excel.
		LookupSubDetailByKpiSubKpi(
			idPengajuan, kpiName, subKpiName string,
		) (idSubDetail, idDetail, rumus, idQualifier string, targetKuantitatifTriwulan float64, err error)

		// ValidateRealisasiKpi menyimpan data realisasi ke data_kpi_subdetail (status 80 = draft realisasi).
		// Juga meng-update data_challenge_detail dan data_method_detail jika ada extended data (TW2/TW4).
		ValidateRealisasiKpi(
			req *dto.ValidateRealisasiKpiRequest,
			kpiRows []dto.RealisasiKpiRow,
			kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
		) error

		// Digunakan oleh endpoint POST /realisasi-kpi/create.
		// CreateRealisasiKpi mengubah status dari 80 → 3 (pending approval realisasi) dan menyimpan approval chain.
		CreateRealisasiKpi(
			req *dto.CreateRealisasiKpiRequest,
		) error

		// Digunakan oleh endpoint POST /realisasi-kpi/revision.
		// RevisionRealisasiKpi meng-update ulang data realisasi di DB.
		// Mengizinkan update dari status 80 (draft) atau 4 (ditolak).
		RevisionRealisasiKpi(
			req *dto.RevisionRealisasiKpiRequest,
			kpiRows []dto.RealisasiKpiRow,
			kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
		) error

		// Digunakan oleh endpoint POST /realisasi-kpi/approve.
		ApproveRealisasiKpi(idPengajuan, approvalList, approvalPosisi, user string) error

		// Digunakan oleh endpoint POST /realisasi-kpi/reject.
		RejectRealisasiKpi(idPengajuan, approvalList, catatan, user string) error

		// Digunakan oleh endpoint POST /realisasi-kpi/approve dan /reject.
		// Mengambil approval_list JSON untuk id_pengajuan jika user adalah approval_posisi aktif.
		GetApprovalListJSON(idPengajuan, userID string) (string, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-all.
		GetAllRealisasiKpi(
			req *dto.GetAllRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-all-approval.
		// GetAllApprovalRealisasiKpi mengembalikan list pengajuan berstatus 3 (pending approval realisasi) yang approval_posisi-nya adalah user yang sedang login.
		GetAllApprovalRealisasiKpi(
			req *dto.GetAllApprovalRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-all-tolakan.
		// GetAllTolakanRealisasiKpi mengembalikan list pengajuan berstatus 4 (realisasi ditolak) milik entry_user_realisasi tertentu.
		GetAllTolakanRealisasiKpi(
			req *dto.GetAllTolakanRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-penyusunan.
		// GetAllDaftarRealisasiKpi mengembalikan semua pengajuan dalam konteks realisasi (status 2, 3, 4, 5, 80) dengan filter opsional.
		GetAllDaftarRealisasiKpi(
			req *dto.GetAllDaftarRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-approval.
		// GetAllDaftarApprovalRealisasiKpi mengembalikan semua pengajuan realisasi yang pernah melewati approval (status 3 atau 5) dengan filter opsional.
		GetAllDaftarApprovalRealisasiKpi(
			req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /realisasi-kpi/get-detail.
		GetDetailRealisasiKpi(
			req *dto.GetDetailRealisasiKpiRequest,
		) (*model.DataKpi, error)

		// Digunakan oleh service RevisionPenyusunanKpi untuk mengambil header dari DB.
		GetKpiHeader(idPengajuan string) (tahun, triwulan, kostl, kostlTx, entryUser, entryName string, status int, statusDesc string, err error)

		// Digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
		GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error)

		// GetLinkFormats mengambil semua url_prefix yang aktif dari mst_link_format.
		// Digunakan untuk memvalidasi kolom "Link Dokumen Sumber" pada Excel upload.
		GetLinkFormats() ([]string, error)

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
