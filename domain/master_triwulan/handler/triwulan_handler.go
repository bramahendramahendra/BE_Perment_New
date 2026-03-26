package handler

import (
	service "permen_api/domain/master_triwulan/service"
	globalDTO "permen_api/dto"
	errors "permen_api/errors"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterTriwulanHandler struct {
	service service.MasterTriwulanServiceInterface
}

func NewMasterTriwulanHandler(service service.MasterTriwulanServiceInterface) *MasterTriwulanHandler {
	return &MasterTriwulanHandler{service: service}
}

func (h *MasterTriwulanHandler) GetAllTriwulan(c *gin.Context) {
	data, err := h.service.GetAllTriwulan()
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data master triwulan berhasil diambil",
		Data:    data,
	})
}
