package service

import (
	"fmt"
	"mime/multipart"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

func (s *penyusunanKpiService) InsertPenyusunanKpi(
	req *dto.InsertPenyusunanKpiRequest,
	files []*multipart.FileHeader,
) (*dto.InsertPenyusunanKpiResult, error) {

	if len(files) == 0 {
		return nil, fmt.Errorf("tidak ada file Excel yang dikirim, harus mengirim tepat 1 file Excel")
	}
	if len(files) > 1 {
		return nil, fmt.Errorf(
			"hanya boleh mengirim 1 file Excel (diterima %d file). "+
				"Semua data sub KPI dari semua KPI harus digabung dalam 1 file",
			len(files),
		)
	}

	file := files[0]

	kpiSubDetails, err := ParseAndValidateExcel(file, req.Triwulan, req.Kpi)
	if err != nil {
		return nil, fmt.Errorf("validasi file Excel '%s' gagal: %w", file.Filename, err)
	}

	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return nil, err
	}

	idPengajuan, err := s.repo.InsertPenyusunanKpi(req, kpiSubDetails)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan data KPI: %w", err)
	}

	return &dto.InsertPenyusunanKpiResult{
		IDPengajuan:   idPengajuan,
		KpiSubDetails: kpiSubDetails,
	}, nil
}

// resolveMasterLookup melakukan lookup mst_kpi dan mst_polarisasi untuk setiap
// baris sub KPI, lalu memvalidasi kesesuaian polarisasi dengan rumus di mst_kpi.
func (s *penyusunanKpiService) resolveMasterLookup(
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) error {
	for i, rows := range kpiSubDetails {
		for j := range rows {
			subRow := &kpiSubDetails[i][j]

			idKpi, kpiFromDB, rumusMstKpi, err := s.repo.LookupSubKpiMaster(subRow.SubKPI)
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
				return fmt.Errorf(
					"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi '%s' tidak valid: %w",
					i+1, j+1, subRow.SubKPI, subRow.Polarisasi, err,
				)
			}
			subRow.IdPolarisasi = idPolarisasi

			if subRow.IdSubKpi != "0" {
				polarisasiMaster := "Maximize"
				if rumusMstKpi == "0" {
					polarisasiMaster = "Minimize"
				}
				if idPolarisasi != rumusMstKpi {
					return fmt.Errorf(
						"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi tidak sesuai master. "+
							"Excel: '%s' (id=%s), master KPI: '%s' (id=%s). "+
							"Periksa kembali kolom D pada file Excel",
						i+1, j+1, subRow.SubKPI,
						subRow.Polarisasi, idPolarisasi,
						polarisasiMaster, rumusMstKpi,
					)
				}
			}
		}
	}
	return nil
}

// sheetName mengembalikan nama sheet Excel berdasarkan triwulan.
func sheetName(triwulan string) string {
	if strings.EqualFold(triwulan, "TW4") {
		return "TW 4"
	}
	return "Selain TW 4"
}
