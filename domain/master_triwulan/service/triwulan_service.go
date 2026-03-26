package service

import (
	dto "permen_api/domain/master_triwulan/dto"
)

func (s *masterTriwulanService) GetAllTriwulan() ([]dto.TriwulanResponse, error) {
	triwulans, err := s.repo.GetAllTriwulan()
	if err != nil {
		return nil, err
	}

	var result []dto.TriwulanResponse
	for _, t := range triwulans {
		result = append(result, dto.TriwulanResponse{
			IdTriwulan: t.IdTriwulan,
			Triwulan:   t.Triwulan,
		})
	}

	return result, nil
}
