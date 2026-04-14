package file_export

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

const mimeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

// ExcelFile mewakili file Excel yang siap ditulis ke buffer atau dikirim sebagai response.
type ExcelFile struct {
	f        *excelize.File
	filename string
}

// NewExcelFile membuat instance ExcelFile baru dengan satu sheet awal bernama sheetName.
func NewExcelFile(sheetName string) (*ExcelFile, error) {
	f := excelize.NewFile()
	defaultSheet := f.GetSheetName(0)
	if err := f.SetSheetName(defaultSheet, sheetName); err != nil {
		return nil, fmt.Errorf("gagal set nama sheet: %w", err)
	}
	return &ExcelFile{f: f, filename: "output.xlsx"}, nil
}

// SetFilename menetapkan nama file untuk header Content-Disposition.
func (e *ExcelFile) SetFilename(filename string) {
	e.filename = filename
}

// File mengembalikan pointer ke excelize.File untuk operasi langsung.
func (e *ExcelFile) File() *excelize.File {
	return e.f
}

// Close melepas resource yang dipegang oleh excelize.File.
func (e *ExcelFile) Close() {
	e.f.Close()
}

// ToBytes menulis file Excel ke byte slice.
func (e *ExcelFile) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := e.f.Write(&buf); err != nil {
		return nil, fmt.Errorf("gagal menulis Excel ke buffer: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteResponse menulis file Excel langsung ke gin.Context sebagai HTTP download response.
func (e *ExcelFile) WriteResponse(c *gin.Context) error {
	fileBytes, err := e.ToBytes()
	if err != nil {
		return err
	}
	SetExcelDownloadHeaders(c, e.filename)
	c.Data(http.StatusOK, mimeExcel, fileBytes)
	return nil
}

// SetExcelDownloadHeaders menyetel header HTTP untuk response file Excel download.
func SetExcelDownloadHeaders(c *gin.Context, filename string) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", mimeExcel)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
}

// SendExcel adalah shorthand untuk mengirim []byte Excel langsung ke response
// tanpa perlu membuat ExcelFile terlebih dahulu.
func SendExcel(c *gin.Context, fileBytes []byte, filename string) {
	SetExcelDownloadHeaders(c, filename)
	c.Data(http.StatusOK, mimeExcel, fileBytes)
}
