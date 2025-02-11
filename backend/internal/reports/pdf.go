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
    largeFont = 12
    largestFont = 20
    fontFamily = "Arial"
) 

var (
    primaryColor = [3]int{255, 255, 255}
    secondaryColor = [3]int{200, 200, 200}
)

type PDFGenerator struct {
    pdf      *fpdf.Fpdf
}

func newPDFGenerator() *PDFGenerator {
    pdf := fpdf.New("L", "mm", "A4", "") // 210 x 297 (mm)
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
		return pkg.Errorf(pkg.INTERNAL_ERROR, "Error getting current working directory")
	}

	imagePath := filepath.Join(currentDir, "afya_credit.png")

    p.pdf.SetXY(-(marginX + 110), marginY+50)
	p.pdf.ImageOptions(imagePath, marginX, marginY-5, 35, 30, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	p.pdf.SetFont(fontFamily, "B", smallFont)
	p.pdf.SetXY(marginX, 15)
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

    p.pdf.SetFont(fontFamily, "", mediumFont)
    w = p.pdf.GetStringWidth(fmt.Sprintf("Generated on: %s", generatedDate))
    p.pdf.SetX(center-w/2)
    p.pdf.Cell(0, lineHt, fmt.Sprintf("Generated on: %s", generatedDate))
    p.pdf.Ln(5)

    p.pdf.SetFont(fontFamily, "", mediumFont)
    w = p.pdf.GetStringWidth(fmt.Sprintf("Date Range: %s - %s", fromDate, toDate))
    p.pdf.SetX(center-w/2)
    p.pdf.Cell(0, lineHt, fmt.Sprintf("Date Range: %s - %s", fromDate, toDate))
    p.pdf.Ln(lineHt)
}

func (p *PDFGenerator) drawBox(x, y, w, h float64, color [3]int) {
    p.pdf.SetFillColor(color[0], color[1], color[2])
    p.pdf.Rect(x, y, w, h, "FD")
    p.pdf.Ln(lineHt*4)
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