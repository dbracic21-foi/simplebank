package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/dbracic21-foi/simplebank/db/sqlc"
	"github.com/dbracic21-foi/simplebank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("Task enqueued")
	return nil

}
func (processor *RedisTaskProcessor) ProcessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}
	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("cannot get user: %w", err)

	}
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RadnomString(32),
	})
	if err != nil {
		return fmt.Errorf("cannot create verify email: %w", err)

	}
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_emails?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	subject := "Confirm your email adress for Simple Bank!"
	content := fmt.Sprintf(`<<html>
	<head>
	  <style>
		body { font-family: Arial, sans-serif; line-height: 1.6; }
		.content { margin: 20px; padding: 20px; border-radius: 10px; background-color: #f4f4f4; }
		.button { background-color: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
	  </style>
	</head>
	<body>
	  <div class="content">
		<p>Thank you for registering with us. Please click the link below to verify your email address and activate your account:</p>
		<a href="%s" class="button">Verify Email</a>
		<p>If you are having trouble clicking the button, copy and paste the URL below into your web browser:</p>
		<p>%s</p>
		<p>Best regards,<br>Simple Bank Team</p>
	  </div>
	</body>
	</html>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}
	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("cannot send email: %w", err)

	}

	log.Info().Str("username", user.Username).Str("email", user.Email).Msg("Email sent")
	return nil
}
