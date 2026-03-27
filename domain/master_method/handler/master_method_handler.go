package handler

import (
	dto "permen_api/domain/master_method/dto"
	service "permen_api/domain/master_method/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type MasterMethodHandler struct {
	service service.MasterMethodServiceInterface
}

func NewMasterMethodHandler(service service.MasterMethodServiceInterface) *MasterMethodHandler {
	return &MasterMethodHandler{service: service}
}

func (h *MasterMethodHandler) GetAllMasterMethod(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllMasterMethodRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetAllMasterMethod(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Method berhasil diambil",
		Data:    data,
	})
}
