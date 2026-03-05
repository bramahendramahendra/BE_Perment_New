package repo

import (
	"fmt"
	"strings"

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
// IMPLEMENTATION
// =============================================

// InsertPenyusunanKpi melakukan insert semua data KPI dalam satu transaksi DB.
// Flow:
//  1. Cek apakah data sudah ada (tahun + triwulan + kostl)
//  2. Ambil orgeh & orgeh_tx dari tabel user
//  3. Build semua query batch insert
//  4. Eksekusi dalam 1 transaksi → commit jika semua sukses, rollback jika ada yang gagal
func (r *penyusunanKpiRepo) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) error {
	// --- 1. Cek data sudah exist ---
	var countExist int
	if err := r.db.Raw(queryCheckExistKpi, req.Tahun, req.Triwulan, req.Kostl).
		Scan(&countExist).Error; err != nil {
		return fmt.Errorf("gagal mengecek data KPI: %w", err)
	}
	if countExist > 0 {
		return fmt.Errorf("data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
			req.Tahun, req.Triwulan, req.Kostl)
	}

	// --- 2. Ambil orgeh & orgeh_tx ---
	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Kostl).Row().Scan(&orgeh, &orgehTx)

	// --- 3. Tentukan status berdasarkan SaveAsDraft ---
	// Status 70 = draft, NULL = normal (insert tanpa kolom status)
	var statusKpi interface{}
	if req.SaveAsDraft == "1" {
		statusKpi = 70
	} else {
		statusKpi = nil
	}

	// --- 4. Build batch INSERT data_kpi_detail ---
	kpiDetailPlaceholders := []string{}
	kpiDetailArgs := []interface{}{}

	for i, kpiItem := range req.Kpi {
		inc := fmt.Sprintf("000%d", i)
		idDetail := req.IDPengajuan + inc[len(inc)-3:]

		kpiDetailPlaceholders = append(kpiDetailPlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		kpiDetailArgs = append(kpiDetailArgs,
			req.IDPengajuan,
			idDetail,
			req.Tahun,
			req.Triwulan,
			kpiItem.IdKpi,
			kpiItem.Kpi,
			kpiItem.Rumus,
			kpiItem.Persfektif,
			kpiItem.KeteranganProject,
		)

	}

	// --- 5. Build batch INSERT data_kpi_subdetail (dari hasil parse Excel) ---
	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}

	for i, kpiItem := range req.Kpi {
		inc := fmt.Sprintf("000%d", i)
		idDetail := req.IDPengajuan + inc[len(inc)-3:]

		rows, ok := kpiSubDetails[i]
		if !ok {
			continue
		}

		for j, subRow := range rows {
			subInc := fmt.Sprintf("000%d", j)
			idSubDetail := idDetail + subInc[len(subInc)-3:]

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
				req.IDPengajuan,
				idDetail,
				idSubDetail,
				req.Tahun,
				req.Triwulan,
				kpiItem.IdKpi,
				subRow.SubKPI,
				subRow.Polarisasi, // rumus diisi polarisasi sesuai mapping lama
				"0",               // otomatis default "0"
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
				kpiItem.KeteranganProject,
				subRow.TerdapatQualifier,
			)
		}
	}

	// --- 6. Build batch INSERT data_challenge_detail ---
	challengePlaceholders := []string{}
	challengeArgs := []interface{}{}

	for _, ch := range req.ChallengeList {
		challengePlaceholders = append(challengePlaceholders, "(?, ?, ?, ?, ?, ?)")
		challengeArgs = append(challengeArgs,
			ch.IdPengajuan,
			ch.IdDetailChallenge,
			ch.Tahun,
			ch.Triwulan,
			ch.NamaChallenge,
			ch.DeskripsiChallenge,
		)
	}

	// --- 7. Build batch INSERT data_method_detail ---
	methodPlaceholders := []string{}
	methodArgs := []interface{}{}

	for _, mt := range req.MethodList {
		methodPlaceholders = append(methodPlaceholders, "(?, ?, ?, ?, ?, ?)")
		methodArgs = append(methodArgs,
			mt.IdPengajuan,
			mt.IdDetailMethod,
			mt.Tahun,
			mt.Triwulan,
			mt.NamaMethod,
			mt.DeskripsiMethod,
		)
	}

	// --- 8. Finalisasi query dengan format batch values ---
	finalQueryKpiDetail := fmt.Sprintf(queryInsertKpiDetail,
		strings.Join(kpiDetailPlaceholders, ","))

	finalQuerySubDetail := fmt.Sprintf(queryInsertKpiSubDetail,
		strings.Join(subDetailPlaceholders, ","))

	finalQueryChallenge := fmt.Sprintf(queryInsertChallengeDetail,
		strings.Join(challengePlaceholders, ","))

	finalQueryMethod := fmt.Sprintf(queryInsertMethodDetail,
		strings.Join(methodPlaceholders, ","))

	// --- 9. Eksekusi dalam 1 transaksi DB ---
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// Insert data_kpi (header)
	if err := tx.Exec(queryInsertKpi,
		req.IDPengajuan,
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
		return fmt.Errorf("gagal insert data_kpi: %w", err)
	}

	// Insert data_kpi_detail
	if len(kpiDetailPlaceholders) > 0 {
		if err := tx.Exec(finalQueryKpiDetail, kpiDetailArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_kpi_detail: %w", err)
		}
	}

	// Insert data_kpi_subdetail
	if len(subDetailPlaceholders) > 0 {
		if err := tx.Exec(finalQuerySubDetail, subDetailArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_kpi_subdetail: %w", err)
		}
	}

	// Insert data_challenge_detail
	if len(challengePlaceholders) > 0 {
		if err := tx.Exec(finalQueryChallenge, challengeArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_challenge_detail: %w", err)
		}
	}

	// Insert data_method_detail
	if len(methodPlaceholders) > 0 {
		if err := tx.Exec(finalQueryMethod, methodArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_method_detail: %w", err)
		}
	}

	// Commit jika semua berhasil
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return nil
}
