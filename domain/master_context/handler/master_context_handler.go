package handler

import (
	dto "permen_api/domain/master_context/dto"
	service "permen_api/domain/master_context/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type MasterContextHandler struct {
	service service.MasterContextServiceInterface
}

func NewMasterContextHandler(service service.MasterContextServiceInterface) *MasterContextHandler {
	return &MasterContextHandler{service: service}
}

func (h *MasterContextHandler) GetAllMasterContext(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllMasterContextRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetAllMasterContext(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Context berhasil diambil",
		Data:    data,
	})
}
