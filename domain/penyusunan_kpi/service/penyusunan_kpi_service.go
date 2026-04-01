package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/utils"
	customErrors "permen_api/errors"

	"github.com/jung-kurt/gofpdf"
)

// =============================================================================
// VALIDATE
// =============================================================================

func (s *penyusunanKpiService) ValidatePenyusunanKpi(
	req *dto.ValidatePenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidatePenyusunanKpiResponse, err error) {

	// User error: tidak mengirim file
	if file == nil {
		return data, &customErrors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files'",
		}
	}

	// User error: format file salah
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("file '%s' bukan format Excel (.xlsx)", file.Filename),
		}
	}

	// User error: isi Excel tidak valid (kolom salah, bobot salah, KPI tidak cocok, dll)
	kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan, req.Kpi)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// User error: polarisasi tidak sesuai master, lookup gagal karena data tidak valid
	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return data, err
	}

	// System error: gagal simpan ke DB — biarkan naik apa adanya
	idPengajuan, err := s.repo.ValidatePenyusunanKpi(req, kpiSubDetails)
	if err != nil {
		return data, err
	}

	data = dto.ValidatePenyusunanKpiResponse{
		IDPengajuan: idPengajuan,
		Tahun:       req.Tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.DivisiResponse{
			Kostl:   req.Divisi.Kostl,
			KostlTx: req.Divisi.KostlTx,
		},
		Entry: dto.EntryResponse{
			EntryUser: req.EntryUser,
			EntryName: req.EntryName,
			EntryTime: req.EntryTime,
		},
		TotalKpi:      len(req.Kpi),
		Kpi:           utils.BuildKpiResponse(idPengajuan, req.Kpi, kpiSubDetails),
		ChallengeList: req.ChallengeList,
		MethodList:    req.MethodList,
	}

	return data, nil
}

// =============================================================================
// CREATE
// =============================================================================

func (s *penyusunanKpiService) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
) (data dto.CreatePenyusunanKpiResponse, err error) {
	// User error (idPengajuan tidak ada) atau system error (DB) — repo sudah wrap dengan tipe yang tepat
	if err = s.repo.CreatePenyusunanKpi(req); err != nil {
		return data, err
	}

	data = dto.CreatePenyusunanKpiResponse{
		IdPengajuan:  req.IdPengajuan,
		ApprovalList: req.ApprovalList,
	}

	return data, nil
}

// =============================================================================
// HELPER — resolveMasterLookup
// =============================================================================

// resolveMasterLookup melakukan lookup mst_kpi dan mst_polarisasi untuk setiap
// baris sub KPI, lalu memvalidasi kesesuaian polarisasi dengan rumus di mst_kpi.
func (s *penyusunanKpiService) resolveMasterLookup(
	kpiSubDetails map[int][]dto.PenyusunanKpiSubDetailRow,
) error {
	for i, rows := range kpiSubDetails {
		for j := range rows {
			subRow := &kpiSubDetails[i][j]

			idKpi, kpiFromDB, rumusMstKpi, err := s.repo.LookupSubKpiMaster(subRow.SubKPI)
			if err != nil {
				// System error: query DB gagal
				return fmt.Errorf(
					"KPI ke-%d, Sub KPI ke-%d ('%s'): gagal lookup master KPI: %w",
					i+1, j+1, subRow.SubKPI, err,
				)
			}
			subRow.IdSubKpi = idKpi
			subRow.SubKPI = kpiFromDB
			if idKpi != "0" {
				subRow.Otomatis = "1"
			} else {
				subRow.Otomatis = "0"
			}

			idPolarisasi, err := s.repo.LookupPolarisasi(subRow.Polarisasi)
			if err != nil {
				// User error: polarisasi yang diisi di Excel tidak valid
				return &customErrors.BadRequestError{
					Message: fmt.Sprintf(
						"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi '%s' tidak valid: %s",
						i+1, j+1, subRow.SubKPI, subRow.Polarisasi, err.Error(),
					),
				}
			}
			subRow.IdPolarisasi = idPolarisasi

			if subRow.IdSubKpi != "0" {
				polarisasiMaster := "Maximize"
				if rumusMstKpi == "0" {
					polarisasiMaster = "Minimize"
				}
				if idPolarisasi != rumusMstKpi {
					// User error: polarisasi tidak cocok dengan master KPI
					return &customErrors.BadRequestError{
						Message: fmt.Sprintf(
							"KPI ke-%d, Sub KPI ke-%d ('%s'): polarisasi tidak sesuai master. "+
								"Excel: '%s' (id=%s), master KPI: '%s' (id=%s). "+
								"Periksa kembali kolom D pada file Excel",
							i+1, j+1, subRow.SubKPI,
							subRow.Polarisasi, idPolarisasi,
							polarisasiMaster, rumusMstKpi,
						),
					}
				}
			}
		}
	}
	return nil
}

// =============================================================================
// GET ALL DRAFT
// =============================================================================

func (s *penyusunanKpiService) GetAllDraftPenyusunanKpi(
	req *dto.GetAllDraftPenyusunanKpiRequest,
) (data []*dto.GetAllDraftPenyusunanKpiResponse, total int64, err error) {
	data, total, err = s.repo.GetAllDraftPenyusunanKpi(req)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

// =============================================================================
// GET DETAIL
// =============================================================================

func (s *penyusunanKpiService) GetDetailPenyusunanKpi(
	req *dto.GetDetailPenyusunanKpiRequest,
) (*dto.GetAllDraftPenyusunanKpiResponse, error) {
	data, err := s.repo.GetDetailPenyusunanKpi(req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// =============================================================================
// GET CSV
// =============================================================================

func (s *penyusunanKpiService) GetCsvPenyusunanKpi(
	req *dto.GetCsvPenyusunanKpiRequest,
) ([]byte, string, error) {
	exportData, err := s.repo.GetKpiExportData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	_ = writer.Write([]string{exportData.NamaDivisi})
	_ = writer.Write([]string{"Tahun " + exportData.Tahun})
	_ = writer.Write([]string{})

	_ = writer.Write([]string{"No", "KPI", "Bobot (%)", "Target Tahunan", "Capping"})

	for _, row := range exportData.Rows {
		_ = writer.Write([]string{
			strconv.Itoa(row.No),
			row.KpiNama,
			row.Bobot,
			row.TargetTahunan,
			row.Capping,
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", fmt.Errorf("gagal menulis CSV: %w", err)
	}

	filename := fmt.Sprintf("KPI_%s_%s_%s.csv",
		exportData.NamaDivisi, exportData.Tahun, exportData.Triwulan)

	return buf.Bytes(), filename, nil
}

// =============================================================================
// GET PDF
// =============================================================================

func (s *penyusunanKpiService) GetPdfPenyusunanKpi(
	req *dto.GetPdfPenyusunanKpiRequest,
) ([]byte, string, error) {
	exportData, err := s.repo.GetKpiExportData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Warna palette (sesuai gambar)
	headerBgR, headerBgG, headerBgB := 31, 73, 125   // biru tua  (#1F497D)
	headerFgR, headerFgG, headerFgB := 255, 255, 255 // putih
	rowBlueR, rowBlueG, rowBlueB := 189, 215, 238    // biru muda (#BDD7EE)
	rowPeachR, rowPeachG, rowPeachB := 252, 228, 214 // peach     (#FCE4D6)
	rowGreenR, rowGreenG, rowGreenB := 226, 240, 217 // hijau muda (#E2F0D9)
	textR, textG, textB := 0, 0, 0

	// Judul
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(textR, textG, textB)
	pdf.CellFormat(0, 7, exportData.NamaDivisi, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "Tahun "+exportData.Tahun, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	// Lebar kolom — total ~267mm untuk A4 landscape (297 - 15*2 margin)
	// No | KPI | Bobot (%) | Target Tahunan | Capping
	colWidths := []float64{12, 100, 25, 80, 25}
	headers := []string{"No", "KPI", "Bobot (%)", "Target Tahunan", "Capping"}
	rowHeight := 8.0

	// Header tabel
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(headerBgR, headerBgG, headerBgB)
	pdf.SetTextColor(headerFgR, headerFgG, headerFgB)
	pdf.SetDrawColor(255, 255, 255)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], rowHeight, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Baris data — alternating per 3 baris
	pdf.SetFont("Arial", "", 9)
	pdf.SetDrawColor(200, 200, 200)
	dataAligns := []string{"C", "L", "C", "L", "C"}

	for _, row := range exportData.Rows {
		group := ((row.No - 1) / 3) % 3
		switch group {
		case 0:
			pdf.SetFillColor(rowBlueR, rowBlueG, rowBlueB)
		case 1:
			pdf.SetFillColor(rowPeachR, rowPeachG, rowPeachB)
		default:
			pdf.SetFillColor(rowGreenR, rowGreenG, rowGreenB)
		}
		pdf.SetTextColor(textR, textG, textB)

		values := []string{
			strconv.Itoa(row.No),
			row.KpiNama,
			row.Bobot,
			row.TargetTahunan,
			row.Capping,
		}
		for i, v := range values {
			pdf.CellFormat(colWidths[i], rowHeight, v, "1", 0, dataAligns[i], true, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("gagal generate PDF: %w", err)
	}

	filename := fmt.Sprintf("KPI_%s_%s_%s.pdf",
		exportData.NamaDivisi, exportData.Tahun, exportData.Triwulan)

	return buf.Bytes(), filename, nil
}
