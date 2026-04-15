package dto

import "permen_api/pkg/excel"

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

type EntryUser struct {
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

type ApprovalUser struct {
	Userid     string `json:"userid"`
	Nama       string `json:"nama"`
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
	Posisi     string `json:"posisi"`
	Level      string `json:"level"`
	Fungsi     string `json:"fungsi"`
	Waktu      string `json:"waktu"`
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

// =============================================================================
// REQUEST DTO
// =============================================================================

// ValidatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiRequest struct {
	Divisi   Divisi `json:"divisi"    validate:"required"`
	Tahun    string `json:"tahun"     validate:"required"`
	Triwulan string `json:"triwulan"  validate:"required"`
	Kostl    string `json:"kostl"`
	KostlTx  string `json:"kostl_tx"`

	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// RevisionPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
	Divisi      Divisi `json:"divisi"      validate:"required"`
	Tahun       string `json:"tahun"       validate:"required"`
	Triwulan    string `json:"triwulan"    validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// CreatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiRequest struct {
	IdPengajuan  string         `json:"id_pengajuan"  validate:"required"`
	ApprovalList []ApprovalUser `json:"approval_list" validate:"required,min=1,dive"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	EntryUser string `json:"entry_user"`
	EntryName string `json:"entry_name"`
	EntryTime string `json:"entry_time"`
}

// BatalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/batal.
type BatalPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	User string `json:"user"`
}

// ApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/approval.
type ApprovalPenyusunanKpiRequest struct {
	IdPengajuan    string `json:"id_pengajuan"      validate:"required"`
	Status         string `json:"status"            validate:"required,oneof=approve reject"`
	ApprovalList   string `json:"approval_list"     validate:"required"`
	ApprovalPosisi string `json:"approval_posisi"`
	CatatanTolak   string `json:"catatan_tolak"`

	// Diisi handler dari header 'userq' (PERNR), tidak boleh dari body.
	User string `json:"user"`
}

// GetAllApprovalPenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/get-all-approval.
type GetAllApprovalPenyusunanKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUser string `json:"approval_user"`
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
	ApprovalUser string `json:"approval_user"`
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
type PenyusunanKpiRow = excel.KpiRow

// PenyusunanKpiSubDetailRow merepresentasikan 1 baris data dari file Excel
// yang sudah diparse dan divalidasi.
type PenyusunanKpiSubDetailRow = excel.KpiSubDetailRow

// =============================================================================
// RESPONSE DTO
// =============================================================================

// ValidatePenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/validate.
type ValidatePenyusunanKpiResponse struct {
	IDPengajuan string              `json:"id_pengajuan"`
	Tahun       string              `json:"tahun"`
	Triwulan    string              `json:"triwulan"`
	Divisi      Divisi              `json:"divisi"`
	Entry       EntryUser           `json:"entry"`
	TotalKpi    int                 `json:"total_kpi"`
	Kpi         []DataKpiDetail     `json:"kpi"`
	ResultList  []PenyusunanResult  `json:"result_list"`
	ProcessList []PenyusunanProcess `json:"process_list"`
	ContextList []PenyusunanContext `json:"context_list"`
}

// RevisionPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/revision.
type RevisionPenyusunanKpiResponse struct {
	IDPengajuan string              `json:"id_pengajuan"`
	Tahun       string              `json:"tahun"`
	Triwulan    string              `json:"triwulan"`
	Divisi      Divisi              `json:"divisi"`
	Entry       EntryUser           `json:"entry"`
	TotalKpi    int                 `json:"total_kpi"`
	Kpi         []DataKpiDetail     `json:"kpi"`
	ResultList  []PenyusunanResult  `json:"result_list"`
	ProcessList []PenyusunanProcess `json:"process_list"`
	ContextList []PenyusunanContext `json:"context_list"`
}

// CreatePenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiResponse struct {
	IdPengajuan  string         `json:"id_pengajuan"`
	ApprovalList []ApprovalUser `json:"approval_list"`
}

// BatalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/batal.
type BatalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
}

// ApprovalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/approval.
type ApprovalPenyusunanKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
}

// GetAllApprovalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-approval.
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

// GetAllTolakanPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
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

// GetAllDaftarPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
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

// GetAllDaftarApprovalPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-all-tolakan.
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

// GetDetailPenyusunanKpiResponse adalah response untuk endpoint POST /penyusunan-kpi/get-detail.
type GetDetailPenyusunanKpiResponse struct {
	IdPengajuan    string              `json:"id_pengajuan"`
	Tahun          string              `json:"tahun"`
	Triwulan       string              `json:"triwulan"`
	Status         string              `json:"status"`
	StatusDesc     string              `json:"status_desc"`
	Divisi         DivisiOrgeh         `json:"divisi"`
	Entry          EntryUser           `json:"entry"`
	ApprovalPosisi string              `json:"approval_posisi"`
	ApprovalList   []ApprovalUser      `json:"approval_list"`
	TotalKpi       int                 `json:"total_kpi"`
	Kpi            []DataKpiDetail     `json:"kpi"`
	TotalResult    int                 `json:"total_result"`
	ResultList     []PenyusunanResult  `json:"result_list"`
	TotalProcess   int                 `json:"total_process"`
	ProcessList    []PenyusunanProcess `json:"process_list"`
	TotalContext   int                 `json:"total_context"`
	ContextList    []PenyusunanContext `json:"context_list"`
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
