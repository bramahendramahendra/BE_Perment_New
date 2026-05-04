package dto

// =============================================================================
// DIVISI DTO
// =============================================================================

type Divisi struct {
	Kostl   string `json:"kostl"   validate:"required"`
	KostlTx string `json:"kostl_tx" validate:"required"`
}

type DivisiOrgeh struct {
	Kostl   string `json:"kostl"`
	KostlTx string `json:"kostl_tx"`
	Orgeh   string `json:"orgeh"`
	OrgehTx string `json:"orgeh_tx"`
}

// =============================================================================
// ENTRY DTO
// =============================================================================

type EntryUserPenyusunan struct {
	EntryUserPenyusunan string `json:"entry_user"`
	EntryNamePenyusunan string `json:"entry_name"`
	EntryTimePenyusunan string `json:"entry_time"`
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

// =============================================================================
// APPROVAL USER DTO
// =============================================================================

type ApprovalUser struct {
	Userid string `json:"userid"`
	Nama   string `json:"nama"`
	Posisi string `json:"posisi"`
}

type ApprovalUserDetail struct {
	Userid     string `json:"userid"`
	Nama       string `json:"nama"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
	Posisi     string `json:"posisi"`
	Level      string `json:"level"`
	Fungsi     string `json:"fungsi"`
	Waktu      string `json:"waktu"`
}

// =============================================================================
// CATATAN DTO
// =============================================================================

type Catatan struct {
	Fungsi    string `json:"fungsi"     validate:"required"`
	EntryNote string `json:"entry_note" validate:"required"`
}

type CatatanDetail struct {
	Fungsi    string `json:"fungsi"`
	EntryUser string `json:"entry_user"`
	EntryTime string `json:"entry_time"`
	EntryNote string `json:"entry_note"`
}

// =============================================================================
// KPI DTO
// =============================================================================

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

type DataResult struct {
	IdDetailResult  string `json:"id_detail_result"`
	Tahun           string `json:"tahun"`
	Triwulan        string `json:"triwulan"`
	NamaResult      string `json:"nama_result"`
	DeskripsiResult string `json:"deskripsi_result"`
}

type DataProcess struct {
	IdDetailProcess  string `json:"id_detail_process"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaProcess      string `json:"nama_process"`
	DeskripsiProcess string `json:"deskripsi_process"`
}

type DataContext struct {
	IdDetailContext  string `json:"id_detail_context"`
	Tahun            string `json:"tahun"`
	Triwulan         string `json:"triwulan"`
	NamaContext      string `json:"nama_context"`
	DeskripsiContext string `json:"deskripsi_context"`
}

// =============================================================================
// REQUEST DTO
// =============================================================================

// ValidatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiRequest struct {
	Divisi   Divisi `json:"divisi"    validate:"required"`
	Triwulan string `json:"triwulan"  validate:"required"`
	Tahun    string `json:"tahun"     validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUserPenyusunan string `json:"entry_user"`
	EntryNamePenyusunan string `json:"entry_name"`
	EntryTimePenyusunan string `json:"entry_time"`
}

// CreatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiRequest struct {
	IdPengajuan  string               `json:"id_pengajuan"  validate:"required"`
	Kostl        string               `json:"kostl"         validate:"required"`
	Triwulan     string               `json:"triwulan"      validate:"required"`
	Tahun        string               `json:"tahun"         validate:"required"`
	ApprovalList []ApprovalUserDetail `json:"approval_list" validate:"required,min=1,dive"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUserPenyusunan string `json:"entry_user"`
	EntryNamePenyusunan string `json:"entry_name"`
	EntryTimePenyusunan string `json:"entry_time"`
}

// RevisionPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl" validate:"required"`
	Triwulan    string `json:"triwulan" validate:"required"`
	Tahun       string `json:"tahun" validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUserPenyusunan string `json:"entry_user"`
	EntryNamePenyusunan string `json:"entry_name"`
	EntryTimePenyusunan string `json:"entry_time"`
}

// ApprovePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/approve.
type ApprovePenyusunanKpiRequest struct {
	IdPengajuan string  `json:"id_pengajuan" validate:"required"`
	Kostl       string  `json:"kostl"        validate:"required"`
	Triwulan    string  `json:"triwulan"     validate:"required"`
	Tahun       string  `json:"tahun"        validate:"required"`
	Catatan     Catatan `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserPenyusunan string `json:"approval_user"`
	ApprovalNamePenyusunan string `json:"approval_name"`
}

// RejectPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/reject.
type RejectPenyusunanKpiRequest struct {
	IdPengajuan string  `json:"id_pengajuan" validate:"required"`
	Kostl       string  `json:"kostl"        validate:"required"`
	Triwulan    string  `json:"triwulan"     validate:"required"`
	Tahun       string  `json:"tahun"        validate:"required"`
	Catatan     Catatan `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserPenyusunan string `json:"approval_user"`
	ApprovalNamePenyusunan string `json:"approval_name"`
}

// GetAllApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-approval.
type GetAllApprovalPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserPenyusunan string `json:"approval_user"`
}

// GetAllTolakanPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
type GetAllTolakanPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
type GetAllDaftarPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUser string `json:"approval_user"`
}

// GetDetailPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-detail.
type GetDetailPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// GetExcelPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-excel.
type GetExcelPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl"        validate:"required"`
	Triwulan    string `json:"triwulan"     validate:"required"`
	Tahun       string `json:"tahun"        validate:"required"`
}

// GetPdfPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-pdf.
type GetPdfPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Kostl       string `json:"kostl"        validate:"required"`
	Triwulan    string `json:"triwulan"     validate:"required"`
	Tahun       string `json:"tahun"        validate:"required"`
}

// =============================================================================
// RESPONSE DTO
// =============================================================================

// ValidatePenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiResponse struct {
	IDPengajuan     string              `json:"id_pengajuan"`
	Divisi          Divisi              `json:"divisi"`
	Triwulan        string              `json:"triwulan"`
	Tahun           string              `json:"tahun"`
	EntryPenyusunan EntryUserPenyusunan `json:"entry"`
	TotalKpi        int                 `json:"total_kpi"`
	KpiList         []DataKpiDetail     `json:"kpi"`
	ResultList      []DataResult        `json:"result_list"`
	ProcessList     []DataProcess       `json:"process_list"`
	ContextList     []DataContext       `json:"context_list"`
}

// CreatePenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiResponse struct {
	IdPengajuan            string         `json:"id_pengajuan"`
	Divisi                 Divisi         `json:"divisi"`
	Triwulan               string         `json:"triwulan"`
	Tahun                  string         `json:"tahun"`
	ApprovalListPenyusunan []ApprovalUser `json:"approval_list"`
}

// RevisionPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiResponse struct {
	IDPengajuan     string              `json:"id_pengajuan"`
	Divisi          Divisi              `json:"divisi"`
	Triwulan        string              `json:"triwulan"`
	Tahun           string              `json:"tahun"`
	EntryPenyusunan EntryUserPenyusunan `json:"entry"`
	TotalKpi        int                 `json:"total_kpi"`
	KpiList         []DataKpiDetail     `json:"kpi"`
	ResultList      []DataResult        `json:"result_list"`
	ProcessList     []DataProcess       `json:"process_list"`
	ContextList     []DataContext       `json:"context_list"`
}

// ApprovePenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/approve.
type ApprovePenyusunanKpiResponse struct {
	IdPengajuan string  `json:"id_pengajuan"`
	Status      string  `json:"status"`
	Catatan     Catatan `json:"catatan"`
}

// RejectPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/reject.
type RejectPenyusunanKpiResponse struct {
	IdPengajuan string  `json:"id_pengajuan"`
	Status      string  `json:"status"`
	Catatan     Catatan `json:"catatan"`
}

// GetAllApprovalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-approval.
type GetAllApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllTolakanPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
type GetAllTolakanPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
type GetAllDaftarPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarApprovalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
type GetAllDaftarApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetDetailPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-detail.
type GetDetailPenyusunanKpiResponse struct {
	IdPengajuan     string               `json:"id_pengajuan"`
	Triwulan        string               `json:"triwulan"`
	Tahun           string               `json:"tahun"`
	Status          string               `json:"status"`
	StatusDesc      string               `json:"status_desc"`
	Divisi          DivisiOrgeh          `json:"divisi"`
	EntryPenyusunan EntryUserPenyusunan  `json:"entry"`
	EntryRealisasi  EntryUserRealisasi   `json:"entry_realisasi"`
	EntryValidasi   EntryUserValidasi    `json:"entry_validasi"`
	ApprovalPosisi  string               `json:"approval_posisi"`
	ApprovalList    []ApprovalUserDetail `json:"approval_list"`
	Catatan         []CatatanDetail      `json:"catatan"`
	TotalKpi        int                  `json:"total_kpi"`
	KpiList         []DataKpiDetail      `json:"kpi"`
	TotalResult     int                  `json:"total_result"`
	ResultList      []DataResult         `json:"result_list"`
	TotalProcess    int                  `json:"total_process"`
	ProcessList     []DataProcess        `json:"process_list"`
	TotalContext    int                  `json:"total_context"`
	ContextList     []DataContext        `json:"context_list"`
}

// =============================================================================
// EXCEL ROW DTO
// =============================================================================

type PenyusunanKpiRow struct {
	KpiIndex int
	IdKpi    string
	Kpi      string
	Rumus    string
}

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
	IsTW24                    bool
	Result                    *string
	DeskripsiResult           *string
	Process                   *string
	DeskripsiProcess          *string
	Context                   *string
	DeskripsiContext          *string
}

// =============================================================================
// EXPORT DTO (digunakan oleh get-excel dan get-pdf)
// =============================================================================
type KpiSubDetailExportRow struct {
	No            int
	KpiNama       string
	Bobot         string
	TargetTahunan string
	Capping       string
}

type KpiExportData struct {
	NamaDivisi string
	Triwulan   string
	Tahun      string
	Rows       []KpiSubDetailExportRow
}
