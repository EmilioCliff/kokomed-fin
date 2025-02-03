package workers

import "github.com/hibiken/asynq"

type TaskDistributor struct {
	client *asynq.Client
}

func NewTaskDistributor(redisOpt asynq.RedisClientOpt) *TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &TaskDistributor{
		client: client,
	}
}