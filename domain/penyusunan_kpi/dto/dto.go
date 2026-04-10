package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

type DivisiItem struct {
	Kostl   string `json:"Kostl"   validate:"required"`
	KostlTx string `json:"KostlTx" validate:"required"`
}

// ValidatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiRequest struct {
	Divisi    DivisiItem `json:"Divisi"    validate:"required"`
	Tahun     string     `json:"Tahun"     validate:"required"`
	Triwulan  string     `json:"Triwulan"  validate:"required"`
	Kostl     string     `json:"Kostl"`
	KostlTx   string     `json:"KostlTx"`
	EntryUser string     `json:"EntryUser"`
	EntryName string     `json:"EntryName"`
	EntryTime string     `json:"EntryTime"`
}

// RevisionPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiRequest struct {
	IdPengajuan string     `json:"IdPengajuan" validate:"required"`
	Divisi      DivisiItem `json:"Divisi"      validate:"required"`
	Tahun       string     `json:"Tahun"       validate:"required"`
	Triwulan    string     `json:"Triwulan"    validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUser string `json:"EntryUser"`
	EntryName string `json:"EntryName"`
	EntryTime string `json:"EntryTime"`
}

// CreatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiRequest struct {
	IdPengajuan  string     `json:"idPengajuan"  validate:"required"`
	ApprovalList []Approval `json:"ApprovalList" validate:"required,min=1,dive"`
}

type Approval struct {
	Userid     string `json:"userid"`
	Nama       string `json:"nama"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
	Posisi     string `json:"posisi"`
	Level      string `json:"level"`
	Fungsi     string `json:"fungsi"`
	Waktu      string `json:"waktu"`
}

type PenyusunanResult struct {
	IdDetailResult  string `json:"idDetailResult"`
	Tahun           string `json:"tahun"`
	Triwulan        string `json:"triwulan"`
	NamaResult      string `json:"namaResult"`
	DeskripsiResult string `json:"deskripsiResult"`
}

type PenyusunanMethod struct {
	IdDetailMethod  string `json:"idDetailMethod"`
	Tahun           string `json:"tahun"`
	Triwulan        string `json:"triwulan"`
	NamaMethod      string `json:"namaMethod"`
	DeskripsiMethod string `json:"deskripsiMethod"`
}

type PenyusunanChallenge struct {
	IdDetailChallenge  string `json:"idDetailChallenge"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	NamaChallenge      string `json:"namaChallenge"`
	DeskripsiChallenge string `json:"deskripsiChallenge"`
}

// GetAllApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-approval.
type GetAllApprovalPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUser string `json:"-"`
}

// GetAllTolakanPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
type GetAllTolakanPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
type GetAllDaftarPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUser string `json:"-"`
}

// GetDetailPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-detail.
type GetDetailPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// GetExcelPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-excel.
type GetExcelPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// GetPdfPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-pdf.
type GetPdfPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// =============================================================================
// EXCEL ROW DTO
// =============================================================================

// PenyusunanKpiRow merepresentasikan 1 KPI unik yang ditemukan dari kolom B Excel,
// beserta hasil lookup ke tabel mst_kpi.
type PenyusunanKpiRow struct {
	// KpiIndex adalah urutan KPI unik (0-based) dari kolom B Excel.
	KpiIndex int
	// IdKpi adalah id_kpi dari mst_kpi. Jika tidak ditemukan, bernilai "0".
	IdKpi string
	// Kpi adalah nama KPI dari kolom B Excel (persis seperti yang diinput user).
	Kpi string
	// Rumus dari mst_kpi. Jika tidak ditemukan, bernilai "0".
	Rumus string
}

// PenyusunanKpiSubDetailRow merepresentasikan 1 baris data dari file Excel
// yang sudah diparse dan divalidasi.
type PenyusunanKpiSubDetailRow struct {
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

// =============================================================================
// RESPONSE DTO
// =============================================================================

// DivisiResponse digunakan di dalam ValidatePenyusunanKpiResponse.
type DivisiResponse struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostlTx"`
}

// EntryResponse digunakan di dalam ValidatePenyusunanKpiResponse.
type EntryResponse struct {
	EntryUser string `json:"entryUser"`
	EntryName string `json:"entryName"`
	EntryTime string `json:"entryTime"`
}

// ValidatePenyusunanKpiResponse adalah response untuk endpoint /validate.
type ValidatePenyusunanKpiResponse struct {
	IDPengajuan   string                        `json:"idPengajuan"`
	Tahun         string                        `json:"tahun"`
	Triwulan      string                        `json:"triwulan"`
	Divisi        DivisiResponse                `json:"divisi"`
	Entry         EntryResponse                 `json:"entry"`
	TotalKpi      int                           `json:"totalKpi"`
	Kpi           []PenyusunanKpiDetailResponse `json:"kpi"`
	ResultList    []PenyusunanResult            `json:"resultList"`
	MethodList    []PenyusunanMethod            `json:"methodList"`
	ChallengeList []PenyusunanChallenge         `json:"challengeList"`
}

// RevisionPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiResponse struct {
	IDPengajuan   string                        `json:"idPengajuan"`
	Tahun         string                        `json:"tahun"`
	Triwulan      string                        `json:"triwulan"`
	Divisi        DivisiResponse                `json:"divisi"`
	Entry         EntryResponse                 `json:"entry"`
	TotalKpi      int                           `json:"totalKpi"`
	Kpi           []PenyusunanKpiDetailResponse `json:"kpi"`
	ResultList    []PenyusunanResult            `json:"resultList"`
	MethodList    []PenyusunanMethod            `json:"methodList"`
	ChallengeList []PenyusunanChallenge         `json:"challengeList"`
}

// CreatePenyusunanKpiResponse adalah response untuk endpoint /create.
type CreatePenyusunanKpiResponse struct {
	IdPengajuan  string     `json:"idPengajuan"`
	ApprovalList []Approval `json:"approvalList"`
}

type PenyusunanKpiDetailResponse struct {
	IdDetail     string                           `json:"idDetail"`
	IdKpi        string                           `json:"idKpi"`
	Kpi          string                           `json:"kpi"`
	Rumus        string                           `json:"rumus"`
	Persfektif   string                           `json:"persfektif"`
	TotalSubKpi  int                              `json:"totalSubKpi"`
	KpiSubDetail []PenyusunanKpiSubDetailResponse `json:"kpiSubDetail"`
}

type PenyusunanKpiSubDetailResponse struct {
	IdSubDetail               string  `json:"idSubDetail"`
	IdSubKpi                  string  `json:"idSubKpi"`
	SubKpi                    string  `json:"subKpi"`
	Otomatis                  string  `json:"otomatis"`
	Polarisasi                string  `json:"polarisasi"`
	IdPolarisasi              string  `json:"idPolarisasi"`
	Capping                   string  `json:"capping"`
	Bobot                     float64 `json:"bobot"`
	Glossary                  string  `json:"glossary"`
	TargetTriwulan            string  `json:"targetTriwulan"`
	TargetKuantitatifTriwulan float64 `json:"targetKuantitatifTriwulan"`
	TargetTahunan             string  `json:"targetTahunan"`
	TargetKuantitatifTahunan  float64 `json:"targetKuantitatifTahunan"`
	TerdapatQualifier         string  `json:"terdapatQualifier"`
	Qualifier                 string  `json:"qualifier"`
	DeskripsiQualifier        string  `json:"deskripsiQualifier"`
	TargetQualifier           string  `json:"targetQualifier"`
}

// GetAllApprovalPenyusunanKpiResponse adalah satu record lengkap (header + nested detail).
type GetAllApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"IdPengajuan"`
	Tahun       string `json:"Tahun"`
	Triwulan    string `json:"Triwulan"`
	Kostl       string `json:"Kostl"`
	KostlTx     string `json:"KostlTx"`
	Orgeh       string `json:"Orgeh"`
	OrgehTx     string `json:"OrgehTx"`
	Status      string `json:"Status"`
	StatusDesc  string `json:"StatusDesc"`
}

// GetAllTolakanPenyusunanKpiResponse adalah satu record lengkap (header + nested detail).
type GetAllTolakanPenyusunanKpiResponse struct {
	IdPengajuan string `json:"IdPengajuan"`
	Tahun       string `json:"Tahun"`
	Triwulan    string `json:"Triwulan"`
	Kostl       string `json:"Kostl"`
	KostlTx     string `json:"KostlTx"`
	Orgeh       string `json:"Orgeh"`
	OrgehTx     string `json:"OrgehTx"`
	Status      string `json:"Status"`
	StatusDesc  string `json:"StatusDesc"`
}

// GetAllDaftarPenyusunanKpiResponse adalah satu record lengkap (header + nested detail).
type GetAllDaftarPenyusunanKpiResponse struct {
	IdPengajuan string `json:"IdPengajuan"`
	Tahun       string `json:"Tahun"`
	Triwulan    string `json:"Triwulan"`
	Kostl       string `json:"Kostl"`
	KostlTx     string `json:"KostlTx"`
	Orgeh       string `json:"Orgeh"`
	OrgehTx     string `json:"OrgehTx"`
	Status      string `json:"Status"`
	StatusDesc  string `json:"StatusDesc"`
}

// GetAllDaftarApprovalPenyusunanKpiResponse adalah satu record lengkap (header + nested detail).
type GetAllDaftarApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"IdPengajuan"`
	Tahun       string `json:"Tahun"`
	Triwulan    string `json:"Triwulan"`
	Kostl       string `json:"Kostl"`
	KostlTx     string `json:"KostlTx"`
	Orgeh       string `json:"Orgeh"`
	OrgehTx     string `json:"OrgehTx"`
	Status      string `json:"Status"`
	StatusDesc  string `json:"StatusDesc"`
}

// GetAllDataPenyusunanKpiResponse adalah satu record lengkap (header + nested detail).
type GetAllDataPenyusunanKpiResponse struct {
	IdPengajuan              string                              `json:"IdPengajuan"`
	Tahun                    string                              `json:"Tahun"`
	Triwulan                 string                              `json:"Triwulan"`
	Kostl                    string                              `json:"Kostl"`
	KostlTx                  string                              `json:"KostlTx"`
	Orgeh                    string                              `json:"Orgeh"`
	OrgehTx                  string                              `json:"OrgehTx"`
	EntryUser                string                              `json:"EntryUser"`
	EntryName                string                              `json:"EntryName"`
	EntryTime                string                              `json:"EntryTime"`
	ApprovalPosisi           string                              `json:"ApprovalPosisi"`
	ApprovalList             string                              `json:"ApprovalList"`
	Status                   string                              `json:"Status"`
	StatusDesc               string                              `json:"StatusDesc"`
	EntryUserRealisasi       string                              `json:"EntryUserRealisasi"`
	EntryNameRealisasi       string                              `json:"EntryNameRealisasi"`
	EntryTimeRealisasi       string                              `json:"EntryTimeRealisasi"`
	ApprovalListRealisasi    string                              `json:"ApprovalListRealisasi"`
	CatatanTolakan           string                              `json:"CatatanTolakan"`
	TotalBobot               string                              `json:"TotalBobot"`
	TotalPencapaian          string                              `json:"TotalPencapaian"`
	TotalBobotPengurang      string                              `json:"TotalBobotPengurang"`
	TotalPencapaianPost      string                              `json:"TotalPencapaianPost"`
	EntryUserValidasi        string                              `json:"EntryUserValidasi"`
	EntryNameValidasi        string                              `json:"EntryNameValidasi"`
	EntryTimeValidasi        string                              `json:"EntryTimeValidasi"`
	ApprovalListValidasi     string                              `json:"ApprovalListValidasi"`
	LampiranValidasi         string                              `json:"LampiranValidasi"`
	QualifierOverallValidasi string                              `json:"QualifierOverallValidasi"`
	KpiDetail                []GetAllDataKpiDetailResponse       `json:"KpiDetail"`
	ResultDetail             []GetAllDataResultDetailResponse    `json:"ResultDetail"`
	MethodDetail             []GetAllDataMethodDetailResponse    `json:"MethodDetail"`
	ChallengeDetail          []GetAllDataChallengeDetailResponse `json:"ChallengeDetail"`
}

// GetAllDataKpiDetailResponse — IdPengajuan, Tahun, Triwulan dihapus (redundan)
type GetAllDataKpiDetailResponse struct {
	IdDetail            string                           `json:"IdDetail"`
	IdKpi               string                           `json:"IdKpi"`
	Kpi                 string                           `json:"Kpi"`
	Rumus               string                           `json:"Rumus"`
	IdPerspektif        string                           `json:"IdPerspektif"`
	Perspektif          string                           `json:"Perspektif"`
	IdKeteranganProject string                           `json:"IdKeteranganProject"`
	KeteranganProject   string                           `json:"KeteranganProject"`
	LampiranFile        string                           `json:"LampiranFile"`
	KpiSubDetail        []GetAllDataKpiSubDetailResponse `json:"KpiSubDetail"`
}

// GetAllDataKpiSubDetailResponse — IdPengajuan, IdDetail, Tahun, Triwulan dihapus (redundan)
type GetAllDataKpiSubDetailResponse struct {
	IdSubDetail                      string `json:"IdSubDetail"`
	IdKpi                            string `json:"IdKpi"`
	Kpi                              string `json:"Kpi"`
	Rumus                            string `json:"Rumus"`
	Otomatis                         string `json:"Otomatis"`
	Bobot                            string `json:"Bobot"`
	Capping                          string `json:"Capping"`
	TargetTriwulan                   string `json:"TargetTriwulan"`
	TargetKuantitatifTriwulan        string `json:"TargetKuantitatifTriwulan"`
	TargetTahunan                    string `json:"TargetTahunan"`
	TargetKuantitatifTahunan         string `json:"TargetKuantitatifTahunan"`
	Realisasi                        string `json:"Realisasi"`
	RealisasiKuantitatif             string `json:"RealisasiKuantitatif"`
	RealisasiKeterangan              string `json:"RealisasiKeterangan"`
	RealisasiValidated               string `json:"RealisasiValidated"`
	RealisasiKuantitatifValidated    string `json:"RealisasiKuantitatifValidated"`
	ValidasiKeterangan               string `json:"ValidasiKeterangan"`
	Pencapaian                       string `json:"Pencapaian"`
	Skor                             string `json:"Skor"`
	DeskripsiGlossary                string `json:"DeskripsiGlossary"`
	ItemQualifier                    string `json:"ItemQualifier"`
	DeskripsiQualifier               string `json:"DeskripsiQualifier"`
	TargetQualifier                  string `json:"TargetQualifier"`
	IdKeteranganProject              string `json:"IdKeteranganProject"`
	KeteranganProject                string `json:"KeteranganProject"`
	IdQualifier                      string `json:"IdQualifier"`
	RealisasiQualifier               string `json:"RealisasiQualifier"`
	RealisasiKuantitatifQualifier    string `json:"RealisasiKuantitatifQualifier"`
	PencapaianQualifierValidated     string `json:"PencapaianQualifierValidated"`
	PencapaianPostQualifierValidated string `json:"PencapaianPostQualifierValidated"`
}

// GetAllDataResultDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataResultDetailResponse struct {
	IdDetailResult  string `json:"IdDetailResult"`
	Tahun           string `json:"Tahun"`
	Triwulan        string `json:"Triwulan"`
	NamaResult      string `json:"NamaResult"`
	DeskripsiResult string `json:"DeskripsiResult"`
}

// GetAllDataMethodDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataMethodDetailResponse struct {
	IdDetailMethod   string `json:"IdDetailMethod"`
	Tahun            string `json:"Tahun"`
	Triwulan         string `json:"Triwulan"`
	NamaMethod       string `json:"NamaMethod"`
	DeskripsiMethod  string `json:"DeskripsiMethod"`
	RealisasiMethod  string `json:"RealisasiMethod"`
	LampiranEvidence string `json:"LampiranEvidence"`
}

// GetAllDataChallengeDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataChallengeDetailResponse struct {
	IdDetailChallenge  string `json:"IdDetailChallenge"`
	Tahun              string `json:"Tahun"`
	Triwulan           string `json:"Triwulan"`
	NamaChallenge      string `json:"NamaChallenge"`
	DeskripsiChallenge string `json:"DeskripsiChallenge"`
	RealisasiChallenge string `json:"RealisasiChallenge"`
	LampiranEvidence   string `json:"LampiranEvidence"`
}

// =============================================================================
// EXPORT DTO (digunakan oleh get-csv dan get-pdf)
// =============================================================================

// KpiSubDetailExportRow merepresentasikan 1 baris data sub KPI untuk keperluan
// ekspor CSV dan PDF. Kolom sesuai tampilan: No, KPI, Bobot(%), Target Tahunan, Capping.
type KpiSubDetailExportRow struct {
	No            int
	KpiNama       string
	Bobot         string
	TargetTahunan string
	Capping       string
}

// KpiExportData berisi header dokumen + daftar baris sub KPI untuk ekspor.
type KpiExportData struct {
	NamaDivisi string // kostl_tx dari data_kpi
	Tahun      string
	Triwulan   string
	Rows       []KpiSubDetailExportRow
}
