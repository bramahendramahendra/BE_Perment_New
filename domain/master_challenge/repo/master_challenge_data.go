package repo

import (
	"strings"

	dto "permen_api/domain/master_challenge/dto"
	model "permen_api/domain/master_challenge/model"
)

const (
	GetAllMasterChallengeBaseQuery = `
		SELECT id_challenge, nama_challenge, desc_challenge, tahun, triwulan, entry_user, entry_name, entry_time
		FROM mst_challenge`
)

func (r *masterChallengeRepo) GetAllMasterChallenge(req *dto.GetAllMasterChallengeRequest) ([]*model.MstChallenge, error) {
	var challenges []*model.MstChallenge

	query := GetAllMasterChallengeBaseQuery

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

	err := r.db.Raw(query, args...).Scan(&challenges).Error
	if err != nil {
		return nil, err
	}

	return challenges, nil
}
