package service

import (
	"fmt"

	dto "permen_api/domain/master_process/dto"
	customErrors "permen_api/errors"
)

func (s *masterProcessService) GetAllMasterProcess(req *dto.GetAllMasterProcessRequest) (data []dto.MasterProcessResponse, err error) {
	exists, err := s.repo.CheckTriwulanExists(req.Triwulan)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak ditemukan", req.Triwulan),
		}
	}

	dataDB, err := s.repo.GetAllMasterProcess(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterProcessResponse{
			IdProcess:   v.IdMethodUse,
			NamaProcess: v.NamaMethod,
			DescProcess: v.DescMethod,
		})
	}

	return data, nil
}
