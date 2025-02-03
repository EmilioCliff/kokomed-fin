package workers

import (
	"context"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"github.com/hibiken/asynq"
)

var _ services.WorkerService = (*WorkerServiceImpl)(nil)

type WorkerServiceImpl struct {
	distributor *TaskDistributor
	processor   *TaskProcessor
}

func NewWorkerService(redisConfig services.RedisConfig, emailSender pkg.EmailSender, repo *mysql.MySQLRepo, maker pkg.JWTMaker) services.WorkerService {
	redisOpt := asynq.RedisClientOpt{
		Addr: redisConfig.Address,
		DB: redisConfig.DB,
		Password: redisConfig.Password,
	}

	return &WorkerServiceImpl{
		distributor: NewTaskDistributor(redisOpt),
		processor:   NewTaskProcessor(redisOpt, emailSender, repo, maker),
	}
}

func (w *WorkerServiceImpl) StartProcessor() error {
	return w.processor.Start()
}

func (w *WorkerServiceImpl) StopProcessor() {
	w.processor.Stop()
}

func (w *WorkerServiceImpl) ProcessSendResetPassword(ctx context.Context, task *asynq.Task) error {
	return w.processor.ProcessSendResetPassword(ctx, task)
}

func (w *WorkerServiceImpl) DistributeTaskSendResetPassword(ctx context.Context, payload services.SendResetPasswordPayload, opt ...asynq.Option, ) error {
	return w.distributor.DistributeTaskSendResetPassword(ctx, payload, opt...)
}