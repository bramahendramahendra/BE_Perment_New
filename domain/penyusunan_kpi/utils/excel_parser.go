package utils

import (
	"mime/multipart"

	"permen_api/domain/penyusunan_kpi/dto"
	"permen_api/pkg/excel"
)

// Konstanta triwulan — re-export dari pkg/excel agar consumer utils tidak perlu ganti import.
const (
	ExcelDataStartRow = excel.DataStartRow
	ExcelMaxDataRows  = excel.MaxDataRows

	SheetTW1 = excel.SheetTW1
	SheetTW2 = excel.SheetTW2
	SheetTW3 = excel.SheetTW3
	SheetTW4 = excel.SheetTW4

	TriwulanTW1 = excel.TriwulanTW1
	TriwulanTW2 = excel.TriwulanTW2
	TriwulanTW3 = excel.TriwulanTW3
	TriwulanTW4 = excel.TriwulanTW4
)

// GetMaxRowsFromEnv mendelegasikan ke pkg/excel.
func GetMaxRowsFromEnv() int {
	return excel.GetMaxRowsFromEnv()
}

// IsExtendedTriwulan mendelegasikan ke pkg/excel.
func IsExtendedTriwulan(triwulan string) bool {
	return excel.IsExtendedTriwulan(triwulan)
}

// ParseAndValidateExcel mendelegasikan ke pkg/excel.
func ParseAndValidateExcel(
	file *multipart.FileHeader,
	triwulan string,
) ([]dto.PenyusunanKpiRow, map[int][]dto.PenyusunanKpiSubDetailRow, error) {
	return excel.ParseAndValidateExcel(file, triwulan)
}
