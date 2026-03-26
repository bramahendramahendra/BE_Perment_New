package service

import (
	dto "permen_api/domain/master_tahun/dto"
	"time"
)

func (s *masterTahunService) GetAllMasterTahun() (data []dto.MasterTahunResponse, err error) {
	config, err := s.repo.GetMasterTahunConfig()
	if err != nil {
		return nil, err
	}

	thisYear := time.Now().Year()

	// Generate tahun dari (thisYear + batas_atas) s.d. (thisYear - batas_bawah)
	for tahun := thisYear + config.BatasAtas; tahun >= thisYear-config.BatasBawah; tahun-- {
		data = append(data, dto.MasterTahunResponse{Tahun: tahun})
	}

	return data, nil
}
