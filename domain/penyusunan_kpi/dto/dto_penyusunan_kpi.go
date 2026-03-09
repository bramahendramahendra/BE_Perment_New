package dto

// =============================================
// REQUEST DTO
// =============================================

type DivisiItem struct {
	Kostl   string `json:"Kostl"   validate:"required"`
	KostlTx string `json:"KostlTx" validate:"required"`
}

type InsertPenyusunanKpiRequest struct {
	Divisi         DivisiItem                `json:"Divisi"         validate:"required"`
	Tahun          string                    `json:"Tahun"          validate:"required"`
	Triwulan       string                    `json:"Triwulan"       validate:"required"`
	Kostl          string                    `json:"Kostl"`
	KostlTx        string                    `json:"KostlTx"`
	EntryUser      string                    `json:"EntryUser"`
	EntryName      string                    `json:"EntryName"`
	EntryTime      string                    `json:"EntryTime"`
	ApprovalPosisi string                    `json:"ApprovalPosisi" validate:"required"`
	ApprovalList   string                    `json:"ApprovalList"   validate:"required"`
	SaveAsDraft    string                    `json:"SaveAsDraft"    validate:"required"`
	Kpi            []PenyusunanKpiDetailItem `json:"Kpi"            validate:"required,min=1,dive"`
	ChallengeList  []PenyusunanChallengeItem `json:"ChallengeList"  validate:"required,min=1,dive"`
	MethodList     []PenyusunanMethodItem    `json:"MethodList"     validate:"required,min=1,dive"`
}

type PenyusunanKpiDetailItem struct {
	IdKpi      string `json:"idKpi"      validate:"required"`
	Kpi        string `json:"kpi"        validate:"required"`
	Rumus      string `json:"rumus"      validate:"required"`
	Persfektif string `json:"persfektif" validate:"required"`
}

type PenyusunanChallengeItem struct {
	IdDetailChallenge  string `json:"idDetailChallenge"  validate:"required"`
	Tahun              string `json:"tahun"              validate:"required"`
	Triwulan           string `json:"triwulan"           validate:"required"`
	NamaChallenge      string `json:"namaChallenge"      validate:"required"`
	DeskripsiChallenge string `json:"deskripsiChallenge" validate:"required"`
}

type PenyusunanMethodItem struct {
	IdDetailMethod  string `json:"idDetailMethod"  validate:"required"`
	Tahun           string `json:"tahun"           validate:"required"`
	Triwulan        string `json:"triwulan"        validate:"required"`
	NamaMethod      string `json:"namaMethod"      validate:"required"`
	DeskripsiMethod string `json:"deskripsiMethod" validate:"required"`
}

// =============================================
// EXCEL ROW DTO
// =============================================

type PenyusunanKpiSubDetailRow struct {
	No           int
	KPI          string
	SubKPI       string
	IdSubKpi     string
	Otomatis     string
	Polarisasi   string
	IdPolarisasi string

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

// =============================================
// RESPONSE DTO
// =============================================

type InsertPenyusunanKpiResponse struct {
	IDPengajuan    string                    `json:"idPengajuan"`
	Tahun          string                    `json:"tahun"`
	Triwulan       string                    `json:"triwulan"`
	Kostl          string                    `json:"kostl"`
	KostlTx        string                    `json:"kostlTx"`
	EntryUser      string                    `json:"entryUser"`
	EntryName      string                    `json:"entryName"`
	EntryTime      string                    `json:"entryTime"`
	ApprovalPosisi string                    `json:"approvalPosisi"`
	SaveAsDraft    string                    `json:"saveAsDraft"`
	TotalKpi       int                       `json:"totalKpi"`
	Kpi            []PenyusunanKpiDetailItem `json:"kpi"`
	ChallengeList  []PenyusunanChallengeItem `json:"challengeList"`
	MethodList     []PenyusunanMethodItem    `json:"methodList"`
}
