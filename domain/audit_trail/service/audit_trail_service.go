package service

import (
	dto "permen_api/domain/audit_trail/dto"
	model "permen_api/domain/audit_trail/model"
)

// SaveAuditTrail menyimpan log audit trail ke database.
// Hanya dieksekusi jika Body tidak kosong (mengikuti logika bisnis BE_Perment_Old).
func (s *auditTrailService) SaveAuditTrail(req *dto.AuditTrailRequest) {
	if req.Body == "" {
		return
	}

	data := &model.LogAudit{
		Ipaddress: req.Ip,
		Userid:    req.Userid,
		Function:  req.Function,
		Body:      req.Body,
		Response:  req.Response,
		Errordesc: req.ErrSis,
	}

	// Fire-and-forget: error tidak di-propagate (mengikuti logika bisnis BE_Perment_Old)
	_ = s.repo.InsertAuditTrail(data)
}
