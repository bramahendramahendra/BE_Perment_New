package repo

import (
	dto "permen_api/domain/penyusunan_kpi/dto"

	"gorm.io/gorm"
)

// =============================================
// INTERFACE & STRUCT DEFINITION
// =============================================

type (
	PenyusunanKpiRepoInterface interface {
		// InsertPenyusunanKpi melakukan insert ke semua tabel terkait dalam 1 transaksi DB:
		//   - data_kpi
		//   - data_kpi_detail
		//   - data_kpi_subdetail  (berasal dari hasil parse Excel + hasil lookup master)
		//   - data_challenge_detail
		//   - data_method_detail
		//
		// Catatan: kpiSubDetails yang diterima sudah berisi IdSubKpi dan IdPolarisasi
		// yang telah di-lookup sebelumnya via LookupSubKpiMaster & LookupPolarisasi.
		//
		// ID di-generate di backend mengikuti pola frontend lama:
		//   - IDPengajuan   = Kostl + Tahun + Triwulan + timestamp(ymdhis)
		//   - id_detail     = IDPengajuan + "P" + index KPI 3 digit (P001, P002, ...)
		//   - id_sub_detail = IDPengajuan + "C" + counter global 3 digit (C001, C002, ...)
		//                     ⚠️ counter TIDAK reset antar KPI
		//
		// Return: idPengajuan yang di-generate, error
		InsertPenyusunanKpi(
			req *dto.InsertPenyusunanKpiRequest,
			kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
		) (string, error)

		// LookupSubKpiMaster mencari id_kpi, nama kpi, dan rumus dari tabel mst_kpi
		// secara case-insensitive berdasarkan teks sub KPI dari Excel (kolom C).
		//
		// Behavior:
		//   - Ditemukan     : return id_kpi (string), kpiFromDB (nama dari DB), rumus, nil
		//   - Tidak ditemukan: return "0", subKpiText (nama asli Excel), "", nil
		//     (bukan error — id_kpi = 0 adalah kondisi valid)
		LookupSubKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error)

		// LookupPolarisasi mencari id_polarisasi dari tabel mst_polarisasi
		// berdasarkan teks polarisasi dari Excel (kolom D): "Maximize" atau "Minimize".
		//
		// Behavior:
		//   - Ditemukan     : return id_polarisasi (string), nil
		//   - Tidak ditemukan: return "", error (polarisasi tidak valid)
		//
		// Mapping yang diharapkan di DB:
		//   Maximize → id_polarisasi = "1"
		//   Minimize → id_polarisasi = "0"
		LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error)

		// GetDB mengembalikan instance *gorm.DB untuk keperluan lain (logging, dsb)
		GetDB() *gorm.DB
	}

	penyusunanKpiRepo struct {
		db *gorm.DB
	}
)

// NewPenyusunanKpiRepo membuat instance baru penyusunanKpiRepo
func NewPenyusunanKpiRepo(db *gorm.DB) *penyusunanKpiRepo {
	return &penyusunanKpiRepo{db: db}
}

// GetDB mengembalikan instance *gorm.DB
func (r *penyusunanKpiRepo) GetDB() *gorm.DB {
	return r.db
}
