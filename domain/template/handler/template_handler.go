package handler

import (
	"net/http"

	dto "permen_api/domain/template/dto"
	service "permen_api/domain/template/service"
	"permen_api/errors"
	binder "permen_api/pkg/binder"
	file_export "permen_api/pkg/file_export"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	service service.TemplateServiceInterface
}

func NewTemplateHandler(service service.TemplateServiceInterface) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// =============================================================================
// GET TEMPLATE PENYUSUNAN
// =============================================================================

// GetAllMasterProcess handles POST /template/format-penyusunan-kpi.
// Menerima application/json dengan JSON biasa.
func (h *TemplateHandler) GetFormatPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.FormatPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GenerateFormatPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SetExcelDownloadHeaders(c, filename)
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// GetRevisionPenyusunanKpi handles POST /template/revision-penyusunan-kpi.
// Menerima application/json dengan JSON biasa.
func (h *TemplateHandler) GetRevisionPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RevisionPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GenerateRevisionPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SetExcelDownloadHeaders(c, filename)
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// =============================================================================
// GET TEMPLATE REALISASI
// =============================================================================

// GetFormatRealisasiKpi handles POST /template/format-realisasi-kpi.
// Menerima application/json dengan JSON biasa.
func (h *TemplateHandler) GetFormatRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.FormatRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GenerateFormatRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SetExcelDownloadHeaders(c, filename)
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// GetRevisionRealisasiKpi handles POST /template/revision-realisasi-kpi.
// Menerima application/json dengan JSON biasa.
func (h *TemplateHandler) GetRevisionRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RevisionRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GenerateRevisionRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SetExcelDownloadHeaders(c, filename)
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}
