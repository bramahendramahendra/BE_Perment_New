package file_export

import (
	"bytes"
	"fmt"
	"strings"

	dto "permen_api/domain/pencapaian_kpi/dto"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePencapaianKpiPDF membuat file PDF dari PencapaianKpiExportData.
func GeneratePencapaianKpiPDF(exportData *dto.PencapaianKpiExportData) ([]byte, string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	headerBgR, headerBgG, headerBgB := 31, 73, 125
	headerFgR, headerFgG, headerFgB := 255, 255, 255
	rowGreenR, rowGreenG, rowGreenB := 226, 240, 217
	textR, textG, textB := 0, 0, 0

	title := buildPencapaianPdfTitle(exportData)
	subtitle := fmt.Sprintf("Tahun %s - TW %s", exportData.Tahun, exportData.TriwulanNum)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(textR, textG, textB)
	pdf.CellFormat(0, 6, title, "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, subtitle, "", 1, "L", false, 0, "")
	pdf.Ln(3)

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
	colWidths := []float64{8, 52, 26, 13, 25, 26, 25, 26, 24, 26, 26}
	colAligns := []string{"C", "L", "L", "C", "C", "C", "C", "C", "C", "C", "C"}
	lineHeight := 5.5
	headerHeight := 12.0

	leftMargin, topMargin, _, _ := pdf.GetMargins()
	pageBreakY := 210.0 - topMargin - 22.0

	drawTableHeader := func() {
		pdf.SetFont("Arial", "B", 8)
		pdf.SetFillColor(headerBgR, headerBgG, headerBgB)
		pdf.SetTextColor(headerFgR, headerFgG, headerFgB)
		pdf.SetDrawColor(200, 200, 200)

		hy := pdf.GetY()
		for i, h := range headers {
			x := leftMargin
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}
			pdf.SetXY(x, hy)
			pdf.CellFormat(colWidths[i], headerHeight, "", "1", 0, "C", true, 0, "")
			textH := lineHeight * float64(len(pdf.SplitLines([]byte(h), colWidths[i]-2)))
			if textH < lineHeight {
				textH = lineHeight
			}
			offsetY := (headerHeight - textH) / 2
			if offsetY < 0 {
				offsetY = 0
			}
			pdf.SetXY(x, hy+offsetY)
			pdf.MultiCell(colWidths[i], lineHeight, h, "", "C", false)
		}
		pdf.SetXY(leftMargin, hy+headerHeight)
	}

	drawTableHeader()

	pdf.SetFont("Arial", "", 8)
	pdf.SetDrawColor(200, 200, 200)

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

		maxLines := 1
		for i, v := range values {
			n := len(pdf.SplitLines([]byte(v), colWidths[i]-2))
			if n < 1 {
				n = 1
			}
			if n > maxLines {
				maxLines = n
			}
		}
		rowHeight := lineHeight*float64(maxLines) + 3
		if rowHeight < lineHeight+3 {
			rowHeight = lineHeight + 3
		}

		if pdf.GetY()+rowHeight > pageBreakY {
			pdf.AddPage()
			drawTableHeader()
			pdf.SetFont("Arial", "", 8)
			pdf.SetDrawColor(200, 200, 200)
		}

		pdf.SetFillColor(rowGreenR, rowGreenG, rowGreenB)
		pdf.SetTextColor(textR, textG, textB)

		rx := leftMargin
		ry := pdf.GetY()

		for i, v := range values {
			x := rx
			for j := 0; j < i; j++ {
				x += colWidths[j]
			}
			pdf.SetXY(x, ry)
			pdf.CellFormat(colWidths[i], rowHeight, "", "1", 0, colAligns[i], true, 0, "")

			if i == 8 || i == 9 || i == 10 {
				drawPencapaianIndicatorCell(pdf, x, ry, colWidths[i], rowHeight, lineHeight, v,
					exportData.Indikator,
					textR, textG, textB,
					rowGreenR, rowGreenG, rowGreenB)
			} else {
				nLines := len(pdf.SplitLines([]byte(v), colWidths[i]-2))
				if nLines < 1 {
					nLines = 1
				}
				textH := lineHeight * float64(nLines)
				offsetY := (rowHeight - textH) / 2
				if offsetY < 0 {
					offsetY = 0
				}
				pdf.SetXY(x, ry+offsetY)
				pdf.MultiCell(colWidths[i], lineHeight, v, "", colAligns[i], false)
			}
		}
		pdf.SetXY(rx, ry+rowHeight)
	}

	footerR, footerG, footerB := 31, 73, 125
	pdf.Ln(6)
	pdf.SetTextColor(footerR, footerG, footerB)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "PT Bank Rakyat Indonesia (Persero) Tbk", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 5, "Planning, Budgeting & Performance Management Group", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, "Gedung BRI II. Jalan Jendral Sudirman Kav. 44-46. Jakarta, Indonesia 10210", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "BI", 9)
	pdf.CellFormat(0, 5, "Integrity, Collaborative, Accountability, Growth Mindset, Customer Focus", "", 1, "L", false, 0, "")

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

// drawPencapaianIndicatorCell menggambar lingkaran indikator warna + teks persen di dalam sel.
func drawPencapaianIndicatorCell(pdf *gofpdf.Fpdf, x, y, w, h, lineH float64, val string, indikator []dto.IndikatorPencapaian, tr, tg, tb, restoreR, restoreG, restoreB int) {
	const circleR = 2.0
	const circlePad = 3.5

	trimmed := strings.TrimSpace(val)
	hasCircle := trimmed != "" && trimmed != "-"

	if hasCircle {
		pctStr := strings.TrimSuffix(trimmed, "%")
		var pct float64
		fmt.Sscanf(pctStr, "%f", &pct)

		cr, cg, cb := 195, 0, 2
		for _, item := range indikator {
			if pct <= item.Value {
				switch item.Warna {
				case "hijau":
					cr, cg, cb = 114, 173, 74
				case "kuning":
					cr, cg, cb = 249, 195, 1
				case "merah":
					cr, cg, cb = 195, 0, 2
				}
			}
		}

		cx := x + circlePad
		cy := y + h/2
		pdf.SetFillColor(cr, cg, cb)
		pdf.SetDrawColor(0, 0, 0)
		pdf.Circle(cx, cy, circleR, "F")
	}

	textX := x
	textW := w
	if hasCircle {
		textX = x + circlePad*2 + circleR
		textW = w - (circlePad*2 + circleR)
	}
	nLines := len(pdf.SplitLines([]byte(trimmed), textW-1))
	if nLines < 1 {
		nLines = 1
	}
	offsetY := (h - lineH*float64(nLines)) / 2
	if offsetY < 0 {
		offsetY = 0
	}
	pdf.SetTextColor(tr, tg, tb)
	pdf.SetXY(textX, y+offsetY)
	pdf.MultiCell(textW, lineH, trimmed, "", "R", false)

	pdf.SetFillColor(restoreR, restoreG, restoreB)
	pdf.SetDrawColor(200, 200, 200)
}

func buildPencapaianPdfTitle(exportData *dto.PencapaianKpiExportData) string {
	prefix := ""
	if exportData.IsDraft {
		prefix = "Draft "
	}
	return fmt.Sprintf("%s%s (Pencapaian: %s%%)", prefix, exportData.NamaDivisi, exportData.TotalPencapaian)
}
