package middleware

import (
	"bytes"
	"fmt"
	"strings"

	dto "permen_api/domain/audit_trail/dto"
	service "permen_api/domain/audit_trail/service"
	request_helper "permen_api/helper/request"

	"github.com/gin-gonic/gin"
)

// isFileContentType mengembalikan true jika Content-Type adalah tipe binary/file.
func isFileContentType(ct string) bool {
	fileTypes := []string{
		"application/vnd.openxmlformats-officedocument",
		"application/vnd.ms-excel",
		"application/pdf",
		"application/zip",
		"application/octet-stream",
	}
	for _, ft := range fileTypes {
		if strings.HasPrefix(ct, ft) {
			return true
		}
	}
	return false
}

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

		// Jika response berupa file binary, simpan keterangan saja (bukan isi binary)
		responseStr := blw.body.String()
		contentType := c.Writer.Header().Get("Content-Type")
		if isFileContentType(contentType) {
			filename := ""
			disposition := c.Writer.Header().Get("Content-Disposition")
			if disposition != "" {
				// Ambil nama file dari: attachment; filename="namafile.xlsx"
				for part := range strings.SplitSeq(disposition, ";") {
					part = strings.TrimSpace(part)
					if val, ok := strings.CutPrefix(part, "filename="); ok {
						filename = strings.Trim(val, "\"")
						break
					}
				}
			}
			responseStr = fmt.Sprintf("[FILE RESPONSE] filename: %s, Content-Type: %s, Size: %d bytes", filename, contentType, blw.body.Len())
		}

		svc.SaveAuditTrail(&dto.AuditTrailRequest{
			Ip:       c.ClientIP(),
			Userid:   userid,
			Function: function,
			Body:     reqBodyStr,
			Response: responseStr,
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
