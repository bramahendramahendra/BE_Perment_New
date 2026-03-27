package service

import (
	"fmt"

	dto "permen_api/domain/master_challenge/dto"
	customErrors "permen_api/errors"
)

func (s *masterChallengeService) GetAllMasterChallenge(req *dto.GetAllMasterChallengeRequest) (data []dto.MasterChallengeResponse, err error) {
	exists, err := s.repo.CheckTriwulanExists(req.Triwulan)
	if err != nil {
		return nil, err
	}
	if !exists {
		// return nil, fmt.Errorf("triwulan '%s' tidak ditemukan", req.Triwulan)
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak ditemukan", req.Triwulan),
		}
	}

	dataDB, err := s.repo.GetAllMasterChallenge(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterChallengeResponse{
			IdChallenge:   v.IdChallenge,
			NamaChallenge: v.NamaChallenge,
			DescChallenge: v.DescChallenge,
		})
	}

	return data, nil
}
