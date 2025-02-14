package reports

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type loanReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.LoanReportData
	adminSummary services.LoanSummary
	userData services.LoanReportDataById
	filters services.ReportFilters
}

func newLoanReport(adminData []services.LoanReportData, userData services.LoanReportDataById, summary services.LoanSummary, format string, filters services.ReportFilters) *loanReport {
	reportGenerator := &loanReport{
		adminData: adminData,
		adminSummary: summary,
		userData: userData,
		filters: filters,
	}

	switch format {
	case "excel":
		if len(adminData) > 0 {
			reportGenerator.excelGenerator = newExcelGenerator()
		} else {
			reportGenerator.PDFGenerator = newPDFGenerator("L", "A4")
		}
	case "pdf":
		if len(adminData) > 0 {
			reportGenerator.PDFGenerator = newPDFGenerator("L", "A3")
		} else {
			reportGenerator.PDFGenerator = newPDFGenerator("L", "A4")
		}
	}

	return reportGenerator
}

func (lr *loanReport) generateExcel(sheetName string) ([]byte, error) {
	lr.createSheet(sheetName)

	if len(lr.adminData) > 0 {
		return lr.adminReportExcel()
	} else {
		return lr.loanReportPDF()
	}
}

func (lr *loanReport) generatePDF() ([]byte, error) {
	if len(lr.adminData) > 0 {
		return lr.adminReportPDF()
	} else {
		return lr.loanReportPDF()
	}
}

func (lr *loanReport) adminReportExcel() ([]byte, error) {
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
			loan.DueDate,
			loan.TotalInstallments,
			loan.PaidInstallments,
			loan.DisbursedDate,
			loan.DefaultRisk,
		}
		lr.writeRow(rowIdx+2, row)
	}

	buffer, err := lr.file.WriteToBuffer()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	if err := lr.closeExcel(); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to close excel file: %v", err)
	}
	// buffer.Bytes(),

	// if err := lr.file.SaveAs("loans_report.xlsx"); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	// }

	return buffer.Bytes(), nil
}

func (lr *loanReport) adminReportPDF() ([]byte, error) {
	if err := lr.addLogo(); err != nil {
		return nil, err
	}

	lr.writeReportMetadata("Admin Clients Report", time.Now().Format("2006-01-02"), lr.filters.StartDate.Format("2006-01-02"), lr.filters.EndDate.Format("2006-01-02"))

	lr.pdf.SetFont(fontFamily, "B", largeFont)
	lr.pdf.Cell(0, lineHt, "Summary:")
	lr.pdf.Ln(lineHt)

	summary := map[string]string {
		"TotalLoans": fmt.Sprintf("LN%03d", lr.adminSummary.TotalLoans),
		"TotalActiveLoans": formatQuantity(lr.adminSummary.TotalActiveLoans),
		"TotalCompletedLoans": formatQuantity(lr.adminSummary.TotalCompletedLoans),
		"TotalDefaultedLoans": formatQuantity(lr.adminSummary.TotalDefaultedLoans),
		"TotalDisbursedAmount": formatMoney(lr.adminSummary.TotalDisbursedAmount),
		"TotalRepaidAmount": formatMoney(lr.adminSummary.TotalRepaidAmount),
		"TotalOutstanding": formatMoney(lr.adminSummary.TotalOutstanding),
		"MostIssuedLoanBranch": lr.adminSummary.MostIssuedLoanBranch,
		"MostLoansOfficer": lr.adminSummary.MostLoansOfficer,
	}

	lr.pdf.SetFontSize(12)
	lr.writeSummary(summary)

	lr.pdf.Ln(lineHt*2)
	headers := []string{"LoanID", "Client Name", "Branch Name", "Loan Officer", "Loan Amount", "Repay Amount", "Paid Amount", "Outstanding Amount", "Status", "Due Date", "No Installments", "Paid Installments", "Disbursed Date", "Default Risk(%)"}
	colWidths := []float64{25, 30, 30, 30, 25, 32, 30, 35, 30, 25, 30, 30, 30, 28}

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
			loan.DueDate,
			loan.TotalInstallments,
			loan.PaidInstallments,
			loan.DisbursedDate,
			loan.DefaultRisk,
		}
		lr.writeTableRow(row, colWidths, colAlignment)
	}

	var buffer bytes.Buffer
	if err := lr.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}

	lr.closePDF()
	// buffer.Bytes()

	// lr.pdf.OutputFileAndClose("loans_report.pdf")

	return buffer.Bytes(), nil
}

func (lr *loanReport) loanReportPDF() ([]byte, error) {
	if err := lr.addLogo(); err != nil {
		return nil, err
	}

	lr.writeReportMetadata(fmt.Sprintf("Loan(LN%03d) Report", lr.userData.LoanID), time.Now().Format("2006-01-02"), lr.filters.StartDate.Format("2006-01-02"), lr.filters.EndDate.Format("2006-01-02"))

	lr.pdf.SetFont(fontFamily, "B", largeFont)
	lr.pdf.Cell(0, lineHt, "Summary:")
	lr.pdf.Ln(lineHt)

	summary := map[string]string {
		"LoanID": fmt.Sprintf("LN%03d", lr.userData.LoanID),
		"ClientName": lr.userData.ClientName,
		"LoanAmount": formatMoney(lr.userData.LoanAmount),
		"RepayAmount": formatMoney(lr.userData.RepayAmount),
		"PaidAmount": formatMoney(lr.userData.PaidAmount),
		"Status": lr.userData.Status,
		"TotalInstallments": formatQuantity(lr.userData.TotalInstallments),
		"PaidInstallments": formatQuantity(lr.userData.PaidInstallments),
		"RemainingInstallments": formatQuantity(lr.userData.RemainingInstallments),
	}

	lr.pdf.SetFontSize(12)
	lr.writeSummary(summary)

	lr.pdf.Ln(lineHt*2)

	
	if len(lr.userData.InstallmentDetails) > 0 {
		lr.pdf.SetFont(fontFamily, "B", largeFont)
		lr.pdf.Cell(10, lineHt, "Loan Installments")
		lr.pdf.Ln(-1)
		headers := []string{"InstallmentNumber", "InstallmentAmount", "RemainingAmount", "DueDate", "Paid", "PaidAt"}
		colWidths := []float64{40, 35, 35, 35, 25, 30}
	
		lr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
		lr.pdf.SetFont("Arial", "B", mediumFont)
		lr.pdf.SetX(marginX)
		lr.writeTableHeaders(headers, colWidths)
		colAlignment := []string{"CM", "R", "R", "R", "CM", "R"}

		lr.pdf.SetFontStyle("")
		lr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
		for _, installment := range lr.userData.InstallmentDetails {
			paidDate := strings.Split(installment.PaidAt, " ")[0]
			dueDate := strings.Split(installment.DueDate, " ")[0]
			paid := "PAID"
			if installment.Paid <= 0 {
				paid = "UNPAID"
			}
			row := []interface{}{
				installment.InstallmentNumber,
				formatMoney(installment.InstallmentAmount),
				formatMoney(installment.RemainingAmount),
				dueDate,
				paid,
				paidDate,
			}
			lr.writeTableRow(row, colWidths, colAlignment)
		}

	}
	lr.pdf.SetFontStyle("")
    lr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])

	var buffer bytes.Buffer
	if err := lr.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}
	lr.closePDF()
	// buffer.Bytes()

	// lr.pdf.OutputFileAndClose(fmt.Sprintf("LN%0d_report.pdf", lr.userData.LoanID))

	return buffer.Bytes(), nil
}