package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

	// queryGetAllDraftKpiHeader digunakan oleh GetAllDraftPenyusunanKpi dan GetDetailPenyusunanKpi
	// untuk mengambil baris header data_kpi.
	queryGetAllDraftKpiHeader = `
        SELECT
            a.id_pengajuan, a.tahun, a.triwulan, a.kostl, a.kostl_tx,
            a.orgeh, a.orgeh_tx, a.entry_user, a.entry_name, a.entry_time,
            a.approval_posisi, a.approval_list, a.status, b.status_desc,
            IFNULL(a.entry_user_realisasi, '')    entry_user_realisasi,
            IFNULL(a.entry_name_realisasi, '')    entry_name_realisasi,
            IFNULL(a.entry_time_realisasi, '')    entry_time_realisasi,
            IFNULL(a.approval_list_realisasi, '') approval_list_realisasi,
            IFNULL(a.catatan_tolakan, '')          catatan_tolakan,
            IFNULL(a.total_bobot, '')              total_bobot,
            IFNULL(a.total_pencapaian, '')         total_pencapaian,
            IFNULL(a.total_bobot_pengurang, '')    total_bobot_pengurang,
            IFNULL(a.total_pencapaian_post, '')    total_pencapaian_post,
            IFNULL(a.entry_user_validasi, '')      entry_user_validasi,
            IFNULL(a.entry_name_validasi, '')      entry_name_validasi,
            IFNULL(a.entry_time_validasi, '')      entry_time_validasi,
            IFNULL(a.approval_list_validasi, '')   approval_list_validasi,
            IFNULL(a.lampiran_validasi, '')        lampiran_validasi,
            IFNULL(a.qualifier_overall_validasi,'') qualifier_overall_validasi
        FROM data_kpi a
        INNER JOIN mst_status b ON a.status = b.id_status`

	queryCountAllDraftKpi = `
        SELECT COUNT(1)
        FROM data_kpi a
        INNER JOIN mst_status b ON a.status = b.id_status`

	queryGetKpiDetail = `
		SELECT
			a.id_detail,
			a.id_kpi, a.kpi, a.rumus,
			a.id_perspektif, b.perspektif,
			a.id_keterangan_project,
			IFNULL(c.keterangan_project, '') keterangan_project,
			IFNULL(a.lampiran_file, '')       lampiran_file
		FROM data_kpi_detail a
		INNER JOIN mst_perspektif b ON a.id_perspektif = b.id_perspektif
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ?`

	queryGetKpiSubDetail = `
		SELECT
			a.id_sub_detail,
			a.id_kpi, a.kpi, a.rumus,
			a.otomatis,
			a.bobot, a.capping,
			a.target_triwulan, a.target_kuantitatif_triwulan,
			a.target_tahunan, a.target_kuantitatif_tahunan,
			IFNULL(a.realisasi, '')                           realisasi,
			IFNULL(a.realisasi_kuantitatif, '')               realisasi_kuantitatif,
			IFNULL(a.realisasi_keterangan, '')                realisasi_keterangan,
			IFNULL(a.realisasi_validated, '')                 realisasi_validated,
			IFNULL(a.realisasi_kuantitatif_validated, '')     realisasi_kuantitatif_validated,
			IFNULL(a.validasi_keterangan, '')                 validasi_keterangan,
			IFNULL(a.pencapaian, '')                          pencapaian,
			IFNULL(a.skor, '')                                skor,
			IFNULL(a.deskripsi_glossary, '')                  deskripsi_glossary,
			IFNULL(a.item_qualifier, '')                      item_qualifier,
			IFNULL(a.deskripsi_qualifier, '')                 deskripsi_qualifier,
			IFNULL(a.target_qualifier, '')                    target_qualifier,
			IFNULL(a.id_keterangan_project, '')               id_keterangan_project,
			IFNULL(a.id_qualifier, '')                        id_qualifier,
			IFNULL(a.realisasi_qualifier, '')                 realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '')     realisasi_kuantitatif_qualifier,
			IFNULL(a.pencapaian_qualifier_validated, '')      pencapaian_qualifier_validated,
			IFNULL(a.pencapaian_post_qualifier_validated, '') pencapaian_post_qualifier_validated,
			IFNULL(c.keterangan_project, '')                  keterangan_project
		FROM data_kpi_subdetail a
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_detail = ?`

	queryGetChallengeDetail = `
		SELECT
			id_detail_challenge, tahun, triwulan,
			nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '')  realisasi_challenge,
			IFNULL(lampiran_evidence, '')    lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?`

	queryGetMethodDetail = `
		SELECT
			id_detail_method, tahun, triwulan,
			nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')   realisasi_method,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?`

	// queryGetKpiHeaderForExport digunakan oleh GetKpiExportData.
	queryGetKpiHeaderForExport = `
        SELECT kostl_tx, tahun, triwulan
        FROM data_kpi
        WHERE id_pengajuan = ?
        LIMIT 1`

	// queryGetSubDetailForExport digunakan oleh GetKpiExportData.
	queryGetSubDetailForExport = `
        SELECT
            a.kpi,
            IFNULL(CAST(a.bobot AS CHAR), '') bobot,
            IFNULL(a.target_tahunan, '')       target_tahunan,
            IFNULL(a.capping, '')              capping
        FROM data_kpi_subdetail a
        WHERE a.id_pengajuan = ?
        ORDER BY a.id_sub_detail ASC`
)

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

	// System error: gagal update ke DB
	if err := r.db.Exec(queryUpdateKpi,
		approvalPosisi, string(approvalListBytes), req.IdPengajuan,
	).Error; err != nil {
		return fmt.Errorf("gagal update data_kpi saat submit: %w", err)
	}

	return nil
}

// =============================================================================
// scanNestedKpi adalah helper yang digunakan oleh GetAllDraftPenyusunanKpi
// dan GetDetailPenyusunanKpi untuk mengisi KpiDetail, ChallengeDetail,
// dan MethodDetail ke dalam header record h.
// =============================================================================

func (r *penyusunanKpiRepo) scanNestedKpi(h *dto.GetAllDraftPenyusunanKpiResponse) error {

	// =====================================================================
	// QUERY KPI DETAIL per id_pengajuan
	// =====================================================================
	detailRows, err := r.db.Raw(queryGetKpiDetail, h.IdPengajuan).Rows()
	if err != nil {
		return fmt.Errorf("gagal mengambil kpi detail [%s]: %w", h.IdPengajuan, err)
	}

	var kpiDetails []dto.GetAllDraftKpiDetailResponse
	for detailRows.Next() {
		var d dto.GetAllDraftKpiDetailResponse
		if err := detailRows.Scan(
			&d.IdDetail,
			&d.IdKpi, &d.Kpi, &d.Rumus,
			&d.IdPerspektif, &d.Perspektif,
			&d.IdKeteranganProject,
			&d.KeteranganProject,
			&d.LampiranFile,
		); err != nil {
			detailRows.Close()
			return fmt.Errorf("gagal scan kpi detail: %w", err)
		}

		// =================================================================
		// QUERY KPI SUB DETAIL per id_detail
		// =================================================================
		subDetailRows, err := r.db.Raw(queryGetKpiSubDetail, d.IdDetail).Rows()
		if err != nil {
			detailRows.Close()
			return fmt.Errorf("gagal mengambil kpi sub detail [%s]: %w", d.IdDetail, err)
		}

		var subDetails []dto.GetAllDraftKpiSubDetailResponse
		for subDetailRows.Next() {
			var s dto.GetAllDraftKpiSubDetailResponse
			if err := subDetailRows.Scan(
				&s.IdSubDetail,
				&s.IdKpi, &s.Kpi, &s.Rumus,
				&s.Otomatis,
				&s.Bobot, &s.Capping,
				&s.TargetTriwulan, &s.TargetKuantitatifTriwulan,
				&s.TargetTahunan, &s.TargetKuantitatifTahunan,
				&s.Realisasi, &s.RealisasiKuantitatif, &s.RealisasiKeterangan,
				&s.RealisasiValidated, &s.RealisasiKuantitatifValidated,
				&s.ValidasiKeterangan,
				&s.Pencapaian, &s.Skor,
				&s.DeskripsiGlossary,
				&s.ItemQualifier, &s.DeskripsiQualifier, &s.TargetQualifier,
				&s.IdKeteranganProject,
				&s.IdQualifier,
				&s.RealisasiQualifier, &s.RealisasiKuantitatifQualifier,
				&s.PencapaianQualifierValidated, &s.PencapaianPostQualifierValidated,
				&s.KeteranganProject,
			); err != nil {
				subDetailRows.Close()
				detailRows.Close()
				return fmt.Errorf("gagal scan kpi sub detail: %w", err)
			}
			subDetails = append(subDetails, s)
		}
		subDetailRows.Close()

		d.KpiSubDetail = subDetails
		kpiDetails = append(kpiDetails, d)
	}
	detailRows.Close()
	h.KpiDetail = kpiDetails

	// =====================================================================
	// QUERY CHALLENGE DETAIL per id_pengajuan
	// =====================================================================
	challengeRows, err := r.db.Raw(queryGetChallengeDetail, h.IdPengajuan).Rows()
	if err != nil {
		return fmt.Errorf("gagal mengambil challenge detail [%s]: %w", h.IdPengajuan, err)
	}

	var challengeDetails []dto.GetAllDraftChallengeDetailResponse
	for challengeRows.Next() {
		var ch dto.GetAllDraftChallengeDetailResponse
		if err := challengeRows.Scan(
			&ch.IdDetailChallenge, &ch.Tahun, &ch.Triwulan,
			&ch.NamaChallenge, &ch.DeskripsiChallenge,
			&ch.RealisasiChallenge, &ch.LampiranEvidence,
		); err != nil {
			challengeRows.Close()
			return fmt.Errorf("gagal scan challenge detail: %w", err)
		}
		challengeDetails = append(challengeDetails, ch)
	}
	challengeRows.Close()
	h.ChallengeDetail = challengeDetails

	// =====================================================================
	// QUERY METHOD DETAIL per id_pengajuan
	// =====================================================================
	methodRows, err := r.db.Raw(queryGetMethodDetail, h.IdPengajuan).Rows()
	if err != nil {
		return fmt.Errorf("gagal mengambil method detail [%s]: %w", h.IdPengajuan, err)
	}

	var methodDetails []dto.GetAllDraftMethodDetailResponse
	for methodRows.Next() {
		var mt dto.GetAllDraftMethodDetailResponse
		if err := methodRows.Scan(
			&mt.IdDetailMethod, &mt.Tahun, &mt.Triwulan,
			&mt.NamaMethod, &mt.DeskripsiMethod,
			&mt.RealisasiMethod, &mt.LampiranEvidence,
		); err != nil {
			methodRows.Close()
			return fmt.Errorf("gagal scan method detail: %w", err)
		}
		methodDetails = append(methodDetails, mt)
	}
	methodRows.Close()
	h.MethodDetail = methodDetails

	return nil
}

// =============================================================================
// GET ALL DRAFT — list dengan filter, pagination, dan nested detail
// =============================================================================

func (r *penyusunanKpiRepo) GetAllDraftPenyusunanKpi(
	req *dto.GetAllDraftPenyusunanKpiRequest,
) ([]*dto.GetAllDraftPenyusunanKpiResponse, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status IN (70, 71)",
		"a.entry_user = ?",
	}
	args := []interface{}{req.EntryUser}

	// =========================================================================
	// Kondisi opsional dari request body
	// =========================================================================
	if req.Divisi != "" {
		conditions = append(conditions, "a.kostl = ?")
		args = append(args, req.Divisi)
	}
	if req.Tahun != "" {
		conditions = append(conditions, "a.tahun = ?")
		args = append(args, req.Tahun)
	}
	if req.Triwulan != "" {
		conditions = append(conditions, "a.triwulan = ?")
		args = append(args, req.Triwulan)
	}
	if req.Status != "" {
		conditions = append(conditions, "a.status = ?")
		args = append(args, req.Status)
	}

	whereClause := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// COUNT TOTAL RECORDS
	// =========================================================================
	var total int64
	countQuery := queryCountAllDraftKpi + whereClause
	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	// =========================================================================
	// PAGINATION
	// =========================================================================
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// =========================================================================
	// QUERY HEADER
	// =========================================================================
	mainQuery := queryGetAllDraftKpiHeader + whereClause + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	headerArgs := append(args, limit, offset)

	headerRows, err := r.db.Raw(mainQuery, headerArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data header KPI: %w", err)
	}
	defer headerRows.Close()

	var results []*dto.GetAllDraftPenyusunanKpiResponse

	for headerRows.Next() {
		var h dto.GetAllDraftPenyusunanKpiResponse

		if err := headerRows.Scan(
			&h.IdPengajuan, &h.Tahun, &h.Triwulan,
			&h.Kostl, &h.KostlTx,
			&h.Orgeh, &h.OrgehTx,
			&h.EntryUser, &h.EntryName, &h.EntryTime,
			&h.ApprovalPosisi, &h.ApprovalList,
			&h.Status, &h.StatusDesc,
			&h.EntryUserRealisasi, &h.EntryNameRealisasi, &h.EntryTimeRealisasi,
			&h.ApprovalListRealisasi,
			&h.CatatanTolakan,
			&h.TotalBobot, &h.TotalPencapaian,
			&h.TotalBobotPengurang, &h.TotalPencapaianPost,
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi,
			&h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		if err := r.scanNestedKpi(&h); err != nil {
			return nil, 0, err
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET DETAIL — 1 record berdasarkan id_pengajuan (tanpa filter status/entry_user)
// =============================================================================

func (r *penyusunanKpiRepo) GetDetailPenyusunanKpi(
	req *dto.GetDetailPenyusunanKpiRequest,
) (*dto.GetAllDraftPenyusunanKpiResponse, error) {

	detailQuery := queryGetAllDraftKpiHeader + " WHERE a.id_pengajuan = ?"

	headerRows, err := r.db.Raw(detailQuery, req.IdPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil detail KPI: %w", err)
	}
	defer headerRows.Close()

	if !headerRows.Next() {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	var h dto.GetAllDraftPenyusunanKpiResponse
	if err := headerRows.Scan(
		&h.IdPengajuan, &h.Tahun, &h.Triwulan,
		&h.Kostl, &h.KostlTx,
		&h.Orgeh, &h.OrgehTx,
		&h.EntryUser, &h.EntryName, &h.EntryTime,
		&h.ApprovalPosisi, &h.ApprovalList,
		&h.Status, &h.StatusDesc,
		&h.EntryUserRealisasi, &h.EntryNameRealisasi, &h.EntryTimeRealisasi,
		&h.ApprovalListRealisasi,
		&h.CatatanTolakan,
		&h.TotalBobot, &h.TotalPencapaian,
		&h.TotalBobotPengurang, &h.TotalPencapaianPost,
		&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
		&h.ApprovalListValidasi,
		&h.LampiranValidasi,
		&h.QualifierOverallValidasi,
	); err != nil {
		return nil, fmt.Errorf("gagal scan detail KPI: %w", err)
	}

	if err := r.scanNestedKpi(&h); err != nil {
		return nil, err
	}

	return &h, nil
}

// =============================================================================
// GET EXPORT DATA — digunakan bersama oleh get-csv dan get-pdf
// =============================================================================

func (r *penyusunanKpiRepo) GetKpiExportData(
	idPengajuan string,
) (*dto.KpiExportData, error) {

	type kpiHeader struct {
		KostlTx  string `gorm:"column:kostl_tx"`
		Tahun    string `gorm:"column:tahun"`
		Triwulan string `gorm:"column:triwulan"`
	}
	var header kpiHeader
	if err := r.db.Raw(queryGetKpiHeaderForExport, idPengajuan).Scan(&header).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil header KPI: %w", err)
	}
	if header.KostlTx == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
		}
	}

	type subDetailRaw struct {
		Kpi           string `gorm:"column:kpi"`
		Bobot         string `gorm:"column:bobot"`
		TargetTahunan string `gorm:"column:target_tahunan"`
		Capping       string `gorm:"column:capping"`
	}
	var rawRows []subDetailRaw
	if err := r.db.Raw(queryGetSubDetailForExport, idPengajuan).Scan(&rawRows).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil sub detail KPI untuk ekspor: %w", err)
	}

	rows := make([]dto.KpiSubDetailExportRow, 0, len(rawRows))
	for i, row := range rawRows {
		rows = append(rows, dto.KpiSubDetailExportRow{
			No:            i + 1,
			KpiNama:       row.Kpi,
			Bobot:         row.Bobot,
			TargetTahunan: row.TargetTahunan,
			Capping:       row.Capping,
		})
	}

	return &dto.KpiExportData{
		NamaDivisi: header.KostlTx,
		Tahun:      header.Tahun,
		Triwulan:   header.Triwulan,
		Rows:       rows,
	}, nil
}
