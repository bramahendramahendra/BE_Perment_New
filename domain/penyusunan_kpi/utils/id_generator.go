package utils

import (
	"fmt"
	"time"
)

// GenerateIDPengajuan membuat ID unik untuk pengajuan KPI berdasarkan kostl, tahun, triwulan, dan timestamp saat ini.
func GenerateIDPengajuan(kostl, tahun, triwulan string) string {
	t := time.Now()
	timestamp := fmt.Sprintf("%02d%02d%02d%02d%02d%02d",
		t.Year()%100,
		int(t.Month()),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	)
	return kostl + tahun + triwulan + timestamp
}

// GenerateIDDetail membuat ID untuk baris detail KPI (format: <idPengajuan>P<index>).
func GenerateIDDetail(idPengajuan string, index int) string {
	return fmt.Sprintf("%sP%03d", idPengajuan, index+1)
}

// GenerateIDSubDetail membuat ID untuk baris sub-detail KPI (format: <idPengajuan>C<index>).
func GenerateIDSubDetail(idPengajuan string, globalIndex int) string {
	return fmt.Sprintf("%sC%03d", idPengajuan, globalIndex)
}
