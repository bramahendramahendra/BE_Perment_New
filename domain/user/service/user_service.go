package service

import (
	dto "permen_api/domain/user/dto"
)

func (s *userService) GetAllUser(req *dto.GetAllUserRequest) (data []dto.UserResponse, err error) {
	dataDB, err := s.repo.GetAllUser(req)
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.UserResponse{
			Pernr: v.PERNR,
			Sname: v.SNAME,
		})
	}

	return data, nil
}
