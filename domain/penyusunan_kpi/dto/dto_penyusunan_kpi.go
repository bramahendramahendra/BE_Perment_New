package dto

// =============================================
// REQUEST DTO
// =============================================

// InsertPenyusunanKpiRequest adalah request utama yang dikirim via multipart form field "REQUEST"
//
// Catatan generate ID (dilakukan di backend, tidak perlu dikirim frontend):
//   - IDPengajuan = Kostl + Tahun + Triwulan + timestamp (ymdhis)
//                   contoh: "PS100012026TW2260304040242"
//   - Id per KPI  = IDPengajuan + "P" + index 3 digit
//                   contoh: "PS100012026TW2260304040242P001"
type InsertPenyusunanKpiRequest struct {
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
	Kpi            []PenyusunanKpiDetailItem `json:"Kpi"           validate:"required,min=1,dive"`
	ChallengeList  []PenyusunanChallengeItem `json:"ChallengeList" validate:"required,min=1,dive"`
	MethodList     []PenyusunanMethodItem    `json:"MethodList"    validate:"required,min=1,dive"`
}

// PenyusunanKpiDetailItem adalah metadata satu KPI pada tabel data_kpi_detail.
// Data sub detail (data_kpi_subdetail) berasal dari file Excel yang di-upload.
//
// Catatan:
//   - Id dan IDPengajuan tidak perlu dikirim — di-generate otomatis oleh backend
//   - KeteranganProject tidak perlu dikirim — backend otomatis mengisi dengan "-"
type PenyusunanKpiDetailItem struct {
	IdKpi      string `json:"idKpi"      validate:"required"`
	Kpi        string `json:"kpi"        validate:"required"`
	Rumus      string `json:"rumus"      validate:"required"`
	Persfektif string `json:"persfektif" validate:"required"`
}

// PenyusunanChallengeItem adalah satu baris data challenge pada tabel data_challenge_detail.
// Jika tidak ada data challenge (non-TW4), frontend mengirim nilai "-" pada semua field.
type PenyusunanChallengeItem struct {
	IdDetailChallenge  string `json:"idDetailChallenge"  validate:"required"`
	Tahun              string `json:"tahun"              validate:"required"`
	Triwulan           string `json:"triwulan"           validate:"required"`
	NamaChallenge      string `json:"namaChallenge"      validate:"required"`
	DeskripsiChallenge string `json:"deskripsiChallenge" validate:"required"`
}

// PenyusunanMethodItem adalah satu baris data method pada tabel data_method_detail.
// Jika tidak ada data method (non-TW4), frontend mengirim nilai "-" pada semua field.
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

// PenyusunanKpiSubDetailRow adalah representasi satu baris data sub KPI dari file Excel.
//
// Sheet & kolom:
//	Sheet "TW 4"        → kolom A (No.) sampai U (Deskripsi Context)  — 21 kolom
//	Sheet "Selain TW 4" → kolom A (No.) sampai O (Target Qualifier)   — 15 kolom
//
// Mapping kolom:
//   Col A  = No                            (angka)
//   Col B  = KPI (text)                    (free text, kunci mapping ke req.Kpi)
//   Col C  = Sub KPI (text)                (free text, tidak boleh blank)
//   Col D  = Polarisasi                    (enum: Maximize / Minimize)
//   Col E  = Capping                       (enum: 100% / 110%)
//   Col F  = Bobot %                       (angka 2 desimal, total per KPI = 100%)
//   Col G  = Glossary                      (free text, tidak boleh blank)
//   Col H  = Target Triwulanan             (free text, tidak boleh blank)
//   Col I  = Target Kuantitatif Triwulanan (angka 2 desimal)
//   Col J  = Target Tahunan                (free text, tidak boleh blank)
//   Col K  = Target Kuantitatif Tahunan    (angka 2 desimal)
//   Col L  = Terdapat Qualifier            (enum: Ya / Tidak)
//   Col M  = Qualifier                     (free text, wajib jika Col L = "Ya")
//   Col N  = Deskripsi Qualifier           (free text, wajib jika Col L = "Ya")
//   Col O  = Target Qualifier              (free text, wajib jika Col L = "Ya")
//   Col P  = Result                        (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//   Col Q  = Deskripsi Result              (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//   Col R  = Process                       (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//   Col S  = Deskripsi Process             (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//   Col T  = Context                       (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//   Col U  = Deskripsi Context             (free text, hanya sheet "TW 4", NULL jika "Selain TW 4")
//

//
// Lookup yang dilakukan backend (setelah parse Excel, sebelum insert DB):
//
//	IdSubKpi     → lookup mst_kpi WHERE LOWER(kpi) = LOWER(SubKPI)
//	               Ditemukan  : ambil id_kpi dari DB (SubKPI juga diupdate ke nama dari DB)
//	               Tidak ditemukan : IdSubKpi = "0", SubKPI tetap dari Excel
//
//	IdPolarisasi → lookup mst_polarisasi WHERE LOWER(polarisasi) = LOWER(Polarisasi)
//	               Maximize = "1", Minimize = "0"
//
//	Validasi rumus (hanya jika IdSubKpi != "0"):
//	               IdPolarisasi harus == rumus dari mst_kpi
//	               Jika tidak cocok → return error
//
// Kolom P–U (Result–DeskripsiContext):
//   - Sheet "TW 4"        → berisi nilai dari Excel (*string != nil)
//   - Sheet "Selain TW 4" → nil (disimpan NULL di DB)
type PenyusunanKpiSubDetailRow struct {
	No           int
	KPI          string // kolom B — nama KPI induk (untuk mapping, tidak disimpan ke subdetail)
	SubKPI       string // kolom C — nama sub KPI; diupdate ke nama dari DB jika ditemukan di mst_kpi
	IdSubKpi     string // hasil lookup mst_kpi; "0" jika tidak ditemukan
	Otomatis     string // "1" jika IdSubKpi ditemukan di mst_kpi, "0" jika tidak
	Polarisasi   string // kolom D — teks "Maximize" atau "Minimize" dari Excel
	IdPolarisasi string // hasil lookup mst_polarisasi; "1"=Maximize, "0"=Minimize
	// Catatan: IdPolarisasi inilah yang disimpan ke kolom `rumus` di data_kpi_subdetail

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
	IsTW4                     bool    // true jika berasal dari sheet "TW 4"
	Result                    *string // nil jika sheet "Selain TW 4" → NULL di DB
	DeskripsiResult           *string // nil jika sheet "Selain TW 4" → NULL di DB
	Process                   *string // nil jika sheet "Selain TW 4" → NULL di DB
	DeskripsiProcess          *string // nil jika sheet "Selain TW 4" → NULL di DB
	Context                   *string // nil jika sheet "Selain TW 4" → NULL di DB
	DeskripsiContext          *string // nil jika sheet "Selain TW 4" → NULL di DB
}

// =============================================
// RESPONSE DTO
// =============================================

// InsertPenyusunanKpiResponse adalah response yang dikembalikan jika insert berhasil
type InsertPenyusunanKpiResponse struct {
	IDPengajuan string `json:"idPengajuan"`
	Message     string `json:"message"`
}
