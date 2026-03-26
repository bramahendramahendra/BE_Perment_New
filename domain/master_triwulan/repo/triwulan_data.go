package repo

import (
	"permen_api/domain/master_triwulan/model"
)

const (
	GetAllTriwulanQuery = `SELECT id_triwulan, triwulan FROM mst_triwulan`
)

func (r *masterTriwulanRepo) GetAllTriwulan() ([]*model.MstTriwulan, error) {
	var triwulans []*model.MstTriwulan
	err := r.db.Raw(GetAllTriwulanQuery).Scan(&triwulans).Error
	if err != nil {
		return nil, err
	}
	return triwulans, nil
}
