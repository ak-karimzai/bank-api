package worker

import (
	"context"
	"encoding/json"
	"fmt"

	errorhandler "github.com/ak-karimzai/bank-api/internel/error_handler"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (td *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	barr, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %v", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, barr, opts...)
	taskInfo, err := td.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue the task: %v", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVeridyEmail(
	ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		finalErr := errorhandler.DbErrorHandler(err)
		if finalErr.Status == errorhandler.NotFound {
			return fmt.Errorf("user doesn't exist: %v", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %v", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")
	return nil
}
