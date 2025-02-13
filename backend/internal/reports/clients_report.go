package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type clientReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.ClientAdminsReportData
	adminSummary services.ClientSummary
	userData services.ClientClientsReportData
	filters services.ReportFilters
}

func newClientReport(adminData []services.ClientAdminsReportData, userData services.ClientClientsReportData, summary services.ClientSummary, format string, filters services.ReportFilters) *clientReport {
	reportGenerator := &clientReport{
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

func (cr *clientReport) generateExcel(sheetName string) (error) {
	cr.createSheet(sheetName)

	if len(cr.adminData) > 0 {
		return cr.adminReportExcel()
	} else {
		return cr.clientReportPDF()
	}
}

func (cr *clientReport) generatePDF() (error) {
	if len(cr.adminData) > 0 {
		return cr.adminReportPDF()
	} else {
		return cr.clientReportPDF()
	}
}

func (cr *clientReport) adminReportExcel() error {
	columns := []string{"Name", "Branch Name", "Phone Number", "Loan Issued", "Defaulted Loans", "Active Loans", "Completed Loans", "Inactive Loans", "Total Paid", "Total Disbursed", "Total Owed", "Overpayment", "Rate Score(%)", "Default Rate(%)"}

	cr.file.SetColWidth(cr.currentSheet, "A", "N", 20)
	cr.file.SetColStyle(cr.currentSheet, "D", cr.createQuantityStyle())
	cr.file.SetColStyle(cr.currentSheet, "E", cr.createQuantityStyle())
	cr.file.SetColStyle(cr.currentSheet, "F", cr.createQuantityStyle())
	cr.file.SetColStyle(cr.currentSheet, "G", cr.createQuantityStyle())
	cr.file.SetColStyle(cr.currentSheet, "H", cr.createQuantityStyle())
	cr.file.SetColStyle(cr.currentSheet, "I", cr.createMoneyStyle())
	cr.file.SetColStyle(cr.currentSheet, "J", cr.createMoneyStyle())
	cr.file.SetColStyle(cr.currentSheet, "K", cr.createMoneyStyle())
	cr.file.SetColStyle(cr.currentSheet, "L", cr.createMoneyStyle())
	cr.file.SetColStyle(cr.currentSheet, "M", cr.createMoneyStyle())
	// ur.file.SetColStyle(ur.currentSheet, "H", ur.createPercentageStyle())

	cr.writeHeader(columns, cr.createHeaderStyle())

	for rowIdx, client := range cr.adminData {
		row := []interface{}{
			client.Name,
			client.BranchName,
			client.PhoneNumber,
			client.TotalLoanGiven,
			client.DefaultedLoans,
			client.ActiveLoans,
			client.CompletedLoans,
			client.InactiveLoans,
			client.TotalPaid,
			client.TotalDisbursed,
			client.TotalOwed,
			client.Overpayment,
			client.RateScore,
			client.DefaultRate,
		}
		cr.writeRow(rowIdx+2, row)
	}

	// buffer, err := ur.file.WriteToBuffer()
	// buffer.Bytes(),

	if err := cr.file.SaveAs("clients_admin_report.xlsx"); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	return cr.closeExcel()
}

func (cr *clientReport) adminReportPDF() error {
	if err := cr.addLogo(); err != nil {
		return err
	}

	cr.writeReportMetadata("Admin Clients Report", time.Now().Format("2006-01-02"), cr.filters.StartDate.Format("2006-01-02"), cr.filters.EndDate.Format("2006-01-02"))

	cr.pdf.SetFont(fontFamily, "B", largeFont)
	cr.pdf.Cell(0, lineHt, "Summary:")
	cr.pdf.Ln(lineHt)

	summary := map[string]string {
		"TotalClients": formatQuantity(cr.adminSummary.TotalClients),
		"MostLoansClient": cr.adminSummary.MostLoansClient,
		"HighestPayingClient": cr.adminSummary.HighestPayingClient,
		"TotalDisbursed": formatMoney(cr.adminSummary.TotalDisbursed),
		"TotalPaid": formatMoney(cr.adminSummary.TotalPaid),
		"TotalOwed": formatMoney(cr.adminSummary.TotalOwed),
	}

	cr.pdf.SetFontSize(12)
	cr.writeSummary(summary)

	cr.pdf.Ln(lineHt*2)
	headers := []string{"Name", "Branch Name", "Phone Number", "Loan Issued", "Defaulted Loans", "Active Loans", "Completed Loans", "Inactive Loans", "Total Paid", "Total Disbursed", "Total Owed", "Overpayment", "Rate Score(%)", "Default Rate(%)"}
	colWidths := []float64{35, 30, 30, 25, 30, 25, 33, 30, 26, 30, 27, 30, 30, 30}

	cr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    cr.pdf.SetFont("Arial", "B", mediumFont)
    cr.pdf.SetX(marginX)
	cr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "L", "CM", "CM", "CM", "CM", "CM", "R", "R", "R", "R", "R", "R", "R"}

	cr.pdf.SetFontStyle("")
    cr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, client := range cr.adminData {
		row := []interface{}{
			client.Name,
			client.BranchName,
			client.PhoneNumber,
			client.TotalLoanGiven,
			client.DefaultedLoans,
			client.ActiveLoans,
			client.CompletedLoans,
			client.InactiveLoans,
			formatMoney(client.TotalPaid),
			formatMoney(client.TotalDisbursed),
			formatMoney(client.TotalOwed),
			formatMoney(client.Overpayment),
			fmt.Sprintf("%.2f", client.RateScore),
			fmt.Sprintf("%.2f", client.DefaultRate),
		}
		cr.writeTableRow(row, colWidths, colAlignment)
	}

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return cr.pdf.OutputFileAndClose("client_admin_report.pdf")
}

func (cr *clientReport) clientReportPDF() error {
	if err := cr.addLogo(); err != nil {
		return err
	}

	cr.writeReportMetadata(fmt.Sprintf("%s Report", cr.userData.Name), time.Now().Format("2006-01-02"), cr.filters.StartDate.Format("2006-01-02"), cr.filters.EndDate.Format("2006-01-02"))

	cr.pdf.SetFont(fontFamily, "B", largeFont)
	cr.pdf.Cell(0, lineHt, "Summary:")
	cr.pdf.Ln(lineHt)

	summary := map[string]string {
		"Name": cr.userData.Name,
		"PhoneNumber": cr.userData.PhoneNumber,
		"IDNumber": *cr.userData.IDNumber,
		"Dob": formatTime(cr.userData.Dob),
		"BranchName": cr.userData.BranchName,
		"AssignedStaff": cr.userData.AssignedStaff,
		"Active": cr.userData.Active,
	}

	cr.pdf.SetFontSize(12)
	cr.writeSummary(summary)

	cr.pdf.Ln(lineHt*2)

	cr.pdf.SetFont(fontFamily, "B", largeFont)
	cr.pdf.Cell(10, lineHt, "User Loans")
	cr.pdf.Ln(-1)
	headers := []string{"LoanId", "Status", "LoanAmount", "RepayAmount", "PaidAmount", "DisbursedOn", "TransactionFee", "ApprovedBy"}
	colWidths := []float64{35, 30, 35, 35, 35, 30, 30, 35}

	cr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    cr.pdf.SetFont("Arial", "B", mediumFont)
    cr.pdf.SetX(marginX)
	colAlignment := []string{"CM", "CM", "R", "R", "R", "R", "CM", "L"}
	
	if len(cr.userData.Loans) > 0 {
		cr.writeTableHeaders(headers, colWidths)

		cr.pdf.SetFontStyle("")
		cr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
		for _, loan := range cr.userData.Loans {
			transactionFee := "PAID"
			if loan.TransactionFee <= 0 {
				transactionFee = "UNPAID"
			}
			row := []interface{}{
				fmt.Sprintf("LN%0d", loan.LoanId),
				loan.Status,
				formatMoney(loan.LoanAmount),
				formatMoney(loan.RepayAmount),
				formatMoney(loan.PaidAmount),
				loan.DisbursedOn,
				transactionFee,
				loan.ApprovedBy,
			}
			cr.writeTableRow(row, colWidths, colAlignment)
		}
	}

	if len(cr.userData.Payments) > 0 {
		cr.pdf.Ln(lineHt*4)
		cr.pdf.SetFont(fontFamily, "B", largeFont)
		cr.pdf.Cell(10, lineHt, "User Payments")
		cr.pdf.Ln(-1)
		headers = []string{"TransactionNumber", "TransactionSource", "AccountNumber", "PayingName", "AssignedBy", "AmountPaid", "PaidDate"}
		colWidths = []float64{40, 35, 35, 35, 35, 35, 35}
	
		cr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
		cr.pdf.SetFont(fontFamily, "B", mediumFont)
		cr.pdf.SetX(marginX)
		cr.writeTableHeaders(headers, colWidths)
		colAlignment = []string{"CM", "CM", "L", "L", "CM", "R", "R"}
	
		cr.pdf.SetFontStyle("")
		cr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
		for _, payment := range cr.userData.Payments {
			paidDate := strings.Split(payment.PaidDate, " ")[0]
			row := []interface{}{
				payment.TransactionNumber,
				payment.TransactionSource,
				payment.AccountNumber,
				payment.PayingName,
				payment.AssignedBy,
				formatMoney(payment.AmountPaid),
				paidDate,
			}
			cr.writeTableRow(row, colWidths, colAlignment)
		}
	}

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return cr.pdf.OutputFileAndClose(fmt.Sprintf("%s_report.pdf", cr.userData.Name))
}