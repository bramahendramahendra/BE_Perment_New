package repo

import (
	model "permen_api/domain/master_divisi/model"
)

const (
	GetAllMasterDivisiQuery = `
		SELECT DISTINCT KOSTL, KOSTL_TX
		FROM ` + "`user`" + `
		WHERE (WERKS = 'KP00' OR (WERKS = 'PL00' AND BTRTL = 'PL01') OR (WERKS = 'KI00' AND BTRTL = 'KI01'))
		  AND HILFM != '098'
		  AND KOSTL NOT IN ('PS98000', 'PS98200')
		  AND LEFT(PERNR, 1) != '9'
		ORDER BY HILFM ASC
	`
)

func (r *masterDivisiRepo) GetAllMasterDivisi() ([]*model.MstDivisi, error) {
	var masterdivisis []*model.MstDivisi
	err := r.db.Raw(GetAllMasterDivisiQuery).Scan(&masterdivisis).Error
	if err != nil {
		return nil, err
	}
	return masterdivisis, nil
}
