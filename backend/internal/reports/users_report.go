package reports

import (
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type userReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.UserAdminsReportData
	userData []services.UserUsersReportData
	filters services.ReportFilters
}

func newUserReport(adminData []services.UserAdminsReportData, userData []services.UserUsersReportData, format string, filters services.ReportFilters) *userReport {
	reportGenerator := &userReport{
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

func (ur *userReport) generateExcel(sheetName string) (error) {
	ur.createSheet(sheetName)

	if len(ur.adminData) > 0 {
		return ur.adminReportExcel()
	} else {
		return ur.userReportPDF()
	}
}

func (ur *userReport) generatePDF() (error) {
	if len(ur.adminData) > 0 {
		return ur.adminReportPDF()
	} else {
		return ur.userReportPDF()
	}
}

func (ur *userReport) adminReportExcel() error {
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

	// buffer, err := ur.file.WriteToBuffer()
	// buffer.Bytes(),

	if err := ur.file.SaveAs("users_admin_report.xlsx"); err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate Excel report")
	}

	return ur.closeExcel()
}

func (ur *userReport) adminReportPDF() error {
	if err := ur.addLogo(); err != nil {
		return err
	}

	ur.writeReportMetadata("Admin User's Report", time.Now().Format("2006-01-02"), ur.filters.StartDate.Format("2006-01-02"), ur.filters.EndDate.Format("2006-01-02"))

	ur.pdf.SetFont(fontFamily, "B", largeFont)
	ur.pdf.Cell(0, lineHt, "Summary:")
	ur.pdf.Ln(lineHt)

	center := ur.getCenterX()
	ur.drawBox(marginX, ur.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Full Name", "Branch Name", "Role", "Clients Registered", "Payments Assigned", "Approved Loans", "Active Loans", "Completed Loans", "Default Rate"}
	colWidths := []float64{35, 35, 25, 34, 34, 32, 30, 32, 30}

	ur.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    ur.pdf.SetFont("Arial", "B", mediumFont)
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

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return ur.pdf.OutputFileAndClose("user_admin_report.pdf")
}

func (ur *userReport) userReportPDF() error {
	log.Println("Generating user report")
	if err := ur.addLogo(); err != nil {
		return err
	}

	ur.writeReportMetadata("uranches Report", time.Now().Format("2006-01-02"), ur.filters.StartDate.Format("2006-01-02"), ur.filters.EndDate.Format("2006-01-02"))

	ur.pdf.SetFont(fontFamily, "B", largeFont)
	ur.pdf.Cell(0, lineHt, "Summary:")
	ur.pdf.Ln(lineHt)

	center := ur.getCenterX()
	ur.drawBox(marginX, ur.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Name", "No Clients", "No Staff", "Loans Issued", "Disbursed Amount", "Collected Amount", "Outstanding Amount", "Default Rate"}
	colWidths := []float64{35, 20, 20, 20, 25, 25, 25, 20}

	ur.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    ur.pdf.SetFont("Arial", "B", mediumFont)
    ur.pdf.SetX(marginX)
	ur.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "CM", "CM", "CM", "R", "R", "R", "CM"}

	ur.pdf.SetFontStyle("")
    ur.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
	for _, user := range ur.userData {
		log.Println(user)
		row := []interface{}{
			// uranch.uranchName,
			// uranch.TotalClients,
			// uranch.TotalUsers,
			// uranch.LoansIssued,
			// formatMoney(uranch.TotalDisbursed),
			// formatMoney(uranch.TotalCollected),
			// formatMoney(uranch.TotalOutstanding),
			// uranch.DefaultRate,
		}
		ur.writeTableRow(row, colWidths, colAlignment)
	}

	// var buffer bytes.Buffer
	// if err := ur.pdf.Output(&buffer); err != nil {
	// 	return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to generate PDF report")
	// }
	// buffer.Bytes()

	return ur.pdf.OutputFileAndClose("uranches_report.pdf")
}