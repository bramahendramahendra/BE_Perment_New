package repo

import (
	"fmt"

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

	// queryGetRevisionHeader mengambil header data_kpi berdasarkan id_pengajuan.
	queryGetRevisionHeader = `
		SELECT triwulan, tahun, kostl_tx
		FROM data_kpi
		WHERE id_pengajuan = ?
		LIMIT 1`

	// queryGetRevisionRows mengambil seluruh baris sub KPI beserta data extended (TW2/TW4)
	// untuk keperluan generate Excel tolakan penyusunan KPI.
	//
	// Mapping field:
	//   kpi_nama                    → data_kpi_detail.kpi                    (kolom B Excel)
	//   sub_kpi                     → data_kpi_subdetail.kpi                 (kolom C Excel)
	//   polarisasi                  → mst_polarisasi.polarisasi via rumus    (kolom D Excel)
	//   capping                     → data_kpi_subdetail.capping             (kolom E Excel)
	//   bobot                       → data_kpi_subdetail.bobot               (kolom F Excel)
	//   deskripsi_glossary          → data_kpi_subdetail.deskripsi_glossary  (kolom G Excel)
	//   target_triwulan             → data_kpi_subdetail.target_triwulan     (kolom H Excel)
	//   target_kuantitatif_triwulan → data_kpi_subdetail.target_kuantitatif_triwulan (kolom I)
	//   target_tahunan              → data_kpi_subdetail.target_tahunan      (kolom J Excel)
	//   target_kuantitatif_tahunan  → data_kpi_subdetail.target_kuantitatif_tahunan  (kolom K)
	//   terdapat_qualifier          → data_kpi_subdetail.id_qualifier ("Ya"/"Tidak")  (kolom L)
	//   item_qualifier              → data_kpi_subdetail.item_qualifier               (kolom M)
	//   deskripsi_qualifier         → data_kpi_subdetail.deskripsi_qualifier          (kolom N)
	//   target_qualifier            → data_kpi_subdetail.target_qualifier             (kolom O)
	//   nama_result/deskripsi_result   → data_result_detail   (kolom P/Q, TW2/TW4)
	//   nama_process/deskripsi_process → data_method_detail   (kolom R/S, TW2/TW4)
	//   nama_context/deskripsi_context → data_challenge_detail (kolom T/U, TW2/TW4)
	//
	// JOIN data_result_detail, data_method_detail, data_challenge_detail menggunakan
	// id_detail_result/id_detail_method/id_detail_challenge = id_sub_detail
	// karena keduanya di-generate menggunakan GenerateIDSubDetail yang sama saat /validate.
	queryGetRevisionRows = `
		SELECT
			a.id_sub_detail,
			b.kpi                                                    AS kpi_nama,
			a.kpi                                                    AS sub_kpi,
			IFNULL(p.polarisasi, '')                                 AS polarisasi,
			IFNULL(a.capping, '')                                    AS capping,
			IFNULL(CAST(a.bobot AS CHAR), '')                        AS bobot,
			IFNULL(a.deskripsi_glossary, '')                         AS deskripsi_glossary,
			IFNULL(a.target_triwulan, '')                            AS target_triwulan,
			IFNULL(CAST(a.target_kuantitatif_triwulan AS CHAR), '')  AS target_kuantitatif_triwulan,
			IFNULL(a.target_tahunan, '')                             AS target_tahunan,
			IFNULL(CAST(a.target_kuantitatif_tahunan AS CHAR), '')   AS target_kuantitatif_tahunan,
			IFNULL(a.id_qualifier, '')                               AS terdapat_qualifier,
			IFNULL(a.item_qualifier, '')                             AS item_qualifier,
			IFNULL(a.deskripsi_qualifier, '')                        AS deskripsi_qualifier,
			IFNULL(a.target_qualifier, '')                           AS target_qualifier,
			IFNULL(rd.nama_result, '')                               AS nama_result,
			IFNULL(rd.deskripsi_result, '')                          AS deskripsi_result,
			IFNULL(md.nama_method, '')                               AS nama_method,
			IFNULL(md.deskripsi_method, '')                          AS deskripsi_method,
			IFNULL(cd.nama_challenge, '')                            AS nama_challenge,
			IFNULL(cd.deskripsi_challenge, '')                       AS deskripsi_challenge
		FROM data_kpi_subdetail a
		INNER JOIN data_kpi_detail b
			ON a.id_detail = b.id_detail
		LEFT JOIN mst_polarisasi p
			ON p.id_polarisasi = a.rumus
		LEFT JOIN data_result_detail rd
			ON rd.id_detail_result = a.id_sub_detail
			AND rd.id_pengajuan = a.id_pengajuan
		LEFT JOIN data_method_detail md
			ON md.id_detail_method = a.id_sub_detail
			AND md.id_pengajuan = a.id_pengajuan
		LEFT JOIN data_challenge_detail cd
			ON cd.id_detail_challenge = a.id_sub_detail
			AND cd.id_pengajuan = a.id_pengajuan
		WHERE a.id_pengajuan = ?
		ORDER BY a.id_sub_detail ASC`
)

func (r *templateRepo) GetKpiWithPolarisasi() ([]*model.MstKpiPolarisasi, error) {
	var templates []*model.MstKpiPolarisasi
	err := r.db.Raw(GetKpiWithPolarisasiQuery).Scan(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// =============================================================================
// GET TOLAKAN PENYUSUNAN KPI DATA
// =============================================================================

// GetRevisionPenyusunanKpiData mengambil header dan seluruh baris sub KPI
// dari DB berdasarkan id_pengajuan untuk keperluan generate Excel tolakan.
func (r *templateRepo) GetRevisionPenyusunanKpiData(idPengajuan string) (*model.RevisionExcelData, error) {

	// =========================================================================
	// 1. Ambil header (triwulan, tahun, kostl_tx)
	// =========================================================================
	type headerRaw struct {
		Triwulan string `gorm:"column:triwulan"`
		Tahun    string `gorm:"column:tahun"`
		KostlTx  string `gorm:"column:kostl_tx"`
	}
	var header headerRaw
	if err := r.db.Raw(queryGetRevisionHeader, idPengajuan).Scan(&header).Error; err != nil {
		return nil, fmt.Errorf("gagal mengambil header tolakan KPI: %w", err)
	}
	if header.Triwulan == "" {
		return nil, nil
	}

	// =========================================================================
	// 2. Ambil seluruh baris sub KPI
	// =========================================================================
	rows, err := r.db.Raw(queryGetRevisionRows, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data sub KPI tolakan: %w", err)
	}
	defer rows.Close()

	var subDetailRows []model.RevisionSubDetailRow
	for rows.Next() {
		var row model.RevisionSubDetailRow
		if err := rows.Scan(
			&row.IdSubDetail,
			&row.KpiNama,
			&row.SubKpi,
			&row.Polarisasi,
			&row.Capping,
			&row.Bobot,
			&row.DeskripsiGlossary,
			&row.TargetTriwulan,
			&row.TargetKuantitatifTriwulan,
			&row.TargetTahunan,
			&row.TargetKuantitatifTahunan,
			&row.TerdapatQualifier,
			&row.ItemQualifier,
			&row.DeskripsiQualifier,
			&row.TargetQualifier,
			&row.NamaResult,
			&row.DeskripsiResult,
			&row.NamaProcess,
			&row.DeskripsiProcess,
			&row.NamaContext,
			&row.DeskripsiContext,
		); err != nil {
			return nil, fmt.Errorf("gagal scan baris sub KPI tolakan: %w", err)
		}
		subDetailRows = append(subDetailRows, row)
	}

	return &model.RevisionExcelData{
		Triwulan: header.Triwulan,
		Tahun:    header.Tahun,
		KostlTx:  header.KostlTx,
		Rows:     subDetailRows,
	}, nil
}
