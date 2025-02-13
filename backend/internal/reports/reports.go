package reports

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"golang.org/x/text/message"
)

var _ services.ReportService = (*ReportServiceImpl)(nil)

func NewReportService(store *mysql.MySQLRepo) services.ReportService {
	return &ReportServiceImpl{
		store: store,
	}
}

type ReportServiceImpl struct {
	store *mysql.MySQLRepo
}

func (r *ReportServiceImpl) GeneratePaymentsReport(ctx context.Context, format string, filters services.ReportFilters) ([]byte, error){
	data, summary, err := r.store.NonPosted.GetReportPaymentData(ctx, filters)
	if err != nil {
		return nil, err
	}

	report := newPaymentReport(data, summary, format, filters)

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func (r *ReportServiceImpl) GenerateBranchesReport(ctx context.Context, format string, filters services.ReportFilters) ([]byte, error) {
	data, summary, err := r.store.Branches.GetReportBranchData(ctx, filters)
	if err != nil {
		return nil, err
	}

	report := newBranchReport(data, summary, format, filters)

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func (r *ReportServiceImpl) GenerateProductsReport(ctx context.Context, format string, filters services.ReportFilters) ([]byte, error) {
	data, summary, err := r.store.Products.GetReportProductData(ctx, filters)
	if err != nil {
		return nil, err
	}

	report := newProductReport(data, summary, format, filters)

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return nil, pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func (r *ReportServiceImpl) GenerateUsersReport(ctx context.Context, format string, filters services.ReportFilters) (error) {
	var report *userReport

	if filters.UserId != nil {
		data, err := r.store.Users.GetReportUserUsersData(ctx, *filters.UserId, filters)
		if err != nil {
			return err
		}

		log.Println(data)

		report = newUserReport([]services.UserAdminsReportData{}, data, services.UserAdminsSummary{}, format, filters)
		format = "pdf"
	} else {
		data, summary, err := r.store.Users.GetReportUserAdminData(ctx, filters)
		if err != nil {
			return err
		}

		report = newUserReport(data, services.UserUsersReportData{}, summary, format, filters)
	}

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func (r *ReportServiceImpl) GenerateClientsReport(ctx context.Context, format string, filters services.ReportFilters) (error) {
	var report *clientReport

	if filters.ClientId != nil {
		data, err := r.store.Clients.GetReportClientClientsData(ctx, *filters.ClientId, filters)
		if err != nil {
			return err
		}

		log.Println(data)

		report = newClientReport([]services.ClientAdminsReportData{}, data, services.ClientSummary{}, format, filters)
		format = "pdf"
	} else {
		data, summary, err := r.store.Clients.GetReportClientAdminData(ctx, filters)
		if err != nil {
			return err
		}

		report = newClientReport(data, services.ClientClientsReportData{}, summary, format, filters)
	}

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func (r *ReportServiceImpl) GenerateLoansReport(ctx context.Context, format string, filters services.ReportFilters) (error) {
	var report *loanReport

	if filters.LoanId != nil {
		data, err := r.store.Loans.GetReportLoanByIdData(ctx, *filters.LoanId)
		if err != nil {
			return err
		}

		log.Println(data)

		report = newLoanReport([]services.LoanReportData{}, data, services.LoanSummary{}, format, filters)
		format = "pdf"
	} else {
		data, summary, err := r.store.Loans.GetReportLoanData(ctx, filters)
		if err != nil {
			return err
		}

		report = newLoanReport(data, services.LoanReportDataById{}, summary, format, filters)
	}

	switch format {
	case "excel":
		return report.generateExcel("Sheet1")
	case "pdf":
		return report.generatePDF()
	default:
		return pkg.Errorf(pkg.INVALID_ERROR, "unsupported format")
	}
}

func formatMoney(amount float64) string {
	p := message.NewPrinter(message.MatchLanguage("en"))
	return p.Sprintf("%.2f", amount)
}

func formatQuantity(quantity int64) string {
	return fmt.Sprintf("%d", quantity)
}

func formatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "N/A" 
	}
	return t.Format("2006-01-02")
}