package handler

import (
	"strings"
	"time"

	dto "permen_api/domain/realisasi_kpi/dto"
	service "permen_api/domain/realisasi_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	binder "permen_api/pkg/binder"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type RealisasiKpiHandler struct {
	service service.RealisasiKpiServiceInterface
}

func NewRealisasiKpiHandler(service service.RealisasiKpiServiceInterface) *RealisasiKpiHandler {
	return &RealisasiKpiHandler{service: service}
}

// =============================================================================
// VALIDATE
// =============================================================================

// ValidateRealisasiKpi handles POST /realisasi-kpi/validate
// Menerima multipart/form-data dengan field REQUEST (JSON) dan files (Excel realisasi).
func (h *RealisasiKpiHandler) ValidateRealisasiKpi(c *gin.Context) {
	req, file, err := binder.BindMultipartJSON[dto.ValidateRealisasiKpiRequest](c, "REQUEST", "files")
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

	req.EntryUserRealisasi = strings.TrimSpace(parts[0])
	req.EntryNameRealisasi = strings.TrimSpace(parts[1])
	req.EntryTimeRealisasi = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ValidateRealisasiKpi(&req, file)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data realisasi KPI berhasil disimpan sebagai draft",
		Data:    data,
	})
}

// =============================================================================
// CREATE
// =============================================================================

// CreateRealisasiKpi handles POST /realisasi-kpi/create
// Menerima application/json. Mengubah status draft (80) → pending approval (3).
func (h *RealisasiKpiHandler) CreateRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.CreateRealisasiKpiRequest](c)
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

	req.EntryUserRealisasi = strings.TrimSpace(parts[0])
	req.EntryNameRealisasi = strings.TrimSpace(parts[1])
	req.EntryTimeRealisasi = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.CreateRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data Realisasi KPI berhasil diajukan",
		Data:    data,
	})
}

// =============================================================================
// REVISION
// =============================================================================

// RevisionRealisasiKpi handles POST /realisasi-kpi/revision
// Menerima multipart/form-data dengan field REQUEST (JSON) dan files (Excel realisasi revisi).
func (h *RealisasiKpiHandler) RevisionRealisasiKpi(c *gin.Context) {
	req, file, err := binder.BindMultipartJSON[dto.RevisionRealisasiKpiRequest](c, "REQUEST", "files")
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

	req.EntryUserRealisasi = strings.TrimSpace(parts[0])
	req.EntryNameRealisasi = strings.TrimSpace(parts[1])
	req.EntryTimeRealisasi = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.RevisionRealisasiKpi(&req, file)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Revisi data realisasi KPI berhasil disimpan",
		Data:    data,
	})
}

// ApproveRealisasiKpi handles POST /realisasi-kpi/approve
func (h *RealisasiKpiHandler) ApproveRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.ApproveRealisasiKpiRequest](c)
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

	req.ApprovalUserRealisasi = strings.TrimSpace(parts[0])
	req.ApprovalNameRealisasi = strings.TrimSpace(parts[1])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.ApproveRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Realisasi KPI berhasil diapprove",
		Data:    data,
	})
}

// RejectRealisasiKpi handles POST /realisasi-kpi/reject
func (h *RealisasiKpiHandler) RejectRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.RejectRealisasiKpiRequest](c)
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

	req.ApprovalUserRealisasi = strings.TrimSpace(parts[0])
	req.ApprovalNameRealisasi = strings.TrimSpace(parts[1])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.RejectRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Realisasi KPI berhasil ditolak",
		Data:    data,
	})
}

// GetAllDaftarRealisasiKpi handles POST /realisasi-kpi/get-all
func (h *RealisasiKpiHandler) GetAllRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	data, total, err := h.service.GetAllRealisasiKpi(&req)
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
		Message:    "Data Realisasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllApprovalRealisasiKpi handles POST /realisasi-kpi/get-all-approval
func (h *RealisasiKpiHandler) GetAllApprovalRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllApprovalRealisasiKpiRequest](c)
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

	req.ApprovalUserRealisasi = strings.TrimSpace(parts[0])

	data, total, err := h.service.GetAllApprovalRealisasiKpi(&req)
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
		Message:    "Data Approval Realisasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllTolakanRealisasiKpi handles POST /realisasi-kpi/get-all-tolakan
// Mengembalikan list pengajuan realisasi yang ditolak milik user tertentu (status 4).
func (h *RealisasiKpiHandler) GetAllTolakanRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllTolakanRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// if err := validator.Validate.Struct(req); err != nil {
	// 	c.Error(err)
	// 	return
	// }

	data, total, err := h.service.GetAllTolakanRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	pagination := response_helper.SetPagination(&globalDTO.FilterRequestParams{
		Page:  req.Page,
		Limit: req.Limit,
	}, total)

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:       "00",
		Status:     true,
		Message:    "Data Penolakan Realisasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllDaftarRealisasiKpi handles POST /realisasi-kpi/get-all-daftar-realisasi
// Mengembalikan semua pengajuan dalam konteks realisasi dengan filter opsional.
func (h *RealisasiKpiHandler) GetAllDaftarRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// if err := validator.Validate.Struct(req); err != nil {
	// 	c.Error(err)
	// 	return
	// }

	data, total, err := h.service.GetAllDaftarRealisasiKpi(&req)
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
		Message:    "Data Daftar Realisasi KPI berhasil diambil",
		Data:       data,
		Pagination: pagination,
	})
}

// GetAllDaftarApprovalRealisasiKpi handles POST /realisasi-kpi/get-all-daftar-approval
// Mengembalikan semua pengajuan realisasi yang sudah masuk proses approval.
func (h *RealisasiKpiHandler) GetAllDaftarApprovalRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetAllDaftarApprovalRealisasiKpiRequest](c)
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

	req.ApprovalUserRealisasi = strings.TrimSpace(parts[0])

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, total, err := h.service.GetAllDaftarApprovalRealisasiKpi(&req)
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

// GetDetailRealisasiKpi handles POST /realisasi-kpi/get-detail
// Mengembalikan detail lengkap satu pengajuan realisasi beserta nested KPI, context, process.
func (h *RealisasiKpiHandler) GetDetailRealisasiKpi(c *gin.Context) {
	req, err := binder.BindJSON[dto.GetDetailRealisasiKpiRequest](c)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	data, err := h.service.GetDetailRealisasiKpi(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data detail Realisasi KPI berhasil diambil",
		Data:    data,
	})
}
