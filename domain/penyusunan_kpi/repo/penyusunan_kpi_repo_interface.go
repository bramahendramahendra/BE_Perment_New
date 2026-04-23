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

		// CheckExistPenyusunan mengecek apakah data KPI sudah ada untuk tahun/triwulan/kostl.
		CheckExistPenyusunan(tahun, triwulan, kostl string) (bool, error)

		// GetExistPenyusunanStatus mengecek keberadaan data KPI dan mengembalikan id_pengajuan + status.
		// Mengembalikan found=false jika tidak ada data.
		GetExistPenyusunanStatus(tahun, triwulan, kostl string) (idPengajuan string, status int, found bool, err error)

		// CheckExistIdPengajuan mengecek apakah id_pengajuan, kostl, tahun, dan triwulan cocok di DB.
		CheckExistIdPengajuan(idPengajuan, kostl, tahun, triwulan string) (bool, error)

		// CheckApprovalExists mengecek apakah user adalah approval_posisi aktif (status=0) untuk id_pengajuan.
		CheckApprovalExists(user, idPengajuan string) (bool, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/validate.
		// idLama diisi jika ada draft (status=70) yang harus di-replace; kosong string jika insert baru.
		ValidatePenyusunanKpi(
			req *dto.ValidatePenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
			idLama string,
		) (string, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/create.
		CreatePenyusunanKpi(
			req *dto.CreatePenyusunanKpiRequest,
		) error

		// Digunakan oleh endpoint POST /penyusunan-kpi/revision.
		RevisionPenyusunanKpi(
			req *dto.RevisionPenyusunanKpiRequest,
			kpiRows []dto.PenyusunanKpiRow,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
			resultList []dto.DataResult,
			processList []dto.DataProcess,
			contextList []dto.DataContext,
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
		GetDetailPenyusunanKpi(
			req *dto.GetDetailPenyusunanKpiRequest,
		) (*model.DataKpi, error)

		// Digunakan oleh endpoint POST /penyusunan-kpi/get-excel dan /get-pdf.
		GetKpiExportData(idPengajuan, kostl, tahun, triwulan string) (*dto.KpiExportData, error)

		// Digunakan oleh service RevisionPenyusunanKpi untuk mengambil header dari DB.
		GetKpiHeader(idPengajuan string) (tahun, triwulan, kostl, kostlTx, entryUser, entryName string, status int, statusDesc string, err error)

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
