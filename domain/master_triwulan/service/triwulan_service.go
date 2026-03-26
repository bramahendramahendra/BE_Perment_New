package service

import (
	dto "permen_api/domain/master_triwulan/dto"
)

func (s *masterTriwulanService) GetAllTriwulan() (data []dto.TriwulanResponse, err error) {
	dataDB, err := s.repo.GetAllTriwulan()
	if err != nil {
		return nil, err
	}

	for _, V := range dataDB {
		data = append(data, dto.TriwulanResponse{
			IdTriwulan: V.IdTriwulan,
			Triwulan:   V.Triwulan,
		})
	}

	return data, nil
}
