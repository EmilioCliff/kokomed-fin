package reports

import (
	"bytes"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type paymentReport struct {
	*excelGenerator
	*PDFGenerator
	data []services.PaymentReportData
	summary services.PaymentSummary
	filters services.ReportFilters
}

func newPaymentReport(data []services.PaymentReportData, summary services.PaymentSummary, format string, filters services.ReportFilters) *paymentReport {
	reportGenerator := &paymentReport{
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

func (pr *paymentReport) generateExcel(sheetName string) ([]byte, error) {
	pr.createSheet(sheetName)

	columns := []string{"Transaction Number", "Paying Name", "Amount", "Account Number", "Transaction Source", "Paid Date", "Assigned To", "Assigned By"}

	pr.file.SetColWidth(pr.currentSheet, "A", "H", 20)
	pr.file.SetColStyle(pr.currentSheet, "C", pr.createMoneyStyle())
	pr.file.SetColStyle(pr.currentSheet, "F", pr.createDateStyle())

	pr.writeHeader(columns, pr.createHeaderStyle())

	for rowIdx, payment := range pr.data {
		row := []interface{}{
			payment.TransactionNumber,
			payment.PayingName,
			payment.Amount,
			payment.AccountNumber,
			payment.TransactionSource,
			payment.PaidDate,
			payment.AssignedTo,
			payment.AssignedBy,
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
	
	return buffer.Bytes(), err

	// if err := pr.file.SaveAs("payments_report.xlsx"); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	// }

	// return pr.closeExcel()
}

func (pr *paymentReport) generatePDF() ([]byte, error) {
	if err := pr.addLogo(); err != nil {
		return nil, err
	}

	pr.writeReportMetadata("Payments Report", time.Now().Format("2006-01-02"), pr.filters.StartDate.Format("2006-01-02"), pr.filters.EndDate.Format("2006-01-02"))

	pr.pdf.SetFont(fontFamily, "B", largeFont)
	pr.pdf.Cell(0, lineHt, "Summary:")
	pr.pdf.Ln(lineHt)

	summary := map[string]string {
		"TotalPayments": formatQuantity(pr.summary.TotalPayments), 
		"TotalAmountReceived": formatMoney(pr.summary.TotalAmountReceived), 
		"MostCommonSource": pr.summary.MostCommonSource, 
		"MostAssignedStaff": pr.summary.MostAssignedStaff, 
	}

	pr.pdf.SetFontSize(12)
	pr.writeSummary(summary)

	pr.pdf.Ln(lineHt*2)
	headers := []string{"Txn No", "Paying Name", "Amount", "Account No", "Txn Source", "Date Paid", "Assigned To", "Assigned By"}
	colWidths := []float64{25, 25, 20, 25, 25, 35, 25, 20}

	pr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    pr.pdf.SetFont("Arial", "B", mediumFont)
    pr.pdf.SetX(marginX)
	pr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "L", "CM", "L", "CM", "CM", "L", "L"}

	pr.pdf.SetFontStyle("")
    pr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, payment := range pr.data {
		paidDate := payment.PaidDate.Format("2006/01/02 03:04PM")
		row := []interface{}{
			payment.TransactionNumber,
			payment.PayingName,
			formatMoney(payment.Amount),
			payment.AccountNumber,
			payment.TransactionSource,
			paidDate,
			payment.AssignedTo,
			payment.AssignedBy,
		}
		pr.writeTableRow(row, colWidths, colAlignment)
	}

	var buffer bytes.Buffer
	if err := pr.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}

	pr.closePDF()

	return buffer.Bytes(), nil

	// return pr.pdf.OutputFileAndClose("payments_report.pdf")
}
