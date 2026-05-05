package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	dto "permen_api/domain/realisasi_kpi/dto"
	"permen_api/domain/realisasi_kpi/utils"
	customErrors "permen_api/errors"
)

// =============================================================================
// VALIDATE
// =============================================================================

// ValidateRealisasiKpi digunakan oleh endpoint POST /realisasi-kpi/validate.
func (s *realisasiKpiService) ValidateRealisasiKpi(
	req *dto.ValidateRealisasiKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidateRealisasiKpiResponse, err error) {

	if err := utils.ValidateExcelFile(file); err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 2 = Penyusunan Disetujui
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

	if err := s.resolveRealisasiLookups(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
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

	data = dto.ValidateRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Triwulan:    req.Triwulan,
		Tahun:       req.Tahun,
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
// CREATE
// =============================================================================

func (s *realisasiKpiService) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) (data dto.CreateRealisasiKpiResponse, err error) {
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbEntryUserRealisasi := existData.EntryUserRealisasi
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 80 = Draft Realisasi
	if dbStatus != 80 {
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

	ApprovalListRealisasi := make([]dto.ApprovalUser, len(req.ApprovalListRealisasi))
	for i, a := range req.ApprovalListRealisasi {
		ApprovalListRealisasi[i] = dto.ApprovalUser{Userid: a.Userid, Nama: a.Nama}
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

	if err := utils.ValidateExcelFile(file); err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbEntryUserRealisasi := existData.EntryUserRealisasi
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 4 = Realisasi Ditolak
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

	// Parse dan validasi file Excel
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateRealisasiExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	if err := s.resolveRealisasiLookups(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
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

	data = dto.RevisionRealisasiKpiResponse{
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
// APPROVE
// =============================================================================

func (s *realisasiKpiService) ApproveRealisasiKpi(
	req *dto.ApproveRealisasiKpiRequest,
) (data dto.ApproveRealisasiKpiResponse, err error) {
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 3 = Approval Realisasi
	if dbStatus != 3 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat diapprove, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
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

	approvalExists, err := s.repo.CheckApprovalRealisasiExists(req.ApprovalUserRealisasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan.[Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserRealisasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	approvalList, nextApprover, err := utils.ProcessApproveApprovalList(approvalList, req.ApprovalUserRealisasi, req.Catatan.EntryNote)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
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
		Status:      "Approve Realisasi",
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
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 3 = Approval Realisasi
	if dbStatus != 3 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat ditolak, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: kostl, triwulan, tahun dari request harus sesuai dengan data di DB
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
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan.[Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserRealisasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	approvalList, err = utils.ProcessRejectApprovalList(approvalList, req.ApprovalUserRealisasi, req.Catatan.EntryNote)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	existingCatatanJSON, err := s.repo.GetCatatanTolakan(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal membaca catatan_tolakan: %w", err)
	}

	nowDisplay := time.Now().Format("02-01-2006 15:04:05")
	entryUserFull := req.ApprovalUserRealisasi + " - " + req.ApprovalNameRealisasi

	catatanTolakanJSON, err := utils.AppendCatatanTolakan(existingCatatanJSON, dto.CatatanDetail{
		Fungsi:    req.Catatan.Fungsi,
		EntryUser: entryUserFull,
		EntryTime: nowDisplay,
		EntryNote: req.Catatan.EntryNote,
	})
	if err != nil {
		return data, err
	}

	if err = s.repo.RejectRealisasiKpi(req.IdPengajuan, string(updatedJSON), catatanTolakanJSON, req.ApprovalUserRealisasi); err != nil {
		return data, err
	}

	data = dto.RejectRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Reject Realisasi",
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
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

	var approvalList []dto.ApprovalUserDetail
	if dataDB.ApprovalList != "" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalList), &approvalList); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list: %w", err)
		}
	}
	if approvalList == nil {
		approvalList = []dto.ApprovalUserDetail{}
	}

	var approvalListRealisasi []dto.ApprovalUserDetail
	if dataDB.ApprovalListRealisasi != "" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalListRealisasi), &approvalListRealisasi); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_realisasi: %w", err)
		}
	}
	if approvalListRealisasi == nil {
		approvalListRealisasi = []dto.ApprovalUserDetail{}
	}

	var catatanList []dto.CatatanDetail
	if dataDB.CatatanTolakan != "" && dataDB.CatatanTolakan != "null" {
		if err = json.Unmarshal([]byte(dataDB.CatatanTolakan), &catatanList); err != nil {
			return nil, fmt.Errorf("gagal parse catatan_tolakan: %w", err)
		}
	}
	if catatanList == nil {
		catatanList = []dto.CatatanDetail{}
	}

	listKpiDetails := make([]dto.DataKpiDetail, len(dataDB.Kpi))
	for i, v := range dataDB.Kpi {
		listKpiSubDetails := make([]dto.DataKpiSubdetail, len(v.KpiSubDetail))
		for j, s := range v.KpiSubDetail {
			listKpiSubDetails[j] = dto.DataKpiSubdetail{
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
				RealisasiQualifier:            s.RealisasiQualifier,
				RealisasiKuantitatifQualifier: s.RealisasiKuantitatifQualifier,
				RealisasiKeterangan:           s.RealisasiKeterangan,
				RealisasiValidated:            s.RealisasiValidated,
				RealisasiKuantitatifValidated: s.RealisasiKuantitatifValidated,
				Pencapaian:                    s.Pencapaian,
				Skor:                          s.Skor,
			}
		}
		listKpiDetails[i] = dto.DataKpiDetail{
			IdDetail:            v.IdDetail,
			IdKpi:               v.IdKpi,
			Kpi:                 v.Kpi,
			Rumus:               v.Rumus,
			IdPerspektif:        v.IdPersfektif,
			Persfektif:          v.Perspektif,
			IdKeteranganProject: v.IdKeteranganProject,
			KeteranganProject:   v.KeteranganProject,
			LinkDokumenSumber:   v.LampiranFile,
			TotalSubKpi:         v.TotalSubKpi,
			KpiSubDetail:        listKpiSubDetails,
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
		Triwulan:    dataDB.Triwulan,
		Tahun:       dataDB.Tahun,
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
		KpiList:               listKpiDetails,
		TotalResult:           dataDB.TotalResult,
		ResultList:            resultList,
		TotalProcess:          dataDB.TotalProcess,
		ProcessList:           processList,
		TotalContext:          dataDB.TotalContext,
		ContextList:           contextList,
	}

	return data, nil
}
