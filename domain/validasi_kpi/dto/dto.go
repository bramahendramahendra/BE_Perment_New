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
	LinkDokumenSumber   string             `json:"link_dokumen_sumber"`
	TotalSubKpi         int                `json:"total_sub_kpi"`
	KpiSubDetail        []DataKpiSubdetail `json:"kpi_sub_detail"`
}

type DataKpiSubdetail struct {
	IdSubDetail                      string      `json:"id_sub_detail"`
	IdSubKpi                         string      `json:"id_sub_kpi"`
	SubKpi                           string      `json:"sub_kpi"`
	Otomatis                         string      `json:"otomatis"`
	Polarisasi                       string      `json:"polarisasi"`
	IdPolarisasi                     string      `json:"id_polarisasi"`
	Capping                          string      `json:"capping"`
	Bobot                            float64     `json:"bobot"`
	Glossary                         string      `json:"glossary"`
	TargetTriwulan                   string      `json:"target_triwulan"`
	TargetKuantitatifTriwulan        float64     `json:"target_kuantitatif_triwulan"`
	TargetTahunan                    string      `json:"target_tahunan"`
	TargetKuantitatifTahunan         float64     `json:"target_kuantitatif_tahunan"`
	TerdapatQualifier                string      `json:"terdapat_qualifier"`
	Qualifier                        string      `json:"qualifier"`
	DeskripsiQualifier               string      `json:"deskripsi_qualifier"`
	TargetQualifier                  string      `json:"target_qualifier"`
	IdKeteranganProject              string      `json:"id_keterangan_project"`
	KeteranganProject                string      `json:"keterangan_project"`
	Realisasi                        string      `json:"realisasi"`
	RealisasiKuantitatif             float64     `json:"realisasi_kuantitatif"`
	RealisasiQualifier               string      `json:"realisasi_qualifier"`
	RealisasiKuantitatifQualifier    string      `json:"realisasi_kuantitatif_qualifier"`
	RealisasiKeterangan              string      `json:"realisasi_keterangan"`
	RealisasiValidated               string      `json:"realisasi_validated"`
	RealisasiKuantitatifValidated    string      `json:"realisasi_kuantitatif_validated"`
	ValidasiKeterangan               string      `json:"validasi_keterangan"`
	Pencapaian                       float64     `json:"pencapaian"`
	Skor                             float64     `json:"skor"`
	PencapaianQualifierValidated     interface{} `json:"pencapaian_qualifier_validated"`
	PencapaianPostQualifierValidated interface{} `json:"pencapaian_post_qualifier_validated"`
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

type DataValidasiQualifierOverall struct {
	IdKpiQualifier                     string `json:"id_kpi_qualifier"`
	KpiQualifier                       string `json:"kpi_qualifier"`
	Parameter                          string `json:"parameter"`
	Deskripsi                          string `json:"deskripsi"`
	BobotPengurang                     string `json:"bobot_pengurang"`
	Tahun                              string `json:"tahun"`
	RealisasiOverallValidated          string `json:"realisasi_overall_validated"`
	RealisasiQualifierOverallValidated string `json:"realisasi_qualifier_overall_validated"`
}

// =============================================================================
// REQUEST DTO
// =============================================================================

// InputValidasiKpiRequest adalah request untuk endpoint POST /realisasi-kpi/validate.
type InputValidasiKpiRequest struct {
	IdPengajuan                  string                         `json:"id_pengajuan"                    validate:"required"`
	Kostl                        string                         `json:"kostl"                   validate:"required"`
	Triwulan                     string                         `json:"triwulan"                validate:"required"`
	Tahun                        string                         `json:"tahun"                   validate:"required"`
	ApprovalListValidasi         []ApprovalUserDetail           `json:"approval_list_validasi"          validate:"required"`
	TotalBobot                   string                         `json:"total_bobot"                     validate:"required"`
	TotalPencapaian              string                         `json:"total_pencapaian"                validate:"required"`
	TotalBobotPengurang          string                         `json:"total_bobot_pengurang"           validate:"required"`
	TotalPencapaianPost          string                         `json:"total_pencapaian_post"           validate:"required"`
	Kpi                          []DataKpiDetail                `json:"kpi"                             validate:"required"`
	DataValidasiQualifierOverall []DataValidasiQualifierOverall `json:"data_validasi_qualifier_overall" validate:"required"`
	LampiranValidasi             []string                       `json:"lampiran_validasi"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	EntryUserValidasi string `json:"entry_user_validasi"`
	EntryNameValidasi string `json:"entry_name_validasi"`
	EntryTimeValidasi string `json:"entry_time_validasi"`
}

// ApproveValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/approve.
type ApproveValidasiKpiRequest struct {
	IdPengajuan string  `json:"id_pengajuan" validate:"required"`
	Kostl       string  `json:"kostl"        validate:"required"`
	Triwulan    string  `json:"triwulan"     validate:"required"`
	Tahun       string  `json:"tahun"        validate:"required"`
	Catatan     Catatan `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserValidasi string `json:"approval_user_validasi"`
	ApprovalNameValidasi string `json:"approval_name_validasi"`
}

// RejectValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/reject.
type RejectValidasiKpiRequest struct {
	IdPengajuan string  `json:"id_pengajuan" validate:"required"`
	Kostl       string  `json:"kostl"        validate:"required"`
	Triwulan    string  `json:"triwulan"     validate:"required"`
	Tahun       string  `json:"tahun"        validate:"required"`
	Catatan     Catatan `json:"catatan"      validate:"required"`

	// Diisi handler dari header 'userq', tidak boleh dari body.
	ApprovalUserValidasi string `json:"approval_user_validasi"`
	ApprovalNameValidasi string `json:"approval_name_validasi"`
}

// GetAllValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-validasi.
type GetAllValidasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllApprovalValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-approval.
type GetAllApprovalValidasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUserValidasi string `json:"approval_user_validasi"`
}

// GetAllTolakanValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-tolakan.
type GetAllTolakanValidasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	EntryUserRealisasi string `json:"entry_user_validasi"`
}

// GetAllDaftarPValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-daftar-penyusunan.
type GetAllDaftarPValidasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarApprovalValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalValidasiKpiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUserValidasi string `json:"approval_user_validasi"`
}

// GetDetailValidasiKpiRequest adalah request untuk endpoint POST /validasi-kpi/get-detail.
type GetDetailValidasiKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}

// =============================================================================
// RESPONSE DTO
// =============================================================================

// InputValidasiKpiResponse adalah response untuk endpoint POST /validasi-kpi/input.
type InputValidasiKpiResponse struct {
	IdPengajuan                  string                         `json:"id_pengajuan"`
	Divisi                       Divisi                         `json:"divisi"`
	Triwulan                     string                         `json:"triwulan"`
	Tahun                        string                         `json:"tahun"`
	EntryValidasi                EntryUserValidasi              `json:"entry_validasi"`
	ApprovalPosisi               string                         `json:"approval_posisi"`
	ApprovalListValidasi         []ApprovalUser                 `json:"approval_list_validasi"`
	TotalBobot                   string                         `json:"total_bobot"`
	TotalPencapaian              string                         `json:"total_pencapaian"`
	TotalBobotPengurang          interface{}                    `json:"total_bobot_pengurang"`
	TotalPencapaianPost          interface{}                    `json:"total_pencapaian_post"`
	KpiList                      []DataKpiDetail                `json:"kpi"                             validate:"required"`
	DataValidasiQualifierOverall []DataValidasiQualifierOverall `json:"data_validasi_qualifier_overall" validate:"required"`
	LampiranValidasi             []string                       `json:"lampiran_validasi"`
}

// ApproveValidasiKpiResponse adalah response untuk endpoint POST /validasi-kpi/approve.
type ApproveValidasiKpiResponse struct {
	IdPengajuan string  `json:"id_pengajuan"`
	Status      string  `json:"status"`
	Catatan     Catatan `json:"catatan"`
}

// RejectValidasiKpiResponse adalah response untuk endpoint POST /validasi-kpi/reject.
type RejectValidasiKpiResponse struct {
	IdPengajuan string  `json:"id_pengajuan"`
	Status      string  `json:"status"`
	Catatan     Catatan `json:"catatan"`
}

// GetAllValidasiKpiResponse adalah response satu record untuk /validasi-kpi/get-all.
type GetAllValidasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllApprovalValidasiKpiResponse adalah response satu record untuk /validasi-kpi/get-all-approval.
type GetAllApprovalValidasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllTolakanValidasiKpiResponse adalah response satu record untuk /validasi-kpi/get-all-tolakan.
type GetAllTolakanValidasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarValidasiKpiResponse adalah response satu record untuk /validasi-kpi/get-all-daftar-validasi.
type GetAllDaftarValidasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetAllDaftarApprovalValidasiKpiResponse adalah response satu record untuk /validasi-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalValidasiKpiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Triwulan    string `json:"triwulan"`
	Tahun       string `json:"tahun"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}

// GetDetailValidasiKpiResponse adalah response untuk endpoint POST /validasi-kpi/get-detail.
type GetDetailValidasiKpiResponse struct {
	IdPengajuan              string                         `json:"id_pengajuan"`
	Triwulan                 string                         `json:"triwulan"`
	Tahun                    string                         `json:"tahun"`
	Status                   int                            `json:"status"`
	StatusDesc               string                         `json:"status_desc"`
	Divisi                   DivisiOrgeh                    `json:"divisi"`
	EntryPenyusunan          EntryUserPenyusunan            `json:"entry_penyusunan"`
	EntryRealisasi           EntryUserRealisasi             `json:"entry_realisasi"`
	EntryValidasi            EntryUserValidasi              `json:"entry_validasi"`
	ApprovalPosisi           string                         `json:"approval_posisi"`
	ApprovalListValidasi     []ApprovalUserDetail           `json:"approval_list_validasi"`
	Catatan                  []CatatanDetail                `json:"catatan"`
	TotalBobot               string                         `json:"total_bobot"`
	TotalPencapaian          string                         `json:"total_pencapaian"`
	TotalBobotPengurang      string                         `json:"total_bobot_pengurang"`
	TotalPencapaianPost      string                         `json:"total_pencapaian_post"`
	LampiranValidasi         []string                       `json:"lampiran_validasi"`
	TotalKpi                 int                            `json:"total_kpi"`
	KpiList                  []DataKpiDetail                `json:"kpi_list"`
	TotalResult              int                            `json:"total_result"`
	ResultList               []DataResult                   `json:"result_list"`
	TotalProcess             int                            `json:"total_process"`
	ProcessList              []DataProcess                  `json:"process_list"`
	TotalContext             int                            `json:"total_context"`
	ContextList              []DataContext                  `json:"context_list"`
	QualifierOverallValidasi []DataValidasiQualifierOverall `json:"data_validasi_qualifier_overall"`
}
