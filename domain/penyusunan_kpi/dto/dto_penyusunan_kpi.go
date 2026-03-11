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
	Divisi        DivisiItem                   `json:"Divisi"        validate:"required"`
	Tahun         string                       `json:"Tahun"         validate:"required"`
	Triwulan      string                       `json:"Triwulan"      validate:"required"`
	Kostl         string                       `json:"Kostl"`
	KostlTx       string                       `json:"KostlTx"`
	EntryUser     string                       `json:"EntryUser"`
	EntryName     string                       `json:"EntryName"`
	EntryTime     string                       `json:"EntryTime"`
	SaveAsDraft   string                       `json:"SaveAsDraft"   validate:"required"`
	Kpi           []PenyusunanKpiDetailRequest `json:"Kpi"           validate:"required,min=1,dive"`
	ChallengeList []PenyusunanChallenge        `json:"ChallengeList" validate:"required,min=1,dive"`
	MethodList    []PenyusunanMethod           `json:"MethodList"    validate:"required,min=1,dive"`
}

// CreatePenyusunanKpiRequest digunakan untuk endpoint POST /penyusunan-kpi/create.
type CreatePenyusunanKpiRequest struct {
	IdPengajuan  string     `json:"idPengajuan"  validate:"required"`
	ApprovalList []Approval `json:"ApprovalList" validate:"required,min=1,dive"`
	SaveAsDraft  string     `json:"SaveAsDraft"  validate:"required"`
}

type PenyusunanKpiDetailRequest struct {
	IdKpi      string `json:"idKpi"      validate:"required"`
	Kpi        string `json:"kpi"        validate:"required"`
	Rumus      string `json:"rumus"      validate:"required"`
	Persfektif string `json:"persfektif" validate:"required"`
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

type PenyusunanChallenge struct {
	IdDetailChallenge  string `json:"idDetailChallenge"  validate:"required"`
	Tahun              string `json:"tahun"              validate:"required"`
	Triwulan           string `json:"triwulan"           validate:"required"`
	NamaChallenge      string `json:"namaChallenge"      validate:"required"`
	DeskripsiChallenge string `json:"deskripsiChallenge" validate:"required"`
}

type PenyusunanMethod struct {
	IdDetailMethod  string `json:"idDetailMethod"  validate:"required"`
	Tahun           string `json:"tahun"           validate:"required"`
	Triwulan        string `json:"triwulan"        validate:"required"`
	NamaMethod      string `json:"namaMethod"      validate:"required"`
	DeskripsiMethod string `json:"deskripsiMethod" validate:"required"`
}

// =============================================================================
// EXCEL ROW DTO
// =============================================================================

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

// ValidatePenyusunanKpiResponse adalah response untuk endpoint /validate.
type ValidatePenyusunanKpiResponse struct {
	IDPengajuan   string                        `json:"idPengajuan"`
	Tahun         string                        `json:"tahun"`
	Triwulan      string                        `json:"triwulan"`
	Kostl         string                        `json:"kostl"`
	KostlTx       string                        `json:"kostlTx"`
	EntryUser     string                        `json:"entryUser"`
	EntryName     string                        `json:"entryName"`
	EntryTime     string                        `json:"entryTime"`
	SaveAsDraft   string                        `json:"saveAsDraft"`
	TotalKpi      int                           `json:"totalKpi"`
	Kpi           []PenyusunanKpiDetailResponse `json:"kpi"`
	ChallengeList []PenyusunanChallenge         `json:"challengeList"`
	MethodList    []PenyusunanMethod            `json:"methodList"`
}

// CreatePenyusunanKpiResponse adalah response untuk endpoint /create.
type CreatePenyusunanKpiResponse struct {
	IdPengajuan  string     `json:"idPengajuan"`
	SaveAsDraft  string     `json:"saveAsDraft"`
	ApprovalList []Approval `json:"approvalList"`
}

type PenyusunanKpiDetailResponse struct {
	IdKpi        string                           `json:"idKpi"`
	Kpi          string                           `json:"kpi"`
	Rumus        string                           `json:"rumus"`
	Persfektif   string                           `json:"persfektif"`
	KpiSubDetail []PenyusunanKpiSubDetailResponse `json:"kpiSubDetail"`
}

type PenyusunanKpiSubDetailResponse struct {
	IdDetail                  string  `json:"idDetail"`
	IdSubDetail               string  `json:"idSubDetail"`
	NamaKpi                   string  `json:"namaKpi"`
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
	Result                    *string `json:"result"`
	DeskripsiResult           *string `json:"deskripsiResult"`
	Process                   *string `json:"process"`
	DeskripsiProcess          *string `json:"deskripsiProcess"`
	Context                   *string `json:"context"`
	DeskripsiContext          *string `json:"deskripsiContext"`
}
