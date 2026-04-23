package utils

import (
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	dto "permen_api/domain/realisasi_kpi/dto"
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

// IsExtendedTriwulan returns true untuk TW2 dan TW4.
func IsExtendedTriwulan(triwulan string) bool {
	upper := strings.ToUpper(strings.TrimSpace(triwulan))
	return upper == "TW2" || upper == "TW4"
}

func ParseAndValidateRealisasiExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.KpiRow, map[int][]dto.KpiSubDetailRow, error) {
	maxRows := GetMaxRowsFromEnv()
	return parseAndValidateExcelInternal(file, triwulan, maxRows)
}

// ParseAndValidateRealisasiExcel membaca file Excel realisasi, memvalidasi, dan mengembalikan
// slice RealisasiKpiRow yang sudah terisi data dari kolom A–S (TW1/TW3) atau A–Y (TW2/TW4).
//
// Aturan kolom:
//
//	A=No, B=KPI, C=SubKPI, D=Polarisasi, E=Capping, F=Bobot%,
//	G=TargetTriwulan, H=Qualifier (auto-fill "-"), I=TargetQualifier (auto-fill "-"),
//	J=Realisasi, K=RealisasiKuantitatif, L=RealisasiQualifier, M=RealisasiQualifierKuantitatif
//	TW1/TW3: N=Result, O=DeskripsiResult, P=Process, Q=DeskripsiProcess, R=Context, S=DeskripsiContext
//	TW2/TW4: N=Result, O=DeskripsiResult, P=RealisasiResult, Q=LinkResult,
//	         R=Process, S=DeskripsiProcess, T=RealisasiProcess, U=LinkProcess,
//	         V=Context, W=DeskripsiContext, X=RealisasiContext, Y=LinkContext
func parseAndValidateExcelInternal(
	file *multipart.FileHeader,
	triwulan string,
	maxRows int,
) ([]dto.KpiRow, map[int][]dto.KpiSubDetailRow, error) {
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

	isTW24 := IsExtendedTriwulan(triwulan)
	targetSheet := strings.ToUpper(strings.TrimSpace(triwulan))

	sheetIndex, err := xlsx.GetSheetIndex(targetSheet)
	if err != nil || sheetIndex < 0 {
		return nil, nil, fmt.Errorf(
			"file Excel '%s' tidak memiliki sheet '%s'. Pastikan nama sheet sesuai triwulan ('%s', '%s', '%s', atau '%s')",
			file.Filename, targetSheet, SheetTW1, SheetTW2, SheetTW3, SheetTW4,
		)
	}

	allRows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", targetSheet, err)
	}

	// Data mulai dari row 2 (index 1 setelah header)
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

	// kpiIndexMap: lowercase nama KPI kolom B → indeks urutan kemunculan pertama
	kpiIndexMap := make(map[string]int)
	// kpiRows: daftar KPI unik sesuai urutan kemunculan di Excel
	kpiRows := []dto.KpiRow{}
	kpiSubDetails := make(map[int][]dto.KpiSubDetailRow)

	// totalBobot: akumulasi seluruh bobot semua baris (perubahan utama — tidak lagi per KPI)
	totalBobot := 0.0

	// Jumlah kolom minimum yang dibutuhkan
	var expectedCols int
	if isTW24 {
		expectedCols = 25 // A–Y
	} else {
		expectedCols = 19 // A–S
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
		colJ := strings.TrimSpace(row[9])  // Realisasi
		colK := strings.TrimSpace(row[10]) // Realisasi Kuantitatif
		colL := strings.TrimSpace(row[11]) // Realisasi Qualifier
		colM := strings.TrimSpace(row[12]) // Realisasi Qualifier Kuantitatif (free text)

		// Lewati baris kosong
		if colA == "" && colB == "" && colC == "" {
			continue
		}

		// Kolom A: No (harus angka)
		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, nil, fmt.Errorf("baris %d, Kolom A (No): harus berupa angka, nilai saat ini: '%s'", displayRow, colA)
		}

		// Kolom B: KPI — wajib
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
			kpiRows = append(kpiRows, dto.KpiRow{
				KpiIndex: kpiIdx,
				IdKpi:    "",
				Kpi:      colB,
				Rumus:    "",
			})
		}

		// Kolom C: Sub KPI — wajib
		if colC == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		// Kolom D: Polarisasi — dropdown [Maximize, Minimize]
		if colD == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != "Maximize" && colD != "Minimize" {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom D (Polarisasi): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colD, PolarisasiMaximize, PolarisasiMinimize,
			)
		}

		// Kolom E: Capping — dropdown [100%, 110%]
		if colE == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != CappingOption1 && colE != CappingOption2 {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom E (Capping): nilai '%s' tidak valid. Gunakan '%s' atau '%s'",
				displayRow, colE, CappingOption1, CappingOption2,
			)
		}

		// Kolom F: Bobot — wajib angka 2 desimal
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

		// Kolom G: Target Triwulanan — wajib
		if colG == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom G (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		// Kolom J: Realisasi — wajib
		if colJ == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom J (Realisasi): tidak boleh kosong", displayRow)
		}

		// Kolom K: Realisasi Kuantitatif — wajib angka 2 desimal
		realisasiKuantitatif, errK := parseFloat2Decimal(colK)
		if errK != nil {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom K (Realisasi Kuantitatif): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colK,
			)
		}

		// Kolom L: Realisasi Qualifier — dropdown [Ya, Tidak]
		if colL == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom L (Realisasi Qualifier): tidak boleh kosong", displayRow)
		}
		if !strings.EqualFold(colL, QualifierYa) && !strings.EqualFold(colL, QualifierTidak) {
			return nil, nil, fmt.Errorf(
				"baris %d, Kolom L (Realisasi Qualifier): harus 'Ya' atau 'Tidak', nilai saat ini: '%s'",
				displayRow, colL,
			)
		}

		// Kolom M: Realisasi Qualifier Kuantitatif — free text, wajib diisi
		if colM == "" {
			return nil, nil, fmt.Errorf("baris %d, Kolom M (Realisasi Qualifier Kuantitatif): tidak boleh kosong", displayRow)
		}

		subRow := dto.KpiSubDetailRow{
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
			RealisasiQualifierVal:         colL,
			RealisasiKuantitatifQualifier: colM,
			IsTW24:                        isTW24,
		}

		// Kolom extended
		if isTW24 {
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

			if colN == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom N (Result): tidak boleh kosong", displayRow)
			}
			if colO == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom O (Deskripsi Result): tidak boleh kosong", displayRow)
			}
			if colP == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom P (Realisasi Result): tidak boleh kosong", displayRow)
			}
			if colQ == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Q (Link Result): tidak boleh kosong", displayRow)
			}
			if colR == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom R (Process): tidak boleh kosong", displayRow)
			}
			if colS == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom S (Deskripsi Process): tidak boleh kosong", displayRow)
			}
			if colT == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom T (Realisasi Process): tidak boleh kosong", displayRow)
			}
			if colU == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom U (Link Process): tidak boleh kosong", displayRow)
			}
			if colV == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom V (Context): tidak boleh kosong", displayRow)
			}
			if colW == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom W (Deskripsi Context): tidak boleh kosong", displayRow)
			}
			if colX == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom X (Realisasi Context): tidak boleh kosong", displayRow)
			}
			if colY == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Y (Link Context): tidak boleh kosong", displayRow)
			}

			subRow.Result = &colN
			subRow.DeskripsiResult = &colO
			subRow.RealisasiResult = &colP
			subRow.LampiranEvidenceResult = &colQ
			subRow.Process = &colR
			subRow.DeskripsiProcess = &colS
			subRow.RealisasiProcess = &colT
			subRow.LampiranEvidenceProcess = &colU
			subRow.Context = &colV
			subRow.DeskripsiContext = &colW
			subRow.RealisasiContext = &colX
			subRow.LampiranEvidenceContext = &colY
		} else {
			// TW1/TW3: N=Result, O=DeskripsiResult, P=Process, Q=DeskripsiProcess, R=Context, S=DeskripsiContext
			colN := strings.TrimSpace(row[13])
			colO := strings.TrimSpace(row[14])
			colP := strings.TrimSpace(row[15])
			colQ := strings.TrimSpace(row[16])
			colR := strings.TrimSpace(row[17])
			colS := strings.TrimSpace(row[18])

			if colN == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom N (Result): tidak boleh kosong", displayRow)
			}
			if colO == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom O (Deskripsi Result): tidak boleh kosong", displayRow)
			}
			if colP == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom P (Process): tidak boleh kosong", displayRow)
			}
			if colQ == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom Q (Deskripsi Process): tidak boleh kosong", displayRow)
			}
			if colR == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom R (Context): tidak boleh kosong", displayRow)
			}
			if colS == "" {
				return nil, nil, fmt.Errorf("baris %d, Kolom S (Deskripsi Context): tidak boleh kosong", displayRow)
			}

			subRow.Result = &colN
			subRow.DeskripsiResult = &colO
			subRow.Process = &colP
			subRow.DeskripsiProcess = &colQ
			subRow.Context = &colR
			subRow.DeskripsiContext = &colS
		}

		kpiSubDetails[kpiIdx] = append(kpiSubDetails[kpiIdx], subRow)
	}

	if len(kpiRows) == 0 {
		return nil, nil, fmt.Errorf("file Excel '%s' sheet '%s' tidak memiliki data yang valid", file.Filename, targetSheet)
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

	// Validasi total bobot harus 100% (toleransi 0.01)
	roundedTotal := math.Round(totalBobot*100) / 100
	if math.Abs(roundedTotal-TotalBobotExpected) > BobotTolerance {
		return nil, nil, fmt.Errorf(
			"total Bobot (Kolom F) semua KPI = %.2f%%, harus tepat 100%%",
			roundedTotal,
		)
	}

	return kpiRows, kpiSubDetails, nil
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
