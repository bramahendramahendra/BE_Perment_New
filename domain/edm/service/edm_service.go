package service

import (
	dto "permen_api/domain/edm/dto"
)

func (s *edmService) GetKpi(req *dto.GetKpiRequest) (interface{}, error) {
	return s.edm.GetKpi(req.Periode)
}
