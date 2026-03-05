package repo

import (
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

	queryInsertKpiSubDetail = `
		INSERT INTO data_kpi_subdetail 
			(id_pengajuan, id_detail, id_sub_detail, tahun, triwulan, 
			 id_kpi, kpi, rumus, otomatis, bobot, capping, 
			 target_triwulan, target_kuantitatif_triwulan, 
			 target_tahunan, target_kuantitatif_tahunan, 
			 deskripsi_glossary, item_qualifier, deskripsi_qualifier, 
			 target_qualifier, id_keterangan_project, id_qualifier) 
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
	// Format: ymd his → "260304" + "040242"
	timestamp := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		t.Year()%100, // 2 digit tahun
		int(t.Month()),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
	return kostl + tahun + triwulan + timestamp
}

// generateIDDetail membuat ID untuk setiap baris KPI (id_detail) mengikuti pola frontend lama:
//
//	Id = IDPengajuan + "P" + index 3 digit (mulai dari 001)
//	Contoh: "PS100012026TW2260304040242" + "P" + "001"
//	      = "PS100012026TW2260304040242P001"
func generateIDDetail(idPengajuan string, index int) string {
	return fmt.Sprintf("%sP%03d", idPengajuan, index+1)
}

// generateIDSubDetail membuat ID untuk setiap baris Sub KPI (id_sub_detail).
//
// Format: IDPengajuan + "C" + counter global 3 digit (counter TIDAK reset antar KPI).
//
// Counter global memastikan id_sub_detail unik tanpa perlu prefix id_detail:
//
//	KPI ke-1 (P001): IDPengajuanC001, IDPengajuanC002, IDPengajuanC003
//	KPI ke-2 (P002): IDPengajuanC004, IDPengajuanC005  ← lanjut dari C003
func generateIDSubDetail(idPengajuan string, globalIndex int) string {
	return fmt.Sprintf("%sC%03d", idPengajuan, globalIndex)
}

// =============================================
// IMPLEMENTATION
// =============================================

// InsertPenyusunanKpi melakukan insert semua data KPI dalam satu transaksi DB.
// Flow:
//  1. Generate IDPengajuan dan semua ID turunannya di backend
//  2. Cek apakah data sudah ada (tahun + triwulan + kostl)
//  3. Ambil orgeh & orgeh_tx dari tabel user
//  4. Build semua query batch insert
//  5. Eksekusi dalam 1 transaksi → commit jika semua sukses, rollback jika ada yang gagal
//
// Catatan generate ID (mengikuti pola frontend lama):
//   - IDPengajuan  = Kostl + Tahun + Triwulan + timestamp(ymdhis)
//   - id_detail    = IDPengajuan + "P" + index KPI 3 digit (P001, P002, ...)
//   - id_sub_detail = IDPengajuan + "C" + counter global (tidak reset antar KPI)
//     contoh: C001, C002, C003 (KPI P001) → C004, C005 (KPI P002)
//   - id_keterangan_project = "-" (backend otomatis)
//   - id_pengajuan challenge/method = IDPengajuan yang di-generate
func (r *penyusunanKpiRepo) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) (string, error) {
	// --- 1. Generate IDPengajuan di backend ---
	idPengajuan := generateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	// --- 2. Cek data sudah exist ---
	var countExist int
	if err := r.db.Raw(queryCheckExistKpi, req.Tahun, req.Triwulan, req.Kostl).
		Scan(&countExist).Error; err != nil {
		return "", fmt.Errorf("gagal mengecek data KPI: %w", err)
	}
	if countExist > 0 {
		return "", fmt.Errorf("data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
			req.Tahun, req.Triwulan, req.Kostl)
	}

	// --- 3. Ambil orgeh & orgeh_tx ---
	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Kostl).Row().Scan(&orgeh, &orgehTx)

	// --- 4. Tentukan status berdasarkan SaveAsDraft ---
	// Status 70 = draft, NULL = normal
	var statusKpi interface{}
	if req.SaveAsDraft == "1" {
		statusKpi = 70
	} else {
		statusKpi = nil
	}

	// --- 5. Build batch INSERT data_kpi_detail ---
	// id_detail             = IDPengajuan + "P" + index 3 digit  → P001, P002, ...
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

	// --- 6. Build batch INSERT data_kpi_subdetail (dari hasil parse Excel) ---
	// id_sub_detail = IDPengajuan + "C" + counter global (lanjut antar KPI, tidak reset)
	//   contoh: KPI P001 → C001, C002, C003 | KPI P002 → C004, C005, C006
	// id_keterangan_project = "-" (backend otomatis)
	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}

	subCounter := 1 // counter global sub detail, tidak reset antar KPI

	for i, kpiItem := range req.Kpi {
		rows, ok := kpiSubDetails[i]
		if !ok {
			continue
		}

		idDetail := idDetailMap[i]

		for _, subRow := range rows {
			idSubDetail := generateIDSubDetail(idPengajuan, subCounter) // counter global, lanjut antar KPI
			subCounter++

			itemQualifier := ""
			deskripsiQualifier := ""
			targetQualifier := ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			subDetailPlaceholders = append(subDetailPlaceholders,
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			subDetailArgs = append(subDetailArgs,
				idPengajuan,
				idDetail,
				idSubDetail,
				req.Tahun,
				req.Triwulan,
				kpiItem.IdKpi,
				subRow.SubKPI,
				subRow.Polarisasi,
				"0", // otomatis default "0"
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
				"-", // id_keterangan_project: backend otomatis isi "-"
				subRow.TerdapatQualifier,
			)
		}
	}

	// --- 7. Build batch INSERT data_challenge_detail ---
	// id_pengajuan : dari IDPengajuan yang di-generate backend
	// tahun & triwulan : dari item (bisa "-" jika non-TW4)
	challengePlaceholders := []string{}
	challengeArgs := []interface{}{}

	for _, ch := range req.ChallengeList {
		challengePlaceholders = append(challengePlaceholders, "(?, ?, ?, ?, ?, ?)")
		challengeArgs = append(challengeArgs,
			idPengajuan, // dari backend, bukan dari frontend
			ch.IdDetailChallenge,
			ch.Tahun,
			ch.Triwulan,
			ch.NamaChallenge,
			ch.DeskripsiChallenge,
		)
	}

	// --- 8. Build batch INSERT data_method_detail ---
	// id_pengajuan : dari IDPengajuan yang di-generate backend
	// tahun & triwulan : dari item (bisa "-" jika non-TW4)
	methodPlaceholders := []string{}
	methodArgs := []interface{}{}

	for _, mt := range req.MethodList {
		methodPlaceholders = append(methodPlaceholders, "(?, ?, ?, ?, ?, ?)")
		methodArgs = append(methodArgs,
			idPengajuan, // dari backend, bukan dari frontend
			mt.IdDetailMethod,
			mt.Tahun,
			mt.Triwulan,
			mt.NamaMethod,
			mt.DeskripsiMethod,
		)
	}

	// --- 9. Finalisasi query dengan format batch values ---
	finalQueryKpiDetail := fmt.Sprintf(queryInsertKpiDetail,
		strings.Join(kpiDetailPlaceholders, ","))

	finalQuerySubDetail := fmt.Sprintf(queryInsertKpiSubDetail,
		strings.Join(subDetailPlaceholders, ","))

	finalQueryChallenge := fmt.Sprintf(queryInsertChallengeDetail,
		strings.Join(challengePlaceholders, ","))

	finalQueryMethod := fmt.Sprintf(queryInsertMethodDetail,
		strings.Join(methodPlaceholders, ","))

	// --- 10. Eksekusi dalam 1 transaksi DB ---
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// Insert data_kpi (header)
	if err := tx.Exec(queryInsertKpi,
		idPengajuan,
		req.Tahun,
		req.Triwulan,
		req.Kostl,
		req.KostlTx,
		orgeh,
		orgehTx,
		req.EntryUser,
		req.EntryName,
		req.EntryTime,
		req.ApprovalPosisi,
		req.ApprovalList,
		statusKpi,
	).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("gagal insert data_kpi: %w", err)
	}

	// Insert data_kpi_detail
	if len(kpiDetailPlaceholders) > 0 {
		if err := tx.Exec(finalQueryKpiDetail, kpiDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_detail: %w", err)
		}
	}

	// Insert data_kpi_subdetail
	if len(subDetailPlaceholders) > 0 {
		if err := tx.Exec(finalQuerySubDetail, subDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_subdetail: %w", err)
		}
	}

	// Insert data_challenge_detail
	if len(challengePlaceholders) > 0 {
		if err := tx.Exec(finalQueryChallenge, challengeArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_challenge_detail: %w", err)
		}
	}

	// Insert data_method_detail
	if len(methodPlaceholders) > 0 {
		if err := tx.Exec(finalQueryMethod, methodArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_method_detail: %w", err)
		}
	}

	// Commit jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return idPengajuan, nil
}
