package handler

import (
	dto "permen_api/domain/master_challenge/dto"
	service "permen_api/domain/master_challenge/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type MasterChallengeHandler struct {
	service service.MasterChallengeServiceInterface
}

func NewMasterChallengeHandler(service service.MasterChallengeServiceInterface) *MasterChallengeHandler {
	return &MasterChallengeHandler{service: service}
}

func (h *MasterChallengeHandler) GetAllMasterChallenge(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllMasterChallengeRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.GetAllMasterChallenge(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Challenge berhasil diambil",
		Data:    data,
	})
}
