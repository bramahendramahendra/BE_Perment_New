package dto

// =============================================================================
// GENERAL DTO
// =============================================================================

type Divisi struct {
	Kostl   string `json:"kostl"   validate:"required"`
	KostlTx string `json:"kostl_tx" validate:"required"`
}

// Berbeda dari Divisi karena menyertakan Orgeh/OrgehTx.
type DivisiOrgeh struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostl_tx"`
	Orgeh   string `json:"orgeh"`
	OrgehTx string `json:"orgeh_tx"`
}

type EntryUserPenyusunan struct {
	EntryUserPenyusunan string `json:"entry_user_penyusunan"`
	EntryNamePenyusunan string `json:"entry_name_penyusunan"`
	EntryTimePenyusunan string `json:"entry_time_penyusunan"`
}

type EntryUserRealisasi struct {
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
}

type EntryUserValidasi struct {
	EntryUserValidasi string `json:"entry_user_validasi"`
	EntryNameValidasi string `json:"entry_name_validasi"`
	EntryTimeValidasi string `json:"entry_time_validasi"`
}

type ApprovalUserRealisasi struct {
	Userid string `json:"userid"`
	Nama   string `json:"nama"`
}

type ApprovalUserRealisasiDetail struct {
	Userid     string `json:"userid"`
	Nama       string `json:"nama"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
	Posisi     string `json:"posisi"`
	Level      string `json:"level"`
	Fungsi     string `json:"fungsi"`
	Waktu      string `json:"waktu"`
}

type CatatanTolakanEntry struct {
	Fungsi    string `json:"fungsi"`
	EntryUser string `json:"entry_user"`
	EntryTime string `json:"entry_time"`
	EntryNote string `json:"entry_note"`
}

type DataKpiDetail struct {
	IdDetail            string             `json:"id_detail"`
	IdKpi               string             `json:"id_kpi"`
	Kpi                 string             `json:"kpi"`
	Rumus               string             `json:"rumus"`
	IdPerspektif        string             `json:"id_perspektif"`
	Persfektif          string             `json:"persfektif"`
	IdKeteranganProject string             `json:"id_keterangan_project"`
	KeteranganProject   string             `json:"keterangan_project"`
	TotalSubKpi         int                `json:"total_sub_kpi"`
	KpiSubDetail        []DataKpiSubdetail `json:"kpi_sub_detail"`
}

type DataKpiSubdetail struct {
	IdSubDetail                   string  `json:"id_sub_detail"`
	IdSubKpi                      string  `json:"id_sub_kpi"`
	SubKpi                        string  `json:"sub_kpi"`
	Otomatis                      string  `json:"otomatis"`
	Polarisasi                    string  `json:"polarisasi"`
	IdPolarisasi                  string  `json:"id_polarisasi"`
	Capping                       string  `json:"capping"`
	Bobot                         float64 `json:"bobot"`
	Glossary                      string  `json:"glossary"`
	TargetTriwulan                string  `json:"target_triwulan"`
	TargetKuantitatifTriwulan     float64 `json:"target_kuantitatif_triwulan"`
	TargetTahunan                 string  `json:"target_tahunan"`
	TargetKuantitatifTahunan      float64 `json:"target_kuantitatif_tahunan"`
	TerdapatQualifier             string  `json:"terdapat_qualifier"`
	IdQualifier                   string  `json:"id_qualifier"`
	Qualifier                     string  `json:"qualifier"`
	DeskripsiQualifier            string  `json:"deskripsi_qualifier"`
	TargetQualifier               string  `json:"target_qualifier"`
	IdKeteranganProject           string  `json:"id_keterangan_project"`
	KeteranganProject             string  `json:"keterangan_project"`
	Realisasi                     string  `json:"realisasi"`
	RealisasiKuantitatif          float64 `json:"realisasi_kuantitatif"`
	RealisasiQualifier            string  `json:"realisasi_qualifier"`
	RealisasiKuantitatifQualifier string  `json:"realisasi_kuantitatif_qualifier"`
	RealisasiKeterangan           string  `json:"realisasi_keterangan"`
	RealisasiValidated            string  `json:"realisasi_validated"`
	RealisasiKuantitatifValidated string  `json:"realisasi_kuantitatif_validated"`
	Pencapaian                    float64 `json:"pencapaian"`
	Skor                          float64 `json:"skor"`
}

type DataResult struct {
	IdDetailResult   string `json:"id_detail_result"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaResult       string `json:"nama_result"`
	DeskripsiResult  string `json:"deskripsi_result"`
	RealisasiResult  string `json:"realisasi_result"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

type DataProcess struct {
	IdDetailProcess  string `json:"id_detail_process"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaProcess      string `json:"nama_process"`
	DeskripsiProcess string `json:"deskripsi_process"`
	RealisasiProcess string `json:"realisasi_process"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

type DataContext struct {
	IdDetailContext  string `json:"id_detail_context"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaContext      string `json:"nama_context"`
	DeskripsiContext string `json:"deskripsi_context"`
	RealisasiContext string `json:"realisasi_context"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

// =============================================================================
// REQUEST DTO
// =============================================================================

// ValidateRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/validate.
type ValidateRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl" validate:"required"`
	Tahun       string `json:"tahun" validate:"required"`
	Triwulan    string `json:"triwulan" validate:"required"`

	// Diisi service dari DB berdasarkan id_pengajuan, tidak boleh dari body.
	// Divisi Divisi `json:"-" validate:"-"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
}

// CreateRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/create.
// Submit realisasi ke approval (status → 3).
type CreateRealisasiKpiRequest struct {
	IdPengajuan           string                        `json:"id_pengajuan"            validate:"required"`
	Kostl                 string                        `json:"kostl"         			validate:"required"`
	Tahun                 string                        `json:"tahun"         			validate:"required"`
	Triwulan              string                        `json:"triwulan"      			validate:"required"`
	ApprovalListRealisasi []ApprovalUserRealisasiDetail `json:"approval_list_realisasi" validate:"required,min=1,dive"`
	// ApprovalPosisi        string `json:"approval_posisi"          validate:"required"`
	// ApprovalListRealisasi string `json:"approval_list_realisasi"  validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
}

// RevisionRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/revision.
type RevisionRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl" validate:"required"`
	Tahun       string `json:"tahun" validate:"required"`
	Triwulan    string `json:"triwulan" validate:"required"`

	// Diisi service dari DB berdasarkan id_pengajuan, tidak boleh dari body.
	// Divisi Divisi `json:"-" validate:"-"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
}

// ApproveRealisasiKpiRequest digunakan untuk endpoint POST /realisasi-kpi/approve.
type ApproveRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl"        validate:"required"`
	Tahun       string `json:"tahun"        validate:"required"`
	Triwulan    string `json:"triwulan"     validate:"required"`
	Catatan     string `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserRealisasi string `json:"approval_user_realisasi"`
	ApprovalNameRealisasi string `json:"approval_name_realisasi"`
}

// RejectRealisasiKpiRequest digunakan untuk endpoint POST /realisasi-kpi/reject.
type RejectRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl"        validate:"required"`
	Tahun       string `json:"tahun"        validate:"required"`
	Triwulan    string `json:"triwulan"     validate:"required"`
	Catatan     string `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserRealisasi string `json:"approval_user_realisasi"`
	ApprovalNameRealisasi string `json:"approval_name_realisasi"`
}

// GetAllRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all.
type GetAllRealisasiKpiRequest struct {
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllApprovalRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-approval.
type GetAllApprovalRealisasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	ApprovalUserRealisasi string `json:"approval_user_realisasi"`
}

// GetAllTolakanRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-tolakan.
type GetAllTolakanRealisasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	EntryUserRealisasi string `json:"entry_user_realisasi"`
}

// GetAllDaftarRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-daftar-realisasi.
type GetAllDaftarRealisasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarApprovalRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalRealisasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	ApprovalUserRealisasi string `json:"approval_user_realisasi"`
}

// GetDetailRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-detail.
type GetDetailRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// =============================================================================
// RESPONSE DTO
// =============================================================================

// ValidateRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/validate.
type ValidateRealisasiKpiResponse struct {
	IdPengajuan    string             `json:"id_pengajuan"`
	Tahun          string             `json:"tahun"`
	Triwulan       string             `json:"triwulan"`
	Divisi         Divisi             `json:"divisi"`
	EntryRealisasi EntryUserRealisasi `json:"entry_realisasi"`
	TotalKpi       int                `json:"total_kpi"`
	KpiList        []DataKpiDetail    `json:"kpi_list"`
	ResultList     []DataResult       `json:"result_list"`
	ProcessList    []DataProcess      `json:"process_list"`
	ContextList    []DataContext      `json:"context_list"`
}

// CreateRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/create.
type CreateRealisasiKpiResponse struct {
	IdPengajuan           string                  `json:"id_pengajuan"`
	Divisi                Divisi                  `json:"divisi"`
	Tahun                 string                  `json:"tahun"`
	Triwulan              string                  `json:"triwulan"`
	ApprovalListRealisasi []ApprovalUserRealisasi `json:"approval_list_realisasi"`
}

// RevisionRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/revision.
type RevisionRealisasiKpiResponse struct {
	IdPengajuan    string             `json:"id_pengajuan"`
	Tahun          string             `json:"tahun"`
	Triwulan       string             `json:"triwulan"`
	Divisi         Divisi             `json:"divisi"`
	EntryRealisasi EntryUserRealisasi `json:"entry_realisasi"`
	TotalKpi       int                `json:"total_kpi"`
	KpiList        []DataKpiDetail    `json:"kpi_list"`
	ResultList     []DataResult       `json:"result_list"`
	ProcessList    []DataProcess      `json:"process_list"`
	ContextList    []DataContext      `json:"context_list"`
}

// ApproveRealisasiKpiResponse adalah response untuk endpoint POST /realisasi-kpi/approve.
type ApproveRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
	Catatan     string `json:"catatan"`
}

// RejectRealisasiKpiResponse adalah response untuk endpoint POST /realisasi-kpi/reject.
type RejectRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
	Catatan     string `json:"catatan"`
}

// GetAllRealisasiKpiResponse adalah response satu record untuk /realisasi-kpi/get-all.
type GetAllRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllApprovalRealisasiKpiResponse adalah response satu record untuk /realisasi-kpi/get-all-approval.
type GetAllApprovalRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllTolakanRealisasiKpiResponse adalah response satu record untuk /realisasi-kpi/get-all-tolakan.
type GetAllTolakanRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarRealisasiKpiResponse adalah response satu record untuk /realisasi-kpi/get-all-daftar-realisasi.
type GetAllDaftarRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarApprovalRealisasiKpiResponse adalah response satu record untuk /realisasi-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetDetailRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/get-detail.
type GetDetailRealisasiKpiResponse struct {
	IdPengajuan           string                        `json:"id_pengajuan"`
	Tahun                 string                        `json:"tahun"`
	Triwulan              string                        `json:"triwulan"`
	Status                int                           `json:"status"`
	StatusDesc            string                        `json:"status_desc"`
	Divisi                DivisiOrgeh                   `json:"divisi"`
	EntryPenyusunan       EntryUserPenyusunan           `json:"entry_penyusunan"`
	EntryRealisasi        EntryUserRealisasi            `json:"entry_realisasi"`
	EntryValidasi         EntryUserValidasi             `json:"entry_validasi"`
	ApprovalPosisi        string                        `json:"approval_posisi"`
	ApprovalListRealisasi []ApprovalUserRealisasiDetail `json:"approval_list_realisasi"`
	Catatan               []CatatanTolakanEntry         `json:"catatan"`
	TotalBobot            string                        `json:"total_bobot"`
	TotalPencapaian       string                        `json:"total_pencapaian"`
	TotalKpi              int                           `json:"total_kpi"`
	KpiList               []DataKpiDetail               `json:"kpi_list"`
	TotalResult           int                           `json:"total_result"`
	ResultList            []DataResult                  `json:"result_list"`
	TotalProcess          int                           `json:"total_process"`
	ProcessList           []DataProcess                 `json:"process_list"`
	TotalContext          int                           `json:"total_context"`
	ContextList           []DataContext                 `json:"context_list"`
}

// =============================================================================
// INTERNAL ROW — hasil parse Excel realisasi
// =============================================================================

type RealisasiKpiRow = KpiRow

type RealisasiKpiSubDetailRow = KpiSubDetailRow

type KpiRow struct {
	KpiIndex int
	IdDetail string
	IdKpi    string
	Kpi      string
	Rumus    string
}

type KpiSubDetailRow struct {
	No                            int
	KPI                           string
	SubKPI                        string
	IdSubKpi                      string
	Polarisasi                    string
	IdPolarisasi                  string
	Capping                       string
	Bobot                         float64
	TargetTriwulan                string
	Qualifier                     string
	TargetQualifier               string
	Realisasi                     string
	RealisasiKuantitatif          float64
	RealisasiQualifierVal         string
	RealisasiKuantitatifQualifier string
	LinkDokumenSumber             *string
	IsTW24                        bool
	Result                        *string
	DeskripsiResult               *string
	RealisasiResult               *string
	LampiranEvidenceResult        *string
	Process                       *string
	DeskripsiProcess              *string
	RealisasiProcess              *string
	LampiranEvidenceProcess       *string
	Context                       *string
	DeskripsiContext              *string
	RealisasiContext              *string
	LampiranEvidenceContext       *string

	// Di-populate dari DB oleh service setelah parse (via enrichRowsFromDB)
	IdSubDetail               string
	IdDetail                  string
	IdQualifier               string
	TargetKuantitatifTriwulan float64
	Rumus                     string
	Pencapaian                float64
	Skor                      float64
}
