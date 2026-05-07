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
	// Get Count
	// =============================================================================

	// Use func : GetAllApprovalPenyusunanKpi, GetAllTolakanPenyusunanKpi, GetAllDaftarPenyusunanKpi, GetAllDaftarApprovalPenyusunanKpi
	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : CheckApprovalPenyusunanExists
	queryGetCountApprovalKpi = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 0 AND approval_posisi = ? AND id_pengajuan = ?`

	// =============================================================================
	// Get Data
	// =============================================================================

	// Use func : CreatePenyusunanKpi, RejectPenyusunanKpi
	queryGetKpiBaseData = `
		SELECT kostl_tx, entry_user
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// Use func : GetCatatanTolakan
	queryGetCatatanTolakan = `
		SELECT IFNULL(catatan_tolakan, '') FROM data_kpi
		WHERE id_pengajuan = ? LIMIT 1`

	// =============================================================================
	// Get Exist
	// =============================================================================

	// Use func : GetExistDataKpi
	queryGetExistDataKpi = `
		SELECT a.tahun, a.triwulan, a.kostl, a.kostl_tx, a.status, b.status_desc, a.entry_user, a.entry_name
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ?
		LIMIT 1`

	// Use func : GetExistDataKpiStatus
	queryGetExistDataKpiStatus = `
		SELECT id_pengajuan, status
		FROM data_kpi
		WHERE tahun = ? AND triwulan = ? AND kostl = ?
		LIMIT 1`

	// =============================================================================
	// Get Helper
	// =============================================================================

	// Use func : GetApprovalListJSON
	queryGetApprovalListJSON = `
		SELECT approval_list FROM data_kpi
		WHERE status = 0 AND approval_posisi = ? AND id_pengajuan = ?`

	// Use func : RevisionPenyusunanKpi
	queryGetApprovalForRevision = `
		SELECT approval_posisi, approval_list FROM data_kpi
		WHERE id_pengajuan = ? LIMIT 1`

	// =============================================================================
	// SERVICE HELPERS
	// =============================================================================

	// Use func : LookupKpiMaster
	queryLookupSubKpi = `
		SELECT id_kpi, kpi, rumus
		FROM mst_kpi
		WHERE LOWER(kpi) = LOWER(?)
		LIMIT 1`

	// Use func : LookupPolarisasi
	queryLookupPolarisasi = `
		SELECT id_polarisasi
		FROM mst_polarisasi
		WHERE LOWER(polarisasi) = LOWER(?)
		LIMIT 1`

	// =============================================================================
	// Get Export
	// =============================================================================

	// Use func : GetKpiExportData
	queryGetKpiBaseDataForExport = `
		SELECT kostl_tx, tahun, triwulan, entry_user
		FROM data_kpi
		WHERE id_pengajuan = ? AND kostl = ? AND tahun = ? AND triwulan = ?
		LIMIT 1`

	// Use func : GetKpiExportData
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

	// =============================================================================
	// Get Other
	// =============================================================================

	// Use func : ValidatePenyusunanKpi
	queryGetOrgeh = `
		SELECT orgeh, orgeh_tx 
		FROM user 
		WHERE kostl = ? 
		ORDER BY HILFM ASC 
		LIMIT 1`

	// =============================================================================
	// Get KPI
	// =============================================================================

	// Use func : GetAllApprovalPenyusunanKpi, GetAllTolakanPenyusunanKpi, GetAllDaftarPenyusunanKpi, GetAllDaftarApprovalPenyusunanKpi, dan GetDetailPenyusunanKpi
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

	// Use func : GetDetailPenyusunanKpi
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

	// Use func : GetDetailPenyusunanKpi
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

	// Use func : GetDetailPenyusunanKpi
	queryGetDataResultDetail = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result
		FROM data_result_detail
		WHERE id_pengajuan = ?`

	// Use func : GetDetailPenyusunanKpi
	queryGetDataProcessDetail = `
		SELECT
			id_detail_method,
			nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')   realisasi_method,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?`

	// Use func : GetDetailPenyusunanKpi
	queryGetDataContextDetail = `
		SELECT
			id_detail_challenge,
			nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '')  realisasi_challenge,
			IFNULL(lampiran_evidence, '')    lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Insert
	// =============================================================================

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertKpi = `
		INSERT INTO data_kpi 
			(id_pengajuan, tahun, triwulan, kostl, kostl_tx, orgeh, orgeh_tx, 
			 entry_user, entry_name, entry_time, approval_posisi, approval_list, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertKpiDetail = `
		INSERT INTO data_kpi_detail 
			(id_pengajuan, id_detail, tahun, triwulan, id_kpi, kpi, rumus, 
			 id_perspektif, id_keterangan_project) 
		VALUES %s`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertKpiSubDetail = `
		INSERT INTO data_kpi_subdetail 
			(id_pengajuan, id_detail, id_sub_detail, tahun, triwulan, 
			id_kpi, kpi, rumus, otomatis, bobot, capping, 
			target_triwulan, target_kuantitatif_triwulan, 
			target_tahunan, target_kuantitatif_tahunan, 
			deskripsi_glossary, item_qualifier, deskripsi_qualifier, 
			target_qualifier, id_keterangan_project, id_qualifier) 
		VALUES %s`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertResultDetail = `
		INSERT INTO data_result_detail
			(id_pengajuan, id_detail_result, tahun, triwulan,
			nama_result, deskripsi_result)
		VALUES %s`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertProcessDetail = `
		INSERT INTO data_method_detail 
			(id_pengajuan, id_detail_method, tahun, triwulan, 
			 nama_method, deskripsi_method) 
		VALUES %s`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryInsertContextDetail = `
		INSERT INTO data_challenge_detail 
			(id_pengajuan, id_detail_challenge, tahun, triwulan, 
			 nama_challenge, deskripsi_challenge) 
		VALUES %s`

	// =============================================================================
	// Update KPI
	// =============================================================================

	// Use func : CreatePenyusunanKpi
	queryUpdateKpi = `
		UPDATE data_kpi 
		SET approval_posisi = ?, approval_list = ?, status = 0
		WHERE id_pengajuan = ?`

	// Use func : RevisionPenyusunanKpi
	queryUpdateKpiRevision = `
		UPDATE data_kpi
		SET entry_time = ?, approval_list = ?, approval_posisi = ?, status = 0
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Update Approval
	// =============================================================================

	// Use func : ApprovePenyusunanKpi
	queryApproveChainPenyusunan = `
		UPDATE data_kpi SET approval_posisi = ?, approval_list = ? WHERE id_pengajuan = ?`

	// Use func : ApprovePenyusunanKpi
	queryApproveFinalPenyusunan = `
		UPDATE data_kpi SET status = 2, approval_list = ? WHERE id_pengajuan = ?`

	// Use func : RejectPenyusunanKpi
	queryRejectPenyusunan = `
		UPDATE data_kpi SET status = 1, approval_list = ?, catatan_tolakan = ? WHERE id_pengajuan = ?`

	// =============================================================================
	// Delete KPI
	// =============================================================================

	// Use func : ValidatePenyusunanKpi
	queryDeleteKpi = `DELETE FROM data_kpi WHERE id_pengajuan = ?`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryDeleteKpiDetail = `DELETE FROM data_kpi_detail WHERE id_pengajuan = ?`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryDeleteKpiSubDetail = `DELETE FROM data_kpi_subdetail WHERE id_pengajuan = ?`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryDeleteResultDetail = `DELETE FROM data_result_detail WHERE id_pengajuan = ?`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryDeleteProcessDetail = `DELETE FROM data_method_detail WHERE id_pengajuan = ?`

	// Use func : ValidatePenyusunanKpi, RevisionPenyusunanKpi
	queryDeleteContextDetail = `DELETE FROM data_challenge_detail WHERE id_pengajuan = ?`
)

// =============================================================================
// VALIDATE
// =============================================================================

// ValidatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/validate.
func (r *penyusunanKpiRepo) ValidatePenyusunanKpi(
	req *dto.ValidatePenyusunanKpiRequest,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
	idLama string,
) (string, error) {

	idPengajuan := idgen.GenerateIDPengajuan(req.Divisi.Kostl, req.Tahun, req.Triwulan)

	var orgeh, orgehTx string
	r.db.Raw(queryGetOrgeh, req.Divisi.Kostl).Row().Scan(&orgeh, &orgehTx)

	// status 70 = draft
	var statusKpi interface{} = 70

	// =========================================================================
	// Build INSERT data_kpi_detail
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
			{queryDeleteKpi, "data_kpi"},
		} {
			if err := tx.Exec(q.query, idLama).Error; err != nil {
				tx.Rollback()
				return "", fmt.Errorf("gagal menghapus draft lama (%s): %w", q.desc, err)
			}
		}
	}

	// approval_posisi dan approval_list dikosongkan — diisi saat CreatePenyusunanKpi
	if err := tx.Exec(queryInsertKpi,
		idPengajuan, req.Tahun, req.Triwulan, req.Divisi.Kostl, req.Divisi.KostlTx,
		orgeh, orgehTx, req.EntryUserPenyusunan, req.EntryNamePenyusunan, req.EntryTimePenyusunan,
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
// CREATE
// =============================================================================

// CreatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/create.
func (r *penyusunanKpiRepo) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
) error {
	// Ambil userid pertama dari ApprovalList sebagai approval_posisi
	approvalPosisi := ""
	if len(req.ApprovalListPenyusunan) > 0 {
		approvalPosisi = req.ApprovalListPenyusunan[0].Userid
	}

	approvalListBytes, err := json.Marshal(req.ApprovalListPenyusunan)
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
		kpiBase.KostlTx,
		req.EntryUserPenyusunan,
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
// REVISION
// =============================================================================

// RevisionPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/revision.
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
		req.EntryTimePenyusunan, updatedApprovalList, firstApprovalPosisi, req.IdPengajuan,
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
// APPROVAL
// =============================================================================

// ApprovePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/approve.
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

// RejectPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/reject.
func (r *penyusunanKpiRepo) RejectPenyusunanKpi(idPengajuan, approvalList, catatan, user string) error {
	// Ambil entry_user untuk dikirim notifikasi penolakan
	var kpiBase struct {
		EntryUser string `gorm:"column:entry_user"`
	}
	if err := r.db.Raw(queryGetKpiBaseData, idPengajuan).Scan(&kpiBase).Error; err != nil {
		return fmt.Errorf("gagal mengambil entry_user: %w", err)
	}

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
		kpiBase.EntryUser,
		"penyusunan_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-approval.
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
	args := []interface{}{req.ApprovalUserPenyusunan}

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

// GetAllTolakanPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-tolakan.
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

// GetAllDaftarPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
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

// GetAllDaftarApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-approval.
func (r *penyusunanKpiRepo) GetAllDaftarApprovalPenyusunanKpi(
	req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status IN (3, 5)",
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

// GetDetailPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
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
// EXPORT DATA
// =============================================================================

// GetKpiExportData digunakan oleh endpoint POST /penyusunan-kpi/get-excel dan /penyusunan-kpi/get-pdf.
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
// APPROVAL HELPER
// =============================================================================

// GetApprovalListJSON digunakan oleh service ApprovePenyusunanKpi dan RejectPenyusunanKpi untuk mengambil daftar approval dalam format JSON.
func (r *penyusunanKpiRepo) GetApprovalListJSON(idPengajuan, userID string) (string, error) {
	var approvalListBytes []byte
	row := r.db.Raw(queryGetApprovalListJSON, userID, idPengajuan).Row()
	if err := row.Scan(&approvalListBytes); err != nil {
		return "", &customErrors.BadRequestError{Message: "Data List Approval tidak ditemukan."}
	}
	approvalList := string(approvalListBytes)
	if approvalList == "" {
		return "", &customErrors.BadRequestError{Message: "Data KPI untuk Pengajuan ini tidak ditemukan."}
	}
	return approvalList, nil
}

// GetCatatanTolakan digunakan oleh service RejectPenyusunanKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
func (r *penyusunanKpiRepo) GetCatatanTolakan(idPengajuan string) (string, error) {
	var val []byte
	row := r.db.Raw(queryGetCatatanTolakan, idPengajuan).Row()
	if err := row.Scan(&val); err != nil {
		return "", err
	}
	return string(val), nil
}

// CheckApprovalPenyusunanExists digunakan oleh service ApprovePenyusunanKpi dan RejectPenyusunanKpi untuk memvalidasi keberadaan approval.
func (r *penyusunanKpiRepo) CheckApprovalPenyusunanExists(user, idPengajuan string) (bool, error) {
	var count int64
	if err := r.db.Raw(queryGetCountApprovalKpi, user, idPengajuan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	return count > 0, nil
}

// =============================================================================
// GET EXIST
// =============================================================================

// GetExistDataKpi digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
func (r *penyusunanKpiRepo) GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error) {
	var result model.DataKpiExist
	db := r.db.Raw(queryGetExistDataKpi, idPengajuan).Scan(&result)
	if db.Error != nil {
		return nil, fmt.Errorf("gagal mengambil header KPI: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return nil, fmt.Errorf("id_pengajuan '%s' tidak ditemukan", idPengajuan)
	}
	return &result, nil
}

// GetExistDataKpiStatus digunakan oleh service untuk mengecek keberadaan data KPI dan mengembalikan id_pengajuan beserta statusnya.
func (r *penyusunanKpiRepo) GetExistDataKpiStatus(tahun, triwulan, kostl string) (idPengajuan string, status int, found bool, err error) {
	row := r.db.Raw(queryGetExistDataKpiStatus, tahun, triwulan, kostl).Row()
	if scanErr := row.Scan(&idPengajuan, &status); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return "", 0, false, nil
		}
		return "", 0, false, fmt.Errorf("gagal mengecek data Penyusunan KPI: %w", scanErr)
	}
	return idPengajuan, status, true, nil
}

// =============================================================================
// SERVICE HELPERS
// =============================================================================

// LookupKpiMaster digunakan oleh service ValidatePenyusunanKpi dan RevisionPenyusunanKpi untuk mencari id_kpi, kpi, dan rumus dari mst_kpi.
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

// LookupPolarisasi digunakan oleh service ValidatePenyusunanKpi dan RevisionPenyusunanKpi untuk mencari id_polarisasi dari mst_polarisasi.
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
