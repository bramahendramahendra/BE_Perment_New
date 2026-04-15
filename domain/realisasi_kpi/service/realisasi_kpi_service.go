package service

import (
	"fmt"
	"math"
	"strings"

	"mime/multipart"

	dto "permen_api/domain/realisasi_kpi/dto"
	"permen_api/domain/realisasi_kpi/utils"
	customErrors "permen_api/errors"
)

// =============================================================================
// VALIDATE
// =============================================================================

func (s *realisasiKpiService) ValidateRealisasiKpi(
	req *dto.ValidateRealisasiKpiRequest,
	file *multipart.FileHeader,
) (data dto.ValidateRealisasiKpiResponse, err error) {

	if file == nil {
		return data, &customErrors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files'",
		}
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("file '%s' bukan format Excel (.xlsx)", file.Filename),
		}
	}

	// Ambil triwulan dari DB berdasarkan id_pengajuan
	triwulan, err := s.repo.GetTriwulanByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return data, err
	}
	req.Triwulan = triwulan

	// Parse dan validasi file Excel
	kpiRows, kpiSubDetails, err := utils.ParseAndValidateRealisasiExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup DB per baris: isi IdSubDetail, IdDetail, TargetKuantitatifTriwulan, Rumus,
	// lalu hitung Pencapaian dan Skor
	if err := s.enrichRowsFromDB(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Ambil header (tahun, kostl, kostlTx) untuk response
	tahun, _, kostl, kostlTx, err := s.repo.GetKpiHeaderByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return data, err
	}

	// Build resultList, processList, contextList (hanya TW2 dan TW4)
	resultList := []dto.RealisasiResult{}
	processList := []dto.RealisasiProcess{}
	contextList := []dto.RealisasiContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(kpiRows, kpiSubDetails)
	}

	// Simpan ke DB (status 80 = draft realisasi)
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

	// Build response
	totalSubKpi := 0
	for _, kpiRow := range kpiRows {
		totalSubKpi += len(kpiSubDetails[kpiRow.KpiIndex])
	}

	data = dto.ValidateRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Tahun:       tahun,
		Triwulan:    req.Triwulan,
		Divisi: dto.DivisiRealisasiResponse{
			Kostl:   kostl,
			KostlTx: kostlTx,
		},
		Entry: dto.EntryRealisasiResponse{
			EntryUser: req.EntryUser,
			EntryName: req.EntryName,
			EntryTime: req.EntryTime,
		},
		TotalSubKpi: totalSubKpi,
		SubKpiList:  utils.BuildSubKpiDetailList(kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
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

	if file == nil {
		return data, &customErrors.BadRequestError{
			Message: "file Excel tidak ditemukan, pastikan mengirim file via field 'files'",
		}
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("file '%s' bukan format Excel (.xlsx)", file.Filename),
		}
	}

	// Ambil triwulan dari DB berdasarkan id_pengajuan
	triwulan, err := s.repo.GetTriwulanByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return data, err
	}
	req.Triwulan = triwulan

	kpiRows, kpiSubDetails, err := utils.ParseAndValidateRealisasiExcel(file, req.Triwulan)
	if err != nil {
		return data, &customErrors.BadRequestError{
			Message: fmt.Sprintf("validasi file Excel '%s' gagal: %s", file.Filename, err.Error()),
		}
	}

	// Lookup DB per baris: isi IdSubDetail, IdDetail, TargetKuantitatifTriwulan, Rumus,
	// lalu hitung Pencapaian dan Skor
	if err := s.enrichRowsFromDB(req.IdPengajuan, kpiRows, kpiSubDetails); err != nil {
		return data, err
	}

	// Ambil header (tahun) untuk response
	tahun, _, _, _, err := s.repo.GetKpiHeaderByIdPengajuan(req.IdPengajuan)
	if err != nil {
		return data, err
	}

	// Build resultList, processList, contextList (hanya TW2 dan TW4)
	resultList := []dto.RealisasiResult{}
	processList := []dto.RealisasiProcess{}
	contextList := []dto.RealisasiContext{}
	if utils.IsExtendedTriwulan(req.Triwulan) {
		resultList = utils.BuildResultList(kpiRows, kpiSubDetails)
		processList = utils.BuildProcessList(kpiRows, kpiSubDetails)
		contextList = utils.BuildContextList(kpiRows, kpiSubDetails)
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

	// Build response
	totalSubKpi := 0
	for _, kpiRow := range kpiRows {
		totalSubKpi += len(kpiSubDetails[kpiRow.KpiIndex])
	}

	data = dto.RevisionRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Tahun:       tahun,
		Triwulan:    req.Triwulan,
		TotalSubKpi: totalSubKpi,
		SubKpiList:  utils.BuildSubKpiDetailList(kpiRows, kpiSubDetails),
		ResultList:  resultList,
		ProcessList: processList,
		ContextList: contextList,
	}

	return data, nil
}

// =============================================================================
// CREATE — submit ke approval
// =============================================================================

func (s *realisasiKpiService) CreateRealisasiKpi(
	req *dto.CreateRealisasiKpiRequest,
) (data dto.CreateRealisasiKpiResponse, err error) {

	if err := s.repo.CreateRealisasiKpi(req); err != nil {
		return data, err
	}

	data = dto.CreateRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Message:     "Realisasi KPI berhasil disubmit untuk approval",
	}

	return data, nil
}

// =============================================================================
// APPROVAL
// =============================================================================

func (s *realisasiKpiService) ApprovalRealisasiKpi(
	req *dto.ApprovalRealisasiKpiRequest,
) (data dto.ApprovalRealisasiKpiResponse, err error) {

	if err := s.repo.ApprovalRealisasiKpi(req); err != nil {
		return data, err
	}

	message := "Realisasi KPI berhasil diapprove"
	if req.Status == "reject" {
		message = "Realisasi KPI berhasil ditolak"
	}

	data = dto.ApprovalRealisasiKpiResponse{
		IdPengajuan: req.IdPengajuan,
		Message:     message,
	}

	return data, nil
}

// =============================================================================
// GET ALL APPROVAL
// =============================================================================

func (s *realisasiKpiService) GetAllApprovalRealisasiKpi(
	req *dto.GetAllApprovalRealisasiKpiRequest,
) ([]*dto.GetAllApprovalRealisasiKpiResponse, int64, error) {

	records, total, err := s.repo.GetAllApprovalRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.GetAllApprovalRealisasiKpiResponse, 0, len(records))
	for _, r := range records {
		result = append(result, &dto.GetAllApprovalRealisasiKpiResponse{
			IdPengajuan:        r.IdPengajuan,
			Tahun:              r.Tahun,
			Triwulan:           r.Triwulan,
			Kostl:              r.Kostl,
			KostlTx:            r.KostlTx,
			EntryUserRealisasi: r.EntryUserRealisasi,
			EntryNameRealisasi: r.EntryNameRealisasi,
			EntryTimeRealisasi: r.EntryTimeRealisasi,
			Status:             r.Status,
			StatusDesc:         r.StatusDesc,
		})
	}

	return result, total, nil
}

// =============================================================================
// GET ALL TOLAKAN
// =============================================================================

func (s *realisasiKpiService) GetAllTolakanRealisasiKpi(
	req *dto.GetAllTolakanRealisasiKpiRequest,
) ([]*dto.GetAllTolakanRealisasiKpiResponse, int64, error) {

	records, total, err := s.repo.GetAllTolakanRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.GetAllTolakanRealisasiKpiResponse, 0, len(records))
	for _, r := range records {
		result = append(result, &dto.GetAllTolakanRealisasiKpiResponse{
			IdPengajuan:        r.IdPengajuan,
			Tahun:              r.Tahun,
			Triwulan:           r.Triwulan,
			Kostl:              r.Kostl,
			KostlTx:            r.KostlTx,
			EntryUserRealisasi: r.EntryUserRealisasi,
			EntryNameRealisasi: r.EntryNameRealisasi,
			EntryTimeRealisasi: r.EntryTimeRealisasi,
			CatatanTolakan:     r.CatatanTolakan,
			Status:             r.Status,
			StatusDesc:         r.StatusDesc,
		})
	}

	return result, total, nil
}

// =============================================================================
// GET ALL DAFTAR REALISASI
// =============================================================================

func (s *realisasiKpiService) GetAllDaftarRealisasiKpi(
	req *dto.GetAllDaftarRealisasiKpiRequest,
) ([]*dto.GetAllDaftarRealisasiKpiResponse, int64, error) {

	records, total, err := s.repo.GetAllDaftarRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.GetAllDaftarRealisasiKpiResponse, 0, len(records))
	for _, r := range records {
		result = append(result, &dto.GetAllDaftarRealisasiKpiResponse{
			IdPengajuan:        r.IdPengajuan,
			Tahun:              r.Tahun,
			Triwulan:           r.Triwulan,
			Kostl:              r.Kostl,
			KostlTx:            r.KostlTx,
			EntryUserRealisasi: r.EntryUserRealisasi,
			EntryNameRealisasi: r.EntryNameRealisasi,
			EntryTimeRealisasi: r.EntryTimeRealisasi,
			Status:             r.Status,
			StatusDesc:         r.StatusDesc,
			TotalBobot:         r.TotalBobot,
			TotalPencapaian:    r.TotalPencapaian,
		})
	}

	return result, total, nil
}

// =============================================================================
// GET ALL DAFTAR APPROVAL
// =============================================================================

func (s *realisasiKpiService) GetAllDaftarApprovalRealisasiKpi(
	req *dto.GetAllDaftarApprovalRealisasiKpiRequest,
) ([]*dto.GetAllDaftarApprovalRealisasiKpiResponse, int64, error) {

	records, total, err := s.repo.GetAllDaftarApprovalRealisasiKpi(req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.GetAllDaftarApprovalRealisasiKpiResponse, 0, len(records))
	for _, r := range records {
		result = append(result, &dto.GetAllDaftarApprovalRealisasiKpiResponse{
			IdPengajuan:        r.IdPengajuan,
			Tahun:              r.Tahun,
			Triwulan:           r.Triwulan,
			Kostl:              r.Kostl,
			KostlTx:            r.KostlTx,
			EntryUserRealisasi: r.EntryUserRealisasi,
			EntryNameRealisasi: r.EntryNameRealisasi,
			EntryTimeRealisasi: r.EntryTimeRealisasi,
			Status:             r.Status,
			StatusDesc:         r.StatusDesc,
		})
	}

	return result, total, nil
}

// =============================================================================
// GET DETAIL
// =============================================================================

func (s *realisasiKpiService) GetDetailRealisasiKpi(
	req *dto.GetDetailRealisasiKpiRequest,
) (*dto.GetDetailRealisasiKpiResponse, error) {
	return s.repo.GetDetailRealisasiKpi(req)
}

// =============================================================================
// PRIVATE HELPERS
// =============================================================================

// enrichRowsFromDB melakukan lookup ke DB untuk setiap baris sub KPI Excel:
//   - Mencari id_sub_detail, id_detail, target_kuantitatif_triwulan, rumus
//     berdasarkan id_pengajuan + kpi_name + sub_kpi_name
//   - Menghitung Pencapaian dan Skor mengikuti logika bisnis BE_Perment_Old
//
// Logika kalkulasi (dari models/M_realisasi.go):
//
//	Rumus "1" (Maximize): Pencapaian = (RealisasiKuantitatif / TargetKuantitatif) * 100
//	Rumus "0" (Minimize): Pencapaian = (TargetKuantitatif / RealisasiKuantitatif) * 100
//	Capping diterapkan jika Pencapaian > nilai Capping
//	Skor = (Pencapaian * Bobot) / 100
func (s *realisasiKpiService) enrichRowsFromDB(
	idPengajuan string,
	kpiRows []dto.KpiRow,
	kpiSubDetails map[int][]dto.KpiSubDetailRow,
) error {
	for _, kpiRow := range kpiRows {
		subRows := kpiSubDetails[kpiRow.KpiIndex]
		for i := range subRows {
			sub := &kpiSubDetails[kpiRow.KpiIndex][i]

			idSubDetail, idDetail, rumus, targetKuantitatif, err :=
				s.repo.LookupSubDetailByKpiSubKpi(idPengajuan, kpiRow.Kpi, sub.SubKPI)
			if err != nil {
				return err
			}

			sub.IdSubDetail = idSubDetail
			sub.IdDetail = idDetail
			sub.Rumus = rumus
			sub.TargetKuantitatifTriwulan = targetKuantitatif

			// Hitung Pencapaian dan Skor
			pencapaian, skor := calculatePencapaianSkor(
				rumus,
				sub.RealisasiKuantitatif,
				targetKuantitatif,
				sub.Capping,
				sub.Bobot,
			)

			sub.Pencapaian = pencapaian
			sub.Skor = skor
		}
	}

	return nil
}

// calculatePencapaianSkor menghitung Pencapaian (%) dan Skor dari nilai realisasi.
//
// Mengikuti logika bisnis models/M_realisasi.go dari BE_Perment_Old secara persis:
//
//	rumus == "1" → Maximize : Pencapaian = (realisasi / target) * 100
//	rumus == "0" → Minimize : Pencapaian = (target / realisasi) * 100
//	Capping diterapkan jika Pencapaian melebihi batas ("100%" = 100, "110%" = 110)
//	Skor = (Pencapaian * bobot) / 100
//	Jika target = 0 atau rumus tidak dikenal → Pencapaian = 0, Skor = 0
func calculatePencapaianSkor(
	rumus string,
	realisasiKuantitatif float64,
	targetKuantitatif float64,
	cappingStr string,
	bobot float64,
) (pencapaian, skor float64) {

	// Parse nilai capping numerik dari string "100%" atau "110%"
	cappingValue := parseCapping(cappingStr)

	switch rumus {
	case "1": // Maximize
		if targetKuantitatif == 0 {
			return 0, 0
		}
		pencapaian = (realisasiKuantitatif / targetKuantitatif) * 100

	case "0": // Minimize
		if realisasiKuantitatif == 0 {
			return 0, 0
		}
		pencapaian = (targetKuantitatif / realisasiKuantitatif) * 100

	default:
		// Rumus tidak dikenal (misal "0" untuk sub KPI lain/non-standar)
		return 0, 0
	}

	// Terapkan capping
	if cappingValue > 0 && pencapaian > cappingValue {
		pencapaian = cappingValue
	}

	skor = (pencapaian * bobot) / 100

	// Bulatkan 2 desimal agar konsisten
	pencapaian = math.Round(pencapaian*100) / 100
	skor = math.Round(skor*100) / 100

	return pencapaian, skor
}

// parseCapping mengubah string capping ("100%" atau "110%") menjadi nilai float64.
// Mengembalikan 0 jika format tidak dikenal (tidak ada capping yang diterapkan).
func parseCapping(cappingStr string) float64 {
	switch strings.TrimSpace(cappingStr) {
	case "100%":
		return 100.0
	case "110%":
		return 110.0
	default:
		return 0
	}
}

