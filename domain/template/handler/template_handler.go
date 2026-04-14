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

// GetFormatPenyusunanKpi handles POST /template/format-penyusunan-kpi
// Menerima JSON body dengan field triwulan (TW1/TW2/TW3/TW4).
// Menghasilkan file Excel template kosong yang langsung diunduh oleh client.
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
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// GetTolakanPenyusunanKpi handles POST /template/tolakan-penyusunan-kpi
// Menerima JSON body dengan field id_pengajuan.
// Menghasilkan file Excel yang sudah terisi data baris sub KPI berdasarkan id_pengajuan,
// sehingga user dapat langsung merevisi dan mengupload ulang via /penyusunan-kpi/revision.
func (h *TemplateHandler) GetTolakanPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.TolakanPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GenerateTolakanPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	file_export.SetExcelDownloadHeaders(c, filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

// GetFormatRealisasiKpi handles POST /template/format-realisasi-kpi
// Menerima JSON body dengan field id_pengajuan dan triwulan (TW1/TW2/TW3/TW4).
// Menghasilkan file Excel template realisasi KPI: kolom A–I terisi data dari DB,
// kolom J–M dikosongkan untuk diisi user.
// Format kolom extended mengikuti triwulan dari request:
//
//	TW1/TW3 → kolom N–S terisi data result/process/context dari DB.
//	TW2/TW4 → kolom N–Y dengan sebagian dikosongkan untuk diisi user.
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
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

