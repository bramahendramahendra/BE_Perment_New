package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

// =============================================
// QUERY CONSTANTS
// =============================================

const (
	queryCheckExistKpi = `
		SELECT COUNT(id_pengajuan) 
		FROM data_kpi 
		WHERE tahun = ? AND triwulan = ? AND kostl = ?`

	queryGetOrgeh = `
		SELECT orgeh, orgeh_tx 
		FROM user 
		WHERE kostl = ? 
		ORDER BY HILFM ASC 
		LIMIT 1`

	queryInsertKpi = `
		INSERT INTO data_kpi 
			(id_pengajuan, tahun, triwulan, kostl, kostl_tx, orgeh, orgeh_tx, 
			 entry_user, entry_name, entry_time, approval_posisi, approval_list, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	queryInsertKpiDetail = `
		INSERT INTO data_kpi_detail 
			(id_pengajuan, id_detail, tahun, triwulan, id_kpi, kpi, rumus, 
			 id_perspektif, id_keterangan_project) 
		VALUES %s`

	// queryInsertKpiSubDetail mencakup kolom P–U (result, deskripsi_result, process,
	// deskripsi_process, context, deskripsi_context) yang:
	//   - Terisi nilai string jika berasal dari sheet "TW 4"
	//   - NULL jika berasal dari sheet "Selain TW 4"
	//
	// Kolom id_kpi   → dari hasil lookup mst_kpi (atau "0" jika tidak ditemukan)
	// Kolom rumus    → id_polarisasi dari mst_polarisasi (1=Maximize, 0=Minimize)
	queryInsertKpiSubDetail = `
		INSERT INTO data_kpi_subdetail 
			(id_pengajuan, id_detail, id_sub_detail, tahun, triwulan, 
			 id_kpi, kpi, rumus, otomatis, bobot, capping, 
			 target_triwulan, target_kuantitatif_triwulan, 
			 target_tahunan, target_kuantitatif_tahunan, 
			 deskripsi_glossary, item_qualifier, deskripsi_qualifier, 
			 target_qualifier, id_keterangan_project, id_qualifier,
			 result, deskripsi_result,
			 process, deskripsi_process,
			 context, deskripsi_context) 
		VALUES %s`

	queryInsertChallengeDetail = `
		INSERT INTO data_challenge_detail 
			(id_pengajuan, id_detail_challenge, tahun, triwulan, 
			 nama_challenge, deskripsi_challenge) 
		VALUES %s`

	queryInsertMethodDetail = `
		INSERT INTO data_method_detail 
			(id_pengajuan, id_detail_method, tahun, triwulan, 
			 nama_method, deskripsi_method) 
		VALUES %s`

	// queryLookupSubKpi mencari id_kpi, kpi, dan rumus dari mst_kpi secara case-insensitive.
	// Digunakan untuk mengisi kolom id_kpi dan validasi rumus pada data_kpi_subdetail.
	//
	// Behavior:
	//   - Ditemukan     : return id_kpi, kpi (nama dari DB), rumus
	//   - Tidak ditemukan: sql.ErrNoRows → di-handle caller dengan id_kpi = "0"
	queryLookupSubKpi = `
		SELECT id_kpi, kpi, rumus
		FROM mst_kpi
		WHERE LOWER(kpi) = LOWER(?)
		LIMIT 1`

	// queryLookupPolarisasi mencari id_polarisasi dari mst_polarisasi berdasarkan teks polarisasi.
	// Mapping: Maximize = "1", Minimize = "0"
	//
	// Behavior:
	//   - Ditemukan     : return id_polarisasi
	//   - Tidak ditemukan: sql.ErrNoRows → di-handle caller sebagai error
	queryLookupPolarisasi = `
		SELECT id_polarisasi
		FROM mst_polarisasi
		WHERE LOWER(polarisasi) = LOWER(?)
		LIMIT 1`
)

// =============================================
// HELPER — GENERATE ID
// =============================================

// generateIDPengajuan membuat IDPengajuan mengikuti pola frontend lama:
//
//	IDPengajuan = Kostl + Tahun + Triwulan + timestamp(ymdhis)
//	Contoh: "PS10001" + "2026" + "TW2" + "260304040242"
//	      = "PS100012026TW2260304040242"
func generateIDPengajuan(kostl, tahun, triwulan string) string {
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

// generateIDDetail membuat ID untuk setiap baris KPI (id_detail):
//
//	Id = IDPengajuan + "P" + index 3 digit (mulai dari 001)
//	Contoh: "PS100012026TW2260304040242P001"
func generateIDDetail(idPengajuan string, index int) string {
	return fmt.Sprintf("%sP%03d", idPengajuan, index+1)
}

// generateIDSubDetail membuat ID untuk setiap baris Sub KPI (id_sub_detail).
//
// Format: IDPengajuan + "C" + counter global 3 digit (counter TIDAK reset antar KPI).
//
//	KPI ke-1 (P001): IDPengajuanC001, IDPengajuanC002, IDPengajuanC003
//	KPI ke-2 (P002): IDPengajuanC004, IDPengajuanC005  ← lanjut dari C003
func generateIDSubDetail(idPengajuan string, globalIndex int) string {
	return fmt.Sprintf("%sC%03d", idPengajuan, globalIndex)
}

// =============================================
// LOOKUP METHODS
// =============================================

// LookupSubKpiMaster mencari id_kpi, nama kpi, dan rumus dari tabel mst_kpi
// secara case-insensitive berdasarkan teks sub KPI dari Excel (kolom C).
//
// Behavior:
//   - Ditemukan      : return id_kpi, kpiFromDB (nama dari DB), rumus, nil
//   - Tidak ditemukan: return "0", subKpiText (nama asli Excel), "", nil
//     (bukan error — id_kpi = 0 adalah kondisi valid, validasi rumus di-skip)
func (r *penyusunanKpiRepo) LookupSubKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error) {
	row := r.db.Raw(queryLookupSubKpi, subKpiText).Row()
	if scanErr := row.Scan(&idKpi, &kpiFromDB, &rumus); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			// Tidak ditemukan — kondisi valid, id_kpi = 0
			return "0", subKpiText, "", nil
		}
		// Error DB yang tidak terduga
		return "0", subKpiText, "", fmt.Errorf("gagal lookup mst_kpi untuk Sub KPI '%s': %w", subKpiText, scanErr)
	}
	return idKpi, kpiFromDB, rumus, nil
}

// LookupPolarisasi mencari id_polarisasi dari tabel mst_polarisasi
// berdasarkan teks polarisasi dari Excel (kolom D).
//
// Behavior:
//   - Ditemukan     : return id_polarisasi (string), nil
//   - Tidak ditemukan: return "", error
//
// Mapping yang diharapkan di DB:
//
//	Maximize → id_polarisasi = "1"
//	Minimize → id_polarisasi = "0"
func (r *penyusunanKpiRepo) LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error) {
	row := r.db.Raw(queryLookupPolarisasi, polarisasiText).Row()
	if scanErr := row.Scan(&idPolarisasi); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return "", fmt.Errorf(
				"polarisasi '%s' tidak ditemukan di tabel mst_polarisasi. "+
					"Nilai yang valid: 'Maximize' atau 'Minimize'",
				polarisasiText,
			)
		}
		return "", fmt.Errorf("gagal lookup mst_polarisasi untuk polarisasi '%s': %w", polarisasiText, scanErr)
	}
	return idPolarisasi, nil
}

// =============================================
// IMPLEMENTATION — INSERT
// =============================================

// InsertPenyusunanKpi melakukan insert semua data KPI dalam satu transaksi DB.
//
// Catatan: kpiSubDetails yang diterima sudah berisi IdSubKpi dan IdPolarisasi
// yang telah di-resolve di service sebelum fungsi ini dipanggil.
//
// Flow:
//  1. Generate IDPengajuan dan semua ID turunannya di backend
//  2. Cek apakah data sudah ada (tahun + triwulan + kostl)
//  3. Ambil orgeh & orgeh_tx dari tabel user
//  4. Tentukan status berdasarkan SaveAsDraft
// Catatan generate ID:
//   - IDPengajuan   = Kostl + Tahun + Triwulan + timestamp(ymdhis)
//   - id_detail     = IDPengajuan + "P" + index KPI 3 digit (P001, P002, ...)
//   - id_sub_detail = IDPengajuan + "C" + counter global (tidak reset antar KPI)
//   - id_keterangan_project = "-" (backend otomatis)
//  5. Build batch INSERT data_kpi_detail
//  6. Build batch INSERT data_kpi_subdetail
//     - id_kpi  = subRow.IdSubKpi  (dari mst_kpi, atau "0")
//     - kpi     = subRow.SubKPI    (nama dari DB, atau nama Excel)
//     - rumus   = subRow.IdPolarisasi (id_polarisasi dari mst_polarisasi)
//  7. Build batch INSERT data_challenge_detail
//  8. Build batch INSERT data_method_detail
//  9. Eksekusi semua INSERT dalam 1 transaksi
func (r *penyusunanKpiRepo) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) (string, error) {

	// --- 1. Generate IDPengajuan di backend ---
	idPengajuan := generateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	// --- 2. Cek data sudah exist (tahun + triwulan + kostl) ---
	var countExist int
	if err := r.db.Raw(queryCheckExistKpi, req.Tahun, req.Triwulan, req.Kostl).
		Scan(&countExist).Error; err != nil {
		return "", fmt.Errorf("gagal mengecek data KPI: %w", err)
	}
	if countExist > 0 {
		return "", fmt.Errorf(
			"data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
			req.Tahun, req.Triwulan, req.Kostl,
		)
	}

	// --- 3. Ambil orgeh & orgeh_tx dari tabel user ---
	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Kostl).Row().Scan(&orgeh, &orgehTx)

	// --- 4. Tentukan status berdasarkan SaveAsDraft ---
	// Status 70 = draft, NULL = submit normal
	var statusKpi interface{}
	if req.SaveAsDraft == "1" {
		statusKpi = 70
	} else {
		statusKpi = nil
	}

	// --- 5. Build batch INSERT data_kpi_detail ---
	// id_detail             = IDPengajuan + "P" + index 3 digit → P001, P002, ...
	// id_keterangan_project = "-" (backend otomatis)
	kpiDetailPlaceholders := []string{}
	kpiDetailArgs := []interface{}{}

	// Simpan idDetail per index KPI agar bisa dipakai saat build sub detail
	idDetailMap := make(map[int]string)

	for i, kpiItem := range req.Kpi {
		idDetail := generateIDDetail(idPengajuan, i)
		idDetailMap[i] = idDetail

		kpiDetailPlaceholders = append(kpiDetailPlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		kpiDetailArgs = append(kpiDetailArgs,
			idPengajuan,
			idDetail,
			req.Tahun,
			req.Triwulan,
			kpiItem.IdKpi,
			kpiItem.Kpi,
			kpiItem.Rumus,
			kpiItem.Persfektif,
			"-", // id_keterangan_project: backend otomatis isi "-"
		)
	}

	// --- 6. Build batch INSERT data_kpi_subdetail ---
	//
	// id_kpi  = subRow.IdSubKpi     → dari mst_kpi (case-insensitive lookup), "0" jika tidak ditemukan
	// kpi     = subRow.SubKPI       → nama dari DB jika ditemukan, nama Excel jika tidak
	// rumus   = subRow.IdPolarisasi → id_polarisasi dari mst_polarisasi (1=Maximize, 0=Minimize)
	//
	// id_sub_detail = IDPengajuan + "C" + counter global (lanjut antar KPI, tidak reset)
	//   contoh: KPI P001 → C001, C002, C003 | KPI P002 → C004, C005 | dst
	//
	// Kolom P–U (result–deskripsi_context):
	//   - *string != nil → insert nilai string (sheet "TW 4")
	//   - *string == nil → insert NULL (sheet "Selain TW 4")
	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}

	subCounter := 1 // counter global, tidak reset antar KPI

	for i := range req.Kpi {
		rows, ok := kpiSubDetails[i]
		if !ok {
			continue
		}

		idDetail := idDetailMap[i]

		for _, subRow := range rows {
			idSubDetail := generateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Qualifier: hanya isi jika TerdapatQualifier = "Ya"
			itemQualifier := ""
			deskripsiQualifier := ""
			targetQualifier := ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			// Konversi *string → interface{} untuk kolom P–U
			// nil *string → nil interface{} → NULL di DB
			var result, deskripsiResult, process, deskripsiProcess, context, deskripsiContext interface{}
			if subRow.Result != nil {
				result = *subRow.Result
			}
			if subRow.DeskripsiResult != nil {
				deskripsiResult = *subRow.DeskripsiResult
			}
			if subRow.Process != nil {
				process = *subRow.Process
			}
			if subRow.DeskripsiProcess != nil {
				deskripsiProcess = *subRow.DeskripsiProcess
			}
			if subRow.Context != nil {
				context = *subRow.Context
			}
			if subRow.DeskripsiContext != nil {
				deskripsiContext = *subRow.DeskripsiContext
			}

			// 27 kolom: 21 kolom lama + 6 kolom P–U baru
			subDetailPlaceholders = append(subDetailPlaceholders,
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			subDetailArgs = append(subDetailArgs,
				idPengajuan,
				idDetail,
				idSubDetail,
				req.Tahun,
				req.Triwulan,
				subRow.IdSubKpi,     // id_kpi dari mst_kpi (atau "0")
				subRow.SubKPI,       // nama dari DB (atau Excel jika tidak ditemukan)
				subRow.IdPolarisasi, // id_polarisasi dari mst_polarisasi (1=Maximize, 0=Minimize)
				"0",                 // otomatis: selalu "0"
				subRow.Bobot,
				subRow.Capping,
				subRow.TargetTriwulan,
				subRow.TargetKuantitatifTriwulan,
				subRow.TargetTahunan,
				subRow.TargetKuantitatifTahunan,
				subRow.Glossary,
				itemQualifier,
				deskripsiQualifier,
				targetQualifier,
				"-",                      // id_keterangan_project: backend otomatis isi "-"
				subRow.TerdapatQualifier, // id_qualifier
				result,
				deskripsiResult,
				process,
				deskripsiProcess,
				context,
				deskripsiContext,
			)
		}
	}

	// --- 7. Build batch INSERT data_challenge_detail ---
	challengePlaceholders := []string{}
	challengeArgs := []interface{}{}

	for _, ch := range req.ChallengeList {
		challengePlaceholders = append(challengePlaceholders, "(?, ?, ?, ?, ?, ?)")
		challengeArgs = append(challengeArgs,
			idPengajuan,
			ch.IdDetailChallenge,
			ch.Tahun,
			ch.Triwulan,
			ch.NamaChallenge,
			ch.DeskripsiChallenge,
		)
	}

	// --- 8. Build batch INSERT data_method_detail ---
	methodPlaceholders := []string{}
	methodArgs := []interface{}{}

	for _, mt := range req.MethodList {
		methodPlaceholders = append(methodPlaceholders, "(?, ?, ?, ?, ?, ?)")
		methodArgs = append(methodArgs,
			idPengajuan,
			mt.IdDetailMethod,
			mt.Tahun,
			mt.Triwulan,
			mt.NamaMethod,
			mt.DeskripsiMethod,
		)
	}

	// --- 9. Eksekusi semua INSERT dalam 1 transaksi ---
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// INSERT data_kpi (1 baris)
	if err := tx.Exec(queryInsertKpi,
		idPengajuan, req.Tahun, req.Triwulan, req.Kostl, req.KostlTx,
		orgeh, orgehTx, req.EntryUser, req.EntryName, req.EntryTime,
		req.ApprovalPosisi, req.ApprovalList, statusKpi,
	).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("gagal insert data_kpi: %w", err)
	}

	// INSERT data_kpi_detail (batch)
	if len(kpiDetailPlaceholders) > 0 {
		queryDetail := fmt.Sprintf(queryInsertKpiDetail, strings.Join(kpiDetailPlaceholders, ", "))
		if err := tx.Exec(queryDetail, kpiDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_detail: %w", err)
		}
	}

	// INSERT data_kpi_subdetail (batch)
	if len(subDetailPlaceholders) > 0 {
		querySubDetail := fmt.Sprintf(queryInsertKpiSubDetail, strings.Join(subDetailPlaceholders, ", "))
		if err := tx.Exec(querySubDetail, subDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_subdetail: %w", err)
		}
	}

	// INSERT data_challenge_detail (batch)
	if len(challengePlaceholders) > 0 {
		queryChallenge := fmt.Sprintf(queryInsertChallengeDetail, strings.Join(challengePlaceholders, ", "))
		if err := tx.Exec(queryChallenge, challengeArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_challenge_detail: %w", err)
		}
	}

	// INSERT data_method_detail (batch)
	if len(methodPlaceholders) > 0 {
		queryMethod := fmt.Sprintf(queryInsertMethodDetail, strings.Join(methodPlaceholders, ", "))
		if err := tx.Exec(queryMethod, methodArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_method_detail: %w", err)
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return idPengajuan, nil
}
