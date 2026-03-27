package service

import (
	dto "permen_api/domain/master_challenge/dto"
)

func (s *masterChallengeService) GetAllMasterChallenge(req *dto.GetAllMasterChallengeRequest) (data []dto.MasterChallengeResponse, err error) {
	dataDB, err := s.repo.GetAllMasterChallenge(req)
	if err != nil {
		return nil, err
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
