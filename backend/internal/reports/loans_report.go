package reports

import (
	"fmt"
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type loanReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.LoanReportData
	userData services.LoanReportDataById
	filters services.ReportFilters
}

func newLoanReport(adminData []services.LoanReportData, userData services.LoanReportDataById, format string, filters services.ReportFilters) *loanReport {
	reportGenerator := &loanReport{
		adminData: adminData,
		userData: userData,
		filters: filters,
	}

	switch format {
	case "excel":
		if len(adminData) > 0 {
			reportGenerator.excelGenerator = newExcelGenerator()
		} else {
			reportGenerator.PDFGenerator = newPDFGenerator()
		}
	case "pdf":
		reportGenerator.PDFGenerator = newPDFGenerator()
	}

	return reportGenerator
}

func (lr *loanReport) generateExcel(sheetName string) (error) {
	lr.createSheet(sheetName)

	if len(lr.adminData) > 0 {
		return lr.adminReportExcel()
	} else {
		return lr.loanReportPDF()
	}
}

func (lr *loanReport) generatePDF() (error) {
	if len(lr.adminData) > 0 {
		return lr.adminReportPDF()
	} else {
		return lr.loanReportPDF()
	}
}

func (lr *loanReport) adminReportExcel() error {
	columns := []string{"Loan ID", "Client Name", "Branch Name", "Loan Officer", "Loan Amount", "Repay Amount", "Paid Amount", "Outstanding Amount", "Status", "Due Date", "No Installments", "Paid Installments", "Disbursed Date", "Default Risk(%)"}

	lr.file.SetColWidth(lr.currentSheet, "A", "N", 20)
	lr.file.SetColStyle(lr.currentSheet, "E", lr.createMoneyStyle())
	lr.file.SetColStyle(lr.currentSheet, "F", lr.createMoneyStyle())
	lr.file.SetColStyle(lr.currentSheet, "G", lr.createMoneyStyle())
	lr.file.SetColStyle(lr.currentSheet, "H", lr.createMoneyStyle())
	lr.file.SetColStyle(lr.currentSheet, "J", lr.createDateStyle())
	lr.file.SetColStyle(lr.currentSheet, "K", lr.createQuantityStyle())
	lr.file.SetColStyle(lr.currentSheet, "L", lr.createQuantityStyle())
	lr.file.SetColStyle(lr.currentSheet, "M", lr.createDateStyle())
	// ur.file.SetColStyle(ur.currentSheet, "H", ur.createPercentageStyle())

	lr.writeHeader(columns, lr.createHeaderStyle())

	for rowIdx, loan := range lr.adminData {
		row := []interface{}{
			fmt.Sprintf("LN%03d", loan.LoanID),
			loan.ClientName,
			loan.BranchName,
			loan.LoanOfficer,
			loan.LoanAmount,
			loan.RepayAmount,
			loan.PaidAmount,
			loan.OutstandingAmount,
			loan.Status,
			formatTime(loan.DueDate),
			loan.TotalInstallments,
			loan.PaidInstallments,
			formatTime(loan.DisbursedDate),
			loan.DefaultRisk,
		}
		lr.writeRow(rowIdx+2, row)
	}

	// buffer, err := ur.file.WriteToBuffer()
	// buffer.Bytes(),

	if err := lr.file.SaveAs("loans_report.xlsx"); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	return lr.closeExcel()
}

func (lr *loanReport) adminReportPDF() error {
	if err := lr.addLogo(); err != nil {
		return err
	}

	lr.writeReportMetadata("Admin Clients Report", time.Now().Format("2006-01-02"), lr.filters.StartDate.Format("2006-01-02"), lr.filters.EndDate.Format("2006-01-02"))

	lr.pdf.SetFont(fontFamily, "B", largeFont)
	lr.pdf.Cell(0, lineHt, "Summary:")
	lr.pdf.Ln(lineHt)

	center := lr.getCenterX()
	lr.drawBox(marginX, lr.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"LoanID", "Client Name", "Branch Name", "Loan Officer", "Loan Amount", "Repay Amount", "Paid Amount", "Outstanding Amount", "Status", "Due Date", "No Installments", "Paid Installments", "Disbursed Date", "Default Risk(%)"}
	colWidths := []float64{35, 35, 25, 34, 34, 32, 30, 32, 30, 30, 30, 30, 30, 30}

	lr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    lr.pdf.SetFont("Arial", "B", mediumFont)
    lr.pdf.SetX(marginX)
	lr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"CM", "L", "L", "L", "R", "R", "R", "R", "CM", "R", "CM", "CM", "R", "CM"}

	lr.pdf.SetFontStyle("")
    lr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, loan := range lr.adminData {
		row := []interface{}{
			fmt.Sprintf("LN%03d", loan.LoanID),
			loan.ClientName,
			loan.BranchName,
			loan.LoanOfficer,
			formatMoney(loan.LoanAmount),
			formatMoney(loan.RepayAmount),
			formatMoney(loan.PaidAmount),
			formatMoney(loan.OutstandingAmount),
			loan.Status,
			loan.DueDate.Format("2006/01/02"),
			loan.TotalInstallments,
			loan.PaidInstallments,
			loan.DisbursedDate.Format("2006/01/02"),
			loan.DefaultRisk,
		}
		lr.writeTableRow(row, colWidths, colAlignment)
	}

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return lr.pdf.OutputFileAndClose("loans_report.pdf")
}

func (lr *loanReport) loanReportPDF() error {
	log.Println("Generating user report")
	if err := lr.addLogo(); err != nil {
		return err
	}

	lr.writeReportMetadata("Client Report", time.Now().Format("2006-01-02"), lr.filters.StartDate.Format("2006-01-02"), lr.filters.EndDate.Format("2006-01-02"))

	lr.pdf.SetFont(fontFamily, "B", largeFont)
	lr.pdf.Cell(0, lineHt, "Summary:")
	lr.pdf.Ln(lineHt)

	center := lr.getCenterX()
	lr.drawBox(marginX, lr.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Name", "No Clients", "No Staff", "Loans Issued", "Disbursed Amount", "Collected Amount", "Outstanding Amount", "Default Rate"}
	colWidths := []float64{35, 20, 20, 20, 25, 25, 25, 20}

	lr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    lr.pdf.SetFont("Arial", "B", mediumFont)
    lr.pdf.SetX(marginX)
	lr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "CM", "CM", "CM", "R", "R", "R", "CM"}
	log.Println(colAlignment)

	lr.pdf.SetFontStyle("")
    lr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	// for _, user := range cr.userData {
	// 	log.Println(user)
	// 	row := []interface{}{
			// uranch.uranchName,
			// uranch.TotalClients,
			// uranch.TotalUsers,
			// uranch.LoansIssued,
			// formatMoney(uranch.TotalDisbursed),
			// formatMoney(uranch.TotalCollected),
			// formatMoney(uranch.TotalOutstanding),
			// uranch.DefaultRate,
	// 	}
	// 	cr.writeTableRow(row, colWidths, colAlignment)
	// }

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return lr.pdf.OutputFileAndClose("clients_report.pdf")
}