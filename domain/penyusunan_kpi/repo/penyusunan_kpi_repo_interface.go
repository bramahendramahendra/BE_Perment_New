package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"
	model "permen_api/domain/penyusunan_kpi/model"

	"gorm.io/gorm"
)

type (
	PenyusunanKpiRepoInterface interface {
		// LookupKpiMaster mencari id_kpi, kpi, dan rumus dari mst_kpi.
		LookupKpiMaster(kpiText string) (idKpi, kpiFromDB, rumus string, err error)

		// LookupPolarisasi mencari id_polarisasi dari mst_polarisasi.
		LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.PenyusunanResult,
			processList []dto.PenyusunanProcess,
			contextList []dto.PenyusunanContext,
		) (string, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/revision.
		RevisionPenyusunanKpi(
			req *dto.RevisionPenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.PenyusunanResult,
			processList []dto.PenyusunanProcess,
			contextList []dto.PenyusunanContext,
		) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/approve.
		ApprovePenyusunanKpi(idPengajuan, approvalList, approvalPosisi, user string) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/reject.
		RejectPenyusunanKpi(idPengajuan, approvalList, catatan, user string) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/approve dan /reject.
		// Mengambil approval_list JSON untuk id_pengajuan jika user adalah approval_posisi aktif.
		GetApprovalListJSON(idPengajuan, userID string) (string, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-approval.
		GetAllApprovalPenyusunanKpi(
			req *dto.GetAllApprovalPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-tolakan.
		GetAllTolakanPenyusunanKpi(
			req *dto.GetAllTolakanPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
		GetAllDaftarPenyusunanKpi(
			req *dto.GetAllDaftarPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-approval.
		GetAllDaftarApprovalPenyusunanKpi(
			req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
		) ([]*model.DataKpi, int64, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
		// Mengembalikan GetDetailPenyusunanKpiResponse dengan approval_list sudah di-unmarshal
		// menjadi []Approval, dan nested KPI/result/process/context sudah terisi.
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (*dto.GetDetailPenyusunanKpiResponse, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-excel dan /get-pdf.
		GetKpiExportData(idPengajuan string) (*dto.KpiExportData, error)

		GetDB() *gorm.DB
	}

	penyusunanKpiRepo struct {
		db *gorm.DB
	}
)

func NewPenyusunanKpiRepo(db *gorm.DB) *penyusunanKpiRepo {
	return &penyusunanKpiRepo{db: db}
}

func (r *penyusunanKpiRepo) GetDB() *gorm.DB {
	return r.db
}
