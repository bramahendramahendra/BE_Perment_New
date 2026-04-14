package handler

import (
	dto "permen_api/domain/master_process/dto"
	service "permen_api/domain/master_process/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type MasterProcessHandler struct {
	service service.MasterProcessServiceInterface
}

func NewMasterProcessHandler(service service.MasterProcessServiceInterface) *MasterProcessHandler {
	return &MasterProcessHandler{service: service}
}

func (h *MasterProcessHandler) GetAllMasterProcess(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllMasterProcessRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetAllMasterProcess(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Process berhasil diambil",
		Data:    data,
	})
}
