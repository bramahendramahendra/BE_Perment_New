package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// NullableString mengembalikan pointer string jika isActive true, nil jika false.
// Nilai nil disimpan sebagai NULL di DB (digunakan untuk kolom P–U sheet "TW 4").
func NullableString(val string, isActive bool) *string {
	if !isActive {
		return nil
	}
	return &val
}

// ParseFloat2Decimal mem-parse string menjadi float64 dengan presisi 2 desimal.
func ParseFloat2Decimal(s string) (float64, error) {
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
