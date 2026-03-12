package service

import (
	"fmt"
	"mime/multipart"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/utils"
	customErrors "permen_api/errors"
)

// =============================================================================
// VALIDATE
// =============================================================================

func (s *penyusunanKpiService) ValidatePenyusunanKpi(
	req *dto.ValidatePenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidatePenyusunanKpiResponse, err error) {

	// User error: tidak mengirim file
	if file == nil {
		return data, &customErrors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files'",
		}
	}

	// User error: format file salah
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("file '%s' bukan format Excel (.xlsx)", file.Filename),
		}
	}

	// User error: isi Excel tidak valid (kolom salah, bobot salah, KPI tidak cocok, dll)
	kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan, req.Kpi)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// User error: polarisasi tidak sesuai master, lookup gagal karena data tidak valid
	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return data, err
	}

	// System error: gagal simpan ke DB — biarkan naik apa adanya
	idPengajuan, err := s.repo.ValidatePenyusunanKpi(req, kpiSubDetails)
	if err != nil {
		return data, err
	}

	data = dto.ValidatePenyusunanKpiResponse{
		IDPengajuan: idPengajuan,
		Tahun:       req.Tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.DivisiResponse{
			Kostl:   req.Divisi.Kostl,
			KostlTx: req.Divisi.KostlTx,
		},
		Entry: dto.EntryResponse{
			EntryUser: req.EntryUser,
			EntryName: req.EntryName,
			EntryTime: req.EntryTime,
		},
		TotalKpi:      len(req.Kpi),
		Kpi:           utils.BuildKpiResponse(idPengajuan, req.Kpi, kpiSubDetails),
		ChallengeList: req.ChallengeList,
		MethodList:    req.MethodList,
	}

	return data, nil
}

// =============================================================================
// CREATE
// =============================================================================

func (s *penyusunanKpiService) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
) (data dto.CreatePenyusunanKpiResponse, err error) {
	// User error (idPengajuan tidak ada) atau system error (DB) — repo sudah wrap dengan tipe yang tepat
	if err := s.repo.CreatePenyusunanKpi(req); err != nil {
		return data, err
	}

	data = dto.CreatePenyusunanKpiResponse{
		IdPengajuan:  req.IdPengajuan,
		ApprovalList: req.ApprovalList,
	}

	return data, nil
}

// =============================================================================
// HELPER — buildKpiResponse
// =============================================================================

// buildKpiResponse membangun slice PenyusunanKpiDetailResponse dengan
// KpiSubDetail yang sudah di-nested ke dalam masing-masing KPI sesuai indeksnya.
func buildKpiResponse(
	idPengajuan string,
	kpiList []dto.PenyusunanKpiDetailRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanKpiDetailResponse {

	result := make([]dto.PenyusunanKpiDetailResponse, 0, len(kpiList))

	subCounter := 1
	for i, kpiItem := range kpiList {
		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)

		rows := kpiSubDetails[i]
		subDetails := make([]dto.PenyusunanKpiSubDetailResponse, 0, len(rows))

		for _, subRow := range rows {
			idSubDetail := fmt.Sprintf("%sC%03d", idPengajuan, subCounter)
			subCounter++

			qualifier, deskripsiQualifier, targetQualifier := "", "", ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				qualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			subDetails = append(subDetails, dto.PenyusunanKpiSubDetailResponse{
				IdSubDetail:               idSubDetail,
				IdSubKpi:                  subRow.IdSubKpi,
				SubKpi:                    subRow.SubKPI,
				Otomatis:                  subRow.Otomatis,
				Polarisasi:                subRow.Polarisasi,
				IdPolarisasi:              subRow.IdPolarisasi,
				Capping:                   subRow.Capping,
				Bobot:                     subRow.Bobot,
				Glossary:                  subRow.Glossary,
				TargetTriwulan:            subRow.TargetTriwulan,
				TargetKuantitatifTriwulan: subRow.TargetKuantitatifTriwulan,
				TargetTahunan:             subRow.TargetTahunan,
				TargetKuantitatifTahunan:  subRow.TargetKuantitatifTahunan,
				TerdapatQualifier:         subRow.TerdapatQualifier,
				Qualifier:                 qualifier,
				DeskripsiQualifier:        deskripsiQualifier,
				TargetQualifier:           targetQualifier,
				Result:                    subRow.Result,
				DeskripsiResult:           subRow.DeskripsiResult,
				Process:                   subRow.Process,
				DeskripsiProcess:          subRow.DeskripsiProcess,
				Context:                   subRow.Context,
				DeskripsiContext:          subRow.DeskripsiContext,
			})
		}

		result = append(result, dto.PenyusunanKpiDetailResponse{
			IdDetail:     idDetail,
			IdKpi:        kpiItem.IdKpi,
			Kpi:          kpiItem.Kpi,
			Rumus:        kpiItem.Rumus,
			Persfektif:   kpiItem.Persfektif,
			TotalSubKpi:  len(rows),
			KpiSubDetail: subDetails,
		})
	}

	return result
}

// =============================================================================
// HELPER — resolveMasterLookup
// =============================================================================

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
				// System error: query DB gagal
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
				// User error: polarisasi yang diisi di Excel tidak valid
				return &customErrors.BadRequestError{
					Message: fmt.Sprintf(
						"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi '%s' tidak valid: %s",
						i+1, j+1, subRow.SubKPI, subRow.Polarisasi, err.Error(),
					),
				}
			}
			subRow.IdPolarisasi = idPolarisasi

			if subRow.IdSubKpi != "0" {
				polarisasiMaster := "Maximize"
				if rumusMstKpi == "0" {
					polarisasiMaster = "Minimize"
				}
				if idPolarisasi != rumusMstKpi {
					// User error: polarisasi tidak cocok dengan master KPI
					return &customErrors.BadRequestError{
						Message: fmt.Sprintf(
							"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi tidak sesuai master. "+
								"Excel: '%s' (id=%s), master KPI: '%s' (id=%s). "+
								"Periksa kembali kolom D pada file Excel",
							i+1, j+1, subRow.SubKPI,
							subRow.Polarisasi, idPolarisasi,
							polarisasiMaster, rumusMstKpi,
						),
					}
				}
			}
		}
	}
	return nil
}
