package excel

// KpiRow merepresentasikan 1 KPI unik yang ditemukan dari kolom B Excel,
// beserta hasil lookup ke tabel mst_kpi.
type KpiRow struct {
	// KpiIndex adalah urutan KPI unik (0-based) dari kolom B Excel.
	KpiIndex int
	// IdKpi adalah id_kpi dari mst_kpi. Jika tidak ditemukan, bernilai "0".
	IdKpi string
	// Kpi adalah nama KPI dari kolom B Excel (persis seperti yang diinput user).
	Kpi string
	// Rumus dari mst_kpi. Jika tidak ditemukan, bernilai "0".
	Rumus string
}

// KpiSubDetailRow merepresentasikan 1 baris data dari file Excel
// yang sudah diparse dan divalidasi.
type KpiSubDetailRow struct {
	No                        int
	KPI                       string
	SubKPI                    string
	IdSubKpi                  string
	Otomatis                  string
	Polarisasi                string
	IdPolarisasi              string
	Capping                   string
	Bobot                     float64
	Glossary                  string
	TargetTriwulan            string
	TargetKuantitatifTriwulan float64
	TargetTahunan             string
	TargetKuantitatifTahunan  float64
	TerdapatQualifier         string
	Qualifier                 string
	DeskripsiQualifier        string
	TargetQualifier           string
	IsTW4                     bool
	Result                    *string
	DeskripsiResult           *string
	Process                   *string
	DeskripsiProcess          *string
	Context                   *string
	DeskripsiContext          *string
}
