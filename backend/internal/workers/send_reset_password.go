package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/hibiken/asynq"
)

const SendResetPasswordTask = "task:send_verify_email"
const resetPasswordExpiriy = 10*time.Minute

func (distributor TaskDistributor) DistributeTaskSendResetPassword(ctx context.Context, payload services.SendResetPasswordPayload, opt ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to marshal payload", err)
	}

	task := asynq.NewTask(SendResetPasswordTask, jsonPayload, opt...)
	_, err = distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return pkg.Errorf(pkg.INTERNAL_ERROR, "failed to enqueue task", err)
	}

	return nil
}

func (processor *TaskProcessor) ProcessSendResetPassword(ctx context.Context, task *asynq.Task) error {
	var payLoad services.SendResetPasswordPayload
	if err := json.Unmarshal(task.Payload(), &payLoad); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	user, err := processor.repo.Users.GetUserByEmail(ctx, payLoad.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no user found: %w", err)
		}
		return fmt.Errorf("internal error: %w", err)
	}

	accessToken, err := processor.maker.CreateToken(user.Email, user.ID, user.BranchID, user.Role, resetPasswordExpiriy)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}

	resetPasswordLink := fmt.Sprintf("http://127.0.0.1:5173/reset-password/%v", accessToken)
	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>We received a request to reset your password. Click the link below to reset it:</p>
	<a href="%s" style="display:inline-block; padding:10px 20px; background-color:#007BFF; color:#fff; text-decoration:none; border-radius:5px;">Reset Password</a>
	<h5>The link is valid for 10 Minutes</h5>
`, user.FullName, resetPasswordLink)

	err = processor.sender.SendMail("Reset Password", emailBody, "application/pdf", []string{"emiliocliff@gmail.com"}, nil, nil, nil, nil) // test email with my email
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}