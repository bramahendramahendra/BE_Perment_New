package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

// InputValidasiRequest adalah request untuk endpoint POST /validasi-kpi/input.
// Menyimpan data validasi KPI sebagai draft (status → 6).
type InputValidasiRequest struct {
	IdPengajuan                  string        `json:"id_pengajuan"                    validate:"required"`
	EntryName                    string        `json:"entry_name"                      validate:"required"`
	ApprovalPosisi               string        `json:"approval_posisi"                 validate:"required"`
	ApprovalListValidasi         string        `json:"approval_list_validasi"          validate:"required"`
	TotalBobot                   string        `json:"total_bobot"                     validate:"required"`
	TotalPencapaian              string        `json:"total_pencapaian"                validate:"required"`
	TotalBobotPengurang          interface{}   `json:"total_bobot_pengurang"           validate:"required"`
	TotalPencapaianPost          interface{}   `json:"total_pencapaian_post"           validate:"required"`
	DataValidasi                 []KpiSubDetail `json:"data_validasi"                  validate:"required"`
	DataValidasiQualifierOverall interface{}   `json:"data_validasi_qualifier_overall" validate:"required"`
	LampiranValidasi             string        `json:"lampiran_validasi"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	EntryUserValidasi string `json:"-"`
	EntryNameValidasi string `json:"-"`
}

// ApprovalValidasiRequest adalah request untuk endpoint POST /validasi-kpi/approval.
type ApprovalValidasiRequest struct {
	IdPengajuan    string `json:"id_pengajuan"  validate:"required"`
	Status         string `json:"status"        validate:"required,oneof=approve reject"`
	ApprovalList   string `json:"approval_list" validate:"required"`
	ApprovalPosisi string `json:"approval_posisi"`
	CatatanTolakan string `json:"catatan_tolakan"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUser string `json:"-"`
	ApprovalName string `json:"-"`
}

// ApproveValidasiRequest adalah request untuk endpoint POST /validasi-kpi/approve.
type ApproveValidasiRequest struct {
	IdPengajuan    string `json:"id_pengajuan"  validate:"required"`
	ApprovalList   string `json:"approval_list" validate:"required"`
	ApprovalPosisi string `json:"approval_posisi"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUser string `json:"-"`
	ApprovalName string `json:"-"`
}

// RejectValidasiRequest adalah request untuk endpoint POST /validasi-kpi/reject.
type RejectValidasiRequest struct {
	IdPengajuan    string `json:"id_pengajuan"  validate:"required"`
	ApprovalList   string `json:"approval_list" validate:"required"`
	CatatanTolakan string `json:"catatan_tolakan"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUser string `json:"-"`
	ApprovalName string `json:"-"`
}

// GetAllApprovalValidasiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-approval.
type GetAllApprovalValidasiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUser string `json:"-"`
}

// GetAllTolakanValidasiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-tolakan.
type GetAllTolakanValidasiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarPenyusunanValidasiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-daftar-penyusunan.
type GetAllDaftarPenyusunanValidasiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// GetAllDaftarApprovalValidasiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-daftar-approval.
type GetAllDaftarApprovalValidasiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	ApprovalUser string `json:"-"`
}

// GetAllValidasiRequest adalah request untuk endpoint POST /validasi-kpi/get-all-validasi.
type GetAllValidasiRequest struct {
	Divisi   string `json:"divisi"`
	Tahun    string `json:"tahun"`
	Triwulan string `json:"triwulan"`
	Status   string `json:"status"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

// ValidasiBatalRequest adalah request untuk endpoint POST /validasi-kpi/batal.
type ValidasiBatalRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`

	// Di-populate dari header "userq" oleh handler, tidak dari body
	User string `json:"-"`
}

// =============================================================================
// SUB DTO
// =============================================================================

// KpiSubDetail adalah data validasi per sub KPI dalam InputValidasiRequest.
type KpiSubDetail struct {
	KeyPengajuan                    string `json:"KeyPengajuan"`
	KeyDetail                       string `json:"KeyDetail"`
	KeySubDetail                    string `json:"KeySubDetail"`
	TargetTriwulanValidated         string `json:"TargetTriwulanValidated"`
	TargetKuantitatifValidated      string `json:"TargetKuantitatifValidated"`
	RealisasiValidated              string `json:"RealisasiValidated"`
	RealisasiKuantitatifValidated   string `json:"RealisasiKuantitatifValidated"`
	Pencapaian                      string `json:"Pencapaian"`
	Skor                            string `json:"Skor"`
	ValidasiKeterangan              string `json:"ValidasiKeterangan"`
	PencapaianQualifierValidated    interface{} `json:"PencapaianQualifierValidated"`
	PencapaianPostQualifierValidated interface{} `json:"PencapaianPostQualifierValidated"`
	TargetQualifierValidated        interface{} `json:"TargetQualifierValidated"`
}

// =============================================================================
// RESPONSE DTO
// =============================================================================

// InputValidasiResponse adalah response untuk endpoint POST /validasi-kpi/input.
type InputValidasiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
}

// ApprovalValidasiResponse adalah response untuk endpoint POST /validasi-kpi/approval.
type ApprovalValidasiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
}

// ValidasiBatalResponse adalah response untuk endpoint POST /validasi-kpi/batal.
type ValidasiBatalResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
}

// ApproveValidasiResponse adalah response untuk endpoint POST /validasi-kpi/approve.
type ApproveValidasiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
}

// RejectValidasiResponse adalah response untuk endpoint POST /validasi-kpi/reject.
type RejectValidasiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Status      string `json:"status"`
}

// GetAllValidasiResponse adalah response untuk endpoint POST /validasi-kpi/get-all-*.
type GetAllValidasiResponse struct {
	IdPengajuan string `json:"id_pengajuan"`
	Tahun       string `json:"tahun"`
	Triwulan    string `json:"triwulan"`
	KostlTx     string `json:"kostl_tx"`
	OrgehTx     string `json:"orgeh_tx"`
	StatusDesc  string `json:"status_desc"`
}
