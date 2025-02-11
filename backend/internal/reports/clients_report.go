package reports

import (
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

type clientReport struct {
	*excelGenerator
	*PDFGenerator
	adminData []services.ClientAdminsReportData
	userData services.ClientClientsReportData
	filters services.ReportFilters
}

func newClientReport(adminData []services.ClientAdminsReportData, userData services.ClientClientsReportData, format string, filters services.ReportFilters) *clientReport {
	reportGenerator := &clientReport{
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

	center := cr.getCenterX()
	cr.drawBox(marginX, cr.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Name", "Branch Name", "Phone Number", "Loan Issued", "Defaulted Loans", "Active Loans", "Completed Loans", "Inactive Loans", "Total Paid", "Total Disbursed", "Total Owed", "Overpayment", "Rate Score(%)", "Default Rate(%)"}
	colWidths := []float64{35, 35, 25, 34, 34, 32, 30, 32, 30, 30, 30, 30, 30, 30}

	cr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    cr.pdf.SetFont("Arial", "B", mediumFont)
    cr.pdf.SetX(marginX)
	cr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "L", "CM", "CM", "CM", "CM", "CM", "CM", "CM", "CM", "CM", "CM", "CM", "CM"}

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
			client.TotalPaid,
			client.TotalDisbursed,
			client.TotalOwed,
			client.Overpayment,
			client.RateScore,
			client.DefaultRate,
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
	log.Println("Generating client report")
	if err := cr.addLogo(); err != nil {
		return err
	}

	cr.writeReportMetadata("Client Report", time.Now().Format("2006-01-02"), cr.filters.StartDate.Format("2006-01-02"), cr.filters.EndDate.Format("2006-01-02"))

	cr.pdf.SetFont(fontFamily, "B", largeFont)
	cr.pdf.Cell(0, lineHt, "Summary:")
	cr.pdf.Ln(lineHt)

	center := cr.getCenterX()
	cr.drawBox(marginX, cr.pdf.GetY(), center/2, lineHt*3, secondaryColor)

	headers := []string{"Name", "No Clients", "No Staff", "Loans Issued", "Disbursed Amount", "Collected Amount", "Outstanding Amount", "Default Rate"}
	colWidths := []float64{35, 20, 20, 20, 25, 25, 25, 20}

	cr.pdf.SetFillColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
    cr.pdf.SetFont("Arial", "B", mediumFont)
    cr.pdf.SetX(marginX)
	cr.writeTableHeaders(headers, colWidths)
	colAlignment := []string{"L", "CM", "CM", "CM", "R", "R", "R", "CM"}
	log.Println(colAlignment)

	cr.pdf.SetFontStyle("")
    cr.pdf.SetFillColor(primaryColor[0], primaryColor[1], primaryColor[2])
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

	return cr.pdf.OutputFileAndClose("clients_report.pdf")
}