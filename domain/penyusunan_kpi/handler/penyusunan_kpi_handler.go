package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
	service "permen_api/domain/penyusunan_kpi/service"
	globalDTO "permen_api/dto"
	"permen_api/errors"
	response_helper "permen_api/helper/response"
	validator "permen_api/validation"

	"github.com/gin-gonic/gin"
)

type PenyusunanKpiHandler struct {
	service service.PenyusunanKpiServiceInterface
}

func NewPenyusunanKpiHandler(service service.PenyusunanKpiServiceInterface) *PenyusunanKpiHandler {
	return &PenyusunanKpiHandler{service: service}
}

func (h *PenyusunanKpiHandler) InsertKPI(c *gin.Context) {
	requestStr := c.PostForm("REQUEST")
	if requestStr == "" {
		c.Error(&errors.BadRequestError{Message: "field 'REQUEST' tidak boleh kosong"})
		return
	}

	sanitizedStr, approvalListRaw, err := extractApprovalList(requestStr)
	if err != nil {
		c.Error(&errors.BadRequestError{
			Message: "format field 'ApprovalList' tidak valid: " + err.Error(),
		})
		return
	}

	var req dto.InsertPenyusunanKpiRequest
	if err := json.Unmarshal([]byte(sanitizedStr), &req); err != nil {
		c.Error(&errors.BadRequestError{
			Message: "format 'REQUEST' tidak valid: " + err.Error(),
		})
		return
	}

	req.ApprovalList = approvalListRaw
	req.Kostl = req.Divisi.Kostl
	req.KostlTx = req.Divisi.KostlTx

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
	req.EntryUser = strings.TrimSpace(parts[0])
	req.EntryName = strings.TrimSpace(parts[1])
	req.EntryTime = time.Now().Format("2006-01-02 15:04:05")

	if err := validator.Validate.Struct(req); err != nil {
		c.Error(err)
		return
	}

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

	if err := validateExcelFiles(files); err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	result, err := h.service.InsertPenyusunanKpi(&req, files)
	if err != nil {
		c.Error(&errors.BadRequestError{Message: err.Error()})
		return
	}

	kpiResp := buildKpiResponse(result.IDPengajuan, req.Kpi, result.KpiSubDetails)

	response_helper.WrapResponse(c, 200, "json", &globalDTO.ResponseParams{
		Code:    "00",
		Status:  true,
		Message: "Data KPI berhasil disimpan",
		Data: dto.InsertPenyusunanKpiResponse{
			IDPengajuan:    result.IDPengajuan,
			Tahun:          req.Tahun,
			Triwulan:       req.Triwulan,
			Kostl:          req.Kostl,
			KostlTx:        req.KostlTx,
			EntryUser:      req.EntryUser,
			EntryName:      req.EntryName,
			EntryTime:      req.EntryTime,
			ApprovalPosisi: req.ApprovalPosisi,
			SaveAsDraft:    req.SaveAsDraft,
			TotalKpi:       len(req.Kpi),
			Kpi:            kpiResp,
			ChallengeList:  req.ChallengeList,
			MethodList:     req.MethodList,
		},
	})
}

// buildKpiResponse membangun slice KpiDetailItemResponse dengan KpiSubDetail
// yang sudah di-nested ke dalam masing-masing KPI sesuai indeksnya.
func buildKpiResponse(
	idPengajuan string,
	kpiList []dto.PenyusunanKpiDetailItemRequest,
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) []dto.PenyusunanKpiDetailItemResponse {
	result := make([]dto.PenyusunanKpiDetailItemResponse, len(kpiList))
	subCounter := 1

	for i, kpiItem := range kpiList {
		idDetail := fmt.Sprintf("%sP%03d", idPengajuan, i+1)
		rows := kpiSubDetails[i]

		subDetails := make([]dto.KpiSubDetailResponse, 0, len(rows))
		for _, row := range rows {
			idSubDetail := fmt.Sprintf("%sC%03d", idPengajuan, subCounter)
			subCounter++

			subDetails = append(subDetails, dto.KpiSubDetailResponse{
				IdDetail:                  idDetail,
				IdSubDetail:               idSubDetail,
				NamaKpi:                   row.KPI,
				IdSubKpi:                  row.IdSubKpi,
				SubKpi:                    row.SubKPI,
				Otomatis:                  row.Otomatis,
				Polarisasi:                row.Polarisasi,
				IdPolarisasi:              row.IdPolarisasi,
				Capping:                   row.Capping,
				Bobot:                     row.Bobot,
				Glossary:                  row.Glossary,
				TargetTriwulan:            row.TargetTriwulan,
				TargetKuantitatifTriwulan: row.TargetKuantitatifTriwulan,
				TargetTahunan:             row.TargetTahunan,
				TargetKuantitatifTahunan:  row.TargetKuantitatifTahunan,
				TerdapatQualifier:         row.TerdapatQualifier,
				Qualifier:                 row.Qualifier,
				DeskripsiQualifier:        row.DeskripsiQualifier,
				TargetQualifier:           row.TargetQualifier,
				Result:                    row.Result,
				DeskripsiResult:           row.DeskripsiResult,
				Process:                   row.Process,
				DeskripsiProcess:          row.DeskripsiProcess,
				Context:                   row.Context,
				DeskripsiContext:          row.DeskripsiContext,
			})
		}

		result[i] = dto.PenyusunanKpiDetailItemResponse{
			IdKpi:        kpiItem.IdKpi,
			Kpi:          kpiItem.Kpi,
			Rumus:        kpiItem.Rumus,
			Persfektif:   kpiItem.Persfektif,
			KpiSubDetail: subDetails,
		}
	}

	return result
}

func extractApprovalList(requestStr string) (sanitizedStr, approvalListRaw string, err error) {
	marker := `"ApprovalList":`
	markerIdx := strings.Index(requestStr, marker)
	if markerIdx == -1 {
		return requestStr, "", nil
	}

	afterMarker := requestStr[markerIdx+len(marker):]
	afterMarker = strings.TrimLeft(afterMarker, " \t\n\r")

	if !strings.HasPrefix(afterMarker, `"[`) {
		return requestStr, "", nil
	}

	startIdx := markerIdx + len(marker) + strings.Index(requestStr[markerIdx+len(marker):], `"[`)
	contentStart := startIdx + 1

	content := requestStr[contentStart:]
	endIdx := strings.Index(content, `]"`)
	if endIdx == -1 {
		return "", "", fmt.Errorf("tidak menemukan penutup ']\"' pada ApprovalList")
	}

	approvalListRaw = content[:endIdx+1]
	placeholder := `"__APPROVAL_PLACEHOLDER__"`
	sanitizedStr = requestStr[:startIdx] + placeholder + requestStr[contentStart+endIdx+2:]

	return sanitizedStr, approvalListRaw, nil
}

func validateExcelFiles(files []*multipart.FileHeader) error {
	for _, f := range files {
		if !strings.HasSuffix(strings.ToLower(f.Filename), ".xlsx") {
			return fmt.Errorf("file '%s' bukan format Excel (.xlsx)", f.Filename)
		}
	}
	return nil
}
