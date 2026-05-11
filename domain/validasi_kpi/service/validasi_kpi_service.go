package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto "permen_api/domain/validasi_kpi/dto"
	"permen_api/domain/validasi_kpi/utils"
	customErrors "permen_api/errors"
	file_export "permen_api/pkg/file_export"
)

// =============================================================================
// INPUT
// =============================================================================

// InputValidasiKpi digunakan oleh endpoint POST /validasi-kpi/input.
func (s *validasiKpiService) InputValidasiKpi(
	req *dto.InputValidasiKpiRequest,
) (data dto.InputValidasiKpiResponse, err error) {
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	// Validasi: status harus 5 = Realisasi Disetujui / 7 = Validasi Ditolak / 90 = Draft Validasi / 91 = Validasi Batal
	if existData.Status != 5 && existData.Status != 7 && existData.Status != 90 && existData.Status != 91 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat direvisi, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	// Validasi: kostl, triwulan, tahun dari request harus sesuai dengan data di DB
	if req.Kostl != existData.Kostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, existData.Kostl),
		}
	}
	if req.Triwulan != existData.Triwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, existData.Triwulan),
		}
	}
	if req.Tahun != existData.Tahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, existData.Tahun),
		}
	}

	if err := s.repo.InputValidasiKpi(req); err != nil {
		return data, err
	}

	approvalList := make([]dto.ApprovalUser, len(req.ApprovalListValidasi))
	for i, a := range req.ApprovalListValidasi {
		approvalList[i] = dto.ApprovalUser{Userid: a.Userid, Nama: a.Nama, Posisi: a.Posisi}
	}

	data = dto.InputValidasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Triwulan:    req.Triwulan,
		Tahun:       req.Tahun,
		Divisi: dto.Divisi{
			Kostl:   req.Kostl,
			KostlTx: existData.KostlTx,
		},
		EntryValidasi: dto.EntryUserValidasi{
			EntryUserValidasi: req.EntryUserValidasi,
			EntryNameValidasi: req.EntryNameValidasi,
			EntryTimeValidasi: req.EntryTimeValidasi,
		},
		ApprovalListValidasi:         approvalList,
		TotalBobot:                   req.TotalBobot,
		TotalPencapaian:              req.TotalPencapaian,
		TotalBobotPengurang:          req.TotalBobotPengurang,
		TotalPencapaianPost:          req.TotalPencapaianPost,
		KpiList:                      req.Kpi,
		DataValidasiQualifierOverall: req.DataValidasiQualifierOverall,
		LampiranValidasi:             req.LampiranValidasi,
	}

	return data, nil
}

// =============================================================================
// APPROVE
// =============================================================================

// ApproveValidasiKpi digunakan oleh endpoint POST /validasi-kpi/approve.
func (s *validasiKpiService) ApproveValidasiKpi(
	req *dto.ApproveValidasiKpiRequest,
) (data dto.ApproveValidasiKpiResponse, err error) {
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}

	// Validasi: status harus 6 = Approval Validasi
	if existData.Status != 6 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat diapprove, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	// Validasi: kostl, triwulan, tahun dari request harus sesuai dengan data di DB
	if req.Kostl != existData.Kostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, existData.Kostl),
		}
	}
	if req.Triwulan != existData.Triwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, existData.Triwulan),
		}
	}
	if req.Tahun != existData.Tahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, existData.Tahun),
		}
	}

	approvalExists, err := s.repo.CheckApprovalValidasiExists(req.ApprovalUserValidasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan. [Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserValidasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list_validasi: %w", err)
	}

	approvalList, nextApprover, err := utils.ProcessApproveApprovalList(approvalList, req.ApprovalUserValidasi, req.Catatan.EntryNote)

	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list_validasi: %w", err)
	}

	if err = s.repo.ApproveValidasiKpi(req.IdPengajuan, string(updatedJSON), nextApprover, req.ApprovalUserValidasi); err != nil {
		return data, err
	}

	data = dto.ApproveValidasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Approve Validasi",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// REJECT VALIDASI
// =============================================================================

// RejectValidasiKpi digunakan oleh endpoint POST /validasi-kpi/reject.
func (s *validasiKpiService) RejectValidasiKpi(
	req *dto.RejectValidasiKpiRequest,
) (data dto.RejectValidasiKpiResponse, err error) {
	// Ambil header dari DB berdasarkan id_pengajuan
	existData, err := s.repo.GetExistDataKpi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan)}
	}

	// Validasi: status harus 6 = Approval Validasi
	if existData.Status != 6 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat ditolak, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	// Validasi: kostl, triwulan, tahun dari request harus sesuai dengan data di DB
	if req.Kostl != existData.Kostl {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("kostl '%s' tidak sesuai dengan data pengajuan (kostl: '%s')", req.Kostl, existData.Kostl),
		}
	}
	if req.Tahun != existData.Tahun {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("tahun '%s' tidak sesuai dengan data pengajuan (tahun: '%s')", req.Tahun, existData.Tahun),
		}
	}
	if req.Triwulan != existData.Triwulan {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("triwulan '%s' tidak sesuai dengan data pengajuan (triwulan: '%s')", req.Triwulan, existData.Triwulan),
		}
	}

	rejectExists, err := s.repo.CheckApprovalValidasiExists(req.ApprovalUserValidasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !rejectExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan. [Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListJSON(req.IdPengajuan, req.ApprovalUserValidasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalUserDetail
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval list validasi: %w", err)
	}

	approvalList, err = utils.ProcessRejectApprovalList(approvalList, req.ApprovalUserValidasi, req.Catatan.EntryNote)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval list validasi: %w", err)
	}

	existingCatatanJSON, err := s.repo.GetCatatanTolakan(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal membaca catatan_tolakan: %w", err)
	}

	nowDisplay := time.Now().Format("02-01-2006 15:04:05")
	entryUserFull := req.ApprovalUserValidasi + " - " + req.ApprovalNameValidasi

	catatanTolakanJSON, err := utils.AppendCatatanTolakan(existingCatatanJSON, dto.CatatanDetail{
		Fungsi:    req.Catatan.Fungsi,
		EntryUser: entryUserFull,
		EntryTime: nowDisplay,
		EntryNote: req.Catatan.EntryNote,
	})
	if err != nil {
		return data, err
	}

	if err = s.repo.RejectValidasiKpi(req.IdPengajuan, string(updatedJSON), catatanTolakanJSON, req.ApprovalUserValidasi); err != nil {
		return data, err
	}

	data = dto.RejectValidasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Reject Validasi",
		Catatan:     req.Catatan,
	}

	return data, nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-validasi.
func (s *validasiKpiService) GetAllValidasiKpi(
	req *dto.GetAllValidasiKpiRequest,
) (data []*dto.GetAllValidasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllValidasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiKpiResponse{
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

// GetAllApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
func (s *validasiKpiService) GetAllApprovalValidasiKpi(
	req *dto.GetAllApprovalValidasiKpiRequest,
) (data []*dto.GetAllApprovalValidasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllApprovalValidasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllApprovalValidasiKpiResponse{
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

// GetAllTolakanValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
func (s *validasiKpiService) GetAllTolakanValidasiKpi(
	req *dto.GetAllTolakanValidasiKpiRequest,
) (data []*dto.GetAllTolakanValidasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllTolakanValidasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllTolakanValidasiKpiResponse{
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

// GetAllDaftarValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-validasi.
func (s *validasiKpiService) GetAllDaftarValidasiKpi(
	req *dto.GetAllDaftarPValidasiKpiRequest,
) (data []*dto.GetAllDaftarValidasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarValidasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarValidasiKpiResponse{
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

// GetAllDaftarApprovalValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
func (s *validasiKpiService) GetAllDaftarApprovalValidasiKpi(
	req *dto.GetAllDaftarApprovalValidasiKpiRequest,
) (data []*dto.GetAllDaftarApprovalValidasiKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarApprovalValidasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllDaftarApprovalValidasiKpiResponse{
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

// GetDetailValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-detail.
func (s *validasiKpiService) GetDetailValidasiKpi(
	req *dto.GetDetailValidasiKpiRequest,
) (data *dto.GetDetailValidasiKpiResponse, err error) {
	dataDB, err := s.repo.GetDetailValidasiKpi(req)
	if err != nil {
		return nil, err
	}
	if dataDB.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	var approvalListValidasi []dto.ApprovalUserDetail
	if dataDB.ApprovalListValidasi != "" && dataDB.ApprovalListValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalListValidasi), &approvalListValidasi); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_validasi: %w", err)
		}
	}
	if approvalListValidasi == nil {
		approvalListValidasi = []dto.ApprovalUserDetail{}
	}

	var lampiranValidasi []string
	if dataDB.LampiranValidasi != "" && dataDB.LampiranValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.LampiranValidasi), &lampiranValidasi); err != nil {
			return nil, fmt.Errorf("gagal parse lampiran_validasi: %w", err)
		}
	}
	if lampiranValidasi == nil {
		lampiranValidasi = []string{}
	}

	var qualifierOverall []dto.DataValidasiQualifierOverall
	if dataDB.QualifierOverallValidasi != "" && dataDB.QualifierOverallValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.QualifierOverallValidasi), &qualifierOverall); err != nil {
			return nil, fmt.Errorf("gagal parse qualifier_overall_validasi: %w", err)
		}
	}
	if qualifierOverall == nil {
		qualifierOverall = []dto.DataValidasiQualifierOverall{}
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

	kpiList := make([]dto.DataKpiDetail, len(dataDB.Kpi))
	for i, kpi := range dataDB.Kpi {
		subList := make([]dto.DataKpiSubdetail, len(kpi.KpiSubDetail))
		for j, sub := range kpi.KpiSubDetail {
			subList[j] = dto.DataKpiSubdetail{
				IdSubDetail:                      sub.IdSubDetail,
				IdSubKpi:                         sub.IdKpi,
				SubKpi:                           sub.Kpi,
				Bobot:                            sub.Bobot,
				TargetTriwulan:                   sub.TargetTriwulan,
				TargetKuantitatifTriwulan:        sub.TargetKuantitatifTriwulan,
				TargetQualifier:                  sub.TargetQualifier,
				RealisasiValidated:               sub.RealisasiValidated,
				RealisasiKuantitatifValidated:    sub.RealisasiKuantitatifValidated,
				IdSumber:                         sub.IdSumber,
				Sumber:                           sub.Sumber,
				ValidasiKeterangan:               sub.ValidasiKeterangan,
				Pencapaian:                       sub.Pencapaian,
				Skor:                             sub.Skor,
				PencapaianQualifierValidated:     sub.PencapaianQualifierValidated,
				PencapaianPostQualifierValidated: sub.PencapaianPostQualifierValidated,
			}
		}
		kpiList[i] = dto.DataKpiDetail{
			IdDetail:     kpi.IdDetail,
			IdKpi:        kpi.IdKpi,
			Kpi:          kpi.Kpi,
			Rumus:        kpi.Rumus,
			TotalSubKpi:  kpi.TotalSubKpi,
			KpiSubDetail: subList,
		}
	}

	data = &dto.GetDetailValidasiKpiResponse{
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
		ApprovalPosisi:           dataDB.ApprovalPosisi,
		ApprovalListValidasi:     approvalListValidasi,
		Catatan:                  catatanList,
		TotalBobot:               dataDB.TotalBobot,
		TotalPencapaian:          dataDB.TotalPencapaian,
		TotalBobotPengurang:      dataDB.TotalBobotPengurang,
		TotalPencapaianPost:      dataDB.TotalPencapaianPost,
		LampiranValidasi:         lampiranValidasi,
		TotalKpi:                 dataDB.TotalKpi,
		KpiList:                  kpiList,
		QualifierOverallValidasi: qualifierOverall,
	}

	return data, nil
}

// =============================================================================
// DOWNLOAD
// =============================================================================

// GetExcelValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-excel.
func (s *validasiKpiService) GetExcelValidasiKpi(
	req *dto.GetExcelValidasiKpiRequest,
) ([]byte, string, error) {
	exportData, err := s.buildValidasiKpiExportData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}
	return file_export.GenerateValidasiKpiExcel(exportData)
}

// GetPdfValidasiKpi digunakan oleh endpoint POST /validasi-kpi/get-pdf.
func (s *validasiKpiService) GetPdfValidasiKpi(
	req *dto.GetPdfValidasiKpiRequest,
) ([]byte, string, error) {
	exportData, err := s.buildValidasiKpiExportData(req.IdPengajuan)
	if err != nil {
		return nil, "", err
	}
	return file_export.GenerateValidasiKpiPDF(exportData)
}

// buildValidasiKpiExportData mengambil data dari DB dan mengubahnya ke ValidasiKpiExportData.
func (s *validasiKpiService) buildValidasiKpiExportData(idPengajuan string) (*dto.ValidasiKpiExportData, error) {
	dataDB, err := s.repo.GetDetailValidasiKpi(&dto.GetDetailValidasiKpiRequest{IdPengajuan: idPengajuan})
	if err != nil {
		return nil, err
	}
	if dataDB.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
		}
	}

	twNum := strings.TrimPrefix(dataDB.Triwulan, "TW")

	// Status 90 = Draft Validasi
	isDraft := dataDB.Status == 90

	var rows []dto.ValidasiKpiExportRow
	no := 1
	for _, kpi := range dataDB.Kpi {
		for _, sub := range kpi.KpiSubDetail {
			noQualifier := strings.EqualFold(sub.IdQualifier, "TIDAK")

			itemQualifier := sub.ItemQualifier
			targetQualifier := sub.TargetQualifier
			realisasiQualifier := sub.RealisasiQualifier
			pencapaianQualifier := fmt.Sprintf("%.2f%%", sub.PencapaianQualifierValidated)
			pencapaianPost := fmt.Sprintf("%.2f%%", sub.PencapaianPostQualifierValidated)

			if noQualifier {
				itemQualifier = "-"
				targetQualifier = "-"
				realisasiQualifier = "-"
				pencapaianQualifier = "-"
				pencapaianPost = "-"
			}

			rows = append(rows, dto.ValidasiKpiExportRow{
				No:                      no,
				Kpi:                     sub.Kpi,
				ItemQualifier:           itemQualifier,
				Bobot:                   sub.Bobot,
				TargetTriwulan:          sub.TargetTriwulan,
				TargetQualifier:         targetQualifier,
				RealisasiValidated:      sub.RealisasiValidated,
				RealisasiQualifier:      realisasiQualifier,
				Pencapaian:              fmt.Sprintf("%.2f%%", sub.Pencapaian),
				PencapaianQualifier:     pencapaianQualifier,
				PencapaianPostQualifier: pencapaianPost,
			})
			no++
		}
	}

	if rows == nil {
		rows = []dto.ValidasiKpiExportRow{}
	}

	return &dto.ValidasiKpiExportData{
		NamaDivisi:      dataDB.KostlTx,
		Triwulan:        dataDB.Triwulan,
		TriwulanNum:     twNum,
		Tahun:           dataDB.Tahun,
		TotalPencapaian: dataDB.TotalPencapaian,
		IsDraft:         isDraft,
		Rows:            rows,
	}, nil
}

// =============================================================================
// HELPER — approval list processing
// =============================================================================

func processApproveValidasiApprovalList(
	list []dto.ApprovalUserDetail,
	userID, catatan string,
) (updated []dto.ApprovalUserDetail, nextApprover string, err error) {
	now := time.Now().Format("02-01-2006 15:04:05")
	found := false

	for i := range list {
		if list[i].Userid == userID && list[i].Status == "" {
			list[i].Status = "approve"
			list[i].Keterangan = catatan
			list[i].Waktu = now
			found = true
			break
		}
	}

	if !found {
		return nil, "", fmt.Errorf("user '%s' tidak ditemukan dalam daftar approval atau sudah memproses", userID)
	}

	for _, item := range list {
		if item.Status == "" {
			nextApprover = item.Userid
			break
		}
	}

	return list, nextApprover, nil
}

func processRejectValidasiApprovalList(
	list []dto.ApprovalUserDetail,
	userID, catatan string,
) (updated []dto.ApprovalUserDetail, err error) {
	now := time.Now().Format("02-01-2006 15:04:05")
	found := false

	for i := range list {
		if list[i].Userid == userID && list[i].Status == "" {
			list[i].Status = "reject"
			list[i].Keterangan = catatan
			list[i].Waktu = now
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("user '%s' tidak ditemukan dalam daftar approval atau sudah memproses", userID)
	}

	return list, nil
}
