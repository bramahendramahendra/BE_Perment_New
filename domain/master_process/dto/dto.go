package dto

// Search bersifat opsional, Triwulan dan Tahun wajib diisi.
type GetAllMasterProcessRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan" validate:"required"`
	Tahun    string `json:"tahun"     validate:"required,numeric,len=4"`
}

type MasterProcessResponse struct {
	IdProcess   int    `json:"id_process"`
	NamaProcess string `json:"nama_process"`
	DescProcess string `json:"desc_process"`
}
