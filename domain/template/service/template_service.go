package service

import (
	"bytes"
	"fmt"
	"strconv"

	dto "permen_api/domain/template/dto"
	"permen_api/errors"

	"github.com/xuri/excelize/v2"
)

// =============================================================================
// Konstanta header kolom
// =============================================================================

// headerRow adalah nomor baris header kolom (row ke-2 di Excel, index 1-based).
const headerRow = 2

// columnsBase adalah header kolom A–O (untuk TW1 dan TW3).
var columnsBase = []string{
	"No.",
	"KPI (text)",
	"Sub KPI (text)",
	"Polarisasi (Maximize atau Minimize)",
	"Capping (100% atau 110%)",
	"Bobot % (bilangan bulat dua angka belakang koma)",
	"Glossary (text)",
	"Target Triwulanan (text)",
	"Target Kuantitatif Triwulanan (angka)",
	"Target Tahunan (text)",
	"Target Kuantitatif Tahunan (angka)",
	"Terdapat Qualifier (Ya/Tidak)",
	"Qualifier (text)",
	"Deskripsi Qualifier (text)",
	"Target Qualifier (text)",
}

// columnsExtended adalah header kolom tambahan P–U khusus TW2 dan TW4.
var columnsExtended = []string{
	"Result (text)",
	"Deskripsi Result (text)",
	"Process (text)",
	"Deskripsi Process (text)",
	"Context (text)",
	"Deskripsi Context (text)",
}

// =============================================================================
// GenerateFormatPenyusunanKpi — TIDAK DIUBAH dari versi asli
// =============================================================================

func (s *templateService) GenerateFormatPenyusunanKpi(req *dto.FormatPenyusunanKpiRequest) ([]byte, string, error) {
	// TW2 dan TW4 menggunakan format kolom A–U (extended).
	// TW1 dan TW3 menggunakan format kolom A–O (base).
	useExtended := req.Triwulan == "TW2" || req.Triwulan == "TW4"

	// Nama sheet mengikuti nilai triwulan dari request (TW1, TW2, TW3, TW4).
	sheetName := req.Triwulan

	f := excelize.NewFile()
	defer f.Close()

	// Rename default sheet "Sheet1" menjadi nama sheet yang sesuai
	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nama sheet: %v", err)}
	}

	// Gabungkan semua kolom header sesuai kondisi
	// Pakai make+copy agar tidak mutasi package-level var columnsBase
	allColumns := make([]string, len(columnsBase))
	copy(allColumns, columnsBase)
	if useExtended {
		allColumns = append(allColumns, columnsExtended...)
	}

	// Buat style untuk row 1 (merge label qualifier) — background kuning
	styleRow1, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style row1: %v", err)}
	}

	// Buat style untuk header kolom (row 2) — background biru muda + bold
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"BDD7EE"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style header: %v", err)}
	}

	// Buat style untuk data cell (row 3 dst) — border tipis
	styleData, err := f.NewStyle(&excelize.Style{
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style data: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: Merge cell M1:O1 → label "Jika Ya, di Terdapat Qualifier (kolom L)"
	// -------------------------------------------------------------------------
	mergeStart := "M1"
	mergeEnd := "O1"
	if err := f.MergeCell(sheetName, mergeStart, mergeEnd); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal merge cell: %v", err)}
	}
	if err := f.SetCellValue(sheetName, mergeStart, "Jika Ya, di Terdapat Qualifier (kolom L)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai merge cell: %v", err)}
	}
	// Set style row1 untuk range M1:O1
	if err := f.SetCellStyle(sheetName, mergeStart, mergeEnd, styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style merge cell: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: E1 = "Total Bobot", F1 = SUM(F3:F100)
	// -------------------------------------------------------------------------
	if err := f.SetCellValue(sheetName, "E1", "Total Bobot"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai E1: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, "E1", "E1", styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style E1: %v", err)}
	}
	if err := f.SetCellFormula(sheetName, "F1", "SUM(F3:F100)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set formula F1: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, "F1", "F1", styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style F1: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 2: Tulis header kolom
	// -------------------------------------------------------------------------
	for colIdx, header := range allColumns {
		cellName, _ := excelize.CoordinatesToCellName(colIdx+1, headerRow)
		if err := f.SetCellValue(sheetName, cellName, header); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set header %s: %v", cellName, err)}
		}
		if err := f.SetCellStyle(sheetName, cellName, cellName, styleHeader); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style header %s: %v", cellName, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Row 3–100: Pre-fill area data dengan style & data validation
	// -------------------------------------------------------------------------
	dataStartRow := 3
	dataEndRow := 100
	lastColIdx := len(allColumns)

	for rowIdx := dataStartRow; rowIdx <= dataEndRow; rowIdx++ {
		for colIdx := 1; colIdx <= lastColIdx; colIdx++ {
			cellName, _ := excelize.CoordinatesToCellName(colIdx, rowIdx)
			if err := f.SetCellStyle(sheetName, cellName, cellName, styleData); err != nil {
				return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style data %s: %v", cellName, err)}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Data Validation per kolom
	// -------------------------------------------------------------------------
	sqrefDataRange := func(col string) string {
		return fmt.Sprintf("%s%d:%s%d", col, dataStartRow, col, dataEndRow)
	}

	// Kolom A (No.) → Angka bulat
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "whole",
		Operator:         "greaterThan",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Kolom No. harus berupa angka bulat positif."),
		Sqref:            sqrefDataRange("A"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom A: %v", err)}
	}

	// Kolom D (Polarisasi) → Dropdown: Maximize / Minimize
	dvPolarisasi := excelize.NewDataValidation(true)
	dvPolarisasi.Sqref = sqrefDataRange("D")
	if err := dvPolarisasi.SetDropList([]string{"Maximize", "Minimize"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Polarisasi: %v", err)}
	}
	dvPolarisasi.ShowErrorMessage = true
	dvPolarisasi.ErrorStyle = strPtr("stop")
	dvPolarisasi.ErrorTitle = strPtr("Input Tidak Valid")
	dvPolarisasi.Error = strPtr("Pilih salah satu: Maximize atau Minimize.")
	if err := f.AddDataValidation(sheetName, dvPolarisasi); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Polarisasi: %v", err)}
	}

	// Kolom E (Capping) → Dropdown: 100% / 110%
	dvCapping := excelize.NewDataValidation(true)
	dvCapping.Sqref = sqrefDataRange("E")
	if err := dvCapping.SetDropList([]string{"100%", "110%"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Capping: %v", err)}
	}
	dvCapping.ShowErrorMessage = true
	dvCapping.ErrorStyle = strPtr("stop")
	dvCapping.ErrorTitle = strPtr("Input Tidak Valid")
	dvCapping.Error = strPtr("Pilih salah satu: 100% atau 110%.")
	if err := f.AddDataValidation(sheetName, dvCapping); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Capping: %v", err)}
	}

	// Kolom F (Bobot %) → Angka desimal 2 digit di belakang koma, range 0–100
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "between",
		Formula1:         "0",
		Formula2:         "100",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Bobot % harus berupa angka antara 0 sampai 100 (maks. 2 angka di belakang koma, tanpa simbol %)."),
		Sqref:            sqrefDataRange("F"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom F: %v", err)}
	}

	// Kolom I (Target Kuantitatif Triwulanan) → Angka desimal
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Target Kuantitatif Triwulanan harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("I"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom I: %v", err)}
	}

	// Kolom K (Target Kuantitatif Tahunan) → Angka desimal
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Target Kuantitatif Tahunan harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("K"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom K: %v", err)}
	}

	// Kolom L (Terdapat Qualifier) → Dropdown: Ya / Tidak
	dvQualifier := excelize.NewDataValidation(true)
	dvQualifier.Sqref = sqrefDataRange("L")
	if err := dvQualifier.SetDropList([]string{"Ya", "Tidak"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Terdapat Qualifier: %v", err)}
	}
	dvQualifier.ShowErrorMessage = true
	dvQualifier.ErrorStyle = strPtr("stop")
	dvQualifier.ErrorTitle = strPtr("Input Tidak Valid")
	dvQualifier.Error = strPtr("Pilih salah satu: Ya atau Tidak.")
	if err := f.AddDataValidation(sheetName, dvQualifier); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Terdapat Qualifier: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Set lebar kolom agar lebih mudah dibaca
	// -------------------------------------------------------------------------
	colWidths := map[string]float64{
		"A": 6,  // No.
		"B": 25, // KPI
		"C": 25, // Sub KPI
		"D": 20, // Polarisasi
		"E": 18, // Capping
		"F": 20, // Bobot %
		"G": 30, // Glossary
		"H": 25, // Target Triwulanan
		"I": 22, // Target Kuantitatif Triwulanan
		"J": 25, // Target Tahunan
		"K": 22, // Target Kuantitatif Tahunan
		"L": 22, // Terdapat Qualifier
		"M": 25, // Qualifier
		"N": 30, // Deskripsi Qualifier
		"O": 25, // Target Qualifier
	}
	if useExtended {
		colWidths["P"] = 25
		colWidths["Q"] = 30
		colWidths["R"] = 25
		colWidths["S"] = 30
		colWidths["T"] = 25
		colWidths["U"] = 30
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom %s: %v", col, err)}
		}
	}

	// Set tinggi row header (row 2)
	if err := f.SetRowHeight(sheetName, headerRow, 40); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi row header: %v", err)}
	}

	// Freeze pane di bawah row header agar header selalu terlihat saat scroll
	if err := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set freeze pane: %v", err)}
	}

	// =========================================================================
	// Sheet 2: "KPI" — data dari mst_kpi join mst_polarisasi
	// =========================================================================
	if err := s.generateSheetKpi(f); err != nil {
		return nil, "", err
	}

	// -------------------------------------------------------------------------
	// Tulis ke buffer bytes
	// -------------------------------------------------------------------------
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal write file Excel: %v", err)}
	}

	filename := fmt.Sprintf("Format Penyusunan KPI Aplikasi Performance Management %s.xlsx", req.Triwulan)
	return buf.Bytes(), filename, nil
}

// =============================================================================
// GenerateRevisionPenyusunanKpi
// =============================================================================

func (s *templateService) GenerateRevisionPenyusunanKpi(req *dto.RevisionPenyusunanKpiRequest) ([]byte, string, error) {

	// Ambil data dari DB (header + seluruh baris sub KPI)
	excelData, err := s.repo.GetRevisionPenyusunanKpiData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}
	if excelData == nil {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	// TW2 dan TW4 menggunakan format kolom A–U (extended).
	// TW1 dan TW3 menggunakan format kolom A–O (base).
	useExtended := excelData.Triwulan == "TW2" || excelData.Triwulan == "TW4"

	// Nama sheet mengikuti nilai triwulan dari db (TW1, TW2, TW3, TW4).
	sheetName := excelData.Triwulan

	f := excelize.NewFile()
	defer f.Close()

	// Rename default sheet "Sheet1" menjadi nama sheet yang sesuai
	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nama sheet: %v", err)}
	}

	// Gabungkan semua kolom header sesuai kondisi
	// Pakai make+copy agar tidak mutasi package-level var columnsBase
	allColumns := make([]string, len(columnsBase))
	copy(allColumns, columnsBase)
	if useExtended {
		allColumns = append(allColumns, columnsExtended...)
	}

	// Buat style untuk row 1 (merge label qualifier) — background kuning
	styleRow1, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style row1: %v", err)}
	}

	// Buat style untuk header kolom (row 2) — background biru muda + bold
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"BDD7EE"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style header: %v", err)}
	}

	// Buat style untuk data cell (row 3 dst) — border tipis
	styleData, err := f.NewStyle(&excelize.Style{
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style data: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: Merge cell M1:O1 → label "Jika Ya, di Terdapat Qualifier (kolom L)"
	// -------------------------------------------------------------------------
	mergeStart := "M1"
	mergeEnd := "O1"
	if err := f.MergeCell(sheetName, mergeStart, mergeEnd); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal merge cell %s:%s: %v", mergeStart, mergeEnd, err)}
	}
	if err := f.SetCellValue(sheetName, mergeStart, "Jika Ya, di Terdapat Qualifier (kolom L)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai merge cell: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, mergeStart, mergeEnd, styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style merge cell: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: E1 = "Total Bobot", F1 = SUM(F3:F100)
	// -------------------------------------------------------------------------
	if err := f.SetCellValue(sheetName, "E1", "Total Bobot"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai E1: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, "E1", "E1", styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style E1: %v", err)}
	}
	if err := f.SetCellFormula(sheetName, "F1", "SUM(F3:F100)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set formula F1: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, "F1", "F1", styleRow1); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style F1: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 2: Tulis header kolom
	// -------------------------------------------------------------------------
	for colIdx, header := range allColumns {
		cellName, _ := excelize.CoordinatesToCellName(colIdx+1, headerRow)
		if err := f.SetCellValue(sheetName, cellName, header); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set header kolom %s: %v", cellName, err)}
		}
		if err := f.SetCellStyle(sheetName, cellName, cellName, styleHeader); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style header kolom %s: %v", cellName, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Row 3+: Pre-fill area data dengan style & data validation. Tulis data baris dari DB
	// -------------------------------------------------------------------------
	dataStartRow := 3
	dataEndRow := 100
	lastColIdx := len(allColumns)

	for rowIdx, row := range excelData.Rows {
		rowNum := dataStartRow + rowIdx

		values := []interface{}{
			rowIdx + 1,                    // A: No.
			row.KpiNama,                   // B: KPI
			row.SubKpi,                    // C: Sub KPI
			row.Polarisasi,                // D: Polarisasi
			row.Capping + "%",             // E: Capping
			parseFloatOrString(row.Bobot), // F: Bobot
			row.DeskripsiGlossary,         // G: Glossary
			row.TargetTriwulan,            // H: Target Triwulanan
			parseFloatOrString(row.TargetKuantitatifTriwulan), // I
			row.TargetTahunan, // J: Target Tahunan
			parseFloatOrString(row.TargetKuantitatifTahunan), // K
			row.TerdapatQualifier,                            // L: Ya/Tidak (dikonversi di repo)
			row.ItemQualifier,                                // M: Qualifier
			row.DeskripsiQualifier,                           // N: Deskripsi Qualifier
			row.TargetQualifier,                              // O: Target Qualifier
		}

		if useExtended {
			values = append(values,
				row.NamaResult,       // P
				row.DeskripsiResult,  // Q
				row.NamaProcess,      // R
				row.DeskripsiProcess, // S
				row.NamaContext,      // T
				row.DeskripsiContext, // U
			)
		}

		for colIdx, val := range values {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cellName, val); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set nilai baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
			if err := f.SetCellStyle(sheetName, cellName, cellName, styleData); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set style baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Pre-fill style untuk baris kosong setelah data
	// agar user bisa menambah baris baru dengan tampilan yang konsisten
	// -------------------------------------------------------------------------
	lastDataRow := dataStartRow + len(excelData.Rows)
	for rowIdx := lastDataRow; rowIdx <= dataEndRow; rowIdx++ {
		for colIdx := 1; colIdx <= lastColIdx; colIdx++ {
			cellName, _ := excelize.CoordinatesToCellName(colIdx, rowIdx)
			if err := f.SetCellStyle(sheetName, cellName, cellName, styleData); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set pre-fill style baris %d kolom %d: %v", rowIdx, colIdx, err),
				}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Data Validation per kolom
	// -------------------------------------------------------------------------
	sqrefDataRange := func(col string) string {
		return fmt.Sprintf("%s%d:%s%d", col, dataStartRow, col, dataEndRow)
	}

	// Kolom A (No.) → Angka bulat
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "whole",
		Operator:         "greaterThan",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Kolom No. harus berupa angka bulat positif."),
		Sqref:            sqrefDataRange("A"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom A: %v", err)}
	}

	// Kolom D (Polarisasi) → Dropdown: Maximize / Minimize
	dvPolarisasi := excelize.NewDataValidation(true)
	dvPolarisasi.Sqref = sqrefDataRange("D")
	if err := dvPolarisasi.SetDropList([]string{"Maximize", "Minimize"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Polarisasi: %v", err)}
	}
	dvPolarisasi.ShowErrorMessage = true
	dvPolarisasi.ErrorStyle = strPtr("stop")
	dvPolarisasi.ErrorTitle = strPtr("Input Tidak Valid")
	dvPolarisasi.Error = strPtr("Pilih salah satu: Maximize atau Minimize.")
	if err := f.AddDataValidation(sheetName, dvPolarisasi); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Polarisasi: %v", err)}
	}

	// Kolom E (Capping) → Dropdown: 100% / 110%
	dvCapping := excelize.NewDataValidation(true)
	dvCapping.Sqref = sqrefDataRange("E")
	if err := dvCapping.SetDropList([]string{"100%", "110%"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Capping: %v", err)}
	}
	dvCapping.ShowErrorMessage = true
	dvCapping.ErrorStyle = strPtr("stop")
	dvCapping.ErrorTitle = strPtr("Input Tidak Valid")
	dvCapping.Error = strPtr("Pilih salah satu: 100% atau 110%.")
	if err := f.AddDataValidation(sheetName, dvCapping); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Capping: %v", err)}
	}

	// Kolom F (Bobot %) → Angka desimal 2 digit di belakang koma, range 0–100
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "between",
		Formula1:         "0",
		Formula2:         "100",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Bobot % harus berupa angka antara 0 sampai 100 (maks. 2 angka di belakang koma, tanpa simbol %)."),
		Sqref:            sqrefDataRange("F"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom F: %v", err)}
	}

	// Kolom I (Target Kuantitatif Triwulanan) → Angka desimal
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Target Kuantitatif Triwulanan harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("I"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom I: %v", err)}
	}

	// Kolom K (Target Kuantitatif Tahunan) → Angka desimal
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Target Kuantitatif Tahunan harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("K"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom K: %v", err)}
	}

	// Kolom L (Terdapat Qualifier) → Dropdown: Ya / Tidak
	dvQualifier := excelize.NewDataValidation(true)
	dvQualifier.Sqref = sqrefDataRange("L")
	if err := dvQualifier.SetDropList([]string{"Ya", "Tidak"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Terdapat Qualifier: %v", err)}
	}
	dvQualifier.ShowErrorMessage = true
	dvQualifier.ErrorStyle = strPtr("stop")
	dvQualifier.ErrorTitle = strPtr("Input Tidak Valid")
	dvQualifier.Error = strPtr("Pilih salah satu: Ya atau Tidak.")
	if err := f.AddDataValidation(sheetName, dvQualifier); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Terdapat Qualifier: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Set lebar kolom
	// -------------------------------------------------------------------------
	colWidths := map[string]float64{
		"A": 6,  // No.
		"B": 25, // KPI
		"C": 25, // Sub KPI
		"D": 20, // Polarisasi
		"E": 18, // Capping
		"F": 20, // Bobot %
		"G": 30, // Glossary
		"H": 25, // Target Triwulanan
		"I": 22, // Target Kuantitatif Triwulanan
		"J": 25, // Target Tahunan
		"K": 22, // Target Kuantitatif Tahunan
		"L": 22, // Terdapat Qualifier
		"M": 25, // Qualifier
		"N": 30, // Deskripsi Qualifier
		"O": 25, // Target Qualifier
	}
	if useExtended {
		colWidths["P"] = 25
		colWidths["Q"] = 30
		colWidths["R"] = 25
		colWidths["S"] = 30
		colWidths["T"] = 25
		colWidths["U"] = 30
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom %s: %v", col, err)}
		}
	}

	// Set tinggi row header (row 2)
	if err := f.SetRowHeight(sheetName, headerRow, 40); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi row header: %v", err)}
	}

	// Freeze pane di bawah row header agar header selalu terlihat saat scroll
	if err := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      2,
		TopLeftCell: "A3",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set freeze pane: %v", err)}
	}

	// =========================================================================
	// Sheet 2: "KPI" — data dari mst_kpi join mst_polarisasi
	// =========================================================================
	if err := s.generateSheetKpi(f); err != nil {
		return nil, "", err
	}

	// -------------------------------------------------------------------------
	// Tulis ke buffer bytes
	// -------------------------------------------------------------------------
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal write file Excel: %v", err)}
	}

	filename := fmt.Sprintf("Revisi Penyusunan KPI Aplikasi Performance Management %s %s %s.xlsx", excelData.Triwulan, excelData.Tahun, excelData.KostlTx)
	return buf.Bytes(), filename, nil
}

// =============================================================================
// GenerateFormatRealisasiKpi
// =============================================================================

// columnsRealisasiBase adalah header kolom A–M (sama untuk semua triwulan).
var columnsRealisasiBase = []string{
	"No",
	"KPI",
	"Sub KPI",
	"Polarisasi",
	"Capping",
	"Bobot %",
	"Target Triwulanan",
	"Qualifier",
	"Target Qualifier",
	"Realisasi",
	"Realisasi Kuantitatif",
	"Realisasi Qualifier",
	"Realisasi Qualifier Kuantitatif",
}

// columnsRealisasiExtendedTW13 adalah header kolom N–S (khusus TW1 dan TW3).
var columnsRealisasiExtendedTW13 = []string{
	"Result",
	"Deskripsi Result",
	"Process",
	"Deskripsi Process",
	"Context",
	"Deskripsi Context",
}

// columnsRealisasiExtendedTW24 adalah header kolom N–Y (khusus TW2 dan TW4).
var columnsRealisasiExtendedTW24 = []string{
	"Result",
	"Deskripsi Result",
	"Realisasi Result",
	"Link Result",
	"Process",
	"Deskripsi Process",
	"Realisasi Process",
	"Link Process",
	"Context",
	"Deskripsi Context",
	"Realisasi Context",
	"Link Context",
}

func (s *templateService) GenerateFormatRealisasiKpi(req *dto.FormatRealisasiKpiRequest) ([]byte, string, error) {
	// Ambil data dari DB (header + seluruh baris sub KPI)
	excelData, err := s.repo.GetRevisionPenyusunanKpiData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}
	if excelData == nil {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	// TW1/TW3 → extended kolom N–S; TW2/TW4 → extended kolom N–Y
	isTW24 := req.Triwulan == "TW2" || req.Triwulan == "TW4"

	sheetName := req.Triwulan

	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nama sheet: %v", err)}
	}

	// Gabungkan header kolom sesuai triwulan
	allColumns := make([]string, len(columnsRealisasiBase))
	copy(allColumns, columnsRealisasiBase)
	if isTW24 {
		allColumns = append(allColumns, columnsRealisasiExtendedTW24...)
	} else {
		allColumns = append(allColumns, columnsRealisasiExtendedTW13...)
	}

	// -------------------------------------------------------------------------
	// Style
	// -------------------------------------------------------------------------
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"BDD7EE"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style header: %v", err)}
	}

	styleData, err := f.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: true,
		},
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style data: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: Header kolom (data mulai row 2)
	// -------------------------------------------------------------------------
	const realisasiHeaderRow = 1
	const realisasiDataStartRow = 2
	const realisasiDataEndRow = 100

	for colIdx, header := range allColumns {
		cellName, _ := excelize.CoordinatesToCellName(colIdx+1, realisasiHeaderRow)
		if err := f.SetCellValue(sheetName, cellName, header); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set header %s: %v", cellName, err)}
		}
		if err := f.SetCellStyle(sheetName, cellName, cellName, styleHeader); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style header %s: %v", cellName, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Row 2+: Tulis data baris dari DB
	// -------------------------------------------------------------------------
	lastColIdx := len(allColumns)

	for rowIdx, row := range excelData.Rows {
		rowNum := realisasiDataStartRow + rowIdx

		var values []interface{}
		if isTW24 {
			values = []interface{}{
				rowIdx + 1,                    // A: No
				row.KpiNama,                   // B: KPI
				row.SubKpi,                    // C: Sub KPI
				row.Polarisasi,                // D: Polarisasi
				row.Capping + "%",             // E: Capping
				parseFloatOrString(row.Bobot), // F: Bobot %
				row.TargetTriwulan,            // G: Target Triwulanan
				realisasiQualifierOrDash(row.ItemQualifier),   // H: Qualifier
				realisasiQualifierOrDash(row.TargetQualifier), // I: Target Qualifier
				"", "", "", "", // J–M: kosong (diisi user)
				row.NamaResult,      // N: Result
				row.DeskripsiResult, // O: Deskripsi Result
				"", "",              // P–Q: kosong (diisi user)
				row.NamaProcess,      // R: Process
				row.DeskripsiProcess, // S: Deskripsi Process
				"", "",               // T–U: kosong (diisi user)
				row.NamaContext,      // V: Context
				row.DeskripsiContext, // W: Deskripsi Context
				"", "",               // X–Y: kosong (diisi user)
			}
		} else {
			values = []interface{}{
				rowIdx + 1,                    // A: No
				row.KpiNama,                   // B: KPI
				row.SubKpi,                    // C: Sub KPI
				row.Polarisasi,                // D: Polarisasi
				row.Capping + "%",             // E: Capping
				parseFloatOrString(row.Bobot), // F: Bobot %
				row.TargetTriwulan,            // G: Target Triwulanan
				realisasiQualifierOrDash(row.ItemQualifier),   // H: Qualifier
				realisasiQualifierOrDash(row.TargetQualifier), // I: Target Qualifier
				"", "", "", "", // J–M: kosong (diisi user)
				row.NamaResult,       // N: Result
				row.DeskripsiResult,  // O: Deskripsi Result
				row.NamaProcess,      // P: Process
				row.DeskripsiProcess, // Q: Deskripsi Process
				row.NamaContext,      // R: Context
				row.DeskripsiContext, // S: Deskripsi Context
			}
		}

		for colIdx, val := range values {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cellName, val); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set nilai baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
			if err := f.SetCellStyle(sheetName, cellName, cellName, styleData); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set style baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
		}
	}

	// Pre-fill style baris kosong setelah data
	lastDataRow := realisasiDataStartRow + len(excelData.Rows)
	for rowIdx := lastDataRow; rowIdx <= realisasiDataEndRow; rowIdx++ {
		for colIdx := 1; colIdx <= lastColIdx; colIdx++ {
			cellName, _ := excelize.CoordinatesToCellName(colIdx, rowIdx)
			if err := f.SetCellStyle(sheetName, cellName, cellName, styleData); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set pre-fill style baris %d kolom %d: %v", rowIdx, colIdx, err),
				}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Data Validation
	// -------------------------------------------------------------------------
	sqrefDataRange := func(col string) string {
		return fmt.Sprintf("%s%d:%s%d", col, realisasiDataStartRow, col, realisasiDataEndRow)
	}

	// Kolom A (No) → Angka bulat positif
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "whole",
		Operator:         "greaterThan",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Kolom No harus berupa angka bulat positif."),
		Sqref:            sqrefDataRange("A"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom A: %v", err)}
	}

	// Kolom D (Polarisasi) → Dropdown: Maximize / Minimize
	dvPolarisasi := excelize.NewDataValidation(true)
	dvPolarisasi.Sqref = sqrefDataRange("D")
	if err := dvPolarisasi.SetDropList([]string{"Maximize", "Minimize"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Polarisasi: %v", err)}
	}
	dvPolarisasi.ShowErrorMessage = true
	dvPolarisasi.ErrorStyle = strPtr("stop")
	dvPolarisasi.ErrorTitle = strPtr("Input Tidak Valid")
	dvPolarisasi.Error = strPtr("Pilih salah satu: Maximize atau Minimize.")
	if err := f.AddDataValidation(sheetName, dvPolarisasi); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Polarisasi: %v", err)}
	}

	// Kolom E (Capping) → Dropdown: 100% / 110%
	dvCapping := excelize.NewDataValidation(true)
	dvCapping.Sqref = sqrefDataRange("E")
	if err := dvCapping.SetDropList([]string{"100%", "110%"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Capping: %v", err)}
	}
	dvCapping.ShowErrorMessage = true
	dvCapping.ErrorStyle = strPtr("stop")
	dvCapping.ErrorTitle = strPtr("Input Tidak Valid")
	dvCapping.Error = strPtr("Pilih salah satu: 100% atau 110%.")
	if err := f.AddDataValidation(sheetName, dvCapping); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Capping: %v", err)}
	}

	// Kolom F (Bobot %) → Angka desimal 0–100
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "between",
		Formula1:         "0",
		Formula2:         "100",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Bobot % harus berupa angka antara 0 sampai 100 (maks. 2 angka di belakang koma, tanpa simbol %)."),
		Sqref:            sqrefDataRange("F"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom F: %v", err)}
	}

	// Kolom K (Realisasi Kuantitatif) → Angka desimal (diisi user)
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Realisasi Kuantitatif harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("K"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom K: %v", err)}
	}

	// Kolom L (Realisasi Qualifier) → Dropdown: Ya / Tidak (diisi user)
	dvRealisasiQualifier := excelize.NewDataValidation(true)
	dvRealisasiQualifier.Sqref = sqrefDataRange("L")
	if err := dvRealisasiQualifier.SetDropList([]string{"Ya", "Tidak"}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set dropdown Realisasi Qualifier: %v", err)}
	}
	dvRealisasiQualifier.ShowErrorMessage = true
	dvRealisasiQualifier.ErrorStyle = strPtr("stop")
	dvRealisasiQualifier.ErrorTitle = strPtr("Input Tidak Valid")
	dvRealisasiQualifier.Error = strPtr("Pilih salah satu: Ya atau Tidak.")
	if err := f.AddDataValidation(sheetName, dvRealisasiQualifier); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi Realisasi Qualifier: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Set lebar kolom
	// -------------------------------------------------------------------------
	colWidths := map[string]float64{
		"A": 6,  // No
		"B": 25, // KPI
		"C": 25, // Sub KPI
		"D": 20, // Polarisasi
		"E": 18, // Capping
		"F": 20, // Bobot %
		"G": 25, // Target Triwulanan
		"H": 25, // Qualifier
		"I": 25, // Target Qualifier
		"J": 25, // Realisasi
		"K": 25, // Realisasi Kuantitatif
		"L": 25, // Realisasi Qualifier
		"M": 30, // Realisasi Qualifier Kuantitatif
	}
	if isTW24 {
		colWidths["N"] = 25 // Result
		colWidths["O"] = 30 // Deskripsi Result
		colWidths["P"] = 25 // Realisasi Result
		colWidths["Q"] = 25 // Link Result
		colWidths["R"] = 25 // Process
		colWidths["S"] = 30 // Deskripsi Process
		colWidths["T"] = 25 // Realisasi Process
		colWidths["U"] = 25 // Link Process
		colWidths["V"] = 25 // Context
		colWidths["W"] = 30 // Deskripsi Context
		colWidths["X"] = 25 // Realisasi Context
		colWidths["Y"] = 25 // Link Context
	} else {
		colWidths["N"] = 25 // Result
		colWidths["O"] = 30 // Deskripsi Result
		colWidths["P"] = 25 // Process
		colWidths["Q"] = 30 // Deskripsi Process
		colWidths["R"] = 25 // Context
		colWidths["S"] = 30 // Deskripsi Context
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom %s: %v", col, err)}
		}
	}

	// Set tinggi row header
	if err := f.SetRowHeight(sheetName, realisasiHeaderRow, 40); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi row header: %v", err)}
	}

	// Freeze pane di bawah row header
	if err := f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set freeze pane: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Sheet Protection: kunci semua sel, unlock hanya kolom input user
	// -------------------------------------------------------------------------
	// Tentukan kolom yang boleh diisi user sesuai triwulan:
	//   TW1/TW3 : J (Realisasi), K (Realisasi Kuantitatif),
	//             L (Realisasi Qualifier), M (Realisasi Qualifier Kuantitatif)
	//   TW2/TW4 : tambah P (Realisasi Result), Q (Link Result),
	//             T (Realisasi Process), U (Link Process),
	//             X (Realisasi Context), Y (Link Context)
	userInputCols := []string{"J", "K", "L", "M"}
	if isTW24 {
		userInputCols = append(userInputCols, "P", "Q", "T", "U", "X", "Y")
	}

	// Jumlah baris data aktual dari DB (baris 2 s.d. lastDataRow-1)
	totalDataRows := len(excelData.Rows)
	if totalDataRows > 0 {
		// Style "unlocked" untuk sel input user
		styleUnlocked, err := f.NewStyle(&excelize.Style{
			Protection: &excelize.Protection{
				Locked: false,
			},
			Border: borderStyle(),
			Alignment: &excelize.Alignment{
				Vertical: "top",
				WrapText: true,
			},
		})
		if err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style unlocked: %v", err)}
		}

		dataEndRow := realisasiDataStartRow + totalDataRows - 1
		for _, col := range userInputCols {
			rangeRef := fmt.Sprintf("%s%d:%s%d", col, realisasiDataStartRow, col, dataEndRow)
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, realisasiDataStartRow),
				fmt.Sprintf("%s%d", col, dataEndRow), styleUnlocked); err != nil {
				return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style unlocked %s: %v", rangeRef, err)}
			}
		}
	}

	// Aktifkan proteksi sheet — semua sel terkunci kecuali yang sudah di-unlock di atas
	if err := f.ProtectSheet(sheetName, &excelize.SheetProtectionOptions{
		SelectLockedCells:   true,
		SelectUnlockedCells: true,
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal protect sheet: %v", err)}
	}

	// =========================================================================
	// Sheet 2: "KPI" — data dari mst_kpi join mst_polarisasi
	// =========================================================================
	if err := s.generateSheetKpi(f); err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal write file Excel: %v", err)}
	}

	filename := fmt.Sprintf("Format Realisasi KPI Aplikasi Performance Management %s %s %s.xlsx", excelData.KostlTx, excelData.Tahun, req.Triwulan)
	return buf.Bytes(), filename, nil
}

// realisasiQualifierOrDash mengembalikan nilai string, atau "-" jika kosong.
func realisasiQualifierOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// =============================================================================
// generateSheetKpi — sheet kedua berisi daftar KPI dan Polarisasi dari DB
// =============================================================================

// generateSheetKpi membuat sheet "KPI" pada file Excel yang diberikan.
// Kolom A1: KPI, Kolom B1: Polarisasi.
// Data diambil dari mst_kpi LEFT JOIN mst_polarisasi.
// Jika polarisasi tidak ditemukan di mst_polarisasi, kolom B dikosongkan.
func (s *templateService) generateSheetKpi(f *excelize.File) error {
	const kpiSheetName = "KPI"

	// Tambahkan sheet baru "KPI"
	if _, err := f.NewSheet(kpiSheetName); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat sheet KPI: %v", err)}
	}

	// Ambil data dari DB
	kpiRows, err := s.repo.GetKpiWithPolarisasi()
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal ambil data mst_kpi: %v", err)}
	}

	// Style header sheet KPI — background biru muda + bold
	styleHeader, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"BDD7EE"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style header sheet KPI: %v", err)}
	}

	// Style data sheet KPI — border tipis
	styleData, err := f.NewStyle(&excelize.Style{
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style data sheet KPI: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: Header — A1 = "KPI", B1 = "Polarisasi"
	// -------------------------------------------------------------------------
	headers := []string{"KPI", "Polarisasi"}
	for colIdx, header := range headers {
		cellName, _ := excelize.CoordinatesToCellName(colIdx+1, 1)
		if err := f.SetCellValue(kpiSheetName, cellName, header); err != nil {
			return &errors.InternalServerError{Message: fmt.Sprintf("gagal set header %s sheet KPI: %v", cellName, err)}
		}
		if err := f.SetCellStyle(kpiSheetName, cellName, cellName, styleHeader); err != nil {
			return &errors.InternalServerError{Message: fmt.Sprintf("gagal set style header %s sheet KPI: %v", cellName, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Row 2 dst: Isi data KPI dan Polarisasi
	// -------------------------------------------------------------------------
	for i, row := range kpiRows {
		rowNum := i + 2 // data mulai row 2 (setelah header row 1)

		cellKpi, _ := excelize.CoordinatesToCellName(1, rowNum)
		cellPolarisasi, _ := excelize.CoordinatesToCellName(2, rowNum)

		if err := f.SetCellValue(kpiSheetName, cellKpi, row.Kpi); err != nil {
			return &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai KPI baris %d sheet KPI: %v", rowNum, err)}
		}
		if err := f.SetCellValue(kpiSheetName, cellPolarisasi, row.Polarisasi); err != nil {
			return &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai Polarisasi baris %d sheet KPI: %v", rowNum, err)}
		}
		if err := f.SetCellStyle(kpiSheetName, cellKpi, cellPolarisasi, styleData); err != nil {
			return &errors.InternalServerError{Message: fmt.Sprintf("gagal set style data baris %d sheet KPI: %v", rowNum, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Set lebar kolom sheet KPI
	// -------------------------------------------------------------------------
	if err := f.SetColWidth(kpiSheetName, "A", "A", 40); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom A sheet KPI: %v", err)}
	}
	if err := f.SetColWidth(kpiSheetName, "B", "B", 20); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom B sheet KPI: %v", err)}
	}

	// Set tinggi row header sheet KPI
	if err := f.SetRowHeight(kpiSheetName, 1, 30); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi header sheet KPI: %v", err)}
	}

	return nil
}

// =============================================================================
// Helper
// =============================================================================

// strPtr mengembalikan pointer ke string — pengganti excelize.Ptr yang tidak tersedia di v2.10.1.
func strPtr(s string) *string {
	return &s
}

// borderStyle mengembalikan konfigurasi border tipis untuk semua sisi cell.
func borderStyle() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}

// parseFloatOrString mencoba parse string sebagai float64.
// Jika berhasil, mengembalikan float64 agar Excel menyimpan sebagai angka.
// Jika gagal (atau string kosong), mengembalikan string aslinya.
func parseFloatOrString(s string) interface{} {
	if s == "" {
		return ""
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return v
}
