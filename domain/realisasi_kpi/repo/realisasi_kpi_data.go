package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dto "permen_api/domain/realisasi_kpi/dto"
	customErrors "permen_api/errors"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Check
	// =============================================================================

	// Use func : ValidateRealisasiKpi
	// yang mengizinkan input realisasi (2=approved penyusunan, 4=tolak realisasi, 80=draft realisasi, 81).
	queryCheckExistRealisasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status IN (2, 4, 80, 81)`

	// queryCheckStatusRevisiRealisasi memvalidasi bahwa id_pengajuan ada dan berstatus
	// yang mengizinkan revisi realisasi (4=tolak realisasi, 80=draft realisasi).
	queryCheckStatusRevisiRealisasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status IN (4, 80)`

	// queryCheckStatusCreateRealisasi memvalidasi status 80 (draft) sebelum create/submit.
	queryCheckStatusCreateRealisasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status = 80`

	// queryCheckStatusApprovalRealisasi memvalidasi status 3 (pending approval) sebelum approval.
	queryCheckStatusApprovalRealisasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status = 3`

	// =============================================================================
	// Lookup
	// =============================================================================

	// queryLookupSubDetail mencari id_sub_detail, id_detail, target_kuantitatif_triwulan,
	// dan rumus (id_polarisasi) berdasarkan id_pengajuan + nama kpi (dari detail) + nama sub_kpi.
	queryLookupSubDetail = `
		SELECT
			s.id_sub_detail,
			s.id_detail,
			s.rumus,
			IFNULL(s.target_kuantitatif_triwulan, 0) AS target_kuantitatif_triwulan
		FROM data_kpi_subdetail s
		INNER JOIN data_kpi_detail d ON d.id_detail = s.id_detail
		WHERE s.id_pengajuan = ?
		  AND LOWER(d.kpi)  = LOWER(?)
		  AND LOWER(s.kpi)  = LOWER(?)
		LIMIT 1`

	// queryGetKpiHeaderRealisasi mengambil kostl_tx, tahun, triwulan untuk keperluan notifikasi.
	queryGetKpiHeaderRealisasi = `
		SELECT kostl_tx, tahun, triwulan, entry_user_realisasi
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// =============================================================================
	// Update — Validate / Revision
	// =============================================================================

	// queryUpdateKpiStatusDraft meng-update header data_kpi ke status 80 (draft realisasi).
	queryUpdateKpiStatusDraft = `
		UPDATE data_kpi
		SET status                  = 80,
		    entry_user_realisasi    = ?,
		    entry_name_realisasi    = ?,
		    entry_time_realisasi    = NOW()
		WHERE id_pengajuan = ?`

	// queryUpdateSubDetailRealisasi meng-update satu baris data_kpi_subdetail dengan nilai realisasi.
	queryUpdateSubDetailRealisasi = `
		UPDATE data_kpi_subdetail
		SET realisasi                      = ?,
		    realisasi_kuantitatif          = ?,
		    realisasi_validated            = ?,
		    realisasi_kuantitatif_validated = ?,
		    realisasi_keterangan           = '',
		    pencapaian                     = ?,
		    skor                           = ?,
		    realisasi_qualifier            = ?,
		    realisasi_kuantitatif_qualifier = ?
		WHERE id_pengajuan = ?
		  AND id_sub_detail = ?`

	// queryUpdateChallengeRealisasi meng-update realisasi pada data_challenge_detail.
	queryUpdateChallengeRealisasi = `
		UPDATE data_challenge_detail
		SET realisasi_challenge = ?
		WHERE id_pengajuan = ?
		  AND id_detail_challenge = ?`

	// queryUpdateMethodRealisasi meng-update realisasi pada data_method_detail.
	queryUpdateMethodRealisasi = `
		UPDATE data_method_detail
		SET realisasi_method = ?
		WHERE id_pengajuan = ?
		  AND id_detail_method = ?`

	// =============================================================================
	// Update — Create (Submit)
	// =============================================================================

	// querySubmitRealisasi meng-update header data_kpi ke status 3 (pending approval realisasi).
	querySubmitRealisasi = `
		UPDATE data_kpi
		SET status                   = 3,
		    approval_posisi          = ?,
		    approval_list_realisasi  = ?
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Update — Approval
	// =============================================================================

	queryApproveChainRealisasi = `
		UPDATE data_kpi
		SET approval_posisi         = ?,
		    approval_list_realisasi = ?
		WHERE id_pengajuan = ?`

	queryApproveFinalRealisasi = `
		UPDATE data_kpi
		SET status                  = 5,
		    approval_list_realisasi = ?
		WHERE id_pengajuan = ?`

	queryRejectRealisasi = `
		UPDATE data_kpi
		SET status                  = 4,
		    approval_list_realisasi = ?,
		    catatan_tolakan         = ?
		WHERE id_pengajuan = ?`

	// =============================================================================
	// GetAll
	// =============================================================================

	queryGetCountDataKpiRealisasi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	queryGetDataKpiRealisasi = `
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

	// queryGetDetailHeader mengambil header satu pengajuan untuk GetDetailRealisasiKpi.
	queryGetDetailHeader = `
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
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ?
		LIMIT 1`

	queryGetDetailKpiList = `
		SELECT id_detail, id_kpi, kpi, rumus
		FROM data_kpi_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail ASC`

	queryGetDetailSubKpiList = `
		SELECT
			a.id_sub_detail, a.id_kpi, a.kpi AS sub_kpi, a.otomatis,
			IFNULL(CAST(a.bobot AS CHAR), '')                        bobot,
			IFNULL(a.capping, '')                                    capping,
			IFNULL(a.target_triwulan, '')                            target_triwulan,
			IFNULL(CAST(a.target_kuantitatif_triwulan AS CHAR), '')  target_kuantitatif_triwulan,
			IFNULL(a.target_tahunan, '')                             target_tahunan,
			IFNULL(CAST(a.target_kuantitatif_tahunan AS CHAR), '')   target_kuantitatif_tahunan,
			IFNULL(a.realisasi, '')                                  realisasi,
			IFNULL(a.realisasi_kuantitatif, '')                      realisasi_kuantitatif,
			IFNULL(a.realisasi_keterangan, '')                       realisasi_keterangan,
			IFNULL(a.realisasi_validated, '')                        realisasi_validated,
			IFNULL(a.realisasi_kuantitatif_validated, '')            realisasi_kuantitatif_validated,
			IFNULL(CAST(a.pencapaian AS CHAR), '')                   pencapaian,
			IFNULL(CAST(a.skor AS CHAR), '')                         skor,
			IFNULL(a.deskripsi_glossary, '')                         deskripsi_glossary,
			IFNULL(a.item_qualifier, '')                             item_qualifier,
			IFNULL(a.deskripsi_qualifier, '')                        deskripsi_qualifier,
			IFNULL(a.target_qualifier, '')                           target_qualifier,
			IFNULL(a.id_qualifier, '')                               id_qualifier,
			IFNULL(a.realisasi_qualifier, '')                        realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '')            realisasi_kuantitatif_qualifier
		FROM data_kpi_subdetail a
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	queryGetDetailContextList = `
		SELECT
			id_detail_challenge, nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '') realisasi_challenge,
			IFNULL(lampiran_evidence, '')   lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_challenge ASC`

	queryGetDetailProcessList = `
		SELECT
			id_detail_method, nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')  realisasi_method,
			IFNULL(lampiran_evidence, '') lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_method ASC`
)

// =============================================================================
// LOOKUP
// =============================================================================

func (r *realisasiKpiRepo) GetTriwulanByIdPengajuan(idPengajuan string) (string, error) {
	row := r.db.Raw(queryGetKpiHeaderRealisasi, idPengajuan).Row()
	var kostlTx, tahun, triwulan, entryUserRealisasi string
	if err := row.Scan(&kostlTx, &tahun, &triwulan, &entryUserRealisasi); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &customErrors.BadRequestError{
				Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
			}
		}
		return "", fmt.Errorf("gagal mengambil triwulan untuk id_pengajuan '%s': %w", idPengajuan, err)
	}
	return triwulan, nil
}

func (r *realisasiKpiRepo) GetKpiHeaderByIdPengajuan(
	idPengajuan string,
) (tahun, triwulan, kostl, kostlTx string, err error) {
	row := r.db.Raw(`
		SELECT kostl_tx, tahun, triwulan, IFNULL(kostl, '') AS kostl
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`, idPengajuan).Row()
	if scanErr := row.Scan(&kostlTx, &tahun, &triwulan, &kostl); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return "", "", "", "", &customErrors.BadRequestError{
				Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
			}
		}
		return "", "", "", "", fmt.Errorf("gagal mengambil header kpi '%s': %w", idPengajuan, scanErr)
	}
	return tahun, triwulan, kostl, kostlTx, nil
}

func (r *realisasiKpiRepo) LookupSubDetailByKpiSubKpi(
	idPengajuan, kpiName, subKpiName string,
) (idSubDetail, idDetail, rumus string, targetKuantitatifTriwulan float64, err error) {
	row := r.db.Raw(queryLookupSubDetail, idPengajuan, kpiName, subKpiName).Row()
	if scanErr := row.Scan(&idSubDetail, &idDetail, &rumus, &targetKuantitatifTriwulan); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return "", "", "", 0, &customErrors.BadRequestError{
				Message: fmt.Sprintf(
					"sub KPI '%s' pada KPI '%s' tidak ditemukan di id_pengajuan '%s'",
					subKpiName, kpiName, idPengajuan,
				),
			}
		}
		return "", "", "", 0, fmt.Errorf("gagal lookup sub detail untuk sub KPI '%s': %w", subKpiName, scanErr)
	}
	return idSubDetail, idDetail, rumus, targetKuantitatifTriwulan, nil
}

// =============================================================================
// VALIDATE — simpan draft realisasi (status 80)
// =============================================================================

func (r *realisasiKpiRepo) ValidateRealisasiKpi(
	req *dto.ValidateRealisasiKpiRequest,
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
	resultList []dto.RealisasiResult,
	processList []dto.RealisasiProcess,
	contextList []dto.RealisasiContext,
) error {
	// Validasi status: harus 2, 4, 80, atau 81 agar bisa input realisasi
	var countExist int
	if err := r.db.Raw(queryCheckExistRealisasi, req.IdPengajuan).Scan(&countExist).Error; err != nil {
		return fmt.Errorf("gagal mengecek data Realisasi KPI: %w", err)
	}
	if countExist == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf(
				"id_pengajuan '%s' tidak ditemukan atau status tidak mengizinkan input realisasi",
				req.IdPengajuan,
			),
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// -------------------------------------------------------------------------
	// UPDATE setiap baris sub KPI
	// -------------------------------------------------------------------------
	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if err := tx.Exec(queryUpdateSubDetailRealisasi,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Realisasi,            // realisasi_validated = sama dengan realisasi
				row.RealisasiKuantitatif, // realisasi_kuantitatif_validated = sama
				row.Pencapaian,
				row.Skor,
				row.RealisasiQualifierVal,
				row.RealisasiKuantitatifQualifier,
				req.IdPengajuan,
				row.IdSubDetail,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update sub detail '%s': %w", row.IdSubDetail, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE result (data_challenge_detail.realisasi_challenge) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, r2 := range resultList {
		if err := tx.Exec(queryUpdateChallengeRealisasi,
			r2.RealisasiResult, req.IdPengajuan, r2.IdDetailResult,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update challenge result '%s': %w", r2.IdDetailResult, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE process (data_method_detail.realisasi_method) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, p := range processList {
		if err := tx.Exec(queryUpdateMethodRealisasi,
			p.RealisasiProcess, req.IdPengajuan, p.IdDetailProcess,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update method process '%s': %w", p.IdDetailProcess, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE context (data_challenge_detail.realisasi_challenge untuk context) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, c := range contextList {
		if err := tx.Exec(queryUpdateChallengeRealisasi,
			c.RealisasiContext, req.IdPengajuan, c.IdDetailContext,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update challenge context '%s': %w", c.IdDetailContext, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE header data_kpi → status 80
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryUpdateKpiStatusDraft,
		req.EntryUser, req.EntryName, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update header realisasi: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi validate realisasi: %w", err)
	}

	return nil
}

// =============================================================================
// REVISION — update ulang realisasi (status 80 atau 4)
// =============================================================================

func (r *realisasiKpiRepo) RevisionRealisasiKpi(
	req *dto.RevisionRealisasiKpiRequest,
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
	resultList []dto.RealisasiResult,
	processList []dto.RealisasiProcess,
	contextList []dto.RealisasiContext,
) error {
	// Validasi status: harus 80 (draft) atau 4 (ditolak) untuk revisi
	var count int
	if err := r.db.Raw(queryCheckStatusRevisiRealisasi, req.IdPengajuan).Scan(&count).Error; err != nil {
		return fmt.Errorf("gagal mengecek status pengajuan: %w", err)
	}
	if count == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf(
				"id_pengajuan '%s' tidak ditemukan atau status tidak mengizinkan revisi realisasi (harus draft atau ditolak)",
				req.IdPengajuan,
			),
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	// -------------------------------------------------------------------------
	// UPDATE setiap baris sub KPI
	// -------------------------------------------------------------------------
	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if err := tx.Exec(queryUpdateSubDetailRealisasi,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Realisasi,
				row.RealisasiKuantitatif,
				row.Pencapaian,
				row.Skor,
				row.RealisasiQualifierVal,
				row.RealisasiKuantitatifQualifier,
				req.IdPengajuan,
				row.IdSubDetail,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update sub detail '%s': %w", row.IdSubDetail, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE result (data_challenge_detail.realisasi_challenge) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, r2 := range resultList {
		if err := tx.Exec(queryUpdateChallengeRealisasi,
			r2.RealisasiResult, req.IdPengajuan, r2.IdDetailResult,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update challenge result '%s': %w", r2.IdDetailResult, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE process (data_method_detail.realisasi_method) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, p := range processList {
		if err := tx.Exec(queryUpdateMethodRealisasi,
			p.RealisasiProcess, req.IdPengajuan, p.IdDetailProcess,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update method process '%s': %w", p.IdDetailProcess, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE context (data_challenge_detail.realisasi_challenge untuk context) — TW2/TW4
	// -------------------------------------------------------------------------
	for _, c := range contextList {
		if err := tx.Exec(queryUpdateChallengeRealisasi,
			c.RealisasiContext, req.IdPengajuan, c.IdDetailContext,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update challenge context '%s': %w", c.IdDetailContext, err)
		}
	}

	// -------------------------------------------------------------------------
	// UPDATE header → status tetap 80 (draft), update entry time
	// -------------------------------------------------------------------------
	if err := tx.Exec(queryUpdateKpiStatusDraft,
		req.EntryUser, req.EntryName, req.IdPengajuan,
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
// CREATE — submit realisasi ke approval (status 80 → 3)
// =============================================================================

func (r *realisasiKpiRepo) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) error {
	var count int
	if err := r.db.Raw(queryCheckStatusCreateRealisasi, req.IdPengajuan).Scan(&count).Error; err != nil {
		return fmt.Errorf("gagal mengecek status pengajuan: %w", err)
	}
	if count == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf(
				"id_pengajuan '%s' tidak ditemukan atau belum dalam status draft realisasi (80)",
				req.IdPengajuan,
			),
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(querySubmitRealisasi,
		req.ApprovalPosisi, req.ApprovalListRealisasi, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal submit realisasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Approval Realisasi, ID : "+req.IdPengajuan,
		req.User,
		req.ApprovalPosisi,
		"approval_realisasi",
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// =============================================================================
// APPROVAL — approve atau reject realisasi
// =============================================================================

func (r *realisasiKpiRepo) ApprovalRealisasiKpi(
	req *dto.ApprovalRealisasiKpiRequest,
) error {
	var count int
	if err := r.db.Raw(queryCheckStatusApprovalRealisasi, req.IdPengajuan).Scan(&count).Error; err != nil {
		return fmt.Errorf("gagal mengecek status pengajuan: %w", err)
	}
	if count == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf(
				"id_pengajuan '%s' tidak ditemukan atau tidak dalam status menunggu approval realisasi (3)",
				req.IdPengajuan,
			),
		}
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if req.Status == "approve" {
		// ApprovalPosisi kosong → approval final (status 5)
		if req.ApprovalPosisi == "" {
			if err := tx.Exec(queryApproveFinalRealisasi,
				req.ApprovalList, req.IdPengajuan,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal approve final realisasi: %w", err)
			}
			tx.Commit()
			return nil
		}

		// Masih ada approver berikutnya → chain approval
		if err := tx.Exec(queryApproveChainRealisasi,
			req.ApprovalPosisi, req.ApprovalList, req.IdPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve chain realisasi: %w", err)
		}

		if err := notif.Insert(tx,
			req.IdPengajuan,
			"Approval Realisasi, ID : "+req.IdPengajuan,
			req.User,
			req.ApprovalPosisi,
			"approval_realisasi",
		); err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}

	if req.Status == "reject" {
		// Ambil entry_user_realisasi untuk kirim notifikasi ke pengaju
		var header struct {
			EntryUserRealisasi string `gorm:"column:entry_user_realisasi"`
		}
		if err := r.db.Raw(queryGetKpiHeaderRealisasi, req.IdPengajuan).Scan(&header).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal mengambil entry_user_realisasi: %w", err)
		}

		if err := tx.Exec(queryRejectRealisasi,
			req.ApprovalList, req.CatatanTolak, req.IdPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal reject realisasi: %w", err)
		}

		if err := notif.Insert(tx,
			req.IdPengajuan,
			"Realisasi Ditolak, ID : "+req.IdPengajuan,
			req.User,
			header.EntryUserRealisasi,
			"realisasi_ditolak",
		); err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}

	tx.Rollback()
	return &customErrors.BadRequestError{Message: "status tidak valid, gunakan 'approve' atau 'reject'"}
}

// =============================================================================
// GET ALL APPROVAL — status 3, approval_posisi = user
// =============================================================================

func (r *realisasiKpiRepo) GetAllApprovalRealisasiKpi(
	req *dto.GetAllApprovalRealisasiKpiRequest,
) ([]*dto.DataKpiRealisasi, int64, error) {
	conditions := []string{
		"a.status = 3",
		"a.approval_posisi = ?",
	}
	args := []interface{}{req.ApprovalUser}

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

	return r.queryPaginatedRealisasi(conditions, args, req.Page, req.Limit)
}

// =============================================================================
// GET ALL TOLAKAN — status 4, entry_user_realisasi = user
// =============================================================================

func (r *realisasiKpiRepo) GetAllTolakanRealisasiKpi(
	req *dto.GetAllTolakanRealisasiKpiRequest,
) ([]*dto.DataKpiRealisasi, int64, error) {
	conditions := []string{
		"a.status = 4",
		"a.entry_user_realisasi = ?",
	}
	args := []interface{}{req.EntryUserRealisasi}

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

	return r.queryPaginatedRealisasi(conditions, args, req.Page, req.Limit)
}

// =============================================================================
// GET ALL DAFTAR REALISASI — semua status realisasi (2, 3, 4, 5, 80)
// =============================================================================

func (r *realisasiKpiRepo) GetAllDaftarRealisasiKpi(
	req *dto.GetAllDaftarRealisasiKpiRequest,
) ([]*dto.DataKpiRealisasi, int64, error) {
	conditions := []string{
		"a.status IN (2, 3, 4, 5, 80)",
	}
	args := []interface{}{}

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

	return r.queryPaginatedRealisasi(conditions, args, req.Page, req.Limit)
}

// =============================================================================
// GET ALL DAFTAR APPROVAL — semua yang ada approval_list_realisasi
// =============================================================================

func (r *realisasiKpiRepo) GetAllDaftarApprovalRealisasiKpi(
	req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
) ([]*dto.DataKpiRealisasi, int64, error) {
	conditions := []string{
		"a.status IN (3, 5)",
		"a.approval_list_realisasi != ''",
	}
	args := []interface{}{}

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

	return r.queryPaginatedRealisasi(conditions, args, req.Page, req.Limit)
}

// =============================================================================
// GET DETAIL
// =============================================================================

func (r *realisasiKpiRepo) GetDetailRealisasiKpi(
	req *dto.GetDetailRealisasiKpiRequest,
) (*dto.GetDetailRealisasiKpiResponse, error) {

	// =========================================================================
	// SCAN HEADER
	// =========================================================================
	headerRow := r.db.Raw(queryGetDetailHeader, req.IdPengajuan).Row()

	var (
		idPengajuan, tahun, triwulan                string
		kostl, kostlTx, orgeh, orgehTx              string
		entryUser, entryName, entryTime             string
		approvalPosisi, approvalListRaw             string
		status                                      int
		statusDesc                                  string
		entryUserR, entryNameR, entryTimeR          string
		approvalListRealisasiRaw                    string
		catatanTolakan, totalBobot, totalPencapaian string
	)

	if err := headerRow.Scan(
		&idPengajuan, &tahun, &triwulan,
		&kostl, &kostlTx, &orgeh, &orgehTx,
		&entryUser, &entryName, &entryTime,
		&approvalPosisi, &approvalListRaw,
		&status, &statusDesc,
		&entryUserR, &entryNameR, &entryTimeR,
		&approvalListRealisasiRaw,
		&catatanTolakan, &totalBobot, &totalPencapaian,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &customErrors.BadRequestError{
				Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
			}
		}
		return nil, fmt.Errorf("gagal scan header detail realisasi: %w", err)
	}

	// =========================================================================
	// UNMARSHAL approval_list dan approval_list_realisasi
	// =========================================================================
	var approvalList []dto.ApprovalDetailResponse
	if approvalListRaw != "" {
		_ = json.Unmarshal([]byte(approvalListRaw), &approvalList)
	}
	if approvalList == nil {
		approvalList = []dto.ApprovalDetailResponse{}
	}

	var approvalListRealisasi []dto.ApprovalDetailResponse
	if approvalListRealisasiRaw != "" {
		_ = json.Unmarshal([]byte(approvalListRealisasiRaw), &approvalListRealisasi)
	}
	if approvalListRealisasi == nil {
		approvalList = []dto.ApprovalDetailResponse{}
	}

	resp := &dto.GetDetailRealisasiKpiResponse{
		IdPengajuan: idPengajuan,
		Tahun:       tahun,
		Triwulan:    triwulan,
		Kostl:       kostl,
		KostlTx:     kostlTx,
		Orgeh:       orgeh,
		OrgehTx:     orgehTx,
		Status:      status,
		StatusDesc:  statusDesc,
		EntryPenyusunan: dto.EntryDetailResponse{
			EntryUser: entryUser,
			EntryName: entryName,
			EntryTime: entryTime,
		},
		EntryRealisasi: dto.EntryDetailResponse{
			EntryUser: entryUserR,
			EntryName: entryNameR,
			EntryTime: entryTimeR,
		},
		ApprovalList:          approvalList,
		ApprovalListRealisasi: approvalListRealisasi,
		CatatanTolakan:        catatanTolakan,
		TotalBobot:            totalBobot,
		TotalPencapaian:       totalPencapaian,
	}

	// =========================================================================
	// SCAN KPI DETAIL LIST
	// =========================================================================
	kpiRows, err := r.db.Raw(queryGetDetailKpiList, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil kpi detail list: %w", err)
	}
	defer kpiRows.Close()

	var kpiList []dto.DetailKpiRealisasi
	totalSubKpi := 0

	for kpiRows.Next() {
		var kpi dto.DetailKpiRealisasi
		if err := kpiRows.Scan(&kpi.IdDetail, &kpi.IdKpi, &kpi.Kpi, &kpi.Rumus); err != nil {
			return nil, fmt.Errorf("gagal scan kpi detail: %w", err)
		}

		// Sub KPI per detail
		subRows, err := r.db.Raw(queryGetDetailSubKpiList, idPengajuan, kpi.IdDetail).Rows()
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil sub kpi list untuk detail '%s': %w", kpi.IdDetail, err)
		}
		defer subRows.Close()

		var subDetailList []dto.DetailSubKpiRealisasi
		for subRows.Next() {
			var s dto.DetailSubKpiRealisasi
			if err := subRows.Scan(
				&s.IdSubDetail, &s.IdKpi, &s.SubKpi, &s.Otomatis,
				&s.Bobot, &s.Capping,
				&s.TargetTriwulan, &s.TargetKuantitatifTriwulan,
				&s.TargetTahunan, &s.TargetKuantitatifTahunan,
				&s.Realisasi, &s.RealisasiKuantitatif, &s.RealisasiKeterangan,
				&s.RealisasiValidated, &s.RealisasiKuantitatifValidated,
				&s.Pencapaian, &s.Skor,
				&s.DeskripsiGlossary, &s.ItemQualifier, &s.DeskripsiQualifier, &s.TargetQualifier,
				&s.IdQualifier, &s.RealisasiQualifier, &s.RealisasiKuantitatifQualifier,
			); err != nil {
				return nil, fmt.Errorf("gagal scan sub kpi detail: %w", err)
			}
			subDetailList = append(subDetailList, s)
		}

		kpi.TotalSubKpi = len(subDetailList)
		kpi.SubDetailList = subDetailList
		totalSubKpi += kpi.TotalSubKpi
		kpiList = append(kpiList, kpi)
	}

	resp.TotalSubKpi = totalSubKpi
	resp.KpiList = kpiList

	// =========================================================================
	// SCAN CONTEXT LIST (data_challenge_detail)
	// =========================================================================
	ctxRows, err := r.db.Raw(queryGetDetailContextList, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil context list: %w", err)
	}
	defer ctxRows.Close()

	var contextList []dto.DetailContextRealisasi
	for ctxRows.Next() {
		var c dto.DetailContextRealisasi
		if err := ctxRows.Scan(
			&c.IdDetailChallenge, &c.NamaChallenge, &c.DeskripsiChallenge,
			&c.RealisasiChallenge, &c.LampiranEvidence,
		); err != nil {
			return nil, fmt.Errorf("gagal scan context list: %w", err)
		}
		contextList = append(contextList, c)
	}
	resp.ContextList = contextList

	// =========================================================================
	// SCAN PROCESS LIST (data_method_detail)
	// =========================================================================
	procRows, err := r.db.Raw(queryGetDetailProcessList, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil process list: %w", err)
	}
	defer procRows.Close()

	var processList []dto.DetailProcessRealisasi
	for procRows.Next() {
		var p dto.DetailProcessRealisasi
		if err := procRows.Scan(
			&p.IdDetailMethod, &p.NamaMethod, &p.DeskripsiMethod,
			&p.RealisasiMethod, &p.LampiranEvidence,
		); err != nil {
			return nil, fmt.Errorf("gagal scan process list: %w", err)
		}
		processList = append(processList, p)
	}
	resp.ProcessList = processList

	return resp, nil
}

// =============================================================================
// HELPER — paginated query reusable untuk semua GetAll
// =============================================================================

func (r *realisasiKpiRepo) queryPaginatedRealisasi(
	conditions []string,
	args []interface{},
	page, limit int,
) ([]*dto.DataKpiRealisasi, int64, error) {
	where := " WHERE " + strings.Join(conditions, " AND ")

	var total int64
	countQuery := queryGetCountDataKpiRealisasi + where
	if err := r.db.Raw(countQuery, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data realisasi: %w", err)
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	listQuery := queryGetDataKpiRealisasi + where +
		" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	dbRows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil data realisasi: %w", err)
	}
	defer dbRows.Close()

	var results []*dto.DataKpiRealisasi
	for dbRows.Next() {
		var h dto.DataKpiRealisasi
		if err := dbRows.Scan(
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
			return nil, 0, fmt.Errorf("gagal scan data realisasi: %w", err)
		}
		results = append(results, &h)
	}

	return results, total, nil
}
