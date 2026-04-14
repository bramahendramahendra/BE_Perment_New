package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	dto "permen_api/domain/realisasi_kpi/dto"
)

// Konstanta format Excel realisasi.
// Header ada di Row 1; data mulai dari Row 2.
const (
	RealisasiDataStartRow = 2
	RealisasiMaxDataRows  = 200

	SheetTW1 = "TW1"
	SheetTW2 = "TW2"
	SheetTW3 = "TW3"
	SheetTW4 = "TW4"
)

// IsExtendedTriwulan returns true untuk TW2 dan TW4.
func IsExtendedTriwulan(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == "TW2" || upper == "TW4"
}

// ParseAndValidateRealisasiExcel membaca file Excel realisasi, memvalidasi, dan mengembalikan
// slice RealisasiKpiRow yang sudah terisi data dari kolom A–M (atau A–Y untuk TW2/TW4).
//
// Kolom template realisasi:
//
//	A=No, B=KPI, C=SubKPI, D=Polarisasi, E=Capping, F=Bobot%,
//	G=TargetTriwulan, H=Qualifier, I=TargetQualifier,
//	J=Realisasi, K=RealisasiKuantitatif, L=RealisasiQualifier, M=RealisasiKuantitatifQualifier
//	TW1/TW3: N=Result, O=DeskripsiResult, P=Process, Q=DeskripsiProcess, R=Context, S=DeskripsiContext
//	TW2/TW4: N=Result, O=DeskripsiResult, P=RealisasiResult, Q=LinkResult,
//	         R=Process, S=DeskripsiProcess, T=RealisasiProcess, U=LinkProcess,
//	         V=Context, W=DeskripsiContext, X=RealisasiContext, Y=LinkContext
func ParseAndValidateRealisasiExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.RealisasiKpiRow, error) {
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

	isTW24 := IsExtendedTriwulan(triwulan)
	targetSheet := strings.ToUpper(strings.TrimSpace(triwulan))

	sheetIndex, err := xlsx.GetSheetIndex(targetSheet)
	if err != nil || sheetIndex < 0 {
		return nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. Pastikan nama sheet sesuai triwulan ('%s', '%s', '%s', atau '%s')",
			file.Filename, targetSheet, SheetTW1, SheetTW2, SheetTW3, SheetTW4,
		)
	}

	allRows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", targetSheet, err)
	}

	// Data mulai dari row 2 (index 1 setelah header)
	dataStartIdx := RealisasiDataStartRow - 1
	if len(allRows) <= dataStartIdx {
		return nil, fmt.Errorf(
			"file Excel '%s' sheet '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, targetSheet, RealisasiDataStartRow,
		)
	}

	dataEndIdx := dataStartIdx + RealisasiMaxDataRows
	if dataEndIdx > len(allRows) {
		dataEndIdx = len(allRows)
	}
	limitedRows := allRows[dataStartIdx:dataEndIdx]

	// Jumlah kolom minimum yang dibutuhkan
	expectedCols := 13 // A–M
	if isTW24 {
		expectedCols = 25 // A–Y
	} else {
		expectedCols = 19 // A–S (TW1/TW3 extended A–S)
	}

	var rows []dto.RealisasiKpiRow

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
		colJ := strings.TrimSpace(row[9])  // Realisasi
		colK := strings.TrimSpace(row[10]) // Realisasi Kuantitatif
		colL := strings.TrimSpace(row[11]) // Realisasi Qualifier
		colM := strings.TrimSpace(row[12]) // Realisasi Kuantitatif Qualifier

		// Lewati baris kosong
		if colA == "" && colB == "" && colC == "" {
			continue
		}

		// Kolom A: No (harus angka)
		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, fmt.Errorf("baris %d, Kolom A (No): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		// Kolom B: KPI — wajib
		if colB == "" {
			return nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}

		// Kolom C: Sub KPI — wajib
		if colC == "" {
			return nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		// Kolom D: Polarisasi — wajib
		if colD == "" {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}

		// Kolom E: Capping — wajib
		if colE == "" {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}

		// Kolom F: Bobot — wajib angka
		bobot, errBobot := parseFloat2Decimal(colF)
		if errBobot != nil {
			return nil, fmt.Errorf(
				"baris %d, Kolom F (Bobot %%): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colF,
			)
		}

		// Kolom K: Realisasi Kuantitatif — wajib angka
		realisasiKuantitatif, errK := parseFloat2Decimal(colK)
		if errK != nil {
			return nil, fmt.Errorf(
				"baris %d, Kolom K (Realisasi Kuantitatif): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colK,
			)
		}

		// Kolom M: Realisasi Kuantitatif Qualifier — boleh kosong (berarti 0)
		realisasiKuantitatifQualifier := 0.0
		if colM != "" {
			realisasiKuantitatifQualifier, err = parseFloat2Decimal(colM)
			if err != nil {
				return nil, fmt.Errorf(
					"baris %d, Kolom M (Realisasi Kuantitatif Qualifier): harus berupa angka jika diisi, nilai saat ini: '%s'",
					displayRow, colM,
				)
			}
		}

		r := dto.RealisasiKpiRow{
			RowIndex:                      rowIdx,
			No:                            no,
			KPI:                           colB,
			SubKPI:                        colC,
			Polarisasi:                    colD,
			Capping:                       colE,
			Bobot:                         bobot,
			TargetTriwulan:                colG,
			Qualifier:                     colH,
			TargetQualifier:               colI,
			Realisasi:                     colJ,
			RealisasiKuantitatif:          realisasiKuantitatif,
			RealisasiQualifierVal:         colL,
			RealisasiKuantitatifQualifier: realisasiKuantitatifQualifier,
			IsTW24:                        isTW24,
		}

		// Kolom extended
		if isTW24 && len(row) >= 25 {
			// TW2/TW4: N=Result, O=DeskripsiResult, P=RealisasiResult, Q=LinkResult,
			//          R=Process, S=DeskripsiProcess, T=RealisasiProcess, U=LinkProcess,
			//          V=Context, W=DeskripsiContext, X=RealisasiContext, Y=LinkContext
			colN := strings.TrimSpace(row[13])
			colO := strings.TrimSpace(row[14])
			colP := strings.TrimSpace(row[15])
			colQ := strings.TrimSpace(row[16])
			colR := strings.TrimSpace(row[17])
			colS := strings.TrimSpace(row[18])
			colT := strings.TrimSpace(row[19])
			colU := strings.TrimSpace(row[20])
			colV := strings.TrimSpace(row[21])
			colW := strings.TrimSpace(row[22])
			colX := strings.TrimSpace(row[23])
			colY := strings.TrimSpace(row[24])

			r.Result = &colN
			r.DeskripsiResult = &colO
			r.RealisasiResult = &colP
			r.LinkResult = &colQ
			r.Process = &colR
			r.DeskripsiProcess = &colS
			r.RealisasiProcess = &colT
			r.LinkProcess = &colU
			r.Context = &colV
			r.DeskripsiContext = &colW
			r.RealisasiContext = &colX
			r.LinkContext = &colY
		} else if !isTW24 && len(row) >= 19 {
			// TW1/TW3: N=Result, O=DeskripsiResult, P=Process, Q=DeskripsiProcess, R=Context, S=DeskripsiContext
			colN := strings.TrimSpace(row[13])
			colO := strings.TrimSpace(row[14])
			colP := strings.TrimSpace(row[15])
			colQ := strings.TrimSpace(row[16])
			colR := strings.TrimSpace(row[17])
			colS := strings.TrimSpace(row[18])

			r.Result = &colN
			r.DeskripsiResult = &colO
			r.Process = &colP
			r.DeskripsiProcess = &colQ
			r.Context = &colR
			r.DeskripsiContext = &colS
		}

		rows = append(rows, r)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("file Excel '%s' sheet '%s' tidak memiliki data yang valid", file.Filename, targetSheet)
	}

	return rows, nil
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
