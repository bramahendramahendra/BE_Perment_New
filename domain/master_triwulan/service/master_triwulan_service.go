package service

import (
	dto "permen_api/domain/master_triwulan/dto"
)

func (s *masterTriwulanService) GetAllMasterTriwulan() (data []dto.MasterTriwulanResponse, err error) {
	dataDB, err := s.repo.GetAllMasterTriwulan()
	if err != nil {
		return nil, err
	}

	for _, V := range dataDB {
		data = append(data, dto.MasterTriwulanResponse{
			IdTriwulan: V.IdTriwulan,
			Triwulan:   V.Triwulan,
		})
	}

	return data, nil
}
