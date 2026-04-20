package repo

import (
	"strings"

	dto "permen_api/domain/master_result/dto"
	model "permen_api/domain/master_result/model"
)

const (
	GetAllMasterResultBaseQuery = `
		SELECT id_result, nama_result, desc_result, tahun, triwulan, entry_user, entry_name, entry_time
		FROM mst_result`

	CheckTriwulanExistsQuery = `
		SELECT COUNT(1)
		FROM mst_triwulan
		WHERE id_triwulan = ?`
)

func (r *masterResultRepo) GetAllMasterResult(req *dto.GetAllMasterResultRequest) ([]*model.MstResult, error) {
	var results []*model.MstResult

	query := GetAllMasterResultBaseQuery

	var conditions []string
	var args []interface{}

	if req.Search != "" {
		conditions = append(conditions, "nama_result LIKE ?")
		args = append(args, "%"+req.Search+"%")
	}

	if req.Triwulan != "" {
		conditions = append(conditions, "triwulan = ?")
		args = append(args, req.Triwulan)
	}

	if req.Tahun != "" {
		conditions = append(conditions, "tahun = ?")
		args = append(args, req.Tahun)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY tahun DESC"

	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *masterResultRepo) CheckTriwulanExists(idTriwulan string) (bool, error) {
	var count int
	err := r.db.Raw(CheckTriwulanExistsQuery, idTriwulan).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
