package file_export

import (
	"bytes"
	"fmt"
	"strings"

	dto "permen_api/domain/validasi_kpi/dto"

	"github.com/jung-kurt/gofpdf"
)

// GenerateValidasiKpiPDF membuat file PDF dari ValidasiKpiExportData.
func GenerateValidasiKpiPDF(exportData *dto.ValidasiKpiExportData) ([]byte, string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Warna palette
	headerBgR, headerBgG, headerBgB := 31, 73, 125   // biru tua  (#1F497D)
	headerFgR, headerFgG, headerFgB := 255, 255, 255 // putih
	rowGreenR, rowGreenG, rowGreenB := 226, 240, 217  // hijau muda (#E2F0D9)
	textR, textG, textB := 0, 0, 0

	// Judul
	title := buildPdfTitle(exportData)
	subtitle := fmt.Sprintf("Tahun %s - TW %s", exportData.Tahun, exportData.TriwulanNum)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(textR, textG, textB)
	pdf.CellFormat(0, 6, title, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, subtitle, "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// Lebar kolom — total ~277mm (A4 landscape 297 - 10*2 margin)
	twNum := exportData.TriwulanNum
	headers := []string{
		"No",
		"KPI",
		"KPI Qualifier",
		"Bobot (%)",
		"Target Triwulan " + twNum,
		"Target Qualifier Triwulan",
		"Realisasi Triwulan " + twNum,
		"Realisasi Qualifier",
		"Pencapaian",
		"Pencapaian Qualifier KPI",
		"Pencapaian KPI Post Qualifier",
	}
	colWidths := []float64{8, 46, 24, 13, 22, 24, 22, 24, 20, 27, 27}
	lineHeight := 6.0

	// Header tabel
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(headerBgR, headerBgG, headerBgB)
	pdf.SetTextColor(headerFgR, headerFgG, headerFgB)
	pdf.SetDrawColor(255, 255, 255)

	startX := pdf.GetX()
	startY := pdf.GetY()

	// Header dengan MultiCell agar teks panjang bisa wrap
	headerHeight := 10.0
	for i, h := range headers {
		x := startX
		for j := 0; j < i; j++ {
			x += colWidths[j]
		}
		pdf.SetXY(x, startY)
		pdf.CellFormat(colWidths[i], headerHeight, "", "1", 0, "C", true, 0, "")
		pdf.SetXY(x, startY)
		pdf.MultiCell(colWidths[i], lineHeight, h, "", "C", false)
	}
	pdf.SetXY(startX, startY+headerHeight)

	// Baris data
	pdf.SetFont("Arial", "", 8)
	pdf.SetDrawColor(180, 180, 180)

	for _, row := range exportData.Rows {
		values := []string{
			fmt.Sprintf("%d", row.No),
			row.Kpi,
			row.ItemQualifier,
			fmt.Sprintf("%.0f", row.Bobot),
			row.TargetTriwulan,
			row.TargetQualifier,
			row.RealisasiValidated,
			row.RealisasiQualifier,
			row.Pencapaian,
			row.PencapaianQualifier,
			row.PencapaianPostQualifier,
		}

		// Hitung tinggi baris
		maxLines := 1
		for i, v := range values {
			n := len(pdf.SplitLines([]byte(v), colWidths[i]-2))
			if n == 0 {
				n = 1
			}
			if n > maxLines {
				maxLines = n
			}
		}
		rowHeight := lineHeight * float64(maxLines)
		if rowHeight < lineHeight {
			rowHeight = lineHeight
		}

		// Cek page break
		if pdf.GetY()+rowHeight > 200 {
			pdf.AddPage()
		}

		pdf.SetFillColor(rowGreenR, rowGreenG, rowGreenB)
		pdf.SetTextColor(textR, textG, textB)

		rx := pdf.GetX()
		ry := pdf.GetY()
		colAligns := []string{"C", "L", "L", "C", "C", "C", "C", "C", "C", "C", "C"}

		for i, v := range values {
			x := rx
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}
			pdf.SetXY(x, ry)
			pdf.CellFormat(colWidths[i], rowHeight, "", "1", 0, colAligns[i], true, 0, "")
			pdf.SetXY(x, ry)
			pdf.MultiCell(colWidths[i], lineHeight, v, "", colAligns[i], false)
		}
		pdf.SetXY(rx, ry+rowHeight)
	}

	// Footer
	pdf.SetY(-20)
	pdf.SetFont("Arial", "B", 8)
	pdf.SetTextColor(textR, textG, textB)
	pdf.CellFormat(0, 5, "PT Bank Rakyat Indonesia (Persero) Tbk", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 7)
	pdf.CellFormat(0, 4, "Planning, Budgeting & Performance Management Group", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Gedung BRI Jl. Jalan Jenderal Sudirman Kav. 44-46. Jakarta, Indonesia 10210", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "I", 7)
	pdf.CellFormat(0, 4, "Integrity, Collaborative, Accountability, Growth Mindset, Customer Focus", "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", fmt.Errorf("gagal generate PDF: %w", err)
	}

	filename := fmt.Sprintf("Pencapaian_KPI_%s_%s_TW%s.pdf",
		strings.ReplaceAll(exportData.NamaDivisi, " ", "_"),
		exportData.Tahun,
		exportData.TriwulanNum,
	)
	return buf.Bytes(), filename, nil
}

func buildPdfTitle(exportData *dto.ValidasiKpiExportData) string {
	prefix := ""
	if exportData.IsDraft {
		prefix = "Draft "
	}
	return fmt.Sprintf("%s%s (Pencapaian: %s%%)", prefix, exportData.NamaDivisi, exportData.TotalPencapaian)
}
