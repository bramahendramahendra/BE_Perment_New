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
//  2. Parse & validasi file Excel:
//     - Pilih sheet berdasarkan req.Triwulan ("TW4" → "TW 4", lainnya → "Selain TW 4")
//     - Mapping baris ke KPI via kolom B (case-insensitive)
//     - Jika kolom B tidak cocok dengan KPI manapun → error
//     - Validasi bobot 100% per KPI
//  3. [TESTING] Log semua data yang akan di-insert per tabel — INSERT DB DINONAKTIFKAN
//  4. Return dummy idPengajuan untuk keperluan testing
//
// TODO: Setelah testing selesai dan log sudah diverifikasi:
//   - Hapus blok logPreviewInsert & return dummy
//   - Uncomment blok INSERT KE DB di bawah
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

	// --- 3. [TESTING] Log semua data yang akan di-insert ---
	// INSERT KE DB DINONAKTIFKAN — hanya log preview saja
	logPreviewInsert(req, kpiSubDetails)

	// --- 4. INSERT KE DB (DINONAKTIFKAN SEMENTARA UNTUK TESTING) ---
	// Uncomment blok ini setelah log sudah diverifikasi & testing selesai:
	//
	// idPengajuan, err := s.repo.InsertPenyusunanKpi(req, kpiSubDetails)
	// if err != nil {
	// 	return "", fmt.Errorf("gagal menyimpan data KPI: %w", err)
	// }
	// return idPengajuan, nil

	// --- 5. Return dummy idPengajuan untuk testing ---
	dummyID := simulateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)
	return dummyID, nil
}

// =============================================
// TESTING HELPER — SIMULASI ID
// =============================================

// simulateIDPengajuan mensimulasikan generate IDPengajuan (sama persis dengan repo).
// Digunakan hanya untuk keperluan log testing.
func simulateIDPengajuan(kostl, tahun, triwulan string) string {
	t := time.Now()
	timestamp := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		t.Year()%100,
		int(t.Month()),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
	return kostl + tahun + triwulan + timestamp
}

// =============================================
// TESTING HELPER — LOG PREVIEW
// =============================================

// logPreviewInsert mencetak preview lengkap semua data yang akan di-insert ke DB.
//
// Output log mencakup:
//   - [TABLE] data_kpi           → 1 baris header
//   - [TABLE] data_kpi_detail    → 1 baris per KPI
//   - [TABLE] data_kpi_subdetail → N baris per KPI dari Excel
//   - [TABLE] data_challenge_detail
//   - [TABLE] data_method_detail
//   - [SUMMARY] ringkasan jumlah baris per tabel
//
// Catatan kolom P–U (result–deskripsi_context):
//   - Sheet "TW 4"        → tampil nilai string
//   - Sheet "Selain TW 4" → tampil "NULL" (akan disimpan NULL di DB)
func logPreviewInsert(
	req *dto.InsertPenyusunanKpiRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) {
	sep := strings.Repeat("=", 75)
	dash := strings.Repeat("-", 75)

	idPengajuan := simulateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	// Tentukan sheet yang digunakan
	sheetUsed := "Selain TW 4"
	if strings.EqualFold(req.Triwulan, "TW4") {
		sheetUsed = "TW 4"
	}

	log.Println(sep)
	log.Println("[TESTING] PREVIEW DATA SEBELUM INSERT KE DB")
	log.Println("[TESTING] INSERT DB DINONAKTIFKAN — DATA BELUM MASUK DB")
	log.Println(sep)
	log.Printf("  Simulated IDPengajuan : %s", idPengajuan)
	log.Printf("  Sheet Excel digunakan : %s (Triwulan: %s)", sheetUsed, req.Triwulan)
	log.Println(sep)

	// -------------------------------------------------------
	// TABLE: data_kpi (1 baris header)
	// -------------------------------------------------------
	log.Println()
	log.Printf("[TABLE] data_kpi — 1 baris")
	log.Println(dash)

	statusVal := interface{}(nil)
	if req.SaveAsDraft == "1" {
		statusVal = 70
	}

	printJSON(map[string]interface{}{
		"id_pengajuan":    idPengajuan,
		"tahun":           req.Tahun,
		"triwulan":        req.Triwulan,
		"kostl":           req.Kostl,
		"kostl_tx":        req.KostlTx,
		"orgeh":           fmt.Sprintf("(SELECT orgeh FROM user WHERE kostl='%s')", req.Kostl),
		"orgeh_tx":        fmt.Sprintf("(SELECT orgeh_tx FROM user WHERE kostl='%s')", req.Kostl),
		"entry_user":      req.EntryUser,
		"entry_name":      req.EntryName,
		"entry_time":      req.EntryTime,
		"approval_posisi": req.ApprovalPosisi,
		"approval_list":   req.ApprovalList,
		"status":          statusVal,
	})

	// -------------------------------------------------------
	// TABLE: data_kpi_detail (1 baris per KPI)
	// -------------------------------------------------------
	log.Println()
	log.Printf("[TABLE] data_kpi_detail — %d baris", len(req.Kpi))
	log.Println(dash)

	for i, kpiItem := range req.Kpi {
		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		log.Printf("  Baris %d (KPI: '%s'):", i+1, kpiItem.Kpi)
		printJSON(map[string]interface{}{
			"id_pengajuan":          idPengajuan,
			"id_detail":             idDetail,
			"tahun":                 req.Tahun,
			"triwulan":              req.Triwulan,
			"id_kpi":                kpiItem.IdKpi,
			"kpi":                   kpiItem.Kpi,
			"rumus":                 kpiItem.Rumus,
			"id_perspektif":         kpiItem.Persfektif,
			"id_keterangan_project": "-",
		})
	}

	// -------------------------------------------------------
	// TABLE: data_kpi_subdetail (dari Excel, dikelompokkan per KPI)
	// -------------------------------------------------------
	totalSubRows := 0
	for _, rows := range kpiSubDetails {
		totalSubRows += len(rows)
	}

	log.Println()
	log.Printf("[TABLE] data_kpi_subdetail — %d baris total (sheet: %s)", totalSubRows, sheetUsed)
	log.Println(dash)

	subCounter := 1
	for i, kpiItem := range req.Kpi {
		rows, ok := kpiSubDetails[i]
		if !ok {
			log.Printf("  [WARN] KPI ke-%d ('%s') tidak memiliki data sub KPI di Excel!", i+1, kpiItem.Kpi)
			continue
		}

		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		log.Printf("  KPI ke-%d: '%s' | id_detail: %s | jumlah sub KPI: %d",
			i+1, kpiItem.Kpi, idDetail, len(rows))

		for j, subRow := range rows {
			idSubDetail := fmt.Sprintf("%sC%03d", idPengajuan, subCounter)

			// Qualifier hanya diisi jika TerdapatQualifier = "Ya"
			itemQualifier := ""
			deskripsiQualifier := ""
			targetQualifier := ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			log.Printf("    Sub KPI %d/%d | id_sub_detail: %s", j+1, len(rows), idSubDetail)
			printJSON(map[string]interface{}{
				// --- Kolom identitas ---
				"id_pengajuan":  idPengajuan,
				"id_detail":     idDetail,
				"id_sub_detail": idSubDetail,
				"tahun":         req.Tahun,
				"triwulan":      req.Triwulan,
				"id_kpi":        kpiItem.IdKpi,
				// --- Kolom A–O (selalu ada, kedua sheet) ---
				"kpi":                         subRow.SubKPI,
				"rumus":                       subRow.Polarisasi,
				"otomatis":                    "0",
				"bobot":                       subRow.Bobot,
				"capping":                     subRow.Capping,
				"target_triwulan":             subRow.TargetTriwulan,
				"target_kuantitatif_triwulan": subRow.TargetKuantitatifTriwulan,
				"target_tahunan":              subRow.TargetTahunan,
				"target_kuantitatif_tahunan":  subRow.TargetKuantitatifTahunan,
				"deskripsi_glossary":          subRow.Glossary,
				"item_qualifier":              itemQualifier,
				"deskripsi_qualifier":         deskripsiQualifier,
				"target_qualifier":            targetQualifier,
				"id_keterangan_project":       "-",
				"id_qualifier":                subRow.TerdapatQualifier,
				// --- Kolom P–U (hanya sheet "TW 4", selain itu NULL) ---
				"result":            nullableStringLog(subRow.Result),
				"deskripsi_result":  nullableStringLog(subRow.DeskripsiResult),
				"process":           nullableStringLog(subRow.Process),
				"deskripsi_process": nullableStringLog(subRow.DeskripsiProcess),
				"context":           nullableStringLog(subRow.Context),
				"deskripsi_context": nullableStringLog(subRow.DeskripsiContext),
			})
			subCounter++
		}
	}

	// -------------------------------------------------------
	// TABLE: data_challenge_detail
	// -------------------------------------------------------
	log.Println()
	log.Printf("[TABLE] data_challenge_detail — %d baris", len(req.ChallengeList))
	log.Println(dash)

	for i, ch := range req.ChallengeList {
		log.Printf("  Baris %d:", i+1)
		printJSON(map[string]interface{}{
			"id_pengajuan":        idPengajuan,
			"id_detail_challenge": ch.IdDetailChallenge,
			"tahun":               ch.Tahun,
			"triwulan":            ch.Triwulan,
			"nama_challenge":      ch.NamaChallenge,
			"deskripsi_challenge": ch.DeskripsiChallenge,
		})
	}

	// -------------------------------------------------------
	// TABLE: data_method_detail
	// -------------------------------------------------------
	log.Println()
	log.Printf("[TABLE] data_method_detail — %d baris", len(req.MethodList))
	log.Println(dash)

	for i, mt := range req.MethodList {
		log.Printf("  Baris %d:", i+1)
		printJSON(map[string]interface{}{
			"id_pengajuan":     idPengajuan,
			"id_detail_method": mt.IdDetailMethod,
			"tahun":            mt.Tahun,
			"triwulan":         mt.Triwulan,
			"nama_method":      mt.NamaMethod,
			"deskripsi_method": mt.DeskripsiMethod,
		})
	}

	// -------------------------------------------------------
	// SUMMARY
	// -------------------------------------------------------
	log.Println()
	log.Println(sep)
	log.Println("[TESTING] SUMMARY — JUMLAH DATA YANG AKAN DI-INSERT")
	log.Println(dash)
	log.Printf("  %-30s : 1 baris", "data_kpi")
	log.Printf("  %-30s : %d baris", "data_kpi_detail", len(req.Kpi))
	log.Printf("  %-30s : %d baris total", "data_kpi_subdetail", totalSubRows)
	for i, kpiItem := range req.Kpi {
		rows := kpiSubDetails[i]
		log.Printf("    └─ KPI %-25s : %d sub KPI", "'"+kpiItem.Kpi+"'", len(rows))
	}
	log.Printf("  %-30s : %d baris", "data_challenge_detail", len(req.ChallengeList))
	log.Printf("  %-30s : %d baris", "data_method_detail", len(req.MethodList))
	log.Println(sep)
	log.Println("[TESTING] END PREVIEW — DATA BELUM DIINSERT KE DB")
	log.Println(sep)
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// nullableStringLog mengkonversi *string menjadi nilai yang siap di-log.
//   - nil  → "NULL" (akan disimpan NULL di DB)
//   - &val → nilai string asli
func nullableStringLog(s *string) interface{} {
	if s == nil {
		return "NULL"
	}
	return *s
}

// printJSON mencetak map sebagai JSON yang diformat rapi ke terminal log.
func printJSON(data map[string]interface{}) {
	b, err := json.MarshalIndent(data, "      ", "  ")
	if err != nil {
		log.Printf("      (gagal format JSON: %v)", err)
		return
	}
	log.Printf("      %s", string(b))
}
