package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	dto "permen_api/domain/realisasi_kpi/dto"
)

// ValidateExcelFile memeriksa file tidak nil dan berekstensi .xlsx.
func ValidateExcelFile(file *multipart.FileHeader) error {
	if file == nil {
		return fmt.Errorf("file Excel tidak ditemukan, pastikan mengirim file via field 'files'")
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
		return fmt.Errorf("file '%s' bukan format Excel (.xlsx)", file.Filename)
	}
	return nil
}

// ProcessApproveApprovalList menandai user sebagai approve dan mencari approver berikutnya.
// Mengembalikan error jika user tidak ditemukan dalam approval list.
func ProcessApproveApprovalList(
	approvalList []dto.ApprovalUserDetail,
	userid string,
	catatan string,
) (updatedList []dto.ApprovalUserDetail, nextApprover string, err error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	currentIdx := -1
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, userid) && approvalList[i].Status == "" {
			approvalList[i].Status = "approve"
			approvalList[i].Keterangan = catatan
			approvalList[i].Waktu = now
			currentIdx = i
			break
		}
	}
	if currentIdx == -1 {
		return nil, "", fmt.Errorf("Data tidak ditemukan.[User Approval Kosong.]")
	}

	for i := currentIdx + 1; i < len(approvalList); i++ {
		if approvalList[i].Status == "" {
			nextApprover = approvalList[i].Userid
			break
		}
	}

	return approvalList, nextApprover, nil
}

// ProcessRejectApprovalList menandai user sebagai reject dalam approval list.
// Mengembalikan error jika user tidak ditemukan.
func ProcessRejectApprovalList(
	approvalList []dto.ApprovalUserDetail,
	userid string,
	catatan string,
) ([]dto.ApprovalUserDetail, error) {
	now := time.Now().Format("2006-01-02 15:04:05")
	found := false
	for i := range approvalList {
		if strings.EqualFold(approvalList[i].Userid, userid) && approvalList[i].Status == "" {
			approvalList[i].Status = "reject"
			approvalList[i].Keterangan = catatan
			approvalList[i].Waktu = now
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("Data tidak ditemukan.[Pastikan User Approval sesuai.]")
	}
	return approvalList, nil
}

// AppendCatatanTolakan menambahkan entry baru ke dalam JSON array catatan_tolakan.
func AppendCatatanTolakan(existingJSON string, entry dto.CatatanDetail) (string, error) {
	var entries []dto.CatatanDetail
	if existingJSON != "" && existingJSON != "null" {
		if err := json.Unmarshal([]byte(existingJSON), &entries); err != nil {
			return "", fmt.Errorf("gagal parse catatan_tolakan: %w", err)
		}
	}
	entries = append(entries, entry)
	result, err := json.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("gagal serialize catatan_tolakan: %w", err)
	}
	return string(result), nil
}

// CalculatePencapaianSkor menghitung Pencapaian (%) dan Skor dari nilai realisasi.
//
// rumus == "1" → Maximize : Pencapaian = (realisasi / target) * 100
// rumus == "0" → Minimize : Pencapaian = (target / realisasi) * 100
// Capping diterapkan jika Pencapaian melebihi batas ("100%" = 100, "110%" = 110)
// Skor = (Pencapaian * bobot) / 100
func CalculatePencapaianSkor(
	rumus string,
	realisasiKuantitatif float64,
	targetKuantitatif float64,
	cappingStr string,
	bobot float64,
) (pencapaian, skor float64) {
	cappingValue := ParseCapping(cappingStr)

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
		return 0, 0
	}

	if cappingValue > 0 && pencapaian > cappingValue {
		pencapaian = cappingValue
	}

	skor = (pencapaian * bobot) / 100
	pencapaian = math.Round(pencapaian*100) / 100
	skor = math.Round(skor*100) / 100

	return pencapaian, skor
}

// ParseCapping mengubah string capping ("100%" atau "110%") menjadi nilai float64.
// Mengembalikan 0 jika format tidak dikenal.
func ParseCapping(cappingStr string) float64 {
	switch strings.TrimSpace(cappingStr) {
	case "100%":
		return 100.0
	case "110%":
		return 110.0
	default:
		return 0
	}
}

// ValidateLampiranEvidence memvalidasi kolom R/V/Z (Lampiran Evidence) untuk TW2/TW4
// terhadap daftar prefix URL yang diizinkan dari DB (mst_link_format).
func ValidateLampiranEvidence(
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
	allowedPrefixes []string,
) error {
	type evidenceField struct {
		kolom string
		value *string
	}

	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if !row.IsTW24 {
				continue
			}

			fields := []evidenceField{
				{"R (Lampiran Evidence Result)", row.LampiranEvidenceResult},
				{"V (Lampiran Evidence Process)", row.LampiranEvidenceProcess},
				{"Z (Lampiran Evidence Context)", row.LampiranEvidenceContext},
			}

			for _, f := range fields {
				if f.value == nil || *f.value == "" {
					return fmt.Errorf(
						"baris No %d, Sub KPI '%s': Kolom %s tidak boleh kosong",
						row.No, row.SubKPI, f.kolom,
					)
				}
				link := strings.TrimSpace(*f.value)
				valid := false
				for _, prefix := range allowedPrefixes {
					if strings.HasPrefix(link, prefix) {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf(
						"baris No %d, Sub KPI '%s': Kolom %s '%s' tidak sesuai format yang diizinkan",
						row.No, row.SubKPI, f.kolom, link,
					)
				}
			}
		}
	}
	return nil
}

// ValidateLinkDokumenSumber memvalidasi kolom N (Link Dokumen Sumber) setiap baris Excel
// terhadap daftar prefix URL yang diizinkan dari DB (mst_link_format).
func ValidateLinkDokumenSumber(
	kpiRows []dto.RealisasiKpiRow,
	kpiSubDetails map[int][]dto.RealisasiKpiSubDetailRow,
	allowedPrefixes []string,
) error {
	for _, kpiRow := range kpiRows {
		for _, row := range kpiSubDetails[kpiRow.KpiIndex] {
			if row.LinkDokumenSumber == nil || *row.LinkDokumenSumber == "" {
				return fmt.Errorf(
					"baris No %d, Sub KPI '%s': Kolom N (Link Dokumen Sumber) tidak boleh kosong",
					row.No, row.SubKPI,
				)
			}
			link := strings.TrimSpace(*row.LinkDokumenSumber)
			valid := false
			for _, prefix := range allowedPrefixes {
				if strings.HasPrefix(link, prefix) {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf(
					"baris No %d, Sub KPI '%s': Link Dokumen Sumber '%s' tidak sesuai format yang diizinkan",
					row.No, row.SubKPI, link,
				)
			}
		}
	}
	return nil
}
