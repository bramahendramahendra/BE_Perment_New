package repo

import (
	model "permen_api/domain/template/model"
)

const (
	// GetKpiWithPolarisasiQuery mengambil semua KPI dari mst_kpi beserta polarisasi dari mst_polarisasi.
	// LEFT JOIN digunakan agar mst_kpi yang rumus-nya tidak ada di mst_polarisasi tetap dimunculkan
	// dengan kolom polarisasi bernilai string kosong (IFNULL).
	// Join condition: mst_polarisasi.id_polarisasi = mst_kpi.rumus
	GetKpiWithPolarisasiQuery = `
		SELECT
			m.kpi,
			IFNULL(p.polarisasi, '') AS polarisasi
		FROM mst_kpi m
		LEFT JOIN mst_polarisasi p ON p.id_polarisasi = m.rumus
		ORDER BY m.id_kpi ASC`
)

func (r *templateRepo) GetKpiWithPolarisasi() ([]*model.MstKpiPolarisasi, error) {
	var templates []*model.MstKpiPolarisasi
	err := r.db.Raw(GetKpiWithPolarisasiQuery).Scan(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}
