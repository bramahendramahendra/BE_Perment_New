package service

import (
	dto "permen_api/domain/edm/dto"
)

func (s *edmService) GetRealisasi(req *dto.GetRealisasiRequest) (interface{}, error) {
	return s.edm.GetDataKPI(req.Tahun, req.Triwulan, req.IdKpi)
}
