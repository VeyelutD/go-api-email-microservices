package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VeyelutD/go-api-microservice/internal/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

type Service struct {
	rabbitClient rabbitmq.RabbitClient
}
type LoginPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type ConfirmationPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewService(rabbitClient *rabbitmq.RabbitClient) *Service {
	return &Service{
		rabbitClient: *rabbitClient,
	}
}

func (es *Service) SendOTP(ctx context.Context, email, code string) error {
	payload, err := json.Marshal(LoginPayload{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}
	if err := es.rabbitClient.Send(ctx, "email_exchange", "email.login", amqp091.Publishing{
		ContentType:   "application/json",
		CorrelationId: email,
		Body:          payload,
	}); err != nil {
		return ErrCouldNotSendOTP
	}
	return nil
}

func (es *Service) SendConfirmationLink(ctx context.Context, email, token string) error {
	payload, err := json.Marshal(ConfirmationPayload{
		Email: email,
		Token: token,
	})
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}
	if err := es.rabbitClient.Send(ctx, "email_exchange", "email.registration", amqp091.Publishing{
		ContentType:   "application/json",
		CorrelationId: email,
		Body:          payload,
	}); err != nil {
		return ErrCouldNotSendConfirmation
	}
	return nil
}
