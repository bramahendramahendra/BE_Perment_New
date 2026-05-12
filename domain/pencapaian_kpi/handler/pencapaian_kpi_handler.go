package handler

import (
	dto "permen_api/domain/pencapaian_kpi/dto"
	service "permen_api/domain/pencapaian_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	file_export "permen_api/pkg/file_export"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type PencapaianKpiHandler struct {
	service service.PencapaianKpiServiceInterface
}

func NewPencapaianKpiHandler(service service.PencapaianKpiServiceInterface) *PencapaianKpiHandler {
	return &PencapaianKpiHandler{service: service}
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllPencapaianKpi handles POST /pencapaian-kpi/get-all-pencapaian
func (h *PencapaianKpiHandler) GetAllPencapaianKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllPencapaianKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllPencapaianKpi(&req)
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
		Message:    "Data Pencapaian KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// =============================================================================
// GET DETAIL
// =============================================================================

// GetDetailPencapaianKpi handles POST /pencapaian-kpi/get-detail
func (h *PencapaianKpiHandler) GetDetailPencapaianKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetDetailPencapaianKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetDetailPencapaianKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data detail Pencapaian KPI berhasil diambil",
		Data:    data,
	})
}

// =============================================================================
// DOWNLOAD
// =============================================================================

// GetExcelPencapaianKpi handles POST /pencapaian-kpi/get-excel
func (h *PencapaianKpiHandler) GetExcelPencapaianKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetExcelPencapaianKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GetExcelPencapaianKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SendExcel(c, fileBytes, filename)
}

// GetPdfPencapaianKpi handles POST /pencapaian-kpi/get-pdf
func (h *PencapaianKpiHandler) GetPdfPencapaianKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetPdfPencapaianKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GetPdfPencapaianKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SendPDF(c, fileBytes, filename)
}
