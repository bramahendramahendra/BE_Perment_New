package repo

import (
	model "permen_api/domain/master_perspektif/model"
)

const (
	GetAllMasterPerspektifQuery = `SELECT id_perspektif, perspektif FROM mst_perspektif`
)

func (r *masterPerspektifRepo) GetAllMasterPerspektif() ([]*model.MstPerspektif, error) {
	var masterperspektifs []*model.MstPerspektif
	err := r.db.Raw(GetAllMasterPerspektifQuery).Scan(&masterperspektifs).Error
	if err != nil {
		return nil, err
	}
	return masterperspektifs, nil
}
