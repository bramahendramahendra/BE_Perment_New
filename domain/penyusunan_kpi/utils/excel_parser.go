package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"

	"github.com/xuri/excelize/v2"
)

const (
	ExcelDataStartRow = 3
	ExcelMaxDataRows  = 13

	SheetTW4       = "TW 4"
	SheetSelainTW4 = "Selain TW 4"

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

func GetMaxRowsFromEnv() int {
	val := os.Getenv("EXCEL_MAX_ROWS")
	if val == "" {
		return ExcelMaxDataRows
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return ExcelMaxDataRows
	}
	return n
}

func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
	kpiList []dto.PenyusunanKpiDetailRequest,
) (map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, kpiList, maxRows)
}

func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	kpiList []dto.PenyusunanKpiDetailRequest,
	maxRows int,
) (map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	if maxRows <= 0 {
		return nil, fmt.Errorf("maxRows harus lebih dari 0, nilai saat ini: %d", maxRows)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file Excel '%s': %w", file.Filename, err)
	}
	defer src.Close()

	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca file Excel '%s': %w", file.Filename, err)
	}
	defer xlsx.Close()

	isTW4 := strings.EqualFold(triwulan, TriwulanTW4)
	targetSheet := SheetSelainTW4
	if isTW4 {
		targetSheet = SheetTW4
	}

	sheetIndex, err := xlsx.GetSheetIndex(targetSheet)
	if err != nil || sheetIndex < 0 {
		return nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. Pastikan file memiliki sheet '%s' dan '%s'",
			file.Filename, targetSheet, SheetTW4, SheetSelainTW4,
		)
	}

	allRows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", targetSheet, err)
	}

	if len(allRows) < ExcelDataStartRow {
		return nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, ExcelDataStartRow,
		)
	}

	dataStartIdx := ExcelDataStartRow - 1
	dataEndIdx := dataStartIdx + maxRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}
	limitedRows := allRows[dataStartIdx:dataEndIdx]

	kpiIndexMap := make(map[string]int, len(kpiList))
	for i, kpiItem := range kpiList {
		kpiIndexMap[strings.ToLower(strings.TrimSpace(kpiItem.Kpi))] = i
	}

	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)
	bobotPerKpi := make(map[int]float64)

	expectedCols := 15
	if isTW4 {
		expectedCols = 21
	}

	for rowIdx, row := range limitedRows {
		displayRow := dataStartIdx + rowIdx + 1

		for len(row) < expectedCols {
			row = append(row, "")
		}

		colA := strings.TrimSpace(row[0])
		colB := strings.TrimSpace(row[1])
		colC := strings.TrimSpace(row[2])
		colD := strings.TrimSpace(row[3])
		colE := strings.TrimSpace(row[4])
		colF := strings.TrimSpace(row[5])
		colG := strings.TrimSpace(row[6])
		colH := strings.TrimSpace(row[7])
		colI := strings.TrimSpace(row[8])
		colJ := strings.TrimSpace(row[9])
		colK := strings.TrimSpace(row[10])
		colL := strings.TrimSpace(row[11])
		colM := strings.TrimSpace(row[12])
		colN := strings.TrimSpace(row[13])
		colO := strings.TrimSpace(row[14])

		var colP, colQ, colR, colS, colT, colU string
		if isTW4 {
			colP = strings.TrimSpace(row[15])
			colQ = strings.TrimSpace(row[16])
			colR = strings.TrimSpace(row[17])
			colS = strings.TrimSpace(row[18])
			colT = strings.TrimSpace(row[19])
			colU = strings.TrimSpace(row[20])
		}

		if colA == "" && colB == "" && colC == "" {
			continue
		}

		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, fmt.Errorf("baris %d, Kolom A (NO): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		if colB == "" {
			return nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}
		kpiIdx, found := kpiIndexMap[strings.ToLower(colB)]
		if !found {
			kpiNames := make([]string, len(kpiList))
			for i, k := range kpiList {
				kpiNames[i] = "'" + k.Kpi + "'"
			}
			return nil, fmt.Errorf(
				"baris %d, Kolom B (KPI): nilai '%s' tidak cocok dengan KPI manapun di REQUEST. KPI yang valid: %s",
				displayRow, colB, strings.Join(kpiNames, ", "),
			)
		}

		if colC == "" {
			return nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		if colD == "" {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != PolarisasiMaximize && colD != PolarisasiMinimize {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): nilai '%s' tidak valid. Gunakan '%s' atau '%s'", displayRow, colD, PolarisasiMaximize, PolarisasiMinimize)
		}

		if colE == "" {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != CappingOption1 && colE != CappingOption2 {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): nilai '%s' tidak valid. Gunakan '%s' atau '%s'", displayRow, colE, CappingOption1, CappingOption2)
		}

		if colF == "" {
			return nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		bobot, errBobot := ParseFloat2Decimal(colF)
		if errBobot != nil {
			return nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): harus berupa angka 2 desimal tanpa simbol persen, nilai saat ini: '%s'", displayRow, colF)
		}
		bobotPerKpi[kpiIdx] += bobot

		if colG == "" {
			return nil, fmt.Errorf("baris %d, Kolom G (Glossary): tidak boleh kosong", displayRow)
		}

		if colH == "" {
			return nil, fmt.Errorf("baris %d, Kolom H (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		targetKuantitatifTriwulan, errI := ParseFloat2Decimal(colI)
		if errI != nil {
			return nil, fmt.Errorf("baris %d, Kolom I (Target Kuantitatif Triwulanan): harus berupa angka 2 desimal, nilai saat ini: '%s'", displayRow, colI)
		}

		if colJ == "" {
			return nil, fmt.Errorf("baris %d, Kolom J (Target Tahunan): tidak boleh kosong", displayRow)
		}

		targetKuantitatifTahunan, errK := ParseFloat2Decimal(colK)
		if errK != nil {
			return nil, fmt.Errorf("baris %d, Kolom K (Target Kuantitatif Tahunan): harus berupa angka 2 desimal, nilai saat ini: '%s'", displayRow, colK)
		}

		if colL == "" {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): tidak boleh kosong", displayRow)
		}

		// if colL != QualifierYa && colL != QualifierTidak {
		if !strings.EqualFold(colL, QualifierYa) && !strings.EqualFold(colL, QualifierTidak) {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): nilai '%s' tidak valid. Gunakan '%s' atau '%s'", displayRow, colL, QualifierYa, QualifierTidak)
		}

		if strings.EqualFold(colL, QualifierYa) {
			if colM == "" {
				return nil, fmt.Errorf("baris %d, Kolom M (Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
			if colN == "" {
				return nil, fmt.Errorf("baris %d, Kolom N (Deskripsi Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
			if colO == "" {
				return nil, fmt.Errorf("baris %d, Kolom O (Target Qualifier): tidak boleh kosong jika Terdapat Qualifier = 'Ya'", displayRow)
			}
		}

		if isTW4 {
			if colP == "" {
				return nil, fmt.Errorf("baris %d, Kolom P (Result): tidak boleh kosong", displayRow)
			}
			if colQ == "" {
				return nil, fmt.Errorf("baris %d, Kolom Q (Deskripsi Result): tidak boleh kosong", displayRow)
			}
			if colR == "" {
				return nil, fmt.Errorf("baris %d, Kolom R (Process): tidak boleh kosong", displayRow)
			}
			if colS == "" {
				return nil, fmt.Errorf("baris %d, Kolom S (Deskripsi Process): tidak boleh kosong", displayRow)
			}
			if colT == "" {
				return nil, fmt.Errorf("baris %d, Kolom T (Context): tidak boleh kosong", displayRow)
			}
			if colU == "" {
				return nil, fmt.Errorf("baris %d, Kolom U (Deskripsi Context): tidak boleh kosong", displayRow)
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
			Capping:                   colE,
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
			IsTW4:                     isTW4,
			Result:                    NullableString(colP, isTW4),
			DeskripsiResult:           NullableString(colQ, isTW4),
			Process:                   NullableString(colR, isTW4),
			DeskripsiProcess:          NullableString(colS, isTW4),
			Context:                   NullableString(colT, isTW4),
			DeskripsiContext:          NullableString(colU, isTW4),
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	for i, kpiItem := range kpiList {
		if _, ok := kpiSubDetails[i]; !ok {
			return nil, fmt.Errorf(
				"KPI '%s' (index %d) tidak memiliki data sub KPI di file Excel '%s'. Pastikan kolom B berisi nama KPI yang sesuai dengan REQUEST",
				kpiItem.Kpi, i+1, file.Filename,
			)
		}
	}

	for i, kpiItem := range kpiList {
		totalBobot := math.Round(bobotPerKpi[i]*100) / 100
		if math.Abs(totalBobot-TotalBobotExpected) > BobotTolerance {
			return nil, fmt.Errorf(
				"KPI '%s': total Bobot (Kolom F) = %.2f%%, harus tepat 100%%",
				kpiItem.Kpi, totalBobot,
			)
		}
	}

	if len(kpiSubDetails) == 0 {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki data yang valid", file.Filename)
	}

	return kpiSubDetails, nil
}
