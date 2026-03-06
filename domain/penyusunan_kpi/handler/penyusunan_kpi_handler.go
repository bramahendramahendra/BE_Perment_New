package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"

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
// Catatan khusus ApprovalList:
//
//	Frontend mengirim ApprovalList sebagai JSON string tanpa escape inner quotes:
//	"ApprovalList": "[{"userid":"xxx",...}]"
//	Handler melakukan sanitize otomatis — extract ApprovalList terlebih dahulu,
//	replace dengan placeholder, parse JSON, lalu set nilai aslinya kembali.
func (h *PenyusunanKpiHandler) InsertKPI(c *gin.Context) {
	// --- 1. Ambil field REQUEST (JSON string) dari multipart form ---
	requestStr := c.PostForm("REQUEST")
	if requestStr == "" {
		c.Error(&errors.BadRequestError{Message: "field 'REQUEST' tidak boleh kosong"})
		return
	}

	// --- 2. Extract & sanitize ApprovalList sebelum parse JSON ---
	// Frontend mengirim ApprovalList dengan inner quotes tanpa escape,
	// sehingga JSON menjadi invalid. Kita extract nilainya terlebih dahulu,
	// ganti dengan placeholder, baru parse JSON sisanya.
	sanitizedStr, approvalListRaw, err := extractApprovalList(requestStr)
	if err != nil {
		c.Error(&errors.BadRequestError{
			Message: "format field 'ApprovalList' tidak valid: " + err.Error(),
		})
		return
	}

	// --- 3. Parse JSON string ke struct ---
	var req dto.InsertPenyusunanKpiRequest
	if err := json.Unmarshal([]byte(sanitizedStr), &req); err != nil {
		c.Error(&errors.BadRequestError{
			Message: "format 'REQUEST' tidak valid: " + err.Error(),
		})
		return
	}

	// --- 4. Set ApprovalList dengan nilai asli yang sudah diextract ---
	// Nilai ini disimpan persis seperti yang dikirim frontend (tanpa escape)
	req.ApprovalList = approvalListRaw

	// --- 5. Validasi struct menggunakan validator ---
	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

	// --- 6. Ambil semua file Excel dari field "files[]" ---
	form, err := c.MultipartForm()
	if err != nil {
		c.Error(&errors.BadRequestError{
			Message: "gagal membaca multipart form: " + err.Error(),
		})
		return
	}

	files, ok := form.File["files"]
	if !ok || len(files) == 0 {
		c.Error(&errors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files'",
		})
		return
	}

	// --- 7. Validasi ekstensi file harus .xlsx ---
	if err := validateExcelFiles(files); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// --- 8. Panggil service ---
	idPengajuan, err := h.service.InsertPenyusunanKpi(&req, files)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	// --- 9. Return response sukses ---
	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil disimpan",
		Data: dto.InsertPenyusunanKpiResponse{
			IDPengajuan: idPengajuan,
			Message:     "Insert KPI berhasil",
		},
	})
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// extractApprovalList mengekstrak nilai ApprovalList dari raw REQUEST string
// dan menggantikannya dengan placeholder yang valid secara JSON.
//
// Masalah yang diselesaikan:
//
//	Frontend mengirim ApprovalList TANPA escape inner quotes:
//	  "ApprovalList": "[{"userid":"xxx","nama":"yyy",...}]"
//	Ini membuat seluruh JSON REQUEST menjadi invalid.
//
// Cara kerja:
//  1. Cari marker: `"ApprovalList":` lalu cari tanda `"[` sebagai awal value
//  2. Scan karakter satu per satu dari `[` sampai menemukan `]"` sebagai penutup
//  3. Ekstrak raw value (isi array JSON asli)
//  4. Ganti di requestStr dengan placeholder `"__APPROVAL_PLACEHOLDER__"`
//  5. Return: sanitizedStr (JSON valid), approvalListRaw (nilai asli), error
//
// Setelah json.Unmarshal pada sanitizedStr berhasil,
// req.ApprovalList di-set manual dengan approvalListRaw.
func extractApprovalList(requestStr string) (sanitizedStr, approvalListRaw string, err error) {
	// Cari posisi key "ApprovalList"
	key := `"ApprovalList"`
	keyIdx := strings.Index(requestStr, key)
	if keyIdx == -1 {
		// ApprovalList tidak ada di request — kembalikan as-is, biarkan validator handle
		return requestStr, "", nil
	}

	// Cari tanda "[ setelah key (value ApprovalList selalu berupa array JSON string)
	afterKey := requestStr[keyIdx+len(key):]

	// Lewati spasi dan ":"
	colonIdx := strings.Index(afterKey, `"[`)
	if colonIdx == -1 {
		return requestStr, "", fmt.Errorf("format ApprovalList tidak ditemukan, pastikan berupa array JSON string")
	}

	// Posisi awal "[" (tanpa tanda kutip pembuka)
	valueStart := keyIdx + len(key) + colonIdx + 1 // +1 untuk skip karakter "

	// Scan dari "[" sampai ketemu "]" yang diikuti tanda kutip penutup `"`
	// Ini untuk handle kasus inner quotes tanpa escape
	scanStr := requestStr[valueStart:]
	bracketDepth := 0
	endIdx := -1

	for i, ch := range scanStr {
		switch ch {
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
			if bracketDepth == 0 {
				// Cek apakah karakter berikutnya adalah " (penutup string JSON)
				if i+1 < len(scanStr) && scanStr[i+1] == '"' {
					endIdx = i
				}
			}
		}
		if endIdx != -1 {
			break
		}
	}

	if endIdx == -1 {
		return requestStr, "", fmt.Errorf("penutup ApprovalList tidak ditemukan, pastikan diakhiri dengan ]\"")
	}

	// Ekstrak raw value: dari "[" sampai "]" (inclusive)
	approvalListRaw = scanStr[:endIdx+1]

	// Ganti seluruh block "ApprovalList": "..." dengan placeholder yang valid JSON
	fullOriginal := requestStr[keyIdx : valueStart+endIdx+2] // +2 untuk include " penutup
	placeholder := `"ApprovalList": "__APPROVAL_PLACEHOLDER__"`
	sanitizedStr = strings.Replace(requestStr, fullOriginal, placeholder, 1)

	return sanitizedStr, approvalListRaw, nil
}

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
