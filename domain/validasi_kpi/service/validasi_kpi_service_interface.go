package service

import (
	dto "permen_api/domain/validasi_kpi/dto"
	repo "permen_api/domain/validasi_kpi/repo"
)

type (
	ValidasiKpiServiceInterface interface {
		// InputValidasi digunakan oleh endpoint POST /validasi-kpi/input.
		// Menyimpan data validasi KPI dan mengirim notifikasi ke approver (status → 6).
		InputValidasi(req *dto.InputValidasiRequest) (data dto.InputValidasiResponse, err error)

		// ApprovalValidasi digunakan oleh endpoint POST /validasi-kpi/approval.
		// Memproses approve (status → 8 atau chain) atau reject (status → 7).
		ApprovalValidasi(req *dto.ApprovalValidasiRequest) (data dto.ApprovalValidasiResponse, err error)

		// ValidasiBatal digunakan oleh endpoint POST /validasi-kpi/batal.
		// Membatalkan proses validasi (status → 91) dan menghapus notifikasi terkait.
		ValidasiBatal(req *dto.ValidasiBatalRequest) (data dto.ValidasiBatalResponse, err error)

		// ApproveValidasi digunakan oleh endpoint POST /validasi-kpi/approve.
		ApproveValidasi(req *dto.ApproveValidasiRequest) (data dto.ApproveValidasiResponse, err error)

		// RejectValidasi digunakan oleh endpoint POST /validasi-kpi/reject.
		RejectValidasi(req *dto.RejectValidasiRequest) (data dto.RejectValidasiResponse, err error)

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
	}

	validasiKpiService struct {
		repo repo.ValidasiKpiRepoInterface
	}
)

func NewValidasiKpiService(repo repo.ValidasiKpiRepoInterface) *validasiKpiService {
	return &validasiKpiService{repo: repo}
}
