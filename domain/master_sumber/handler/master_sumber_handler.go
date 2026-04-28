package handler

import (
	service "permen_api/domain/master_sumber/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterSumberHandler struct {
	service service.MasterSumberServiceInterface
}

func NewMasterSumberHandler(service service.MasterSumberServiceInterface) *MasterSumberHandler {
	return &MasterSumberHandler{service: service}
}

func (h *MasterSumberHandler) GetAllMasterSumber(c *gin.Context) {
	data, err := h.service.GetAllMasterSumber()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Sumber berhasil diambil",
		Data:    data,
	})
}
