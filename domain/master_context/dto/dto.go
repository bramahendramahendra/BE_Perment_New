package dto

// Search bersifat opsional, Triwulan dan Tahun wajib diisi.
type GetAllMasterContextRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan" validate:"required"`
	Tahun    string `json:"tahun"     validate:"required,numeric,len=4"`
}

type MasterContextResponse struct {
	IdContext   int    `json:"id_Context"`
	NamaContext string `json:"nama_context"`
	DescContext string `json:"desc_context"`
}
