package service

import (
	dto "permen_api/domain/master_divisi/dto"
)

func (s *masterDivisiService) GetAllMasterDivisi() (data []dto.MasterDivisiResponse, err error) {
	dataDB, err := s.repo.GetAllMasterDivisi()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterDivisiResponse{
			Kostl:   v.Kostl,
			KostlTx: v.KostlTx,
		})
	}

	return data, nil
}
