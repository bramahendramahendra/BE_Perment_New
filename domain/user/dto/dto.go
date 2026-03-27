package dto

type (
	// GetAllUserRequest adalah request body untuk endpoint POST /user/get-all.
	// Branch bersifat opsional, jika diisi akan digunakan sebagai filter data user.
	GetAllUserRequest struct {
		Branch string `json:"branch"`
	}

	// UserResponse adalah representasi data user yang dikembalikan ke client.
	UserResponse struct {
		Pernr string `json:"pernr"`
		Sname string `json:"sname"`
	}
)
