package repo

import (
	model "permen_api/domain/master_link_format/model"
)

const (
	GetAllMasterLinkFormatQuery = `SELECT id_link_format, url_prefix, keterangan FROM mst_link_format`
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterLinkFormat digunakan oleh endpoint POST /master-link_format/get-all.
func (r *masterLinkFormatRepo) GetAllMasterLinkFormat() ([]*model.MstLinkFormat, error) {
	var masterlinkformats []*model.MstLinkFormat
	err := r.db.Raw(GetAllMasterLinkFormatQuery).Scan(&masterlinkformats).Error
	if err != nil {
		return nil, err
	}
	return masterlinkformats, nil
}
