package service

import (
	dto "permen_api/domain/master_status/dto"
)

func (s *masterStatusService) GetAllMasterStatus() (data []dto.MasterStatusResponse, err error) {
	dataDB, err := s.repo.GetAllMasterStatus()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterStatusResponse{
			IdStatus:   v.IdStatus,
			StatusDesc: v.StatusDesc,
		})
	}

	return data, nil
}
