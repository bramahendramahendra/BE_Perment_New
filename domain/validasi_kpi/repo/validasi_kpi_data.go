package repo

import (
	"encoding/json"
	"fmt"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Get Count
	// =============================================================================

	// Use func : GetAllApprovalValidasi, GetAllTolakanValidasi, GetAllDaftarPenyusunanValidasi,
	//            GetAllDaftarApprovalValidasi, GetAllValidasi
	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// Use func : CheckApprovalValidasiExists
	queryCheckApprovalValidasiExists = `
		SELECT COUNT(*) FROM data_kpi
		WHERE status = 6 AND approval_posisi = ? AND id_pengajuan = ?`

	// =============================================================================
	// Get Exist
	// =============================================================================

	// Use func : GetExistDataValidasi
	queryGetExistDataValidasi = `
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

	// Use func : GetApprovalListValidasiJSON, CheckApprovalValidasiExists
	queryGetApprovalListValidasiJSON = `
		SELECT IFNULL(approval_list_validasi, '') FROM data_kpi
		WHERE status = 6 AND approval_posisi = ? AND id_pengajuan = ?`

	// Use func : ValidasiBatal
	queryGetApprovalPosisiValidasi = `
		SELECT IFNULL(approval_posisi, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// Use func : RejectValidasiKpi
	queryGetEntryUserValidasi = `
		SELECT IFNULL(entry_user_validasi, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// Use func : InputValidasi
	queryGetKostlTxValidasi = `
		SELECT IFNULL(kostl_tx, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// =============================================================================
	// Get Data (untuk GetAll)
	// =============================================================================

	// Use func : GetAllApprovalValidasi, GetAllTolakanValidasi, GetAllDaftarPenyusunanValidasi,
	//            GetAllDaftarApprovalValidasi, GetAllValidasi
	queryGetDataKpiValidasi = `
		SELECT
			a.id_pengajuan, a.tahun, a.triwulan, a.kostl, a.kostl_tx,
			a.orgeh, a.orgeh_tx, a.status, b.status_desc
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// =============================================================================
	// Get Data (untuk GetDetail)
	// =============================================================================

	// Use func : GetDetailValidasiKpi — header data_kpi
	queryGetDetailDataKpiValidasi = `
		SELECT
			a.id_pengajuan, a.tahun, a.triwulan, a.kostl, a.kostl_tx,
			a.orgeh, a.orgeh_tx, a.entry_user, a.entry_name, a.entry_time,
			a.approval_posisi, a.status, b.status_desc,
			IFNULL(a.entry_user_realisasi, '')       entry_user_realisasi,
			IFNULL(a.entry_name_realisasi, '')       entry_name_realisasi,
			IFNULL(a.entry_time_realisasi, '')       entry_time_realisasi,
			IFNULL(a.catatan_tolakan, '')             catatan_tolakan,
			IFNULL(a.total_bobot, '')                 total_bobot,
			IFNULL(a.total_pencapaian, '')            total_pencapaian,
			IFNULL(a.total_bobot_pengurang, '')       total_bobot_pengurang,
			IFNULL(a.total_pencapaian_post, '')       total_pencapaian_post,
			IFNULL(a.entry_user_validasi, '')         entry_user_validasi,
			IFNULL(a.entry_name_validasi, '')         entry_name_validasi,
			IFNULL(a.entry_time_validasi, '')         entry_time_validasi,
			IFNULL(a.approval_list_validasi, '')      approval_list_validasi,
			IFNULL(a.lampiran_validasi, '')           lampiran_validasi,
			IFNULL(a.qualifier_overall_validasi, '')  qualifier_overall_validasi
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status
		WHERE a.id_pengajuan = ?
		LIMIT 1`

	// Use func : GetDetailValidasiKpi — daftar KPI detail
	queryGetDetailKpiDetail = `
		SELECT
			a.id_detail,
			a.id_kpi, a.kpi, a.rumus
		FROM data_kpi_detail a
		WHERE a.id_pengajuan = ?
		ORDER BY a.id_detail ASC`

	// Use func : GetDetailValidasiKpi — sub detail per KPI
	queryGetDetailKpiSubDetail = `
		SELECT
			a.id_sub_detail,
			a.id_kpi, a.kpi,
			COALESCE(NULLIF(a.bobot, ''), 0)                              bobot,
			IFNULL(a.target_triwulan, '')                                 target_triwulan,
			COALESCE(NULLIF(a.target_kuantitatif_triwulan, ''), 0)        target_kuantitatif_triwulan,
			IFNULL(a.target_qualifier, '')                                target_qualifier,
			IFNULL(a.realisasi_validated, '')                             realisasi_validated,
			COALESCE(NULLIF(a.realisasi_kuantitatif_validated, ''), 0)    realisasi_kuantitatif_validated,
			IFNULL(a.validasi_keterangan, '')                             validasi_keterangan,
			COALESCE(NULLIF(a.pencapaian, ''), 0)                         pencapaian,
			COALESCE(NULLIF(a.skor, ''), 0)                               skor,
			COALESCE(NULLIF(a.pencapaian_qualifier_validated, ''), 0)     pencapaian_qualifier_validated,
			COALESCE(NULLIF(a.pencapaian_post_qualifier_validated, ''), 0) pencapaian_post_qualifier_validated
		FROM data_kpi_subdetail a
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	// =============================================================================
	// Update — InputValidasi
	// =============================================================================

	// Use func : InputValidasi
	queryUpdateKpiInputValidasi = `
		UPDATE data_kpi
		SET status                        = 6,
		    entry_user_validasi           = ?,
		    entry_name_validasi           = ?,
		    entry_time_validasi           = NOW(),
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
	// Update — Approve / Reject
	// =============================================================================

	// Use func : ApproveValidasiKpi (final — status = 8)
	queryApproveFinalValidasi = `
		UPDATE data_kpi
		SET status                = 8,
		    approval_list_validasi = ?
		WHERE id_pengajuan = ?`

	// Use func : ApproveValidasiKpi (chain — lanjut ke approver berikutnya)
	queryApproveChainValidasi = `
		UPDATE data_kpi
		SET approval_posisi        = ?,
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
	var kostlTxStr string
	if err := r.db.Raw(queryGetKostlTxValidasi, req.IdPengajuan).Scan(&kostlTxStr).Error; err != nil {
		return fmt.Errorf("gagal mengambil kostl_tx: %w", err)
	}

	approvalListJSON, err := json.Marshal(req.ApprovalListValidasi)
	if err != nil {
		return fmt.Errorf("gagal marshal approval_list_validasi: %w", err)
	}

	lampiranJSON, err := json.Marshal(req.LampiranValidasi)
	if err != nil {
		return fmt.Errorf("gagal marshal lampiran_validasi: %w", err)
	}

	qualifierJSON, err := json.Marshal(req.DataValidasiQualifierOverall)
	if err != nil {
		return fmt.Errorf("gagal marshal data_validasi_qualifier_overall: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryUpdateKpiInputValidasi,
		req.EntryUserValidasi,
		req.EntryNameValidasi,
		req.ApprovalPosisi,
		string(approvalListJSON),
		req.TotalBobot,
		req.TotalPencapaian,
		string(lampiranJSON),
		req.TotalBobotPengurang,
		req.TotalPencapaianPost,
		string(qualifierJSON),
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
		req.ApprovalPosisi,
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
func (r *validasiKpiRepo) ApproveValidasiKpi(idPengajuan, approvalList, nextApprover, user string) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if nextApprover == "" {
		// Approve final: tidak ada approver berikutnya → status = 8
		if err := tx.Exec(queryApproveFinalValidasi, approvalList, idPengajuan).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve final validasi: %w", err)
		}
		return tx.Commit().Error
	}

	// Chain approval: masih ada approver berikutnya → update posisi + kirim notif
	if err := tx.Exec(queryApproveChainValidasi, nextApprover, approvalList, idPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal approve chain validasi: %w", err)
	}

	if err := notif.Insert(tx,
		idPengajuan,
		"Approval Validasi, ID : "+idPengajuan,
		user,
		nextApprover,
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

	if err := tx.Exec(queryRejectValidasi, approvalList, catatan, idPengajuan).Error; err != nil {
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
// GET DETAIL
// =============================================================================

// GetDetailValidasiKpi mengambil detail lengkap satu pengajuan validasi KPI.
func (r *validasiKpiRepo) GetDetailValidasiKpi(req *dto.GetDetailValidasiKpiRequest) (*model.DataKpi, error) {
	var h model.DataKpi

	// -------------------------------------------------------------------------
	// Query header data_kpi
	// -------------------------------------------------------------------------
	headerRow := r.db.Raw(queryGetDetailDataKpiValidasi, req.IdPengajuan).Row()
	if err := headerRow.Scan(
		&h.IdPengajuan, &h.Tahun, &h.Triwulan,
		&h.Kostl, &h.KostlTx,
		&h.Orgeh, &h.OrgehTx,
		&h.EntryUser, &h.EntryName, &h.EntryTime,
		&h.ApprovalPosisi,
		&h.Status, &h.StatusDesc,
		&h.EntryUserRealisasi, &h.EntryNameRealisasi, &h.EntryTimeRealisasi,
		&h.CatatanTolakan,
		&h.TotalBobot, &h.TotalPencapaian,
		&h.TotalBobotPengurang, &h.TotalPencapaianPost,
		&h.EntryUserValidasi, &h.EntryNameValidasi, &h.EntryTimeValidasi,
		&h.ApprovalListValidasi,
		&h.LampiranValidasi,
		&h.QualifierOverallValidasi,
	); err != nil {
		return nil, fmt.Errorf("gagal scan header validasi KPI: %w", err)
	}

	if h.IdPengajuan == "" {
		return &h, nil
	}

	// -------------------------------------------------------------------------
	// Query data_kpi_detail
	// -------------------------------------------------------------------------
	detailRows, err := r.db.Raw(queryGetDetailKpiDetail, req.IdPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil KPI detail: %w", err)
	}
	defer detailRows.Close()

	var kpiDetails []model.DataKpiDetailValidasi
	for detailRows.Next() {
		var d model.DataKpiDetailValidasi
		if err := detailRows.Scan(&d.IdDetail, &d.IdKpi, &d.Kpi, &d.Rumus); err != nil {
			return nil, fmt.Errorf("gagal scan KPI detail: %w", err)
		}
		kpiDetails = append(kpiDetails, d)
	}

	// -------------------------------------------------------------------------
	// Query data_kpi_subdetail per detail
	// -------------------------------------------------------------------------
	for i, kpi := range kpiDetails {
		subRows, err := r.db.Raw(queryGetDetailKpiSubDetail, req.IdPengajuan, kpi.IdDetail).Rows()
		if err != nil {
			return nil, fmt.Errorf("gagal mengambil sub detail KPI '%s': %w", kpi.IdDetail, err)
		}

		var subDetails []model.DataKpiSubDetailValidasi
		for subRows.Next() {
			var s model.DataKpiSubDetailValidasi
			if err := subRows.Scan(
				&s.IdSubDetail, &s.IdKpi, &s.Kpi, &s.Bobot,
				&s.TargetTriwulan, &s.TargetKuantitatifTriwulan,
				&s.TargetQualifier,
				&s.RealisasiValidated, &s.RealisasiKuantitatifValidated,
				&s.ValidasiKeterangan,
				&s.Pencapaian, &s.Skor,
				&s.PencapaianQualifierValidated, &s.PencapaianPostQualifierValidated,
			); err != nil {
				subRows.Close()
				return nil, fmt.Errorf("gagal scan sub detail '%s': %w", kpi.IdDetail, err)
			}
			subDetails = append(subDetails, s)
		}
		subRows.Close()

		kpiDetails[i].KpiSubDetail = subDetails
		kpiDetails[i].TotalSubKpi = len(subDetails)
	}

	h.Kpi = kpiDetails
	h.TotalKpi = len(kpiDetails)

	return &h, nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all.
func (r *validasiKpiRepo) GetAllValidasiKPi(
	req *dto.GetAllValidasiKpiRequest,
) ([]*model.DataKpi, int64, error) {
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
	listArgs := append(args, limit, offset)
	rows, err := r.db.Raw(queryGetDataKpiValidasi+where+" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?", listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	results, err := r.scanGetAllRows(rows)
	return results, total, err
}

// =============================================================================
// GET ALL APPROVAL VALIDASI — status=6, approval_posisi=user
// =============================================================================

func (r *validasiKpiRepo) GetAllApprovalValidasi(
	req *dto.GetAllApprovalValidasiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{"a.status = 6", "a.approval_posisi = ?"}
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

	where := " WHERE " + strings.Join(conditions, " AND ")

	var total int64
	if err := r.db.Raw(queryGetCountDataKpiValidasi+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	_, limit, offset := paginateValidasi(req.Page, req.Limit)
	listArgs := append(args, limit, offset)
	rows, err := r.db.Raw(queryGetDataKpiValidasi+where+" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?", listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	results, err := r.scanGetAllRows(rows)
	return results, total, err
}

// =============================================================================
// GET ALL TOLAKAN VALIDASI — status=7
// =============================================================================

func (r *validasiKpiRepo) GetAllTolakanValidasi(
	req *dto.GetAllTolakanValidasiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{"a.status = 7"}
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

	where := " WHERE " + strings.Join(conditions, " AND ")

	var total int64
	if err := r.db.Raw(queryGetCountDataKpiValidasi+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	_, limit, offset := paginateValidasi(req.Page, req.Limit)
	listArgs := append(args, limit, offset)
	rows, err := r.db.Raw(queryGetDataKpiValidasi+where+" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?", listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	results, err := r.scanGetAllRows(rows)
	return results, total, err
}

// =============================================================================
// GET ALL DAFTAR PENYUSUNAN VALIDASI — semua data dengan filter opsional
// =============================================================================

func (r *validasiKpiRepo) GetAllDaftarPenyusunanValidasi(
	req *dto.GetAllDaftarPenyusunanValidasiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{"1=1"}
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

	where := " WHERE " + strings.Join(conditions, " AND ")

	var total int64
	if err := r.db.Raw(queryGetCountDataKpiValidasi+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	_, limit, offset := paginateValidasi(req.Page, req.Limit)
	listArgs := append(args, limit, offset)
	rows, err := r.db.Raw(queryGetDataKpiValidasi+where+" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?", listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	results, err := r.scanGetAllRows(rows)
	return results, total, err
}

// =============================================================================
// GET ALL DAFTAR APPROVAL VALIDASI — approval_list_validasi LIKE %user%
// =============================================================================

func (r *validasiKpiRepo) GetAllDaftarApprovalValidasi(
	req *dto.GetAllDaftarApprovalValidasiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{"a.approval_list_validasi LIKE ?"}
	args := []interface{}{"%" + req.ApprovalUser + "%"}

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

	var total int64
	if err := r.db.Raw(queryGetCountDataKpiValidasi+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	_, limit, offset := paginateValidasi(req.Page, req.Limit)
	listArgs := append(args, limit, offset)
	rows, err := r.db.Raw(queryGetDataKpiValidasi+where+" ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?", listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar KPI: %w", err)
	}
	results, err := r.scanGetAllRows(rows)
	return results, total, err
}

// =============================================================================
// GET EXIST
// =============================================================================

// GetExistDataValidasi mengambil header KPI untuk keperluan validasi status sebelum operasi.
func (r *validasiKpiRepo) GetExistDataValidasi(idPengajuan string) (*model.DataKpiExist, error) {
	var exist model.DataKpiExist
	row := r.db.Raw(queryGetExistDataValidasi, idPengajuan).Row()
	if err := row.Scan(
		&exist.Tahun, &exist.Triwulan, &exist.Kostl, &exist.KostlTx,
		&exist.Status, &exist.StatusDesc,
		&exist.EntryUserValidasi, &exist.EntryNameValidasi,
	); err != nil {
		return nil, fmt.Errorf("gagal mengambil data pengajuan '%s': %w", idPengajuan, err)
	}
	return &exist, nil
}

// =============================================================================
// APPROVAL HELPER
// =============================================================================

// CheckApprovalValidasiExists memvalidasi bahwa user adalah approver aktif.
func (r *validasiKpiRepo) CheckApprovalValidasiExists(user, idPengajuan string) (bool, error) {
	var count int64
	err := r.db.Raw(queryCheckApprovalValidasiExists, user, idPengajuan).Scan(&count).Error
	return count > 0, err
}

// GetApprovalListValidasiJSON mengambil approval_list_validasi dalam format JSON string.
func (r *validasiKpiRepo) GetApprovalListValidasiJSON(idPengajuan, userID string) (string, error) {
	var result string
	err := r.db.Raw(queryGetApprovalListValidasiJSON, userID, idPengajuan).Scan(&result).Error
	return result, err
}
