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
	// Check
	// =============================================================================

	// Use func : ValidateRealisasiKpi
	// yang mengizinkan input realisasi (2=approved penyusunan, 4=tolak realisasi, 80=draft realisasi, 81).
	queryCheckExistRealisasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND tahun = ? AND triwulan = ? AND kostl = ? AND status IN (2, 4, 80, 81)`

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

	// queryGetKpiBaseData digunakan oleh: SubmitPenyusunanKpi (kostl_tx),
	// ApprovePenyusunanKpi (entry_user), dan GetKpiExportData (kostl_tx, tahun, triwulan).
	queryGetKpiBaseData = `
		SELECT kostl_tx, tahun, triwulan, entry_user
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

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
	// Update
	// =============================================================================
	// queryUpdateKpiRealisasi digunakan oleh CreateRealisasiKpi untuk mengisi approval dan mengubah status.
	queryUpdateKpiRealisasi = `
		UPDATE data_kpi 
		SET approval_posisi = ?, approval_list_realisasi = ?, status = 2
		WHERE id_pengajuan = ?`

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

	// queryCheckApprovalRealisasi memvalidasi bahwa user adalah approval_posisi aktif dan status = 3.
	queryCheckApprovalRealisasi = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 3 AND approval_posisi = ? AND id_pengajuan = ?`

	// Use func : GetAllApprovalRealisasiKpi, GetAllTolakanRealisasiKpi, GetAllDaftarRealisasiKpi, GetAllDaftarApprovalRealisasiKpi
	queryGetCountDataKpiRealisasi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// =============================================================================
	// GetAll
	// =============================================================================

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
		WHERE a.id_pengajuan = ?
		ORDER BY a.id_detail ASC`

	queryGetDetailSubKpiList = `
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
			IFNULL(a.realisasi, '')                       realisasi,
			IFNULL(a.realisasi_kuantitatif, 0)            realisasi_kuantitatif,
			IFNULL(a.realisasi_keterangan, '')            realisasi_keterangan,
			IFNULL(a.realisasi_validated, '')             realisasi_validated,
			IFNULL(a.realisasi_kuantitatif_validated, '') realisasi_kuantitatif_validated,
			IFNULL(a.pencapaian, 0)                       pencapaian,
			IFNULL(a.skor, 0)                             skor,
			IFNULL(a.realisasi_qualifier, '')             realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '') realisasi_kuantitatif_qualifier
		FROM data_kpi_subdetail a
		LEFT JOIN mst_polarisasi p ON a.rumus = p.id_polarisasi
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	queryGetDetailResultList = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result,
			IFNULL(realisasi_result, '')   realisasi_result,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_result_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_result ASC`

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

func (r *realisasiKpiRepo) GetKpiHeader(idPengajuan string) (tahun, triwulan, kostl, kostlTx, entryUser, entryName string, status int, statusDesc string, err error) {
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
	if err = r.db.Raw(`
		SELECT a.tahun, a.triwulan, a.kostl, a.kostl_tx, a.status, b.status_desc, a.entry_user, a.entry_name
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ?
		LIMIT 1`, idPengajuan).Scan(&h).Error; err != nil {
		return "", "", "", "", "", "", 0, "", fmt.Errorf("gagal mengambil header KPI: %w", err)
	}
	if h.Tahun == "" {
		return "", "", "", "", "", "", 0, "", fmt.Errorf("id_pengajuan '%s' tidak ditemukan", idPengajuan)
	}
	return h.Tahun, h.Triwulan, h.Kostl, h.KostlTx, h.EntryUser, h.EntryName, h.Status, h.StatusDesc, nil
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
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
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
		req.EntryUserRealisasi, req.EntryNameRealisasi, req.IdPengajuan,
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
// CREATE — submit realisasi ke approval (status 80 → 3)
// =============================================================================

func (r *realisasiKpiRepo) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) error {
	var countExist int
	if err := r.db.Raw(queryCheckStatusCreateRealisasi, req.IdPengajuan).
		Scan(&countExist).Error; err != nil {
		return fmt.Errorf("gagal mengecek status pengajuan: %w", err)
	}
	if countExist == 0 {
		return &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau belum dalam status draft realisasi (80)", req.IdPengajuan),
		}
	}

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
	kostlTx := kpiBase.KostlTx

	// =========================================================================
	// Jalankan dalam transaksi agar update + notif atomic
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(querySubmitRealisasi,
		approvalPosisi, string(approvalListBytes), req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal submit realisasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		kostlTx,
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
// REVISION — update ulang realisasi (status 80 atau 4)
// =============================================================================

func (r *realisasiKpiRepo) RevisionRealisasiKpi(
	req *dto.RevisionRealisasiKpiRequest,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
	resultList []dto.DataResult,
	processList []dto.DataProcess,
	contextList []dto.DataContext,
) error {
	// Validasi status: harus 80 (draft) atau 4 (ditolak) untuk revisi
	var countExist int
	if err := r.db.Raw(queryCheckStatusRevisiRealisasi, req.IdPengajuan).
		Scan(&countExist).Error; err != nil {
		return fmt.Errorf("gagal mengecek status pengajuan: %w", err)
	}
	if countExist == 0 {
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
		req.EntryUserRealisasi, req.EntryNameRealisasi, req.IdPengajuan,
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
// APPROVE REALISASI KPI
// =============================================================================

func (r *realisasiKpiRepo) ApproveRealisasiKpi(idPengajuan, approvalList, approvalPosisi, user string) error {
	var count int64
	if err := r.db.Raw(queryCheckApprovalRealisasi, user, idPengajuan).Scan(&count).Error; err != nil {
		return fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	if count == 0 {
		return &customErrors.BadRequestError{Message: "Data Not Found"}
	}

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
		tx.Commit()
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

	tx.Commit()
	return nil
}

// =============================================================================
// REJECT REALISASI KPI
// =============================================================================

func (r *realisasiKpiRepo) RejectRealisasiKpi(idPengajuan, approvalList, catatan, user string) error {
	var count int64
	if err := r.db.Raw(queryCheckApprovalRealisasi, user, idPengajuan).Scan(&count).Error; err != nil {
		return fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	if count == 0 {
		return &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	// Ambil entry_user_realisasi untuk notifikasi penolakan ke pengaju
	var header struct {
		EntryUserRealisasi string `gorm:"column:entry_user_realisasi"`
	}
	if err := r.db.Raw(queryGetKpiHeaderRealisasi, idPengajuan).Scan(&header).Error; err != nil {
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
		header.EntryUserRealisasi,
		"realisasi_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// =============================================================================
// GET APPROVAL LIST JSON
// =============================================================================

func (r *realisasiKpiRepo) GetApprovalListJSON(idPengajuan, userID string) (string, error) {
	var count int64
	if err := r.db.Raw(queryCheckApprovalRealisasi, userID, idPengajuan).Scan(&count).Error; err != nil {
		return "", fmt.Errorf("gagal mengecek data pengajuan: %w", err)
	}
	if count == 0 {
		return "", &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	var approvalListJSON string
	row := r.db.Raw(`SELECT approval_list_realisasi FROM data_kpi WHERE id_pengajuan = ? LIMIT 1`, idPengajuan).Row()
	if err := row.Scan(&approvalListJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", &customErrors.BadRequestError{
				Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
			}
		}
		return "", fmt.Errorf("gagal mengambil approval_list_realisasi: %w", err)
	}
	return approvalListJSON, nil
}

// =============================================================================
// GET ALL REALISASI — status 2
// =============================================================================

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
	countQuery := queryGetCountDataKpiRealisasi + where
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
	listQuery := queryGetDataKpiRealisasi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
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

// =============================================================================
// GET ALL APPROVAL — status 3, approval_posisi = user
// =============================================================================

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
	countQuery := queryGetCountDataKpiRealisasi + where
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
	listQuery := queryGetDataKpiRealisasi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
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

// =============================================================================
// GET ALL TOLAKAN — status 4, entry_user_realisasi = user
// =============================================================================

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
	countQuery := queryGetCountDataKpiRealisasi + where
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
	listQuery := queryGetDataKpiRealisasi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
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
// GET ALL DAFTAR REALISASI — semua status realisasi (2, 3, 4, 5, 80)
// =============================================================================

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
	countQuery := queryGetCountDataKpiRealisasi + where
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
	listQuery := queryGetDataKpiRealisasi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
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
// GET ALL DAFTAR APPROVAL — semua yang ada approval_list_realisasi
// =============================================================================

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
	countQuery := queryGetCountDataKpiRealisasi + where
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
	listQuery := queryGetDataKpiRealisasi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
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

func (r *realisasiKpiRepo) GetDetailRealisasiKpi(
	req *dto.GetDetailRealisasiKpiRequest,
) (*model.DataKpi, error) {

	// =========================================================================
	// QUERY HEADER
	// =========================================================================
	var result model.DataKpi
	if err := r.db.Raw(queryGetDetailHeader, req.IdPengajuan).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil detail realisasi KPI: %w", err)
	}
	if result.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	// =========================================================================
	// KPI DETAIL + SUB DETAIL
	// =========================================================================
	var kpiDetails []model.DataKpiDetail
	if err := r.db.Raw(queryGetDetailKpiList, result.IdPengajuan).Scan(&kpiDetails).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil kpi detail: %w", err)
	}
	for i := range kpiDetails {
		var subDetails []model.DataKpiSubDetail
		if err := r.db.Raw(queryGetDetailSubKpiList, result.IdPengajuan, kpiDetails[i].IdDetail).Scan(&subDetails).Error; err != nil {
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
	if err := r.db.Raw(queryGetDetailResultList, result.IdPengajuan).Scan(&resultList).Error; err != nil {
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
	if err := r.db.Raw(queryGetDetailProcessList, result.IdPengajuan).Scan(&processList).Error; err != nil {
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
	if err := r.db.Raw(queryGetDetailContextList, result.IdPengajuan).Scan(&contextList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil context detail: %w", err)
	}
	if contextList == nil {
		contextList = []model.DataChallengeDetail{}
	}
	result.ContextList = contextList
	result.TotalContext = len(contextList)

	return &result, nil
}
