package repo

import (
	"strings"

	dto "permen_api/domain/master_method/dto"
	model "permen_api/domain/master_method/model"
)

const (
	GetAllMasterMethodBaseQuery = `
		SELECT id_method_use, nama_method, desc_method, tahun, triwulan, entry_user, entry_name, entry_time
		FROM mst_method_use`

	CheckTriwulanExistsQuery = `
		SELECT COUNT(1)
		FROM mst_triwulan
		WHERE id_triwulan = ?`
)

func (r *masterMethodRepo) GetAllMasterMethod(req *dto.GetAllMasterMethodRequest) ([]*model.MstMethod, error) {
	var methods []*model.MstMethod

	query := GetAllMasterMethodBaseQuery

	var conditions []string
	var args []interface{}

	if req.Search != "" {
		conditions = append(conditions, "nama_method LIKE ?")
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

	err := r.db.Raw(query, args...).Scan(&methods).Error
	if err != nil {
		return nil, err
	}

	return methods, nil
}

func (r *masterMethodRepo) CheckTriwulanExists(idTriwulan string) (bool, error) {
	var count int
	err := r.db.Raw(CheckTriwulanExistsQuery, idTriwulan).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
