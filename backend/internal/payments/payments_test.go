package payments

import (
	"context"
	"database/sql"
	"testing"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mock"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/generated"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql/mockdb"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/repository"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"go.uber.org/mock/gomock"
)

type TestGRPCServer struct {
	*PaymentService
	LoansRepository mock.MockLoanRepository
}

func mockCreateLoan(ctx context.Context, loan *repository.Loan) (repository.LoanFullData, error) {
	return repository.LoanFullData{}, nil
}

func TestProcessCallback(t *testing.T) {
	server := &TestGRPCServer{}

	server.mySQL.Loans = &server.LoansRepository

	server.LoansRepository.MockCreateLoan = mockCreateLoan

	store := mysql.NewStore(pkg.Config{})
	db := mysql.NewMySQLRepo(store)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// store.

	mockQueries := mockdb.NewMockQuerier(ctrl)

	store.NewQuerierFn = func(tx *sql.Tx) generated.Querier {
		return mockQueries
	}

	payments := NewPaymentService(db, store)

	callbackData := &services.MpesaCallbackData{
		TransactionSource: "MPESA",
		TransactionID:     "1234567890",
		AccountNumber:     "1234567890",
		PhoneNumber:       "1234567890",
		PayingName:        "John Doe",
		Amount:            1000,
		AssignedTo:        &[]uint32{1}[0],
	}

	loanId, err := payments.ProcessCallback(context.Background(), callbackData)
	if err != nil {
		t.Fatalf("failed to process callback: %v", err)
	}

	t.Logf("loanId: %d", loanId)
}
