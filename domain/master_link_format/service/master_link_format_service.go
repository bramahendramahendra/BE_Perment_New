package service

import (
	dto "permen_api/domain/master_link_format/dto"
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterLinkFormat digunakan oleh endpoint POST /master-perspektif/get-all.
func (s *masterLinkFormatService) GetAllMasterLinkFormat() (data []dto.MasterLinkFormatResponse, err error) {
	dataDB, err := s.repo.GetAllMasterLinkFormat()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterLinkFormatResponse{
			IdLinkFormat: v.IdLinkFormat,
			UrlPrefix:    v.UrlPrefix,
			Keterangan:   v.Keterangan,
		})
	}

	return data, nil
}
