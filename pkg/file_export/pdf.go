package file_export

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const mimePDF = "application/pdf"

// SetPDFDownloadHeaders menyetel header HTTP untuk response file PDF download.
func SetPDFDownloadHeaders(c *gin.Context, filename string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", mimePDF)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
}

// SendPDF mengirim []byte PDF langsung ke response sebagai HTTP download.
func SendPDF(c *gin.Context, fileBytes []byte, filename string) {
	SetPDFDownloadHeaders(c, filename)
	c.Data(http.StatusOK, mimePDF, fileBytes)
}
