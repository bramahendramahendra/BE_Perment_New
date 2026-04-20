package dto

// Search bersifat opsional, Triwulan dan Tahun wajib diisi.
type GetAllMasterResultRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan" validate:"required"`
	Tahun    string `json:"tahun"     validate:"required,numeric,len=4"`
}

type MasterResultResponse struct {
	IdResult   int    `json:"id_result"`
	NamaResult string `json:"nama_result"`
	DescResult string `json:"desc_result"`
}
