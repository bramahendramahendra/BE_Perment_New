package service

import (
	"fmt"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"
	model "permen_api/domain/validasi_kpi/model"
	customErrors "permen_api/errors"
)

// buildValidasiKpiExportData mengambil data dari DB dan mengubahnya ke ValidasiKpiExportData.
func (s *validasiKpiService) buildValidasiKpiExportData(idPengajuan string, indikatorDB []*model.IndikatorPencapaian) (*dto.ValidasiKpiExportData, error) {
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

	indikator := make([]dto.IndikatorPencapaian, 0, len(indikatorDB))
	for _, item := range indikatorDB {
		indikator = append(indikator, dto.IndikatorPencapaian{
			Warna: item.IndikatorWarna,
			Value: item.IndikatorValue,
		})
	}

	return &dto.ValidasiKpiExportData{
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
