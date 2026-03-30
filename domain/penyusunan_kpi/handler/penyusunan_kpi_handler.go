package handler

import (
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
	service "permen_api/domain/penyusunan_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type PenyusunanKpiHandler struct {
	service service.PenyusunanKpiServiceInterface
}

func NewPenyusunanKpiHandler(service service.PenyusunanKpiServiceInterface) *PenyusunanKpiHandler {
	return &PenyusunanKpiHandler{service: service}
}

// ValidatePenyusunanKpi handles POST /penyusunan-kpi/validate
// Menerima multipart/form-data dengan field REQUEST (JSON) dan files (Excel).
func (h *PenyusunanKpiHandler) ValidatePenyusunanKpi(c *gin.Context) {
	req, file, err := binder.BindMultipartJSON[dto.ValidatePenyusunanKpiRequest](c, "REQUEST", "files")
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userq := c.GetHeader("userq")
	if userq == "" {
		c.Error(&errors.BadRequestError{Message: "header 'userq' tidak ditemukan"})
		return
	}

	parts := strings.SplitN(userq, " | ", 2)
	if len(parts) != 2 {
		c.Error(&errors.BadRequestError{Message: "format header 'userq' tidak valid"})
		return
	}

	req.Kostl = req.Divisi.Kostl
	req.KostlTx = req.Divisi.KostlTx
	req.EntryUser = strings.TrimSpace(parts[0])
	req.EntryName = strings.TrimSpace(parts[1])
	req.EntryTime = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ValidatePenyusunanKpi(&req, file)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil disimpan",
		Data:    data,
	})
}

// CreatePenyusunanKpi handles POST /penyusunan-kpi/create
// Menerima JSON biasa (bukan multipart) dengan idPengajuan, ApprovalList, dan SaveAsDraft.
func (h *PenyusunanKpiHandler) CreatePenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreatePenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.CreatePenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil diajukan",
		Data:    data,
	})
}

// GetAllDraftPenyusunanKpi handles GET /penyusunan-kpi/get-all
// Mendukung filter opsional via query params: tahun, triwulan, kostl, status, page, limit.
func (h *PenyusunanKpiHandler) GetAllDraftPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindQuery[dto.GetAllDraftPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAllDraftPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	pagination := response_helper.SetPagination(&globalDTO.FilterRequestParams{
		Page:  req.Page,
		Limit: req.Limit,
	}, total)

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:       "00",
		Status:     true,
		Message:    "Data KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}
