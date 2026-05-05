package handler

import (
	service "permen_api/domain/master_divisi/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type MasterDivisiHandler struct {
	service service.MasterDivisiServiceInterface
}

func NewMasterDivisiHandler(service service.MasterDivisiServiceInterface) *MasterDivisiHandler {
	return &MasterDivisiHandler{service: service}
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllMasterDivisi handles POST /master-divisi/get-all
// Menerima application/json dengan JSON biasa.
func (h *MasterDivisiHandler) GetAllMasterDivisi(c *gin.Context) {
	data, err := h.service.GetAllMasterDivisi()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Divisi berhasil diambil",
		Data:    data,
	})
}
