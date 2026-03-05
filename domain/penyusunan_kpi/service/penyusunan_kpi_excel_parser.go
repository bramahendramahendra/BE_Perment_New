package service

import (
	"fmt"
	"math"
	"mime/multipart"
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

	polarisasiMaximize = "Maximize"
	polarisasiMinimize = "Minimize"

	cappingOption1 = "100%"
	cappingOption2 = "110%"

	qualifierYa    = "Ya"
	qualifierTidak = "Tidak"

	totalBobotExpected = 100.0
	bobotTolerance     = 0.01 // toleransi floating point untuk total bobot
)

// =============================================
// MAIN PARSER FUNCTION
// =============================================

// ParseAndValidateExcel membaca file Excel dari multipart.FileHeader,
// memvalidasi setiap baris mulai dari baris ke-3 (baris ke-2 = header),
// dan mengembalikan slice PenyusunanKpiSubDetailRow beserta error jika ada.
//
// Aturan validasi per kolom:
//   - Col A  : angka
//   - Col B  : free text, tidak boleh blank
//   - Col C  : free text, tidak boleh blank
//   - Col D  : enum Maximize / Minimize, tidak boleh blank
//   - Col E  : enum 100% / 110%, tidak boleh blank
//   - Col F  : angka 2 desimal, total semua baris harus = 100%
//   - Col G  : free text, tidak boleh blank
//   - Col H  : free text, tidak boleh blank
//   - Col I  : angka 2 desimal
//   - Col J  : free text, tidak boleh blank
//   - Col K  : angka 2 desimal
//   - Col L  : enum Ya / Tidak, tidak boleh blank
//   - Col M  : wajib diisi jika Col L = "Ya"
//   - Col N  : wajib diisi jika Col L = "Ya"
//   - Col O  : wajib diisi jika Col L = "Ya"
//   - Col P  : free text, tidak boleh blank
//   - Col Q  : free text, tidak boleh blank
//   - Col R  : free text, tidak boleh blank
//   - Col S  : free text, tidak boleh blank
//   - Col T  : free text, tidak boleh blank
//   - Col U  : free text, tidak boleh blank
func ParseAndValidateExcel(file *multipart.FileHeader) ([]dto.PenyusunanKpiSubDetailRow, error) {
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

	// Ambil sheet pertama
	sheetName := xlsx.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki sheet", file.Filename)
	}

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", sheetName, err)
	}

	// Pastikan ada data mulai baris ke-3 (index 2)
	if len(rows) < excelDataStartRow {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki data (data dimulai dari baris %d)",
			file.Filename, excelDataStartRow)
	}

	// --- Loop setiap baris data ---
	var result []dto.PenyusunanKpiSubDetailRow
	var totalBobot float64

	for rowIdx := excelDataStartRow - 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]
		displayRow := rowIdx + 1 // nomor baris yang tampil ke user (1-based)

		// Pastikan baris memiliki cukup kolom (minimal 21 kolom A-U)
		// Padding jika kolom kurang
		for len(row) < 21 {
			row = append(row, "")
		}

		// Ambil nilai tiap kolom
		colA := strings.TrimSpace(row[0])  // NO
		colB := strings.TrimSpace(row[1])  // KPI
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
		colP := strings.TrimSpace(row[15]) // Result
		colQ := strings.TrimSpace(row[16]) // Deskripsi Result
		colR := strings.TrimSpace(row[17]) // Process
		colS := strings.TrimSpace(row[18]) // Deskripsi Process
		colT := strings.TrimSpace(row[19]) // Context
		colU := strings.TrimSpace(row[20]) // Deskripsi Context

		// Lewati baris kosong (semua kolom kosong)
		if colA == "" && colB == "" && colC == "" {
			continue
		}

		// --- Validasi Col A: NO (angka) ---
		no, errNo := strconv.Atoi(colA)
		if errNo != nil {
			return nil, fmt.Errorf("baris %d, Kolom A (NO): harus berupa angka, nilai saat ini: '%s'",
				displayRow, colA)
		}

		// --- Validasi Col B: KPI (free text, tidak boleh blank) ---
		if colB == "" {
			return nil, fmt.Errorf("baris %d, Kolom B (KPI): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col C: Sub KPI (free text, tidak boleh blank) ---
		if colC == "" {
			return nil, fmt.Errorf("baris %d, Kolom C (Sub KPI): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col D: Polarisasi (enum) ---
		if colD == "" {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): tidak boleh kosong", displayRow)
		}
		if colD != polarisasiMaximize && colD != polarisasiMinimize {
			return nil, fmt.Errorf("baris %d, Kolom D (Polarisasi): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colD, polarisasiMaximize, polarisasiMinimize)
		}

		// --- Validasi Col E: Capping (enum) ---
		if colE == "" {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): tidak boleh kosong", displayRow)
		}
		if colE != cappingOption1 && colE != cappingOption2 {
			return nil, fmt.Errorf("baris %d, Kolom E (Capping): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colE, cappingOption1, cappingOption2)
		}

		// --- Validasi Col F: Bobot % (angka 2 desimal) ---
		if colF == "" {
			return nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): tidak boleh kosong", displayRow)
		}
		bobot, errBobot := parseFloat2Decimal(colF)
		if errBobot != nil {
			return nil, fmt.Errorf("baris %d, Kolom F (Bobot %%): harus berupa angka 2 desimal tanpa simbol persen, nilai saat ini: '%s'",
				displayRow, colF)
		}
		totalBobot += bobot

		// --- Validasi Col G: Glossary ---
		if colG == "" {
			return nil, fmt.Errorf("baris %d, Kolom G (Glossary): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col H: Target Triwulanan ---
		if colH == "" {
			return nil, fmt.Errorf("baris %d, Kolom H (Target Triwulanan): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col I: Target Kuantitatif Triwulanan (angka 2 desimal) ---
		targetKuantitatifTriwulan, errI := parseFloat2Decimal(colI)
		if errI != nil {
			return nil, fmt.Errorf("baris %d, Kolom I (Target Kuantitatif Triwulanan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colI)
		}

		// --- Validasi Col J: Target Tahunan ---
		if colJ == "" {
			return nil, fmt.Errorf("baris %d, Kolom J (Target Tahunan): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col K: Target Kuantitatif Tahunan (angka 2 desimal) ---
		targetKuantitatifTahunan, errK := parseFloat2Decimal(colK)
		if errK != nil {
			return nil, fmt.Errorf("baris %d, Kolom K (Target Kuantitatif Tahunan): harus berupa angka 2 desimal, nilai saat ini: '%s'",
				displayRow, colK)
		}

		// --- Validasi Col L: Terdapat Qualifier (enum) ---
		if colL == "" {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): tidak boleh kosong", displayRow)
		}
		if colL != qualifierYa && colL != qualifierTidak {
			return nil, fmt.Errorf("baris %d, Kolom L (Terdapat Qualifier): nilai tidak valid '%s', harus '%s' atau '%s'",
				displayRow, colL, qualifierYa, qualifierTidak)
		}

		// --- Validasi Col M, N, O: wajib jika Kolom L = "Ya" ---
		if strings.EqualFold(colL, qualifierYa) {
			if colM == "" {
				return nil, fmt.Errorf("baris %d, Kolom M (Qualifier): tidak boleh kosong karena Kolom L = 'Ya'",
					displayRow)
			}
			if colN == "" {
				return nil, fmt.Errorf("baris %d, Kolom N (Deskripsi Qualifier): tidak boleh kosong karena Kolom L = 'Ya'",
					displayRow)
			}
			if colO == "" {
				return nil, fmt.Errorf("baris %d, Kolom O (Target Qualifier): tidak boleh kosong karena Kolom L = 'Ya'",
					displayRow)
			}
		}

		// --- Validasi Col P: Result ---
		if colP == "" {
			return nil, fmt.Errorf("baris %d, Kolom P (Result): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col Q: Deskripsi Result ---
		if colQ == "" {
			return nil, fmt.Errorf("baris %d, Kolom Q (Deskripsi Result): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col R: Process ---
		if colR == "" {
			return nil, fmt.Errorf("baris %d, Kolom R (Process): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col S: Deskripsi Process ---
		if colS == "" {
			return nil, fmt.Errorf("baris %d, Kolom S (Deskripsi Process): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col T: Context ---
		if colT == "" {
			return nil, fmt.Errorf("baris %d, Kolom T (Context): tidak boleh kosong", displayRow)
		}

		// --- Validasi Col U: Deskripsi Context ---
		if colU == "" {
			return nil, fmt.Errorf("baris %d, Kolom U (Deskripsi Context): tidak boleh kosong", displayRow)
		}

		// --- Semua validasi per baris lolos, tambahkan ke result ---
		result = append(result, dto.PenyusunanKpiSubDetailRow{
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
			Qualifier:                 colM,
			DeskripsiQualifier:        colN,
			TargetQualifier:           colO,
			Result:                    colP,
			DeskripsiResult:           colQ,
			Process:                   colR,
			DeskripsiProcess:          colS,
			Context:                   colT,
			DeskripsiContext:          colU,
		})
	}

	// --- Validasi total bobot setelah semua baris diproses ---
	if len(result) == 0 {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki data yang valid", file.Filename)
	}

	totalBobotRounded := math.Round(totalBobot*100) / 100
	if math.Abs(totalBobotRounded-totalBobotExpected) > bobotTolerance {
		return nil, fmt.Errorf(
			"file Excel '%s': total Bobot (Kolom F) = %.2f%%, harus tepat 100%%",
			file.Filename, totalBobotRounded,
		)
	}

	return result, nil
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// parseFloat2Decimal mem-parse string menjadi float64 dengan 2 angka di belakang koma.
// Mengembalikan error jika string bukan angka valid.
func parseFloat2Decimal(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}

	// Hapus simbol % jika ada (misal jika user salah format)
	cleaned := strings.ReplaceAll(s, "%", "")
	cleaned = strings.TrimSpace(cleaned)

	val, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, fmt.Errorf("'%s' bukan angka valid", s)
	}

	// Bulatkan ke 2 desimal
	rounded := math.Round(val*100) / 100
	return rounded, nil
}
