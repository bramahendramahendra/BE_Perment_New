package service

import (
	"strings"

	dto "permen_api/domain/realisasi_kpi/dto"
	"permen_api/domain/realisasi_kpi/utils"
	customErrors "permen_api/errors"
)

// resolveRealisasiLookups memvalidasi link dokumen sumber dan melakukan enrich data dari DB.
func (s *realisasiKpiService) resolveRealisasiLookups(
	idPengajuan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) error {
	linkFormats, err := s.repo.GetLinkFormats()
	if err != nil {
		return err
	}
	if err := utils.ValidateLinkDokumenSumber(kpiRows, kpiSubDetails, linkFormats); err != nil {
		return &customErrors.BadRequestError{Message: err.Error()}
	}
	return s.enrichRowsFromDB(idPengajuan, kpiRows, kpiSubDetails)
}

// enrichRowsFromDB melakukan lookup ke DB untuk setiap baris sub KPI Excel dan menghitung
// Pencapaian + Skor mengikuti logika bisnis BE_Perment_Old.
func (s *realisasiKpiService) enrichRowsFromDB(
	idPengajuan string,
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
) error {
	for ki := range kpiRows {
		for i := range kpiSubDetails[kpiRows[ki].KpiIndex] {
			sub := &kpiSubDetails[kpiRows[ki].KpiIndex][i]

			lookup, err := s.repo.LookupSubDetailByKpiSubKpi(idPengajuan, kpiRows[ki].Kpi, sub.SubKPI)
			if err != nil {
				return &customErrors.BadRequestError{Message: err.Error()}
			}

			sub.IdSubDetail = lookup.IdSubDetail
			sub.IdDetail = lookup.IdDetail
			sub.IdSubKpi = lookup.IdSubKpi
			sub.Otomatis = lookup.Otomatis
			sub.TerdapatQualifier = lookup.IdQualifier
			sub.IdPolarisasi = lookup.Rumus
			sub.Glossary = lookup.Glossary
			sub.TargetKuantitatifTriwulan = lookup.TargetKuantitatifTriwulan
			sub.TargetTahunan = lookup.TargetTahunan
			sub.TargetKuantitatifTahunan = lookup.TargetKuantitatifTahunan
			sub.DeskripsiQualifier = lookup.DeskripsiQualifier

			if kpiRows[ki].IdDetail == "" {
				kpiRows[ki].IdDetail = lookup.IdDetail
				kpiRows[ki].IdKpi = lookup.IdKpi
				kpiRows[ki].Rumus = lookup.DetailRumus
			}

			// Kolom L dan M hanya disimpan jika id_qualifier = "ya"
			if strings.ToLower(strings.TrimSpace(lookup.IdQualifier)) != "ya" {
				sub.RealisasiQualifier = ""
				sub.RealisasiKuantitatifQualifier = ""
			}

			sub.Pencapaian, sub.Skor = utils.CalculatePencapaianSkor(
				lookup.Rumus,
				sub.RealisasiKuantitatif,
				lookup.TargetKuantitatifTriwulan,
				sub.Capping,
				sub.Bobot,
			)
		}

		for _, sub := range kpiSubDetails[kpiRows[ki].KpiIndex] {
			if sub.LinkDokumenSumber != nil && *sub.LinkDokumenSumber != "" {
				kpiRows[ki].LinkDokumenSumber = *sub.LinkDokumenSumber
				break
			}
		}
	}

	return nil
}
