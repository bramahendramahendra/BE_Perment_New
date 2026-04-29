package repo

import (
	"database/sql"
	"errors"
	"fmt"

	model "permen_api/domain/template/model"
)

const (
	// queryGetExistPenyusunanStatus mengecek keberadaan record di data_kpi
	// berdasarkan tahun, triwulan, dan kostl.
	queryGetExistPenyusunanStatus = `
		SELECT status
		FROM data_kpi
		WHERE tahun = ? AND triwulan = ? AND kostl = ?
		LIMIT 1`

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

	// queryGetRows mengambil seluruh baris sub KPI beserta data extended (TW2/TW4)
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
	// queryExistData mengecek keberadaan record data_kpi
	// berdasarkan id_pengajuan, kostl, tahun, dan triwulan.
	queryExistData = `
		SELECT COUNT(1)
		FROM data_kpi
		WHERE id_pengajuan = ? AND kostl = ? AND tahun = ? AND triwulan = ?
		LIMIT 1`

	// queryCheckRevisiRealisasi mengecek keberadaan data_kpi yang dapat direvisi realisasi-nya
	// (status 4 = tolak realisasi, 80 = draft realisasi).
	queryCheckRevisiRealisasi = `
		SELECT COUNT(1)
		FROM data_kpi
		WHERE id_pengajuan = ? AND kostl = ? AND tahun = ? AND triwulan = ? AND status IN (4, 80)
		LIMIT 1`

	// queryGetRealisasiRows mengambil seluruh baris realisasi KPI termasuk data realisasi yang sudah diisi user.
	// Digunakan untuk generate Excel revisi realisasi KPI (revision-realisasi-kpi).
	queryGetRealisasiRows = `
		SELECT
			b.kpi                                                        AS kpi_nama,
			a.kpi                                                        AS sub_kpi,
			IFNULL(p.polarisasi, '')                                     AS polarisasi,
			IFNULL(a.capping, '')                                        AS capping,
			IFNULL(CAST(a.bobot AS CHAR), '')                            AS bobot,
			IFNULL(a.target_triwulan, '')                                AS target_triwulan,
			IFNULL(a.item_qualifier, '')                                 AS item_qualifier,
			IFNULL(a.target_qualifier, '')                               AS target_qualifier,
			IFNULL(a.id_qualifier, '')                                   AS terdapat_qualifier,
			IFNULL(a.realisasi, '')                                      AS realisasi,
			IFNULL(CAST(a.realisasi_kuantitatif AS CHAR), '')            AS realisasi_kuantitatif,
			IFNULL(a.realisasi_qualifier, '')                            AS realisasi_qualifier,
			IFNULL(a.realisasi_kuantitatif_qualifier, '')                AS realisasi_kuantitatif_qualifier,
			IFNULL(b.lampiran_file, '')                                  AS link_dokumen_sumber,
			IFNULL(rd.nama_result, '')                                   AS nama_result,
			IFNULL(rd.deskripsi_result, '')                              AS deskripsi_result,
			IFNULL(md.nama_method, '')                                   AS nama_method,
			IFNULL(md.deskripsi_method, '')                              AS deskripsi_method,
			IFNULL(cd.nama_challenge, '')                                AS nama_challenge,
			IFNULL(cd.deskripsi_challenge, '')                           AS deskripsi_challenge,
			IFNULL(rd.realisasi_result, '')                              AS realisasi_result,
			IFNULL(rd.lampiran_evidence, '')                             AS link_result,
			IFNULL(md.realisasi_method, '')                              AS realisasi_process,
			IFNULL(md.lampiran_evidence, '')                             AS link_process,
			IFNULL(cd.realisasi_challenge, '')                           AS realisasi_context,
			IFNULL(cd.lampiran_evidence, '')                             AS link_context
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

	queryGetRows = `
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

func (r *templateRepo) GetExistPenyusunanStatus(tahun, triwulan, kostl string) (status int, found bool, err error) {
	row := r.db.Raw(queryGetExistPenyusunanStatus, tahun, triwulan, kostl).Row()
	if scanErr := row.Scan(&status); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("gagal mengecek data penyusunan KPI: %w", scanErr)
	}
	return status, true, nil
}

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

func (r *templateRepo) CheckDataExist(idPengajuan, kostl, tahun, triwulan string) (bool, error) {
	var count int
	if err := r.db.Raw(queryExistData, idPengajuan, kostl, tahun, triwulan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal memvalidasi data tolakan KPI: %w", err)
	}
	return count > 0, nil
}

// GetPenyusunanKpiData mengambil header dan seluruh baris sub KPI
// dari DB berdasarkan id_pengajuan untuk keperluan generate Excel tolakan.
func (r *templateRepo) GetPenyusunanKpiData(idPengajuan string) (*model.ExcelData, error) {

	rows, err := r.db.Raw(queryGetRows, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data sub KPI tolakan: %w", err)
	}
	defer rows.Close()

	var subDetailRows []model.SubDetailRow
	for rows.Next() {
		var row model.SubDetailRow
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

	if len(subDetailRows) == 0 {
		return nil, nil
	}

	return &model.ExcelData{
		Rows: subDetailRows,
	}, nil
}

func (r *templateRepo) CheckRevisiRealisasiExist(idPengajuan, kostl, tahun, triwulan string) (bool, error) {
	var count int
	if err := r.db.Raw(queryCheckRevisiRealisasi, idPengajuan, kostl, tahun, triwulan).Scan(&count).Error; err != nil {
		return false, fmt.Errorf("gagal memvalidasi data revisi realisasi KPI: %w", err)
	}
	return count > 0, nil
}

// GetRealisasiKpiData mengambil seluruh baris realisasi KPI dari DB
// berdasarkan id_pengajuan, termasuk data realisasi yang sudah diisi user sebelumnya.
func (r *templateRepo) GetRealisasiKpiData(idPengajuan string) (*model.RealisasiExcelData, error) {
	rows, err := r.db.Raw(queryGetRealisasiRows, idPengajuan).Rows()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data realisasi KPI: %w", err)
	}
	defer rows.Close()

	var subDetailRows []model.RealisasiSubDetailRow
	for rows.Next() {
		var row model.RealisasiSubDetailRow
		if err := rows.Scan(
			&row.KpiNama,
			&row.SubKpi,
			&row.Polarisasi,
			&row.Capping,
			&row.Bobot,
			&row.TargetTriwulan,
			&row.ItemQualifier,
			&row.TargetQualifier,
			&row.TerdapatQualifier,
			&row.Realisasi,
			&row.RealisasiKuantitatif,
			&row.RealisasiQualifier,
			&row.RealisasiKuantitatifQualifier,
			&row.LinkDokumenSumber,
			&row.NamaResult,
			&row.DeskripsiResult,
			&row.NamaProcess,
			&row.DeskripsiProcess,
			&row.NamaContext,
			&row.DeskripsiContext,
			&row.RealisasiResult,
			&row.LinkResult,
			&row.RealisasiProcess,
			&row.LinkProcess,
			&row.RealisasiContext,
			&row.LinkContext,
		); err != nil {
			return nil, fmt.Errorf("gagal scan baris realisasi KPI: %w", err)
		}
		subDetailRows = append(subDetailRows, row)
	}

	if len(subDetailRows) == 0 {
		return nil, nil
	}

	return &model.RealisasiExcelData{
		Rows: subDetailRows,
	}, nil
}
