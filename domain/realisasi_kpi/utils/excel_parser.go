package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	dto "permen_api/domain/realisasi_kpi/dto"
	"permen_api/pkg/excel"
)

// Konstanta format Excel realisasi.
// Header ada di Row 1; data mulai dari Row 2.
const (
	DataStartRow = 2
	MaxDataRows  = 200

	SheetTW1 = "TW1"
	SheetTW2 = "TW2"
	SheetTW3 = "TW3"
	SheetTW4 = "TW4"

	PolarisasiMaximize = "Maximize"
	PolarisasiMinimize = "Minimize"

	CappingOption1 = "100%"
	CappingOption2 = "110%"

	QualifierYa    = "Ya"
	QualifierTidak = "Tidak"

	TotalBobotExpected = 100.0
	BobotTolerance     = 0.01
)

// GetMaxRowsFromEnv membaca batas maksimum baris data dari environment variable EXCEL_MAX_ROWS.
func GetMaxRowsFromEnv() int {
	val := os.Getenv("EXCEL_MAX_ROWS")
	if val == "" {
		return MaxDataRows
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return MaxDataRows
	}
	return n
}

// IsExtendedTriwulan returns true untuk TW2 dan TW4.
func IsExtendedTriwulan(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == "TW2" || upper == "TW4"
}

func ParseAndValidateRealisasiExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.RealisasiKpiRow, map[int][]dto.RealisasiKpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, maxRows)
}

// parseAndValidateExcelInternal membaca file Excel realisasi, memvalidasi, dan mengembalikan
// slice KpiRow dan map KpiSubDetailRow yang sudah terisi data.
//
// Aturan kolom:
//
//	A=No, B=KPI, C=SubKPI, D=Polarisasi, E=Capping, F=Bobot%,
//	G=TargetTriwulan, H=Qualifier (auto-fill "-"), I=TargetQualifier (auto-fill "-"),
//	J=Realisasi, K=RealisasiKuantitatif, L=RealisasiQualifier, M=RealisasiQualifierKuantitatif,
//	N=LinkDokumenSumber (semua triwulan)
//	TW2/TW4: O=Result, P=DeskripsiResult, Q=RealisasiResult, R=LinkResult,
//	         S=Process, T=DeskripsiProcess, U=RealisasiProcess, V=LinkProcess,
//	         W=Context, X=DeskripsiContext, Y=RealisasiContext, Z=LinkContext
func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	maxRows int,
) ([]dto.RealisasiKpiRow, map[int][]dto.RealisasiKpiSubDetailRow, error) {
	if maxRows <= 0 {
		return nil, nil, fmt.Errorf("maxRows harus lebih dari 0, nilai saat ini: %d", maxRows)
	}

	isTW24 := IsExtendedTriwulan(triwulan)
	targetSheet := strings.ToUpper(strings.TrimSpace(triwulan))

	allRows, err := excel.ReadSheet(file, targetSheet)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"%w. Pastikan nama sheet sesuai triwulan ('%s', '%s', '%s', atau '%s')",
			err, SheetTW1, SheetTW2, SheetTW3, SheetTW4,
		)
	}

	if len(allRows) <= DataStartRow {
		return nil, nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, DataStartRow,
		)
	}

	dataStartIdx := DataStartRow - 1
	dataEndIdx := dataStartIdx + MaxDataRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}
	limitedRows := allRows[dataStartIdx:dataEndIdx]

	kpiIndexMap := make(map[string]int)
	kpiRows := []dto.RealisasiKpiRow{}
	kpiSubDetails := make(map[int][]dto.RealisasiKpiSubDetailRow)
	totalBobot := 0.0

	prevNo := 0
	prevKpiName := ""
	prevLinkDokumen := ""

	var expectedCols int
	if isTW24 {
		expectedCols = 26 // A–Z
	} else {
		expectedCols = 14 // A–N
	}

	for rowIdx, row := range limitedRows {
		displayRow := dataStartIdx + rowIdx + 1

		for len(row) < expectedCols {
			row = append(row, "")
		}

		colA := excel.GetCell(row, 0)
		colB := excel.GetCell(row, 1)
		colC := excel.GetCell(row, 2)
		colD := excel.GetCell(row, 3)
		colE := excel.GetCell(row, 4)
		colF := excel.GetCell(row, 5)
		colG := excel.GetCell(row, 6)
		colH := excel.GetCell(row, 7)
		colI := excel.GetCell(row, 8)
		colJ := excel.GetCell(row, 9)
		colK := excel.GetCell(row, 10)
		colL := excel.GetCell(row, 11)
		colM := excel.GetCell(row, 12)
		colN := excel.GetCell(row, 13)

		if colC == "" {
			continue
		}

		if colA == "" {
			colA = strconv.Itoa(prevNo)
		}
		if colB == "" {
			colB = prevKpiName
		}
		if colN == "" {
			colN = prevLinkDokumen
		}

		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, nil, fmt.Errorf("baris %d, Kolom A (No): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		if colB == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}

		kpiKey := strings.ToLower(strings.TrimSpace(colB))
		kpiIdx, found := kpiIndexMap[kpiKey]
		if !found {
			kpiIdx = len(kpiRows)
			kpiIndexMap[kpiKey] = kpiIdx
			kpiRows = append(kpiRows, dto.RealisasiKpiRow{
				KpiIndex: kpiIdx,
				IdKpi:    "",
				Kpi:      colB,
				Rumus:    "",
			})
		}

		if colC == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		if colD == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != "Maximize" && colD != "Minimize" {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom D (Polarisasi): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colD, PolarisasiMaximize, PolarisasiMinimize,
			)
		}

		if colE == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != CappingOption1 && colE != CappingOption2 {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom E (Capping): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colE, CappingOption1, CappingOption2,
			)
		}

		if colF == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		if !strings.HasSuffix(colF, "%") {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom F (Bobot %%): harus menggunakan simbol persen, contoh: '25%%', nilai saat ini: '%s'",
				displayRow, colF,
			)
		}
		bobot, errBobot := excel.ParseFloat(strings.TrimSuffix(colF, "%"))
		if errBobot != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom F (Bobot %%): harus berupa angka 2 desimal dengan simbol persen, contoh: '25%%', nilai saat ini: '%s'",
				displayRow, colF,
			)
		}
		totalBobot += bobot

		if colG == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom G (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		if colJ == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom J (Realisasi): tidak boleh kosong", displayRow)
		}

		realisasiKuantitatif, errK := excel.ParseFloat(colK)
		if errK != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom K (Realisasi Kuantitatif): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colK,
			)
		}

		if colM != "" && colM != "-" {
			if _, errM := excel.ParseFloat(colM); errM != nil {
				return nil, nil, fmt.Errorf(
					"baris %d, Kolom M (Realisasi Qualifier Kuantitatif): harus berupa angka, nilai saat ini: '%s'",
					displayRow, colM,
				)
			}
		}

		prevNo = no
		prevKpiName = colB
		prevLinkDokumen = colN

		subRow := dto.RealisasiKpiSubDetailRow{
			No:                            no,
			KPI:                           colB,
			SubKPI:                        colC,
			Polarisasi:                    colD,
			Capping:                       strings.TrimSuffix(colE, "%"),
			Bobot:                         bobot,
			TargetTriwulan:                colG,
			Qualifier:                     colH,
			TargetQualifier:               colI,
			Realisasi:                     colJ,
			RealisasiKuantitatif:          realisasiKuantitatif,
			RealisasiQualifier:            colL,
			RealisasiKuantitatifQualifier: colM,
			LinkDokumenSumber:             &colN,
			IsTW24:                        isTW24,
		}

		if isTW24 {
			colO := excel.GetCell(row, 14)
			colP := excel.GetCell(row, 15)
			colQ := excel.GetCell(row, 16)
			colR := excel.GetCell(row, 17)
			colS := excel.GetCell(row, 18)
			colT := excel.GetCell(row, 19)
			colU := excel.GetCell(row, 20)
			colV := excel.GetCell(row, 21)
			colW := excel.GetCell(row, 22)
			colX := excel.GetCell(row, 23)
			colY := excel.GetCell(row, 24)
			colZ := excel.GetCell(row, 25)

			if colO == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom O (Result): tidak boleh kosong", displayRow)
			}
			if colP == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom P (Deskripsi Result): tidak boleh kosong", displayRow)
			}
			if colQ == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Q (Realisasi Result): tidak boleh kosong", displayRow)
			}
			if colR == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom R (Link Result): tidak boleh kosong", displayRow)
			}
			if colS == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom S (Process): tidak boleh kosong", displayRow)
			}
			if colT == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom T (Deskripsi Process): tidak boleh kosong", displayRow)
			}
			if colU == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom U (Realisasi Process): tidak boleh kosong", displayRow)
			}
			if colV == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom V (Link Process): tidak boleh kosong", displayRow)
			}
			if colW == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom W (Context): tidak boleh kosong", displayRow)
			}
			if colX == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom X (Deskripsi Context): tidak boleh kosong", displayRow)
			}
			if colY == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Y (Realisasi Context): tidak boleh kosong", displayRow)
			}
			if colZ == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Z (Link Context): tidak boleh kosong", displayRow)
			}

			subRow.Result = &colO
			subRow.DeskripsiResult = &colP
			subRow.RealisasiResult = &colQ
			subRow.LampiranEvidenceResult = &colR
			subRow.Process = &colS
			subRow.DeskripsiProcess = &colT
			subRow.RealisasiProcess = &colU
			subRow.LampiranEvidenceProcess = &colV
			subRow.Context = &colW
			subRow.DeskripsiContext = &colX
			subRow.RealisasiContext = &colY
			subRow.LampiranEvidenceContext = &colZ
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	if len(kpiRows) == 0 {
		return nil, nil, fmt.Errorf("file Excel '%s' sheet '%s' tidak memiliki data yang valid", file.Filename, targetSheet)
	}

	for _, kpiRow := range kpiRows {
		if _, ok := kpiSubDetails[kpiRow.KpiIndex]; !ok {
			return nil, nil, fmt.Errorf(
				"KPI '%s' tidak memiliki data sub KPI di file Excel '%s'",
				kpiRow.Kpi, file.Filename,
			)
		}
	}

	roundedTotal := math.Round(totalBobot*100) / 100
	if math.Abs(roundedTotal-TotalBobotExpected) > BobotTolerance {
		return nil, nil, fmt.Errorf(
			"total Bobot (Kolom F) semua KPI = %.2f%%, harus tepat 100%%",
			roundedTotal,
		)
	}

	return kpiRows, kpiSubDetails, nil
}
