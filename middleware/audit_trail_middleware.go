package middleware

import (
	"bytes"
	"strings"

	dto "permen_api/domain/audit_trail/dto"
	service "permen_api/domain/audit_trail/service"
	request_helper "permen_api/helper/request"

	"github.com/gin-gonic/gin"
)

// AuditTrailMiddleware mencatat setiap request ke log_request setelah response dikirim.
// userid diambil dari header "userq" (format: "pernr | nama"), diambil bagian pernr-nya.
// function diambil dari full path route dengan prefix "/api" dihilangkan.
func AuditTrailMiddleware(svc service.AuditTrailServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBodyStr, _ := request_helper.ReadRequestBody(c)

		blw := &auditBodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		// Ambil userid dari header "userq" format: "pernr | nama"
		userid := ""
		userq := c.GetHeader("userq")
		if userq != "" {
			parts := strings.SplitN(userq, " | ", 2)
			userid = strings.TrimSpace(parts[0])
		}

		// Hilangkan prefix "/api" dari full path
		fullPath := c.FullPath()
		function := strings.TrimPrefix(fullPath, "/api")

		svc.SaveAuditTrail(&dto.AuditTrailRequest{
			Ip:       c.ClientIP(),
			Userid:   userid,
			Function: function,
			Body:     reqBodyStr,
			Response: blw.body.String(),
			ErrSis:   "",
		})
	}
}

type auditBodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *auditBodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *auditBodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
