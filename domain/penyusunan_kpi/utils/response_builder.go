package utils

import (
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

// BuildKpiResponse membangun slice PenyusunanKpiDetailResponse dengan
// KpiSubDetail yang sudah di-nested ke dalam masing-masing KPI sesuai indeksnya.
func BuildKpiResponse(
	idPengajuan string,
	kpiList []dto.PenyusunanKpiDetailRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanKpiDetailResponse {

	result := make([]dto.PenyusunanKpiDetailResponse, 0, len(kpiList))

	subCounter := 1
	for i, kpiItem := range kpiList {
		idDetail := GenerateIDDetail(idPengajuan, i)

		rows := kpiSubDetails[i]
		subDetails := make([]dto.PenyusunanKpiSubDetailResponse, 0, len(rows))

		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
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
