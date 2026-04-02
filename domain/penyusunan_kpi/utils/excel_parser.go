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

	// Nama sheet mengikuti nilai triwulan dari template (TW1, TW2, TW3, TW4).
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

// IsTriwulanWithChallengeMethod mengembalikan true jika triwulan adalah TW2 atau TW4,
// yaitu triwulan yang memerlukan insert ChallengeList dan MethodList ke DB.
func IsTriwulanWithChallengeMethod(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == strings.ToUpper(TriwulanTW2) || upper == strings.ToUpper(TriwulanTW4)
}

// ParseAndValidateExcel membaca file Excel, memvalidasi isi, dan mengembalikan:
//   - kpiRows     : slice KPI unik dari kolom B (urutan kemunculan pertama)
//   - kpiSubDetails : map[kpiIndex] -> []PenyusunanKpiSubDetailRow
//
// Perubahan dari versi sebelumnya:
//  1. Tidak lagi menerima kpiList dari REQUEST — KPI diambil langsung dari kolom B Excel.
//  2. Validasi bobot: total akumulasi SEMUA baris kolom F harus tepat 100% (bukan per KPI).
//  3. Kolom R,S,T,U (Process/Deskripsi Process/Context/Deskripsi Context) hanya diparse
//     pada TW2 dan TW4 (bukan hanya TW4 seperti sebelumnya).
func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.PenyusunanKpiRow, map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, maxRows)
}

func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	maxRows int,
) ([]dto.PenyusunanKpiRow, map[int][]dto.PenyusunanKpiSubDetailRow, error) {
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

	// Tentukan sheet berdasarkan triwulan.
	// Nama sheet mengikuti nilai triwulan: TW1, TW2, TW3, atau TW4.
	// TW2 dan TW4 menggunakan kolom extended (R,S,T,U tersedia).
	// TW1 dan TW3 menggunakan kolom base (A-O saja).
	isChallengeMethodTriwulan := IsTriwulanWithChallengeMethod(triwulan)
	isTW4 := strings.EqualFold(triwulan, TriwulanTW4)

	// targetSheet = nama sheet yang sesuai dengan triwulan yang dikirim
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

	if len(allRows) < ExcelDataStartRow {
		return nil, nil, fmt.Errorf(
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

	// kpiIndexMap: lowercase nama KPI kolom B → indeks urutan kemunculan pertama
	kpiIndexMap := make(map[string]int)
	// kpiRows: daftar KPI unik sesuai urutan kemunculan di Excel
	kpiRows := []dto.PenyusunanKpiRow{}

	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)

	// totalBobot: akumulasi seluruh bobot semua baris (perubahan utama — tidak lagi per KPI)
	totalBobot := 0.0

	expectedCols := 15
	if isChallengeMethodTriwulan {
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
		if isChallengeMethodTriwulan {
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
			// IdKpi dan Rumus akan diisi oleh service saat lookup mst_kpi.
			// Untuk sementara diisi string asli dari kolom B.
			kpiRows = append(kpiRows, dto.PenyusunanKpiRow{
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

		// Kolom F: Bobot — akumulasikan ke totalBobot (perubahan: tidak lagi per KPI)
		if colF == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		bobot, errBobot := ParseFloat2Decimal(colF)
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
		targetKuantitatifTriwulan, errI := ParseFloat2Decimal(colI)
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
		targetKuantitatifTahunan, errK := ParseFloat2Decimal(colK)
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
		if isChallengeMethodTriwulan {
			// Kolom P,Q (Result) hanya wajib pada TW4
			if isTW4 {
				if colP == "" {
					return nil, nil, fmt.Errorf("baris %d, Kolom P (Result): tidak boleh kosong", displayRow)
				}
				if colQ == "" {
					return nil, nil, fmt.Errorf("baris %d, Kolom Q (Deskripsi Result): tidak boleh kosong", displayRow)
				}
			}
			// Kolom R,S,T,U (Process & Context) wajib pada TW2 dan TW4
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
			// Process dan Context diisi untuk TW2 dan TW4
			Process:          NullableString(colR, isChallengeMethodTriwulan),
			DeskripsiProcess: NullableString(colS, isChallengeMethodTriwulan),
			Context:          NullableString(colT, isChallengeMethodTriwulan),
			DeskripsiContext: NullableString(colU, isChallengeMethodTriwulan),
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
	// (perubahan utama: sebelumnya per KPI, sekarang akumulasi semua baris)
	roundedTotal := math.Round(totalBobot*100) / 100
	if math.Abs(roundedTotal-TotalBobotExpected) > BobotTolerance {
		return nil, nil, fmt.Errorf(
			"total Bobot (Kolom F) semua KPI = %.2f%%, harus tepat 100%%",
			roundedTotal,
		)
	}

	return kpiRows, kpiSubDetails, nil
}
