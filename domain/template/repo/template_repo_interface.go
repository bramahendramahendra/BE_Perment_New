package repo

import (
	model "permen_api/domain/template/model"

	"gorm.io/gorm"
)

type (
	TemplateRepoInterface interface {
		// GetKpiWithPolarisasi mengambil semua data mst_kpi beserta polarisasi-nya.
		// Jika rumus pada mst_kpi tidak ditemukan di mst_polarisasi, kolom polarisasi dikosongkan.
		GetKpiWithPolarisasi() ([]*model.MstKpiPolarisasi, error)

		// CheckDataExist mengecek keberadaan data tolakan KPI
		// berdasarkan id_pengajuan, kostl, tahun, dan triwulan.
		CheckDataExist(idPengajuan, kostl, tahun, triwulan string) (bool, error)

		// GetPenyusunanKpiData mengambil seluruh baris sub KPI dari DB
		// berdasarkan id_pengajuan, untuk keperluan generate Excel tolakan penyusunan KPI.
		GetPenyusunanKpiData(idPengajuan string) (*model.ExcelData, error)

		// GetExistPenyusunanStatus mengecek apakah sudah ada record di data_kpi
		// berdasarkan tahun, triwulan, dan kostl.
		// Mengembalikan status jika ditemukan, dan found=false jika tidak ada.
		GetExistPenyusunanStatus(tahun, triwulan, kostl string) (status int, found bool, err error)

		// CheckRevisiRealisasiExist mengecek apakah pengajuan memiliki data realisasi
		// yang dapat direvisi (status 4=tolak realisasi atau 80=draft realisasi).
		CheckRevisiRealisasiExist(idPengajuan, kostl, tahun, triwulan string) (bool, error)

		// GetRealisasiKpiData mengambil seluruh baris realisasi KPI dari DB
		// berdasarkan id_pengajuan, termasuk data realisasi yang sudah diisi user sebelumnya.
		GetRealisasiKpiData(idPengajuan string) (*model.RealisasiExcelData, error)

		GetDB() *gorm.DB
	}

	templateRepo struct {
		db *gorm.DB
	}
)

func NewTemplateRepo(db *gorm.DB) *templateRepo {
	return &templateRepo{db: db}
}

func (r *templateRepo) GetDB() *gorm.DB {
	return r.db
}
