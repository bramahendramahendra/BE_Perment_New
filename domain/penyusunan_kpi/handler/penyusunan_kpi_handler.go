package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	dto "permen_api/domain/penyusunan_kpi/dto"
	service "permen_api/domain/penyusunan_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

// =============================================
// HANDLER STRUCT & CONSTRUCTOR
// =============================================

type PenyusunanKpiHandler struct {
	service service.PenyusunanKpiServiceInterface
}

func NewPenyusunanKpiHandler(service service.PenyusunanKpiServiceInterface) *PenyusunanKpiHandler {
	return &PenyusunanKpiHandler{service: service}
}

// =============================================
// HANDLER METHODS
// =============================================

// InsertKPI menerima request multipart/form-data dengan:
//   - Field "REQUEST" : JSON string berisi InsertPenyusunanKpiRequest
//   - Field "files[]" : satu atau lebih file Excel (.xlsx)
//     urutan file harus sesuai urutan array Kpi di REQUEST
//
// Contoh request:
//
//	POST /api/penyusunan-kpi/insert
//	Content-Type: multipart/form-data
//	Authorization: Bearer <token>
//
//	REQUEST = '{"IDPengajuan":"...","Tahun":"2025",...,"Kpi":[...],...}'
//	files[]  = [excel_kpi_1.xlsx, excel_kpi_2.xlsx]
func (h *PenyusunanKpiHandler) InsertKPI(c *gin.Context) {
	// --- 1. Ambil field REQUEST (JSON string) dari multipart form ---
	requestStr := c.PostForm("REQUEST")
	if requestStr == "" {
		c.Error(&errors.BadRequestError{Message: "field 'REQUEST' tidak boleh kosong"})
		return
	}

	// --- 2. Parse JSON string ke struct ---
	var req dto.InsertPenyusunanKpiRequest
	if err := json.Unmarshal([]byte(requestStr), &req); err != nil {
		c.Error(&errors.BadRequestError{
			Message: "format 'REQUEST' tidak valid, pastikan berupa JSON yang benar: " + err.Error(),
		})
		return
	}

	// --- 3. Validasi struct menggunakan validator ---
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	// --- 4. Ambil semua file Excel dari field "files[]" ---
	form, err := c.MultipartForm()
	if err != nil {
		c.Error(&errors.BadRequestError{
			Message: "gagal membaca multipart form: " + err.Error(),
		})
		return
	}

	files, ok := form.File["files[]"]
	if !ok || len(files) == 0 {
		c.Error(&errors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files[]'",
		})
		return
	}

	// --- 5. Validasi ekstensi file harus .xlsx ---
	if err := validateExcelFiles(files); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// --- 6. Panggil service ---
	if err := h.service.InsertPenyusunanKpi(&req, files); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// --- 7. Return response sukses ---
	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil disimpan",
		Data: dto.InsertPenyusunanKpiResponse{
			IDPengajuan: req.IDPengajuan,
			Message:     "Insert KPI berhasil",
		},
	})
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// validateExcelFiles memastikan semua file yang di-upload berekstensi .xlsx
func validateExcelFiles(files []*multipart.FileHeader) error {
	for _, file := range files {
		filename := file.Filename
		if len(filename) < 5 || filename[len(filename)-5:] != ".xlsx" {
			return fmt.Errorf(
				"file '%s' tidak valid, hanya file Excel (.xlsx) yang diizinkan",
				filename,
			)
		}
	}
	return nil
}
