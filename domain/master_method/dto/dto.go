package dto

// Search bersifat opsional, Triwulan dan Tahun wajib diisi.
type GetAllMasterMethodRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan" validate:"required"`
	Tahun    string `json:"tahun"     validate:"required,numeric,len=4"`
}

type MasterMethodResponse struct {
	IdMethodUse int    `json:"idMethodUse"`
	NamaMethod  string `json:"namaMethod"`
	DescMethod  string `json:"descMethod"`
}
