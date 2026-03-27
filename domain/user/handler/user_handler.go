package handler

import (
	dto "permen_api/domain/user/dto"
	service "permen_api/domain/user/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUser menangani request POST /api/user/get-all.
// Mengambil semua data user, filter Branch bersifat opsional.
func (h *UserHandler) GetAllUser(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllUserRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, err := h.service.GetAllUser(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data user berhasil diambil",
		Data:    data,
	})
}
