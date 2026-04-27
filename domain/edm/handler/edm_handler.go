package handler

import (
	dto "permen_api/domain/edm/dto"
	service "permen_api/domain/edm/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type EdmHandler struct {
	service service.EdmServiceInterface
}

func NewEdmHandler(service service.EdmServiceInterface) *EdmHandler {
	return &EdmHandler{service: service}
}

func (h *EdmHandler) GetRealisasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetRealisasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetRealisasi(&req)
	if err != nil {
		c.Error(&errors.InternalServerError{Message: err.Error()})
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Berhasil mendapatkan data realisasi EDM",
		Data:    data,
	})
}
