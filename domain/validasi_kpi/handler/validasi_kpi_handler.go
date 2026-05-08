package handler

import (
	"strings"
	"time"

	dto "permen_api/domain/validasi_kpi/dto"
	service "permen_api/domain/validasi_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type ValidasiKpiHandler struct {
	service service.ValidasiKpiServiceInterface
}

func NewValidasiKpiHandler(service service.ValidasiKpiServiceInterface) *ValidasiKpiHandler {
	return &ValidasiKpiHandler{service: service}
}

// =============================================================================
// INPUT
// =============================================================================

// InputValidasiKpi handles POST /validasi-kpi/input
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) InputValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.InputValidasiKpiRequest](c)
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
	req.EntryUserValidasi = strings.TrimSpace(parts[0])
	req.EntryNameValidasi = strings.TrimSpace(parts[1])
	req.EntryTimeValidasi = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.InputValidasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data validasi KPI berhasil disimpan",
		Data:    data,
	})
}

// =============================================================================
// APPROVAL
// =============================================================================

// ApproveValidasiKpi handles POST /validasi-kpi/approve
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) ApproveValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.ApproveValidasiKpiRequest](c)
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

	req.ApprovalUserValidasi = strings.TrimSpace(parts[0])
	req.ApprovalNameValidasi = strings.TrimSpace(parts[1])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ApproveValidasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Validasi KPI berhasil diapprove",
		Data:    data,
	})
}

// RejectValidasiKpi handles POST /validasi-kpi/reject
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) RejectValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RejectValidasiKpiRequest](c)
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
	req.ApprovalUserValidasi = strings.TrimSpace(parts[0])
	req.ApprovalNameValidasi = strings.TrimSpace(parts[1])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.RejectValidasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Validasi KPI berhasil ditolak",
		Data:    data,
	})
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllValidasiKpi handles POST /validasi-kpi/get-all-validasi
// Menerima application/json.
func (h *ValidasiKpiHandler) GetAllValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllValidasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllValidasiKpi(&req)
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
		Message:    "Data Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllApprovalValidasi handles POST /validasi-kpi/get-all-approval
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) GetAllApprovalValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllApprovalValidasiKpiRequest](c)
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

	req.ApprovalUserValidasi = strings.TrimSpace(parts[0])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllApprovalValidasiKpi(&req)
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
		Message:    "Data Approval Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllTolakanValidasiKpi handles POST /validasi-kpi/get-all-tolakan
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) GetAllTolakanValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllTolakanValidasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllTolakanValidasiKpi(&req)
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
		Message:    "Data Penolakan Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllDaftarValidasiKpi handles POST /validasi-kpi/get-all-daftar-penyusunan
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) GetAllDaftarValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarPValidasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllDaftarValidasiKpi(&req)
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
		Message:    "Data Daftar Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllDaftarApprovalValidasi handles POST /validasi-kpi/get-all-daftar-approval
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) GetAllDaftarApprovalValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarApprovalValidasiKpiRequest](c)
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

	req.ApprovalUserValidasi = strings.TrimSpace(parts[0])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllDaftarApprovalValidasiKpi(&req)
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
		Message:    "Data Daftar Approval Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// =============================================================================
// GET DETAIL VALIDASI
// =============================================================================

// GetDetailValidasiKpi handles POST /validasi-kpi/get-detail
// Menerima application/json dengan JSON biasa.
func (h *ValidasiKpiHandler) GetDetailValidasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetDetailValidasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetDetailValidasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data detail Validasi KPI berhasil diambil",
		Data:    data,
	})
}
