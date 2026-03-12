package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/utils"
	customErrors "permen_api/errors"
)

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

	// queryInsertKpi digunakan oleh ValidatePenyusunanKpi.
	queryInsertKpi = `
		INSERT INTO data_kpi 
			(id_pengajuan, tahun, triwulan, kostl, kostl_tx, orgeh, orgeh_tx, 
			 entry_user, entry_name, entry_time, approval_posisi, approval_list, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// queryUpdateKpi digunakan oleh CreatePenyusunanKpi untuk mengisi approval dan mengubah status.
	queryUpdateKpi = `
		UPDATE data_kpi 
		SET approval_posisi = ?, approval_list = ?, status = 0
		WHERE id_pengajuan = ?`

	queryCheckExistIdPengajuan = `
		SELECT COUNT(id_pengajuan) 
		FROM data_kpi 
		WHERE id_pengajuan = ?`

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

	queryLookupSubKpi = `
		SELECT id_kpi, kpi, rumus
		FROM mst_kpi
		WHERE LOWER(kpi) = LOWER(?)
		LIMIT 1`

	queryLookupPolarisasi = `
		SELECT id_polarisasi
		FROM mst_polarisasi
		WHERE LOWER(polarisasi) = LOWER(?)
		LIMIT 1`
)

// =============================================================================
// HELPER
// =============================================================================

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

func generateIDDetail(idPengajuan string, index int) string {
	return fmt.Sprintf("%sP%03d", idPengajuan, index+1)
}

func generateIDSubDetail(idPengajuan string, globalIndex int) string {
	return fmt.Sprintf("%sC%03d", idPengajuan, globalIndex)
}

// =============================================================================
// LOOKUP
// =============================================================================

func (r *penyusunanKpiRepo) LookupSubKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error) {
	row := r.db.Raw(queryLookupSubKpi, subKpiText).Row()
	if scanErr := row.Scan(&idKpi, &kpiFromDB, &rumus); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			// Sub KPI tidak ditemukan di master → dianggap KPI lain (id = "0")
			return "0", subKpiText, "", nil
		}
		return "0", subKpiText, "", fmt.Errorf("gagal lookup mst_kpi untuk Sub KPI '%s': %w", subKpiText, scanErr)
	}
	return idKpi, kpiFromDB, rumus, nil
}

func (r *penyusunanKpiRepo) LookupPolarisasi(polarisasiText string) (idPolarisasi string, err error) {
	row := r.db.Raw(queryLookupPolarisasi, polarisasiText).Row()
	if scanErr := row.Scan(&idPolarisasi); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			// User error: nilai polarisasi dari Excel tidak ada di master
			return "", &customErrors.BadRequestError{
				Message: fmt.Sprintf(
					"polarisasi '%s' tidak ditemukan di master polarisasi. Nilai yang valid: 'Maximize' atau 'Minimize'",
					polarisasiText,
				),
			}
		}
		// System error: query DB gagal
		return "", fmt.Errorf("gagal lookup master polarisasi untuk polarisasi '%s': %w", polarisasiText, scanErr)
	}
	return idPolarisasi, nil
}

// =============================================================================
// VALIDATE — simpan data KPI tanpa approval
// =============================================================================

func (r *penyusunanKpiRepo) ValidatePenyusunanKpi(
	req *dto.ValidatePenyusunanKpiRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) (string, error) {

	idPengajuan := utils.GenerateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	var countExist int
	if err := r.db.Raw(queryCheckExistKpi, req.Tahun, req.Triwulan, req.Kostl).
		Scan(&countExist).Error; err != nil {
		return "", fmt.Errorf("gagal mengecek data KPI: %w", err)
	}
	if countExist > 0 {
		return "", &customErrors.BadRequestError{
			Message: fmt.Sprintf(
				"data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
				req.Tahun, req.Triwulan, req.Kostl,
			),
		}
	}

	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Kostl).Row().Scan(&orgeh, &orgehTx)

	// status 70 = draft
	var statusKpi interface{} = 70

	kpiDetailPlaceholders := []string{}
	kpiDetailArgs := []interface{}{}
	idDetailMap := make(map[int]string)

	for i, kpiItem := range req.Kpi {
		idDetail := utils.GenerateIDDetail(idPengajuan, i)
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
			"",
		)
	}

	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}
	subCounter := 1

	for i := range req.Kpi {
		rows, ok := kpiSubDetails[i]
		if !ok {
			continue
		}

		idDetail := idDetailMap[i]

		for _, subRow := range rows {
			idSubDetail := utils.GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			itemQualifier, deskripsiQualifier, targetQualifier := "", "", ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

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

			subDetailPlaceholders = append(subDetailPlaceholders,
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			subDetailArgs = append(subDetailArgs,
				idPengajuan,
				idDetail,
				idSubDetail,
				req.Tahun,
				req.Triwulan,
				subRow.IdSubKpi,
				subRow.SubKPI,
				subRow.IdPolarisasi,
				subRow.Otomatis,
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
				"",
				subRow.TerdapatQualifier,
				result,
				deskripsiResult,
				process,
				deskripsiProcess,
				context,
				deskripsiContext,
			)
		}
	}

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

	tx := r.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// approval_posisi dan approval_list dikosongkan — diisi saat CreatePenyusunanKpi
	if err := tx.Exec(queryInsertKpi,
		idPengajuan, req.Tahun, req.Triwulan, req.Kostl, req.KostlTx,
		orgeh, orgehTx, req.EntryUser, req.EntryName, req.EntryTime,
		"", "[]", statusKpi,
	).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("gagal insert data_kpi: %w", err)
	}

	if len(kpiDetailPlaceholders) > 0 {
		queryDetail := fmt.Sprintf(queryInsertKpiDetail, strings.Join(kpiDetailPlaceholders, ", "))
		if err := tx.Exec(queryDetail, kpiDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_detail: %w", err)
		}
	}

	if len(subDetailPlaceholders) > 0 {
		querySubDetail := fmt.Sprintf(queryInsertKpiSubDetail, strings.Join(subDetailPlaceholders, ", "))
		if err := tx.Exec(querySubDetail, subDetailArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_kpi_subdetail: %w", err)
		}
	}

	if len(challengePlaceholders) > 0 {
		queryChallenge := fmt.Sprintf(queryInsertChallengeDetail, strings.Join(challengePlaceholders, ", "))
		if err := tx.Exec(queryChallenge, challengeArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_challenge_detail: %w", err)
		}
	}

	if len(methodPlaceholders) > 0 {
		queryMethod := fmt.Sprintf(queryInsertMethodDetail, strings.Join(methodPlaceholders, ", "))
		if err := tx.Exec(queryMethod, methodArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_method_detail: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return idPengajuan, nil
}

// =============================================================================
// CREATE — update approval pada data KPI yang sudah ada
// =============================================================================

func (r *penyusunanKpiRepo) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
) error {

	// Cek apakah idPengajuan benar-benar ada di DB
	var countExist int
	if err := r.db.Raw(queryCheckExistIdPengajuan, req.IdPengajuan).
		Scan(&countExist).Error; err != nil {
		return fmt.Errorf("gagal mengecek id_pengajuan: %w", err)
	}
	if countExist == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	// Ambil userid signer sebagai approval_posisi (posisi = "SIGNER")
	approvalPosisi := ""
	for _, a := range req.ApprovalList {
		if strings.EqualFold(a.Posisi, "SIGNER") {
			approvalPosisi = a.Userid
			break
		}
	}
	// Fallback: jika tidak ada SIGNER, pakai userid pertama
	if approvalPosisi == "" && len(req.ApprovalList) > 0 {
		approvalPosisi = req.ApprovalList[0].Userid
	}

	approvalListBytes, err := json.Marshal(req.ApprovalList)
	if err != nil {
		// System error: seharusnya tidak terjadi karena struct sudah tervalidasi
		return fmt.Errorf("gagal serialize approval_list: %w", err)
	}
	approvalListStr := string(approvalListBytes)

	// System error: gagal update ke DB
	if err := r.db.Exec(queryUpdateKpi,
		approvalPosisi, approvalListStr, req.IdPengajuan,
	).Error; err != nil {
		return fmt.Errorf("gagal update data_kpi saat submit: %w", err)
	}

	return nil
}
