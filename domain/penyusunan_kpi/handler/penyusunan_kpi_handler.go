package handler

import (
	"fmt"
	"net/http"
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
// Menerima JSON biasa (bukan multipart) dengan idPengajuan dan ApprovalList.
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

// GetAllApprovalPenyusunanKpi handles POST /penyusunan-kpi/get-all-approval
func (h *PenyusunanKpiHandler) GetAllApprovalPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllApprovalPenyusunanKpiRequest](c)
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

	req.ApprovalUser = strings.TrimSpace(parts[0])

	data, total, err := h.service.GetAllApprovalPenyusunanKpi(&req)
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

// GetAllTolakanPenyusunanKpi handles POST /penyusunan-kpi/get-all-tolakan
func (h *PenyusunanKpiHandler) GetAllTolakanPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllTolakanPenyusunanKpiRequest](c)
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

	data, total, err := h.service.GetAllTolakanPenyusunanKpi(&req)
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

// GetAllDaftarPenyusunanKpi handles POST /penyusunan-kpi/get-all-daftar-penyusunan
func (h *PenyusunanKpiHandler) GetAllDaftarPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarPenyusunanKpiRequest](c)
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

	data, total, err := h.service.GetAllDaftarPenyusunanKpi(&req)
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

// GetAllApprovalPenyusunanKpi handles POST /penyusunan-kpi/get-all-daftar-approval
func (h *PenyusunanKpiHandler) GetAllDaftarApprovalPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarApprovalPenyusunanKpiRequest](c)
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

	req.ApprovalUser = strings.TrimSpace(parts[0])

	data, total, err := h.service.GetAllDaftarApprovalPenyusunanKpi(&req)
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

// GetDetailPenyusunanKpi handles POST /penyusunan-kpi/get-detail
func (h *PenyusunanKpiHandler) GetDetailPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetDetailPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetDetailPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Detail KPI berhasil diambil",
		Data:    data,
	})
}

// GetCsvPenyusunanKpi handles POST /penyusunan-kpi/get-csv
func (h *PenyusunanKpiHandler) GetCsvPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetCsvPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GetCsvPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/csv", fileBytes)
}

// GetPdfPenyusunanKpi handles POST /penyusunan-kpi/get-pdf
func (h *PenyusunanKpiHandler) GetPdfPenyusunanKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetPdfPenyusunanKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	fileBytes, filename, err := h.service.GetPdfPenyusunanKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "application/pdf", fileBytes)
}
