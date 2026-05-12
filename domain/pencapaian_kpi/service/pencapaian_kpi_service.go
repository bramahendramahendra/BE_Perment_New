package service

import (
	"encoding/json"
	"fmt"
	"strings"

	dto "permen_api/domain/pencapaian_kpi/dto"
	model "permen_api/domain/pencapaian_kpi/model"
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
				IndikatorPencapaian:              resolveIndikator(sub.Pencapaian, indikatorDB),
				Skor:                             sub.Skor,
				PencapaianQualifierValidated:     sub.PencapaianQualifierValidated,
				IndikatorPencapaianQualifier:     resolveIndikator(sub.PencapaianQualifierValidated, indikatorDB),
				PencapaianPostQualifierValidated: sub.PencapaianPostQualifierValidated,
				IndikatorPencapaianPostQualifier: resolveIndikator(sub.PencapaianPostQualifierValidated, indikatorDB),
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

// buildPencapaianKpiExportData mengambil data dari DB dan mengubahnya ke PencapaianKpiExportData.
func (s *pencapaianKpiService) buildPencapaianKpiExportData(idPengajuan string, indikatorDB []*model.IndikatorPencapaian) (*dto.PencapaianKpiExportData, error) {
	dataDB, err := s.repo.GetDetailPencapaianKpi(&dto.GetDetailPencapaianKpiRequest{IdPengajuan: idPengajuan})
	if err != nil {
		return nil, err
	}
	if dataDB.IdPengajuan == "" {
		return nil, &customErrors.BadRequestError{
			Message: fmt.Sprintf("id_pengajuan '%s' tidak ditemukan", idPengajuan),
		}
	}

	twNum := strings.TrimPrefix(dataDB.Triwulan, "TW")

	// Status 8 = Validasi Final Disetujui (bukan draft)
	isDraft := false

	var rows []dto.PencapaianKpiExportRow
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

			rows = append(rows, dto.PencapaianKpiExportRow{
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
		rows = []dto.PencapaianKpiExportRow{}
	}

	indikator := make([]dto.IndikatorPencapaian, 0, len(indikatorDB))
	for _, item := range indikatorDB {
		indikator = append(indikator, dto.IndikatorPencapaian{
			Warna: item.IndikatorWarna,
			Value: item.IndikatorValue,
		})
	}

	return &dto.PencapaianKpiExportData{
		NamaDivisi:      dataDB.KostlTx,
		Triwulan:        dataDB.Triwulan,
		TriwulanNum:     twNum,
		Tahun:           dataDB.Tahun,
		TotalPencapaian: dataDB.TotalPencapaian,
		IsDraft:         isDraft,
		Rows:            rows,
		Indikator:       indikator,
	}, nil
}

// resolveIndikator menentukan warna indikator berdasarkan nilai pct dan daftar indikator dari DB.
// Logika: iterasi descending by value, last match wins. Default "merah" jika tidak ada match.
func resolveIndikator(pct float64, indikator []*model.IndikatorPencapaian) string {
	warna := "merah"
	for _, item := range indikator {
		if pct <= item.IndikatorValue {
			warna = item.IndikatorWarna
		}
	}
	return warna
}
