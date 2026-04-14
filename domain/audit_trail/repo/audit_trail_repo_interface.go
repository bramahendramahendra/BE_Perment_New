package repo

import (
	model "permen_api/domain/audit_trail/model"

	"gorm.io/gorm"
)

type (
	AuditTrailRepoInterface interface {
		InsertAuditTrail(data *model.LogAudit) error
		GetDB() *gorm.DB
	}

	auditTrailRepo struct {
		db *gorm.DB
	}
)

func NewAuditTrailRepo(db *gorm.DB) *auditTrailRepo {
	return &auditTrailRepo{db: db}
}

func (r *auditTrailRepo) GetDB() *gorm.DB {
	return r.db
}
