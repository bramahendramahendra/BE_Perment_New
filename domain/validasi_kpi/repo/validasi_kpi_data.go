package repo

import (
	"fmt"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"
	notif "permen_api/pkg/notif"
)

const (
	// =============================================================================
	// Check
	// =============================================================================

	// queryCheckExistInputValidasi memeriksa status yang mengizinkan input validasi (5, 7, 90, 91).
	queryCheckExistInputValidasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status IN (5, 7, 90, 91)`

	// queryCheckExistApprovalValidasi memeriksa status 6 (pending approval validasi).
	queryCheckExistApprovalValidasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ? AND status = 6`

	// queryCheckExistBatalValidasi memeriksa keberadaan id_pengajuan.
	queryCheckExistBatalValidasi = `
		SELECT COUNT(id_pengajuan)
		FROM data_kpi
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Lookup
	// =============================================================================

	queryGetKostlTx = `
		SELECT IFNULL(kostl_tx, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	queryGetEntryUserValidasi = `
		SELECT IFNULL(entry_user_validasi, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	queryGetApprovalPosisiValidasi = `
		SELECT IFNULL(approval_posisi, '')
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// =============================================================================
	// Update — InputValidasi
	// =============================================================================

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

	queryUpdateSubDetailValidasi = `
		UPDATE data_kpi_subdetail
		SET target_triwulan                      = ?,
		    target_kuantitatif_triwulan          = ?,
		    realisasi_validated                  = ?,
		    realisasi_kuantitatif_validated      = ?,
		    pencapaian                           = ?,
		    skor                                 = ?,
		    validasi_keterangan                  = ?,
		    pencapaian_qualifier_validated       = ?,
		    pencapaian_post_qualifier_validated  = ?,
		    target_qualifier                     = ?
		WHERE id_pengajuan = ?
		  AND id_detail    = ?
		  AND id_sub_detail = ?`

	// =============================================================================
	// Update — ApprovalValidasi
	// =============================================================================

	// Approve final: approval_posisi kosong → set status = 8
	queryApproveFinalValidasi = `
		UPDATE data_kpi
		SET status                = 8,
		    approval_list_validasi = ?
		WHERE id_pengajuan = ?`

	// Approve chain: masih ada approver berikutnya
	queryApproveChainValidasi = `
		UPDATE data_kpi
		SET approval_posisi       = ?,
		    approval_list_validasi = ?
		WHERE id_pengajuan = ?`

	// Reject: set status = 7
	queryRejectValidasi = `
		UPDATE data_kpi
		SET status                = 7,
		    approval_list_validasi = ?,
		    catatan_tolakan       = ?
		WHERE id_pengajuan = ?`

	// =============================================================================
	// Update — ValidasiBatal
	// =============================================================================

	queryUpdateKpiBatalValidasi = `
		UPDATE data_kpi
		SET status             = 91,
		    entry_time_validasi = NOW()
		WHERE id_pengajuan = ?`

	queryDeleteNotifBatalValidasi = `
		DELETE FROM log_notif
		WHERE key_notif      = ?
		  AND user_penerima  = ?
		  AND status         = 0`
)

// =============================================================================
// CHECK
// =============================================================================

func (r *validasiKpiRepo) CheckExistInputValidasi(idPengajuan string) (bool, error) {
	var count int64
	err := r.db.Raw(queryCheckExistInputValidasi, idPengajuan).Scan(&count).Error
	return count > 0, err
}

func (r *validasiKpiRepo) CheckExistApprovalValidasi(idPengajuan string) (bool, error) {
	var count int64
	err := r.db.Raw(queryCheckExistApprovalValidasi, idPengajuan).Scan(&count).Error
	return count > 0, err
}

func (r *validasiKpiRepo) CheckExistBatalValidasi(idPengajuan string) (bool, error) {
	var count int64
	err := r.db.Raw(queryCheckExistBatalValidasi, idPengajuan).Scan(&count).Error
	return count > 0, err
}

// =============================================================================
// LOOKUP
// =============================================================================

func (r *validasiKpiRepo) GetKostlTxByIdPengajuan(idPengajuan string) (string, error) {
	var kostlTx string
	err := r.db.Raw(queryGetKostlTx, idPengajuan).Scan(&kostlTx).Error
	return kostlTx, err
}

func (r *validasiKpiRepo) GetEntryUserValidasiByIdPengajuan(idPengajuan string) (string, error) {
	var entryUser string
	err := r.db.Raw(queryGetEntryUserValidasi, idPengajuan).Scan(&entryUser).Error
	return entryUser, err
}

// =============================================================================
// INPUT VALIDASI
// =============================================================================

func (r *validasiKpiRepo) InputValidasi(req *dto.InputValidasiRequest) error {
	kostlTx, err := r.GetKostlTxByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return fmt.Errorf("gagal mengambil kostl_tx: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryUpdateKpiInputValidasi,
		req.EntryUserValidasi,
		req.EntryNameValidasi,
		req.ApprovalPosisi,
		req.ApprovalListValidasi,
		req.TotalBobot,
		req.TotalPencapaian,
		req.LampiranValidasi,
		req.TotalBobotPengurang,
		req.TotalPencapaianPost,
		req.DataValidasiQualifierOverall,
		req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update data_kpi input validasi: %w", err)
	}

	for _, detail := range req.DataValidasi {
		if err := tx.Exec(queryUpdateSubDetailValidasi,
			detail.TargetTriwulanValidated,
			detail.TargetKuantitatifValidated,
			detail.RealisasiValidated,
			detail.RealisasiKuantitatifValidated,
			detail.Pencapaian,
			detail.Skor,
			detail.ValidasiKeterangan,
			detail.PencapaianQualifierValidated,
			detail.PencapaianPostQualifierValidated,
			detail.TargetQualifierValidated,
			detail.KeyPengajuan,
			detail.KeyDetail,
			detail.KeySubDetail,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal update sub detail '%s': %w", detail.KeySubDetail, err)
		}
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Validasi "+kostlTx,
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
// APPROVAL VALIDASI
// =============================================================================

func (r *validasiKpiRepo) ApprovalValidasi(req *dto.ApprovalValidasiRequest) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if req.Status == "approve" {
		if req.ApprovalPosisi == "" {
			// Approve final: tidak ada approver berikutnya → status = 8
			if err := tx.Exec(queryApproveFinalValidasi,
				req.ApprovalList, req.IdPengajuan,
			).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal approve final validasi: %w", err)
			}
			tx.Commit()
			return nil
		}

		// Approve chain: masih ada approver berikutnya + kirim notif
		if err := tx.Exec(queryApproveChainValidasi,
			req.ApprovalPosisi, req.ApprovalList, req.IdPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve chain validasi: %w", err)
		}

		if err := notif.Insert(tx,
			req.IdPengajuan,
			"Approval Validasi, ID : "+req.IdPengajuan,
			req.ApprovalUser,
			req.ApprovalPosisi,
			"approval_validasi",
		); err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}

	// Reject: status = 7, kirim notif ke pengaju
	entryUserValidasi, err := r.GetEntryUserValidasiByIdPengajuan(req.IdPengajuan)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mengambil entry_user_validasi: %w", err)
	}

	if err := tx.Exec(queryRejectValidasi,
		req.ApprovalList, req.CatatanTolakan, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal reject validasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Validasi Ditolak, ID : "+req.IdPengajuan,
		req.ApprovalUser,
		entryUserValidasi,
		"validasi_ditolak",
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi approval validasi: %w", err)
	}
	return nil
}

// =============================================================================
// VALIDASI BATAL
// =============================================================================

func (r *validasiKpiRepo) ValidasiBatal(req *dto.ValidasiBatalRequest) error {
	// Ambil approval_posisi sebelum update agar bisa hapus notif yang sudah dikirim
	var approvalPosisi string
	if err := r.db.Raw(queryGetApprovalPosisiValidasi, req.IdPengajuan).Scan(&approvalPosisi).Error; err != nil {
		return fmt.Errorf("gagal mengambil approval_posisi: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryUpdateKpiBatalValidasi, req.IdPengajuan).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal update status batal validasi: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("gagal commit transaksi batal validasi: %w", err)
	}

	// Hapus notifikasi yang belum dibaca (dilakukan di luar transaksi utama, non-critical)
	if approvalPosisi != "" {
		r.db.Exec(queryDeleteNotifBatalValidasi, req.IdPengajuan, approvalPosisi)
	}

	return nil
}

// =============================================================================
// APPROVE VALIDASI (terpisah dari reject)
// =============================================================================

func (r *validasiKpiRepo) ApproveValidasi(req *dto.ApproveValidasiRequest) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if req.ApprovalPosisi == "" {
		// Approve final: tidak ada approver berikutnya → status = 8
		if err := tx.Exec(queryApproveFinalValidasi,
			req.ApprovalList, req.IdPengajuan,
		).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal approve final validasi: %w", err)
		}
		return tx.Commit().Error
	}

	// Approve chain: masih ada approver berikutnya + kirim notif
	if err := tx.Exec(queryApproveChainValidasi,
		req.ApprovalPosisi, req.ApprovalList, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal approve chain validasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Approval Validasi, ID : "+req.IdPengajuan,
		req.ApprovalUser,
		req.ApprovalPosisi,
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

// =============================================================================
// REJECT VALIDASI (terpisah dari approve)
// =============================================================================

func (r *validasiKpiRepo) RejectValidasi(req *dto.RejectValidasiRequest) error {
	entryUserValidasi, err := r.GetEntryUserValidasiByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return fmt.Errorf("gagal mengambil entry_user_validasi: %w", err)
	}

	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", tx.Error)
	}

	if err := tx.Exec(queryRejectValidasi,
		req.ApprovalList, req.CatatanTolakan, req.IdPengajuan,
	).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal reject validasi: %w", err)
	}

	if err := notif.Insert(tx,
		req.IdPengajuan,
		"Validasi Ditolak, ID : "+req.IdPengajuan,
		req.ApprovalUser,
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
// GET ALL — shared query constants (validasi)
// =============================================================================

const (
	queryGetCountDataKpiValidasi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	queryGetDataKpiValidasi = `
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
)

func (r *validasiKpiRepo) scanDataKpiRows(rows interface{ Next() bool; Scan(...interface{}) error; Close() error }) ([]*model.DataKpi, error) {
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
			return nil, fmt.Errorf("gagal scan header KPI: %w", err)
		}
		results = append(results, &h)
	}
	return results, nil
}

func paginateValidasi(page, limit int) (int, int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return page, limit, (page - 1) * limit
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
	results, err := r.scanDataKpiRows(rows)
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
	results, err := r.scanDataKpiRows(rows)
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
	results, err := r.scanDataKpiRows(rows)
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
	results, err := r.scanDataKpiRows(rows)
	return results, total, err
}

// =============================================================================
// GET ALL VALIDASI — semua data tanpa filter user (mirip get-all realisasi)
// =============================================================================

func (r *validasiKpiRepo) GetAllValidasi(
	req *dto.GetAllValidasiRequest,
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
	results, err := r.scanDataKpiRows(rows)
	return results, total, err
}
