package dto

// GetAllMasterChallengeRequest adalah request body untuk endpoint GET-ALL master challenge.
// Semua field bersifat opsional — digunakan sebagai filter.
type GetAllMasterChallengeRequest struct {
	Search   string `json:"search"`
	Triwulan string `json:"triwulan"`
	Tahun    string `json:"tahun"`
}

// MasterChallengeResponse adalah response data per item challenge.
type MasterChallengeResponse struct {
	IdChallenge   string `json:"idChallenge"`
	NamaChallenge string `json:"namaChallenge"`
	DescChallenge string `json:"descChallenge"`
	Tahun         string `json:"tahun"`
	Triwulan      string `json:"triwulan"`
	EntryUser     string `json:"entryUser"`
	EntryName     string `json:"entryName"`
	EntryTime     string `json:"entryTime"`
}
