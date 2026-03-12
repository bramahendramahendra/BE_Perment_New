package dto

// =============================================================================
// REQUEST DTO
// =============================================================================

// FormatPenyusunanKpiRequest digunakan untuk endpoint GET /template/format-penyusunan-kpi.
// Menerima JSON body meskipun menggunakan method GET.
type FormatPenyusunanKpiRequest struct {
	Triwulan string `json:"triwulan" validate:"required,oneof=TW1 TW2 TW3 TW4"`
}
