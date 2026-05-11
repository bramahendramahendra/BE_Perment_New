package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dto "permen_api/domain/realisasi_kpi/dto"
	model "permen_api/domain/realisasi_kpi/model"
	customErrors "permen_api/errors"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Get Count
	// =============================================================================

	// Use func : GetAllRealisasiKpi, GetAllApprovalRealisasiKpi, GetAllTolakanRealisasiKpi, GetAllDaftarRealisasiKpi, GetAllDaftarApprovalRealisasiKpi
	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : CheckApprovalRealisasiExists
	queryGetCountApprovalKpi = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 3 AND approval_posisi = ? AND id_pengajuan = ?`

	// =============================================================================
	// Get Data
	// =============================================================================

	// Use func : CreateRealisasiKpi, RejectRealisasiKpi
	queryGetKpiBaseData = `
		SELECT kostl_tx, entry_user_realisasi
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
		SELECT a.tahun, a.triwulan, a.kostl, a.kostl_tx, a.status, b.status_desc, a.entry_user_realisasi, a.entry_name_realisasi
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ? AND status IN (2, 3, 4, 80, 81)
		LIMIT 1`

	// =============================================================================
	// Get Helper
	// =============================================================================

	// Use func : GetApprovalListJSON
	queryGetApprovalListJSON = `
		SELECT approval_list_realisasi FROM data_kpi
		WHERE status = 3 AND approval_posisi = ? AND id_pengajuan = ?`

	// Use func : RevisionRealisasiKpi
	queryGetApprovalForRevision = `
		SELECT approval_posisi, approval_list_realisasi FROM data_kpi
		WHERE id_pengajuan = ? LIMIT 1`

	// =============================================================================
	// SERVICE HELPERS
	// =============================================================================

	// Use func : GetLinkFormats
	queryGetLinkFormats = `
		SELECT url_prefix
		FROM mst_link_format
		WHERE is_active = 1
		ORDER BY id_link_format ASC`

	queryLookupSubDetail = `
		SELECT
			s.id_sub_detail,
			s.id_detail,
			d.id_kpi,
			d.rumus    AS detail_rumus,
			s.id_kpi   AS id_sub_kpi,
			s.rumus,
			s.otomatis,
			IFNULL(s.deskripsi_glossary, '') glossary,
			s.target_triwulan,
			IFNULL(s.target_kuantitatif_triwulan, 0) target_kuantitatif_triwulan,
			s.target_tahunan,
			IFNULL(s.target_kuantitatif_tahunan, 0) target_kuantitatif_tahunan,
			IFNULL(s.id_qualifier, '') AS id_qualifier,
			IFNULL(s.item_qualifier, '')                  qualifier,
			IFNULL(s.deskripsi_qualifier, '')             deskripsi_qualifier
		FROM data_kpi_subdetail s
		INNER JOIN data_kpi_detail d ON d.id_detail = s.id_detail
		WHERE s.id_pengajuan = ?
		  AND LOWER(d.kpi)  = LOWER(?)
		  AND LOWER(s.kpi)  = LOWER(?)
		LIMIT 1`

	// =============================================================================
	// Get KPI
	// =============================================================================

	// Use func : GetAllRealisasiKpi, GetAllApprovalRealisasiKpi, GetAllTolakanRealisasiKpi, GetAllDaftarRealisasiKpi, GetAllDaftarApprovalRealisasiKpi, dan GetDetailRealisasiKpi
	queryGetDataKpi = `
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
			IFNULL(a.total_pencapaian, '')         total_pencapaian
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : GetDetailRealisasiKpi
	queryGetDataKpiDetail = `
		SELECT
			a.id_detail,
			a.id_kpi, a.kpi, a.rumus,
			IFNULL(a.id_perspektif, '')         id_perspektif,
			IFNULL(b.perspektif, '')            perspektif,
			IFNULL(a.id_keterangan_project, '') id_keterangan_project,
			IFNULL(c.keterangan_project, '')    keterangan_project,
			IFNULL(a.lampiran_file, '') 		lampiran_file
		FROM data_kpi_detail a
		LEFT JOIN mst_perspektif b ON a.id_perspektif = b.id_perspektif
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ?
		ORDER BY a.id_detail ASC`

	// Use func : GetDetailRealisasiKpi
	queryGetDataKpiSubDetail = `
		SELECT
			a.id_sub_detail,
			a.id_kpi, a.kpi, a.rumus,
			a.otomatis,
			a.bobot, a.capping,
			a.target_triwulan,  a.target_kuantitatif_triwulan,
			a.target_tahunan,   a.target_kuantitatif_tahunan,
			IFNULL(a.deskripsi_glossary, '')              deskripsi_glossary,
			IFNULL(a.rumus, '')                           id_polarisasi,
			IFNULL(p.polarisasi, '')                      polarisasi,
			IFNULL(a.id_qualifier, '')                    id_qualifier,
			IFNULL(a.item_qualifier, '')                  item_qualifier,
			IFNULL(a.deskripsi_qualifier, '')             deskripsi_qualifier,
			IFNULL(a.target_qualifier, '')                target_qualifier,
			IFNULL(a.id_keterangan_project, '')           id_keterangan_project,
			IFNULL(c.keterangan_project, '')              keterangan_project,
			IFNULL(a.realisasi, '')                            realisasi,
			COALESCE(NULLIF(a.realisasi_kuantitatif, ''), 0)  realisasi_kuantitatif,
			IFNULL(a.realisasi_keterangan, '')                realisasi_keterangan,
			IFNULL(a.realisasi_validated, '')                 realisasi_validated,
			IFNULL(a.realisasi_kuantitatif_validated, '')     realisasi_kuantitatif_validated,
			COALESCE(NULLIF(a.pencapaian, ''), 0)             pencapaian,
			COALESCE(NULLIF(a.skor, ''), 0)                   skor,
			IFNULL(a.realisasi_qualifier, '')             realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '') realisasi_kuantitatif_qualifier
		FROM data_kpi_subdetail a
		LEFT JOIN mst_polarisasi p ON a.rumus = p.id_polarisasi
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	// Use func : GetDetailRealisasiKpi
	queryGetDataResultDetail = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result,
			IFNULL(realisasi_result, '')   realisasi_result,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_result_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_result ASC`

	// Use func : GetDetailRealisasiKpi
	queryGetDataProcessDetail = `
		SELECT
			id_detail_method, nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')  realisasi_method,
			IFNULL(lampiran_evidence, '') lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_method ASC`

	// Use func : GetDetailRealisasiKpi
	queryGetDataContextDetail = `
		SELECT
			id_detail_challenge, nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '') realisasi_challenge,
			IFNULL(lampiran_evidence, '')   lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_challenge ASC`

	// =============================================================================
	// Update KPI
	// =============================================================================

	// Use func : ValidateRealisasiKpi
	queryUpdateKpiRealisasi = `
		UPDATE data_kpi
		SET status                  = ?,
		    entry_user_realisasi    = ?,
		    entry_name_realisasi    = ?,
		    entry_time_realisasi    = ?
		WHERE id_pengajuan = ?`

	// Use func : CreateRealisasiKpi
	queryUpdateKpiRealisasiCreate = `
		UPDATE data_kpi 
		SET approval_posisi = ?, approval_list_realisasi = ?, status = 3
		WHERE id_pengajuan = ?`

	// Use func : RevisionRealisasiKpi
	queryUpdateKpiRealisasiRevision = `
		UPDATE data_kpi
		SET entry_time_realisasi = ?, approval_posisi = ?, approval_list_realisasi = ?, status = 3
		WHERE id_pengajuan = ?`

	// Use func : ValidateRealisasiKpi, RevisionRealisasiKpi
	queryUpdateKpiSubDetailRealisasi = `
		UPDATE data_kpi_subdetail
		SET realisasi                      	= ?,
		    realisasi_kuantitatif          	= ?,
		    realisasi_validated            	= ?,
		    realisasi_kuantitatif_validated = ?,
		    pencapaian                     	= ?,
		    skor                           	= ?,
		    realisasi_qualifier            	= ?,
		    realisasi_kuantitatif_qualifier = ?
		WHERE id_pengajuan = ?
		  AND id_sub_detail = ?`

	// Use func : ValidateRealisasiKpi, RevisionRealisasiKpi
	queryUpdateKpiDetailLampiranFile = `
		UPDATE data_kpi_detail
		SET lampiran_file = ?
		WHERE id_pengajuan = ?
		  AND id_detail = ?`

	// Use func : ValidateRealisasiKpi, RevisionRealisasiKpi
	queryUpdateResultDetailRealisasi = `
		UPDATE data_result_detail
		SET realisasi_result   = ?,
		    lampiran_evidence  = ?
		WHERE id_pengajuan = ?
		  AND id_detail_result = ?`

	// Use func : ValidateRealisasiKpi, RevisionRealisasiKpi
	queryUpdateProcessDetailRealisasi = `
		UPDATE data_method_detail
		SET realisasi_method  = ?,
		    lampiran_evidence = ?
		WHERE id_pengajuan = ?
		  AND id_detail_method = ?`

	// Use func : ValidateRealisasiKpi, RevisionRealisasiKpi
	queryUpdateContextDetailRealisasi = `
		UPDATE data_challenge_detail
		SET realisasi_challenge = ?,
		    lampiran_evidence   = ?
		WHERE id_pengajuan = ?
		  AND id_detail_challenge = ?`

	// =============================================================================
	// Update Approval
	// =============================================================================

	// Use func : ApproveRealisasiKpi
	queryApproveChainRealisasi = `
		UPDATE data_kpi
		SET approval_posisi         = ?,
		    approval_list_realisasi = ?
		WHERE id_pengajuan = ?`

	// Use func : ApproveRealisasiKpi
	queryApproveFinalRealisasi = `
		UPDATE data_kpi
		SET status                  = 5,
		    approval_list_realisasi = ?
		WHERE id_pengajuan = ?`

	// Use func : RejectRealisasiKpi
	queryRejectRealisasi = `
		UPDATE data_kpi
		SET status                  = 4,
		    approval_list_realisasi = ?,
		    catatan_tolakan         = ?
		WHERE id_pengajuan = ?`
)

// =============================================================================
// VALIDATE
// =============================================================================

// ValidateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/validate.
func (r *realisasiKpiRepo) ValidateRealisasiKpi(
	req *dto.ValidateRealisasiKpiRequest,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
) error {

	// status 80 = draft realisasi
	var statusKpi any = 80

	// =========================================================================
	// Eksekusi semua UPDATE dalam satu transaksi
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// -------------------------------------------------------------------------
	// UPDATE setiap baris sub KPI
	// -------------------------------------------------------------------------
	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if err := tx.Exec(queryUpdateKpiSubDetailRealisasi,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Realisasi,            // realisasi_validated = sama dengan realisasi
				row.RealisasiKuantitatif, // realisasi_kuantitatif_validated = sama
				row.Pencapaian,
				row.Skor,
				row.RealisasiQualifier,
				row.RealisasiKuantitatifQualifier,
				req.IdPengajuan,
				row.IdSubDetail,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update sub detail '%s': %w", row.IdSubDetail, err)
			}
		}
		// UPDATE lampiran_file pada data_kpi_detail (1 Link Dokumen Sumber per KPI)
		// Ambil link dari sub-detail pertama yang tidak kosong dalam grup KPI ini
		linkDokumen := ""
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if row.LinkDokumenSumber != nil && *row.LinkDokumenSumber != "" {
				linkDokumen = *row.LinkDokumenSumber
				break
			}
		}
		if err := tx.Exec(queryUpdateKpiDetailLampiranFile,
			linkDokumen, req.IdPengajuan, kpiRow.IdDetail,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update lampiran_file detail '%s': %w", kpiRow.IdDetail, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE result (data_result_detail.realisasi_result) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, r2 := range resultList {
		if err := tx.Exec(queryUpdateResultDetailRealisasi,
			r2.RealisasiResult, r2.LampiranEvidence, req.IdPengajuan, r2.IdDetailResult,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update result '%s': %w", r2.IdDetailResult, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE process (data_method_detail.realisasi_method) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, p := range processList {
		if err := tx.Exec(queryUpdateProcessDetailRealisasi,
			p.RealisasiProcess, p.LampiranEvidence, req.IdPengajuan, p.IdDetailProcess,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update process '%s': %w", p.IdDetailProcess, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE context (data_challenge_detail.realisasi_challenge) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, c := range contextList {
		if err := tx.Exec(queryUpdateContextDetailRealisasi,
			c.RealisasiContext, c.LampiranEvidence, req.IdPengajuan, c.IdDetailContext,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update context '%s': %w", c.IdDetailContext, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE header data_kpi → status 80
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryUpdateKpiRealisasi,
		statusKpi, req.EntryUserRealisasi, req.EntryNameRealisasi, req.EntryTimeRealisasi, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update kpi realisasi: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi validate realisasi: %w", err)
	}

	return nil
}

// =============================================================================
// CREATE
// =============================================================================

// CreateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/create.
func (r *realisasiKpiRepo) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) error {
	// Ambil userid pertama dari ApprovalList sebagai approval_posisi
	approvalPosisi := ""
	if len(req.ApprovalListRealisasi) > 0 {
		approvalPosisi = req.ApprovalListRealisasi[0].Userid
	}

	approvalListBytes, err := json.Marshal(req.ApprovalListRealisasi)
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

	if err := tx.Exec(queryUpdateKpiRealisasiCreate,
		approvalPosisi, string(approvalListBytes), req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal submit realisasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		kpiBase.KostlTx,
		req.EntryUserRealisasi,
		approvalPosisi,
		"approval_realisasi",
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

// RevisionRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/revision.
func (r *realisasiKpiRepo) RevisionRealisasiKpi(
	req *dto.RevisionRealisasiKpiRequest,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
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

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// -------------------------------------------------------------------------
	// UPDATE setiap baris sub KPI
	// -------------------------------------------------------------------------
	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if err := tx.Exec(queryUpdateKpiSubDetailRealisasi,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Pencapaian,
				row.Skor,
				row.RealisasiQualifier,
				row.RealisasiKuantitatifQualifier,
				req.IdPengajuan,
				row.IdSubDetail,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update sub detail '%s': %w", row.IdSubDetail, err)
			}
		}
		// UPDATE lampiran_file pada data_kpi_detail (1 Link Dokumen Sumber per KPI)
		linkDokumen := ""
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if row.LinkDokumenSumber != nil && *row.LinkDokumenSumber != "" {
				linkDokumen = *row.LinkDokumenSumber
				break
			}
		}
		if err := tx.Exec(queryUpdateKpiDetailLampiranFile,
			linkDokumen, req.IdPengajuan, kpiRow.IdDetail,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update lampiran_file detail '%s': %w", kpiRow.IdDetail, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE result (data_result_detail.realisasi_result) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, r2 := range resultList {
		if err := tx.Exec(queryUpdateResultDetailRealisasi,
			r2.RealisasiResult, r2.LampiranEvidence, req.IdPengajuan, r2.IdDetailResult,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update result '%s': %w", r2.IdDetailResult, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE process (data_method_detail.realisasi_method) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, p := range processList {
		if err := tx.Exec(queryUpdateProcessDetailRealisasi,
			p.RealisasiProcess, p.LampiranEvidence, req.IdPengajuan, p.IdDetailProcess,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update process '%s': %w", p.IdDetailProcess, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE context (data_challenge_detail.realisasi_challenge) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, c := range contextList {
		if err := tx.Exec(queryUpdateContextDetailRealisasi,
			c.RealisasiContext, c.LampiranEvidence, req.IdPengajuan, c.IdDetailContext,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update context '%s': %w", c.IdDetailContext, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE header → status tetap 80 (draft), update entry time
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryUpdateKpiRealisasiRevision,
		req.EntryTimeRealisasi, firstApprovalPosisi, updatedApprovalList, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update header revisi realisasi: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi revision realisasi: %w", err)
	}

	return nil
}

// =============================================================================
// APPROVAL
// =============================================================================

// ApproveRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/approve.
func (r *realisasiKpiRepo) ApproveRealisasiKpi(idPengajuan, approvalList, approvalPosisi, user string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if approvalPosisi == "" {
		// Approve final: set status = 5
		if err := tx.Exec(queryApproveFinalRealisasi,
			approvalList, idPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve final realisasi: %w", err)
		}
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("gagal commit approve final realisasi: %w", err)
		}
		return nil
	}

	// Masih ada approver berikutnya → chain approval + notif
	if err := tx.Exec(queryApproveChainRealisasi,
		approvalPosisi, approvalList, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal approve chain realisasi: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Approval Realisasi, ID : "+idPengajuan,
		user,
		approvalPosisi,
		"approval_realisasi",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi approve realisasi: %w", err)
	}
	return nil
}

// RejectRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/reject.
func (r *realisasiKpiRepo) RejectRealisasiKpi(idPengajuan, approvalList, catatan, user string) error {
	// Ambil entry_user_realisasi untuk notifikasi penolakan ke pengaju
	var kpiBase struct {
		EntryUserRealisasi string `gorm:"column:entry_user_realisasi"`
	}
	if err := r.db.Raw(queryGetKpiBaseData, idPengajuan).Scan(&kpiBase).Error; err != nil {
		return fmt.Errorf("gagal mengambil entry_user_realisasi: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryRejectRealisasi,
		approvalList, catatan, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal reject realisasi: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Realisasi Ditolak, ID : "+idPengajuan,
		user,
		kpiBase.EntryUserRealisasi,
		"realisasi_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi reject realisasi: %w", err)
	}
	return nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all.
func (r *realisasiKpiRepo) GetAllRealisasiKpi(
	req *dto.GetAllRealisasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"a.status = 2",
	}
	args := []interface{}{}

	// =========================================================================
	// Kondisi opsional dari request body
	// =========================================================================
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
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-approval.
func (r *realisasiKpiRepo) GetAllApprovalRealisasiKpi(
	req *dto.GetAllApprovalRealisasiKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status = 3",
		"a.approval_posisi = ?",
	}
	args := []interface{}{req.ApprovalUserRealisasi}

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
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllTolakanRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-tolakan.
func (r *realisasiKpiRepo) GetAllTolakanRealisasiKpi(
	req *dto.GetAllTolakanRealisasiKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status = 4",
		"a.entry_user_realisasi = ?",
	}
	args := []interface{}{req.EntryUserRealisasi}

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

// GetAllDaftarRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-realisasi.
func (r *realisasiKpiRepo) GetAllDaftarRealisasiKpi(
	req *dto.GetAllDaftarRealisasiKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status IN (2, 3, 4, 5, 80)",
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

// GetAllDaftarApprovalRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-all-daftar-approval.
func (r *realisasiKpiRepo) GetAllDaftarApprovalRealisasiKpi(
	req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
) ([]*model.DataKpi, int64, error) {

	// =========================================================================
	// BUILD DYNAMIC WHERE
	// =========================================================================
	conditions := []string{
		"a.status IN (3, 5)",
		"a.approval_list_realisasi != ''",
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
// GET DETAIL
// =============================================================================

// GetDetailRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/get-detail.
func (r *realisasiKpiRepo) GetDetailRealisasiKpi(
	req *dto.GetDetailRealisasiKpiRequest,
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
	// QUERY DATA KPI
	// =========================================================================
	var result model.DataKpi
	headerQuery := queryGetDataKpi + where + " LIMIT 1"
	if err := r.db.Raw(headerQuery, args...).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil detail realisasi KPI: %w", err)
	}

	// =========================================================================
	// QUERY KPI DETAIL
	// =========================================================================
	var kpiDetails []model.DataKpiDetail
	if err := r.db.Raw(queryGetDataKpiDetail, result.IdPengajuan).Scan(&kpiDetails).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil kpi detail: %w", err)
	}

	// =========================================================================
	// QUERY KPI SUB DETAIL
	// =========================================================================
	for i := range kpiDetails {
		var subDetails []model.DataKpiSubDetail
		if err := r.db.Raw(queryGetDataKpiSubDetail, result.IdPengajuan, kpiDetails[i].IdDetail).Scan(&subDetails).Error; err != nil {
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
// APPROVAL HELPER
// =============================================================================

// GetApprovalListJSON digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi untuk mengambil daftar approval dalam format JSON.
func (r *realisasiKpiRepo) GetApprovalListJSON(idPengajuan, userID string) (string, error) {
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

// CheckApprovalRealisasiExists digunakan oleh service ApproveRealisasiKpi dan RejectRealisasiKpi untuk memvalidasi keberadaan approval.
func (r *realisasiKpiRepo) CheckApprovalRealisasiExists(user, idPengajuan string) (bool, error) {
	var count int64
	if err := r.db.Raw(queryGetCountApprovalKpi, user, idPengajuan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	return count > 0, nil
}

// GetCatatanTolakan digunakan oleh service RejectRealisasiKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
func (r *realisasiKpiRepo) GetCatatanTolakan(idPengajuan string) (string, error) {
	var val []byte
	row := r.db.Raw(queryGetCatatanTolakan, idPengajuan).Row()
	if err := row.Scan(&val); err != nil {
		return "", err
	}
	return string(val), nil
}

// =============================================================================
// GET EXIST
// =============================================================================

// GetExistDataKpi digunakan oleh service untuk mengambil header KPI berdasarkan id_pengajuan.
func (r *realisasiKpiRepo) GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error) {
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

// =============================================================================
// SERVICE HELPERS
// =============================================================================

// GetLinkFormats digunakan oleh service ValidateRealisasiKpi dan RevisionPenyusunanKpi untuk mengambil daftar format link yang valid dari master data.
func (r *realisasiKpiRepo) GetLinkFormats() ([]string, error) {
	rows, err := r.db.Raw(queryGetLinkFormats).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data format link: %w", err)
	}
	defer rows.Close()

	var prefixes []string
	for rows.Next() {
		var prefix string
		if err := rows.Scan(&prefix); err != nil {
			return nil, fmt.Errorf("gagal scan format link: %w", err)
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes, nil
}

// LookupSubDetailByKpiSubKpi digunakan oleh service ValidateRealisasiKpi dan RevisionPenyusunanKpi untuk mencari data sub detail berdasarkan id_pengajuan, kpi_name, dan sub_kpi_name dari Excel.
func (r *realisasiKpiRepo) LookupSubDetailByKpiSubKpi(
	idPengajuan, kpiName, subKpiName string,
) (*model.SubDetailLookup, error) {
	var result model.SubDetailLookup
	if err := r.db.Raw(queryLookupSubDetail, idPengajuan, kpiName, subKpiName).Scan(&result).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf(
				"sub KPI '%s' pada KPI '%s' tidak ditemukan di id_pengajuan '%s'",
				subKpiName, kpiName, idPengajuan,
			)
		}
		return nil, fmt.Errorf("gagal lookup sub detail untuk sub KPI '%s': %w", subKpiName, err)
	}
	if result.IdSubDetail == "" {
		return nil, fmt.Errorf(
			"sub KPI '%s' pada KPI '%s' tidak ditemukan di id_pengajuan '%s'",
			subKpiName, kpiName, idPengajuan,
		)
	}
	return &result, nil
}
