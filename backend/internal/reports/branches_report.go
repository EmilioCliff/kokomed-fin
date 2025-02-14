package reports

import (
	"bytes"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type branchReport struct {
	*excelGenerator
	*PDFGenerator
	data []services.BranchReportData
	summary services.BranchSummary
	filters services.ReportFilters
}

func newBranchReport(data []services.BranchReportData, summary services.BranchSummary, format string, filters services.ReportFilters) *branchReport {
	reportGenerator := &branchReport{
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

func (br *branchReport) generateExcel(sheetName string) ([]byte, error) {
	br.createSheet(sheetName)

	columns := []string{"Branch Name", "Total Clients", "Total Staff", "Loans Issued", "Total Disbursed Amount", "Total Collected Amount", "Total Outstanding Amount", "Default Rate(%)"}

	br.file.SetColWidth(br.currentSheet, "A", "H", 20)
	br.file.SetColStyle(br.currentSheet, "B", br.createQuantityStyle())
	br.file.SetColStyle(br.currentSheet, "C", br.createQuantityStyle())
	br.file.SetColStyle(br.currentSheet, "D", br.createQuantityStyle())
	br.file.SetColStyle(br.currentSheet, "E", br.createMoneyStyle())
	br.file.SetColStyle(br.currentSheet, "F", br.createMoneyStyle())
	br.file.SetColStyle(br.currentSheet, "G", br.createMoneyStyle())

	br.writeHeader(columns, br.createHeaderStyle())

	for rowIdx, branch := range br.data {
		row := []interface{}{
			branch.BranchName,
			branch.TotalClients,
			branch.TotalUsers,
			branch.LoansIssued,
			branch.TotalDisbursed,
			branch.TotalCollected,
			branch.TotalOutstanding,
			branch.DefaultRate,
		}
		br.writeRow(rowIdx+2, row)
	}

	buffer, err := br.file.WriteToBuffer()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	if err := br.closeExcel(); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to close excel file: %v", err)
	}

	// buffer.Bytes(),

	// if err := br.file.SaveAs("branches_report.xlsx"); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	// }

	return buffer.Bytes(), nil
}

func (br *branchReport) generatePDF() ([]byte, error) {
	if err := br.addLogo(); err != nil {
		return nil, err
	}

	br.writeReportMetadata("Branches Report", time.Now().Format("2006-01-02"), br.filters.StartDate.Format("2006-01-02"), br.filters.EndDate.Format("2006-01-02"))

	br.pdf.SetFont(fontFamily, "B", largeFont)
	br.pdf.Cell(0, lineHt, "Summary:")
	br.pdf.Ln(lineHt)

	summary := map[string]string{
		"TotalBranches":  formatQuantity(br.summary.TotalBranches),
		"HighestPerformingBranch":  br.summary.HighestPerformingBranch,
		"MostClientsBranch":  br.summary.MostClientsBranch,
		"TotalClients":  formatQuantity(br.summary.TotalClients),
		"TotalUsers":  formatQuantity(br.summary.TotalUsers),
	}

	br.pdf.SetFontSize(12)
	br.writeSummary(summary)

	br.pdf.Ln(lineHt*2)

	headers := []string{"Name", "No Clients", "No Staff", "Loans Issued", "Disbursed Amount", "Collected Amount", "Outstanding Amount", "Default Rate"}
	colWidths := []float64{35, 25, 25, 35, 40, 40, 40, 25}

	br.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    br.pdf.SetFont("Arial", "B", mediumFont)
    br.pdf.SetX(marginX)
	br.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "CM", "CM", "CM", "R", "R", "R", "CM"}

	br.pdf.SetFontStyle("")
    br.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, branch := range br.data {
		row := []interface{}{
			branch.BranchName,
			branch.TotalClients,
			branch.TotalUsers,
			branch.LoansIssued,
			formatMoney(branch.TotalDisbursed),
			formatMoney(branch.TotalCollected),
			formatMoney(branch.TotalOutstanding),
			branch.DefaultRate,
		}
		br.writeTableRow(row, colWidths, colAlignment)
	}

	var buffer bytes.Buffer
	if err := br.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}

	br.closePDF()
	// buffer.Bytes()
	// br.pdf.OutputFileAndClose("branches_report.pdf")

	return buffer.Bytes(), nil
}