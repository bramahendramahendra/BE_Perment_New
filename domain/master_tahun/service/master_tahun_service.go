package service

import (
	dto "permen_api/domain/master_tahun/dto"
	"time"
)

func (s *masterTahunService) GetAllMasterTahun() (data []dto.MasterTahunResponse, err error) {
	thisYear := time.Now().Year()

	// 1 tahun ke depan
	data = append(data, dto.MasterTahunResponse{Tahun: thisYear + 1})

	// Tahun sekarang s.d. 4 tahun ke belakang (total 5 tahun)
	for i := 0; i < 5; i++ {
		data = append(data, dto.MasterTahunResponse{Tahun: thisYear - i})
	}

	return data, nil
}
