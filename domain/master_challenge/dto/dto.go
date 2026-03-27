package dto

// GetAllMasterChallengeRequest adalah request body untuk endpoint GET-ALL master challenge.
// Search bersifat opsional, Triwulan dan Tahun wajib diisi.
type GetAllMasterChallengeRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan" validate:"required"`
	Tahun    string `json:"tahun"     validate:"required,numeric,len=4"`
}

// MasterChallengeResponse adalah response data per item challenge.
type MasterChallengeResponse struct {
	IdChallenge   string `json:"idChallenge"`
	NamaChallenge string `json:"namaChallenge"`
	DescChallenge string `json:"descChallenge"`
}
