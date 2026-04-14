package excel

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
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
// Jika tidak di-set atau tidak valid, mengembalikan MaxDataRows (13).
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

// IsExtendedTriwulan mengembalikan true jika triwulan adalah TW2 atau TW4,
// yaitu triwulan yang memerlukan kolom extended (P–U).
func IsExtendedTriwulan(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == strings.ToUpper(TriwulanTW2) || upper == strings.ToUpper(TriwulanTW4)
}

// ParseAndValidateExcel membaca file Excel, memvalidasi isi, dan mengembalikan:
//   - kpiRows       : slice KPI unik dari kolom B (urutan kemunculan pertama)
//   - kpiSubDetails : map[kpiIndex] -> []KpiSubDetailRow
func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]KpiRow, map[int][]KpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseInternal(file, triwulan, maxRows)
}

func parseInternal(
	file *multipart.FileHeader,
	triwulan string,
	maxRows int,
) ([]KpiRow, map[int][]KpiSubDetailRow, error) {
	if maxRows <= 0 {
		return nil, nil, fmt.Errorf("maxRows harus lebih dari 0, nilai saat ini: %d", maxRows)
	}

	src, err := file.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membuka file Excel '%s': %w", file.Filename, err)
	}
	defer src.Close()

	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membaca file Excel '%s': %w", file.Filename, err)
	}
	defer xlsx.Close()

	isExtendedTriwulan := IsExtendedTriwulan(triwulan)
	isTW4 := strings.EqualFold(triwulan, TriwulanTW4)

	targetSheet := strings.ToUpper(strings.TrimSpace(triwulan))

	sheetIndex, err := xlsx.GetSheetIndex(targetSheet)
	if err != nil || sheetIndex < 0 {
		return nil, nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. Pastikan nama sheet di file sesuai dengan triwulan ('%s', '%s', '%s', atau '%s')",
			file.Filename, targetSheet, SheetTW1, SheetTW2, SheetTW3, SheetTW4,
		)
	}

	allRows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", targetSheet, err)
	}

	if len(allRows) < DataStartRow {
		return nil, nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, DataStartRow,
		)
	}

	dataStartIdx := DataStartRow - 1
	dataEndIdx := dataStartIdx + maxRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}
	limitedRows := allRows[dataStartIdx:dataEndIdx]

	kpiIndexMap := make(map[string]int)
	kpiRows := []KpiRow{}
	kpiSubDetails := make(map[int][]KpiSubDetailRow)
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
		if isExtendedTriwulan {
			colP = strings.TrimSpace(row[15])
			colQ = strings.TrimSpace(row[16])
			colR = strings.TrimSpace(row[17])
			colS = strings.TrimSpace(row[18])
			colT = strings.TrimSpace(row[19])
			colU = strings.TrimSpace(row[20])
		}

		// Lewati baris kosong
		if colA == "" && colB == "" && colC == "" {
			continue
		}

		// Kolom A: No (harus angka)
		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, nil, fmt.Errorf("baris %d, Kolom A (NO): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		// Kolom B: KPI — registrasi sebagai KPI unik jika belum ada
		if colB == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}
		kpiKey := strings.ToLower(strings.TrimSpace(colB))
		kpiIdx, found := kpiIndexMap[kpiKey]
		if !found {
			kpiIdx = len(kpiRows)
			kpiIndexMap[kpiKey] = kpiIdx
			kpiRows = append(kpiRows, KpiRow{
				KpiIndex: kpiIdx,
				IdKpi:    "",
				Kpi:      colB,
				Rumus:    "",
			})
		}

		// Kolom C: Sub KPI
		if colC == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		// Kolom D: Polarisasi
		if colD == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != PolarisasiMaximize && colD != PolarisasiMinimize {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom D (Polarisasi): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colD, PolarisasiMaximize, PolarisasiMinimize,
			)
		}

		// Kolom E: Capping
		if colE == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != CappingOption1 && colE != CappingOption2 {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom E (Capping): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colE, CappingOption1, CappingOption2,
			)
		}

		// Kolom F: Bobot — akumulasikan ke totalBobot
		if colF == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		bobot, errBobot := parseFloat2Decimal(colF)
		if errBobot != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom F (Bobot %%): harus berupa angka 2 desimal tanpa simbol persen, nilai saat ini: '%s'",
				displayRow, colF,
			)
		}
		totalBobot += bobot

		// Kolom G: Glossary
		if colG == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom G (Glossary): tidak boleh kosong", displayRow)
		}

		// Kolom H: Target Triwulanan
		if colH == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom H (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		// Kolom I: Target Kuantitatif Triwulanan
		targetKuantitatifTriwulan, errI := parseFloat2Decimal(colI)
		if errI != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom I (Target Kuantitatif Triwulanan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colI,
			)
		}

		// Kolom J: Target Tahunan
		if colJ == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom J (Target Tahunan): tidak boleh kosong", displayRow)
		}

		// Kolom K: Target Kuantitatif Tahunan
		targetKuantitatifTahunan, errK := parseFloat2Decimal(colK)
		if errK != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom K (Target Kuantitatif Tahunan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colK,
			)
		}

		// Kolom L: Terdapat Qualifier
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

		// Kolom P-U: hanya divalidasi pada TW2 dan TW4
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

		subRow := KpiSubDetailRow{
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
			Result:                    nullableString(colP, isExtendedTriwulan),
			DeskripsiResult:           nullableString(colQ, isExtendedTriwulan),
			Process:                   nullableString(colR, isExtendedTriwulan),
			DeskripsiProcess:          nullableString(colS, isExtendedTriwulan),
			Context:                   nullableString(colT, isExtendedTriwulan),
			DeskripsiContext:          nullableString(colU, isExtendedTriwulan),
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	// Validasi: minimal harus ada 1 KPI
	if len(kpiRows) == 0 {
		return nil, nil, fmt.Errorf("file Excel '%s' tidak memiliki data yang valid", file.Filename)
	}

	// Validasi: setiap KPI unik harus memiliki minimal 1 sub KPI
	for _, kpiRow := range kpiRows {
		if _, ok := kpiSubDetails[kpiRow.KpiIndex]; !ok {
			return nil, nil, fmt.Errorf(
				"KPI '%s' tidak memiliki data sub KPI di file Excel '%s'",
				kpiRow.Kpi, file.Filename,
			)
		}
	}

	// Validasi bobot: akumulasi TOTAL semua baris harus tepat 100%
	roundedTotal := math.Round(totalBobot*100) / 100
	if math.Abs(roundedTotal-TotalBobotExpected) > BobotTolerance {
		return nil, nil, fmt.Errorf(
			"total Bobot (Kolom F) semua KPI = %.2f%%, harus tepat 100%%",
			roundedTotal,
		)
	}

	return kpiRows, kpiSubDetails, nil
}

// nullableString mengembalikan pointer string jika isActive true, nil jika false.
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
