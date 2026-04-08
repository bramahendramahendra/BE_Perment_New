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
) []dto.PenyusunanKpiDetailResponse {

	result := make([]dto.PenyusunanKpiDetailResponse, 0, len(kpiRows))

	subCounter := 1
	for i, kpiRow := range kpiRows {
		idDetail := GenerateIDDetail(idPengajuan, i)

		rows := kpiSubDetails[kpiRow.KpiIndex]
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
			})
		}

		result = append(result, dto.PenyusunanKpiDetailResponse{
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

// BuildMethodList membangun slice PenyusunanMethod dari data sub KPI Excel.
// MethodList diambil dari kolom R (Process) dan S (Deskripsi Process).
// idDetailMethod menggunakan GenerateIDSubDetail yang sama dengan id_sub_detail baris tersebut.
// Hanya diisi untuk TW2 dan TW4 (isExtendedTriwulan).
func BuildMethodList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanMethod {
	methods := []dto.PenyusunanMethod{}

	subCounter := 1
	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Kolom R (Process) = namaMethod, kolom S (Deskripsi Process) = deskripsiMethod
			// Hanya insert jika Process tidak kosong/nil
			if subRow.Process != nil && *subRow.Process != "" {
				methods = append(methods, dto.PenyusunanMethod{
					IdDetailMethod:  idSubDetail,
					Tahun:           tahun,
					Triwulan:        triwulan,
					NamaMethod:      *subRow.Process,
					DeskripsiMethod: safeDeref(subRow.DeskripsiProcess),
				})
			}
		}
	}

	return methods
}

// BuildChallengeList membangun slice PenyusunanChallenge dari data sub KPI Excel.
// ChallengeList diambil dari kolom T (Context) dan U (Deskripsi Context).
// idDetailChallenge menggunakan GenerateIDSubDetail yang sama dengan id_sub_detail baris tersebut.
// Hanya diisi untuk TW2 dan TW4 (isExtendedTriwulan).
func BuildChallengeList(
	idPengajuan string,
	tahun string,
	triwulan string,
	kpiRows []dto.PenyusunanKpiRow,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanChallenge {
	challenges := []dto.PenyusunanChallenge{}

	subCounter := 1
	for _, kpiRow := range kpiRows {
		rows := kpiSubDetails[kpiRow.KpiIndex]
		for _, subRow := range rows {
			idSubDetail := GenerateIDSubDetail(idPengajuan, subCounter)
			subCounter++

			// Kolom T (Context) = namaChallenge, kolom U (Deskripsi Context) = deskripsiChallenge
			// Hanya insert jika Context tidak kosong/nil
			if subRow.Context != nil && *subRow.Context != "" {
				challenges = append(challenges, dto.PenyusunanChallenge{
					IdDetailChallenge:  idSubDetail,
					Tahun:              tahun,
					Triwulan:           triwulan,
					NamaChallenge:      *subRow.Context,
					DeskripsiChallenge: safeDeref(subRow.DeskripsiContext),
				})
			}
		}
	}

	return challenges
}

// safeDeref mengembalikan value dari pointer string, atau string kosong jika nil.
func safeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
