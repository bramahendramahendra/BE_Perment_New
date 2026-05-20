package dto

type GetKpiRequest struct {
	Periode string `json:"periode" validate:"required"`
}
