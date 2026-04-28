package model

// SubDetailRow merepresentasikan 1 baris data sub KPI dari DB
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
type SubDetailRow struct {
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

// ExcelData berisi header dokumen + daftar baris sub KPI
// untuk keperluan generate Excel tolakan penyusunan KPI.
type ExcelData struct {
	Rows []SubDetailRow
}

// RealisasiSubDetailRow merepresentasikan 1 baris data realisasi KPI dari DB
// untuk keperluan generate Excel revisi realisasi KPI.
//
// Kolom Excel format-realisasi-kpi (TW1/TW3):
//
//	A  = No (sequential)
//	B  = KpiNama        → data_kpi_detail.kpi
//	C  = SubKpi         → data_kpi_subdetail.kpi
//	D  = Polarisasi     → mst_polarisasi.polarisasi
//	E  = Capping        → data_kpi_subdetail.capping
//	F  = Bobot          → data_kpi_subdetail.bobot
//	G  = TargetTriwulan → data_kpi_subdetail.target_triwulan
//	H  = ItemQualifier  → data_kpi_subdetail.item_qualifier
//	I  = TargetQualifier → data_kpi_subdetail.target_qualifier
//	J  = Realisasi      → data_kpi_subdetail.realisasi          (pre-filled)
//	K  = RealisasiKuantitatif → data_kpi_subdetail.realisasi_kuantitatif (pre-filled)
//	L  = RealisasiQualifier   → data_kpi_subdetail.realisasi_qualifier   (pre-filled, jika ada qualifier)
//	M  = RealisasiKuantitatifQualifier → data_kpi_subdetail.realisasi_kuantitatif_qualifier (pre-filled)
//
// Kolom TW2/TW4 extended:
//
//	N  = NamaResult         → data_result_detail.nama_result
//	O  = DeskripsiResult    → data_result_detail.deskripsi_result
//	P  = RealisasiResult    → data_result_detail.realisasi_result (pre-filled)
//	Q  = LinkResult         → data_result_detail.lampiran_evidence (pre-filled)
//	R  = NamaProcess        → data_method_detail.nama_method
//	S  = DeskripsiProcess   → data_method_detail.deskripsi_method
//	T  = RealisasiProcess   → data_method_detail.realisasi_method (pre-filled)
//	U  = LinkProcess        → data_method_detail.lampiran_evidence (pre-filled)
//	V  = NamaContext        → data_challenge_detail.nama_challenge
//	W  = DeskripsiContext   → data_challenge_detail.deskripsi_challenge
//	X  = RealisasiContext   → data_challenge_detail.realisasi_challenge (pre-filled)
//	Y  = LinkContext        → data_challenge_detail.lampiran_evidence (pre-filled)
type RealisasiSubDetailRow struct {
	KpiNama       string
	SubKpi        string
	Polarisasi    string
	Capping       string
	Bobot         string
	TargetTriwulan  string
	ItemQualifier   string
	TargetQualifier string
	TerdapatQualifier string
	// Kolom realisasi (J–M) — pre-filled dari DB
	Realisasi                    string
	RealisasiKuantitatif         string
	RealisasiQualifier           string
	RealisasiKuantitatifQualifier string
	// TW2/TW4 penyusunan (N,O,R,S,V,W)
	NamaResult       string
	DeskripsiResult  string
	NamaProcess      string
	DeskripsiProcess string
	NamaContext      string
	DeskripsiContext string
	// TW2/TW4 realisasi extended (P,Q,T,U,X,Y) — pre-filled dari DB
	RealisasiResult  string
	LinkResult       string
	RealisasiProcess string
	LinkProcess      string
	RealisasiContext string
	LinkContext      string
}

// RealisasiExcelData berisi daftar baris realisasi KPI
// untuk keperluan generate Excel revisi realisasi KPI.
type RealisasiExcelData struct {
	Rows []RealisasiSubDetailRow
}
