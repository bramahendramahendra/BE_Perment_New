package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"
	customErrors "permen_api/errors"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Get Count
	// =============================================================================

	// Use func : GetAllValidasi, GetAllApprovalValidasiKpi, GetAllTolakanValidasiKpi, GetAllDaftarValidasiKpi, GetAllDaftarApprovalValidasiKpi
	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : CheckApprovalValidasiExists
	queryCheckApprovalValidasiExists = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 6 AND approval_posisi = ? AND id_pengajuan = ?`

	// =============================================================================
	// Get Data
	// =============================================================================

	// Use func : RejectValidasiKpi
	queryGetEntryUserValidasi = `
		SELECT IFNULL(entry_user_validasi, '')
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
		SELECT a.tahun, a.triwulan, a.kostl, a.kostl_tx, a.status, b.status_desc,
		       IFNULL(a.entry_user_validasi, '') entry_user_validasi,
		       IFNULL(a.entry_name_validasi, '') entry_name_validasi
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ? AND a.status IN (5, 6, 7, 8, 90, 91)
		LIMIT 1`

	// =============================================================================
	// Get Helper
	// =============================================================================

	// Use func : GetApprovalListValidasiJSON
	queryGetApprovalListJSON = `
		SELECT IFNULL(approval_list_validasi, '') FROM data_kpi
		WHERE status = 6 AND approval_posisi = ? AND id_pengajuan = ?`

	// Use func : InputValidasi
	queryGetKostlTxValidasi = `
		SELECT IFNULL(kostl_tx, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// =============================================================================
	// Get Data (untuk GetAll)
	// =============================================================================

	// Use func : GetAllValidasiKpi, GetAllApprovalValidasiKpi, GetAllTolakanValidasiKpi, GetAllDaftarValidasiKpi, GetAllDaftarApprovalValidasiKpi, dan GetDetailValidasiKpi
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
			IFNULL(a.total_pencapaian, '')         total_pencapaian,
			IFNULL(a.entry_user_validasi, '')       entry_user_validasi,
			IFNULL(a.entry_name_validasi, '')       entry_name_validasi,
			IFNULL(a.entry_time_validasi, '')       entry_time_validasi,
			IFNULL(a.approval_list_validasi, '')    approval_list_validasi,
			IFNULL(a.total_bobot_pengurang, '')     total_bobot_pengurang,
			IFNULL(a.total_pencapaian_post, '')     total_pencapaian_post,
			IFNULL(a.lampiran_validasi, '')         lampiran_validasi,
			IFNULL(a.qualifier_overall_validasi,'') qualifier_overall_validasi
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : GetDetailValidasiKpi
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

	// Use func : GetDetailValidasiKpi
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
			COALESCE(NULLIF(a.realisasi_kuantitatif_validated, ''), 0) realisasi_kuantitatif_validated,
			COALESCE(NULLIF(a.pencapaian, ''), 0)             pencapaian,
			COALESCE(NULLIF(a.skor, ''), 0)                   skor,
			IFNULL(a.realisasi_qualifier, '')             realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '') realisasi_kuantitatif_qualifier,
			IFNULL(a.validasi_keterangan, '')                             validasi_keterangan,
			COALESCE(NULLIF(a.pencapaian_qualifier_validated, ''), 0)     pencapaian_qualifier_validated,
			COALESCE(NULLIF(a.pencapaian_post_qualifier_validated, ''), 0) pencapaian_post_qualifier_validated
		FROM data_kpi_subdetail a
		LEFT JOIN mst_polarisasi p ON a.rumus = p.id_polarisasi
		LEFT JOIN mst_keterangan_project c ON a.id_keterangan_project = c.id
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	// Use func : GetDetailValidasiKpi
	queryGetDataResultDetail = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result,
			IFNULL(realisasi_result, '')   realisasi_result,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_result_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_result ASC`

	// Use func : GetDetailValidasiKpi
	queryGetDataProcessDetail = `
		SELECT
			id_detail_method, nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')  realisasi_method,
			IFNULL(lampiran_evidence, '') lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_method ASC`

	// Use func : GetDetailValidasiKpi
	queryGetDataContextDetail = `
		SELECT
			id_detail_challenge, nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '') realisasi_challenge,
			IFNULL(lampiran_evidence, '')   lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_challenge ASC`

	// =============================================================================
	// Update — InputValidasi
	// =============================================================================

	// Use func : InputValidasi
	queryUpdateKpiInputValidasi = `
		UPDATE data_kpi
		SET status                        = 6,
		    entry_user_validasi           = ?,
		    entry_name_validasi           = ?,
		    entry_time_validasi           = ?,
		    approval_posisi               = ?,
		    approval_list_validasi        = ?,
		    total_bobot                   = ?,
		    total_pencapaian              = ?,
		    lampiran_validasi             = ?,
		    total_bobot_pengurang         = ?,
		    total_pencapaian_post         = ?,
		    qualifier_overall_validasi    = ?
		WHERE id_pengajuan = ?`

	// Use func : InputValidasi
	queryUpdateSubDetailValidasi = `
		UPDATE data_kpi_subdetail
		SET target_triwulan                     = ?,
		    target_kuantitatif_triwulan         = ?,
		    realisasi_validated                 = ?,
		    realisasi_kuantitatif_validated     = ?,
		    pencapaian                          = ?,
		    skor                                = ?,
		    validasi_keterangan                 = ?,
		    pencapaian_qualifier_validated      = ?,
		    pencapaian_post_qualifier_validated = ?,
		    target_qualifier                    = ?
		WHERE id_pengajuan  = ?
		  AND id_detail     = ?
		  AND id_sub_detail = ?`

	// =============================================================================
	// Update Approval
	// =============================================================================

	// Use func : ApproveValidasiKpi
	queryApproveChainValidasi = `
		UPDATE data_kpi
		SET approval_posisi        = ?,
		    approval_list_validasi = ?
		WHERE id_pengajuan = ?`

	// Use func : ApproveValidasiKpi
	queryApproveFinalValidasi = `
		UPDATE data_kpi
		SET status                = 8,
		    approval_list_validasi = ?
		WHERE id_pengajuan = ?`

	// Use func : RejectValidasiKpi
	queryRejectValidasi = `
		UPDATE data_kpi
		SET status                = 7,
		    approval_list_validasi = ?,
		    catatan_tolakan       = ?
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Update — ValidasiBatal
	// =============================================================================

	// Use func : ValidasiBatal
	queryUpdateKpiBatalValidasi = `
		UPDATE data_kpi
		SET status              = 91,
		    entry_time_validasi = NOW()
		WHERE id_pengajuan = ?`

	// Use func : ValidasiBatal
	queryDeleteNotifBatalValidasi = `
		DELETE FROM log_notif
		WHERE key_notif     = ?
		  AND user_penerima = ?
		  AND status        = 0`
)

// =============================================================================
// INPUT VALIDASI
// =============================================================================

// InputValidasiKpi digunakan oleh endpoint POST /validasi-kpi/input.
func (r *validasiKpiRepo) InputValidasiKpi(req *dto.InputValidasiKpiRequest) error {
	// Ambil userid pertama dari ApprovalList sebagai approval_posisi
	approvalPosisi := ""
	if len(req.ApprovalListValidasi) > 0 {
		approvalPosisi = req.ApprovalListValidasi[0].Userid
	}

	approvalListBytes, err := json.Marshal(req.ApprovalListValidasi)
	if err != nil {
		return fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	lampiranJSON, err := json.Marshal(req.LampiranValidasi)
	if err != nil {
		return fmt.Errorf("gagal marshal lampiran_validasi: %w", err)
	}

	qualifierJSON, err := json.Marshal(req.DataValidasiQualifierOverall)
	if err != nil {
		return fmt.Errorf("gagal marshal data_validasi_qualifier_overall: %w", err)
	}

	// Ambil kostl_tx sebelum transaksi dimulai
	var kostlTxStr string
	if err := r.db.Raw(queryGetKostlTxValidasi, req.IdPengajuan).Scan(&kostlTxStr).Error; err != nil {
		return fmt.Errorf("gagal mengambil kostl_tx: %w", err)
	}

	// =========================================================================
	// Jalankan dalam transaksi agar update + notif atomic
	// =========================================================================
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryUpdateKpiInputValidasi,
		req.EntryUserValidasi,
		req.EntryNameValidasi,
		req.EntryTimeValidasi,
		approvalPosisi,
		approvalListBytes,
		req.TotalBobot,
		req.TotalPencapaian,
		lampiranJSON,
		req.TotalBobotPengurang,
		req.TotalPencapaianPost,
		qualifierJSON,
		req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update data_kpi input validasi: %w", err)
	}

	for _, kpi := range req.Kpi {
		for _, sub := range kpi.KpiSubDetail {
			if err := tx.Exec(queryUpdateSubDetailValidasi,
				sub.TargetTriwulan,
				sub.TargetKuantitatifTriwulan,
				sub.RealisasiValidated,
				sub.RealisasiKuantitatifValidated,
				sub.Pencapaian,
				sub.Skor,
				sub.ValidasiKeterangan,
				sub.PencapaianQualifierValidated,
				sub.PencapaianPostQualifierValidated,
				sub.TargetQualifier,
				req.IdPengajuan,
				kpi.IdDetail,
				sub.IdSubDetail,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal update sub detail '%s': %w", sub.IdSubDetail, err)
			}
		}
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Validasi "+kostlTxStr,
		req.EntryUserValidasi,
		approvalPosisi,
		"approval_validasi",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi input validasi: %w", err)
	}
	return nil
}

// =============================================================================
// APPROVAL
// =============================================================================

// ApproveValidasiKpi digunakan oleh endpoint POST /validasi-kpi/approve.
func (r *validasiKpiRepo) ApproveValidasiKpi(idPengajuan, approvalList, approvalPosisi, user string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if approvalPosisi == "" {
		// Approve final: status = 8
		if err := tx.Exec(queryApproveFinalValidasi,
			approvalList, idPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve final validasi: %w", err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("gagal commit approve final validasi: %w", err)
		}
		return nil
	}

	// Masih ada approver berikutnya → chain approval + notif
	if err := tx.Exec(queryApproveChainValidasi,
		approvalPosisi, approvalList, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal approve chain validasi: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Approval Validasi, ID : "+idPengajuan,
		user,
		approvalPosisi,
		"approval_validasi",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi approve validasi: %w", err)
	}
	return nil
}

// RejectValidasiKpi digunakan oleh endpoint POST /validasi-kpi/reject.
func (r *validasiKpiRepo) RejectValidasiKpi(idPengajuan, approvalList, catatan, user string) error {
	var entryUserValidasi string
	if err := r.db.Raw(queryGetEntryUserValidasi, idPengajuan).Scan(&entryUserValidasi).Error; err != nil {
		return fmt.Errorf("gagal mengambil entry_user_validasi: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryRejectValidasi,
		approvalList, catatan, idPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal reject validasi: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Validasi Ditolak, ID : "+idPengajuan,
		user,
		entryUserValidasi,
		"validasi_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi reject validasi: %w", err)
	}
	return nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all.
func (r *validasiKpiRepo) GetAllValidasiKpi(
	req *dto.GetAllValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"1=1",
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
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi, &h.TotalBobotPengurang, &h.TotalPencapaianPost, &h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
func (r *validasiKpiRepo) GetAllApprovalValidasiKpi(
	req *dto.GetAllApprovalValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"a.status = 6",
		"a.approval_posisi = ?",
	}
	args := []interface{}{req.ApprovalUserValidasi}

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
	if err := r.db.Raw(queryGetCountDataKpi+where, args...).Scan(&total).Error; err != nil {
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
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi, &h.TotalBobotPengurang, &h.TotalPencapaianPost, &h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllTolakanValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
func (r *validasiKpiRepo) GetAllTolakanValidasiKpi(
	req *dto.GetAllTolakanValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"a.status = 7",
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
	if err := r.db.Raw(queryGetCountDataKpi+where, args...).Scan(&total).Error; err != nil {
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
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi, &h.TotalBobotPengurang, &h.TotalPencapaianPost, &h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllDaftarValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-validasi.
func (r *validasiKpiRepo) GetAllDaftarValidasiKpi(
	req *dto.GetAllDaftarPValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"1=1",
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
	if err := r.db.Raw(queryGetCountDataKpi+where, args...).Scan(&total).Error; err != nil {
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
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi, &h.TotalBobotPengurang, &h.TotalPencapaianPost, &h.QualifierOverallValidasi,
		); err != nil {
			return nil, 0, fmt.Errorf("gagal scan header KPI: %w", err)
		}

		results = append(results, &h)
	}

	return results, total, nil
}

// GetAllDaftarApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
func (r *validasiKpiRepo) GetAllDaftarApprovalValidasiKpi(
	req *dto.GetAllDaftarApprovalValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"a.approval_list_validasi LIKE ?",
	}
	args := []interface{}{"%" + req.ApprovalUserValidasi + "%"}

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
	if err := r.db.Raw(queryGetCountDataKpi+where, args...).Scan(&total).Error; err != nil {
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
			&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
			&h.ApprovalListValidasi,
			&h.LampiranValidasi, &h.TotalBobotPengurang, &h.TotalPencapaianPost, &h.QualifierOverallValidasi,
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

// GetDetailValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-detail.
func (r *validasiKpiRepo) GetDetailValidasiKpi(
	req *dto.GetDetailValidasiKpiRequest,
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

// GetApprovalListJSON digunakan oleh service ApproveValidasiKpi dan RejectValidasiKpi untuk mengambil daftar approval dalam format JSON.
func (r *validasiKpiRepo) GetApprovalListJSON(idPengajuan, userID string) (string, error) {
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

// CheckApprovalValidasiExists digunakan oleh service ApproveValidasiKpi dan RejectValidasiKpi untuk memvalidasi keberadaan approval.
func (r *validasiKpiRepo) CheckApprovalValidasiExists(user, idPengajuan string) (bool, error) {
	var count int64
	err := r.db.Raw(queryCheckApprovalValidasiExists, user, idPengajuan).Scan(&count).Error
	return count > 0, err
}

// GetCatatanTolakan digunakan oleh service RejectValidasiKpi untuk mengambil catatan tolakan berdasarkan id_pengajuan.
func (r *validasiKpiRepo) GetCatatanTolakan(idPengajuan string) (string, error) {
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

// GetExistDataKpi mengambil header KPI untuk keperluan validasi status sebelum operasi.
func (r *validasiKpiRepo) GetExistDataKpi(idPengajuan string) (*model.DataKpiExist, error) {
	var exist model.DataKpiExist
	db := r.db.Raw(queryGetExistDataKpi, idPengajuan).Scan(&exist)
	if db.Error != nil {
		return nil, fmt.Errorf("gagal mengambil header KPI: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return nil, fmt.Errorf("id_pengajuan '%s' tidak ditemukan", idPengajuan)
	}
	return &exist, nil
}
