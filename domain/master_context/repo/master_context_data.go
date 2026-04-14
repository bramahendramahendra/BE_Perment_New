package repo

import (
	"strings"

	dto "permen_api/domain/master_context/dto"
	model "permen_api/domain/master_context/model"
)

const (
	GetAllMasterContextBaseQuery = `
		SELECT id_challenge, nama_challenge, desc_challenge, tahun, triwulan, entry_user, entry_name, entry_time
		FROM mst_challenge`

	CheckTriwulanExistsQuery = `
		SELECT COUNT(1)
		FROM mst_triwulan
		WHERE id_triwulan = ?`
)

func (r *masterContextRepo) GetAllMasterContext(req *dto.GetAllMasterContextRequest) ([]*model.MstChallenge, error) {
	var contexts []*model.MstChallenge

	query := GetAllMasterContextBaseQuery

	var conditions []string
	var args []interface{}

	if req.Search != "" {
		conditions = append(conditions, "nama_challenge LIKE ?")
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

	err := r.db.Raw(query, args...).Scan(&contexts).Error
	if err != nil {
		return nil, err
	}

	return contexts, nil
}

func (r *masterContextRepo) CheckTriwulanExists(idTriwulan string) (bool, error) {
	var count int
	err := r.db.Raw(CheckTriwulanExistsQuery, idTriwulan).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
