package dto

type AuditTrailRequest struct {
	Ip       string
	Userid   string
	Function string
	Body     string
	Response string
	ErrSis   string
}
