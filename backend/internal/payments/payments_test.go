package payments_test

import (
	"context"
	"testing"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func TestProcessCallback(t *testing.T) {
	db := mysql.NewMySQLRepo(&pkg.Config{})
	store := mysql.NewStore(&pkg.Config{})
	payments := payments.NewPaymentService(db, store)

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
