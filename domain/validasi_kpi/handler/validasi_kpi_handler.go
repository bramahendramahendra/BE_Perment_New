package handler

import (
	"strings"

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

// parseUserqHeader mengekstrak userid dan nama dari header "userq" (format: "userid | nama").
func parseUserqHeader(c *gin.Context) (userid, nama string, ok bool) {
	userq := c.GetHeader("userq")
	if userq == "" {
		c.Error(&errors.BadRequestError{Message: "header 'userq' tidak ditemukan"})
		return "", "", false
	}
	parts := strings.SplitN(userq, " | ", 2)
	if len(parts) != 2 {
		c.Error(&errors.BadRequestError{Message: "format header 'userq' tidak valid"})
		return "", "", false
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
}

// =============================================================================
// INPUT VALIDASI
// =============================================================================

// InputValidasi handles POST /validasi-kpi/input
// Menerima application/json. Menyimpan data validasi KPI dan mengirim notifikasi ke approver (status → 6).
func (h *ValidasiKpiHandler) InputValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.InputValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, nama, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.EntryUserValidasi = userid
	req.EntryNameValidasi = nama

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.InputValidasi(&req)
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
// APPROVAL VALIDASI (legacy — approve atau reject dalam satu endpoint)
// =============================================================================

// ApprovalValidasi handles POST /validasi-kpi/approval
// Menerima application/json. Memproses approve atau reject validasi KPI.
func (h *ValidasiKpiHandler) ApprovalValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.ApprovalValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, nama, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.ApprovalUser = userid
	req.ApprovalName = nama

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ApprovalValidasi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	message := "Validasi KPI berhasil diapprove"
	if req.Status == "reject" {
		message = "Validasi KPI berhasil ditolak"
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// =============================================================================
// APPROVE VALIDASI
// =============================================================================

// ApproveValidasi handles POST /validasi-kpi/approve
func (h *ValidasiKpiHandler) ApproveValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.ApproveValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, nama, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.ApprovalUser = userid
	req.ApprovalName = nama

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ApproveValidasi(&req)
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

// =============================================================================
// REJECT VALIDASI
// =============================================================================

// RejectValidasi handles POST /validasi-kpi/reject
func (h *ValidasiKpiHandler) RejectValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RejectValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, nama, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.ApprovalUser = userid
	req.ApprovalName = nama

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.RejectValidasi(&req)
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
// GET ALL APPROVAL VALIDASI
// =============================================================================

// GetAllApprovalValidasi handles POST /validasi-kpi/get-all-approval
func (h *ValidasiKpiHandler) GetAllApprovalValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllApprovalValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, _, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.ApprovalUser = userid

	data, total, err := h.service.GetAllApprovalValidasi(&req)
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

// =============================================================================
// GET ALL TOLAKAN VALIDASI
// =============================================================================

// GetAllTolakanValidasi handles POST /validasi-kpi/get-all-tolakan
func (h *ValidasiKpiHandler) GetAllTolakanValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllTolakanValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAllTolakanValidasi(&req)
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

// =============================================================================
// GET ALL DAFTAR PENYUSUNAN VALIDASI
// =============================================================================

// GetAllDaftarPenyusunanValidasi handles POST /validasi-kpi/get-all-daftar-penyusunan
func (h *ValidasiKpiHandler) GetAllDaftarPenyusunanValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarPenyusunanValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAllDaftarPenyusunanValidasi(&req)
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
		Message:    "Data Daftar Penyusunan Validasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// =============================================================================
// GET ALL DAFTAR APPROVAL VALIDASI
// =============================================================================

// GetAllDaftarApprovalValidasi handles POST /validasi-kpi/get-all-daftar-approval
func (h *ValidasiKpiHandler) GetAllDaftarApprovalValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarApprovalValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, _, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.ApprovalUser = userid

	data, total, err := h.service.GetAllDaftarApprovalValidasi(&req)
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
// GET ALL VALIDASI
// =============================================================================

// GetAllValidasi handles POST /validasi-kpi/get-all-validasi
func (h *ValidasiKpiHandler) GetAllValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllValidasiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAllValidasi(&req)
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

// =============================================================================
// VALIDASI BATAL
// =============================================================================

// ValidasiBatal handles POST /validasi-kpi/batal
// Menerima application/json. Membatalkan proses validasi (status → 91).
func (h *ValidasiKpiHandler) ValidasiBatal(c *gin.Context) {
	req, err := binder.BindJSON[dto.ValidasiBatalRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	userid, _, ok := parseUserqHeader(c)
	if !ok {
		return
	}
	req.User = userid

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ValidasiBatal(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Validasi KPI berhasil dibatalkan",
		Data:    data,
	})
}
