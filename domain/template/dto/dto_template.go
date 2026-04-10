package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

// FormatPenyusunanKpiRequest digunakan untuk endpoint POST /template/format-penyusunan-kpi.
// Menerima JSON body meskipun menggunakan method POST.
type FormatPenyusunanKpiRequest struct {
	Triwulan string `json:"triwulan" validate:"required,oneof=TW1 TW2 TW3 TW4"`
}

// TolakanPenyusunanKpiRequest digunakan untuk endpoint POST /template/tolakan-penyusunan-kpi.
// Menghasilkan file Excel yang sudah terisi data baris sub KPI berdasarkan id_pengajuan.
// Format kolom mengikuti triwulan dari DB: TW1/TW3 → A–O, TW2/TW4 → A–U.
type TolakanPenyusunanKpiRequest struct {
	IdPengajuan string `json:"id_pengajuan" validate:"required"`
}
