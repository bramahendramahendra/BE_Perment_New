package service

import (
	"encoding/json"
	"fmt"
	"time"

	dto "permen_api/domain/validasi_kpi/dto"
	customErrors "permen_api/errors"
)

// =============================================================================
// INPUT VALIDASI (validate + create + revision)
// =============================================================================

// InputValidasi digunakan oleh endpoint POST /validasi-kpi/input.
func (s *validasiKpiService) InputValidasi(
	req *dto.InputValidasiRequest,
) (data dto.InputValidasiResponse, err error) {
	existData, err := s.repo.GetExistDataValidasi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau status tidak mengizinkan input validasi", req.IdPengajuan),
		}
	}

	// Status yang diizinkan: 5 (baru), 7 (revisi setelah tolak), 90/91 (ulang setelah batal)
	if existData.Status != 5 && existData.Status != 7 && existData.Status != 90 && existData.Status != 91 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat diinput, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	if err := s.repo.InputValidasi(req); err != nil {
		return data, err
	}

	return dto.InputValidasiResponse{IdPengajuan: req.IdPengajuan}, nil
}

// =============================================================================
// APPROVE VALIDASI
// =============================================================================

// ApproveValidasi digunakan oleh endpoint POST /validasi-kpi/approve.
func (s *validasiKpiService) ApproveValidasi(
	req *dto.ApproveValidasiRequest,
) (data dto.ApproveValidasiResponse, err error) {
	existData, err := s.repo.GetExistDataValidasi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	if existData.Status != 6 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat diapprove, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	approvalExists, err := s.repo.CheckApprovalValidasiExists(req.ApprovalUserValidasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !approvalExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan. [Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListValidasiJSON(req.IdPengajuan, req.ApprovalUserValidasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalListItem
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list_validasi: %w", err)
	}

	approvalList, nextApprover, err := processApproveValidasiApprovalList(approvalList, req.ApprovalUserValidasi, req.Catatan.EntryNote)
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

	return dto.ApproveValidasiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Approve Validasi",
		Catatan:     req.Catatan,
	}, nil
}

// =============================================================================
// REJECT VALIDASI
// =============================================================================

// RejectValidasi digunakan oleh endpoint POST /validasi-kpi/reject.
func (s *validasiKpiService) RejectValidasi(
	req *dto.RejectValidasiRequest,
) (data dto.RejectValidasiResponse, err error) {
	existData, err := s.repo.GetExistDataValidasi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	if existData.Status != 6 {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("pengajuan '%s' tidak dapat ditolak, status saat ini '%s'", req.IdPengajuan, existData.StatusDesc),
		}
	}

	rejectExists, err := s.repo.CheckApprovalValidasiExists(req.ApprovalUserValidasi, req.IdPengajuan)
	if err != nil {
		return data, err
	}
	if !rejectExists {
		return data, &customErrors.BadRequestError{Message: "Data tidak ditemukan. [Pastikan User Approval sesuai.]"}
	}

	approvalListJSON, err := s.repo.GetApprovalListValidasiJSON(req.IdPengajuan, req.ApprovalUserValidasi)
	if err != nil {
		return data, err
	}

	var approvalList []dto.ApprovalListItem
	if err = json.Unmarshal([]byte(approvalListJSON), &approvalList); err != nil {
		return data, fmt.Errorf("gagal parse approval_list_validasi: %w", err)
	}

	approvalList, err = processRejectValidasiApprovalList(approvalList, req.ApprovalUserValidasi, req.Catatan.EntryNote)
	if err != nil {
		return data, &customErrors.BadRequestError{Message: err.Error()}
	}

	updatedJSON, err := json.Marshal(approvalList)
	if err != nil {
		return data, fmt.Errorf("gagal serialize approval_list_validasi: %w", err)
	}

	if err = s.repo.RejectValidasiKpi(req.IdPengajuan, string(updatedJSON), req.Catatan.EntryNote, req.ApprovalUserValidasi); err != nil {
		return data, err
	}

	return dto.RejectValidasiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "Reject Validasi",
		Catatan:     req.Catatan,
	}, nil
}

// =============================================================================
// VALIDASI BATAL
// =============================================================================

// ValidasiBatal digunakan oleh endpoint POST /validasi-kpi/batal.
func (s *validasiKpiService) ValidasiBatal(
	req *dto.ValidasiBatalRequest,
) (data dto.ValidasiBatalResponse, err error) {
	_, err = s.repo.GetExistDataValidasi(req.IdPengajuan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	if err := s.repo.ValidasiBatal(req); err != nil {
		return data, err
	}

	return dto.ValidasiBatalResponse{IdPengajuan: req.IdPengajuan}, nil
}

// =============================================================================
// GET ALL
// =============================================================================

// GetAllApprovalValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-approval.
func (s *validasiKpiService) GetAllApprovalValidasi(
	req *dto.GetAllApprovalValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllApprovalValidasi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiResponse{
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

// GetAllTolakanValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-tolakan.
func (s *validasiKpiService) GetAllTolakanValidasi(
	req *dto.GetAllTolakanValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllTolakanValidasi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiResponse{
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

// GetAllDaftarPenyusunanValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-penyusunan.
func (s *validasiKpiService) GetAllDaftarPenyusunanValidasi(
	req *dto.GetAllDaftarPenyusunanValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarPenyusunanValidasi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiResponse{
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

// GetAllDaftarApprovalValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-daftar-approval.
func (s *validasiKpiService) GetAllDaftarApprovalValidasi(
	req *dto.GetAllDaftarApprovalValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarApprovalValidasi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiResponse{
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

// GetAllValidasi digunakan oleh endpoint POST /validasi-kpi/get-all-validasi.
func (s *validasiKpiService) GetAllValidasi(
	req *dto.GetAllValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllValidasi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllValidasiResponse{
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

	// Parse approval_list_validasi (JSON → []ApprovalListItem)
	var approvalListValidasi []dto.ApprovalListItem
	if dataDB.ApprovalListValidasi != "" && dataDB.ApprovalListValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalListValidasi), &approvalListValidasi); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_validasi: %w", err)
		}
	}
	if approvalListValidasi == nil {
		approvalListValidasi = []dto.ApprovalListItem{}
	}

	// Parse lampiran_validasi (JSON → []string)
	var lampiranValidasi []string
	if dataDB.LampiranValidasi != "" && dataDB.LampiranValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.LampiranValidasi), &lampiranValidasi); err != nil {
			return nil, fmt.Errorf("gagal parse lampiran_validasi: %w", err)
		}
	}
	if lampiranValidasi == nil {
		lampiranValidasi = []string{}
	}

	// Parse qualifier_overall_validasi (JSON → []QualifierOverallItem)
	var qualifierOverall []dto.QualifierOverallItem
	if dataDB.QualifierOverallValidasi != "" && dataDB.QualifierOverallValidasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.QualifierOverallValidasi), &qualifierOverall); err != nil {
			return nil, fmt.Errorf("gagal parse qualifier_overall_validasi: %w", err)
		}
	}
	if qualifierOverall == nil {
		qualifierOverall = []dto.QualifierOverallItem{}
	}

	// Build KPI list dari nested model
	kpiList := make([]dto.DataKpiDetailValidasiResponse, len(dataDB.Kpi))
	for i, kpi := range dataDB.Kpi {
		subList := make([]dto.DataKpiSubDetailValidasiResponse, len(kpi.KpiSubDetail))
		for j, sub := range kpi.KpiSubDetail {
			subList[j] = dto.DataKpiSubDetailValidasiResponse{
				IdSubDetail:                      sub.IdSubDetail,
				IdSubKpi:                         sub.IdKpi,
				SubKpi:                           sub.Kpi,
				Bobot:                            sub.Bobot,
				TargetTriwulan:                   sub.TargetTriwulan,
				TargetKuantitatifTriwulan:        sub.TargetKuantitatifTriwulan,
				TargetQualifier:                  sub.TargetQualifier,
				RealisasiValidated:               sub.RealisasiValidated,
				RealisasiKuantitatifValidated:    sub.RealisasiKuantitatifValidated,
				ValidasiKeterangan:               sub.ValidasiKeterangan,
				Pencapaian:                       sub.Pencapaian,
				Skor:                             sub.Skor,
				PencapaianQualifierValidated:     sub.PencapaianQualifierValidated,
				PencapaianPostQualifierValidated: sub.PencapaianPostQualifierValidated,
			}
		}
		kpiList[i] = dto.DataKpiDetailValidasiResponse{
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
		CatatanTolakan:           dataDB.CatatanTolakan,
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
// HELPER — approval list processing
// =============================================================================

// processApproveValidasiApprovalList menandai entry approver sebagai "approve",
// lalu mengembalikan nextApprover (userid approver berikutnya, "" jika sudah final).
func processApproveValidasiApprovalList(
	list []dto.ApprovalListItem,
	userID, catatan string,
) (updated []dto.ApprovalListItem, nextApprover string, err error) {
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

	// Cari approver berikutnya (entry pertama yang status-nya masih kosong)
	for _, item := range list {
		if item.Status == "" {
			nextApprover = item.Userid
			break
		}
	}

	return list, nextApprover, nil
}

// processRejectValidasiApprovalList menandai entry approver sebagai "reject".
func processRejectValidasiApprovalList(
	list []dto.ApprovalListItem,
	userID, catatan string,
) (updated []dto.ApprovalListItem, err error) {
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
