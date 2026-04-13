package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

type DivisiItem struct {
	Kostl   string `json:"kostl"   validate:"required"`
	KostlTx string `json:"kostl_tx" validate:"required"`
}

// ValidatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiRequest struct {
	Divisi    DivisiItem `json:"divisi"    validate:"required"`
	Tahun     string     `json:"tahun"     validate:"required"`
	Triwulan  string     `json:"triwulan"  validate:"required"`
	Kostl     string     `json:"kostl"`
	KostlTx   string     `json:"kostl_tx"`
	EntryUser string     `json:"entry_user"`
	EntryName string     `json:"entry_name"`
	EntryTime string     `json:"entry_time"`
}

// RevisionPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiRequest struct {
	IdPengajuan string     `json:"id_pengajuan" validate:"required"`
	Divisi      DivisiItem `json:"divisi"      validate:"required"`
	Tahun       string     `json:"tahun"       validate:"required"`
	Triwulan    string     `json:"triwulan"    validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// CreatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiRequest struct {
	IdPengajuan  string     `json:"id_pengajuan"  validate:"required"`
	ApprovalList []Approval `json:"approval_list" validate:"required,min=1,dive"`
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
	IdDetailResult  string `json:"id_detail_result"`
	Tahun           string `json:"tahun"`
	Triwulan        string `json:"triwulan"`
	NamaResult      string `json:"nama_result"`
	DeskripsiResult string `json:"deskripsi_result"`
}

type PenyusunanProcess struct {
	IdDetailProcess  string `json:"id_detail_process"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaProcess      string `json:"nama_process"`
	DeskripsiProcess string `json:"deskripsi_process"`
}

type PenyusunanContext struct {
	IdDetailContext  string `json:"id_detail_context"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaContext      string `json:"nama_context"`
	DeskripsiContext string `json:"deskripsi_context"`
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

// DivisiResponse digunakan di dalam ValidatePenyusunanKpiResponse dan RevisionPenyusunanKpiResponse.
type DivisiResponse struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostl_tx"`
}

// DivisiDetailResponse digunakan di dalam GetDetailPenyusunanKpiResponse.
// Berbeda dari DivisiResponse karena menyertakan Orgeh/OrgehTx dan menggunakan PascalCase json tag.
type DivisiDetailResponse struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostl_tx"`
	Orgeh   string `json:"orgeh"`
	OrgehTx string `json:"orgeh_tx"`
}

// EntryResponse digunakan di dalam ValidatePenyusunanKpiResponse.
type EntryResponse struct {
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// ValidatePenyusunanKpiResponse adalah response untuk endpoint /validate.
type ValidatePenyusunanKpiResponse struct {
	IDPengajuan string                        `json:"id_pengajuan"`
	Tahun       string                        `json:"tahun"`
	Triwulan    string                        `json:"triwulan"`
	Divisi      DivisiResponse                `json:"divisi"`
	Entry       EntryResponse                 `json:"entry"`
	TotalKpi    int                           `json:"total_kpi"`
	Kpi         []PenyusunanKpiDetailResponse `json:"kpi"`
	ResultList  []PenyusunanResult            `json:"result_list"`
	ProcessList []PenyusunanProcess           `json:"process_list"`
	ContextList []PenyusunanContext           `json:"context_list"`
}

// RevisionPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiResponse struct {
	IDPengajuan string                        `json:"id_pengajuan"`
	Tahun       string                        `json:"tahun"`
	Triwulan    string                        `json:"triwulan"`
	Divisi      DivisiResponse                `json:"divisi"`
	Entry       EntryResponse                 `json:"entry"`
	TotalKpi    int                           `json:"total_kpi"`
	Kpi         []PenyusunanKpiDetailResponse `json:"kpi"`
	ResultList  []PenyusunanResult            `json:"result_list"`
	ProcessList []PenyusunanProcess           `json:"process_list"`
	ContextList []PenyusunanContext           `json:"context_list"`
}

// CreatePenyusunanKpiResponse adalah response untuk endpoint /create.
type CreatePenyusunanKpiResponse struct {
	IdPengajuan  string     `json:"id_pengajuan"`
	ApprovalList []Approval `json:"approval_list"`
}

// PenyusunanKpiDetailResponse merepresentasikan 1 KPI beserta sub-detail-nya.
// Digunakan oleh validate, revision, dan get-detail.
type PenyusunanKpiDetailResponse struct {
	IdDetail     string                           `json:"id_detail"`
	IdKpi        string                           `json:"id_kpi"`
	Kpi          string                           `json:"kpi"`
	Rumus        string                           `json:"rumus"`
	IdPerspektif string                           `json:"id_perspektif"`
	Persfektif   string                           `json:"persfektif"`
	TotalSubKpi  int                              `json:"total_sub_kpi"`
	KpiSubDetail []PenyusunanKpiSubDetailResponse `json:"kpi_sub_detail"`
}

// PenyusunanKpiSubDetailResponse merepresentasikan 1 sub KPI.
// Digunakan oleh validate, revision, dan get-detail.
type PenyusunanKpiSubDetailResponse struct {
	IdSubDetail               string  `json:"id_sub_detail"`
	IdSubKpi                  string  `json:"id_sub_kpi"`
	SubKpi                    string  `json:"sub_kpi"`
	Otomatis                  string  `json:"otomatis"`
	Polarisasi                string  `json:"polarisasi"`
	IdPolarisasi              string  `json:"id_polarisasi"`
	Capping                   string  `json:"capping"`
	Bobot                     float64 `json:"bobot"`
	Glossary                  string  `json:"glossary"`
	TargetTriwulan            string  `json:"target_triwulan"`
	TargetKuantitatifTriwulan float64 `json:"target_kuantitatif_triwulan"`
	TargetTahunan             string  `json:"target_tahunan"`
	TargetKuantitatifTahunan  float64 `json:"target_kuantitatif_tahunan"`
	TerdapatQualifier         string  `json:"terdapat_qualifier"`
	Qualifier                 string  `json:"qualifier"`
	DeskripsiQualifier        string  `json:"deskripsi_qualifier"`
	TargetQualifier           string  `json:"target_qualifier"`
	IdKeteranganProject       string  `json:"id_keterangan_project"`
	KeteranganProject         string  `json:"keterangan_project"`
}

// GetDetailPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-detail.
type GetDetailPenyusunanKpiResponse struct {
	IdPengajuan    string                        `json:"id_pengajuan"`
	Tahun          string                        `json:"tahun"`
	Triwulan       string                        `json:"triwulan"`
	Status         string                        `json:"status"`
	StatusDesc     string                        `json:"status_desc"`
	Divisi         DivisiDetailResponse          `json:"divisi"`
	Entry          EntryResponse                 `json:"entry"`
	ApprovalPosisi string                        `json:"approval_posisi"`
	ApprovalList   []Approval                    `json:"approval_list"`
	TotalKpi       int                           `json:"total_kpi"`
	Kpi            []PenyusunanKpiDetailResponse `json:"kpi"`
	TotalResult    int                           `json:"total_result"`
	ResultList     []PenyusunanResult            `json:"result_list"`
	TotalProcess   int                           `json:"total_process"`
	ProcessList    []PenyusunanProcess           `json:"process_list"`
	TotalContext   int                           `json:"total_context"`
	ContextList    []PenyusunanContext           `json:"context_list"`
}

// GetAllApprovalPenyusunanKpiResponse adalah satu record untuk list approval.
type GetAllApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	Kostl       string `json:"kostl"`
	KostlTx     string `json:"kostl_tx"`
	Orgeh       string `json:"orgeh"`
	OrgehTx     string `json:"orgeh_tx"`
	Status      string `json:"status"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllTolakanPenyusunanKpiResponse adalah satu record untuk list tolakan.
type GetAllTolakanPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	Kostl       string `json:"kostl"`
	KostlTx     string `json:"kostl_tx"`
	Orgeh       string `json:"orgeh"`
	OrgehTx     string `json:"orgeh_tx"`
	Status      string `json:"status"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarPenyusunanKpiResponse adalah satu record untuk daftar penyusunan.
type GetAllDaftarPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	Kostl       string `json:"kostl"`
	KostlTx     string `json:"kostl_tx"`
	Orgeh       string `json:"orgeh"`
	OrgehTx     string `json:"orgeh_tx"`
	Status      string `json:"status"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarApprovalPenyusunanKpiResponse adalah satu record untuk daftar approval.
type GetAllDaftarApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	Kostl       string `json:"kostl"`
	KostlTx     string `json:"kostl_tx"`
	Orgeh       string `json:"orgeh"`
	OrgehTx     string `json:"orgeh_tx"`
	Status      string `json:"status"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDataPenyusunanKpiResponse adalah response lengkap (header + nested detail) untuk konteks validasi.
type GetAllDataPenyusunanKpiResponse struct {
	IdPengajuan              string                            `json:"id_pengajuan"`
	Tahun                    string                            `json:"tahun"`
	Triwulan                 string                            `json:"triwulan"`
	Kostl                    string                            `json:"kostl"`
	KostlTx                  string                            `json:"kostl_tx"`
	Orgeh                    string                            `json:"orgeh"`
	OrgehTx                  string                            `json:"orgeh_tx"`
	EntryUser                string                            `json:"entry_user"`
	EntryName                string                            `json:"entry_name"`
	EntryTime                string                            `json:"entry_time"`
	ApprovalPosisi           string                            `json:"approval_posisi"`
	ApprovalList             string                            `json:"approval_list"`
	Status                   string                            `json:"status"`
	StatusDesc               string                            `json:"status_desc"`
	EntryUserRealisasi       string                            `json:"entry_user_realisasi"`
	EntryNameRealisasi       string                            `json:"entry_name_realisasi"`
	EntryTimeRealisasi       string                            `json:"entry_time_realisasi"`
	ApprovalListRealisasi    string                            `json:"approval_list_realisasi"`
	CatatanTolakan           string                            `json:"catatan_tolakan"`
	TotalBobot               string                            `json:"total_bobot"`
	TotalPencapaian          string                            `json:"total_pencapaian"`
	TotalBobotPengurang      string                            `json:"total_bobot_pengurang"`
	TotalPencapaianPost      string                            `json:"total_pencapaian_post"`
	EntryUserValidasi        string                            `json:"entry_user_validasi"`
	EntryNameValidasi        string                            `json:"entry_name_validasi"`
	EntryTimeValidasi        string                            `json:"entry_time_validasi"`
	ApprovalListValidasi     string                            `json:"approval_list_validasi"`
	LampiranValidasi         string                            `json:"lampiran_validasi"`
	QualifierOverallValidasi string                            `json:"qualifier_overall_validasi"`
	KpiDetail                []GetAllDataKpiDetailResponse     `json:"kpi_detail"`
	ResultDetail             []GetAllDataResultDetailResponse  `json:"result_detail"`
	ProcessDetail            []GetAllDataProcessDetailResponse `json:"process_detail"`
	ContextDetail            []GetAllDataContextDetailResponse `json:"context_detail"`
}

// GetAllDataKpiDetailResponse — IdPengajuan, Tahun, Triwulan dihapus (redundan)
type GetAllDataKpiDetailResponse struct {
	IdDetail            string                           `json:"id_detail"`
	IdKpi               string                           `json:"id_kpi"`
	Kpi                 string                           `json:"kpi"`
	Rumus               string                           `json:"rumus"`
	IdPerspektif        string                           `json:"id_perspektif"`
	Perspektif          string                           `json:"persfektif"`
	IdKeteranganProject string                           `json:"id_keterangan_project"`
	KeteranganProject   string                           `json:"keterangan_project"`
	LampiranFile        string                           `json:"lampiran_file"`
	KpiSubDetail        []GetAllDataKpiSubDetailResponse `json:"kpi_sub_detail"`
}

// GetAllDataKpiSubDetailResponse — IdPengajuan, IdDetail, Tahun, Triwulan dihapus (redundan)
type GetAllDataKpiSubDetailResponse struct {
	IdSubDetail                      string `json:"id_sub_detail"`
	IdKpi                            string `json:"id_kpi"`
	Kpi                              string `json:"kpi"`
	Rumus                            string `json:"rumus"`
	Otomatis                         string `json:"otomatis"`
	Bobot                            string `json:"bobot"`
	Capping                          string `json:"capping"`
	TargetTriwulan                   string `json:"target_triwulan"`
	TargetKuantitatifTriwulan        string `json:"target_kuantitatif_triwulan"`
	TargetTahunan                    string `json:"target_tahunan"`
	TargetKuantitatifTahunan         string `json:"target_kuantitatif_tahunan"`
	Realisasi                        string `json:"realisasi"`
	RealisasiKuantitatif             string `json:"realisasi_kuantitatif"`
	RealisasiKeterangan              string `json:"realisasi_keterangan"`
	RealisasiValidated               string `json:"realisasi_validated"`
	RealisasiKuantitatifValidated    string `json:"realisasi_kuantitatif_validated"`
	ValidasiKeterangan               string `json:"validasi_keterangan"`
	Pencapaian                       string `json:"pencapaian"`
	Skor                             string `json:"skor"`
	DeskripsiGlossary                string `json:"deskripsi_glossary"`
	ItemQualifier                    string `json:"item_qualifier"`
	DeskripsiQualifier               string `json:"deskripsi_qualifier"`
	TargetQualifier                  string `json:"target_qualifier"`
	IdKeteranganProject              string `json:"id_keterangan_project"`
	KeteranganProject                string `json:"keterangan_project"`
	IdQualifier                      string `json:"id_qualifier"`
	RealisasiQualifier               string `json:"realisasi_qualifier"`
	RealisasiKuantitatifQualifier    string `json:"realisasi_kuantitatif_qualifier"`
	PencapaianQualifierValidated     string `json:"pencapaian_qualifier_validated"`
	PencapaianPostQualifierValidated string `json:"pencapaian_post_qualifier_validated"`
}

// GetAllDataResultDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataResultDetailResponse struct {
	IdDetailResult  string `json:"Id_detail_result"`
	Tahun           string `json:"tahun"`
	Triwulan        string `json:"triwulan"`
	NamaResult      string `json:"nama_result"`
	DeskripsiResult string `json:"deskripsi_result"`
}

// GetAllDataProcessDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataProcessDetailResponse struct {
	IdDetailProcess  string `json:"id_detail_process"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaProcess      string `json:"nama_process"`
	DeskripsiProcess string `json:"deskripsi_process"`
	RealisasiProcess string `json:"realisasi_process"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

// GetAllDataContextDetailResponse — IdPengajuan dihapus (redundan)
type GetAllDataContextDetailResponse struct {
	IdDetailContext  string `json:"id_detail_context"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaContext      string `json:"nama_context"`
	DeskripsiContext string `json:"deskripsi_context"`
	RealisasiContext string `json:"realisasi_context"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

// =============================================================================
// EXPORT DTO (digunakan oleh get-excel dan get-pdf)
// =============================================================================

// KpiSubDetailExportRow merepresentasikan 1 baris data sub KPI untuk keperluan ekspor.
type KpiSubDetailExportRow struct {
	No            int
	KpiNama       string
	Bobot         string
	TargetTahunan string
	Capping       string
}

// KpiExportData berisi header dokumen + daftar baris sub KPI untuk ekspor.
type KpiExportData struct {
	NamaDivisi string
	Tahun      string
	Triwulan   string
	Rows       []KpiSubDetailExportRow
}
