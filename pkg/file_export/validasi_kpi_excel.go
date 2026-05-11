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
	// A=No, B=KPI, C=KPI Qualifier, D=Bobot, E=Target TW, F=Target Qualifier,
	// G=Realisasi TW, H=Realisasi Qualifier, I=Pencapaian, J=Pencapaian Qualifier KPI, K=Post Qualifier
	colWidths := map[string]float64{
		"A": 6,
		"B": 35,
		"C": 22,
		"D": 12,
		"E": 18,
		"F": 20,
		"G": 18,
		"H": 20,
		"I": 16,
		"J": 24,
		"K": 26,
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", fmt.Errorf("gagal set lebar kolom %s: %w", col, err)
		}
	}

	// Baris 1: Nama Divisi + Pencapaian
	title := buildTitle(exportData)
	if err := f.MergeCell(sheetName, "A1", "K1"); err != nil {
		return nil, "", err
	}
	if err := f.SetCellValue(sheetName, "A1", title); err != nil {
		return nil, "", err
	}

	// Baris 2: Tahun - Triwulan
	subtitle := fmt.Sprintf("Tahun %s - TW %s", exportData.Tahun, exportData.TriwulanNum)
	if err := f.MergeCell(sheetName, "A2", "K2"); err != nil {
		return nil, "", err
	}
	if err := f.SetCellValue(sheetName, "A2", subtitle); err != nil {
		return nil, "", err
	}

	// Style: border tipis
	borderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{WrapText: true, Vertical: "top"},
	})
	if err != nil {
		return nil, "", fmt.Errorf("gagal membuat border style: %w", err)
	}

	// Style: header (biru tua, teks putih, bold)
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
		return nil, "", fmt.Errorf("gagal membuat header style: %w", err)
	}

	const numCols = 11
	applyStyleRow := func(rowNum, styleID int) error {
		for colIdx := 1; colIdx <= numCols; colIdx++ {
			cell, _ := excelize.CoordinatesToCellName(colIdx, rowNum)
			if err := f.SetCellStyle(sheetName, cell, cell, styleID); err != nil {
				return fmt.Errorf("gagal set style baris %d kolom %d: %w", rowNum, colIdx, err)
			}
		}
		return nil
	}

	// Baris 4: Header
	twNum := exportData.TriwulanNum
	headers := []string{
		"No", "KPI", "KPI Qualifier", "Bobot (%)",
		"Target Triwulan " + twNum, "Target Qualifier Triwulan",
		"Realisasi Triwulan " + twNum, "Realisasi Qualifier",
		"Pencapaian", "Pencapaian Qualifier KPI", "Pencapaian KPI Post Qualifier",
	}
	if err := f.SetRowHeight(sheetName, 4, 30); err != nil {
		return nil, "", err
	}
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 4)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", fmt.Errorf("gagal menulis header '%s': %w", header, err)
		}
	}
	if err := applyStyleRow(4, headerStyle); err != nil {
		return nil, "", err
	}

	// Baris data
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
		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return nil, "", fmt.Errorf("gagal menulis data baris %d kolom %d: %w", rowNum, colIdx+1, err)
			}
		}
		if err := applyStyleRow(rowNum, borderStyle); err != nil {
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

func buildTitle(exportData *dto.ValidasiKpiExportData) string {
	prefix := ""
	if exportData.IsDraft {
		prefix = "Draft "
	}
	return fmt.Sprintf("%s%s (Pencapaian: %s%%)", prefix, exportData.NamaDivisi, exportData.TotalPencapaian)
}
