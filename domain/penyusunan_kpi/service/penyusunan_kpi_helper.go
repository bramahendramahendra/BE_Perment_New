package service

import (
	"fmt"

	dto "permen_api/domain/penyusunan_kpi/dto"
	customErrors "permen_api/errors"
)

// resolvePenyusunanLookups menjalankan semua lookup master untuk kpiRows dan kpiSubDetails.
func (s *penyusunanKpiService) resolvePenyusunanLookups(
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) error {
	if err := s.resolveKpiMasterLookup(kpiRows); err != nil {
		return err
	}
	return s.resolveMasterLookup(kpiSubDetails)
}

// resolveKpiMasterLookup melakukan lookup mst_kpi untuk setiap KPI unik dari kolom B Excel.
// Jika ditemukan → idKpi dan rumus dari DB. Jika tidak → idKpi = "0", rumus = "0".
func (s *penyusunanKpiService) resolveKpiMasterLookup(
	kpiRows []dto.PenyusunanKpiRow,
) error {
	for i := range kpiRows {
		idKpi, _, rumus, err := s.repo.LookupKpiMaster(kpiRows[i].Kpi)
		if err != nil {
			return fmt.Errorf(
				"KPI '%s': gagal lookup master KPI: %w",
				kpiRows[i].Kpi, err,
			)
		}

		if idKpi == "0" {
			kpiRows[i].IdKpi = "0"
			kpiRows[i].Rumus = "0"
		} else {
			kpiRows[i].IdKpi = idKpi
			kpiRows[i].Rumus = rumus
		}
	}
	return nil
}

// resolveMasterLookup melakukan lookup mst_kpi dan mst_polarisasi untuk setiap baris sub KPI,
// lalu memvalidasi kesesuaian polarisasi dengan rumus di mst_kpi.
func (s *penyusunanKpiService) resolveMasterLookup(
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) error {
	for i, rows := range kpiSubDetails {
		for j := range rows {
			subRow := &kpiSubDetails[i][j]

			idKpi, kpiFromDB, _, err := s.repo.LookupKpiMaster(subRow.SubKPI)
			if err != nil {
				return fmt.Errorf(
					"KPI ke-%d, Sub KPI ke-%d ('%s'): gagal lookup master KPI: %w",
					i+1, j+1, subRow.SubKPI, err,
				)
			}
			subRow.IdSubKpi = idKpi
			subRow.SubKPI = kpiFromDB
			if idKpi != "0" {
				subRow.Otomatis = "1"
			} else {
				subRow.Otomatis = "0"
			}

			idPolarisasi, err := s.repo.LookupPolarisasi(subRow.Polarisasi)
			if err != nil {
				return &customErrors.BadRequestError{
					Message: fmt.Sprintf(
						"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi '%s' tidak valid: %s",
						i+1, j+1, subRow.SubKPI, subRow.Polarisasi, err.Error(),
					),
				}
			}
			subRow.IdPolarisasi = idPolarisasi
		}
	}
	return nil
}
