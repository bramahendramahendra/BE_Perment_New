package file_export

import (
	"bytes"
	"fmt"
	"strconv"

	dto "permen_api/domain/penyusunan_kpi/dto"

	"github.com/jung-kurt/gofpdf"
)

// GenerateKpiPDF membuat file PDF dari KpiExportData dan mengembalikan byte slice
// serta nama file yang direkomendasikan untuk header Content-Disposition.
func GenerateKpiPDF(exportData *dto.KpiExportData) ([]byte, string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Warna palette (sesuai gambar)
	headerBgR, headerBgG, headerBgB := 31, 73, 125   // biru tua  (#1F497D)
	headerFgR, headerFgG, headerFgB := 255, 255, 255 // putih
	rowBlueR, rowBlueG, rowBlueB := 189, 215, 238    // biru muda (#BDD7EE)
	rowPeachR, rowPeachG, rowPeachB := 252, 228, 214 // peach     (#FCE4D6)
	rowGreenR, rowGreenG, rowGreenB := 226, 240, 217 // hijau muda (#E2F0D9)
	textR, textG, textB := 0, 0, 0

	// Judul
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(textR, textG, textB)
	pdf.CellFormat(0, 7, exportData.NamaDivisi, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "Tahun "+exportData.Tahun, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	// Lebar kolom — total ~267mm untuk A4 landscape (297 - 15*2 margin)
	// No | KPI | Bobot (%) | Target Tahunan | Capping
	colWidths := []float64{12, 100, 25, 80, 25}
	headers := []string{"No", "KPI", "Bobot (%)", "Target Tahunan", "Capping"}
	rowHeight := 8.0

	// Header tabel
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(headerBgR, headerBgG, headerBgB)
	pdf.SetTextColor(headerFgR, headerFgG, headerFgB)
	pdf.SetDrawColor(255, 255, 255)
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], rowHeight, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Baris data — alternating per 3 baris
	pdf.SetFont("Arial", "", 9)
	pdf.SetDrawColor(200, 200, 200)
	dataAligns := []string{"C", "L", "C", "L", "C"}

	for _, row := range exportData.Rows {
		group := ((row.No - 1) / 3) % 3
		switch group {
		case 0:
			pdf.SetFillColor(rowBlueR, rowBlueG, rowBlueB)
		case 1:
			pdf.SetFillColor(rowPeachR, rowPeachG, rowPeachB)
		default:
			pdf.SetFillColor(rowGreenR, rowGreenG, rowGreenB)
		}
		pdf.SetTextColor(textR, textG, textB)

		values := []string{
			strconv.Itoa(row.No),
			row.KpiNama,
			row.Bobot,
			row.TargetTahunan,
			row.Capping,
		}
		for i, v := range values {
			pdf.CellFormat(colWidths[i], rowHeight, v, "1", 0, dataAligns[i], true, 0, "")
		}
		pdf.Ln(-1)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("gagal generate PDF: %w", err)
	}

	filename := fmt.Sprintf("KPI_%s_%s_%s.pdf",
		exportData.NamaDivisi, exportData.Tahun, exportData.Triwulan)

	return buf.Bytes(), filename, nil
}
