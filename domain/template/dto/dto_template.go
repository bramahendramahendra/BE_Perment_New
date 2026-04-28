package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

// DivisiRequest berisi informasi divisi (kostl dan nama kostl).
type DivisiRequest struct {
	Kostl   string `json:"kostl"    validate:"required"`
	KostlTx string `json:"kostl_tx" validate:"required"`
}

// FormatPenyusunanKpiRequest digunakan untuk endpoint POST /template/format-penyusunan-kpi.
// Menerima JSON body meskipun menggunakan method POST. Tanpa id_pengajuan.
type FormatPenyusunanKpiRequest struct {
	Divisi   DivisiRequest `json:"divisi"    validate:"required"`
	Tahun    string        `json:"tahun"     validate:"required"`
	Triwulan string        `json:"triwulan"  validate:"required,oneof=TW1 TW2 TW3 TW4"`
}

// RevisionPenyusunanKpiRequest digunakan untuk endpoint POST /template/revision-penyusunan-kpi.
// Menghasilkan file Excel yang sudah terisi data baris sub KPI berdasarkan id_pengajuan.
// Format kolom mengikuti triwulan dari DB: TW1/TW3 → A–O, TW2/TW4 → A–U.
type RevisionPenyusunanKpiRequest struct {
	IdPengajuan string        `json:"id_pengajuan" validate:"required"`
	Divisi      DivisiRequest `json:"divisi"       validate:"required"`
	Tahun       string        `json:"tahun"        validate:"required"`
	Triwulan    string        `json:"triwulan"     validate:"required,oneof=TW1 TW2 TW3 TW4"`
}

// RevisionRealisasiKpiRequest digunakan untuk endpoint POST /template/revision-realisasi-kpi.
// Menghasilkan file Excel realisasi KPI yang sudah terisi data baris realisasi berdasarkan id_pengajuan,
// sehingga user dapat langsung merevisi dan mengupload ulang via /realisasi-kpi/revision.
// Format kolom mengikuti triwulan dari request:
//
//	TW1/TW3 → A–M (kolom A–I dari DB, J–M pre-filled data realisasi).
//	TW2/TW4 → A–Y (kolom extended N,O,R,S,V,W dari DB; P,Q,T,U,X,Y pre-filled realisasi).
type RevisionRealisasiKpiRequest struct {
	IdPengajuan string        `json:"id_pengajuan" validate:"required"`
	Divisi      DivisiRequest `json:"divisi"       validate:"required"`
	Tahun       string        `json:"tahun"        validate:"required"`
	Triwulan    string        `json:"triwulan"     validate:"required,oneof=TW1 TW2 TW3 TW4"`
}

// FormatRealisasiKpiRequest digunakan untuk endpoint POST /template/format-realisasi-kpi.
// Menghasilkan file Excel template realisasi KPI yang sudah terisi data A–I (dari DB),
// dengan kolom J–M dikosongkan untuk diisi user.
// Format kolom mengikuti triwulan dari request:
//
//	TW1/TW3 → A–S (kolom N–S terisi data result/process/context dari DB)
//	TW2/TW4 → A–Y (kolom N, O, R, S, V, W dari DB; kolom P, Q, T, U, X, Y kosong untuk user)
type FormatRealisasiKpiRequest struct {
	IdPengajuan string        `json:"id_pengajuan" validate:"required"`
	Divisi      DivisiRequest `json:"divisi"       validate:"required"`
	Tahun       string        `json:"tahun"        validate:"required"`
	Triwulan    string        `json:"triwulan"     validate:"required,oneof=TW1 TW2 TW3 TW4"`
}
