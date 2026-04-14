package middleware

import (
	"bytes"

	dto "permen_api/domain/audit_trail/dto"
	service "permen_api/domain/audit_trail/service"
	request_helper "permen_api/helper/request"

	"github.com/gin-gonic/gin"
)

// AuditTrailMiddleware mencatat setiap request ke log_request setelah response dikirim.
// Logika mengikuti BE_Perment_Old: hanya log jika request body tidak kosong,
// userid diambil dari header "User", IP dari client, function dari full path route.
func AuditTrailMiddleware(svc service.AuditTrailServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Baca request body sebelum c.Next() agar tidak hilang
		reqBodyStr, _ := request_helper.ReadRequestBody(c)

		// Wrap response writer untuk menangkap response body
		blw := &auditBodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		c.Next()

		svc.SaveAuditTrail(&dto.AuditTrailRequest{
			Ip:       c.ClientIP(),
			Userid:   c.GetHeader("User"),
			Function: c.FullPath(),
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
