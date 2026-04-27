package service

import (
	"fmt"

	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"
	customErrors "permen_api/errors"
)

// =============================================================================
// INPUT VALIDASI
// =============================================================================

func (s *validasiKpiService) InputValidasi(req *dto.InputValidasiRequest) (data dto.InputValidasiResponse, err error) {
	exists, err := s.repo.CheckExistInputValidasi(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal memeriksa status pengajuan: %w", err)
	}
	if !exists {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau status tidak mengizinkan input validasi", req.IdPengajuan),
		}
	}

	if err := s.repo.InputValidasi(req); err != nil {
		return data, err
	}

	return dto.InputValidasiResponse{IdPengajuan: req.IdPengajuan}, nil
}

// =============================================================================
// APPROVAL VALIDASI
// =============================================================================

func (s *validasiKpiService) ApprovalValidasi(req *dto.ApprovalValidasiRequest) (data dto.ApprovalValidasiResponse, err error) {
	exists, err := s.repo.CheckExistApprovalValidasi(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal memeriksa status pengajuan: %w", err)
	}
	if !exists {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau bukan dalam status pending approval validasi", req.IdPengajuan),
		}
	}

	if err := s.repo.ApprovalValidasi(req); err != nil {
		return data, err
	}

	return dto.ApprovalValidasiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      req.Status,
	}, nil
}

// =============================================================================
// VALIDASI BATAL
// =============================================================================

func (s *validasiKpiService) ValidasiBatal(req *dto.ValidasiBatalRequest) (data dto.ValidasiBatalResponse, err error) {
	exists, err := s.repo.CheckExistBatalValidasi(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal memeriksa pengajuan: %w", err)
	}
	if !exists {
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
// APPROVE VALIDASI
// =============================================================================

func (s *validasiKpiService) ApproveValidasi(req *dto.ApproveValidasiRequest) (data dto.ApproveValidasiResponse, err error) {
	exists, err := s.repo.CheckExistApprovalValidasi(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal memeriksa status pengajuan: %w", err)
	}
	if !exists {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau bukan dalam status pending approval validasi", req.IdPengajuan),
		}
	}

	if err := s.repo.ApproveValidasi(req); err != nil {
		return data, err
	}

	return dto.ApproveValidasiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "approve",
	}, nil
}

// =============================================================================
// REJECT VALIDASI
// =============================================================================

func (s *validasiKpiService) RejectValidasi(req *dto.RejectValidasiRequest) (data dto.RejectValidasiResponse, err error) {
	exists, err := s.repo.CheckExistApprovalValidasi(req.IdPengajuan)
	if err != nil {
		return data, fmt.Errorf("gagal memeriksa status pengajuan: %w", err)
	}
	if !exists {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan atau bukan dalam status pending approval validasi", req.IdPengajuan),
		}
	}

	if err := s.repo.RejectValidasi(req); err != nil {
		return data, err
	}

	return dto.RejectValidasiResponse{
		IdPengajuan: req.IdPengajuan,
		Status:      "reject",
	}, nil
}

// =============================================================================
// GET ALL — helper mapping
// =============================================================================

func mapToGetAllValidasiResponse(dataDB []*model.DataKpi) []*dto.GetAllValidasiResponse {
	var result []*dto.GetAllValidasiResponse
	for _, v := range dataDB {
		result = append(result, &dto.GetAllValidasiResponse{
			IdPengajuan: v.IdPengajuan,
			Tahun:       v.Tahun,
			Triwulan:    v.Triwulan,
			KostlTx:     v.KostlTx,
			OrgehTx:     v.OrgehTx,
			StatusDesc:  v.StatusDesc,
		})
	}
	return result
}

// =============================================================================
// GET ALL APPROVAL VALIDASI
// =============================================================================

func (s *validasiKpiService) GetAllApprovalValidasi(
	req *dto.GetAllApprovalValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllApprovalValidasi(req)
	if err != nil {
		return nil, 0, err
	}
	return mapToGetAllValidasiResponse(dataDB), total, nil
}

// =============================================================================
// GET ALL TOLAKAN VALIDASI
// =============================================================================

func (s *validasiKpiService) GetAllTolakanValidasi(
	req *dto.GetAllTolakanValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllTolakanValidasi(req)
	if err != nil {
		return nil, 0, err
	}
	return mapToGetAllValidasiResponse(dataDB), total, nil
}

// =============================================================================
// GET ALL DAFTAR PENYUSUNAN VALIDASI
// =============================================================================

func (s *validasiKpiService) GetAllDaftarPenyusunanValidasi(
	req *dto.GetAllDaftarPenyusunanValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarPenyusunanValidasi(req)
	if err != nil {
		return nil, 0, err
	}
	return mapToGetAllValidasiResponse(dataDB), total, nil
}

// =============================================================================
// GET ALL DAFTAR APPROVAL VALIDASI
// =============================================================================

func (s *validasiKpiService) GetAllDaftarApprovalValidasi(
	req *dto.GetAllDaftarApprovalValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllDaftarApprovalValidasi(req)
	if err != nil {
		return nil, 0, err
	}
	return mapToGetAllValidasiResponse(dataDB), total, nil
}

// =============================================================================
// GET ALL VALIDASI
// =============================================================================

func (s *validasiKpiService) GetAllValidasi(
	req *dto.GetAllValidasiRequest,
) (data []*dto.GetAllValidasiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllValidasi(req)
	if err != nil {
		return nil, 0, err
	}
	return mapToGetAllValidasiResponse(dataDB), total, nil
}
