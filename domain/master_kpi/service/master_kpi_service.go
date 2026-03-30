package service

import (
	dto "permen_api/domain/master_kpi/dto"
)

func (s *masterKpiService) GetAllMasterKpi() (data []dto.MasterKpiResponse, err error) {
	dataDB, err := s.repo.GetAllMasterKpi()
	if err != nil {
		return data, err
	}

	for _, v := range dataDB {
		data = append(data, dto.MasterKpiResponse{
			IdKpi: v.IdKpi,
			Kpi:   v.Kpi,
			Rumus: v.Rumus,
		})
	}

	return data, nil
}
