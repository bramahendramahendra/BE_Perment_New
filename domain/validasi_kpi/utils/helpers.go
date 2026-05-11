package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto "permen_api/domain/validasi_kpi/dto"
)

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
