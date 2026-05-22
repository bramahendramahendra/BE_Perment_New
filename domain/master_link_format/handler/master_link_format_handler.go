package handler

import (
	service "permen_api/domain/master_link_format/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterLinkFormatHandler struct {
	service service.MasterLinkFormatServiceInterface
}

func NewMasterLinkFormatHandler(service service.MasterLinkFormatServiceInterface) *MasterLinkFormatHandler {
	return &MasterLinkFormatHandler{service: service}
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterKpi handles POST /master-link_format/get-all.
// Menerima application/json dengan JSON biasa.
func (h *MasterLinkFormatHandler) GetAllMasterLinkFormat(c *gin.Context) {
	data, err := h.service.GetAllMasterLinkFormat()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data LinkFormat berhasil diambil",
		Data:    data,
	})
}
