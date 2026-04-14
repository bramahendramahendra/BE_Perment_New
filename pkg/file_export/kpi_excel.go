package file_export

import (
	"fmt"
	"strconv"

	dto "permen_api/domain/penyusunan_kpi/dto"

	"github.com/xuri/excelize/v2"
)

// GenerateKpiExcel membuat file Excel dari KpiExportData dan mengembalikan byte slice
// serta nama file yang direkomendasikan untuk header Content-Disposition.
func GenerateKpiExcel(exportData *dto.KpiExportData) ([]byte, string, error) {
	const sheetName = "Data KPI"

	ef, err := NewExcelFile(sheetName)
	if err != nil {
		return nil, "", err
	}
	defer ef.Close()

	f := ef.File()

	// -------------------------------------------------------------------------
	// Lebar kolom
	// A  = No             → sempit
	// B  = KPI            → lebar (konten teks panjang)
	// C  = Bobot (%)      → sedang
	// D  = Target Tahunan → sedang
	// E  = Capping        → sedang
	// -------------------------------------------------------------------------
	colWidths := map[string]float64{
		"A": 6,
		"B": 40,
		"C": 14,
		"D": 20,
		"E": 14,
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", fmt.Errorf("gagal set lebar kolom %s: %w", col, err)
		}
	}

	// -------------------------------------------------------------------------
	// Baris 1: Nama Divisi — merge A1:E1
	// -------------------------------------------------------------------------
	if err := f.MergeCell(sheetName, "A1", "E1"); err != nil {
		return nil, "", fmt.Errorf("gagal merge cell baris 1: %w", err)
	}
	if err := f.SetCellValue(sheetName, "A1", exportData.NamaDivisi); err != nil {
		return nil, "", fmt.Errorf("gagal menulis nama divisi: %w", err)
	}

	// -------------------------------------------------------------------------
	// Baris 2: Tahun — merge A2:E2
	// -------------------------------------------------------------------------
	if err := f.MergeCell(sheetName, "A2", "E2"); err != nil {
		return nil, "", fmt.Errorf("gagal merge cell baris 2: %w", err)
	}
	if err := f.SetCellValue(sheetName, "A2", "Tahun "+exportData.Tahun); err != nil {
		return nil, "", fmt.Errorf("gagal menulis tahun: %w", err)
	}

	// -------------------------------------------------------------------------
	// Baris 3: Kosong
	// -------------------------------------------------------------------------
	// (tidak perlu set value — cell default kosong)

	// -------------------------------------------------------------------------
	// Baris 4: Header kolom
	// -------------------------------------------------------------------------
	headers := []string{"No", "KPI", "Bobot (%)", "Target Tahunan", "Capping"}
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 4)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", fmt.Errorf("gagal menulis header kolom '%s': %w", header, err)
		}
	}

	// -------------------------------------------------------------------------
	// Baris 5 dst: Data rows
	// -------------------------------------------------------------------------
	for i, row := range exportData.Rows {
		rowNum := 5 + i

		values := []interface{}{
			strconv.Itoa(row.No),
			row.KpiNama,
			row.Bobot,
			row.TargetTahunan,
			row.Capping,
		}

		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return nil, "", fmt.Errorf("gagal menulis data baris %d kolom %d: %w", rowNum, colIdx+1, err)
			}
		}
	}

	fileBytes, err := ef.ToBytes()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("KPI_%s_%s_%s.xlsx",
		exportData.NamaDivisi, exportData.Tahun, exportData.Triwulan)

	return fileBytes, filename, nil
}
