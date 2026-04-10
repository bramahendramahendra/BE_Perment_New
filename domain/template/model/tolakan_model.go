package model

// TolakanSubDetailRow merepresentasikan 1 baris data sub KPI dari DB
// untuk keperluan generate Excel tolakan penyusunan KPI.
//
// Mapping kolom Excel:
//
//	A  = No (sequential)
//	B  = KpiNama       → data_kpi_detail.kpi
//	C  = SubKpi        → data_kpi_subdetail.kpi
//	D  = Polarisasi    → mst_polarisasi.polarisasi (join via data_kpi_subdetail.rumus)
//	E  = Capping       → data_kpi_subdetail.capping
//	F  = Bobot         → data_kpi_subdetail.bobot
//	G  = Glossary      → data_kpi_subdetail.deskripsi_glossary
//	H  = TargetTriwulan            → data_kpi_subdetail.target_triwulan
//	I  = TargetKuantitatifTriwulan → data_kpi_subdetail.target_kuantitatif_triwulan
//	J  = TargetTahunan             → data_kpi_subdetail.target_tahunan
//	K  = TargetKuantitatifTahunan  → data_kpi_subdetail.target_kuantitatif_tahunan
//	L  = TerdapatQualifier → data_kpi_subdetail.id_qualifier  ("Ya"/"Tidak")
//	M  = ItemQualifier     → data_kpi_subdetail.item_qualifier
//	N  = DeskripsiQualifier → data_kpi_subdetail.deskripsi_qualifier
//	O  = TargetQualifier   → data_kpi_subdetail.target_qualifier
//	P  = NamaResult        → data_result_detail.nama_result      (TW2/TW4)
//	Q  = DeskripsiResult   → data_result_detail.deskripsi_result (TW2/TW4)
//	R  = NamaProcess       → data_method_detail.nama_method      (TW2/TW4)
//	S  = DeskripsiProcess  → data_method_detail.deskripsi_method (TW2/TW4)
//	T  = NamaContext       → data_challenge_detail.nama_challenge      (TW2/TW4)
//	U  = DeskripsiContext  → data_challenge_detail.deskripsi_challenge (TW2/TW4)
type TolakanSubDetailRow struct {
	IdSubDetail               string
	KpiNama                   string
	SubKpi                    string
	Polarisasi                string
	Capping                   string
	Bobot                     string
	DeskripsiGlossary         string
	TargetTriwulan            string
	TargetKuantitatifTriwulan string
	TargetTahunan             string
	TargetKuantitatifTahunan  string
	TerdapatQualifier         string
	ItemQualifier             string
	DeskripsiQualifier        string
	TargetQualifier           string
	// Extended — hanya diisi untuk TW2 dan TW4
	NamaResult       string
	DeskripsiResult  string
	NamaProcess      string
	DeskripsiProcess string
	NamaContext      string
	DeskripsiContext string
}

// TolakanExcelData berisi header dokumen + daftar baris sub KPI
// untuk keperluan generate Excel tolakan penyusunan KPI.
type TolakanExcelData struct {
	Triwulan string
	Tahun    string
	KostlTx  string
	Rows     []TolakanSubDetailRow
}
