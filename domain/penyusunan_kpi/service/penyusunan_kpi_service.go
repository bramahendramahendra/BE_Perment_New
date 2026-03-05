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

// InsertPenyusunanKpi memproses insert KPI dengan file Excel.
//
// Flow:
//  1. Validasi jumlah file harus sama dengan jumlah item di req.Kpi
//  2. Parse & validasi SEMUA file Excel terlebih dahulu
//     → Jika ada 1 file gagal validasi, langsung return error (tidak ada yang masuk DB)
//  3. [TESTING] Log semua data yang akan di-insert per tabel
//  4. Setelah semua valid, panggil repo untuk insert dalam 1 transaksi DB
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

	// --- 2. Tentukan batas baris Excel ---
	// Jika MaxRowsExcel tidak dikirim atau 0, gunakan default (ExcelMaxDataRows = 20)
	maxRows := 13
	if maxRows <= 0 {
		maxRows = ExcelMaxDataRows
	}

	// --- 3. Parse & validasi semua file Excel sebelum insert DB ---
	// Semua file harus valid dulu — jika ada yang gagal, tidak ada yang masuk DB
	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)

	for i, file := range files {
		kpiName := req.Kpi[i].Kpi

		rows, err := ParseAndValidateExcelWithLimit(file, maxRows)
		if err != nil {
			return "", fmt.Errorf(
				"validasi gagal pada file Excel KPI ke-%d ('%s' — KPI: '%s'): %w",
				i+1, file.Filename, kpiName, err,
			)
		}

		kpiSubDetails[i] = rows
	}

	// --- 4. [TESTING] Simulasi generate ID & log data per tabel ---
	// TODO: Hapus blok log ini setelah testing selesai,
	//       lalu uncomment pemanggilan repo di bagian bawah
	logPreviewInsert(req, kpiSubDetails)

	// --- 5. INSERT KE DB (DINONAKTIFKAN SEMENTARA UNTUK TESTING) ---
	// idPengajuan, err := s.repo.InsertPenyusunanKpi(req, kpiSubDetails)
	// if err != nil {
	// 	return "", fmt.Errorf("gagal menyimpan data KPI: %w", err)
	// }
	// return idPengajuan, nil

	// --- 6. Sementara return dummy idPengajuan dari simulasi ---
	idPengajuan := simulateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)
	return idPengajuan, nil
}

// =============================================
// TESTING HELPER — SIMULASI & LOG
// =============================================

// simulateIDPengajuan mensimulasikan generate IDPengajuan (sama persis dengan repo)
// Digunakan hanya untuk keperluan log testing
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

// logPreviewInsert mencetak preview data yang akan di-insert ke setiap tabel DB.
// Semua data ditampilkan dalam format yang mudah dibaca di terminal.
func logPreviewInsert(req *dto.InsertPenyusunanKpiRequest, kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow) {
	sep := strings.Repeat("=", 70)
	dash := strings.Repeat("-", 70)

	idPengajuan := simulateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	log.Println(sep)
	log.Println("[TESTING] PREVIEW DATA SEBELUM INSERT KE DB")
	log.Println(sep)

	// -------------------------------------------------------
	// TABLE: data_kpi
	// -------------------------------------------------------
	log.Println("[TABLE] data_kpi")
	log.Println(dash)

	dataKpi := map[string]interface{}{
		"id_pengajuan":    idPengajuan,
		"tahun":           req.Tahun,
		"triwulan":        req.Triwulan,
		"kostl":           req.Kostl,
		"kostl_tx":        req.KostlTx,
		"orgeh":           "(diambil dari tabel user by kostl saat insert)",
		"orgeh_tx":        "(diambil dari tabel user by kostl saat insert)",
		"entry_user":      req.EntryUser,
		"entry_name":      req.EntryName,
		"entry_time":      req.EntryTime,
		"approval_posisi": req.ApprovalPosisi,
		"approval_list":   req.ApprovalList,
		"status": func() interface{} {
			if req.SaveAsDraft == "1" {
				return 70
			}
			return nil
		}(),
	}
	printJSON(dataKpi)

	// -------------------------------------------------------
	// TABLE: data_kpi_detail
	// -------------------------------------------------------
	log.Println(dash)
	log.Printf("[TABLE] data_kpi_detail (%d baris)", len(req.Kpi))
	log.Println(dash)

	for i, kpiItem := range req.Kpi {
		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		dataKpiDetail := map[string]interface{}{
			"id_pengajuan":          idPengajuan,
			"id_detail":             idDetail,
			"tahun":                 req.Tahun,
			"triwulan":              req.Triwulan,
			"id_kpi":                kpiItem.IdKpi,
			"kpi":                   kpiItem.Kpi,
			"rumus":                 kpiItem.Rumus,
			"id_perspektif":         kpiItem.Persfektif,
			"id_keterangan_project": "-",
		}
		log.Printf("  → KPI ke-%d:", i+1)
		printJSON(dataKpiDetail)
	}

	// -------------------------------------------------------
	// TABLE: data_kpi_subdetail
	// -------------------------------------------------------
	totalSubRows := 0
	for _, rows := range kpiSubDetails {
		totalSubRows += len(rows)
	}

	log.Println(dash)
	log.Printf("[TABLE] data_kpi_subdetail (%d baris total dari semua Excel)", totalSubRows)
	log.Println(dash)

	subCounter := 1
	for i, kpiItem := range req.Kpi {
		rows, ok := kpiSubDetails[i]
		if !ok {
			continue
		}

		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		log.Printf("  → KPI ke-%d ('%s') — %d sub detail:", i+1, kpiItem.Kpi, len(rows))

		for _, subRow := range rows {
			idSubDetail := fmt.Sprintf("%sC%03d", idPengajuan, subCounter)

			itemQualifier := ""
			deskripsiQualifier := ""
			targetQualifier := ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			dataSubDetail := map[string]interface{}{
				"id_pengajuan":                idPengajuan,
				"id_detail":                   idDetail,
				"id_sub_detail":               idSubDetail,
				"tahun":                       req.Tahun,
				"triwulan":                    req.Triwulan,
				"id_kpi":                      kpiItem.IdKpi,
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
			}
			log.Printf("    Sub Detail %s:", idSubDetail)
			printJSON(dataSubDetail)
			subCounter++
		}
	}

	// -------------------------------------------------------
	// TABLE: data_challenge_detail
	// -------------------------------------------------------
	log.Println(dash)
	log.Printf("[TABLE] data_challenge_detail (%d baris)", len(req.ChallengeList))
	log.Println(dash)

	for i, ch := range req.ChallengeList {
		dataChallenge := map[string]interface{}{
			"id_pengajuan":        idPengajuan,
			"id_detail_challenge": ch.IdDetailChallenge,
			"tahun":               ch.Tahun,
			"triwulan":            ch.Triwulan,
			"nama_challenge":      ch.NamaChallenge,
			"deskripsi_challenge": ch.DeskripsiChallenge,
		}
		log.Printf("  → Challenge ke-%d:", i+1)
		printJSON(dataChallenge)
	}

	// -------------------------------------------------------
	// TABLE: data_method_detail
	// -------------------------------------------------------
	log.Println(dash)
	log.Printf("[TABLE] data_method_detail (%d baris)", len(req.MethodList))
	log.Println(dash)

	for i, mt := range req.MethodList {
		dataMethod := map[string]interface{}{
			"id_pengajuan":     idPengajuan,
			"id_detail_method": mt.IdDetailMethod,
			"tahun":            mt.Tahun,
			"triwulan":         mt.Triwulan,
			"nama_method":      mt.NamaMethod,
			"deskripsi_method": mt.DeskripsiMethod,
		}
		log.Printf("  → Method ke-%d:", i+1)
		printJSON(dataMethod)
	}

	log.Println(sep)
	log.Println("[TESTING] END PREVIEW")
	log.Println(sep)
}

// printJSON mencetak map sebagai JSON yang diformat rapi ke terminal
func printJSON(data map[string]interface{}) {
	b, err := json.MarshalIndent(data, "    ", "  ")
	if err != nil {
		log.Printf("    (gagal format JSON: %v)", err)
		return
	}
	log.Printf("    %s", string(b))
}
