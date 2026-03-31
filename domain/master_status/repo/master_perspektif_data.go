package repo

import (
	model "permen_api/domain/master_status/model"
)

const (
	GetAllMasterStatusQuery = `SELECT id_status, status_desc FROM mst_status`
)

func (r *masterStatusRepo) GetAllMasterStatus() ([]*model.MstStatus, error) {
	var masterstatuss []*model.MstStatus
	err := r.db.Raw(GetAllMasterStatusQuery).Scan(&masterstatuss).Error
	if err != nil {
		return nil, err
	}
	return masterstatuss, nil
}
