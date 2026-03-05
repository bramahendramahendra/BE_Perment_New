package service

import (
	"fmt"
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

// =============================================
// IMPLEMENTATION
// =============================================

// InsertPenyusunanKpi memproses insert KPI dengan file Excel.
//
// Flow:
//  1. Validasi jumlah file harus sama dengan jumlah item di req.Kpi
//  2. Parse & validasi SEMUA file Excel terlebih dahulu
//     → Jika ada 1 file gagal validasi, langsung return error (tidak ada yang masuk DB)
//  3. Setelah semua valid, panggil repo untuk insert dalam 1 transaksi DB
//     → IDPengajuan dan semua ID turunan di-generate di repo
//     → Jika ada yang gagal saat insert, semua di-rollback
//
// Return: idPengajuan yang di-generate backend, error
func (s *penyusunanKpiService) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	files []*multipart.FileHeader,
) (string, error) {
	// --- 1. Validasi jumlah file harus sama dengan jumlah KPI ---
	if len(files) != len(req.Kpi) {
		return "", fmt.Errorf(
			"jumlah file Excel (%d) tidak sesuai dengan jumlah KPI (%d). "+
				"Setiap KPI harus memiliki 1 file Excel dengan urutan yang sama",
			len(files), len(req.Kpi),
		)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("tidak ada file Excel yang dikirim")
	}

	// --- 2. Parse & validasi semua file Excel sebelum insert DB ---
	// Semua file harus valid dulu — jika ada yang gagal, tidak ada yang masuk DB
	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)

	for i, file := range files {
		kpiName := req.Kpi[i].Kpi

		rows, err := ParseAndValidateExcel(file)
		if err != nil {
			return "", fmt.Errorf(
				"validasi gagal pada file Excel KPI ke-%d ('%s' — KPI: '%s'): %w",
				i+1, file.Filename, kpiName, err,
			)
		}

		kpiSubDetails[i] = rows
	}

	// --- 3. Semua Excel valid → insert ke DB dalam 1 transaksi ---
	// IDPengajuan dan semua ID turunan di-generate di repo
	idPengajuan, err := s.repo.InsertPenyusunanKpi(req, kpiSubDetails)
	if err != nil {
		return "", fmt.Errorf("gagal menyimpan data KPI: %w", err)
	}

	return idPengajuan, nil
}
