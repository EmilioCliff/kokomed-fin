package reports

import (
	"bytes"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type productReport struct {
	*excelGenerator
	*PDFGenerator
	data []services.ProductReportData
	summary services.ProductSummary
	filters services.ReportFilters
}

func newProductReport(data []services.ProductReportData, summary services.ProductSummary, format string, filters services.ReportFilters) *productReport {
	reportGenerator := &productReport{
		data:           data,
		summary: summary,
		filters: filters,
	}

	switch format {
	case "excel":
		reportGenerator.excelGenerator = newExcelGenerator()
	case "pdf":
		reportGenerator.PDFGenerator = newPDFGenerator("L", "A4")
	}

	return reportGenerator
}

func (pr *productReport) generateExcel(sheetName string) ([]byte, error) {
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

	buffer, err := pr.file.WriteToBuffer()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error writing to buffer excel: %v", err)
	}

	if err := pr.closeExcel(); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "error closing excel file: %v", err)
	}

	// buffer.Bytes(),

	// if err := pr.file.SaveAs("products_report.xlsx"); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	// }

	return buffer.Bytes(), nil
}

func (pr *productReport) generatePDF() ([]byte, error) {
	if err := pr.addLogo(); err != nil {
		return nil, err
	}

	pr.writeReportMetadata("Products Report", time.Now().Format("2006-01-02"), pr.filters.StartDate.Format("2006-01-02"), pr.filters.EndDate.Format("2006-01-02"))

	pr.pdf.SetFont(fontFamily, "B", largeFont)
	pr.pdf.Cell(0, lineHt, "Summary:")
	pr.pdf.Ln(lineHt)

	summary := map[string]string{
		"TotalProducts": formatQuantity(pr.summary.TotalProducts),
		"MostPopularProduct": pr.summary.MostPopularProduct,
		"MaxLoans": formatQuantity(pr.summary.MaxLoans),
		"TotalActiveLoanAmount": formatQuantity(pr.summary.TotalActiveLoanAmount),
	}

	pr.pdf.SetFontSize(12)
	pr.writeSummary(summary)

	pr.pdf.Ln(lineHt*2)
	headers := []string{"Product Name", "Loans Issued", "Active Loans", "Completed Loans", "Defaulted Loans", "Amount Disbursed", "Amount Repaid", "Outstanding Amount", "Default Rate(%)"}
	colWidths := []float64{45, 25, 25, 32, 32, 34, 30, 37, 30}

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

	var buffer bytes.Buffer
	if err := pr.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}
	// buffer.Bytes()
	// pr.pdf.OutputFileAndClose("products_report.pdf")

	return buffer.Bytes(), nil
}