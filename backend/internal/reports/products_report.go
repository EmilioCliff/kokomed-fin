package reports

import (
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type productReport struct {
	*excelGenerator
	*PDFGenerator
	data []services.ProductReportData
	filters services.ReportFilters
}

func newProductReport(data []services.ProductReportData, format string, filters services.ReportFilters) *productReport {
	reportGenerator := &productReport{
		data:           data,
		filters: filters,
	}

	switch format {
	case "excel":
		reportGenerator.excelGenerator = newExcelGenerator()
	case "pdf":
		reportGenerator.PDFGenerator = newPDFGenerator()
	}

	return reportGenerator
}

func (pr *productReport) generateExcel(sheetName string) (error) {
	pr.createSheet(sheetName)

	columns := []string{"Product Name", "Loans Issued", "Active Loans", "Completed Loans", "Defaulted Loans", "Amount Disbursed", "Amount Repaid", "Outstanding Amount", "Default Rate(%)"}

	pr.file.SetColWidth(pr.currentSheet, "A", "I", 20)
	pr.file.SetColStyle(pr.currentSheet, "B", pr.createQuantityStyle())
	pr.file.SetColStyle(pr.currentSheet, "C", pr.createQuantityStyle())
	pr.file.SetColStyle(pr.currentSheet, "D", pr.createQuantityStyle())
	pr.file.SetColStyle(pr.currentSheet, "E", pr.createQuantityStyle())
	pr.file.SetColStyle(pr.currentSheet, "F", pr.createMoneyStyle())
	pr.file.SetColStyle(pr.currentSheet, "G", pr.createMoneyStyle())
	pr.file.SetColStyle(pr.currentSheet, "H", pr.createMoneyStyle())

	pr.writeHeader(columns, pr.createHeaderStyle())

	for rowIdx, product := range pr.data {
		row := []interface{}{
			product.ProductName,
			product.LoansIssued,
			product.ActiveLoans,
			product.CompletedLoans,
			product.DefaultedLoans,
			product.AmountDisbursed,
			product.AmountRepaid,
			product.OutstandingAmount,
			product.DefaultRate,
		}
		pr.writeRow(rowIdx+2, row)
	}

	// buffer, err := br.file.WriteToBuffer()
	// buffer.Bytes(),

	if err := pr.file.SaveAs("products_report.xlsx"); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	return pr.closeExcel()
}

func (pr *productReport) generatePDF() (error) {
	if err := pr.addLogo(); err != nil {
		return err
	}

	pr.writeReportMetadata("Branches Report", time.Now().Format("2006-01-02"), pr.filters.StartDate.Format("2006-01-02"), pr.filters.EndDate.Format("2006-01-02"))

	pr.pdf.SetFont(fontFamily, "B", largeFont)
	pr.pdf.Cell(0, lineHt, "Summary:")
	pr.pdf.Ln(lineHt)

	center := pr.getCenterX()
	pr.drawBox(marginX, pr.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Product Name", "Loans Issued", "Active Loans", "Completed Loans", "Defaulted Loans", "Amount Disbursed", "Amount Repaid", "Outstanding Amount", "Default Rate(%)"}
	colWidths := []float64{35, 20, 20, 20, 25, 25, 25, 20, 20}

	pr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    pr.pdf.SetFont("Arial", "B", mediumFont)
    pr.pdf.SetX(marginX)
	pr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "CM", "CM", "CM", "CM", "R", "R", "R", "R"}

	pr.pdf.SetFontStyle("")
    pr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, product := range pr.data {
		row := []interface{}{
			product.ProductName,
			product.LoansIssued,
			product.ActiveLoans,
			product.CompletedLoans,
			product.DefaultedLoans,
			formatMoney(product.AmountDisbursed),
			formatMoney(product.AmountRepaid),
			formatMoney(product.OutstandingAmount),
			formatMoney(product.DefaultRate),
		}
		pr.writeTableRow(row, colWidths, colAlignment)
	}

	// var buffer bytes.Buffer
	// if err := br.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return pr.pdf.OutputFileAndClose("products_report.pdf")
}