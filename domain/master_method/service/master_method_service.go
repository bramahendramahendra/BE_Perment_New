package service

import (
	"fmt"

	dto "permen_api/domain/master_method/dto"
	customErrors "permen_api/errors"
)

func (s *masterMethodService) GetAllMasterMethod(req *dto.GetAllMasterMethodRequest) (data []dto.MasterMethodResponse, err error) {
	exists, err := s.repo.CheckTriwulanExists(req.Triwulan)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak ditemukan", req.Triwulan),
		}
	}

	dataDB, err := s.repo.GetAllMasterMethod(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterMethodResponse{
			IdMethodUse: v.IdMethodUse,
			NamaMethod:  v.NamaMethod,
			DescMethod:  v.DescMethod,
		})
	}

	return data, nil
}
