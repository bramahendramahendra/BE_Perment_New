package excel

import (
	"fmt"
	"math"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadSheet membuka file Excel dan membaca semua baris dari sheet yang diminta.
// Mengembalikan raw [][]string untuk diproses oleh masing-masing domain.
func ReadSheet(file *multipart.FileHeader, sheetName string) ([][]string, error) {
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

	sheetIndex, err := xlsx.GetSheetIndex(sheetName)
	if err != nil || sheetIndex < 0 {
		return nil, fmt.Errorf("file Excel '%s' tidak memiliki sheet '%s'", file.Filename, sheetName)
	}

	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca baris sheet '%s': %w", sheetName, err)
	}

	return rows, nil
}

// GetCell mengambil nilai cell dari row secara aman (tidak panic jika index melebihi panjang row).
func GetCell(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

// ParseFloat mem-parse string menjadi float64 dengan presisi 2 desimal.
func ParseFloat(s string) (float64, error) {
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

// NullableString mengembalikan pointer string jika isActive true, nil jika false.
func NullableString(val string, isActive bool) *string {
	if !isActive {
		return nil
	}
	return &val
}
