package file_export

import (
	"fmt"
	"strconv"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"

	"github.com/xuri/excelize/v2"
)

// GenerateValidasiKpiExcel membuat file Excel dari ValidasiKpiExportData.
func GenerateValidasiKpiExcel(exportData *dto.ValidasiKpiExportData) ([]byte, string, error) {
	const sheetName = "Pencapaian KPI"

	ef, err := NewExcelFile(sheetName)
	if err != nil {
		return nil, "", err
	}
	defer ef.Close()

	f := ef.File()

	// Lebar kolom
	colWidths := map[string]float64{
		"A": 6,
		"B": 35,
		"C": 22,
		"D": 12,
		"E": 20,
		"F": 22,
		"G": 20,
		"H": 22,
		"I": 18,
		"J": 26,
		"K": 28,
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", fmt.Errorf("gagal set lebar kolom %s: %w", col, err)
		}
	}

	// --- Style definitions ---

	titleStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12, Color: "1F497D"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, "", err
	}

	subtitleStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "1F497D"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, "", err
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"1F497D"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Color: "FFFFFF", Size: 9},
		Border: []excelize.Border{
			{Type: "left", Color: "FFFFFF", Style: 1},
			{Type: "right", Color: "FFFFFF", Style: 1},
			{Type: "top", Color: "FFFFFF", Style: 1},
			{Type: "bottom", Color: "FFFFFF", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Font: &excelize.Font{Size: 9},
		Border: []excelize.Border{
			{Type: "left", Color: "B4B4B4", Style: 1},
			{Type: "right", Color: "B4B4B4", Style: 1},
			{Type: "top", Color: "B4B4B4", Style: 1},
			{Type: "bottom", Color: "B4B4B4", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	dataLeftStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Font: &excelize.Font{Size: 9},
		Border: []excelize.Border{
			{Type: "left", Color: "B4B4B4", Style: 1},
			{Type: "right", Color: "B4B4B4", Style: 1},
			{Type: "top", Color: "B4B4B4", Style: 1},
			{Type: "bottom", Color: "B4B4B4", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	indicatorGreenStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Size: 9, Color: "72AD4A"},
		Border: []excelize.Border{
			{Type: "left", Color: "B4B4B4", Style: 1},
			{Type: "right", Color: "B4B4B4", Style: 1},
			{Type: "top", Color: "B4B4B4", Style: 1},
			{Type: "bottom", Color: "B4B4B4", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	indicatorYellowStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Size: 9, Color: "F9C301"},
		Border: []excelize.Border{
			{Type: "left", Color: "B4B4B4", Style: 1},
			{Type: "right", Color: "B4B4B4", Style: 1},
			{Type: "top", Color: "B4B4B4", Style: 1},
			{Type: "bottom", Color: "B4B4B4", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	indicatorRedStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"FFFFFF"}, Pattern: 1},
		Font: &excelize.Font{Bold: true, Size: 9, Color: "C30002"},
		Border: []excelize.Border{
			{Type: "left", Color: "B4B4B4", Style: 1},
			{Type: "right", Color: "B4B4B4", Style: 1},
			{Type: "top", Color: "B4B4B4", Style: 1},
			{Type: "bottom", Color: "B4B4B4", Style: 1},
		},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	if err != nil {
		return nil, "", err
	}

	footerBoldStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Color: "1F497D"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, "", err
	}

	footerRegularStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 9, Color: "1F497D"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, "", err
	}

	footerItalicStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Italic: true, Size: 9, Color: "1F497D"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	if err != nil {
		return nil, "", err
	}

	// --- Baris 1: Judul ---
	title := buildTitle(exportData)
	if err := f.MergeCell(sheetName, "A1", "K1"); err != nil {
		return nil, "", err
	}
	if err := f.SetCellValue(sheetName, "A1", title); err != nil {
		return nil, "", err
	}
	if err := f.SetCellStyle(sheetName, "A1", "A1", titleStyle); err != nil {
		return nil, "", err
	}
	if err := f.SetRowHeight(sheetName, 1, 18); err != nil {
		return nil, "", err
	}

	// --- Baris 2: Subtitle ---
	subtitle := fmt.Sprintf("Tahun %s - TW %s", exportData.Tahun, exportData.TriwulanNum)
	if err := f.MergeCell(sheetName, "A2", "K2"); err != nil {
		return nil, "", err
	}
	if err := f.SetCellValue(sheetName, "A2", subtitle); err != nil {
		return nil, "", err
	}
	if err := f.SetCellStyle(sheetName, "A2", "A2", subtitleStyle); err != nil {
		return nil, "", err
	}
	if err := f.SetRowHeight(sheetName, 2, 16); err != nil {
		return nil, "", err
	}

	// --- Baris 3: kosong ---
	if err := f.SetRowHeight(sheetName, 3, 6); err != nil {
		return nil, "", err
	}

	// --- Baris 4: Header tabel ---
	twNum := exportData.TriwulanNum
	headers := []string{
		"No", "KPI", "KPI Qualifier", "Bobot (%)",
		"Target Triwulan " + twNum, "Target Qualifier Triwulan",
		"Realisasi Triwulan " + twNum, "Realisasi Qualifier",
		"Pencapaian", "Pencapaian Qualifier KPI", "Pencapaian KPI Post Qualifier",
	}
	if err := f.SetRowHeight(sheetName, 4, 32); err != nil {
		return nil, "", err
	}
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 4)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", fmt.Errorf("gagal menulis header '%s': %w", header, err)
		}
		if err := f.SetCellStyle(sheetName, cell, cell, headerStyle); err != nil {
			return nil, "", err
		}
	}

	// --- Baris data ---
	// Kolom kiri (B, C) pakai left-align; kolom indikator (I, J, K) pakai warna persen
	leftAlignCols := map[int]bool{2: true, 3: true} // kolom B=2, C=3 (1-based)
	indicatorCols := map[int]bool{9: true, 10: true, 11: true} // I=9, J=10, K=11

	for i, row := range exportData.Rows {
		rowNum := 5 + i
		values := []interface{}{
			strconv.Itoa(row.No),
			row.Kpi,
			row.ItemQualifier,
			row.Bobot,
			row.TargetTriwulan,
			row.TargetQualifier,
			row.RealisasiValidated,
			row.RealisasiQualifier,
			row.Pencapaian,
			row.PencapaianQualifier,
			row.PencapaianPostQualifier,
		}
		if err := f.SetRowHeight(sheetName, rowNum, 22); err != nil {
			return nil, "", err
		}
		for colIdx, val := range values {
			colNum := colIdx + 1
			cell, _ := excelize.CoordinatesToCellName(colNum, rowNum)
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return nil, "", fmt.Errorf("gagal menulis data baris %d kolom %d: %w", rowNum, colNum, err)
			}

			var styleID int
			switch {
			case indicatorCols[colNum]:
				styleID = indicatorStyleID(val.(string), exportData.Indikator, indicatorGreenStyle, indicatorYellowStyle, indicatorRedStyle, dataStyle)
			case leftAlignCols[colNum]:
				styleID = dataLeftStyle
			default:
				styleID = dataStyle
			}
			if err := f.SetCellStyle(sheetName, cell, cell, styleID); err != nil {
				return nil, "", err
			}
		}
	}

	// --- Footer ---
	footerStart := 5 + len(exportData.Rows) + 1
	footerRows := []struct {
		text    string
		styleID int
		height  float64
	}{
		{"PT Bank Rakyat Indonesia (Persero) Tbk", footerBoldStyle, 16},
		{"Planning, Budgeting & Performance Management Group", footerBoldStyle, 14},
		{"Gedung BRI II. Jalan Jendral Sudirman Kav. 44-46. Jakarta, Indonesia 10210", footerRegularStyle, 14},
		{"", footerRegularStyle, 8},
		{"Integrity, Collaborative, Accountability, Growth Mindset, Customer Focus", footerItalicStyle, 14},
	}
	for j, fr := range footerRows {
		rowNum := footerStart + j
		cellRef, _ := excelize.CoordinatesToCellName(1, rowNum)
		if err := f.MergeCell(sheetName, cellRef, func() string {
			end, _ := excelize.CoordinatesToCellName(11, rowNum)
			return end
		}()); err != nil {
			return nil, "", err
		}
		if err := f.SetCellValue(sheetName, cellRef, fr.text); err != nil {
			return nil, "", err
		}
		if err := f.SetCellStyle(sheetName, cellRef, cellRef, fr.styleID); err != nil {
			return nil, "", err
		}
		if err := f.SetRowHeight(sheetName, rowNum, fr.height); err != nil {
			return nil, "", err
		}
	}

	fileBytes, err := ef.ToBytes()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("Pencapaian_KPI_%s_%s_TW%s.xlsx",
		strings.ReplaceAll(exportData.NamaDivisi, " ", "_"),
		exportData.Tahun,
		exportData.TriwulanNum,
	)
	return fileBytes, filename, nil
}

// indicatorStyleID mengembalikan style warna berdasarkan nilai persen dan indikator dari DB.
// Logika: iterasi indikator (descending by value), warna terakhir yang memenuhi pct <= value dipakai.
func indicatorStyleID(val string, indikator []dto.IndikatorPencapaian, greenStyle, yellowStyle, redStyle, defaultStyle int) int {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" || trimmed == "-" {
		return defaultStyle
	}
	pctStr := strings.TrimSuffix(trimmed, "%")
	var pct float64
	fmt.Sscanf(pctStr, "%f", &pct)

	styleID := defaultStyle
	for _, item := range indikator {
		if pct <= item.Value {
			switch item.Warna {
			case "hijau":
				styleID = greenStyle
			case "kuning":
				styleID = yellowStyle
			case "merah":
				styleID = redStyle
			}
		}
	}
	return styleID
}

func buildTitle(exportData *dto.ValidasiKpiExportData) string {
	prefix := ""
	if exportData.IsDraft {
		prefix = "Draft "
	}
	return fmt.Sprintf("%s%s (Pencapaian: %s%%)", prefix, exportData.NamaDivisi, exportData.TotalPencapaian)
}
