package repo

import (
	model "permen_api/domain/master_kpi/model"
)

const (
	GetAllMasterKpiQuery = `SELECT id_kpi, kpi, rumus FROM mst_kpi ORDER BY id_kpi ASC`
)

func (r *masterKpiRepo) GetAllMasterKpi() ([]*model.MstKpi, error) {
	var masterkpis []*model.MstKpi
	err := r.db.Raw(GetAllMasterKpiQuery).Scan(&masterkpis).Error
	if err != nil {
		return nil, err
	}
	return masterkpis, nil
}
