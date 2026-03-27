package handler

import (
	service "permen_api/domain/master_perspektif/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterPerspektifHandler struct {
	service service.MasterPerspektifServiceInterface
}

func NewMasterPerspektifHandler(service service.MasterPerspektifServiceInterface) *MasterPerspektifHandler {
	return &MasterPerspektifHandler{service: service}
}

func (h *MasterPerspektifHandler) GetAllMasterPerspektif(c *gin.Context) {
	data, err := h.service.GetAllMasterPerspektif()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Perspektif berhasil diambil",
		Data:    data,
	})
}
