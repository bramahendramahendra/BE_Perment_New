package utils

import (
	"strconv"

	"github.com/xuri/excelize/v2"
)

// StrPtr mengembalikan pointer ke string.
func StrPtr(s string) *string {
	return &s
}

// BorderStyle mengembalikan konfigurasi border tipis untuk semua sisi cell.
func BorderStyle() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
	}
}

// AppendPercent menambahkan simbol "%" di akhir string, atau string kosong jika input kosong.
func AppendPercent(s string) string {
	if s == "" {
		return ""
	}
	return s + "%"
}

// ParseFloatOrString mencoba parse string sebagai float64.
// Jika berhasil, mengembalikan float64 agar Excel menyimpan sebagai angka.
// Jika gagal atau kosong, mengembalikan string aslinya.
func ParseFloatOrString(s string) interface{} {
	if s == "" {
		return ""
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return v
}

// RealisasiQualifierOrDash mengembalikan nilai string, atau "-" jika kosong.
func RealisasiQualifierOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
