package service

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	"mime/multipart"

	dto "permen_api/domain/realisasi_kpi/dto"
	"permen_api/domain/realisasi_kpi/utils"
	customErrors "permen_api/errors"
)

// =============================================================================
// VALIDATE
// =============================================================================

func (s *realisasiKpiService) ValidateRealisasiKpi(
	req *dto.ValidateRealisasiKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidateRealisasiKpiResponse, err error) {

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
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTahun := existData.Tahun
	dbTriwulan := existData.Triwulan
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 2
	if dbStatus != 2 {
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

	divisi := dto.Divisi{Kostl: req.Kostl, KostlTx: dbKostlTx}

	// Parse dan validasi file Excel
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateRealisasiExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup DB per baris: isi IdSubDetail, IdDetail, TargetKuantitatifTriwulan, Rumus,
	// lalu hitung Pencapaian dan Skor
	if err := s.enrichRowsFromDB(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Build resultList, processList, contextList (hanya TW2 dan TW4)
	resultList := []dto.DataResult{}
	processList := []dto.DataProcess{}
	contextList := []dto.DataContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	// Simpan ke DB (status 80 = draft realisasi)
	if err := s.repo.ValidateRealisasiKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		processList,
		contextList,
	); err != nil {
		return data, err
	}

	// Build response
	totalSubKpi := 0
	for _, kpiRow := range kpiRows {
		totalSubKpi += len(kpiSubDetails[kpiRow.KpiIndex])
	}

	data = dto.ValidateRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Tahun:       req.Tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.Divisi{
			Kostl:   divisi.Kostl,
			KostlTx: divisi.KostlTx,
		},
		EntryRealisasi: dto.EntryUserRealisasi{
			EntryUserRealisasi: req.EntryUserRealisasi,
			EntryNameRealisasi: req.EntryNameRealisasi,
			EntryTimeRealisasi: req.EntryTimeRealisasi,
		},
		TotalKpi:    len(kpiRows),
		KpiList:     utils.BuildKpiResponse(req.IdPengajuan, kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
	}

	return data, nil
}

// =============================================================================
// CREATE — submit ke approval
// =============================================================================

func (s *realisasiKpiService) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) (data dto.CreateRealisasiKpiResponse, err error) {
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTahun := existData.Tahun
	dbTriwulan := existData.Triwulan
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbEntryUserRealisasi := existData.EntryNameRealisasi
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 4
	if dbStatus != 4 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: hanya pembuat pengajuan yang boleh merevisi
	if req.EntryUserRealisasi != dbEntryUserRealisasi {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("user '%s' tidak berhak merevisi pengajuan ini", req.EntryUserRealisasi),
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

	divisi := dto.Divisi{Kostl: req.Kostl, KostlTx: dbKostlTx}

	if err := s.repo.CreateRealisasiKpi(req); err != nil {
		return data, err
	}

	ApprovalListRealisasi := make([]dto.ApprovalUserRealisasi, len(req.ApprovalListRealisasi))
	for i, a := range req.ApprovalListRealisasi {
		ApprovalListRealisasi[i] = dto.ApprovalUserRealisasi{Userid: a.Userid, Nama: a.Nama}
	}

	data = dto.CreateRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Divisi: dto.Divisi{
			Kostl:   divisi.Kostl,
			KostlTx: divisi.KostlTx,
		},
		Tahun:                 req.Tahun,
		Triwulan:              req.Triwulan,
		ApprovalListRealisasi: ApprovalListRealisasi,
	}

	return data, nil
}

// =============================================================================
// REVISION
// =============================================================================

func (s *realisasiKpiService) RevisionRealisasiKpi(
	req *dto.RevisionRealisasiKpiRequest,
	file *multipart.FileHeader,
) (data dto.RevisionRealisasiKpiResponse, err error) {

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
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTahun := existData.Tahun
	dbTriwulan := existData.Triwulan
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbEntryUser := existData.EntryUserRealisasi
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 3
	if dbStatus != 3 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: hanya pembuat pengajuan yang boleh merevisi
	if req.EntryUserRealisasi != dbEntryUser {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("user '%s' tidak berhak merevisi pengajuan ini", req.EntryUserRealisasi),
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

	divisi := dto.Divisi{Kostl: req.Kostl, KostlTx: dbKostlTx}

	kpiRows, kpiSubDetails, err := utils.ParseAndValidateRealisasiExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup DB per baris: isi IdSubDetail, IdDetail, TargetKuantitatifTriwulan, Rumus,
	// lalu hitung Pencapaian dan Skor
	if err := s.enrichRowsFromDB(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Ambil header (tahun) untuk response
	tahun, _, _, _, err := s.repo.GetKpiHeaderByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return data, err
	}

	// Build resultList, processList, contextList (hanya TW2 dan TW4)
	resultList := []dto.DataResult{}
	processList := []dto.DataProcess{}
	contextList := []dto.DataContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

	if err := s.repo.RevisionRealisasiKpi(
		req,
		kpiRows,
		kpiSubDetails,
		resultList,
		processList,
		contextList,
	); err != nil {
		return data, err
	}

	// Build response
	totalSubKpi := 0
	for _, kpiRow := range kpiRows {
		totalSubKpi += len(kpiSubDetails[kpiRow.KpiIndex])
	}

	data = dto.RevisionRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Tahun:       tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.Divisi{
			Kostl:   divisi.Kostl,
			KostlTx: divisi.KostlTx,
		},
		EntryRealisasi: dto.EntryUserRealisasi{
			EntryUserRealisasi: req.EntryUserRealisasi,
			EntryNameRealisasi: req.EntryNameRealisasi,
			EntryTimeRealisasi: req.EntryTimeRealisasi,
		},
		TotalKpi:    len(kpiRows),
		KpiList:     utils.BuildKpiResponse(req.IdPengajuan, kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
	}

	return data, nil
}

// =============================================================================
// APPROVE
// =============================================================================

func (s *realisasiKpiService) ApproveRealisasiKpi(
	req *dto.ApproveRealisasiKpiRequest,
) (data dto.ApproveRealisasiKpiResponse, err error) {
	dbTahun, dbTriwulan, dbKostl, _, _, _, _, _, err := s.repo.GetKpiHeader(req.IdPengajuan)
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

	approvalExists, err := s.repo.CheckApprovalRealisasiExists(req.ApprovalUserRealisasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserRealisasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserRealisasiDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	currentIdx := -1
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.ApprovalUserRealisasi) && approvalList[i].Status == "" {
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

	if err = s.repo.ApproveRealisasiKpi(req.IdPengajuan, string(updatedJSON), nextApprover, req.ApprovalUserRealisasi); err != nil {
		return data, err
	}

	data = dto.ApproveRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "approve",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// REJECT
// =============================================================================

func (s *realisasiKpiService) RejectRealisasiKpi(
	req *dto.RejectRealisasiKpiRequest,
) (data dto.RejectRealisasiKpiResponse, err error) {
	dbTahun, dbTriwulan, dbKostl, _, _, _, _, _, err := s.repo.GetKpiHeader(req.IdPengajuan)
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

	rejectApprovalExists, err := s.repo.CheckApprovalRealisasiExists(req.ApprovalUserRealisasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !rejectApprovalExists {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserRealisasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserRealisasiDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	found := false
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.ApprovalUserRealisasi) && approvalList[i].Status == "" {
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

	if err = s.repo.RejectRealisasiKpi(req.IdPengajuan, string(updatedJSON), req.Catatan, req.ApprovalUserRealisasi); err != nil {
		return data, err
	}

	data = dto.RejectRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "reject",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// GET ALL
// =============================================================================

func (s *realisasiKpiService) GetAllRealisasiKpi(
	req *dto.GetAllRealisasiKpiRequest,
) (data []*dto.GetAllRealisasiKpiResponse, total int64, err error) {

	dataDB, total, err := s.repo.GetAllRealisasiKpi(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllRealisasiKpiResponse{
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
// GET ALL APPROVAL
// =============================================================================

func (s *realisasiKpiService) GetAllApprovalRealisasiKpi(
	req *dto.GetAllApprovalRealisasiKpiRequest,
) (data []*dto.GetAllApprovalRealisasiKpiResponse, total int64, err error) {

	dataDB, total, err := s.repo.GetAllApprovalRealisasiKpi(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllApprovalRealisasiKpiResponse{
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

func (s *realisasiKpiService) GetAllTolakanRealisasiKpi(
	req *dto.GetAllTolakanRealisasiKpiRequest,
) (data []*dto.GetAllTolakanRealisasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllTolakanRealisasiKpi(req)
	if err != nil {
		return data, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllTolakanRealisasiKpiResponse{
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
// GET ALL DAFTAR REALISASI
// =============================================================================

func (s *realisasiKpiService) GetAllDaftarRealisasiKpi(
	req *dto.GetAllDaftarRealisasiKpiRequest,
) (data []*dto.GetAllDaftarRealisasiKpiResponse, total int64, err error) {

	dataDB, total, err := s.repo.GetAllDaftarRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarRealisasiKpiResponse{
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

func (s *realisasiKpiService) GetAllDaftarApprovalRealisasiKpi(
	req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
) (data []*dto.GetAllDaftarApprovalRealisasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarApprovalRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarApprovalRealisasiKpiResponse{
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

func (s *realisasiKpiService) GetDetailRealisasiKpi(
	req *dto.GetDetailRealisasiKpiRequest,
) (data *dto.GetDetailRealisasiKpiResponse, err error) {
	dataDB, err := s.repo.GetDetailRealisasiKpi(req)
	if err != nil {
		return nil, err
	}
	if dataDB.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	var approvalList []dto.ApprovalUserRealisasiDetail
	if dataDB.ApprovalList != "" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalList), &approvalList); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list: %w", err)
		}
	}
	if approvalList == nil {
		approvalList = []dto.ApprovalUserRealisasiDetail{}
	}

	var approvalListRealisasi []dto.ApprovalUserRealisasiDetail
	if dataDB.ApprovalListRealisasi != "" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalListRealisasi), &approvalListRealisasi); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_realisasi: %w", err)
		}
	}
	if approvalListRealisasi == nil {
		approvalListRealisasi = []dto.ApprovalUserRealisasiDetail{}
	}

	var catatanList []dto.CatatanTolakanEntry
	if dataDB.CatatanTolakan != "" && dataDB.CatatanTolakan != "null" {
		if err = json.Unmarshal([]byte(dataDB.CatatanTolakan), &catatanList); err != nil {
			return nil, fmt.Errorf("gagal parse catatan_tolakan: %w", err)
		}
	}
	if catatanList == nil {
		catatanList = []dto.CatatanTolakanEntry{}
	}

	kpiList := make([]dto.DataKpiDetail, len(dataDB.Kpi))
	for i, v := range dataDB.Kpi {
		subDetails := make([]dto.DataKpiSubdetail, len(v.KpiSubDetail))
		for j, s := range v.KpiSubDetail {
			subDetails[j] = dto.DataKpiSubdetail{
				IdSubDetail:                   s.IdSubDetail,
				IdSubKpi:                      s.IdKpi,
				SubKpi:                        s.Kpi,
				Otomatis:                      s.Otomatis,
				IdPolarisasi:                  s.IdPolarisasi,
				Polarisasi:                    s.Polarisasi,
				Capping:                       s.Capping,
				Bobot:                         s.Bobot,
				Glossary:                      s.DeskripsiGlossary,
				TargetTriwulan:                s.TargetTriwulan,
				TargetKuantitatifTriwulan:     s.TargetKuantitatifTriwulan,
				TargetTahunan:                 s.TargetTahunan,
				TargetKuantitatifTahunan:      s.TargetKuantitatifTahunan,
				TerdapatQualifier:             s.IdQualifier,
				Qualifier:                     s.ItemQualifier,
				DeskripsiQualifier:            s.DeskripsiQualifier,
				TargetQualifier:               s.TargetQualifier,
				IdKeteranganProject:           s.IdKeteranganProject,
				KeteranganProject:             s.KeteranganProject,
				Realisasi:                     s.Realisasi,
				RealisasiKuantitatif:          s.RealisasiKuantitatif,
				RealisasiKeterangan:           s.RealisasiKeterangan,
				RealisasiValidated:            s.RealisasiValidated,
				RealisasiKuantitatifValidated: s.RealisasiKuantitatifValidated,
				Pencapaian:                    s.Pencapaian,
				Skor:                          s.Skor,
				RealisasiQualifier:            s.RealisasiQualifier,
				RealisasiKuantitatifQualifier: s.RealisasiKuantitatifQualifier,
			}
		}
		kpiList[i] = dto.DataKpiDetail{
			IdDetail:            v.IdDetail,
			IdKpi:               v.IdKpi,
			Kpi:                 v.Kpi,
			Rumus:               v.Rumus,
			IdPerspektif:        v.IdPersfektif,
			Persfektif:          v.Perspektif,
			IdKeteranganProject: v.IdKeteranganProject,
			KeteranganProject:   v.KeteranganProject,
			TotalSubKpi:         v.TotalSubKpi,
			KpiSubDetail:        subDetails,
		}
	}

	resultList := make([]dto.DataResult, len(dataDB.ResultList))
	for i, v := range dataDB.ResultList {
		resultList[i] = dto.DataResult{
			IdDetailResult:   v.IdDetailResult,
			NamaResult:       v.NamaResult,
			DeskripsiResult:  v.DeskripsiResult,
			RealisasiResult:  v.RealisasiResult,
			LampiranEvidence: v.LampiranEvidence,
		}
	}

	processList := make([]dto.DataProcess, len(dataDB.ProcessList))
	for i, v := range dataDB.ProcessList {
		processList[i] = dto.DataProcess{
			IdDetailProcess:  v.IdDetailMethod,
			NamaProcess:      v.NamaMethod,
			DeskripsiProcess: v.DeskripsiMethod,
			RealisasiProcess: v.RealisasiMethod,
			LampiranEvidence: v.LampiranEvidence,
		}
	}

	contextList := make([]dto.DataContext, len(dataDB.ContextList))
	for i, v := range dataDB.ContextList {
		contextList[i] = dto.DataContext{
			IdDetailContext:  v.IdDetailChallenge,
			NamaContext:      v.NamaChallenge,
			DeskripsiContext: v.DeskripsiChallenge,
			RealisasiContext: v.RealisasiChallenge,
			LampiranEvidence: v.LampiranEvidence,
		}
	}

	data = &dto.GetDetailRealisasiKpiResponse{
		IdPengajuan: dataDB.IdPengajuan,
		Tahun:       dataDB.Tahun,
		Triwulan:    dataDB.Triwulan,
		Status:      dataDB.Status,
		StatusDesc:  dataDB.StatusDesc,
		Divisi: dto.DivisiOrgeh{
			Kostl:   dataDB.Kostl,
			KostlTx: dataDB.KostlTx,
			Orgeh:   dataDB.Orgeh,
			OrgehTx: dataDB.OrgehTx,
		},
		EntryPenyusunan: dto.EntryUserPenyusunan{
			EntryUserPenyusunan: dataDB.EntryUser,
			EntryNamePenyusunan: dataDB.EntryName,
			EntryTimePenyusunan: dataDB.EntryTime,
		},
		EntryRealisasi: dto.EntryUserRealisasi{
			EntryUserRealisasi: dataDB.EntryUserRealisasi,
			EntryNameRealisasi: dataDB.EntryNameRealisasi,
			EntryTimeRealisasi: dataDB.EntryTimeRealisasi,
		},
		EntryValidasi: dto.EntryUserValidasi{
			EntryUserValidasi: dataDB.EntryUserValidasi,
			EntryNameValidasi: dataDB.EntryNameValidasi,
			EntryTimeValidasi: dataDB.EntryTimeValidasi,
		},
		ApprovalPosisi:        dataDB.ApprovalPosisi,
		ApprovalListRealisasi: approvalListRealisasi,
		Catatan:               catatanList,
		TotalBobot:            dataDB.TotalBobot,
		TotalPencapaian:       dataDB.TotalPencapaian,
		TotalKpi:              dataDB.TotalKpi,
		KpiList:               kpiList,
		TotalResult:           dataDB.TotalResult,
		ResultList:            resultList,
		TotalProcess:          dataDB.TotalProcess,
		ProcessList:           processList,
		TotalContext:          dataDB.TotalContext,
		ContextList:           contextList,
	}

	return data, nil
}

// =============================================================================
// PRIVATE HELPERS
// =============================================================================

// enrichRowsFromDB melakukan lookup ke DB untuk setiap baris sub KPI Excel:
//   - Mencari id_sub_detail, id_detail, target_kuantitatif_triwulan, rumus
//     berdasarkan id_pengajuan + kpi_name + sub_kpi_name
//   - Menghitung Pencapaian dan Skor mengikuti logika bisnis BE_Perment_Old
//
// Logika kalkulasi (dari models/M_realisasi.go):
//
//	Rumus "1" (Maximize): Pencapaian = (RealisasiKuantitatif / TargetKuantitatif) * 100
//	Rumus "0" (Minimize): Pencapaian = (TargetKuantitatif / RealisasiKuantitatif) * 100
//	Capping diterapkan jika Pencapaian > nilai Capping
//	Skor = (Pencapaian * Bobot) / 100
func (s *realisasiKpiService) enrichRowsFromDB(
	idPengajuan string,
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) error {
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for i := range subRows {
			sub := &kpiSubDetails[kpiRow.KpiIndex][i]

			idSubDetail, idDetail, rumus, targetKuantitatif, err :=
				s.repo.LookupSubDetailByKpiSubKpi(idPengajuan, kpiRow.Kpi, sub.SubKPI)
			if err != nil {
				return &customErrors.BadRequestError{Message: err.Error()}
			}

			sub.IdSubDetail = idSubDetail
			sub.IdDetail = idDetail
			sub.Rumus = rumus
			sub.TargetKuantitatifTriwulan = targetKuantitatif

			// Hitung Pencapaian dan Skor
			pencapaian, skor := calculatePencapaianSkor(
				rumus,
				sub.RealisasiKuantitatif,
				targetKuantitatif,
				sub.Capping,
				sub.Bobot,
			)

			sub.Pencapaian = pencapaian
			sub.Skor = skor
		}
	}

	return nil
}

// calculatePencapaianSkor menghitung Pencapaian (%) dan Skor dari nilai realisasi.
//
// Mengikuti logika bisnis models/M_realisasi.go dari BE_Perment_Old secara persis:
//
//	rumus == "1" → Maximize : Pencapaian = (realisasi / target) * 100
//	rumus == "0" → Minimize : Pencapaian = (target / realisasi) * 100
//	Capping diterapkan jika Pencapaian melebihi batas ("100%" = 100, "110%" = 110)
//	Skor = (Pencapaian * bobot) / 100
//	Jika target = 0 atau rumus tidak dikenal → Pencapaian = 0, Skor = 0
func calculatePencapaianSkor(
	rumus string,
	realisasiKuantitatif float64,
	targetKuantitatif float64,
	cappingStr string,
	bobot float64,
) (pencapaian, skor float64) {

	// Parse nilai capping numerik dari string "100%" atau "110%"
	cappingValue := parseCapping(cappingStr)

	switch rumus {
	case "1": // Maximize
		if targetKuantitatif == 0 {
			return 0, 0
		}
		pencapaian = (realisasiKuantitatif / targetKuantitatif) * 100

	case "0": // Minimize
		if realisasiKuantitatif == 0 {
			return 0, 0
		}
		pencapaian = (targetKuantitatif / realisasiKuantitatif) * 100

	default:
		// Rumus tidak dikenal (misal "0" untuk sub KPI lain/non-standar)
		return 0, 0
	}

	// Terapkan capping
	if cappingValue > 0 && pencapaian > cappingValue {
		pencapaian = cappingValue
	}

	skor = (pencapaian * bobot) / 100

	// Bulatkan 2 desimal agar konsisten
	pencapaian = math.Round(pencapaian*100) / 100
	skor = math.Round(skor*100) / 100

	return pencapaian, skor
}

// parseCapping mengubah string capping ("100%" atau "110%") menjadi nilai float64.
// Mengembalikan 0 jika format tidak dikenal (tidak ada capping yang diterapkan).
func parseCapping(cappingStr string) float64 {
	switch strings.TrimSpace(cappingStr) {
	case "100%":
		return 100.0
	case "110%":
		return 110.0
	default:
		return 0
	}
}
