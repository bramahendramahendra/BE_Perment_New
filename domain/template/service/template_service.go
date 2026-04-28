package service

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

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
	// Cek apakah data KPI untuk kombinasi tahun+triwulan+kostl sudah ada di DB.
	// Status 70 (draft) dan 71 (tolakan) masih boleh download template baru karena data lama akan di-replace.
	// Status selain itu berarti data sudah diproses → tolak.
	status, found, err := s.repo.GetExistPenyusunanStatus(req.Tahun, req.Triwulan, req.Divisi.Kostl)
	if err != nil {
		return nil, "", err
	}
	if found && status != 70 && status != 71 {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf(
				"data KPI tahun %s triwulan %s untuk divisi %s sudah ada", req.Tahun, req.Triwulan, req.Divisi.KostlTx,
			),
		}
	}

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

	filename := fmt.Sprintf("Format Penyusunan KPI Aplikasi Performance Management %s %s.xlsx", req.Tahun, req.Triwulan)
	return buf.Bytes(), filename, nil
}

// =============================================================================
// GenerateRevisionPenyusunanKpi
// =============================================================================

func (s *templateService) GenerateRevisionPenyusunanKpi(req *dto.RevisionPenyusunanKpiRequest) ([]byte, string, error) {

	exists, err := s.repo.CheckDataExist(req.IdPengajuan, req.Divisi.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf("data KPI  tahun '%s' triwulan '%s' dengan id pengajuan '%s' untuk divisi '%s', tidak ditemukan", req.Tahun, req.Triwulan, req.IdPengajuan, req.Divisi.KostlTx),
		}
	}

	excelData, err := s.repo.GetPenyusunanKpiData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}

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

	filename := fmt.Sprintf("Revisi Penyusunan KPI Aplikasi Performance Management %s %s %s.xlsx", req.Triwulan, req.Tahun, req.Divisi.KostlTx)
	return buf.Bytes(), filename, nil
}

// =============================================================================
// GenerateFormatRealisasiKpi
// =============================================================================

// columnsRealisasiBase adalah header kolom A–N (sama untuk semua triwulan).
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
	"Link Dokumen Sumber",
}

// columnsRealisasiExtendedTW24 adalah header kolom O–Z (khusus TW2 dan TW4).
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
	exists, err := s.repo.CheckDataExist(req.IdPengajuan, req.Divisi.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf("data KPI  tahun '%s' triwulan '%s' dengan id pengajuan '%s' untuk divisi '%s', tidak ditemukan", req.Tahun, req.Triwulan, req.IdPengajuan, req.Divisi.KostlTx),
		}
	}

	excelData, err := s.repo.GetPenyusunanKpiData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
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
	// TW1/TW3: A–N; TW2/TW4: diperluas hingga Z
	allColumns := make([]string, len(columnsRealisasiBase))
	copy(allColumns, columnsRealisasiBase)
	if isTW24 {
		allColumns = append(allColumns, columnsRealisasiExtendedTW24...)
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

	// Style kuning untuk sel yang datanya berasal dari DB (penyusunan KPI)
	styleDBData, err := f.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style DB data: %v", err)}
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

	// dbColSet: indeks kolom (0-based) yang datanya dari DB → warna kuning
	// TW1/TW3: B–I (idx 1–8)
	// TW2/TW4: B–I (idx 1–8) + O,P (14,15) + S,T (18,19) + W,X (22,23)
	dbColSet := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true}
	if isTW24 {
		dbColSet[14] = true
		dbColSet[15] = true
		dbColSet[18] = true
		dbColSet[19] = true
		dbColSet[22] = true
		dbColSet[23] = true
	}

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
				"", "", "", "",      // J–M: kosong (diisi user)
				"",                  // N: Link Dokumen Sumber (diisi user)
				row.NamaResult,      // O: Result
				row.DeskripsiResult, // P: Deskripsi Result
				"", "",              // Q–R: kosong (diisi user)
				row.NamaProcess,      // S: Process
				row.DeskripsiProcess, // T: Deskripsi Process
				"", "",               // U–V: kosong (diisi user)
				row.NamaContext,      // W: Context
				row.DeskripsiContext, // X: Deskripsi Context
				"", "",               // Y–Z: kosong (diisi user)
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
				"",             // N: Link Dokumen Sumber (diisi user)
			}
		}

		for colIdx, val := range values {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cellName, val); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set nilai baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
			cellStyle := styleData
			if dbColSet[colIdx] {
				cellStyle = styleDBData
			}
			if err := f.SetCellStyle(sheetName, cellName, cellName, cellStyle); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set style baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Legenda warna kuning di bawah tabel
	// -------------------------------------------------------------------------
	legendRow := realisasiDataStartRow + len(excelData.Rows) + 1

	styleYellowLegend, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style legenda kuning: %v", err)}
	}
	styleTextLegend, err := f.NewStyle(&excelize.Style{
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style teks legenda: %v", err)}
	}

	legendColorCell, _ := excelize.CoordinatesToCellName(1, legendRow)
	legendTextCell, _ := excelize.CoordinatesToCellName(2, legendRow)
	if err := f.SetCellValue(sheetName, legendColorCell, ""); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set legenda warna: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendColorCell, legendColorCell, styleYellowLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style legenda warna: %v", err)}
	}
	if err := f.SetCellValue(sheetName, legendTextCell, "Data yang didapat dari penyusunan KPI"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set teks legenda: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendTextCell, legendTextCell, styleTextLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style teks legenda: %v", err)}
	}

	// Legenda merah — kolom L/M tidak berlaku (tidak ada qualifier)
	legendRedRow := legendRow + 1
	styleRedLegend, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FF0000"},
			Pattern: 1,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style legenda merah: %v", err)}
	}
	legendRedColorCell, _ := excelize.CoordinatesToCellName(1, legendRedRow)
	legendRedTextCell, _ := excelize.CoordinatesToCellName(2, legendRedRow)
	if err := f.SetCellValue(sheetName, legendRedColorCell, ""); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set legenda merah: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendRedColorCell, legendRedColorCell, styleRedLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style legenda merah: %v", err)}
	}
	if err := f.SetCellValue(sheetName, legendRedTextCell, "Kolom Realisasi Qualifier dan Realisasi Qualifier Kuantitatif tidak berlaku (tidak ada qualifier)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set teks legenda merah: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendRedTextCell, legendRedTextCell, styleTextLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style teks legenda merah: %v", err)}
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

	// Kolom M (Realisasi Qualifier Kuantitatif) → Angka desimal (diisi user, seperti kolom K)
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Realisasi Qualifier Kuantitatif harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("M"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom M: %v", err)}
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
		"N": 45, // Link Dokumen Sumber
	}
	if isTW24 {
		colWidths["O"] = 25 // Result
		colWidths["P"] = 30 // Deskripsi Result
		colWidths["Q"] = 25 // Realisasi Result
		colWidths["R"] = 45 // Link Result
		colWidths["S"] = 25 // Process
		colWidths["T"] = 30 // Deskripsi Process
		colWidths["U"] = 25 // Realisasi Process
		colWidths["V"] = 45 // Link Process
		colWidths["W"] = 25 // Context
		colWidths["X"] = 30 // Deskripsi Context
		colWidths["Y"] = 25 // Realisasi Context
		colWidths["Z"] = 45 // Link Context
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
	//   Semua TW : J (Realisasi), K (Realisasi Kuantitatif),
	//             L (Realisasi Qualifier), M (Realisasi Qualifier Kuantitatif),
	//             N (Link Dokumen Sumber)
	//   TW2/TW4 : tambah Q (Realisasi Result), R (Link Result),
	//             U (Realisasi Process), V (Link Process),
	//             Y (Realisasi Context), Z (Link Context)
	// L dan M dikecualikan dari range unlock massal — diproses per-baris
	userInputCols := []string{"J", "K", "N"}
	if isTW24 {
		userInputCols = append(userInputCols, "Q", "R", "U", "V", "Y", "Z")
	}

	// Jumlah baris data aktual dari DB (baris 2 s.d. lastDataRow-1)
	totalDataRows := len(excelData.Rows)
	if totalDataRows > 0 {
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

		// Style merah terkunci untuk L/M ketika qualifier tidak berlaku
		styleLockedRed, err := f.NewStyle(&excelize.Style{
			Protection: &excelize.Protection{
				Locked: true,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"FF0000"},
				Pattern: 1,
			},
			Border: borderStyle(),
			Alignment: &excelize.Alignment{
				Vertical:   "top",
				Horizontal: "center",
				WrapText:   true,
			},
		})
		if err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style locked red: %v", err)}
		}

		dataEndRow := realisasiDataStartRow + totalDataRows - 1

		// Unlock kolom J, K (dan TW2/TW4: P,Q,T,U,X,Y) secara massal
		for _, col := range userInputCols {
			rangeRef := fmt.Sprintf("%s%d:%s%d", col, realisasiDataStartRow, col, dataEndRow)
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, realisasiDataStartRow),
				fmt.Sprintf("%s%d", col, dataEndRow), styleUnlocked); err != nil {
				return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style unlocked %s: %v", rangeRef, err)}
			}
		}

		// Kolom L dan M diproses per-baris berdasarkan qualifier
		for rowIdx, row := range excelData.Rows {
			rowNum := realisasiDataStartRow + rowIdx
			hasQualifier := strings.EqualFold(strings.TrimSpace(row.TerdapatQualifier), "ya")

			for _, colNum := range []int{12, 13} { // L=12, M=13
				cellName, _ := excelize.CoordinatesToCellName(colNum, rowNum)
				if hasQualifier {
					if err := f.SetCellStyle(sheetName, cellName, cellName, styleUnlocked); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style unlock L/M baris %d: %v", rowNum, err)}
					}
				} else {
					if err := f.SetCellValue(sheetName, cellName, "-"); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai - L/M baris %d: %v", rowNum, err)}
					}
					if err := f.SetCellStyle(sheetName, cellName, cellName, styleLockedRed); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style merah L/M baris %d: %v", rowNum, err)}
					}
				}
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

	filename := fmt.Sprintf("Format Realisasi KPI Aplikasi Performance Management %s %s %s.xlsx", req.Divisi.KostlTx, req.Tahun, req.Triwulan)
	return buf.Bytes(), filename, nil
}

func (s *templateService) GenerateRevisionRealisasiKpi(req *dto.RevisionRealisasiKpiRequest) ([]byte, string, error) {
	exists, err := s.repo.CheckRevisiRealisasiExist(req.IdPengajuan, req.Divisi.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return nil, "", err
	}
	if !exists {
		return nil, "", &errors.BadRequestError{
			Message: fmt.Sprintf("data realisasi KPI tahun '%s' triwulan '%s' dengan id pengajuan '%s' untuk divisi '%s', tidak ditemukan atau tidak dalam status yang dapat direvisi", req.Tahun, req.Triwulan, req.IdPengajuan, req.Divisi.KostlTx),
		}
	}

	excelData, err := s.repo.GetRealisasiKpiData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}

	isTW24 := req.Triwulan == "TW2" || req.Triwulan == "TW4"

	sheetName := req.Triwulan

	f := excelize.NewFile()
	defer f.Close()

	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nama sheet: %v", err)}
	}

	allColumns := make([]string, len(columnsRealisasiBase))
	copy(allColumns, columnsRealisasiBase)
	if isTW24 {
		allColumns = append(allColumns, columnsRealisasiExtendedTW24...)
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

	// Style kuning untuk sel dari DB (penyusunan KPI)
	styleDBData, err := f.NewStyle(&excelize.Style{
		Protection: &excelize.Protection{
			Locked: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style DB data: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Row 1: Header kolom
	// -------------------------------------------------------------------------
	const revRealisasiHeaderRow = 1
	const revRealisasiDataStartRow = 2
	const revRealisasiDataEndRow = 100

	for colIdx, header := range allColumns {
		cellName, _ := excelize.CoordinatesToCellName(colIdx+1, revRealisasiHeaderRow)
		if err := f.SetCellValue(sheetName, cellName, header); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set header %s: %v", cellName, err)}
		}
		if err := f.SetCellStyle(sheetName, cellName, cellName, styleHeader); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style header %s: %v", cellName, err)}
		}
	}

	// -------------------------------------------------------------------------
	// Row 2+: Tulis data baris dari DB (pre-filled termasuk realisasi)
	// -------------------------------------------------------------------------

	// dbColSet: indeks kolom (0-based) yang datanya dari DB → warna kuning
	// TW1/TW3: B–I (idx 1–8)
	// TW2/TW4: B–I (idx 1–8) + O,P (14,15) + S,T (18,19) + W,X (22,23)
	dbColSet := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true}
	if isTW24 {
		dbColSet[14] = true
		dbColSet[15] = true
		dbColSet[18] = true
		dbColSet[19] = true
		dbColSet[22] = true
		dbColSet[23] = true
	}

	for rowIdx, row := range excelData.Rows {
		rowNum := revRealisasiDataStartRow + rowIdx

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
				row.Realisasi,                                  // J: Realisasi (pre-filled)
				parseFloatOrString(row.RealisasiKuantitatif),  // K: Realisasi Kuantitatif (pre-filled)
				row.RealisasiQualifier,                         // L: Realisasi Qualifier (pre-filled, diproses per-baris)
				row.RealisasiKuantitatifQualifier,              // M: Realisasi Qualifier Kuantitatif (pre-filled)
				row.LinkDokumenSumber,                          // N: Link Dokumen Sumber (pre-filled)
				row.NamaResult,      // O: Result
				row.DeskripsiResult, // P: Deskripsi Result
				row.RealisasiResult, // Q: Realisasi Result (pre-filled)
				row.LinkResult,      // R: Link Result (pre-filled)
				row.NamaProcess,      // S: Process
				row.DeskripsiProcess, // T: Deskripsi Process
				row.RealisasiProcess, // U: Realisasi Process (pre-filled)
				row.LinkProcess,      // V: Link Process (pre-filled)
				row.NamaContext,      // W: Context
				row.DeskripsiContext, // X: Deskripsi Context
				row.RealisasiContext, // Y: Realisasi Context (pre-filled)
				row.LinkContext,      // Z: Link Context (pre-filled)
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
				row.Realisasi,                                  // J: Realisasi (pre-filled)
				parseFloatOrString(row.RealisasiKuantitatif),  // K: Realisasi Kuantitatif (pre-filled)
				row.RealisasiQualifier,                         // L: Realisasi Qualifier (pre-filled, diproses per-baris)
				row.RealisasiKuantitatifQualifier,              // M: Realisasi Qualifier Kuantitatif (pre-filled)
				row.LinkDokumenSumber,                          // N: Link Dokumen Sumber (pre-filled)
			}
		}

		for colIdx, val := range values {
			cellName, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cellName, val); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set nilai baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
			cellStyle := styleData
			if dbColSet[colIdx] {
				cellStyle = styleDBData
			}
			if err := f.SetCellStyle(sheetName, cellName, cellName, cellStyle); err != nil {
				return nil, "", &errors.InternalServerError{
					Message: fmt.Sprintf("gagal set style baris %d kolom %d: %v", rowNum, colIdx+1, err),
				}
			}
		}
	}

	// -------------------------------------------------------------------------
	// Legenda warna kuning di bawah tabel
	// -------------------------------------------------------------------------
	legendRow := revRealisasiDataStartRow + len(excelData.Rows) + 1

	styleYellowLegend, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FFFF00"},
			Pattern: 1,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style legenda kuning: %v", err)}
	}
	styleTextLegend, err := f.NewStyle(&excelize.Style{
		Border: borderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style teks legenda: %v", err)}
	}

	legendColorCell, _ := excelize.CoordinatesToCellName(1, legendRow)
	legendTextCell, _ := excelize.CoordinatesToCellName(2, legendRow)
	if err := f.SetCellValue(sheetName, legendColorCell, ""); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set legenda warna: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendColorCell, legendColorCell, styleYellowLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style legenda warna: %v", err)}
	}
	if err := f.SetCellValue(sheetName, legendTextCell, "Data yang didapat dari penyusunan KPI"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set teks legenda: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendTextCell, legendTextCell, styleTextLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style teks legenda: %v", err)}
	}

	legendRedRow := legendRow + 1
	styleRedLegend, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"FF0000"},
			Pattern: 1,
		},
		Border: borderStyle(),
	})
	if err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style legenda merah: %v", err)}
	}
	legendRedColorCell, _ := excelize.CoordinatesToCellName(1, legendRedRow)
	legendRedTextCell, _ := excelize.CoordinatesToCellName(2, legendRedRow)
	if err := f.SetCellValue(sheetName, legendRedColorCell, ""); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set legenda merah: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendRedColorCell, legendRedColorCell, styleRedLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style legenda merah: %v", err)}
	}
	if err := f.SetCellValue(sheetName, legendRedTextCell, "Kolom Realisasi Qualifier dan Realisasi Qualifier Kuantitatif tidak berlaku (tidak ada qualifier)"); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set teks legenda merah: %v", err)}
	}
	if err := f.SetCellStyle(sheetName, legendRedTextCell, legendRedTextCell, styleTextLegend); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style teks legenda merah: %v", err)}
	}

	// -------------------------------------------------------------------------
	// Data Validation
	// -------------------------------------------------------------------------
	sqrefDataRange := func(col string) string {
		return fmt.Sprintf("%s%d:%s%d", col, revRealisasiDataStartRow, col, revRealisasiDataEndRow)
	}

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

	// Kolom M (Realisasi Qualifier Kuantitatif) → Angka desimal (diisi user, seperti kolom K)
	if err := f.AddDataValidation(sheetName, &excelize.DataValidation{
		Type:             "decimal",
		Operator:         "greaterThanOrEqual",
		Formula1:         "0",
		ShowErrorMessage: true,
		ErrorStyle:       strPtr("stop"),
		ErrorTitle:       strPtr("Input Tidak Valid"),
		Error:            strPtr("Realisasi Qualifier Kuantitatif harus berupa angka (maks. 2 angka di belakang koma)."),
		Sqref:            sqrefDataRange("M"),
	}); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal tambah validasi kolom M: %v", err)}
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
		"N": 45, // Link Dokumen Sumber
	}
	if isTW24 {
		colWidths["O"] = 25 // Result
		colWidths["P"] = 30 // Deskripsi Result
		colWidths["Q"] = 25 // Realisasi Result
		colWidths["R"] = 45 // Link Result
		colWidths["S"] = 25 // Process
		colWidths["T"] = 30 // Deskripsi Process
		colWidths["U"] = 25 // Realisasi Process
		colWidths["V"] = 45 // Link Process
		colWidths["W"] = 25 // Context
		colWidths["X"] = 30 // Deskripsi Context
		colWidths["Y"] = 25 // Realisasi Context
		colWidths["Z"] = 45 // Link Context
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom %s: %v", col, err)}
		}
	}

	if err := f.SetRowHeight(sheetName, revRealisasiHeaderRow, 40); err != nil {
		return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi row header: %v", err)}
	}

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
	userInputCols := []string{"J", "K"}
	if isTW24 {
		userInputCols = append(userInputCols, "P", "Q", "T", "U", "X", "Y")
	}

	totalDataRows := len(excelData.Rows)
	if totalDataRows > 0 {
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

		styleLockedRed, err := f.NewStyle(&excelize.Style{
			Protection: &excelize.Protection{
				Locked: true,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"FF0000"},
				Pattern: 1,
			},
			Border: borderStyle(),
			Alignment: &excelize.Alignment{
				Vertical:   "top",
				Horizontal: "center",
				WrapText:   true,
			},
		})
		if err != nil {
			return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style locked red: %v", err)}
		}

		dataEndRow := revRealisasiDataStartRow + totalDataRows - 1

		for _, col := range userInputCols {
			rangeRef := fmt.Sprintf("%s%d:%s%d", col, revRealisasiDataStartRow, col, dataEndRow)
			if err := f.SetCellStyle(sheetName, fmt.Sprintf("%s%d", col, revRealisasiDataStartRow),
				fmt.Sprintf("%s%d", col, dataEndRow), styleUnlocked); err != nil {
				return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style unlocked %s: %v", rangeRef, err)}
			}
		}

		// Kolom L dan M diproses per-baris berdasarkan qualifier
		for rowIdx, row := range excelData.Rows {
			rowNum := revRealisasiDataStartRow + rowIdx
			hasQualifier := strings.EqualFold(strings.TrimSpace(row.TerdapatQualifier), "ya")

			for _, colNum := range []int{12, 13} { // L=12, M=13
				cellName, _ := excelize.CoordinatesToCellName(colNum, rowNum)
				if hasQualifier {
					if err := f.SetCellStyle(sheetName, cellName, cellName, styleUnlocked); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style unlock L/M baris %d: %v", rowNum, err)}
					}
				} else {
					if err := f.SetCellValue(sheetName, cellName, "-"); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set nilai - L/M baris %d: %v", rowNum, err)}
					}
					if err := f.SetCellStyle(sheetName, cellName, cellName, styleLockedRed); err != nil {
						return nil, "", &errors.InternalServerError{Message: fmt.Sprintf("gagal set style merah L/M baris %d: %v", rowNum, err)}
					}
				}
			}
		}
	}

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

	filename := fmt.Sprintf("Revisi Realisasi KPI Aplikasi Performance Management %s %s %s.xlsx", req.Divisi.KostlTx, req.Tahun, req.Triwulan)
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
