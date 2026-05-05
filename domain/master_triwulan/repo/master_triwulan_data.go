package repo

import (
	model "permen_api/domain/master_triwulan/model"
)

const (
	GetAllMasterTriwulanQuery = `SELECT id_triwulan, triwulan FROM mst_triwulan`
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterTriwulan digunakan oleh endpoint POST /master-triwulan/get-all.
func (r *masterTriwulanRepo) GetAllMasterTriwulan() ([]*model.MstTriwulan, error) {
	var mastertriwulans []*model.MstTriwulan
	err := r.db.Raw(GetAllMasterTriwulanQuery).Scan(&mastertriwulans).Error
	if err != nil {
		return nil, err
	}
	return mastertriwulans, nil
}
