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

// =============================================================================
// INPUT VALIDASI (validate + create + revision dalam satu endpoint)
// =============================================================================

// InputValidasi handles POST /validasi-kpi/input
// Menerima application/json. Menyimpan data validasi KPI dan mengirim notifikasi ke approver (status → 6).
// Berlaku untuk status 5 (baru), 7 (revisi setelah tolak), 90/91 (ulang setelah batal).
func (h *ValidasiKpiHandler) InputValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.InputValidasiRequest](c)
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
// APPROVE VALIDASI
// =============================================================================

// ApproveValidasi handles POST /validasi-kpi/approve
// Menerima application/json. Memproses approve validasi KPI dalam rantai approval.
func (h *ValidasiKpiHandler) ApproveValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.ApproveValidasiRequest](c)
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
// Menerima application/json. Memproses penolakan validasi KPI (status → 7).
func (h *ValidasiKpiHandler) RejectValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RejectValidasiRequest](c)
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
	req.User = strings.TrimSpace(parts[0])

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

// =============================================================================
// GET ALL APPROVAL VALIDASI
// =============================================================================

// GetAllApprovalValidasi handles POST /validasi-kpi/get-all-approval
// Menerima application/json. Mengambil list pengajuan yang menunggu approval user (status=6).
func (h *ValidasiKpiHandler) GetAllApprovalValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllApprovalValidasiRequest](c)
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
// Menerima application/json. Mengambil list pengajuan yang ditolak (status=7).
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
// Menerima application/json.
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
// Menerima application/json. Mengambil semua pengajuan yang pernah melibatkan user dalam approval validasi.
func (h *ValidasiKpiHandler) GetAllDaftarApprovalValidasi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarApprovalValidasiRequest](c)
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
// Menerima application/json.
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
// GET DETAIL VALIDASI
// =============================================================================

// GetDetailValidasiKpi handles POST /validasi-kpi/get-detail
// Menerima application/json. Mengambil detail lengkap satu pengajuan validasi KPI.
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
