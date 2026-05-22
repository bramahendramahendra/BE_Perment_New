package service

import (
	"fmt"

	"permen_api/domain/template/utils"
	"permen_api/errors"

	"github.com/xuri/excelize/v2"
)

// generateSheetKpi membuat sheet "KPI" pada file Excel yang diberikan.
// Kolom A1: KPI, Kolom B1: Polarisasi.
// Data diambil dari mst_kpi LEFT JOIN mst_polarisasi.
func (s *templateService) generateSheetKpi(f *excelize.File) error {
	const kpiSheetName = "KPI"

	if _, err := f.NewSheet(kpiSheetName); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat sheet KPI: %v", err)}
	}

	kpiRows, err := s.repo.GetKpiWithPolarisasi()
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal ambil data mst_kpi: %v", err)}
	}

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
		Border: utils.BorderStyle(),
	})
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style header sheet KPI: %v", err)}
	}

	styleData, err := f.NewStyle(&excelize.Style{
		Border: utils.BorderStyle(),
		Alignment: &excelize.Alignment{
			Vertical: "top",
			WrapText: true,
		},
	})
	if err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal buat style data sheet KPI: %v", err)}
	}

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

	for i, row := range kpiRows {
		rowNum := i + 2

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

	if err := f.SetColWidth(kpiSheetName, "A", "A", 40); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom A sheet KPI: %v", err)}
	}
	if err := f.SetColWidth(kpiSheetName, "B", "B", 20); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set lebar kolom B sheet KPI: %v", err)}
	}
	if err := f.SetRowHeight(kpiSheetName, 1, 30); err != nil {
		return &errors.InternalServerError{Message: fmt.Sprintf("gagal set tinggi header sheet KPI: %v", err)}
	}

	return nil
}
