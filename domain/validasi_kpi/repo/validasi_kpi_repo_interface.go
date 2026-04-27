package repo

import (
	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"

	"gorm.io/gorm"
)

type (
	ValidasiKpiRepoInterface interface {
		// CheckExistInputValidasi memeriksa apakah id_pengajuan ada dengan status yang mengizinkan input validasi (5, 7, 90, 91).
		CheckExistInputValidasi(idPengajuan string) (bool, error)

		// CheckExistApprovalValidasi memeriksa apakah id_pengajuan ada dengan status 6 (pending approval validasi).
		CheckExistApprovalValidasi(idPengajuan string) (bool, error)

		// CheckExistBatalValidasi memeriksa apakah id_pengajuan ada di tabel data_kpi.
		CheckExistBatalValidasi(idPengajuan string) (bool, error)

		// GetKostlTxByIdPengajuan mengambil kostl_tx dari data_kpi berdasarkan id_pengajuan.
		GetKostlTxByIdPengajuan(idPengajuan string) (string, error)

		// GetEntryUserValidasiByIdPengajuan mengambil entry_user_validasi untuk keperluan notifikasi tolakan.
		GetEntryUserValidasiByIdPengajuan(idPengajuan string) (string, error)

		// InputValidasi menyimpan data validasi ke data_kpi dan data_kpi_subdetail (status → 6).
		InputValidasi(req *dto.InputValidasiRequest) error

		// ApprovalValidasi memproses approve atau reject validasi KPI.
		ApprovalValidasi(req *dto.ApprovalValidasiRequest) error

		// ValidasiBatal membatalkan proses validasi (status → 91) dan menghapus notifikasi terkait.
		ValidasiBatal(req *dto.ValidasiBatalRequest) error

		// ApproveValidasi memproses approve validasi KPI (status → 8 atau chain).
		ApproveValidasi(req *dto.ApproveValidasiRequest) error

		// RejectValidasi memproses reject validasi KPI (status → 7).
		RejectValidasi(req *dto.RejectValidasiRequest) error

		// GetAllApprovalValidasi mengambil list pengajuan validasi yang menunggu approval user (status=6, approval_posisi=user).
		GetAllApprovalValidasi(req *dto.GetAllApprovalValidasiRequest) ([]*model.DataKpi, int64, error)

		// GetAllTolakanValidasi mengambil list pengajuan validasi yang ditolak (status=7).
		GetAllTolakanValidasi(req *dto.GetAllTolakanValidasiRequest) ([]*model.DataKpi, int64, error)

		// GetAllDaftarPenyusunanValidasi mengambil semua pengajuan dalam konteks validasi dengan filter opsional.
		GetAllDaftarPenyusunanValidasi(req *dto.GetAllDaftarPenyusunanValidasiRequest) ([]*model.DataKpi, int64, error)

		// GetAllDaftarApprovalValidasi mengambil pengajuan validasi yang pernah melibatkan user dalam approval.
		GetAllDaftarApprovalValidasi(req *dto.GetAllDaftarApprovalValidasiRequest) ([]*model.DataKpi, int64, error)

		// GetAllValidasi mengambil semua pengajuan validasi tanpa filter user (mirip get-all realisasi).
		GetAllValidasi(req *dto.GetAllValidasiRequest) ([]*model.DataKpi, int64, error)
	}

	validasiKpiRepo struct {
		db *gorm.DB
	}
)

func NewValidasiKpiRepo(db *gorm.DB) *validasiKpiRepo {
	return &validasiKpiRepo{db: db}
}
