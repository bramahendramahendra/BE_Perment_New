package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/pkg/excel"
)

const (
	DataStartRow = 3
	MaxDataRows  = 13

	SheetTW1 = "TW1"
	SheetTW2 = "TW2"
	SheetTW3 = "TW3"
	SheetTW4 = "TW4"

	TriwulanTW1 = "TW1"
	TriwulanTW2 = "TW2"
	TriwulanTW3 = "TW3"
	TriwulanTW4 = "TW4"

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

// IsExtendedTriwulan mengembalikan true jika triwulan adalah TW2 atau TW4.
func IsExtendedTriwulan(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == strings.ToUpper(TriwulanTW2) || upper == strings.ToUpper(TriwulanTW4)
}

// ParseAndValidateExcel membaca file Excel, memvalidasi isi, dan mengembalikan:
//   - kpiRows       : slice KPI unik dari kolom B (urutan kemunculan pertama)
//   - kpiSubDetails : map[kpiIndex] -> []PenyusunanKpiSubDetailRow
func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.PenyusunanKpiRow, map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseInternal(file, triwulan, maxRows)
}

func parseInternal(
	file *multipart.FileHeader,
	triwulan string,
	maxRows int,
) ([]dto.PenyusunanKpiRow, map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	if maxRows <= 0 {
		return nil, nil, fmt.Errorf("maxRows harus lebih dari 0, nilai saat ini: %d", maxRows)
	}

	targetSheet := strings.ToUpper(strings.TrimSpace(triwulan))
	allRows, err := excel.ReadSheet(file, targetSheet)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"%w. Pastikan nama sheet di file sesuai dengan triwulan ('%s', '%s', '%s', atau '%s')",
			err, SheetTW1, SheetTW2, SheetTW3, SheetTW4,
		)
	}

	if len(allRows) < DataStartRow {
		return nil, nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, DataStartRow,
		)
	}

	isExtendedTriwulan := IsExtendedTriwulan(triwulan)
	isTW4 := strings.EqualFold(triwulan, TriwulanTW4)

	dataStartIdx := DataStartRow - 1
	dataEndIdx := dataStartIdx + maxRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}
	limitedRows := allRows[dataStartIdx:dataEndIdx]

	kpiIndexMap := make(map[string]int)
	kpiRows := []dto.PenyusunanKpiRow{}
	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)
	totalBobot := 0.0

	expectedCols := 15
	if isExtendedTriwulan {
		expectedCols = 21
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
		colO := excel.GetCell(row, 14)

		var colP, colQ, colR, colS, colT, colU string
		if isExtendedTriwulan {
			colP = excel.GetCell(row, 15)
			colQ = excel.GetCell(row, 16)
			colR = excel.GetCell(row, 17)
			colS = excel.GetCell(row, 18)
			colT = excel.GetCell(row, 19)
			colU = excel.GetCell(row, 20)
		}

		if colA == "" && colB == "" && colC == "" {
			continue
		}

		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, nil, fmt.Errorf("baris %d, Kolom A (NO): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		if colB == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}
		kpiKey := strings.ToLower(strings.TrimSpace(colB))
		kpiIdx, found := kpiIndexMap[kpiKey]
		if !found {
			kpiIdx = len(kpiRows)
			kpiIndexMap[kpiKey] = kpiIdx
			kpiRows = append(kpiRows, dto.PenyusunanKpiRow{
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
		if colD != PolarisasiMaximize && colD != PolarisasiMinimize {
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
			return nil, nil, fmt.Errorf("baris %d, Kolom G (Glossary): tidak boleh kosong", displayRow)
		}

		if colH == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom H (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		targetKuantitatifTriwulan, errI := excel.ParseFloat(colI)
		if errI != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom I (Target Kuantitatif Triwulanan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colI,
			)
		}

		if colJ == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom J (Target Tahunan): tidak boleh kosong", displayRow)
		}

		targetKuantitatifTahunan, errK := excel.ParseFloat(colK)
		if errK != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom K (Target Kuantitatif Tahunan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colK,
			)
		}

		if colL == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): tidak boleh kosong", displayRow)
		}
		if !strings.EqualFold(colL, QualifierYa) && !strings.EqualFold(colL, QualifierTidak) {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom L (Terdapat Qualifier): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colL, QualifierYa, QualifierTidak,
			)
		}

		if strings.EqualFold(colL, QualifierYa) {
			if colM == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom M (Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
			if colN == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom N (Deskripsi Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
			if colO == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom O (Target Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
		}

		if isExtendedTriwulan {
			if colP == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom P (Result): tidak boleh kosong", displayRow)
			}
			if colQ == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Q (Deskripsi Result): tidak boleh kosong", displayRow)
			}
			if colR == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom R (Process): tidak boleh kosong", displayRow)
			}
			if colS == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom S (Deskripsi Process): tidak boleh kosong", displayRow)
			}
			if colT == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom T (Context): tidak boleh kosong", displayRow)
			}
			if colU == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom U (Deskripsi Context): tidak boleh kosong", displayRow)
			}
		}

		itemQualifier, deskripsiQualifier, targetQualifier := "", "", ""
		if strings.EqualFold(colL, QualifierYa) {
			itemQualifier = colM
			deskripsiQualifier = colN
			targetQualifier = colO
		}

		subRow := dto.PenyusunanKpiSubDetailRow{
			No:                        no,
			KPI:                       colB,
			SubKPI:                    colC,
			Polarisasi:                colD,
			Capping:                   strings.TrimSuffix(colE, "%"),
			Bobot:                     bobot,
			Glossary:                  colG,
			TargetTriwulan:            colH,
			TargetKuantitatifTriwulan: targetKuantitatifTriwulan,
			TargetTahunan:             colJ,
			TargetKuantitatifTahunan:  targetKuantitatifTahunan,
			TerdapatQualifier:         colL,
			Qualifier:                 itemQualifier,
			DeskripsiQualifier:        deskripsiQualifier,
			TargetQualifier:           targetQualifier,
			IsTW24:                    isTW4,
			Result:                    excel.NullableString(colP, isExtendedTriwulan),
			DeskripsiResult:           excel.NullableString(colQ, isExtendedTriwulan),
			Process:                   excel.NullableString(colR, isExtendedTriwulan),
			DeskripsiProcess:          excel.NullableString(colS, isExtendedTriwulan),
			Context:                   excel.NullableString(colT, isExtendedTriwulan),
			DeskripsiContext:          excel.NullableString(colU, isExtendedTriwulan),
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	if len(kpiRows) == 0 {
		return nil, nil, fmt.Errorf("file Excel '%s' tidak memiliki data yang valid", file.Filename)
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
