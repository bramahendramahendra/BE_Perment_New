package utils

import (
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
)

// BuildKpiResponse membangun slice PenyusunanKpiDetailResponse dari kpiRows (hasil parse Excel + lookup mst_kpi)
// dan kpiSubDetails yang sudah di-nested ke dalam masing-masing KPI sesuai indeksnya.
//
// Perubahan dari versi sebelumnya:
//   - Parameter kpiList []dto.PenyusunanKpiDetailRequest diganti dengan kpiRows []dto.PenyusunanKpiRow
//     karena KPI kini diambil dari Excel + lookup mst_kpi, bukan dari REQUEST.
//   - Persfektif diisi string kosong karena sudah tidak digunakan.
func BuildKpiResponse(
	idPengajuan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.DataKpiDetail {

	result := make([]dto.DataKpiDetail, 0, len(kpiRows))

	subCounter := 1
	for i, kpiRow := range kpiRows {
		idDetail := GenerateIDDetail(idPengajuan, i)

		rows := kpiSubDetails[kpiRow.KpiIndex]
		subDetails := make([]dto.DataKpiSubdetail, 0, len(rows))

		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			qualifier, deskripsiQualifier, targetQualifier := "", "", ""
			if strings.EqualFold(subRow.TerdapatQualifier, "Ya") {
				qualifier = subRow.Qualifier
				deskripsiQualifier = subRow.DeskripsiQualifier
				targetQualifier = subRow.TargetQualifier
			}

			subDetails = append(subDetails, dto.DataKpiSubdetail{
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
			})
		}

		result = append(result, dto.DataKpiDetail{
			IdDetail:     idDetail,
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

// BuildResultList membangun slice PenyusunanResult dari data sub KPI Excel.
// ResultList diambil dari kolom P (Result) dan Q (Deskripsi Result).
// idDetailResult menggunakan GenerateIDSubDetail yang sama dengan id_sub_detail baris tersebut.
// Hanya diisi untuk TW2 dan TW4 (isExtendedTriwulan).
func BuildResultList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanResult {
	results := []dto.PenyusunanResult{}

	subCounter := 1
	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Kolom P (Result) = namaResult, kolom Q (Deskripsi Result) = deskripsiResult
			// Hanya insert jika Result tidak kosong/nil
			if subRow.Result != nil && *subRow.Result != "" {
				results = append(results, dto.PenyusunanResult{
					IdDetailResult:  idSubDetail,
					Tahun:           tahun,
					Triwulan:        triwulan,
					NamaResult:      *subRow.Result,
					DeskripsiResult: safeDeref(subRow.DeskripsiResult),
				})
			}
		}
	}

	return results
}

// BuildProcessList membangun slice PenyusunanProcess dari data sub KPI Excel.
// ProcessList diambil dari kolom R (Process) dan S (Deskripsi Process).
// idDetailProcess menggunakan GenerateIDSubDetail yang sama dengan id_sub_detail baris tersebut.
// Hanya diisi untuk TW2 dan TW4 (isExtendedTriwulan).
func BuildProcessList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanProcess {
	processses := []dto.PenyusunanProcess{}

	subCounter := 1
	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Kolom R (Process) = namaProcess, kolom S (Deskripsi Process) = deskripsiProcess
			// Hanya insert jika Process tidak kosong/nil
			if subRow.Process != nil && *subRow.Process != "" {
				processses = append(processses, dto.PenyusunanProcess{
					IdDetailProcess:  idSubDetail,
					Tahun:            tahun,
					Triwulan:         triwulan,
					NamaProcess:      *subRow.Process,
					DeskripsiProcess: safeDeref(subRow.DeskripsiProcess),
				})
			}
		}
	}

	return processses
}

// BuildContextList membangun slice PenyusunanContext dari data sub KPI Excel.
// ContextList diambil dari kolom T (Context) dan U (Deskripsi Context).
// idDetailContext menggunakan GenerateIDSubDetail yang sama dengan id_sub_detail baris tersebut.
// Hanya diisi untuk TW2 dan TW4 (isExtendedTriwulan).
func BuildContextList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanContext {
	contexts := []dto.PenyusunanContext{}

	subCounter := 1
	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Kolom T (Context) = namaContext, kolom U (Deskripsi Context) = deskripsiContext
			// Hanya insert jika Context tidak kosong/nil
			if subRow.Context != nil && *subRow.Context != "" {
				contexts = append(contexts, dto.PenyusunanContext{
					IdDetailContext:  idSubDetail,
					Tahun:            tahun,
					Triwulan:         triwulan,
					NamaContext:      *subRow.Context,
					DeskripsiContext: safeDeref(subRow.DeskripsiContext),
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
