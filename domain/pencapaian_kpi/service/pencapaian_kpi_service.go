package service

import (
	"encoding/json"
	"fmt"

	dto "permen_api/domain/pencapaian_kpi/dto"
	"permen_api/domain/pencapaian_kpi/utils"
	customErrors "permen_api/errors"
	file_export "permen_api/pkg/file_export"
)

// =============================================================================
// GET ALL
// =============================================================================

// GetAllPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-all-pencapaian.
func (s *pencapaianKpiService) GetAllPencapaianKpi(
	req *dto.GetAllPencapaianKpiRequest,
) (data []*dto.GetAllPencapaianKpiResponse, total int64, err error) {
	dataDB, total, err := s.repo.GetAllPencapaianKpi(req)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range dataDB {
		data = append(data, &dto.GetAllPencapaianKpiResponse{
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

// GetDetailPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-detail.
func (s *pencapaianKpiService) GetDetailPencapaianKpi(
	req *dto.GetDetailPencapaianKpiRequest,
) (data *dto.GetDetailPencapaianKpiResponse, err error) {
	indikatorDB, err := s.repo.GetIndikatorPencapaian()
	if err != nil {
		return nil, err
	}

	dataDB, err := s.repo.GetDetailPencapaianKpi(req)
	if err != nil {
		return nil, err
	}
	if dataDB.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", req.IdPengajuan),
		}
	}

	var approvalListPenyusunan []dto.ApprovalUserDetail
	if dataDB.ApprovalList != "" && dataDB.ApprovalList != "null" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalList), &approvalListPenyusunan); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_penyusunan: %w", err)
		}
	}
	if approvalListPenyusunan == nil {
		approvalListPenyusunan = []dto.ApprovalUserDetail{}
	}

	var approvalListRealisasi []dto.ApprovalUserDetail
	if dataDB.ApprovalListRealisasi != "" && dataDB.ApprovalListRealisasi != "null" {
		if err = json.Unmarshal([]byte(dataDB.ApprovalListRealisasi), &approvalListRealisasi); err != nil {
			return nil, fmt.Errorf("gagal parse approval_list_realisasi: %w", err)
		}
	}
	if approvalListRealisasi == nil {
		approvalListRealisasi = []dto.ApprovalUserDetail{}
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
	if dataDB.QualifierOverallValidasi != "" && dataDB.QualifierOverallValidasi != "null" && dataDB.QualifierOverallValidasi != "-" {
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
				Otomatis:                         sub.Otomatis,
				Polarisasi:                       sub.Polarisasi,
				IdPolarisasi:                     sub.IdPolarisasi,
				Capping:                          sub.Capping,
				Bobot:                            sub.Bobot,
				Glossary:                         sub.DeskripsiGlossary,
				TargetTriwulan:                   sub.TargetTriwulan,
				TargetKuantitatifTriwulan:        sub.TargetKuantitatifTriwulan,
				TargetTahunan:                    sub.TargetTahunan,
				TargetKuantitatifTahunan:         sub.TargetKuantitatifTahunan,
				TerdapatQualifier:                sub.IdQualifier,
				Qualifier:                        sub.ItemQualifier,
				DeskripsiQualifier:               sub.DeskripsiQualifier,
				TargetQualifier:                  sub.TargetQualifier,
				IdKeteranganProject:              sub.IdKeteranganProject,
				KeteranganProject:                sub.KeteranganProject,
				Realisasi:                        sub.Realisasi,
				RealisasiKuantitatif:             sub.RealisasiKuantitatif,
				RealisasiQualifier:               sub.RealisasiQualifier,
				RealisasiKuantitatifQualifier:    sub.RealisasiKuantitatifQualifier,
				RealisasiKeterangan:              sub.RealisasiKeterangan,
				RealisasiValidated:               sub.RealisasiValidated,
				RealisasiKuantitatifValidated:    sub.RealisasiKuantitatifValidated,
				IdSumber:                         sub.IdSumber,
				Sumber:                           sub.Sumber,
				ValidasiKeterangan:               sub.ValidasiKeterangan,
				Pencapaian:                       sub.Pencapaian,
				IndikatorPencapaian:              utils.ResolveIndikator(sub.Pencapaian, indikatorDB),
				Skor:                             sub.Skor,
				PencapaianQualifierValidated:     sub.PencapaianQualifierValidated,
				IndikatorPencapaianQualifier:     utils.ResolveIndikator(sub.PencapaianQualifierValidated, indikatorDB),
				PencapaianPostQualifierValidated: sub.PencapaianPostQualifierValidated,
				IndikatorPencapaianPostQualifier: utils.ResolveIndikator(sub.PencapaianPostQualifierValidated, indikatorDB),
			}
		}
		kpiList[i] = dto.DataKpiDetail{
			IdDetail:            kpi.IdDetail,
			IdKpi:               kpi.IdKpi,
			Kpi:                 kpi.Kpi,
			Rumus:               kpi.Rumus,
			IdPerspektif:        kpi.IdPersfektif,
			Persfektif:          kpi.Perspektif,
			IdKeteranganProject: kpi.IdKeteranganProject,
			KeteranganProject:   kpi.KeteranganProject,
			LinkDokumenSumber:   kpi.LampiranFile,
			TotalSubKpi:         kpi.TotalSubKpi,
			KpiSubDetail:        subList,
		}
	}

	resultList := make([]dto.DataResult, len(dataDB.ResultList))
	for i, r := range dataDB.ResultList {
		resultList[i] = dto.DataResult{
			IdDetailResult:   r.IdDetailResult,
			NamaResult:       r.NamaResult,
			DeskripsiResult:  r.DeskripsiResult,
			RealisasiResult:  r.RealisasiResult,
			LampiranEvidence: r.LampiranEvidence,
		}
	}

	processList := make([]dto.DataProcess, len(dataDB.ProcessList))
	for i, p := range dataDB.ProcessList {
		processList[i] = dto.DataProcess{
			IdDetailProcess:  p.IdDetailMethod,
			NamaProcess:      p.NamaMethod,
			DeskripsiProcess: p.DeskripsiMethod,
			RealisasiProcess: p.RealisasiMethod,
			LampiranEvidence: p.LampiranEvidence,
		}
	}

	contextList := make([]dto.DataContext, len(dataDB.ContextList))
	for i, ctx := range dataDB.ContextList {
		contextList[i] = dto.DataContext{
			IdDetailContext:  ctx.IdDetailChallenge,
			NamaContext:      ctx.NamaChallenge,
			DeskripsiContext: ctx.DeskripsiChallenge,
			RealisasiContext: ctx.RealisasiChallenge,
			LampiranEvidence: ctx.LampiranEvidence,
		}
	}

	data = &dto.GetDetailPencapaianKpiResponse{
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
		ApprovalListPenyusunan:   approvalListPenyusunan,
		ApprovalListRealisasi:    approvalListRealisasi,
		ApprovalListValidasi:     approvalListValidasi,
		Catatan:                  catatanList,
		TotalBobot:               dataDB.TotalBobot,
		TotalPencapaian:          dataDB.TotalPencapaian,
		TotalBobotPengurang:      dataDB.TotalBobotPengurang,
		TotalPencapaianPost:      dataDB.TotalPencapaianPost,
		LampiranValidasi:         lampiranValidasi,
		TotalKpi:                 dataDB.TotalKpi,
		KpiList:                  kpiList,
		TotalResult:              dataDB.TotalResult,
		ResultList:               resultList,
		TotalProcess:             dataDB.TotalProcess,
		ProcessList:              processList,
		TotalContext:             dataDB.TotalContext,
		ContextList:              contextList,
		QualifierOverallValidasi: qualifierOverall,
	}

	return data, nil
}

// =============================================================================
// DOWNLOAD
// =============================================================================

// GetExcelPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-excel.
func (s *pencapaianKpiService) GetExcelPencapaianKpi(
	req *dto.GetExcelPencapaianKpiRequest,
) ([]byte, string, error) {
	indikatorDB, err := s.repo.GetIndikatorPencapaian()
	if err != nil {
		return nil, "", err
	}
	exportData, err := s.buildPencapaianKpiExportData(req.IdPengajuan, indikatorDB)
	if err != nil {
		return nil, "", err
	}
	return file_export.GeneratePencapaianKpiExcel(exportData)
}

// GetPdfPencapaianKpi digunakan oleh endpoint POST /pencapaian-kpi/get-pdf.
func (s *pencapaianKpiService) GetPdfPencapaianKpi(
	req *dto.GetPdfPencapaianKpiRequest,
) ([]byte, string, error) {
	indikatorDB, err := s.repo.GetIndikatorPencapaian()
	if err != nil {
		return nil, "", err
	}
	exportData, err := s.buildPencapaianKpiExportData(req.IdPengajuan, indikatorDB)
	if err != nil {
		return nil, "", err
	}
	return file_export.GeneratePencapaianKpiPDF(exportData)
}

