package dto

// =============================================================================
// INTERNAL ROW — hasil parse Excel realisasi
// =============================================================================

// RealisasiKpiRow merepresentasikan satu baris data dari file Excel realisasi
// yang sudah diparse, divalidasi, dan diperkaya dari DB (lookup id_sub_detail, target, dll).

type KpiRow struct {
	KpiIndex int
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
	IsTW24                        bool
	Result                        *string
	DeskripsiResult               *string
	RealisasiResult               *string
	LinkResult                    *string
	Process                       *string
	DeskripsiProcess              *string
	RealisasiProcess              *string
	LinkProcess                   *string
	Context                       *string
	DeskripsiContext              *string
	RealisasiContext              *string
	LinkContext                   *string

	// Di-populate dari DB oleh service setelah parse (via enrichRowsFromDB)
	IdSubDetail               string
	IdDetail                  string
	TargetKuantitatifTriwulan float64
	Rumus                     string
	Pencapaian                float64
	Skor                      float64
}

// =============================================================================
// INTERNAL LIST — result / process / context untuk update DB (TW2 & TW4)
// =============================================================================

// RealisasiResult berisi data untuk UPDATE data_challenge_detail (result).
type RealisasiResult struct {
	IdDetailResult  string // = id_sub_detail (FK di data_challenge_detail)
	RealisasiResult string
	LinkResult      string
}

// RealisasiProcess berisi data untuk UPDATE data_method_detail (process).
type RealisasiProcess struct {
	IdDetailProcess  string // = id_sub_detail (FK di data_method_detail)
	RealisasiProcess string
	LinkProcess      string
}

// RealisasiContext berisi data untuk UPDATE data_challenge_detail (context).
type RealisasiContext struct {
	IdDetailContext  string // = id_sub_detail (FK di data_challenge_detail)
	RealisasiContext string
	LinkContext      string
}

type RealisasiKpiRow struct {
	RowIndex int
	No       int

	// Dari Excel (pre-filled oleh template)
	KPI             string
	SubKPI          string
	Polarisasi      string
	Capping         string
	Bobot           float64
	TargetTriwulan  string
	Qualifier       string
	TargetQualifier string

	// Dari Excel (user input)
	Realisasi                     string
	RealisasiKuantitatif          float64
	RealisasiQualifierVal         string
	RealisasiKuantitatifQualifier string

	// Extended kolom (semua triwulan: N–S atau N–Y)
	Result           *string
	DeskripsiResult  *string
	Process          *string
	DeskripsiProcess *string
	Context          *string
	DeskripsiContext *string

	// Extended TW2/TW4 saja
	RealisasiResult  *string
	LinkResult       *string
	RealisasiProcess *string
	LinkProcess      *string
	RealisasiContext *string
	LinkContext      *string

	IsTW24 bool

	// Dari DB (di-populate oleh service setelah parse)
	IdSubDetail               string
	IdDetail                  string
	TargetKuantitatifTriwulan float64
	Rumus                     string // id_polarisasi: "1"=Maximize, "0"=Minimize

	// Hasil kalkulasi
	Pencapaian float64
	Skor       float64
}

// =============================================================================
// REQUEST DTO
// =============================================================================

// ValidateRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/validate.
type ValidateRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`

	// Di-populate dari DB berdasarkan IdPengajuan oleh service, tidak dari body
	Triwulan string `json:"triwulan"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// RevisionRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/revision.
type RevisionRealisasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`

	// Di-populate dari DB berdasarkan IdPengajuan oleh service, tidak dari body
	Triwulan string `json:"-"`

	EntryUser string `json:"-"`
	EntryName string `json:"-"`
	EntryTime string `json:"-"`
}

// CreateRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/create.
// Submit realisasi ke approval (status → 3).
type CreateRealisasiKpiRequest struct {
	IdPengajuan           string `json:"id_pengajuan"             validate:"required"`
	ApprovalPosisi        string `json:"approval_posisi"          validate:"required"`
	ApprovalListRealisasi string `json:"approval_list_realisasi"  validate:"required"`

	User string `json:"-"` // di-populate dari header
}

// ApprovalRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/approval.
type ApprovalRealisasiKpiRequest struct {
	IdPengajuan    string `json:"id_pengajuan"   validate:"required"`
	Status         string `json:"status"         validate:"required,oneof=approve reject"`
	ApprovalList   string `json:"approval_list"  validate:"required"`
	ApprovalPosisi string `json:"approval_posisi"`
	CatatanTolak   string `json:"catatan_tolakan"`

	User string `json:"-"`
}

// GetAllApprovalRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-approval.
type GetAllApprovalRealisasiKpiRequest struct {
	ApprovalUser string `json:"approval_user" validate:"required"`
	Divisi       string `json:"divisi"`
	Tahun        string `json:"tahun"`
	Triwulan     string `json:"triwulan"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
}

// GetAllTolakanRealisasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/get-all-tolakan.
type GetAllTolakanRealisasiKpiRequest struct {
	EntryUserRealisasi string `json:"entry_user_realisasi" validate:"required"`
	Divisi             string `json:"divisi"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	Page               int    `json:"page"`
	Limit              int    `json:"limit"`
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
	IdPengajuan string                  `json:"id_pengajuan"`
	Tahun       string                  `json:"tahun"`
	Triwulan    string                  `json:"triwulan"`
	Divisi      DivisiRealisasiResponse `json:"divisi"`
	Entry       EntryRealisasiResponse  `json:"entry"`
	TotalSubKpi int                     `json:"total_sub_kpi"`
	SubKpiList  []RealisasiSubKpiDetail `json:"sub_kpi_list"`
	ResultList  []RealisasiResult       `json:"result_list"`
	ProcessList []RealisasiProcess      `json:"process_list"`
	ContextList []RealisasiContext      `json:"context_list"`
}

// RevisionRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/revision.
type RevisionRealisasiKpiResponse struct {
	IdPengajuan string                  `json:"id_pengajuan"`
	Tahun       string                  `json:"tahun"`
	Triwulan    string                  `json:"triwulan"`
	TotalSubKpi int                     `json:"total_sub_kpi"`
	SubKpiList  []RealisasiSubKpiDetail `json:"sub_kpi_list"`
	ResultList  []RealisasiResult       `json:"result_list"`
	ProcessList []RealisasiProcess      `json:"process_list"`
	ContextList []RealisasiContext      `json:"context_list"`
}

// CreateRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/create.
type CreateRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Message     string `json:"message"`
}

// ApprovalRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/approval.
type ApprovalRealisasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Message     string `json:"message"`
}

// DivisiRealisasiResponse digunakan di dalam response validate/revision.
type DivisiRealisasiResponse struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostl_tx"`
}

// EntryRealisasiResponse digunakan di dalam response validate/revision.
type EntryRealisasiResponse struct {
	EntryUser string `json:"entry_user_realisasi"`
	EntryName string `json:"entry_name_realisasi"`
	EntryTime string `json:"entry_time_realisasi"`
}

// RealisasiSubKpiDetail adalah satu baris preview sub KPI realisasi.
type RealisasiSubKpiDetail struct {
	IdSubDetail                   string  `json:"id_sub_detail"`
	IdDetail                      string  `json:"id_detail"`
	KPI                           string  `json:"kpi"`
	SubKPI                        string  `json:"sub_kpi"`
	Polarisasi                    string  `json:"polarisasi"`
	Capping                       string  `json:"capping"`
	Bobot                         float64 `json:"bobot"`
	TargetTriwulan                string  `json:"target_triwulan"`
	TargetKuantitatifTriwulan     float64 `json:"target_kuantitatif_triwulan"`
	Qualifier                     string  `json:"qualifier"`
	TargetQualifier               string  `json:"target_qualifier"`
	Realisasi                     string  `json:"realisasi"`
	RealisasiKuantitatif          float64 `json:"realisasi_kuantitatif"`
	RealisasiQualifier            string  `json:"realisasi_qualifier"`
	RealisasiKuantitatifQualifier string  `json:"realisasi_kuantitatif_qualifier"`
	Pencapaian                    float64 `json:"pencapaian"`
	Skor                          float64 `json:"skor"`
}

// =============================================================================
// RESPONSE GetAll
// =============================================================================

// GetAllApprovalRealisasiKpiResponse adalah response satu record untuk get-all-approval.
type GetAllApprovalRealisasiKpiResponse struct {
	IdPengajuan        string `json:"id_pengajuan"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	Kostl              string `json:"kostl"`
	KostlTx            string `json:"kostl_tx"`
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
	Status             int    `json:"status"`
	StatusDesc         string `json:"status_desc"`
}

// GetAllTolakanRealisasiKpiResponse adalah response satu record untuk get-all-tolakan.
type GetAllTolakanRealisasiKpiResponse struct {
	IdPengajuan        string `json:"id_pengajuan"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	Kostl              string `json:"kostl"`
	KostlTx            string `json:"kostl_tx"`
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
	CatatanTolakan     string `json:"catatan_tolakan"`
	Status             int    `json:"status"`
	StatusDesc         string `json:"status_desc"`
}

// GetAllDaftarRealisasiKpiResponse adalah response satu record untuk get-all-daftar-realisasi.
type GetAllDaftarRealisasiKpiResponse struct {
	IdPengajuan        string `json:"id_pengajuan"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	Kostl              string `json:"kostl"`
	KostlTx            string `json:"kostl_tx"`
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
	Status             int    `json:"status"`
	StatusDesc         string `json:"status_desc"`
	TotalBobot         string `json:"total_bobot"`
	TotalPencapaian    string `json:"total_pencapaian"`
}

// GetAllDaftarApprovalRealisasiKpiResponse adalah response satu record untuk get-all-daftar-approval.
type GetAllDaftarApprovalRealisasiKpiResponse struct {
	IdPengajuan        string `json:"id_pengajuan"`
	Tahun              string `json:"tahun"`
	Triwulan           string `json:"triwulan"`
	Kostl              string `json:"kostl"`
	KostlTx            string `json:"kostl_tx"`
	EntryUserRealisasi string `json:"entry_user_realisasi"`
	EntryNameRealisasi string `json:"entry_name_realisasi"`
	EntryTimeRealisasi string `json:"entry_time_realisasi"`
	Status             int    `json:"status"`
	StatusDesc         string `json:"status_desc"`
}

// =============================================================================
// RESPONSE GetDetail
// =============================================================================

// GetDetailRealisasiKpiResponse adalah response untuk endpoint /realisasi-kpi/get-detail.
type GetDetailRealisasiKpiResponse struct {
	IdPengajuan           string                   `json:"id_pengajuan"`
	Tahun                 string                   `json:"tahun"`
	Triwulan              string                   `json:"triwulan"`
	Kostl                 string                   `json:"kostl"`
	KostlTx               string                   `json:"kostl_tx"`
	Orgeh                 string                   `json:"orgeh"`
	OrgehTx               string                   `json:"orgeh_tx"`
	Status                int                      `json:"status"`
	StatusDesc            string                   `json:"status_desc"`
	EntryPenyusunan       EntryDetailResponse      `json:"entry_penyusunan"`
	EntryRealisasi        EntryDetailResponse      `json:"entry_realisasi"`
	ApprovalList          []ApprovalDetailResponse `json:"approval_list"`
	ApprovalListRealisasi []ApprovalDetailResponse `json:"approval_list_realisasi"`
	CatatanTolakan        string                   `json:"catatan_tolakan"`
	TotalBobot            string                   `json:"total_bobot"`
	TotalPencapaian       string                   `json:"total_pencapaian"`
	TotalSubKpi           int                      `json:"total_sub_kpi"`
	KpiList               []DetailKpiRealisasi     `json:"kpi_list"`
	ContextList           []DetailContextRealisasi `json:"context_list"`
	ProcessList           []DetailProcessRealisasi `json:"process_list"`
}

// EntryDetailResponse adalah nested entry user dalam GetDetail response.
type EntryDetailResponse struct {
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// ApprovalDetailResponse adalah satu entri dalam approval_list.
type ApprovalDetailResponse struct {
	Userid     string `json:"userid"`
	Nama       string `json:"nama"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
	Posisi     string `json:"posisi"`
	Level      string `json:"level"`
	Fungsi     string `json:"fungsi"`
	Waktu      string `json:"waktu"`
}

// DetailKpiRealisasi adalah nested KPI + sub_detail dalam GetDetail.
type DetailKpiRealisasi struct {
	IdDetail      string                  `json:"id_detail"`
	IdKpi         string                  `json:"id_kpi"`
	Kpi           string                  `json:"kpi"`
	Rumus         string                  `json:"rumus"`
	TotalSubKpi   int                     `json:"total_sub_kpi"`
	SubDetailList []DetailSubKpiRealisasi `json:"sub_detail_list"`
}

// DetailSubKpiRealisasi adalah satu baris sub KPI dalam GetDetail.
type DetailSubKpiRealisasi struct {
	IdSubDetail                   string `json:"id_sub_detail"`
	IdKpi                         string `json:"id_kpi"`
	SubKpi                        string `json:"sub_kpi"`
	Otomatis                      string `json:"otomatis"`
	Bobot                         string `json:"bobot"`
	Capping                       string `json:"capping"`
	TargetTriwulan                string `json:"target_triwulan"`
	TargetKuantitatifTriwulan     string `json:"target_kuantitatif_triwulan"`
	TargetTahunan                 string `json:"target_tahunan"`
	TargetKuantitatifTahunan      string `json:"target_kuantitatif_tahunan"`
	Realisasi                     string `json:"realisasi"`
	RealisasiKuantitatif          string `json:"realisasi_kuantitatif"`
	RealisasiKeterangan           string `json:"realisasi_keterangan"`
	RealisasiValidated            string `json:"realisasi_validated"`
	RealisasiKuantitatifValidated string `json:"realisasi_kuantitatif_validated"`
	Pencapaian                    string `json:"pencapaian"`
	Skor                          string `json:"skor"`
	DeskripsiGlossary             string `json:"deskripsi_glossary"`
	ItemQualifier                 string `json:"item_qualifier"`
	DeskripsiQualifier            string `json:"deskripsi_qualifier"`
	TargetQualifier               string `json:"target_qualifier"`
	IdQualifier                   string `json:"id_qualifier"`
	RealisasiQualifier            string `json:"realisasi_qualifier"`
	RealisasiKuantitatifQualifier string `json:"realisasi_kuantitatif_qualifier"`
}

// DetailContextRealisasi adalah satu baris context/challenge dalam GetDetail.
type DetailContextRealisasi struct {
	IdDetailChallenge  string `json:"id_detail_challenge"`
	NamaChallenge      string `json:"nama_challenge"`
	DeskripsiChallenge string `json:"deskripsi_challenge"`
	RealisasiChallenge string `json:"realisasi_challenge"`
	LampiranEvidence   string `json:"lampiran_evidence"`
}

// DetailProcessRealisasi adalah satu baris process/method dalam GetDetail.
type DetailProcessRealisasi struct {
	IdDetailMethod   string `json:"id_detail_method"`
	NamaMethod       string `json:"nama_method"`
	DeskripsiMethod  string `json:"deskripsi_method"`
	RealisasiMethod  string `json:"realisasi_method"`
	LampiranEvidence string `json:"lampiran_evidence"`
}

// =============================================================================
// MODEL — mirroring dari domain/penyusunan_kpi/model
// =============================================================================

// DataKpiRealisasi adalah model scan untuk query GetAll.
type DataKpiRealisasi struct {
	IdPengajuan           string `gorm:"column:id_pengajuan"`
	Tahun                 string `gorm:"column:tahun"`
	Triwulan              string `gorm:"column:triwulan"`
	Kostl                 string `gorm:"column:kostl"`
	KostlTx               string `gorm:"column:kostl_tx"`
	Orgeh                 string `gorm:"column:orgeh"`
	OrgehTx               string `gorm:"column:orgeh_tx"`
	EntryUser             string `gorm:"column:entry_user"`
	EntryName             string `gorm:"column:entry_name"`
	EntryTime             string `gorm:"column:entry_time"`
	ApprovalPosisi        string `gorm:"column:approval_posisi"`
	ApprovalList          string `gorm:"column:approval_list"`
	Status                int    `gorm:"column:status"`
	StatusDesc            string `gorm:"column:status_desc"`
	EntryUserRealisasi    string `gorm:"column:entry_user_realisasi"`
	EntryNameRealisasi    string `gorm:"column:entry_name_realisasi"`
	EntryTimeRealisasi    string `gorm:"column:entry_time_realisasi"`
	ApprovalListRealisasi string `gorm:"column:approval_list_realisasi"`
	CatatanTolakan        string `gorm:"column:catatan_tolakan"`
	TotalBobot            string `gorm:"column:total_bobot"`
	TotalPencapaian       string `gorm:"column:total_pencapaian"`
}
