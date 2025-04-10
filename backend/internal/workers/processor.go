package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server *asynq.Server
	sender pkg.EmailSender
	repo   *mysql.MySQLRepo
	maker  pkg.JWTMaker
}

func NewTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	sender pkg.EmailSender,
	repo *mysql.MySQLRepo,
	maker pkg.JWTMaker,
) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			services.QueueCritical: 10,
			services.QueueDefault:  5,
			services.QueueLow:      2,
		},
		RetryDelayFunc: CustomRetryDelayFunc,
		ErrorHandler:   asynq.ErrorHandlerFunc(ReportError),
		LogLevel:       asynq.WarnLevel,
	})

	return &TaskProcessor{
		server: server,
		sender: sender,
		repo:   repo,
		maker:  maker,
	}
}

func (processor *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendResetPasswordTask, processor.ProcessSendResetPassword)

	return processor.server.Start(mux)
}

func (processor *TaskProcessor) Stop() {
	processor.server.Shutdown()
	log.Println("Task processor stopped successfully.")
}

func CustomRetryDelayFunc(_ int, _ error, _ *asynq.Task) time.Duration {
	return 2 * time.Second
}

func ReportError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		err = fmt.Errorf("retry exhausted for task %s: %w", task.Type(), err)
	}
	log.Println(err)
	// log it or something
	// errorReportingService.Notify(err)
}
