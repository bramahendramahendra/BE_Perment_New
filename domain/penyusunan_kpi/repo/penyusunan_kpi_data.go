package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/model"
	customErrors "permen_api/errors"
	"permen_api/pkg/idgen"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Check
	// =============================================================================

	// Use func : ValidatePenyusunanKpi
	queryCheckExistPenyusunan = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE tahun = ? AND triwulan = ? AND kostl = ?`

	// Use func : GetExistPenyusunanStatus
	queryGetExistPenyusunanStatus = `
		SELECT id_pengajuan, status
		FROM data_kpi
		WHERE tahun = ? AND triwulan = ? AND kostl = ?
		LIMIT 1`

	// Use func : RevisionPenyusunanKpi, CreatePenyusunanKpi
	queryCheckExistIdPengajuan = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND kostl = ? AND tahun = ? AND triwulan = ?`

	// Use func : GetAllApprovalPenyusunanKpi, GetAllTolakanPenyusunanKpi, GetAllDaftarPenyusunanKpi, GetAllDaftarApprovalPenyusunanKpi
	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// queryGetKpiBaseData digunakan oleh: SubmitPenyusunanKpi (kostl_tx),
	// ApprovePenyusunanKpi (entry_user), dan GetKpiExportData (kostl_tx, tahun, triwulan).
	queryGetKpiBaseData = `
		SELECT kostl_tx, tahun, triwulan, entry_user
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// queryGetKpiHeader digunakan oleh: RevisionPenyusunanKpi (service) untuk mengambil
	// tahun, triwulan, kostl, kostl_tx berdasarkan id_pengajuan.
	queryGetKpiHeader = `
		SELECT a.tahun, a.triwulan, a.kostl, a.kostl_tx, a.status, b.status_desc, a.entry_user, a.entry_name
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ?
		LIMIT 1`

	// Use func : ValidatePenyusunanKpi
	queryGetOrgeh = `
		SELECT orgeh, orgeh_tx 
		FROM user 
		WHERE kostl = ? 
		ORDER BY HILFM ASC 
		LIMIT 1`

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

	queryGetDataKpi = `
		SELECT
			a.id_pengajuan, a.tahun, a.triwulan, a.kostl, a.kostl_tx,
			a.orgeh, a.orgeh_tx, a.entry_user, a.entry_name, a.entry_time,
			a.approval_posisi, a.approval_list, a.status, b.status_desc,
			IFNULL(a.entry_user_realisasi, '')     entry_user_realisasi,
			IFNULL(a.entry_name_realisasi, '')     entry_name_realisasi,
			IFNULL(a.entry_time_realisasi, '')     entry_time_realisasi,
			IFNULL(a.approval_list_realisasi, '')  approval_list_realisasi,
			IFNULL(a.catatan_tolakan, '')           catatan_tolakan,
			IFNULL(a.total_bobot, '')               total_bobot,
			IFNULL(a.total_pencapaian, '')          total_pencapaian,
			IFNULL(a.total_bobot_pengurang, '')     total_bobot_pengurang,
			IFNULL(a.total_pencapaian_post, '')     total_pencapaian_post,
			IFNULL(a.entry_user_validasi, '')       entry_user_validasi,
			IFNULL(a.entry_name_validasi, '')       entry_name_validasi,
			IFNULL(a.entry_time_validasi, '')       entry_time_validasi,
			IFNULL(a.approval_list_validasi, '')    approval_list_validasi,
			IFNULL(a.lampiran_validasi, '')         lampiran_validasi,
			IFNULL(a.qualifier_overall_validasi,'') qualifier_overall_validasi
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	queryGetDataKpiDetail = `
		SELECT
			a.id_detail,
			a.id_kpi, a.kpi, a.rumus,
			IFNULL(a.id_perspektif, '')         id_perspektif,
			IFNULL(b.perspektif, '')            perspektif,
			IFNULL(a.id_keterangan_project, '') id_keterangan_project,
			IFNULL(c.keterangan_project, '')    keterangan_project
		FROM data_kpi_detail a
		LEFT JOIN mst_perspektif b ON a.id_perspektif = b.id_perspektif
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ?`

	// queryGetDataKpiSubDetailPenyusunan digunakan oleh GetDetailPenyusunanKpi (versi ringan).
	// JOIN ke mst_polarisasi via a.rumus untuk mendapatkan polarisasi dan id_polarisasi.
	queryGetDataKpiSubDetail = `
		SELECT
			a.id_sub_detail,
			a.id_kpi, a.kpi, a.rumus,
			a.otomatis,
			a.bobot, a.capping,
			a.target_triwulan,  a.target_kuantitatif_triwulan,
			a.target_tahunan,   a.target_kuantitatif_tahunan,
			IFNULL(a.deskripsi_glossary, '')    deskripsi_glossary,
			IFNULL(a.rumus, '')                 id_polarisasi,
			IFNULL(p.polarisasi, '')            polarisasi,
			IFNULL(a.id_qualifier, '')          id_qualifier,
			IFNULL(a.item_qualifier, '')        item_qualifier,
			IFNULL(a.deskripsi_qualifier, '')   deskripsi_qualifier,
			IFNULL(a.target_qualifier, '')      target_qualifier,
			IFNULL(a.id_keterangan_project, '') id_keterangan_project,
			IFNULL(c.keterangan_project, '')    keterangan_project
		FROM data_kpi_subdetail a
		LEFT JOIN mst_polarisasi p ON a.rumus = p.id_polarisasi
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_detail = ?`

	queryGetDataResultDetail = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result
		FROM data_result_detail
		WHERE id_pengajuan = ?`

	queryGetDataProcessDetail = `
		SELECT
			id_detail_method,
			nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')   realisasi_method,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?`

	queryGetDataContextDetail = `
		SELECT
			id_detail_challenge,
			nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '')  realisasi_challenge,
			IFNULL(lampiran_evidence, '')    lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?`

	queryGetKpiBaseDataForExport = `
		SELECT kostl_tx, tahun, triwulan, entry_user
		FROM data_kpi
		WHERE id_pengajuan = ? AND kostl = ? AND tahun = ? AND triwulan = ?
		LIMIT 1`

	queryGetSubDetailForExport = `
		SELECT
			a.kpi,
			IFNULL(CAST(a.bobot AS CHAR), '') bobot,
			IFNULL(a.target_tahunan, '')       target_tahunan,
			IFNULL(a.capping, '')              capping
		FROM data_kpi_subdetail a
		INNER JOIN data_kpi b ON a.id_pengajuan = b.id_pengajuan
		WHERE a.id_pengajuan = ? AND b.kostl = ? AND b.tahun = ? AND b.triwulan = ?
		ORDER BY a.id_sub_detail ASC`

	queryCheckApprovalPenyusunan = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 0 AND approval_posisi = ? AND id_pengajuan = ?`

	queryGetApprovalListJSON = `
		SELECT approval_list FROM data_kpi
		WHERE status = 0 AND approval_posisi = ? AND id_pengajuan = ?`

	queryGetCatatanTolakan = `
		SELECT IFNULL(catatan_tolakan, '') FROM data_kpi
		WHERE id_pengajuan = ? LIMIT 1`
	// =============================================================================
	// Insert
	// =============================================================================
	// queryInsertKpi digunakan oleh ValidatePenyusunanKpi.
	queryInsertKpi = `
		INSERT INTO data_kpi 
			(id_pengajuan, tahun, triwulan, kostl, kostl_tx, orgeh, orgeh_tx, 
			 entry_user, entry_name, entry_time, approval_posisi, approval_list, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// queryInsertKpiDetail: id_perspektif diisi NULL karena sudah tidak digunakan.
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

	queryInsertResultDetail = `
		INSERT INTO data_result_detail
			(id_pengajuan, id_detail_result, tahun, triwulan,
			nama_result, deskripsi_result)
		VALUES %s`

	queryInsertProcessDetail = `
		INSERT INTO data_method_detail 
			(id_pengajuan, id_detail_method, tahun, triwulan, 
			 nama_method, deskripsi_method) 
		VALUES %s`

	queryInsertContextDetail = `
		INSERT INTO data_challenge_detail 
			(id_pengajuan, id_detail_challenge, tahun, triwulan, 
			 nama_challenge, deskripsi_challenge) 
		VALUES %s`

	// =============================================================================
	// Update
	// =============================================================================
	// queryUpdateKpi digunakan oleh CreatePenyusunanKpi untuk mengisi approval dan mengubah status.
	queryUpdateKpi = `
		UPDATE data_kpi 
		SET approval_posisi = ?, approval_list = ?, status = 0
		WHERE id_pengajuan = ?`

	queryApproveChainPenyusunan = `
		UPDATE data_kpi SET approval_posisi = ?, approval_list = ? WHERE id_pengajuan = ?`

	queryApproveFinalPenyusunan = `
		UPDATE data_kpi SET status = 2, approval_list = ? WHERE id_pengajuan = ?`

	queryRejectPenyusunan = `
		UPDATE data_kpi SET status = 1, approval_list = ?, catatan_tolakan = ? WHERE id_pengajuan = ?`

	// queryGetApprovalForRevision digunakan oleh RevisionPenyusunanKpi untuk membaca
	// approval_posisi dan approval_list sebelum dikosongkan.
	queryGetApprovalForRevision = `
		SELECT approval_posisi, approval_list FROM data_kpi
		WHERE id_pengajuan = ? LIMIT 1`

	// queryUpdateKpiRevision digunakan oleh RevisionPenyusunanKpi untuk update header data_kpi.
	// status = 0 → langsung ke approval (tidak perlu draft lagi).
	queryUpdateKpiRevision = `
		UPDATE data_kpi
		SET entry_time = ?, approval_list = ?, approval_posisi = ?, status = 0
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Delete
	// =============================================================================

	// Query-query DELETE berikut digunakan oleh RevisionPenyusunanKpi dan ValidatePenyusunanKpi (replace draft).
	queryDeleteKpiHeader     = `DELETE FROM data_kpi WHERE id_pengajuan = ?`
	queryDeleteKpiDetail     = `DELETE FROM data_kpi_detail WHERE id_pengajuan = ?`
	queryDeleteKpiSubDetail  = `DELETE FROM data_kpi_subdetail WHERE id_pengajuan = ?`
	queryDeleteResultDetail  = `DELETE FROM data_result_detail WHERE id_pengajuan = ?`
	queryDeleteProcessDetail = `DELETE FROM data_method_detail WHERE id_pengajuan = ?`
	queryDeleteContextDetail = `DELETE FROM data_challenge_detail WHERE id_pengajuan = ?`
)

// =============================================================================
// CHECK
// =============================================================================

func (r *penyusunanKpiRepo) CheckExistPenyusunan(tahun, triwulan, kostl string) (bool, error) {
	var count int
	if err := r.db.Raw(queryCheckExistPenyusunan, tahun, triwulan, kostl).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal mengecek data Penyusunan KPI: %w", err)
	}
	return count > 0, nil
}

func (r *penyusunanKpiRepo) GetExistPenyusunanStatus(tahun, triwulan, kostl string) (idPengajuan string, status int, found bool, err error) {
	row := r.db.Raw(queryGetExistPenyusunanStatus, tahun, triwulan, kostl).Row()
	if scanErr := row.Scan(&idPengajuan, &status); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return "", 0, false, nil
		}
		return "", 0, false, fmt.Errorf("gagal mengecek data Penyusunan KPI: %w", scanErr)
	}
	return idPengajuan, status, true, nil
}

func (r *penyusunanKpiRepo) CheckExistIdPengajuan(idPengajuan, kostl, tahun, triwulan string) (bool, error) {
	var count int
	if err := r.db.Raw(queryCheckExistIdPengajuan, idPengajuan, kostl, tahun, triwulan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal mengecek id_pengajuan: %w", err)
	}
	return count > 0, nil
}

func (r *penyusunanKpiRepo) CheckApprovalExists(user, idPengajuan string) (bool, error) {
	var count int64
	if err := r.db.Raw(queryCheckApprovalPenyusunan, user, idPengajuan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	return count > 0, nil
}

// =============================================================================
// LOOKUP
// =============================================================================

func (r *penyusunanKpiRepo) LookupKpiMaster(subKpiText string) (idKpi, kpiFromDB, rumus string, err error) {
	row := r.db.Raw(queryLookupSubKpi, subKpiText).Row()
	if scanErr := row.Scan(&idKpi, &kpiFromDB, &rumus); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			// KPI/Sub KPI tidak ditemukan di master → dianggap KPI lain (id = "0"), teks asli dikembalikan
			return "0", subKpiText, "", nil
		}
		return "0", subKpiText, "", fmt.Errorf("gagal lookup mst_kpi untuk '%s': %w", subKpiText, scanErr)
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
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
	idLama string,
) (string, error) {

	idPengajuan := idgen.GenerateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Kostl).Row().Scan(&orgeh, &orgehTx)

	// status 70 = draft
	var statusKpi interface{} = 70

	// =========================================================================
	// Build INSERT data_kpi_detail
	// Perubahan: id_perspektif diisi NULL (sudah tidak digunakan)
	// =========================================================================
	kpiDetailPlaceholders := []string{}
	kpiDetailArgs := []interface{}{}
	idDetailMap := make(map[int]string) // kpiIndex → idDetail

	for i, kpiRow := range kpiRows {
		idDetail := idgen.GenerateIDDetail(idPengajuan, i)
		idDetailMap[kpiRow.KpiIndex] = idDetail

		kpiDetailPlaceholders = append(kpiDetailPlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		kpiDetailArgs = append(kpiDetailArgs,
			idPengajuan,
			idDetail,
			req.Tahun,
			req.Triwulan,
			kpiRow.IdKpi,
			kpiRow.Kpi,
			kpiRow.Rumus,
			nil, // id_perspektif = NULL (sudah tidak digunakan)
			nil, // id_keterangan_project
		)
	}

	// =========================================================================
	// Build INSERT data_kpi_subdetail
	// =========================================================================
	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}
	subCounter := 1

	for _, kpiRow := range kpiRows {
		rows, ok := kpiSubDetails[kpiRow.KpiIndex]
		if !ok {
			continue
		}

		idDetail := idDetailMap[kpiRow.KpiIndex]

		for _, subRow := range rows {
			idSubDetail := idgen.GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			itemQualifier, deskripsiQualifier, targetQualifier := "", "", ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			subDetailPlaceholders = append(subDetailPlaceholders,
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			subDetailArgs = append(subDetailArgs,
				idPengajuan,
				idDetail, idSubDetail, req.Tahun, req.Triwulan,
				subRow.IdSubKpi, subRow.SubKPI, subRow.IdPolarisasi, subRow.Otomatis, subRow.Bobot,
				subRow.Capping, subRow.TargetTriwulan, subRow.TargetKuantitatifTriwulan,
				subRow.TargetTahunan, subRow.TargetKuantitatifTahunan, subRow.Glossary,
				itemQualifier, deskripsiQualifier, targetQualifier,
				"", subRow.TerdapatQualifier,
			)
		}
	}

	// =========================================================================
	// Build INSERT data_result_detail
	// resultList sudah dibangun oleh service dari kolom P/Q Excel.
	// Hanya berisi data jika triwulan TW2 atau TW4.
	// =========================================================================
	resultPlaceholders := []string{}
	resultArgs := []interface{}{}
	for _, rs := range resultList {
		resultPlaceholders = append(resultPlaceholders, "(?, ?, ?, ?, ?, ?)")
		resultArgs = append(resultArgs,
			idPengajuan,
			rs.IdDetailResult,
			req.Tahun,
			req.Triwulan,
			rs.NamaResult,
			rs.DeskripsiResult,
		)
	}

	// =========================================================================
	// Build INSERT data_method_detail
	// processList sudah dibangun oleh service dari kolom R/S Excel.
	// Hanya berisi data jika triwulan TW2 atau TW4.
	// =========================================================================
	processPlaceholders := []string{}
	processArgs := []interface{}{}
	for _, mt := range processList {
		processPlaceholders = append(processPlaceholders, "(?, ?, ?, ?, ?, ?)")
		processArgs = append(processArgs,
			idPengajuan,
			mt.IdDetailProcess,
			req.Tahun,
			req.Triwulan,
			mt.NamaProcess,
			mt.DeskripsiProcess,
		)
	}

	// =========================================================================
	// Build INSERT data_challenge_detail
	// contextList sudah dibangun oleh service dari kolom T/U Excel.
	// Hanya berisi data jika triwulan TW2 atau TW4.
	// =========================================================================
	contextPlaceholders := []string{}
	contextArgs := []interface{}{}
	for _, ch := range contextList {
		contextPlaceholders = append(contextPlaceholders, "(?, ?, ?, ?, ?, ?)")
		contextArgs = append(contextArgs,
			idPengajuan,
			ch.IdDetailContext,
			req.Tahun,
			req.Triwulan,
			ch.NamaContext,
			ch.DeskripsiContext,
		)
	}

	// =========================================================================
	// Eksekusi semua INSERT dalam satu transaksi
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return "", fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// Jika ada draft lama (status=70), hapus seluruh datanya sebelum insert baru
	if idLama != "" {
		for _, q := range []struct {
			query string
			desc  string
		}{
			{queryDeleteKpiDetail, "data_kpi_detail"},
			{queryDeleteKpiSubDetail, "data_kpi_subdetail"},
			{queryDeleteResultDetail, "data_result_detail"},
			{queryDeleteProcessDetail, "data_method_detail"},
			{queryDeleteContextDetail, "data_challenge_detail"},
			{queryDeleteKpiHeader, "data_kpi"},
		} {
			if err := tx.Exec(q.query, idLama).Error; err != nil {
				tx.Rollback()
				return "", fmt.Errorf("gagal menghapus draft lama (%s): %w", q.desc, err)
			}
		}
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

	if len(resultPlaceholders) > 0 {
		queryResult := fmt.Sprintf(queryInsertResultDetail, strings.Join(resultPlaceholders, ", "))
		if err := tx.Exec(queryResult, resultArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_result_detail: %w", err)
		}
	}

	if len(processPlaceholders) > 0 {
		queryProcess := fmt.Sprintf(queryInsertProcessDetail, strings.Join(processPlaceholders, ", "))
		if err := tx.Exec(queryProcess, processArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_method_detail: %w", err)
		}
	}

	if len(contextPlaceholders) > 0 {
		queryContext := fmt.Sprintf(queryInsertContextDetail, strings.Join(contextPlaceholders, ", "))
		if err := tx.Exec(queryContext, contextArgs...).Error; err != nil {
			tx.Rollback()
			return "", fmt.Errorf("gagal insert data_challenge_detail: %w", err)
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

	// Ambil userid pertama dari ApprovalList sebagai approval_posisi
	approvalPosisi := ""
	if len(req.ApprovalList) > 0 {
		approvalPosisi = req.ApprovalList[0].Userid
	}

	approvalListBytes, err := json.Marshal(req.ApprovalList)
	if err != nil {
		return fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	// Ambil kostl_tx sebelum transaksi dimulai
	var kpiBase struct {
		KostlTx string `gorm:"column:kostl_tx"`
	}
	if err := r.db.Raw(queryGetKpiBaseData, req.IdPengajuan).Scan(&kpiBase).Error; err != nil {
		return fmt.Errorf("gagal mengambil kostl_tx: %w", err)
	}
	kostlTx := kpiBase.KostlTx

	// =========================================================================
	// Jalankan dalam transaksi agar update + notif atomic
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryUpdateKpi,
		approvalPosisi, string(approvalListBytes), req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update data_kpi saat submit: %w", err)
	}

	if err := notif.Insert(
		tx,
		req.IdPengajuan,
		kostlTx,
		req.EntryUser,
		approvalPosisi,
		"approval_penyusunan",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return nil
}

// =============================================================================
// Revision — update data kpi
// =============================================================================

func (r *penyusunanKpiRepo) RevisionPenyusunanKpi(
	req *dto.RevisionPenyusunanKpiRequest,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
) error {
	// =========================================================================
	// Ambil approval_posisi dan approval_list lama, lalu kosongkan entry
	// yang posisi-nya sama dengan approval_posisi (yaitu approver yang reject).
	// =========================================================================
	var approvalPosisi string
	var approvalListRaw []byte
	row := r.db.Raw(queryGetApprovalForRevision, req.IdPengajuan).Row()
	if err := row.Scan(&approvalPosisi, &approvalListRaw); err != nil {
		return fmt.Errorf("gagal membaca approval data: %w", err)
	}

	var approvalList []dto.ApprovalUserDetail
	if err := json.Unmarshal(approvalListRaw, &approvalList); err != nil {
		return fmt.Errorf("gagal parse approval_list: %w", err)
	}

	// Reset semua entry approval_list ke posisi awal
	for i := range approvalList {
		approvalList[i].Status = ""
		approvalList[i].Keterangan = ""
		approvalList[i].Waktu = ""
	}

	// approval_posisi dikembalikan ke approver pertama dalam list
	firstApprovalPosisi := approvalPosisi
	if len(approvalList) > 0 {
		firstApprovalPosisi = approvalList[0].Userid
	}

	updatedApprovalListBytes, err := json.Marshal(approvalList)
	if err != nil {
		return fmt.Errorf("gagal serialize approval_list: %w", err)
	}
	updatedApprovalList := string(updatedApprovalListBytes)

	// =========================================================================
	// Build INSERT data_kpi_detail
	// =========================================================================
	kpiDetailPlaceholders := []string{}
	kpiDetailArgs := []interface{}{}
	idDetailMap := make(map[int]string) // kpiIndex → idDetail

	for i, kpiRow := range kpiRows {
		idDetail := idgen.GenerateIDDetail(req.IdPengajuan, i)
		idDetailMap[kpiRow.KpiIndex] = idDetail

		kpiDetailPlaceholders = append(kpiDetailPlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		kpiDetailArgs = append(kpiDetailArgs,
			req.IdPengajuan,
			idDetail,
			req.Tahun,
			req.Triwulan,
			kpiRow.IdKpi,
			kpiRow.Kpi,
			kpiRow.Rumus,
			nil, // id_perspektif = NULL (sudah tidak digunakan)
			nil, // id_keterangan_project
		)
	}

	// =========================================================================
	// Build INSERT data_kpi_subdetail
	// =========================================================================
	subDetailPlaceholders := []string{}
	subDetailArgs := []interface{}{}
	subCounter := 1

	for _, kpiRow := range kpiRows {
		rows, ok := kpiSubDetails[kpiRow.KpiIndex]
		if !ok {
			continue
		}

		idDetail := idDetailMap[kpiRow.KpiIndex]

		for _, subRow := range rows {
			idSubDetail := idgen.GenerateIDSubDetail(req.IdPengajuan, subCounter)
			subCounter++

			itemQualifier, deskripsiQualifier, targetQualifier := "", "", ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				itemQualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			subDetailPlaceholders = append(subDetailPlaceholders,
				"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			subDetailArgs = append(subDetailArgs,
				req.IdPengajuan,
				idDetail, idSubDetail, req.Tahun, req.Triwulan,
				subRow.IdSubKpi, subRow.SubKPI, subRow.IdPolarisasi, subRow.Otomatis, subRow.Bobot,
				subRow.Capping, subRow.TargetTriwulan, subRow.TargetKuantitatifTriwulan,
				subRow.TargetTahunan, subRow.TargetKuantitatifTahunan, subRow.Glossary,
				itemQualifier, deskripsiQualifier, targetQualifier,
				"", subRow.TerdapatQualifier,
			)
		}
	}

	// =========================================================================
	// Build INSERT data_result_detail (hanya TW2/TW4)
	// =========================================================================
	resultPlaceholders := []string{}
	resultArgs := []interface{}{}
	for _, rs := range resultList {
		resultPlaceholders = append(resultPlaceholders, "(?, ?, ?, ?, ?, ?)")
		resultArgs = append(resultArgs,
			req.IdPengajuan,
			rs.IdDetailResult,
			req.Tahun,
			req.Triwulan,
			rs.NamaResult,
			rs.DeskripsiResult,
		)
	}

	// =========================================================================
	// Build INSERT data_method_detail (hanya TW2/TW4)
	// =========================================================================
	processPlaceholders := []string{}
	processArgs := []interface{}{}
	for _, mt := range processList {
		processPlaceholders = append(processPlaceholders, "(?, ?, ?, ?, ?, ?)")
		processArgs = append(processArgs,
			req.IdPengajuan,
			mt.IdDetailProcess,
			req.Tahun,
			req.Triwulan,
			mt.NamaProcess,
			mt.DeskripsiProcess,
		)
	}

	// =========================================================================
	// Build INSERT data_challenge_detail (hanya TW2/TW4)
	// =========================================================================
	contextPlaceholders := []string{}
	contextArgs := []interface{}{}
	for _, ch := range contextList {
		contextPlaceholders = append(contextPlaceholders, "(?, ?, ?, ?, ?, ?)")
		contextArgs = append(contextArgs,
			req.IdPengajuan,
			ch.IdDetailContext,
			req.Tahun,
			req.Triwulan,
			ch.NamaContext,
			ch.DeskripsiContext,
		)
	}

	// =========================================================================
	// Eksekusi dalam satu transaksi:
	//   1. DELETE data child lama
	//   2. INSERT data child baru
	//   3. UPDATE header data_kpi
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi revision: %w", tx.Error)
	}

	// -------------------------------------------------------------------------
	// DELETE data child lama
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryDeleteKpiDetail, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal delete data_kpi_detail: %w", err)
	}
	if err := tx.Exec(queryDeleteKpiSubDetail, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal delete data_kpi_subdetail: %w", err)
	}
	if err := tx.Exec(queryDeleteResultDetail, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal delete data_result_detail: %w", err)
	}
	if err := tx.Exec(queryDeleteProcessDetail, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal delete data_method_detail: %w", err)
	}
	if err := tx.Exec(queryDeleteContextDetail, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal delete data_challenge_detail: %w", err)
	}

	// -------------------------------------------------------------------------
	// INSERT data_kpi_detail baru
	// -------------------------------------------------------------------------
	if len(kpiDetailPlaceholders) > 0 {
		queryDetail := fmt.Sprintf(queryInsertKpiDetail, strings.Join(kpiDetailPlaceholders, ", "))
		if err := tx.Exec(queryDetail, kpiDetailArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_kpi_detail saat revision: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// INSERT data_kpi_subdetail baru
	// -------------------------------------------------------------------------
	if len(subDetailPlaceholders) > 0 {
		querySubDetail := fmt.Sprintf(queryInsertKpiSubDetail, strings.Join(subDetailPlaceholders, ", "))
		if err := tx.Exec(querySubDetail, subDetailArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_kpi_subdetail saat revision: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// INSERT data_result_detail baru (hanya TW2/TW4)
	// -------------------------------------------------------------------------
	if len(resultPlaceholders) > 0 {
		queryResult := fmt.Sprintf(queryInsertResultDetail, strings.Join(resultPlaceholders, ", "))
		if err := tx.Exec(queryResult, resultArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_result_detail saat revision: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// INSERT data_method_detail baru (hanya TW2/TW4)
	// -------------------------------------------------------------------------
	if len(processPlaceholders) > 0 {
		queryProcess := fmt.Sprintf(queryInsertProcessDetail, strings.Join(processPlaceholders, ", "))
		if err := tx.Exec(queryProcess, processArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_method_detail saat revision: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// INSERT data_challenge_detail baru (hanya TW2/TW4)
	// -------------------------------------------------------------------------
	if len(contextPlaceholders) > 0 {
		queryContext := fmt.Sprintf(queryInsertContextDetail, strings.Join(contextPlaceholders, ", "))
		if err := tx.Exec(queryContext, contextArgs...).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal insert data_challenge_detail saat revision: %w", err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE header data_kpi
	// status = 0 → langsung ke approval
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryUpdateKpiRevision,
		req.EntryTime, updatedApprovalList, firstApprovalPosisi, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update data_kpi saat revision: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi revision: %w", err)
	}

	return nil
}

// =============================================================================
// APPROVE PENYUSUNAN KPI
// =============================================================================

// ApprovePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/approve.
// Menerima approval_list (JSON string sudah diupdate) dan approval_posisi (next approver, kosong jika final).
func (r *penyusunanKpiRepo) ApprovePenyusunanKpi(idPengajuan, approvalList, approvalPosisi, user string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if approvalPosisi == "" {
		// ── Approve final: set status=2 ──
		if err := tx.Exec(queryApproveFinalPenyusunan,
			approvalList, idPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve final penyusunan: %w", err)
		}
		tx.Commit()
		return nil
	}

	// ── Approve chain: pindah approval_posisi ke level berikutnya + notif ──
	if err := tx.Exec(queryApproveChainPenyusunan,
		approvalPosisi, approvalList, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal approve chain penyusunan: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Approval Penyusunan, ID : "+idPengajuan,
		user,
		approvalPosisi,
		"approval_penyusunan",
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// =============================================================================
// REJECT PENYUSUNAN KPI
// =============================================================================

// RejectPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/reject.
// Menerima approval_list (JSON string sudah diupdate) dan catatan penolakan.
func (r *penyusunanKpiRepo) RejectPenyusunanKpi(idPengajuan, approvalList, catatan, user string) error {
	// Ambil entry_user untuk dikirim notifikasi penolakan
	var kpiBase struct {
		EntryUser string `gorm:"column:entry_user"`
	}
	if err := r.db.Raw(queryGetKpiBaseData, idPengajuan).Scan(&kpiBase).Error; err != nil {
		return fmt.Errorf("gagal mengambil entry_user: %w", err)
	}
	entryUser := kpiBase.EntryUser

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryRejectPenyusunan,
		approvalList, catatan, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal reject penyusunan: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Penyusunan Ditolak, ID : "+idPengajuan,
		user,
		entryUser,
		"penyusunan_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *penyusunanKpiRepo) GetApprovalListJSON(idPengajuan, userID string) (string, error) {
	var approvalListBytes []byte
	row := r.db.Raw(queryGetApprovalListJSON, userID, idPengajuan).Row()
	if err := row.Scan(&approvalListBytes); err != nil {
		return "", &customErrors.BadRequestError{Message: "Data Not Found"}
	}
	approvalList := string(approvalListBytes)
	if approvalList == "" {
		return "", &customErrors.BadRequestError{Message: "Data Not Found"}
	}
	return approvalList, nil
}

func (r *penyusunanKpiRepo) GetCatatanTolakan(idPengajuan string) (string, error) {
	var val []byte
	row := r.db.Raw(queryGetCatatanTolakan, idPengajuan).Row()
	if err := row.Scan(&val); err != nil {
		return "", err
	}
	return string(val), nil
}

// =============================================================================
// GET ALL APPROVAL — list dengan filter, pagination, dan nested detail
// =============================================================================

func (r *penyusunanKpiRepo) GetAllApprovalPenyusunanKpi(
	req *dto.GetAllApprovalPenyusunanKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status = 0",
		"a.approval_posisi = ?",
	}
	args := []interface{}{req.ApprovalUser}

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

	where := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// COUNT TOTAL RECORDS
	// =========================================================================
	var total int64
	countQuery := queryGetCountDataKpi + where
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
	listQuery := queryGetDataKpi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	rows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	defer rows.Close()

	var results []*model.DataKpi

	for rows.Next() {
		var h model.DataKpi

		if err := rows.Scan(
			&h.IdPengajuan, &h.Tahun, &h.Triwulan,
			&h.Kostl, &h.KostlTx,
			&h.Orgeh, &h.OrgehTx,
			&h.EntryUser, &h.EntryName, &h.EntryTime,
			&h.ApprovalPosisi, &h.ApprovalList,
			&h.Status, &h.StatusDesc,
			&h.EntryUserRealisasi, &h.EntryNameRealisasi, &h.EntryTimeRealisasi,
			&h.ApprovalListRealisasi, &h.CatatanTolakan,
			&h.TotalBobot, &h.TotalPencapaian,
			&h.TotalBobotPengurang, &h.TotalPencapaianPost,
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi, &h.LampiranValidasi, &h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET ALL TOLAKAN — list dengan filter, pagination, dan nested detail
// =============================================================================

func (r *penyusunanKpiRepo) GetAllTolakanPenyusunanKpi(
	req *dto.GetAllTolakanPenyusunanKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status = 1",
	}
	args := []interface{}{}

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

	where := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// COUNT TOTAL RECORDS
	// =========================================================================
	var total int64
	countQuery := queryGetCountDataKpi + where
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
	listQuery := queryGetDataKpi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	rows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	defer rows.Close()

	var results []*model.DataKpi

	for rows.Next() {
		var h model.DataKpi

		if err := rows.Scan(
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

		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET ALL DAFTAR PENYUSUNAN — list dengan filter, pagination, dan nested detail
// =============================================================================

func (r *penyusunanKpiRepo) GetAllDaftarPenyusunanKpi(
	req *dto.GetAllDaftarPenyusunanKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{"1=1"}
	args := []interface{}{}

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

	where := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// COUNT TOTAL RECORDS
	// =========================================================================
	var total int64
	countQuery := queryGetCountDataKpi + where
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
	listQuery := queryGetDataKpi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	rows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	defer rows.Close()

	var results []*model.DataKpi

	for rows.Next() {
		var h model.DataKpi

		if err := rows.Scan(
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

		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET ALL DAFTAR APPROVAL — list dengan filter, pagination, dan nested detail
// =============================================================================

func (r *penyusunanKpiRepo) GetAllDaftarApprovalPenyusunanKpi(
	req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.approval_list LIKE ?",
	}
	args := []interface{}{"%" + req.ApprovalUser + "%"}

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

	where := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// COUNT TOTAL RECORDS
	// =========================================================================
	var total int64
	countQuery := queryGetCountDataKpi + where
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
	listQuery := queryGetDataKpi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	rows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	defer rows.Close()

	var results []*model.DataKpi

	for rows.Next() {
		var h model.DataKpi

		if err := rows.Scan(
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

		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET DETAIL
// =============================================================================

func (r *penyusunanKpiRepo) GetDetailPenyusunanKpi(
	req *dto.GetDetailPenyusunanKpiRequest,
) (*model.DataKpi, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.id_pengajuan = ?",
	}
	args := []interface{}{req.IdPengajuan}

	where := " WHERE " + strings.Join(conditions, " AND ")

	// =========================================================================
	// QUERY HEADER
	// =========================================================================
	var result model.DataKpi
	headerQuery := queryGetDataKpi + where + " LIMIT 1"
	if err := r.db.Raw(headerQuery, args...).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil detail KPI: %w", err)
	}

	// =========================================================================
	// KPI DETAIL + SUB DETAIL
	// =========================================================================
	var kpiDetails []model.DataKpiDetail
	if err := r.db.Raw(queryGetDataKpiDetail, result.IdPengajuan).Scan(&kpiDetails).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil kpi detail: %w", err)
	}
	for i := range kpiDetails {
		var subDetails []model.DataKpiSubDetail
		if err := r.db.Raw(queryGetDataKpiSubDetail, kpiDetails[i].IdDetail).Scan(&subDetails).Error; err != nil {
			return nil, fmt.Errorf("gagal mengambil kpi sub detail: %w", err)
		}
		if subDetails == nil {
			subDetails = []model.DataKpiSubDetail{}
		}
		kpiDetails[i].KpiSubDetail = subDetails
		kpiDetails[i].TotalSubKpi = len(subDetails)
	}
	if kpiDetails == nil {
		kpiDetails = []model.DataKpiDetail{}
	}
	result.Kpi = kpiDetails
	result.TotalKpi = len(kpiDetails)

	// =========================================================================
	// RESULT DETAIL
	// =========================================================================
	var resultList []model.DataResultDetail
	if err := r.db.Raw(queryGetDataResultDetail, result.IdPengajuan).Scan(&resultList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil result detail: %w", err)
	}
	if resultList == nil {
		resultList = []model.DataResultDetail{}
	}
	result.ResultList = resultList
	result.TotalResult = len(resultList)

	// =========================================================================
	// PROCESS DETAIL
	// =========================================================================
	var processList []model.DataMethodDetail
	if err := r.db.Raw(queryGetDataProcessDetail, result.IdPengajuan).Scan(&processList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil process detail: %w", err)
	}
	if processList == nil {
		processList = []model.DataMethodDetail{}
	}
	result.ProcessList = processList
	result.TotalProcess = len(processList)

	// =========================================================================
	// CONTEXT DETAIL
	// =========================================================================
	var contextList []model.DataChallengeDetail
	if err := r.db.Raw(queryGetDataContextDetail, result.IdPengajuan).Scan(&contextList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil context detail: %w", err)
	}
	if contextList == nil {
		contextList = []model.DataChallengeDetail{}
	}
	result.ContextList = contextList
	result.TotalContext = len(contextList)

	return &result, nil
}

// =============================================================================
// GET EXPORT DATA — digunakan bersama oleh get-csv dan get-pdf
// =============================================================================

func (r *penyusunanKpiRepo) GetKpiExportData(
	idPengajuan, kostl, tahun, triwulan string,
) (*dto.KpiExportData, error) {

	type kpiHeader struct {
		KostlTx  string `gorm:"column:kostl_tx"`
		Tahun    string `gorm:"column:tahun"`
		Triwulan string `gorm:"column:triwulan"`
	}
	var header kpiHeader
	if err := r.db.Raw(queryGetKpiBaseDataForExport, idPengajuan, kostl, tahun, triwulan).Scan(&header).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil header KPI: %w", err)
	}

	type subDetailRaw struct {
		Kpi           string `gorm:"column:kpi"`
		Bobot         string `gorm:"column:bobot"`
		TargetTahunan string `gorm:"column:target_tahunan"`
		Capping       string `gorm:"column:capping"`
	}
	var rawRows []subDetailRaw
	if err := r.db.Raw(queryGetSubDetailForExport, idPengajuan, kostl, tahun, triwulan).Scan(&rawRows).Error; err != nil {
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

// =============================================================================
// GetKpiHeader
// =============================================================================

func (r *penyusunanKpiRepo) GetKpiHeader(idPengajuan string) (tahun, triwulan, kostl, kostlTx, entryUser, entryName string, status int, statusDesc string, err error) {
	type kpiHeader struct {
		Tahun      string `gorm:"column:tahun"`
		Triwulan   string `gorm:"column:triwulan"`
		Kostl      string `gorm:"column:kostl"`
		KostlTx    string `gorm:"column:kostl_tx"`
		Status     int    `gorm:"column:status"`
		StatusDesc string `gorm:"column:status_desc"`
		EntryUser  string `gorm:"column:entry_user"`
		EntryName  string `gorm:"column:entry_name"`
	}
	var h kpiHeader
	if err = r.db.Raw(queryGetKpiHeader, idPengajuan).Scan(&h).Error; err != nil {
		return "", "", "", "", "", "", 0, "", fmt.Errorf("gagal mengambil header KPI: %w", err)
	}
	if h.Tahun == "" {
		return "", "", "", "", "", "", 0, "", fmt.Errorf("id_pengajuan '%s' tidak ditemukan", idPengajuan)
	}
	return h.Tahun, h.Triwulan, h.Kostl, h.KostlTx, h.EntryUser, h.EntryName, h.Status, h.StatusDesc, nil
}
