package service

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
	excelDataStartRow = 3
	ExcelMaxDataRows  = 13

	sheetTW4       = "TW 4"
	sheetSelainTW4 = "Selain TW 4"

	triwulanTW4 = "TW4"

	polarisasiMaximize = "Maximize"
	polarisasiMinimize = "Minimize"

	cappingOption1 = "100%"
	cappingOption2 = "110%"

	qualifierYa    = "Ya"
	qualifierTidak = "Tidak"

	totalBobotExpected = 100.0
	bobotTolerance     = 0.01
)

func getMaxRowsFromEnv() int {
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
	kpiList []dto.PenyusunanKpiDetailItem,
) (map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	maxRows := getMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, kpiList, maxRows)
}

func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	kpiList []dto.PenyusunanKpiDetailItem,
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

	isTW4 := strings.EqualFold(triwulan, triwulanTW4)
	targetSheet := sheetSelainTW4
	if isTW4 {
		targetSheet = sheetTW4
	}

	sheetIndex, err := xlsx.GetSheetIndex(targetSheet)
	if err != nil || sheetIndex < 0 {
		return nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. Pastikan file memiliki sheet '%s' dan '%s'",
			file.Filename, targetSheet, sheetTW4, sheetSelainTW4,
		)
	}

	allRows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", targetSheet, err)
	}

	if len(allRows) < excelDataStartRow {
		return nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, excelDataStartRow,
		)
	}

	dataStartIdx := excelDataStartRow - 1
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
		if colD != polarisasiMaximize && colD != polarisasiMinimize {
			return nil, fmt.Errorf(
				"baris %d, Kolom D (Polarisasi): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colD, polarisasiMaximize, polarisasiMinimize)
		}

		if colE == "" {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != cappingOption1 && colE != cappingOption2 {
			return nil, fmt.Errorf(
				"baris %d, Kolom E (Capping): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colE, cappingOption1, cappingOption2)
		}

		if colF == "" {
			return nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		bobot, errBobot := parseFloat2Decimal(colF)
		if errBobot != nil {
			return nil, fmt.Errorf(
				"baris %d, Kolom F (Bobot %%): harus berupa angka 2 desimal tanpa simbol persen, nilai saat ini: '%s'",
				displayRow, colF)
		}
		bobotPerKpi[kpiIdx] += bobot

		if colG == "" {
			return nil, fmt.Errorf("baris %d, Kolom G (Glossary): tidak boleh kosong", displayRow)
		}
		if colH == "" {
			return nil, fmt.Errorf("baris %d, Kolom H (Target Triwulanan): tidak boleh kosong", displayRow)
		}
		targetKuantitatifTriwulan, errI := parseFloat2Decimal(colI)
		if errI != nil {
			return nil, fmt.Errorf(
				"baris %d, Kolom I (Target Kuantitatif Triwulanan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colI)
		}
		if colJ == "" {
			return nil, fmt.Errorf("baris %d, Kolom J (Target Tahunan): tidak boleh kosong", displayRow)
		}
		targetKuantitatifTahunan, errK := parseFloat2Decimal(colK)
		if errK != nil {
			return nil, fmt.Errorf(
				"baris %d, Kolom K (Target Kuantitatif Tahunan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colK)
		}

		if colL == "" {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): tidak boleh kosong", displayRow)
		}
		if colL != qualifierYa && colL != qualifierTidak {
			return nil, fmt.Errorf(
				"baris %d, Kolom L (Terdapat Qualifier): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colL, qualifierYa, qualifierTidak)
		}

		if strings.EqualFold(colL, qualifierYa) {
			if colM == "" {
				return nil, fmt.Errorf("baris %d, Kolom M (Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
			}
			if colN == "" {
				return nil, fmt.Errorf("baris %d, Kolom N (Deskripsi Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
			}
			if colO == "" {
				return nil, fmt.Errorf("baris %d, Kolom O (Target Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
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
		if strings.EqualFold(colL, qualifierYa) {
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
			Result:                    nullableString(colP, isTW4),
			DeskripsiResult:           nullableString(colQ, isTW4),
			Process:                   nullableString(colR, isTW4),
			DeskripsiProcess:          nullableString(colS, isTW4),
			Context:                   nullableString(colT, isTW4),
			DeskripsiContext:          nullableString(colU, isTW4),
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
		if math.Abs(totalBobot-totalBobotExpected) > bobotTolerance {
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

// nullableString mengembalikan pointer string jika isActive true, nil jika false.
// Nilai nil disimpan sebagai NULL di DB (digunakan untuk kolom P–U sheet "TW 4").
func nullableString(val string, isActive bool) *string {
	if !isActive {
		return nil
	}
	return &val
}

// parseFloat2Decimal mem-parse string menjadi float64 dengan presisi 2 desimal.
func parseFloat2Decimal(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	cleaned := strings.TrimSpace(strings.ReplaceAll(s, "%", ""))
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' bukan angka valid", s)
	}
	return math.Round(val*100) / 100, nil
}
