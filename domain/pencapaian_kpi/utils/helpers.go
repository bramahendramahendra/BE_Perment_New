package utils

import (
	model "permen_api/domain/pencapaian_kpi/model"
)

// ResolveIndikator menentukan warna indikator berdasarkan nilai pct dan daftar indikator dari DB.
// Logika: iterasi semua item, last match wins. Default "merah" jika tidak ada match.
func ResolveIndikator(pct float64, indikator []*model.IndikatorPencapaian) string {
	warna := "merah"
	for _, item := range indikator {
		if pct <= item.IndikatorValue {
			warna = item.IndikatorWarna
		}
	}
	return warna
}
