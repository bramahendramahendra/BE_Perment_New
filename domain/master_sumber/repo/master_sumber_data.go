package repo

import (
	model "permen_api/domain/master_sumber/model"
)

const (
	GetAllMasterSumberQuery = `SELECT id_sumber, sumber FROM mst_sumber`
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterSumber digunakan oleh endpoint POST /master-sumber/get-all.
func (r *masterSumberRepo) GetAllMasterSumber() ([]*model.MstSumber, error) {
	var mastersumbers []*model.MstSumber
	err := r.db.Raw(GetAllMasterSumberQuery).Scan(&mastersumbers).Error
	if err != nil {
		return nil, err
	}
	return mastersumbers, nil
}
