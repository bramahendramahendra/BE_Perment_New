package repo

import (
	"strings"

	dto "permen_api/domain/user/dto"
	model "permen_api/domain/user/model"
)

const (
	GetAllUserBaseQuery = `
		SELECT PERNR, SNAME, IFNULL(JGPG, '') JGPG, IFNULL(ESELON, '') ESELON,
			WERKS, WERKS_TX, BTRTL, BTRTL_TX, KOSTL, KOSTL_TX, ORGEH, ORGEH_TX,
			IFNULL(STELL, '') STELL, IFNULL(STELL_TX, '') STELL_TX,
			IFNULL(PLANS, '') PLANS, IFNULL(PLANS_TX, '') PLANS_TX,
			IFNULL(HILFM, '') HILFM, HTEXT, BRANCH, MAINBR, IS_PEMIMPIN,
			IFNULL(ADMIN_LEVEL, '') ADMIN_LEVEL,
			IFNULL(ORGEH_PGS, '') ORGEH_PGS, IFNULL(ORGEH_PGS_TX, '') ORGEH_PGS_TX,
			IFNULL(PLANS_PGS, '') PLANS_PGS, IFNULL(PLANS_PGS_TX, '') PLANS_PGS_TX,
			IFNULL(BRANCH_PGS, '') BRANCH_PGS, IFNULL(HILFM_PGS, '') HILFM_PGS,
			IFNULL(HTEXT_PGS, '') HTEXT_PGS, TIPE_UKER, REKENING,
			IFNULL(NPWP, '') NPWP, REGION, RGDESC, BRDESC, MBDESC
		FROM ` + "`user`"
)

func (r *userRepo) GetAllUser(req *dto.GetAllUserRequest) ([]*model.User, error) {
	var users []*model.User

	query := GetAllUserBaseQuery

	var conditions []string
	var args []interface{}

	if req.Branch != "" {
		conditions = append(conditions, "BRANCH = ?")
		args = append(args, req.Branch)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY SNAME ASC"

	err := r.db.Raw(query, args...).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
