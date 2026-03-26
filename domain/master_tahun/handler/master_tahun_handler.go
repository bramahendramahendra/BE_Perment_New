package handler

import (
	service "permen_api/domain/master_tahun/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterTahunHandler struct {
	service service.MasterTahunServiceInterface
}

func NewMasterTahunHandler(service service.MasterTahunServiceInterface) *MasterTahunHandler {
	return &MasterTahunHandler{service: service}
}

func (h *MasterTahunHandler) GetAllMasterTahun(c *gin.Context) {
	data, err := h.service.GetAllMasterTahun()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Tahun berhasil diambil",
		Data:    data,
	})
}
