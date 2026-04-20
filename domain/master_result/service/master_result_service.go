package service

import (
	"fmt"

	dto "permen_api/domain/master_result/dto"
	customErrors "permen_api/errors"
)

func (s *masterResultService) GetAllMasterResult(req *dto.GetAllMasterResultRequest) (data []dto.MasterResultResponse, err error) {
	exists, err := s.repo.CheckTriwulanExists(req.Triwulan)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak ditemukan", req.Triwulan),
		}
	}

	dataDB, err := s.repo.GetAllMasterResult(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterResultResponse{
			IdResult:   v.IdResult,
			NamaResult: v.NamaResult,
			DescResult: v.DescResult,
		})
	}

	return data, nil
}
