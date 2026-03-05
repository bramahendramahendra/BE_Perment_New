package dto

// =============================================
// REQUEST DTO
// =============================================

// InsertPenyusunanKpiRequest adalah request utama yang dikirim via multipart form field "REQUEST"
type InsertPenyusunanKpiRequest struct {
	IDPengajuan    string                    `json:"IDPengajuan"    validate:"required"`
	Divisi         string                    `json:"Divisi"         validate:"required"`
	Tahun          string                    `json:"Tahun"          validate:"required"`
	Triwulan       string                    `json:"Triwulan"       validate:"required"`
	Kostl          string                    `json:"Kostl"          validate:"required"`
	KostlTx        string                    `json:"KostlTx"        validate:"required"`
	EntryUser      string                    `json:"EntryUser"      validate:"required"`
	EntryName      string                    `json:"EntryName"      validate:"required"`
	EntryTime      string                    `json:"EntryTime"      validate:"required"`
	ApprovalPosisi string                    `json:"ApprovalPosisi" validate:"required"`
	ApprovalList   string                    `json:"ApprovalList"   validate:"required"`
	SaveAsDraft    string                    `json:"SaveAsDraft"    validate:"required"`
	Kpi            []PenyusunanKpiDetailItem `json:"Kpi"            validate:"required,min=1,dive"`
	ChallengeList  []PenyusunanChallengeItem `json:"ChallengeList"  validate:"required,min=1,dive"`
	MethodList     []PenyusunanMethodItem    `json:"MethodList"     validate:"required,min=1,dive"`
}

// PenyusunanKpiDetailItem adalah metadata satu KPI pada tabel data_kpi_detail.
// Data sub detail (data_kpi_subdetail) berasal dari file Excel yang di-upload,
// dengan urutan file Excel sesuai urutan array Kpi ini.
type PenyusunanKpiDetailItem struct {
	IdKpi             string `json:"idKpi"             validate:"required"`
	Kpi               string `json:"kpi"               validate:"required"`
	Rumus             string `json:"rumus"             validate:"required"`
	Persfektif        string `json:"persfektif"        validate:"required"`
	KeteranganProject string `json:"keteranganProject"`
}

// PenyusunanChallengeItem adalah satu baris data challenge pada tabel data_challenge_detail
type PenyusunanChallengeItem struct {
	IdPengajuan        string `json:"idPengajuan"       validate:"required"`
	IdDetailChallenge  string `json:"idDetailChallenge" validate:"required"`
	Tahun              string `json:"tahun"             validate:"required"`
	Triwulan           string `json:"triwulan"          validate:"required"`
	NamaChallenge      string `json:"namaChallenge"     validate:"required"`
	DeskripsiChallenge string `json:"deskripsiChallenge" validate:"required"`
}

// PenyusunanMethodItem adalah satu baris data method pada tabel data_method_detail
type PenyusunanMethodItem struct {
	IdPengajuan     string `json:"idPengajuan"  validate:"required"`
	IdDetailMethod  string `json:"idDetailMethod" validate:"required"`
	Tahun           string `json:"tahun"        validate:"required"`
	Triwulan        string `json:"triwulan"     validate:"required"`
	NamaMethod      string `json:"namaMethod"   validate:"required"`
	DeskripsiMethod string `json:"deskripsiMethod" validate:"required"`
}

// =============================================
// EXCEL PARSED ROW
// =============================================

// PenyusunanKpiSubDetailRow adalah representasi satu baris data dari file Excel
// yang akan di-insert ke tabel data_kpi_subdetail.
// Header Excel berada di baris 2, data dimulai dari baris 3.
//
// Mapping kolom Excel:
//   Col A  = No                        (angka)
//   Col B  = KPI                       (free text, tidak boleh blank)
//   Col C  = Sub KPI                   (free text, tidak boleh blank)
//   Col D  = Polarisasi                (enum: Maximize / Minimize)
//   Col E  = Capping                   (enum: 100% / 110%)
//   Col F  = Bobot %                   (angka 2 desimal, total semua baris = 100)
//   Col G  = Glossary                  (free text, tidak boleh blank)
//   Col H  = Target Triwulanan         (free text, tidak boleh blank)
//   Col I  = Target Kuantitatif Triwulanan (angka 2 desimal)
//   Col J  = Target Tahunan            (free text, tidak boleh blank)
//   Col K  = Target Kuantitatif Tahunan (angka 2 desimal)
//   Col L  = Terdapat Qualifier        (enum: Ya / Tidak)
//   Col M  = Qualifier                 (free text, wajib jika Col L = "Ya")
//   Col N  = Deskripsi Qualifier       (free text, wajib jika Col L = "Ya")
//   Col O  = Target Qualifier          (free text, wajib jika Col L = "Ya")
//   Col P  = Result                    (free text, tidak boleh blank)
//   Col Q  = Deskripsi Result          (free text, tidak boleh blank)
//   Col R  = Process                   (free text, tidak boleh blank)
//   Col S  = Deskripsi Process         (free text, tidak boleh blank)
//   Col T  = Context                   (free text, tidak boleh blank)
//   Col U  = Deskripsi Context         (free text, tidak boleh blank)
type PenyusunanKpiSubDetailRow struct {
	No                        int
	KPI                       string
	SubKPI                    string
	Polarisasi                string
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
	Result                    string
	DeskripsiResult           string
	Process                   string
	DeskripsiProcess          string
	Context                   string
	DeskripsiContext          string
}

// =============================================
// RESPONSE DTO
// =============================================

// InsertPenyusunanKpiResponse adalah response yang dikembalikan jika insert berhasil
type InsertPenyusunanKpiResponse struct {
	IDPengajuan string `json:"idPengajuan"`
	Message     string `json:"message"`
}
