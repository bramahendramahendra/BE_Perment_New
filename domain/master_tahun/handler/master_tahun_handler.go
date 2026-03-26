package handler

import (
	service "permen_api/domain/sample/service"
	globalDTO "permen_api/dto"
	response_helper "permen_api/helper/response"

	"github.com/gin-gonic/gin"
)

type SampleHandler struct {
	service service.UserIntegrationServiceInterface
}

func NewSampleHandler(service service.UserIntegrationServiceInterface) *SampleHandler {
	return &SampleHandler{service: service}
}

func (h *SampleHandler) GetAllUserIntegrations(c *gin.Context) {
	data, err := h.service.GetAllUserIntegrations()
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "All user integrations retrieved successfully",
		Data:    data,
	})
}
