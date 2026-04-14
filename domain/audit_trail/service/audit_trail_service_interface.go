package service

import (
	dto "permen_api/domain/audit_trail/dto"
	repo "permen_api/domain/audit_trail/repo"
)

type (
	AuditTrailServiceInterface interface {
		SaveAuditTrail(req *dto.AuditTrailRequest)
	}

	auditTrailService struct {
		repo repo.AuditTrailRepoInterface
	}
)

func NewAuditTrailService(repo repo.AuditTrailRepoInterface) *auditTrailService {
	return &auditTrailService{repo: repo}
}
