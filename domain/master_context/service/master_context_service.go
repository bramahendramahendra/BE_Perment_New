package service

import (
	"fmt"

	dto "permen_api/domain/master_context/dto"
	customErrors "permen_api/errors"
)

func (s *masterContextService) GetAllMasterContext(req *dto.GetAllMasterContextRequest) (data []dto.MasterContextResponse, err error) {
	exists, err := s.repo.CheckTriwulanExists(req.Triwulan)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak ditemukan", req.Triwulan),
		}
	}

	dataDB, err := s.repo.GetAllMasterContext(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterContextResponse{
			IdContext:   v.IdChallenge,
			NamaContext: v.NamaChallenge,
			DescContext: v.DescChallenge,
		})
	}

	return data, nil
}
