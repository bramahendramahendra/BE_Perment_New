package service

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

// =============================================
// IMPLEMENTATION
// =============================================

// InsertPenyusunanKpi memproses insert KPI dengan 1 file Excel.
//
// Flow:
//  1. Validasi harus ada tepat 1 file Excel
//  2. Baca maxRows dari env (EXCEL_MAX_ROWS), fallback ke default
//  3. Parse & validasi file Excel:
//     - Pilih sheet berdasarkan req.Triwulan ("TW4" → "TW 4", lainnya → "Selain TW 4")
//     - Mapping baris ke KPI via kolom B (case-insensitive)
//     - Jika kolom B tidak cocok dengan KPI manapun → error
//     - Validasi bobot 100% per KPI
//  4. Panggil repo untuk insert dalam 1 transaksi DB
func (s *penyusunanKpiService) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	files []*multipart.FileHeader,
) (string, error) {

	// --- 1. Validasi harus ada tepat 1 file Excel ---
	if len(files) == 0 {
		return "", fmt.Errorf("tidak ada file Excel yang dikirim, harus mengirim tepat 1 file Excel")
	}
	if len(files) > 1 {
		return "", fmt.Errorf(
			"hanya boleh mengirim 1 file Excel (diterima %d file). "+
				"Semua data sub KPI dari semua KPI harus digabung dalam 1 file",
			len(files),
		)
	}

	file := files[0]

	// --- 2. Parse & validasi file Excel ---
	// Parser akan:
	//   - Menentukan sheet berdasarkan req.Triwulan
	//   - Mapping baris ke KPI via kolom B (case-insensitive)
	//   - Memvalidasi bobot 100% per KPI
	//   - Return map[kpiIndex][]SubDetailRow
	kpiSubDetails, err := ParseAndValidateExcel(file, req.Triwulan, req.Kpi)
	if err != nil {
		return "", fmt.Errorf("validasi file Excel '%s' gagal: %w", file.Filename, err)
	}

	// --- 3. [TESTING] Simulasi generate ID & log data per tabel ---
	// TODO: Hapus blok log ini setelah testing selesai,
	//       lalu uncomment pemanggilan repo di bagian bawah
	logPreviewInsert(req, kpiSubDetails)

	// --- 4. INSERT KE DB (DINONAKTIFKAN SEMENTARA UNTUK TESTING) ---
	// idPengajuan, err := s.repo.InsertPenyusunanKpi(req, kpiSubDetails)
	// if err != nil {
	// 	return "", fmt.Errorf("gagal menyimpan data KPI: %w", err)
	// }
	// return idPengajuan, nil

	// --- 5. Sementara return dummy ID untuk testing ---
	dummyID := fmt.Sprintf("%s%s%s%s",
		req.Kostl, req.Tahun, req.Triwulan,
		time.Now().Format("060102150405"),
	)
	return dummyID, nil
}

// =============================================
// LOG PREVIEW (TESTING ONLY)
// =============================================

func logPreviewInsert(req *dto.InsertPenyusunanKpiRequest, kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow) {
	log.Println("========== [PREVIEW INSERT] ==========")
	log.Printf("  Kostl       : %s", req.Kostl)
	log.Printf("  Tahun       : %s", req.Tahun)
	log.Printf("  Triwulan    : %s", req.Triwulan)
	log.Printf("  SaveAsDraft : %s", req.SaveAsDraft)
	log.Printf("  Jumlah KPI  : %d", len(req.Kpi))

	for i, kpiItem := range req.Kpi {
		rows := kpiSubDetails[i]
		log.Printf("  KPI[%d] id=%s | nama='%s' | jumlah sub KPI: %d",
			i+1, kpiItem.IdKpi, kpiItem.Kpi, len(rows))

		for j, row := range rows {
			isTW4Label := "Selain TW4"
			if row.IsTW4 {
				isTW4Label = "TW4"
			}
			log.Printf("    SubKPI[%d] No=%d | SubKPI='%s' | Bobot=%.2f | Sheet=%s",
				j+1, row.No, row.SubKPI, row.Bobot, isTW4Label)
		}
	}

	approvalPreview := req.ApprovalList
	if len(approvalPreview) > 80 {
		approvalPreview = approvalPreview[:80] + "..."
	}
	rawJSON, _ := json.MarshalIndent(req, "", "  ")
	_ = strings.Contains(string(rawJSON), "")
	log.Printf("  ApprovalList (preview): %s", approvalPreview)
	log.Println("=======================================")
}
