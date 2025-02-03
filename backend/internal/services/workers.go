package services

import (
	"context"

	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type RedisConfig struct {
	Address string 
	Password  string
	DB int
}

type WorkerService interface {
	// processes
	StartProcessor() error
	StopProcessor()


	ProcessSendResetPassword(ctx context.Context, task *asynq.Task) error

	// distributor
	DistributeTaskSendResetPassword(ctx context.Context, payload SendResetPasswordPayload, opt ...asynq.Option, ) error
}

type SendResetPasswordPayload struct {
	Email string `json:"email"`
}