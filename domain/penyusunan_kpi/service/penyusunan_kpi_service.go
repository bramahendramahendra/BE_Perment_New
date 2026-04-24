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

	// Cek data KPI untuk tahun/triwulan/kostl sudah ada
	idLama, statusLama, found, err := s.repo.GetExistPenyusunanStatus(req.Tahun, req.Triwulan, req.Kostl)
	if err != nil {
		return data, err
	}
	if found {
		// Status 70 = draft, 71 = revisi → boleh di-replace, selain itu ditolak sebagai duplikasi
		if statusLama != 70 && statusLama != 71 {
			return data, &customErrors.BadRequestError{
				Message: fmt.Sprintf(
					"data KPI untuk tahun %s, triwulan %s, kostl %s sudah ada",
					req.Tahun, req.Triwulan, req.Kostl,
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
	exists, err := s.repo.CheckExistIdPengajuan(req.IdPengajuan, req.Kostl, req.Tahun, req.Triwulan)
	if err != nil {
		return data, err
	}
	if !exists {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' dengan kostl '%s', tahun '%s', triwulan '%s' tidak ditemukan", req.IdPengajuan, req.Kostl, req.Tahun, req.Triwulan),
		}
	}

	_, _, _, kostlTx, _, _, _, _, err := s.repo.GetKpiHeader(req.IdPengajuan)
	if err != nil {
		return data, err
	}

	if err = s.repo.CreatePenyusunanKpi(req); err != nil {
		return data, err
	}

	approvalList := make([]dto.ApprovalUser, len(req.ApprovalList))
	for i, a := range req.ApprovalList {
		approvalList[i] = dto.ApprovalUser{Userid: a.Userid, Nama: a.Nama, Posisi: a.Posisi}
	}
	data = dto.CreatePenyusunanKpiResponse{
		IdPengajuan:  req.IdPengajuan,
		Divisi:       kostlTx,
		Tahun:        req.Tahun,
		Triwulan:     req.Triwulan,
		ApprovalList: approvalList,
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
	dbTahun, dbTriwulan, dbKostl, kostlTx, dbEntryUser, _, dbStatus, dbStatusDesc, err := s.repo.GetKpiHeader(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Validasi: status harus 1
	if dbStatus != 1 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, dbStatusDesc),
		}
	}

	// Validasi: hanya pembuat pengajuan yang boleh merevisi
	if req.EntryUser != dbEntryUser {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("user '%s' tidak berhak merevisi pengajuan ini", req.EntryUser),
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
	resultList := []dto.DataResult{}
	processList := []dto.DataProcess{}
	contextList := []dto.DataContext{}
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

	approvalExists, err := s.repo.CheckApprovalExists(req.ApprovalUser, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUser)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	keterangan := req.Catatan.EntryNote
	currentIdx := -1
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.ApprovalUser) && approvalList[i].Status == "" {
			approvalList[i].Status = "approve"
			approvalList[i].Keterangan = keterangan
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

	if err = s.repo.ApprovePenyusunanKpi(req.IdPengajuan, string(updatedJSON), nextApprover, req.ApprovalUser); err != nil {
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

	rejectApprovalExists, err := s.repo.CheckApprovalExists(req.ApprovalUser, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !rejectApprovalExists {
		return data, &customErrors.BadRequestError{Message: "Data Not Found"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUser)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list: %w", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	nowDisplay := time.Now().Format("02-01-2006 15:04:05")
	entryUserFull := req.ApprovalUser + " - " + req.ApprovalName

	keterangan := req.Catatan.EntryNote
	found := false
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, req.ApprovalUser) && approvalList[i].Status == "" {
			approvalList[i].Status = "reject"
			approvalList[i].Keterangan = keterangan
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

	existingCatatanJSON, err := s.repo.GetCatatanTolakan(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal membaca catatan_tolakan: %w", err)
	}

	var catatanTolakanEntries []dto.CatatanTolakanEntry
	if existingCatatanJSON != "" && existingCatatanJSON != "null" {
		if err = json.Unmarshal([]byte(existingCatatanJSON), &catatanTolakanEntries); err != nil {
			return data, fmt.Errorf("gagal parse catatan_tolakan: %w", err)
		}
	}

	catatanTolakanEntries = append(catatanTolakanEntries, dto.CatatanTolakanEntry{
		Fungsi:    req.Catatan.Fungsi,
		EntryUser: entryUserFull,
		EntryTime: nowDisplay,
		EntryNote: req.Catatan.EntryNote,
	})

	catatanTolakanJSON, err := json.Marshal(catatanTolakanEntries)
	if err != nil {
		return data, fmt.Errorf("gagal serialize catatan_tolakan: %w", err)
	}

	if err = s.repo.RejectPenyusunanKpi(req.IdPengajuan, string(updatedJSON), string(catatanTolakanJSON), req.ApprovalUser); err != nil {
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

	var catatanList []dto.CatatanTolakanEntry
	if dataDB.CatatanTolakan != "" && dataDB.CatatanTolakan != "null" {
		if err = json.Unmarshal([]byte(dataDB.CatatanTolakan), &catatanList); err != nil {
			return nil, fmt.Errorf("gagal parse catatan_tolakan: %w", err)
		}
	}
	if catatanList == nil {
		catatanList = []dto.CatatanTolakanEntry{}
	}

	kpiDetails := make([]dto.DataKpiDetail, len(dataDB.Kpi))
	for i, v := range dataDB.Kpi {
		subDetails := make([]dto.DataKpiSubdetail, len(v.KpiSubDetail))
		for j, s := range v.KpiSubDetail {
			subDetails[j] = dto.DataKpiSubdetail{
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
		kpiDetails[i] = dto.DataKpiDetail{
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
		Entry: dto.EntryUser{
			EntryUser: dataDB.EntryUser,
			EntryName: dataDB.EntryName,
			EntryTime: dataDB.EntryTime,
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
		Kpi:          kpiDetails,
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
// GET EXCEL
// =============================================================================

func (s *penyusunanKpiService) GetExcelPenyusunanKpi(
	req *dto.GetExcelPenyusunanKpiRequest,
) ([]byte, string, error) {
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

// =============================================================================
// GET PDF
// =============================================================================

func (s *penyusunanKpiService) GetPdfPenyusunanKpi(
	req *dto.GetPdfPenyusunanKpiRequest,
) ([]byte, string, error) {
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
