package service

import (
	dto "permen_api/domain/sample/dto"
)

func (s *userIntegrationService) GetAllUserIntegrations() (data []dto.UserIntegrationResponse, err error) {
	dataDB, err := s.repo.GetAllUserIntegrations()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.UserIntegrationResponse{
			Username:    v.Username,
			Credentials: v.Credentials,
			ChannelName: v.ChannelName,
			CreatedBy:   v.CreatedBy,
			IsActive:    v.IsActive,
		})
	}

	return data, nil
}
