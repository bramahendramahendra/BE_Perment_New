package handler

import (
	service "permen_api/domain/master_status/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterStatusHandler struct {
	service service.MasterStatusServiceInterface
}

func NewMasterStatusHandler(service service.MasterStatusServiceInterface) *MasterStatusHandler {
	return &MasterStatusHandler{service: service}
}

func (h *MasterStatusHandler) GetAllMasterStatus(c *gin.Context) {
	data, err := h.service.GetAllMasterStatus()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Status berhasil diambil",
		Data:    data,
	})
}

func (h *MasterStatusHandler) GetDraftMasterStatus(c *gin.Context) {
	data, err := h.service.GetDraftMasterStatus()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Status Draft berhasil diambil",
		Data:    data,
	})
}
