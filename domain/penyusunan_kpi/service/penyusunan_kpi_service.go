package service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/utils"
	customErrors "permen_api/errors"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
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

	// Parse dan validasi file Excel.
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup mst_kpi untuk setiap KPI unik dari kolom B Excel.
	// Jika tidak ditemukan: idKpi = "0", rumus = "0".
	if err := s.resolveKpiMasterLookup(kpiRows); err != nil {
		return data, err
	}

	// Lookup mst_kpi dan mst_polarisasi untuk setiap baris sub KPI (kolom C).
	// Validasi polarisasi vs rumus mst_kpi juga dilakukan di sini.
	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return data, err
	}

	// Bangun idPengajuan di service agar bisa digunakan untuk build response sebelum repo insert.
	idPengajuan := utils.GenerateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	// Build ChallengeList dan MethodList dari data Excel (kolom P,Q,R,S,T,U).
	// Hanya diisi untuk TW2 dan TW4; untuk TW1 dan TW3 list kosong (tidak diinsert ke DB).
	resultList := []dto.PenyusunanResult{}
	methodList := []dto.PenyusunanMethod{}
	challengeList := []dto.PenyusunanChallenge{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		methodList = utils.BuildMethodList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		challengeList = utils.BuildChallengeList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	idPengajuan, err = s.repo.ValidatePenyusunanKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		methodList,
		challengeList,
	)
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
		TotalKpi:      len(kpiRows),
		Kpi:           utils.BuildKpiResponse(idPengajuan, kpiRows, kpiSubDetails),
		ResultList:    resultList,
		MethodList:    methodList,
		ChallengeList: challengeList,
	}

	return data, nil
}

// =============================================================================
// REVISION
// =============================================================================

func (s *penyusunanKpiService) RevisionPenyusunanKpi(
	req *dto.RevisionPenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.RevisionPenyusunanKpiResponse, err error) {

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

	// Parse dan validasi file Excel (logika sama dengan /validate)
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup mst_kpi untuk setiap KPI unik dari kolom B Excel
	if err := s.resolveKpiMasterLookup(kpiRows); err != nil {
		return data, err
	}

	// Lookup mst_kpi dan mst_polarisasi untuk setiap baris sub KPI (kolom C)
	if err := s.resolveMasterLookup(kpiSubDetails); err != nil {
		return data, err
	}

	// Build ChallengeList, MethodList, ResultList dari kolom P–U Excel
	// Hanya diisi untuk TW2 dan TW4
	resultList := []dto.PenyusunanResult{}
	methodList := []dto.PenyusunanMethod{}
	challengeList := []dto.PenyusunanChallenge{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		methodList = utils.BuildMethodList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		challengeList = utils.BuildChallengeList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	// Simpan ke DB: DELETE lama + INSERT baru + UPDATE header
	if err := s.repo.RevisionPenyusunanKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		methodList,
		challengeList,
	); err != nil {
		return data, err
	}

	data = dto.RevisionPenyusunanKpiResponse{
		IDPengajuan: req.IdPengajuan,
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
		TotalKpi:      len(kpiRows),
		Kpi:           utils.BuildKpiResponse(req.IdPengajuan, kpiRows, kpiSubDetails),
		ResultList:    resultList,
		MethodList:    methodList,
		ChallengeList: challengeList,
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
// HELPER — resolveKpiMasterLookup
// =============================================================================

// resolveKpiMasterLookup melakukan lookup mst_kpi untuk setiap KPI unik dari kolom B Excel.
// Aturan:
//   - Jika ditemukan → idKpi dan rumus dari DB
//   - Jika tidak ditemukan → idKpi = "0", rumus = "0"
func (s *penyusunanKpiService) resolveKpiMasterLookup(
	kpiRows []dto.PenyusunanKpiRow,
) error {
	for i := range kpiRows {
		idKpi, _, rumus, err := s.repo.LookupKpiMaster(kpiRows[i].Kpi)
		if err != nil {
			return fmt.Errorf(
				"KPI '%s': gagal lookup master KPI: %w",
				kpiRows[i].Kpi, err,
			)
		}

		if idKpi == "0" {
			// Tidak ditemukan di mst_kpi → idKpi = "0", rumus = "0"
			kpiRows[i].IdKpi = "0"
			kpiRows[i].Rumus = "0"
		} else {
			kpiRows[i].IdKpi = idKpi
			kpiRows[i].Rumus = rumus
		}
	}
	return nil
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

			idKpi, kpiFromDB, rumusMstKpi, err := s.repo.LookupKpiMaster(subRow.SubKPI)
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
// GET ALL APPROVAL
// =============================================================================

func (s *penyusunanKpiService) GetAllApprovalPenyusunanKpi(
	req *dto.GetAllApprovalPenyusunanKpiRequest,
) (data []*dto.GetAllApprovalPenyusunanKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllApprovalPenyusunanKpi(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllApprovalPenyusunanKpiResponse{
			IdPengajuan: v.IdPengajuan,
			Tahun:       v.Tahun,
			Triwulan:    v.Triwulan,
			Kostl:       v.Kostl,
			KostlTx:     v.KostlTx,
			Orgeh:       v.Orgeh,
			OrgehTx:     v.OrgehTx,
			Status:      v.Status,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// =============================================================================
// GET ALL TOLAKAN
// =============================================================================

func (s *penyusunanKpiService) GetAllTolakanPenyusunanKpi(
	req *dto.GetAllTolakanPenyusunanKpiRequest,
) (data []*dto.GetAllTolakanPenyusunanKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllTolakanPenyusunanKpi(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllTolakanPenyusunanKpiResponse{
			IdPengajuan: v.IdPengajuan,
			Tahun:       v.Tahun,
			Triwulan:    v.Triwulan,
			Kostl:       v.Kostl,
			KostlTx:     v.KostlTx,
			Orgeh:       v.Orgeh,
			OrgehTx:     v.OrgehTx,
			Status:      v.Status,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// =============================================================================
// GET ALL DAFTAR PENYUSUNAN
// =============================================================================

func (s *penyusunanKpiService) GetAllDaftarPenyusunanKpi(
	req *dto.GetAllDaftarPenyusunanKpiRequest,
) (data []*dto.GetAllDaftarPenyusunanKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarPenyusunanKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarPenyusunanKpiResponse{
			IdPengajuan: v.IdPengajuan,
			Tahun:       v.Tahun,
			Triwulan:    v.Triwulan,
			Kostl:       v.Kostl,
			KostlTx:     v.KostlTx,
			Orgeh:       v.Orgeh,
			OrgehTx:     v.OrgehTx,
			Status:      v.Status,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// =============================================================================
// GET ALL DAFTAR APPROVAL
// =============================================================================

func (s *penyusunanKpiService) GetAllDaftarApprovalPenyusunanKpi(
	req *dto.GetAllDaftarApprovalPenyusunanKpiRequest,
) (data []*dto.GetAllDaftarApprovalPenyusunanKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarApprovalPenyusunanKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarApprovalPenyusunanKpiResponse{
			IdPengajuan: v.IdPengajuan,
			Tahun:       v.Tahun,
			Triwulan:    v.Triwulan,
			Kostl:       v.Kostl,
			KostlTx:     v.KostlTx,
			Orgeh:       v.Orgeh,
			OrgehTx:     v.OrgehTx,
			Status:      v.Status,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// =============================================================================
// GET DETAIL
// =============================================================================

func (s *penyusunanKpiService) GetDetailPenyusunanKpi(
	req *dto.GetDetailPenyusunanKpiRequest,
) (data *dto.GetAllDataPenyusunanKpiResponse, err error) {
	dataDB, err := s.repo.GetDetailPenyusunanKpi(req)
	if err != nil {
		return nil, err
	}
	return dataDB, nil
}

// =============================================================================
// GET EXCEL
// =============================================================================

func (s *penyusunanKpiService) GetExcelPenyusunanKpi(
	req *dto.GetExcelPenyusunanKpiRequest,
) ([]byte, string, error) {
	exportData, err := s.repo.GetKpiExportData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}

	const sheetName = "Data KPI"

	f := excelize.NewFile()
	defer f.Close()

	// Rename default sheet "Sheet1" → "Data KPI"
	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, "", fmt.Errorf("gagal set nama sheet: %w", err)
	}

	// -------------------------------------------------------------------------
	// Lebar kolom
	// A  = No             → sempit
	// B  = KPI            → lebar (konten teks panjang)
	// C  = Bobot (%)      → sedang
	// D  = Target Tahunan → sedang
	// E  = Capping        → sedang
	// -------------------------------------------------------------------------
	colWidths := map[string]float64{
		"A": 6,
		"B": 40,
		"C": 14,
		"D": 20,
		"E": 14,
	}
	for col, width := range colWidths {
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return nil, "", fmt.Errorf("gagal set lebar kolom %s: %w", col, err)
		}
	}

	// -------------------------------------------------------------------------
	// Baris 1: Nama Divisi — merge A1:E1
	// -------------------------------------------------------------------------
	if err := f.MergeCell(sheetName, "A1", "E1"); err != nil {
		return nil, "", fmt.Errorf("gagal merge cell baris 1: %w", err)
	}
	if err := f.SetCellValue(sheetName, "A1", exportData.NamaDivisi); err != nil {
		return nil, "", fmt.Errorf("gagal menulis nama divisi: %w", err)
	}

	// -------------------------------------------------------------------------
	// Baris 2: Tahun — merge A2:E2
	// -------------------------------------------------------------------------
	if err := f.MergeCell(sheetName, "A2", "E2"); err != nil {
		return nil, "", fmt.Errorf("gagal merge cell baris 2: %w", err)
	}
	if err := f.SetCellValue(sheetName, "A2", "Tahun "+exportData.Tahun); err != nil {
		return nil, "", fmt.Errorf("gagal menulis tahun: %w", err)
	}

	// -------------------------------------------------------------------------
	// Baris 3: Kosong
	// -------------------------------------------------------------------------
	// (tidak perlu set value — cell default kosong)

	// -------------------------------------------------------------------------
	// Baris 4: Header kolom
	// -------------------------------------------------------------------------
	headers := []string{"No", "KPI", "Bobot (%)", "Target Tahunan", "Capping"}
	for colIdx, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(colIdx+1, 4)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return nil, "", fmt.Errorf("gagal menulis header kolom '%s': %w", header, err)
		}
	}

	// -------------------------------------------------------------------------
	// Baris 5 dst: Data rows
	// -------------------------------------------------------------------------
	for i, row := range exportData.Rows {
		rowNum := 5 + i

		values := []interface{}{
			strconv.Itoa(row.No),
			row.KpiNama,
			row.Bobot,
			row.TargetTahunan,
			row.Capping,
		}

		for colIdx, val := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowNum)
			if err := f.SetCellValue(sheetName, cell, val); err != nil {
				return nil, "", fmt.Errorf("gagal menulis data baris %d kolom %d: %w", rowNum, colIdx+1, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Write ke buffer
	// -------------------------------------------------------------------------
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", fmt.Errorf("gagal menulis Excel ke buffer: %w", err)
	}

	filename := fmt.Sprintf("KPI_%s_%s_%s.xlsx",
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
