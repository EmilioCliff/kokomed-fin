package reports

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/go-pdf/fpdf"
)

const (
    companyName = "Afya Credit"
    marginX = 5.0
    marginY = 10.0
    lineHt = 10.0
    smallFont = 8
    mediumFont = 10
    subtitleFont = 14
    largeFont = 16
    largestFont = 30
    fontFamily = "Arial"
) 

var (
    primaryColor = [3]int{255, 255, 255}
    secondaryColor = [3]int{200, 200, 200}
    secondaryAltColor = [3]int{230, 230, 230}
)

type PDFGenerator struct {
    pdf      *fpdf.Fpdf
}

func newPDFGenerator(orientation, size string) *PDFGenerator {
    pdf := fpdf.New(orientation, "mm", size, "") // 210 x 297 (mm)
    pdf.SetMargins(marginX, marginY, marginX)
    pdf.SetAutoPageBreak(true, 15)
    pdf.AddPage()
    pdf.SetFooterFunc(func () {
        pdf.SetXY(marginX, -15)
        pdf.SetFont(fontFamily, "I", smallFont)
        pdf.Cell(marginX+10, lineHt, "This report is system-generated and does not require a signature.")

        pdf.SetXY(-(marginX+7), -15)
        pdf.SetFont(fontFamily, "B", smallFont)
        pdf.CellFormat(0, lineHt, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
    })
    return &PDFGenerator{
        pdf: pdf,
    }
}

func (p *PDFGenerator) addLogo() error {
    currentDir, err := os.Getwd()
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Error getting current working directory: %v", err)
	}

	imagePath := filepath.Join(currentDir, "afya_credit.png")

    p.pdf.SetXY(-(marginX + 110), marginY+50)
	p.pdf.ImageOptions(imagePath, marginX, marginY-5, 55, 45, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	p.pdf.SetFont(fontFamily, "B", smallFont)
	p.pdf.SetXY(marginX+8, 27)
	p.pdf.Cell(70, 35, "Bridging Your Finance Needs")

    return nil
}

func (p *PDFGenerator) writeReportMetadata(title, generatedDate, fromDate, toDate string) {
    center := p.getCenterX()

    p.pdf.SetFont(fontFamily, "B", largestFont)
    w := p.pdf.GetStringWidth(companyName)
    p.pdf.SetXY(center-w/2, marginY)
    p.pdf.Cell(0, lineHt, companyName)
    p.pdf.Ln(-1)

    p.pdf.SetFont(fontFamily, "B", largeFont)
    w = p.pdf.GetStringWidth(title)
    p.pdf.SetX(center-w/2)
    p.pdf.Cell(0, lineHt, title)
    p.pdf.Ln(-1)

    p.pdf.SetFont(fontFamily, "", subtitleFont)
    w = p.pdf.GetStringWidth(fmt.Sprintf("Generated on: %s", generatedDate))
    p.pdf.SetX(center-w/2)
    p.pdf.Cell(0, lineHt, fmt.Sprintf("Generated on: %s", generatedDate))
    p.pdf.Ln(8)

    w = p.pdf.GetStringWidth(fmt.Sprintf("Date Range: %s - %s", fromDate, toDate))
    p.pdf.SetX(center-w/2)
    p.pdf.Cell(0, lineHt, fmt.Sprintf("Date Range: %s    to    %s", fromDate, toDate))
    p.pdf.Ln(lineHt*2)
}

func (p *PDFGenerator) drawBox(x, y, w, h float64, color [3]int) {
    p.pdf.SetFillColor(color[0], color[1], color[2])
    p.pdf.Rect(x, y, w, h, "FD")
    // p.pdf.Ln(lineHt*4)
}

func (p *PDFGenerator) writeTableHeaders(headers []string, colWidths []float64) {
    if len(headers) != len(colWidths) {
        return
    }

    for i, header := range headers {
        p.pdf.CellFormat(colWidths[i], lineHt, header, "1", 0, "CM", true, 0, "")
    }
    p.pdf.Ln(-1)
}

func (p *PDFGenerator) writeTableRow(data []interface{}, colWidths []float64, colAlignment []string) {
    if len(data) != len(colWidths) || len(colAlignment) != len(data) {
        return
    }

    for i, val := range data {
        p.pdf.CellFormat(colWidths[i], lineHt, fmt.Sprintf("%v", val), "1", 0, colAlignment[i], true, 0, "")
    }
    p.pdf.Ln(-1)
}

func (p *PDFGenerator) writeSummary(data map[string]string) {
    var maxWidth float64
    lineSpacing := 6.0 
    boxPadding := 3.0   
    startY := p.pdf.GetY()

    for k, v := range data {
        kWidth := p.pdf.GetStringWidth(k)
        vWidth := p.pdf.GetStringWidth(v)

        totalWidth := kWidth + vWidth + 10
        if totalWidth > maxWidth {
            maxWidth = totalWidth
        }
    }

    boxWidth := maxWidth + (2 * boxPadding)
    boxHeight := float64(len(data))*float64(lineSpacing) + (2 * boxPadding)

    p.drawBox(marginX, startY, boxWidth, boxHeight, secondaryAltColor)

    p.pdf.Ln(boxPadding)
    for k, v := range data {
        p.pdf.SetFontStyle("B")
        kWidth := p.pdf.GetStringWidth(k)
        p.pdf.SetX(marginX + boxPadding + 3)
        p.pdf.Cell(kWidth, lineSpacing, fmt.Sprintf("%s:", k))

        p.pdf.SetFontStyle("")
        p.pdf.SetX(marginX + boxPadding + kWidth + 10)
        p.pdf.Cell(p.pdf.GetStringWidth(v), lineSpacing, v)

        p.pdf.Ln(lineSpacing)
    }
}

func (p *PDFGenerator) writeFooter() {
    p.pdf.SetXY(marginX, -15)
    p.pdf.SetFont(fontFamily, "I", smallFont)
    p.pdf.Cell(0, lineHt, "This report is system-generated and does not require a signature.")
}

func (p *PDFGenerator) getCenterX() float64 {
    pageWidth, _ := p.pdf.GetPageSize()
    return (pageWidth - marginX*2) / 2
}

func (p *PDFGenerator) closePDF() {
    p.pdf.Close()
}