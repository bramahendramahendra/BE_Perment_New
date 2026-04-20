package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	dto "permen_api/domain/penyusunan_kpi/dto"
	"permen_api/domain/penyusunan_kpi/utils"
	customErrors "permen_api/errors"
	file_export "permen_api/pkg/file_export"
	"permen_api/pkg/idgen"
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
	idPengajuan := idgen.GenerateIDPengajuan(req.Kostl, req.Tahun, req.Triwulan)

	// Build ProcessList dan ContextList dari data Excel (kolom P,Q,R,S,T,U).
	// Hanya diisi untuk TW2 dan TW4; untuk TW1 dan TW3 list kosong (tidak diinsert ke DB).
	resultList := []dto.PenyusunanResult{}
	processList := []dto.PenyusunanProcess{}
	contextList := []dto.PenyusunanContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(idPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	idPengajuan, err = s.repo.ValidatePenyusunanKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		processList,
		contextList,
	)
	if err != nil {
		return data, err
	}

	data = dto.ValidatePenyusunanKpiResponse{
		IDPengajuan: idPengajuan,
		Tahun:       req.Tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.Divisi{
			Kostl:   req.Divisi.Kostl,
			KostlTx: req.Divisi.KostlTx,
		},
		Entry: dto.EntryUser{
			EntryUser: req.EntryUser,
			EntryName: req.EntryName,
			EntryTime: req.EntryTime,
		},
		TotalKpi:    len(kpiRows),
		Kpi:         utils.BuildKpiResponse(idPengajuan, kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
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

	simpleList := make([]dto.ApprovalUserSimple, len(req.ApprovalList))
	for i, a := range req.ApprovalList {
		simpleList[i] = dto.ApprovalUserSimple{Userid: a.Userid, Nama: a.Nama}
	}
	data = dto.CreatePenyusunanKpiResponse{
		IdPengajuan:  req.IdPengajuan,
		ApprovalList: simpleList,
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

	// Ambil header dari DB berdasarkan id_pengajuan
	dbTahun, dbTriwulan, dbKostl, kostlTx, dbStatus, dbStatusDesc, err := s.repo.GetKpiHeader(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Validasi: status harus 1
	if dbStatus != 1 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: kostl, triwulan, tahun dari request harus sesuai dengan data di DB
	if req.Kostl != dbKostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, dbKostl),
		}
	}
	if req.Triwulan != dbTriwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, dbTriwulan),
		}
	}
	if req.Tahun != dbTahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, dbTahun),
		}
	}

	req.Divisi = dto.Divisi{Kostl: req.Kostl, KostlTx: kostlTx}

	// Parse dan validasi file Excel
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

	// Build ContextList, ProcessList, ResultList dari kolom P–U Excel
	// Hanya diisi untuk TW2 dan TW4
	resultList := []dto.PenyusunanResult{}
	processList := []dto.PenyusunanProcess{}
	contextList := []dto.PenyusunanContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	// Simpan ke DB: DELETE lama + INSERT baru + UPDATE header
	if err := s.repo.RevisionPenyusunanKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		processList,
		contextList,
	); err != nil {
		return data, err
	}

	data = dto.RevisionPenyusunanKpiResponse{
		IDPengajuan: req.IdPengajuan,
		Tahun:       req.Tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.Divisi{
			Kostl:   req.Divisi.Kostl,
			KostlTx: req.Divisi.KostlTx,
		},
		Entry: dto.EntryUser{
			EntryUser: req.EntryUser,
			EntryName: req.EntryName,
			EntryTime: req.EntryTime,
		},
		TotalKpi:    len(kpiRows),
		Kpi:         utils.BuildKpiResponse(req.IdPengajuan, kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
	}

	return data, nil
}

// =============================================================================
// APPROVE
// =============================================================================

func (s *penyusunanKpiService) ApprovePenyusunanKpi(
	req *dto.ApprovePenyusunanKpiRequest,
) (data dto.ApprovePenyusunanKpiResponse, err error) {
	dbTahun, dbTriwulan, dbKostl, _, _, _, err := s.repo.GetKpiHeader(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}
	if req.Kostl != dbKostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, dbKostl),
		}
	}
	if req.Tahun != dbTahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, dbTahun),
		}
	}
	if req.Triwulan != dbTriwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, dbTriwulan),
		}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.User)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUser
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	currentIdx := -1
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.User) && approvalList[i].Status == "" {
			approvalList[i].Status = "approve"
			approvalList[i].Keterangan = req.Catatan
			approvalList[i].Waktu = now
			currentIdx = i
			break
		}
	}
	if currentIdx == -1 {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	// Cari approver berikutnya yang belum approve
	nextApprover := ""
	for i := currentIdx + 1; i < len(approvalList); i++ {
		if approvalList[i].Status == "" {
			nextApprover = approvalList[i].Userid
			break
		}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	if err = s.repo.ApprovePenyusunanKpi(req.IdPengajuan, string(updatedJSON), nextApprover, req.User); err != nil {
		return data, err
	}

	data = dto.ApprovePenyusunanKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "approve",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// REJECT
// =============================================================================

func (s *penyusunanKpiService) RejectPenyusunanKpi(
	req *dto.RejectPenyusunanKpiRequest,
) (data dto.RejectPenyusunanKpiResponse, err error) {
	dbTahun, dbTriwulan, dbKostl, _, _, _, err := s.repo.GetKpiHeader(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}
	if req.Kostl != dbKostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, dbKostl),
		}
	}
	if req.Tahun != dbTahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, dbTahun),
		}
	}
	if req.Triwulan != dbTriwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, dbTriwulan),
		}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.User)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUser
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	found := false
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.User) && approvalList[i].Status == "" {
			approvalList[i].Status = "reject"
			approvalList[i].Keterangan = req.Catatan
			approvalList[i].Waktu = now
			found = true
			break
		}
	}
	if !found {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	if err = s.repo.RejectPenyusunanKpi(req.IdPengajuan, string(updatedJSON), req.Catatan, req.User); err != nil {
		return data, err
	}

	data = dto.RejectPenyusunanKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "reject",
		Catatan:     req.Catatan,
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
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
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
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
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
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
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
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
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
) (data *dto.GetDetailPenyusunanKpiResponse, err error) {
	return s.repo.GetDetailPenyusunanKpi(req)
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

	return file_export.GenerateKpiExcel(exportData)
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

	return file_export.GenerateKpiPDF(exportData)
}
