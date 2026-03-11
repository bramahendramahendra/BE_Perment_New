package service

import (
	"fmt"
	"mime/multipart"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

func (s *penyusunanKpiService) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.CreatePenyusunanKpiResponse, err error) {

	if file == nil {
		return data, fmt.Errorf("file Excel tidak ditemukan, pastikan mengirim file via field 'files'")
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return data, fmt.Errorf("file '%s' bukan format Excel (.xlsx)", file.Filename)
	}

	kpiSubDetails, err := ParseAndValidateExcel(file, req.Triwulan, req.Kpi)
	if err != nil {
		return data, fmt.Errorf("validasi file Excel '%s' gagal: %w", file.Filename, err)
	}

	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return data, err
	}

	idPengajuan, err := s.repo.CreatePenyusunanKpi(req, kpiSubDetails)
	if err != nil {
		return data, err
	}

	data = dto.CreatePenyusunanKpiResponse{
		IDPengajuan:    idPengajuan,
		Tahun:          req.Tahun,
		Triwulan:       req.Triwulan,
		Kostl:          req.Kostl,
		KostlTx:        req.KostlTx,
		EntryUser:      req.EntryUser,
		EntryName:      req.EntryName,
		EntryTime:      req.EntryTime,
		ApprovalPosisi: req.ApprovalPosisi,
		SaveAsDraft:    req.SaveAsDraft,
		TotalKpi:       len(req.Kpi),
		Kpi:            buildKpiResponse(idPengajuan, req.Kpi, kpiSubDetails),
		ChallengeList:  req.ChallengeList,
		MethodList:     req.MethodList,
	}

	return data, nil
}

// buildKpiResponse membangun slice PenyusunanKpiDetailItemResponse dengan
// KpiSubDetail yang sudah di-nested ke dalam masing-masing KPI sesuai indeksnya.
func buildKpiResponse(
	idPengajuan string,
	kpiList []dto.PenyusunanKpiDetailRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanKpiDetailResponse {
	result := make([]dto.PenyusunanKpiDetailResponse, len(kpiList))
	subCounter := 1

	for i, kpiItem := range kpiList {
		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		rows := kpiSubDetails[i]

		subDetails := make([]dto.PenyusunanKpiSubDetailResponse, 0, len(rows))
		for _, row := range rows {
			idSubDetail := fmt.Sprintf("%sC%03d", idPengajuan, subCounter)
			subCounter++

			subDetails = append(subDetails, dto.PenyusunanKpiSubDetailResponse{
				IdDetail:                  idDetail,
				IdSubDetail:               idSubDetail,
				NamaKpi:                   row.KPI,
				IdSubKpi:                  row.IdSubKpi,
				SubKpi:                    row.SubKPI,
				Otomatis:                  row.Otomatis,
				Polarisasi:                row.Polarisasi,
				IdPolarisasi:              row.IdPolarisasi,
				Capping:                   row.Capping,
				Bobot:                     row.Bobot,
				Glossary:                  row.Glossary,
				TargetTriwulan:            row.TargetTriwulan,
				TargetKuantitatifTriwulan: row.TargetKuantitatifTriwulan,
				TargetTahunan:             row.TargetTahunan,
				TargetKuantitatifTahunan:  row.TargetKuantitatifTahunan,
				TerdapatQualifier:         row.TerdapatQualifier,
				Qualifier:                 row.Qualifier,
				DeskripsiQualifier:        row.DeskripsiQualifier,
				TargetQualifier:           row.TargetQualifier,
				Result:                    row.Result,
				DeskripsiResult:           row.DeskripsiResult,
				Process:                   row.Process,
				DeskripsiProcess:          row.DeskripsiProcess,
				Context:                   row.Context,
				DeskripsiContext:          row.DeskripsiContext,
			})
		}

		result[i] = dto.PenyusunanKpiDetailResponse{
			IdKpi:        kpiItem.IdKpi,
			Kpi:          kpiItem.Kpi,
			Rumus:        kpiItem.Rumus,
			Persfektif:   kpiItem.Persfektif,
			KpiSubDetail: subDetails,
		}
	}

	return result
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
