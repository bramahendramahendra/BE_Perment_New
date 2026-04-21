package utils

import (
	dto "permen_api/domain/realisasi_kpi/dto"
)

// BuildSubKpiDetailList membangun slice RealisasiSubKpiDetail dari kpiRows + kpiSubDetails
// yang sudah di-enrich oleh service (IdSubDetail, Pencapaian, Skor sudah terisi).
func BuildKpiResponse(
	idPengajuan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) []dto.DataKpiDetail {

	result := make([]dto.DataKpiDetail, 0, len(kpiRows))

	for _, kpiRow := range kpiRows {

		rows := kpiSubDetails[kpiRow.KpiIndex]
		subDetails := make([]dto.DataKpiSubdetail, 0, len(rows))

		for _, subRow := range rows {

			subDetails = append(subDetails, dto.DataKpiSubdetail{
				IdSubDetail:                   subRow.IdSubDetail,
				IdSubKpi:                      subRow.IdDetail,
				SubKpi:                        subRow.SubKPI,
				Polarisasi:                    subRow.Polarisasi,
				Capping:                       subRow.Capping,
				Bobot:                         subRow.Bobot,
				TargetTriwulan:                subRow.TargetTriwulan,
				TargetKuantitatifTriwulan:     subRow.TargetKuantitatifTriwulan,
				Qualifier:                     subRow.Qualifier,
				TargetQualifier:               subRow.TargetQualifier,
				Realisasi:                     subRow.Realisasi,
				RealisasiKuantitatif:          subRow.RealisasiKuantitatif,
				RealisasiQualifier:            subRow.RealisasiQualifierVal,
				RealisasiKuantitatifQualifier: subRow.RealisasiKuantitatifQualifier,
				Pencapaian:                    subRow.Pencapaian,
				Skor:                          subRow.Skor,
			})
		}

		result = append(result, dto.DataKpiDetail{
			IdDetail:     kpiRow.IdDetail,
			IdKpi:        kpiRow.IdKpi,
			Kpi:          kpiRow.Kpi,
			Rumus:        kpiRow.Rumus,
			Persfektif:   "",
			TotalSubKpi:  len(rows),
			KpiSubDetail: subDetails,
		})
	}
	return result
}

// BuildResultList membangun slice DataResult dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki DataResult (kolom P TW2/TW4) yang dimasukkan.
func BuildResultList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) []dto.DataResult {
	results := []dto.DataResult{}

	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			if subRow.Result != nil && *subRow.Result != "" {
				results = append(results, dto.DataResult{
					IdDetailResult:   subRow.IdSubDetail,
					Tahun:            tahun,
					Triwulan:         triwulan,
					NamaResult:       *subRow.Result,
					DeskripsiResult:  safeDeref(subRow.DeskripsiResult),
					RealisasiResult:  safeDeref(subRow.RealisasiResult),
					LampiranEvidence: safeDeref(subRow.LampiranEvidenceResult),
				})
			}
		}
	}

	return results
}

// BuildProcessList membangun slice DataProcess dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki DataProcess (kolom T TW2/TW4) yang dimasukkan.
func BuildProcessList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) []dto.DataProcess {
	processses := []dto.DataProcess{}

	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			if subRow.Process != nil && *subRow.Process != "" {
				processses = append(processses, dto.DataProcess{
					IdDetailProcess:  subRow.IdSubDetail,
					Tahun:            tahun,
					Triwulan:         triwulan,
					NamaProcess:      *subRow.Process,
					DeskripsiProcess: safeDeref(subRow.DeskripsiProcess),
					RealisasiProcess: safeDeref(subRow.RealisasiProcess),
					LampiranEvidence: safeDeref(subRow.LampiranEvidenceProcess),
				})
			}
		}
	}

	return processses
}

// BuildContextList membangun slice DataContext dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki DataContext (kolom X TW2/TW4) yang dimasukkan.
func BuildContextList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) []dto.DataContext {
	var contexts []dto.DataContext

	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			if subRow.Context != nil && *subRow.Context != "" {
				contexts = append(contexts, dto.DataContext{
					IdDetailContext:  subRow.IdSubDetail,
					Tahun:            tahun,
					Triwulan:         triwulan,
					NamaContext:      *subRow.Context,
					DeskripsiContext: safeDeref(subRow.DeskripsiContext),
					RealisasiContext: safeDeref(subRow.RealisasiContext),
					LampiranEvidence: safeDeref(subRow.LampiranEvidenceContext),
				})
			}
		}
	}
	return contexts
}

// safeDeref mengembalikan value dari pointer string, atau string kosong jika nil.
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
