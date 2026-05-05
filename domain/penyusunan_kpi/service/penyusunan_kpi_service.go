package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
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

// ValidatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/validate.
func (s *penyusunanKpiService) ValidatePenyusunanKpi(
	req *dto.ValidatePenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidatePenyusunanKpiResponse, err error) {

	if err := utils.ValidateExcelFile(file); err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Cek data KPI untuk tahun/triwulan/kostl sudah ada
	idLama, statusLama, found, err := s.repo.GetExistDataKpiStatus(req.Tahun, req.Triwulan, req.Divisi.Kostl)
	if err != nil {
		return data, err
	}
	if found {
		// Status 70 (Penyusunan Draft) dan 71 (Penyusunan Batal) → replace pengajuan penyusunan KPI
		if statusLama != 70 && statusLama != 71 {
			return data, &customErrors.BadRequestError{
				Message: fmt.Sprintf(
					"data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
					req.Tahun, req.Triwulan, req.Divisi.Kostl,
				),
			}
		}
		// idLama akan digunakan di repo untuk menghapus draft lama sebelum insert baru
	} else {
		idLama = ""
	}

	// Parse dan validasi file Excel.
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	if err := s.resolvePenyusunanLookups(kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Bangun idPengajuan di service agar bisa digunakan untuk build response sebelum repo insert.
	idPengajuan := idgen.GenerateIDPengajuan(req.Divisi.Kostl, req.Tahun, req.Triwulan)

	// Build ProcessList dan ContextList dari data Excel (kolom P,Q,R,S,T,U).
	// Hanya diisi untuk TW2 dan TW4; untuk TW1 dan TW3 list kosong (tidak diinsert ke DB).
	resultList := []dto.DataResult{}
	processList := []dto.DataProcess{}
	contextList := []dto.DataContext{}
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
		idLama,
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
		EntryPenyusunan: dto.EntryUserPenyusunan{
			EntryUserPenyusunan: req.EntryUserPenyusunan,
			EntryNamePenyusunan: req.EntryNamePenyusunan,
			EntryTimePenyusunan: req.EntryTimePenyusunan,
		},
		TotalKpi:    len(kpiRows),
		KpiList:     utils.BuildKpiResponse(idPengajuan, kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
	}

	return data, nil
}

// =============================================================================
// CREATE
// =============================================================================

// CreatePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/create.
func (s *penyusunanKpiService) CreatePenyusunanKpi(
	req *dto.CreatePenyusunanKpiRequest,
) (data dto.CreatePenyusunanKpiResponse, err error) {
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	dbTriwulan := existData.Triwulan
	dbTahun := existData.Tahun
	dbKostl := existData.Kostl
	dbKostlTx := existData.KostlTx
	dbEntryUserPenyusunan := existData.EntryUser
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 70 = Penyusunan Draft
	if dbStatus != 70 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: hanya pembuat pengajuan yang boleh merevisi
	if req.EntryUserPenyusunan != dbEntryUserPenyusunan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("user '%s' tidak berhak merevisi pengajuan ini", req.EntryUserPenyusunan),
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

	if err = s.repo.CreatePenyusunanKpi(req); err != nil {
		return data, err
	}

	approvalListPenyusunan := make([]dto.ApprovalUser, len(req.ApprovalListPenyusunan))
	for i, a := range req.ApprovalListPenyusunan {
		approvalListPenyusunan[i] = dto.ApprovalUser{Userid: a.Userid, Nama: a.Nama, Posisi: a.Posisi}
	}
	data = dto.CreatePenyusunanKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Divisi: dto.Divisi{
			Kostl:   divisi.Kostl,
			KostlTx: divisi.KostlTx,
		},
		Tahun:                  req.Tahun,
		Triwulan:               req.Triwulan,
		ApprovalListPenyusunan: approvalListPenyusunan,
	}

	return data, nil
}

// =============================================================================
// REVISION
// =============================================================================

// RevisionPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/revision.
func (s *penyusunanKpiService) RevisionPenyusunanKpi(
	req *dto.RevisionPenyusunanKpiRequest,
	file *multipart.FileHeader,
) (data dto.RevisionPenyusunanKpiResponse, err error) {

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
	dbEntryUserPenyusunan := existData.EntryUser
	dbStatus := existData.Status
	dbStatusDesc := existData.StatusDesc

	// Validasi: status harus 1 = Penyusunan Ditolak
	if dbStatus != 1 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: hanya pembuat pengajuan yang boleh merevisi
	if req.EntryUserPenyusunan != dbEntryUserPenyusunan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("user '%s' tidak berhak merevisi pengajuan ini", req.EntryUserPenyusunan),
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
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	if err := s.resolvePenyusunanLookups(kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Build ContextList, ProcessList, ResultList (hanya TW2 dan TW4)
	resultList := []dto.DataResult{}
	processList := []dto.DataProcess{}
	contextList := []dto.DataContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(req.IdPengajuan, req.Tahun, req.Triwulan, kpiRows, kpiSubDetails)
	}

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
		Triwulan:    req.Triwulan,
		Tahun:       req.Tahun,
		Divisi: dto.Divisi{
			Kostl:   divisi.Kostl,
			KostlTx: divisi.KostlTx,
		},
		EntryPenyusunan: dto.EntryUserPenyusunan{
			EntryUserPenyusunan: req.EntryUserPenyusunan,
			EntryNamePenyusunan: req.EntryNamePenyusunan,
			EntryTimePenyusunan: req.EntryTimePenyusunan,
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
// APPROVAL
// =============================================================================

// ApprovePenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/approve.
func (s *penyusunanKpiService) ApprovePenyusunanKpi(
	req *dto.ApprovePenyusunanKpiRequest,
) (data dto.ApprovePenyusunanKpiResponse, err error) {
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

	// Validasi: status harus 0 = Approval Penyusunan
	if dbStatus != 0 {
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

	approvalExists, err := s.repo.CheckApprovalPenyusunanExists(req.ApprovalUserPenyusunan, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan.[Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserPenyusunan)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	approvalList, nextApprover, err := utils.ProcessApproveApprovalList(approvalList, req.ApprovalUserPenyusunan, req.Catatan.EntryNote)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list: %w", err)
	}

	if err = s.repo.ApprovePenyusunanKpi(req.IdPengajuan, string(updatedJSON), nextApprover, req.ApprovalUserPenyusunan); err != nil {
		return data, err
	}

	data = dto.ApprovePenyusunanKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Approve Penyusunan",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// RejectPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/reject.
func (s *penyusunanKpiService) RejectPenyusunanKpi(
	req *dto.RejectPenyusunanKpiRequest,
) (data dto.RejectPenyusunanKpiResponse, err error) {
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

	// Validasi: status harus 0 = Approval Penyusunan
	if dbStatus != 0 {
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

	rejectApprovalExists, err := s.repo.CheckApprovalPenyusunanExists(req.ApprovalUserPenyusunan, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !rejectApprovalExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan.[User Approval Kosong.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserPenyusunan)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	approvalList, err = utils.ProcessRejectApprovalList(approvalList, req.ApprovalUserPenyusunan, req.Catatan.EntryNote)
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
	entryUserFull := req.ApprovalUserPenyusunan + " - " + req.ApprovalNamePenyusunan

	catatanTolakanJSON, err := utils.AppendCatatanTolakan(existingCatatanJSON, dto.CatatanDetail{
		Fungsi:    req.Catatan.Fungsi,
		EntryUser: entryUserFull,
		EntryTime: nowDisplay,
		EntryNote: req.Catatan.EntryNote,
	})
	if err != nil {
		return data, err
	}

	if err = s.repo.RejectPenyusunanKpi(req.IdPengajuan, string(updatedJSON), catatanTolakanJSON, req.ApprovalUserPenyusunan); err != nil {
		return data, err
	}

	data = dto.RejectPenyusunanKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Reject Penyusunan",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-approval.
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// GetAllTolakanPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-tolakan.
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// GetAllDaftarPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-penyusunan.
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
			Triwulan:    v.Triwulan,
			Tahun:       v.Tahun,
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
			StatusDesc:  v.StatusDesc,
		})
	}

	return data, total, nil
}

// GetAllDaftarApprovalPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-all-daftar-approval.
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

// GetDetailPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-detail.
func (s *penyusunanKpiService) GetDetailPenyusunanKpi(
	req *dto.GetDetailPenyusunanKpiRequest,
) (data *dto.GetDetailPenyusunanKpiResponse, err error) {
	dataDB, err := s.repo.GetDetailPenyusunanKpi(req)
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
				IdSubDetail:               s.IdSubDetail,
				IdSubKpi:                  s.IdKpi,
				SubKpi:                    s.Kpi,
				Otomatis:                  s.Otomatis,
				IdPolarisasi:              s.IdPolarisasi,
				Polarisasi:                s.Polarisasi,
				Capping:                   s.Capping,
				Bobot:                     s.Bobot,
				Glossary:                  s.DeskripsiGlossary,
				TargetTriwulan:            s.TargetTriwulan,
				TargetKuantitatifTriwulan: s.TargetKuantitatifTriwulan,
				TargetTahunan:             s.TargetTahunan,
				TargetKuantitatifTahunan:  s.TargetKuantitatifTahunan,
				TerdapatQualifier:         s.IdQualifier,
				Qualifier:                 s.ItemQualifier,
				DeskripsiQualifier:        s.DeskripsiQualifier,
				TargetQualifier:           s.TargetQualifier,
				IdKeteranganProject:       s.IdKeteranganProject,
				KeteranganProject:         s.KeteranganProject,
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
			TotalSubKpi:         v.TotalSubKpi,
			KpiSubDetail:        listKpiSubDetails,
		}
	}

	resultList := make([]dto.DataResult, len(dataDB.ResultList))
	for i, v := range dataDB.ResultList {
		resultList[i] = dto.DataResult{
			IdDetailResult:  v.IdDetailResult,
			NamaResult:      v.NamaResult,
			DeskripsiResult: v.DeskripsiResult,
		}
	}

	processList := make([]dto.DataProcess, len(dataDB.ProcessList))
	for i, v := range dataDB.ProcessList {
		processList[i] = dto.DataProcess{
			IdDetailProcess:  v.IdDetailMethod,
			NamaProcess:      v.NamaMethod,
			DeskripsiProcess: v.DeskripsiMethod,
		}
	}

	contextList := make([]dto.DataContext, len(dataDB.ContextList))
	for i, v := range dataDB.ContextList {
		contextList[i] = dto.DataContext{
			IdDetailContext:  v.IdDetailChallenge,
			NamaContext:      v.NamaChallenge,
			DeskripsiContext: v.DeskripsiChallenge,
		}
	}

	data = &dto.GetDetailPenyusunanKpiResponse{
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
		ApprovalPosisi: dataDB.ApprovalPosisi,
		ApprovalList:   approvalList,
		Catatan:        catatanList,
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
		TotalKpi:     dataDB.TotalKpi,
		KpiList:      listKpiDetails,
		TotalResult:  dataDB.TotalResult,
		ResultList:   resultList,
		TotalProcess: dataDB.TotalProcess,
		ProcessList:  processList,
		TotalContext: dataDB.TotalContext,
		ContextList:  contextList,
	}

	return data, nil
}

// =============================================================================
// GET EXPORT DATA
// =============================================================================

// GetExcelPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-excel.
func (s *penyusunanKpiService) GetExcelPenyusunanKpi(
	req *dto.GetExcelPenyusunanKpiRequest,
) ([]byte, string, error) {
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}
	status := existData.Status
	statusDesc := existData.StatusDesc
	if status == 0 || status == 1 {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("File tidak dapat diunduh, status pengajuan: %s", statusDesc),
		}
	}

	exportData, err := s.repo.GetKpiExportData(req.IdPengajuan, req.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return nil, "", err
	}
	if exportData.NamaDivisi == "" {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	return file_export.GenerateKpiExcel(exportData)
}

// GetPdfPenyusunanKpi digunakan oleh endpoint POST /penyusunan-kpi/get-pdf.
func (s *penyusunanKpiService) GetPdfPenyusunanKpi(
	req *dto.GetPdfPenyusunanKpiRequest,
) ([]byte, string, error) {
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}
	status := existData.Status
	statusDesc := existData.StatusDesc
	if status == 0 || status == 1 {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("File tidak dapat diunduh, status pengajuan: %s", statusDesc),
		}
	}

	exportData, err := s.repo.GetKpiExportData(req.IdPengajuan, req.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return nil, "", err
	}
	if exportData.NamaDivisi == "" {
		return nil, "", &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	return file_export.GenerateKpiPDF(exportData)
}
