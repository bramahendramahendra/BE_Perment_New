package handler

import (
	service "permen_api/domain/master_triwulan/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterTriwulanHandler struct {
	service service.MasterTriwulanServiceInterface
}

func NewMasterTriwulanHandler(service service.MasterTriwulanServiceInterface) *MasterTriwulanHandler {
	return &MasterTriwulanHandler{service: service}
}

func (h *MasterTriwulanHandler) GetAllMasterTriwulan(c *gin.Context) {
	data, err := h.service.GetAllMasterTriwulan()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Triwulan berhasil diambil",
		Data:    data,
	})
}
