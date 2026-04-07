package service

import (
	"bytes"
	"fmt"

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
// Implementasi service
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
		Font: &excelize.Font{
			Bold: true,
		},
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
