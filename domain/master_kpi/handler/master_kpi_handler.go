package handler

import (
	service "permen_api/domain/master_kpi/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterKpiHandler struct {
	service service.MasterKpiServiceInterface
}

func NewMasterKpiHandler(service service.MasterKpiServiceInterface) *MasterKpiHandler {
	return &MasterKpiHandler{service: service}
}

func (h *MasterKpiHandler) GetAllMasterKpi(c *gin.Context) {
	data, err := h.service.GetAllMasterKpi()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil diambil",
		Data:    data,
	})
}
