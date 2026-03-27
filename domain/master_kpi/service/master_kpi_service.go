package service

import (
	dto "permen_api/domain/master_perspektif/dto"
)

func (s *masterPerspektifService) GetAllMasterPerspektif() (data []dto.MasterPerspektifResponse, err error) {
	dataDB, err := s.repo.GetAllMasterPerspektif()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterPerspektifResponse{
			IdPerspektif: v.IdPerspektif,
			Perspektif:   v.Perspektif,
		})
	}

	return data, nil
}
