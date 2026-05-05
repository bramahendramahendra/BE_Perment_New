package repo

import (
	model "permen_api/domain/master_tahun/model"
)

const (
	GetMasterTahunConfigQuery = `SELECT id, batas_atas, batas_bawah FROM mst_tahun LIMIT 1`
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterTahun digunakan oleh endpoint POST /master-tahun/get-all.
func (r *masterTahunRepo) GetAllMasterTahun() (*model.MstTahun, error) {
	var mstTahun model.MstTahun
	err := r.db.Raw(GetMasterTahunConfigQuery).Scan(&mstTahun).Error
	if err != nil {
		return nil, err
	}
	return &mstTahun, nil
}
