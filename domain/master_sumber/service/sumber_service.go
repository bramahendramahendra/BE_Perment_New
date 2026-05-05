package service

import (
	dto "permen_api/domain/master_sumber/dto"
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterSumber digunakan oleh endpoint POST /master-sumber/get-all.
func (s *masterSumberService) GetAllMasterSumber() (data []dto.MasterSumberResponse, err error) {
	dataDB, err := s.repo.GetAllMasterSumber()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterSumberResponse{
			IdSumber: v.IdSumber,
			Sumber:   v.Sumber,
		})
	}

	return data, nil
}
