package handler

import (
	dto "permen_api/domain/master_result/dto"
	service "permen_api/domain/master_result/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type MasterResultHandler struct {
	service service.MasterResultServiceInterface
}

func NewMasterResultHandler(service service.MasterResultServiceInterface) *MasterResultHandler {
	return &MasterResultHandler{service: service}
}

func (h *MasterResultHandler) GetAllMasterResult(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllMasterResultRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetAllMasterResult(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Result berhasil diambil",
		Data:    data,
	})
}
