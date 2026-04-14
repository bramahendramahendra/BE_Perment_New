package repo

import (
	model "permen_api/domain/audit_trail/model"
	"time"
)

func (r *auditTrailRepo) InsertAuditTrail(data *model.LogAudit) error {
	return r.db.Exec(
		"INSERT INTO log_request (ipaddress, userid, `function`, body, response, errordesc, time_request) VALUES (?, ?, ?, ?, ?, ?, ?)",
		data.Ipaddress,
		data.Userid,
		data.Function,
		data.Body,
		data.Response,
		data.Errordesc,
		time.Now(),
	).Error
}
