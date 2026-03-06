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

// =============================================
// KONSTANTA VALIDASI
// =============================================

const (
	excelDataStartRow = 3 // Baris 2 = header, data mulai baris 3

	// ExcelMaxDataRows adalah nilai default batas baris jika EXCEL_MAX_ROWS tidak di-set di .env
	ExcelMaxDataRows = 13

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

// =============================================
// HELPER: BACA maxRows DARI ENV
// =============================================

// getMaxRowsFromEnv membaca batas baris Excel dari environment variable EXCEL_MAX_ROWS.
// Jika tidak di-set atau tidak valid, gunakan nilai default ExcelMaxDataRows.
func getMaxRowsFromEnv() int {
	val := os.Getenv("EXCEL_MAX_ROWS")
	if val == "" {
		return ExcelMaxDataRows
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		fmt.Printf("[WARN] EXCEL_MAX_ROWS='%s' tidak valid, menggunakan default %d\n", val, ExcelMaxDataRows)
		return ExcelMaxDataRows
	}
	return n
}

// =============================================
// MAIN PARSER FUNCTION
// =============================================

// ParseAndValidateExcel membaca 1 file Excel dengan 2 sheet ("TW 4" / "Selain TW 4"),
// menentukan sheet berdasarkan triwulan, lalu memetakan setiap baris ke KPI
// berdasarkan nilai kolom B (case-insensitive).
//
// Parameter:
//   - file      : file Excel (.xlsx) dari multipart form
//   - triwulan  : nilai Triwulan dari REQUEST (misal "TW4", "TW1", "TW2", "TW3")
//   - kpiList   : daftar KPI dari req.Kpi, digunakan untuk mapping kolom B
//
// Return:
//   - map[int][]dto.PenyusunanKpiSubDetailRow : key = index KPI di kpiList, value = baris sub KPI
//   - error
//
// Aturan sheet:
//   - Triwulan == "TW4" → gunakan sheet "TW 4"    (kolom A–U, 21 kolom)
//   - Selain itu        → gunakan sheet "Selain TW 4" (kolom A–O, 15 kolom)
//
// Aturan mapping KPI:
//   - Kolom B berisi nama KPI (free text)
//   - Dicocokkan case-insensitive ke kpiList[i].Kpi
//   - Jika kolom B tidak cocok dengan KPI manapun → return error
//
// Aturan bobot:
//   - Total bobot (kolom F) dihitung PER KPI, masing-masing harus = 100%
func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
	kpiList []dto.PenyusunanKpiDetailItem,
) (map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	maxRows := getMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, kpiList, maxRows)
}

// parseAndValidateExcelInternal adalah implementasi utama parser Excel.
func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	kpiList []dto.PenyusunanKpiDetailItem,
	maxRows int,
) (map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	if maxRows <= 0 {
		return nil, fmt.Errorf("maxRows harus lebih dari 0, nilai saat ini: %d", maxRows)
	}

	// --- Buka file dari memory ---
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

	// --- Tentukan nama sheet berdasarkan triwulan ---
	isTW4 := strings.EqualFold(triwulan, triwulanTW4)
	var sheetName string
	if isTW4 {
		sheetName = sheetTW4
	} else {
		sheetName = sheetSelainTW4
	}

	// --- Validasi sheet ada di file Excel ---
	sheetIndex, err := xlsx.GetSheetIndex(sheetName)
	if err != nil || sheetIndex < 0 {
		return nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. "+
				"Pastikan file memiliki sheet '%s' dan '%s'",
			file.Filename, sheetName, sheetTW4, sheetSelainTW4,
		)
	}

	allRows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", sheetName, err)
	}

	// Pastikan ada data mulai baris ke-3 (index 2)
	if len(allRows) < excelDataStartRow {
		return nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, sheetName, excelDataStartRow,
		)
	}

	// --- [DEBUG] Log info sheet & baris ---
	fmt.Printf("[DEBUG] File: '%s' | Sheet: '%s' | Total baris: %d | maxRows: %d\n",
		file.Filename, sheetName, len(allRows), maxRows)

	// --- Tentukan batas baris yang akan dibaca ---
	dataStartIdx := excelDataStartRow - 1
	dataEndIdx := dataStartIdx + maxRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}

	limitedRows := allRows[dataStartIdx:dataEndIdx]
	skippedRows := (len(allRows) - dataStartIdx) - len(limitedRows)

	// --- Buat map index KPI (lowercase nama KPI → index di kpiList) ---
	// Untuk lookup O(1) saat mapping kolom B
	kpiIndexMap := make(map[string]int, len(kpiList))
	for i, kpiItem := range kpiList {
		kpiIndexMap[strings.ToLower(strings.TrimSpace(kpiItem.Kpi))] = i
	}

	// --- Hasil: map index KPI → slice baris sub KPI ---
	kpiSubDetails := make(map[int][]dto.PenyusunanKpiSubDetailRow)

	// --- Akumulasi bobot per KPI untuk validasi total 100% ---
	bobotPerKpi := make(map[int]float64)

	// --- Tentukan jumlah kolom yang diharapkan berdasarkan sheet ---
	// TW 4: 21 kolom (A–U), Selain TW 4: 15 kolom (A–O)
	expectedCols := 15
	if isTW4 {
		expectedCols = 21
	}

	// --- Loop setiap baris data ---
	for rowIdx, row := range limitedRows {
		displayRow := dataStartIdx + rowIdx + 1 // nomor baris 1-based untuk pesan error

		// Padding jika kolom kurang dari jumlah yang diharapkan
		for len(row) < expectedCols {
			row = append(row, "")
		}

		// Ambil kolom A–O (selalu ada di kedua sheet)
		colA := strings.TrimSpace(row[0])  // NO
		colB := strings.TrimSpace(row[1])  // KPI (text) — kunci mapping
		colC := strings.TrimSpace(row[2])  // Sub KPI
		colD := strings.TrimSpace(row[3])  // Polarisasi
		colE := strings.TrimSpace(row[4])  // Capping
		colF := strings.TrimSpace(row[5])  // Bobot %
		colG := strings.TrimSpace(row[6])  // Glossary
		colH := strings.TrimSpace(row[7])  // Target Triwulanan
		colI := strings.TrimSpace(row[8])  // Target Kuantitatif Triwulanan
		colJ := strings.TrimSpace(row[9])  // Target Tahunan
		colK := strings.TrimSpace(row[10]) // Target Kuantitatif Tahunan
		colL := strings.TrimSpace(row[11]) // Terdapat Qualifier
		colM := strings.TrimSpace(row[12]) // Qualifier
		colN := strings.TrimSpace(row[13]) // Deskripsi Qualifier
		colO := strings.TrimSpace(row[14]) // Target Qualifier

		// Kolom P–U hanya tersedia di sheet "TW 4"
		var colP, colQ, colR, colS, colT, colU string
		if isTW4 {
			colP = strings.TrimSpace(row[15]) // Result
			colQ = strings.TrimSpace(row[16]) // Deskripsi Result
			colR = strings.TrimSpace(row[17]) // Process
			colS = strings.TrimSpace(row[18]) // Deskripsi Process
			colT = strings.TrimSpace(row[19]) // Context
			colU = strings.TrimSpace(row[20]) // Deskripsi Context
		}
		// Untuk sheet "Selain TW 4", colP–colU tetap "" (akan disimpan NULL di DB)

		// Lewati baris kosong
		if colA == "" && colB == "" && colC == "" {
			continue
		}

		// =============================================
		// VALIDASI KOLOM A: NO (angka)
		// =============================================
		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, fmt.Errorf("baris %d, Kolom A (NO): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colA)
		}

		// =============================================
		// MAPPING KOLOM B → INDEX KPI (case-insensitive)
		// =============================================
		if colB == "" {
			return nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}
		kpiIdx, found := kpiIndexMap[strings.ToLower(colB)]
		if !found {
			// Kolom B tidak cocok dengan KPI manapun di request → error
			kpiNames := make([]string, len(kpiList))
			for i, k := range kpiList {
				kpiNames[i] = "'" + k.Kpi + "'"
			}
			return nil, fmt.Errorf(
				"baris %d, Kolom B (KPI): nilai '%s' tidak cocok dengan KPI manapun di REQUEST. "+
					"KPI yang valid: %s",
				displayRow, colB, strings.Join(kpiNames, ", "),
			)
		}

		// =============================================
		// VALIDASI KOLOM C: Sub KPI
		// =============================================
		if colC == "" {
			return nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		// =============================================
		// VALIDASI KOLOM D: Polarisasi
		// =============================================
		if colD == "" {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != polarisasiMaximize && colD != polarisasiMinimize {
			return nil, fmt.Errorf(
				"baris %d, Kolom D (Polarisasi): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colD, polarisasiMaximize, polarisasiMinimize)
		}

		// =============================================
		// VALIDASI KOLOM E: Capping
		// =============================================
		if colE == "" {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != cappingOption1 && colE != cappingOption2 {
			return nil, fmt.Errorf(
				"baris %d, Kolom E (Capping): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colE, cappingOption1, cappingOption2)
		}

		// =============================================
		// VALIDASI KOLOM F: Bobot % (akumulasi per KPI)
		// =============================================
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

		// =============================================
		// VALIDASI KOLOM G–K
		// =============================================
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

		// =============================================
		// VALIDASI KOLOM L: Terdapat Qualifier
		// =============================================
		if colL == "" {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): tidak boleh kosong", displayRow)
		}
		if colL != qualifierYa && colL != qualifierTidak {
			return nil, fmt.Errorf(
				"baris %d, Kolom L (Terdapat Qualifier): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colL, qualifierYa, qualifierTidak)
		}

		// =============================================
		// VALIDASI KOLOM M, N, O: Wajib jika Qualifier = "Ya"
		// =============================================
		if strings.EqualFold(colL, qualifierYa) {
			if colM == "" {
				return nil, fmt.Errorf(
					"baris %d, Kolom M (Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
			}
			if colN == "" {
				return nil, fmt.Errorf(
					"baris %d, Kolom N (Deskripsi Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
			}
			if colO == "" {
				return nil, fmt.Errorf(
					"baris %d, Kolom O (Target Qualifier): tidak boleh kosong jika Kolom L = 'Ya'", displayRow)
			}
		}

		// =============================================
		// VALIDASI KOLOM P–U: Hanya untuk sheet "TW 4"
		// =============================================
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

		// =============================================
		// SEMUA VALIDASI LOLOS — tambahkan ke map berdasarkan index KPI
		// =============================================

		// Qualifier: hanya diisi jika Terdapat Qualifier = "Ya"
		itemQualifier := ""
		deskripsiQualifier := ""
		targetQualifier := ""
		if strings.EqualFold(colL, qualifierYa) {
			itemQualifier = colM
			deskripsiQualifier = colN
			targetQualifier = colO
		}

		// Untuk sheet "Selain TW 4", kolom P–U adalah "" → disimpan NULL di DB
		// (nil pada *string akan di-handle repo, cukup kirim "" dan biarkan repo konversi)
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
			// Kolom P–U: diisi hanya jika TW4, selain itu nil (NULL di DB)
			Result:           nullableString(colP, isTW4),
			DeskripsiResult:  nullableString(colQ, isTW4),
			Process:          nullableString(colR, isTW4),
			DeskripsiProcess: nullableString(colS, isTW4),
			Context:          nullableString(colT, isTW4),
			DeskripsiContext: nullableString(colU, isTW4),
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	// --- Peringatan jika ada baris yang dipotong ---
	if skippedRows > 0 {
		fmt.Printf("[WARN] file '%s': %d baris melebihi batas maksimal (%d baris), diabaikan\n",
			file.Filename, skippedRows, maxRows)
	}

	// --- Validasi: pastikan setiap KPI di request punya minimal 1 baris sub KPI ---
	for i, kpiItem := range kpiList {
		if _, ok := kpiSubDetails[i]; !ok {
			return nil, fmt.Errorf(
				"KPI '%s' (index %d) tidak memiliki data sub KPI di file Excel '%s'. "+
					"Pastikan kolom B berisi nama KPI yang sesuai dengan REQUEST",
				kpiItem.Kpi, i+1, file.Filename,
			)
		}
	}

	// --- Validasi total bobot per KPI harus = 100% ---
	for i, kpiItem := range kpiList {
		totalBobot := bobotPerKpi[i]
		totalBobotRounded := math.Round(totalBobot*100) / 100
		if math.Abs(totalBobotRounded-totalBobotExpected) > bobotTolerance {
			return nil, fmt.Errorf(
				"KPI '%s': total Bobot (Kolom F) = %.2f%%, harus tepat 100%%",
				kpiItem.Kpi, totalBobotRounded,
			)
		}
	}

	// --- Validasi result tidak kosong ---
	if len(kpiSubDetails) == 0 {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki data yang valid", file.Filename)
	}

	return kpiSubDetails, nil
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// nullableString mengembalikan pointer string jika isActive true, atau nil jika false.
// Digunakan untuk kolom P–U yang hanya ada di sheet "TW 4".
// Nilai nil akan disimpan sebagai NULL di DB.
func nullableString(val string, isActive bool) *string {
	if !isActive {
		return nil
	}
	return &val
}

// parseFloat2Decimal mem-parse string menjadi float64 dengan 2 angka di belakang koma.
func parseFloat2Decimal(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	cleaned := strings.ReplaceAll(s, "%", "")
	cleaned = strings.TrimSpace(cleaned)
	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' bukan angka valid", s)
	}
	rounded := math.Round(val*100) / 100
	return rounded, nil
}
