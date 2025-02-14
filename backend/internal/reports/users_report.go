package reports

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type userReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.UserAdminsReportData
	adminSummary services.UserAdminsSummary
	userData services.UserUsersReportData
	filters services.ReportFilters
}

func newUserReport(adminData []services.UserAdminsReportData, userData services.UserUsersReportData, adminSummary services.UserAdminsSummary, format string, filters services.ReportFilters) *userReport {
	reportGenerator := &userReport{
		adminData: adminData,
		adminSummary: adminSummary,
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
		reportGenerator.PDFGenerator = newPDFGenerator("L", "A4")
	}

	return reportGenerator
}

func (ur *userReport) generateExcel(sheetName string) ([]byte,error) {
	ur.createSheet(sheetName)

	if len(ur.adminData) > 0 {
		return ur.adminReportExcel()
	} else {
		return ur.userReportPDF()
	}
}

func (ur *userReport) generatePDF() ([]byte,error) {
	if len(ur.adminData) > 0 {
		return ur.adminReportPDF()
	} else {
		return ur.userReportPDF()
	}
}

func (ur *userReport) adminReportExcel() ([]byte, error) {
	columns := []string{"Fullname", "Branch Name", "Role", "Clients Registered", "Payments Assigned", "Approved Loans", "Active Clients Loans", "Completed Clients Loans", "Default Rate(%)"}

	ur.file.SetColWidth(ur.currentSheet, "A", "I", 20)
	ur.file.SetColStyle(ur.currentSheet, "D", ur.createQuantityStyle())
	ur.file.SetColStyle(ur.currentSheet, "E", ur.createQuantityStyle())
	ur.file.SetColStyle(ur.currentSheet, "F", ur.createQuantityStyle())
	ur.file.SetColStyle(ur.currentSheet, "G", ur.createQuantityStyle())
	ur.file.SetColStyle(ur.currentSheet, "H", ur.createQuantityStyle())
	// ur.file.SetColStyle(ur.currentSheet, "H", ur.createPercentageStyle())

	ur.writeHeader(columns, ur.createHeaderStyle())

	for rowIdx, user := range ur.adminData {
		row := []interface{}{
			user.FullName,
			user.Roles,
			user.BranchName,
			user.ClientsRegistered,
			user.PaymentsAssigned,
			user.ApprovedLoans,
			user.ActiveLoans,
			user.CompletedLoans,
			user.DefaultRate,
		}
		ur.writeRow(rowIdx+2, row)
	}

	buffer, err := ur.file.WriteToBuffer()
	if err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}
	// buffer.Bytes(),

	// if err := ur.file.SaveAs("users_admin_report.xlsx"); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	// }
	if err := ur.closeExcel(); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to close excel file: %v", err)
	}

	return buffer.Bytes(), nil
}

func (ur *userReport) adminReportPDF() ([]byte, error) {
	if err := ur.addLogo(); err != nil {
		return nil, err
	}

	ur.writeReportMetadata("Admin User's Report", time.Now().Format("2006-01-02"), ur.filters.StartDate.Format("2006-01-02"), ur.filters.EndDate.Format("2006-01-02"))

	ur.pdf.SetFont(fontFamily, "B", largeFont)
	ur.pdf.Cell(0, lineHt, "Summary:")
	ur.pdf.Ln(lineHt)

	summary := map[string]string {
		"TotalUsers": formatQuantity(ur.adminSummary.TotalUsers),
		"TotalClients": formatQuantity(ur.adminSummary.TotalClients),
		"TotalPayments": formatQuantity(ur.adminSummary.TotalPayments),
		"HighestLoanApprovalUser": ur.adminSummary.HighestLoanApprovalUser,
	}

	ur.pdf.SetFontSize(12)
	ur.writeSummary(summary)

	ur.pdf.Ln(lineHt*2)
	headers := []string{"Full Name", "Branch Name", "Role", "Clients Registered", "Payments Assigned", "Approved Loans", "Active Loans", "Completed Loans", "Default Rate"}
	colWidths := []float64{35, 35, 25, 34, 34, 32, 30, 32, 30}

	ur.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    ur.pdf.SetFont(fontFamily, "B", mediumFont)
    ur.pdf.SetX(marginX)
	ur.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "L", "CM", "CM", "CM", "CM", "CM", "CM", "CM"}

	ur.pdf.SetFontStyle("")
    ur.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, user := range ur.adminData {
		row := []interface{}{
			user.FullName,
			user.BranchName,
			user.Roles,
			user.ClientsRegistered,
			user.PaymentsAssigned,
			user.ApprovedLoans,
			user.ActiveLoans,
			user.CompletedLoans,
			user.DefaultRate,
		}
		ur.writeTableRow(row, colWidths, colAlignment)
	}

	var buffer bytes.Buffer
	if err := ur.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}
	// buffer.Bytes()
	// ur.pdf.OutputFileAndClose("user_admin_report.pdf")
	ur.closePDF()

	return buffer.Bytes(), nil
}

func (ur *userReport) userReportPDF() ([]byte, error) {
	if err := ur.addLogo(); err != nil {
		return nil, err
	}

	ur.writeReportMetadata(fmt.Sprintf("%s Report", ur.userData.Name), time.Now().Format("2006-01-02"), ur.filters.StartDate.Format("2006-01-02"), ur.filters.EndDate.Format("2006-01-02"))

	ur.pdf.SetFont(fontFamily, "B", largeFont)
	ur.pdf.Cell(0, lineHt, "Summary:")
	ur.pdf.Ln(lineHt)

	summary := map[string]string {
		"Name": ur.userData.Name,
		"Role": ur.userData.Role,
		"Branch": ur.userData.Branch,
		"TotalClientsHandled": formatQuantity(ur.userData.TotalClientsHandled),
		"LoansApproved": formatQuantity(ur.userData.LoansApproved),
		"TotalLoanAmountManaged": formatMoney(ur.userData.TotalLoanAmountManaged),
		"TotalCollectedAmount": formatMoney(ur.userData.TotalCollectedAmount),
		"DefaultRate": fmt.Sprintf("%.2f%%", ur.userData.DefaultRate),
		"AssignedPayments": formatQuantity(ur.userData.AssignedPayments),
	}

	ur.pdf.SetFontSize(12)
	ur.writeSummary(summary)

	ur.pdf.Ln(lineHt*2)

	ur.pdf.SetFont(fontFamily, "B", largeFont)
	ur.pdf.Cell(10, lineHt, "Assigned Loans")
	ur.pdf.Ln(-1)
	headers := []string{"LoanId", "ClientName", "Status", "LoanAmount", "RepayAmount", "PaidAmount", "DisbursedOn"}
	colWidths := []float64{30, 35, 30, 35, 35, 35, 35}

	ur.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    ur.pdf.SetFont(fontFamily, "B", mediumFont)
    ur.pdf.SetX(marginX)
	colAlignment := []string{"CM", "L", "CM", "R", "R", "R", "R"}
	if len(ur.userData.AssignedLoans) > 0 {
		ur.writeTableHeaders(headers, colWidths)
	
		ur.pdf.SetFontStyle("")
		ur.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
		for _, loans := range ur.userData.AssignedLoans {
			row := []interface{}{
				fmt.Sprintf("LN%03d", loans.LoanId),
				loans.ClientName,
				loans.Status,
				formatMoney(loans.LoanAmount),
				formatMoney(loans.RepayAmount),
				formatMoney(loans.PaidAmount),
				loans.DisbursedOn,
			}
			ur.writeTableRow(row, colWidths, colAlignment)
		}
	}

	if len(ur.userData.AssignedPaymentsList) > 0 {
		ur.pdf.Ln(lineHt*4)
		ur.pdf.SetFont(fontFamily, "B", largeFont)
		ur.pdf.Cell(10, lineHt, "Assigned Payments")
		ur.pdf.Ln(-1)
		headers = []string{"TransactionNumber", "ClientName", "AmountPaid", "PaidDate"}
		colWidths = []float64{40, 35, 35, 35}
	
		ur.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
		ur.pdf.SetFont(fontFamily, "B", mediumFont)
		ur.pdf.SetX(marginX)
		ur.writeTableHeaders(headers, colWidths)
		colAlignment = []string{"CM", "L", "R", "R"}
	
		ur.pdf.SetFontStyle("")
		ur.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
		for _, payment := range ur.userData.AssignedPaymentsList {
			paidDate := strings.Split(payment.PaidDate, " ")[0]
			row := []interface{}{
				payment.TransactionNumber,
				payment.ClientName,
				formatMoney(payment.AmountPaid),
				paidDate,
			}
			ur.writeTableRow(row, colWidths, colAlignment)
		}
	}

	var buffer bytes.Buffer
	if err := ur.pdf.Output(&buffer); err != nil {
		return nil, pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	}

	ur.closePDF()
	// buffer.Bytes()
	// ur.pdf.OutputFileAndClose(fmt.Sprintf("%s_report.pdf", ur.userData.Name))

	return buffer.Bytes(), nil
}