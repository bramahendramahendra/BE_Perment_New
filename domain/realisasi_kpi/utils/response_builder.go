package utils

import (
	dto "permen_api/domain/realisasi_kpi/dto"
)

// BuildSubKpiDetailList membangun slice RealisasiSubKpiDetail dari kpiRows + kpiSubDetails
// yang sudah di-enrich oleh service (IdSubDetail, Pencapaian, Skor sudah terisi).
func BuildSubKpiDetailList(
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) []dto.RealisasiSubKpiDetail {
	var result []dto.RealisasiSubKpiDetail
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for _, sub := range subRows {
			result = append(result, dto.RealisasiSubKpiDetail{
				IdSubDetail:                   sub.IdSubDetail,
				IdDetail:                      sub.IdDetail,
				KPI:                           kpiRow.Kpi,
				SubKPI:                        sub.SubKPI,
				Polarisasi:                    sub.Polarisasi,
				Capping:                       sub.Capping,
				Bobot:                         sub.Bobot,
				TargetTriwulan:                sub.TargetTriwulan,
				TargetKuantitatifTriwulan:     sub.TargetKuantitatifTriwulan,
				Qualifier:                     sub.Qualifier,
				TargetQualifier:               sub.TargetQualifier,
				Realisasi:                     sub.Realisasi,
				RealisasiKuantitatif:          sub.RealisasiKuantitatif,
				RealisasiQualifier:            sub.RealisasiQualifierVal,
				RealisasiKuantitatifQualifier: sub.RealisasiKuantitatifQualifier,
				Pencapaian:                    sub.Pencapaian,
				Skor:                          sub.Skor,
			})
		}
	}
	return result
}

// BuildResultList membangun slice RealisasiResult dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki RealisasiResult (kolom P TW2/TW4) yang dimasukkan.
func BuildResultList(
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) []dto.RealisasiResult {
	var results []dto.RealisasiResult
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for _, sub := range subRows {
			if sub.RealisasiResult != nil && *sub.RealisasiResult != "" {
				results = append(results, dto.RealisasiResult{
					IdDetailResult:  sub.IdSubDetail,
					RealisasiResult: safeDeref(sub.RealisasiResult),
					LinkResult:      safeDeref(sub.LinkResult),
				})
			}
		}
	}
	return results
}

// BuildProcessList membangun slice RealisasiProcess dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki RealisasiProcess (kolom T TW2/TW4) yang dimasukkan.
func BuildProcessList(
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) []dto.RealisasiProcess {
	var processes []dto.RealisasiProcess
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for _, sub := range subRows {
			if sub.RealisasiProcess != nil && *sub.RealisasiProcess != "" {
				processes = append(processes, dto.RealisasiProcess{
					IdDetailProcess:  sub.IdSubDetail,
					RealisasiProcess: safeDeref(sub.RealisasiProcess),
					LinkProcess:      safeDeref(sub.LinkProcess),
				})
			}
		}
	}
	return processes
}

// BuildContextList membangun slice RealisasiContext dari kpiSubDetails yang sudah di-enrich.
// Hanya baris yang memiliki RealisasiContext (kolom X TW2/TW4) yang dimasukkan.
func BuildContextList(
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) []dto.RealisasiContext {
	var contexts []dto.RealisasiContext
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for _, sub := range subRows {
			if sub.RealisasiContext != nil && *sub.RealisasiContext != "" {
				contexts = append(contexts, dto.RealisasiContext{
					IdDetailContext:  sub.IdSubDetail,
					RealisasiContext: safeDeref(sub.RealisasiContext),
					LinkContext:      safeDeref(sub.LinkContext),
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
