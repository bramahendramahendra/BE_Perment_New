package repo

import (
	"fmt"
	"strings"

	dto "permen_api/domain/pencapaian_kpi/dto"
	model "permen_api/domain/pencapaian_kpi/model"
)

const (
	// =============================================================================
	// Get Count
	// =============================================================================

	queryGetCountDataKpi = `
		SELECT COUNT(1)
		FROM data_kpi a
		INNER JOIN mst_status b ON a.status = b.id_status`

	// =============================================================================
	// Get Data (untuk GetAll dan GetDetail)
	// =============================================================================

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
			IFNULL(a.id_sumber, '')                           id_sumber,
			IFNULL(s.sumber, '')                              sumber,
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
		LEFT JOIN mst_sumber s ON a.id_sumber = s.id_sumber
		WHERE a.id_pengajuan = ? AND a.id_detail = ?
		ORDER BY a.id_sub_detail ASC`

	queryGetDataResultDetail = `
		SELECT
			id_detail_result,
			nama_result, deskripsi_result,
			IFNULL(realisasi_result, '')   realisasi_result,
			IFNULL(lampiran_evidence, '')  lampiran_evidence
		FROM data_result_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_result ASC`

	queryGetDataProcessDetail = `
		SELECT
			id_detail_method, nama_method, deskripsi_method,
			IFNULL(realisasi_method, '')  realisasi_method,
			IFNULL(lampiran_evidence, '') lampiran_evidence
		FROM data_method_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_method ASC`

	queryGetDataContextDetail = `
		SELECT
			id_detail_challenge, nama_challenge, deskripsi_challenge,
			IFNULL(realisasi_challenge, '') realisasi_challenge,
			IFNULL(lampiran_evidence, '')   lampiran_evidence
		FROM data_challenge_detail
		WHERE id_pengajuan = ?
		ORDER BY id_detail_challenge ASC`

	queryGetIndikatorPencapaian = `
		SELECT indikator_warna, indikator_value
		FROM indikator_pencapaian
		ORDER BY indikator_value DESC`
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-all-pencapaian.
func (r *pencapaianKpiRepo) GetAllPencapaianKpi(
	req *dto.GetAllPencapaianKpiRequest,
) ([]*model.DataKpi, int64, error) {
	conditions := []string{
		"a.status = 8",
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

	where := " WHERE " + strings.Join(conditions, " AND ")

	var total int64
	if err := r.db.Raw(queryGetCountDataKpi+where, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("gagal menghitung total data: %w", err)
	}

	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	listQuery := queryGetDataKpi + where + " ORDER BY a.tahun DESC, a.triwulan DESC LIMIT ? OFFSET ?"
	listArgs := append(args, limit, offset)

	rows, err := r.db.Raw(listQuery, listArgs...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("gagal mengambil daftar pencapaian KPI: %w", err)
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
			return nil, 0, fmt.Errorf("gagal scan header pencapaian KPI: %w", err)
		}
		results = append(results, &h)
	}

	return results, total, nil
}

// =============================================================================
// GET DETAIL
// =============================================================================

// GetDetailPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-detail.
func (r *pencapaianKpiRepo) GetDetailPencapaianKpi(
	req *dto.GetDetailPencapaianKpiRequest,
) (*model.DataKpi, error) {
	where := " WHERE a.id_pengajuan = ?"
	args := []interface{}{req.IdPengajuan}

	var result model.DataKpi
	headerQuery := queryGetDataKpi + where + " LIMIT 1"
	if err := r.db.Raw(headerQuery, args...).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil detail pencapaian KPI: %w", err)
	}

	var kpiDetails []model.DataKpiDetail
	if err := r.db.Raw(queryGetDataKpiDetail, result.IdPengajuan).Scan(&kpiDetails).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil kpi detail: %w", err)
	}

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

	var resultList []model.DataResultDetail
	if err := r.db.Raw(queryGetDataResultDetail, result.IdPengajuan).Scan(&resultList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil result detail: %w", err)
	}
	if resultList == nil {
		resultList = []model.DataResultDetail{}
	}
	result.ResultList = resultList
	result.TotalResult = len(resultList)

	var processList []model.DataMethodDetail
	if err := r.db.Raw(queryGetDataProcessDetail, result.IdPengajuan).Scan(&processList).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil process detail: %w", err)
	}
	if processList == nil {
		processList = []model.DataMethodDetail{}
	}
	result.ProcessList = processList
	result.TotalProcess = len(processList)

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
// GET INDIKATOR
// =============================================================================

// GetIndikatorPencapaian mengambil semua indikator warna dari tabel indikator_pencapaian, diurutkan descending.
func (r *pencapaianKpiRepo) GetIndikatorPencapaian() ([]*model.IndikatorPencapaian, error) {
	var result []*model.IndikatorPencapaian
	if err := r.db.Raw(queryGetIndikatorPencapaian).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil indikator pencapaian: %w", err)
	}
	return result, nil
}
