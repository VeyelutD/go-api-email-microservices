package main

import (
	"encoding/json"
	"fmt"
	"github.com/VeyelutD/go-email-microservice/rabbitmq"
	"gopkg.in/gomail.v2"
	"os"
)

type LoginOTPPayload struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type ConfirmPayload struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

var emailFrom = os.Getenv("EMAIL_FROM")

func sendLoginOTP(d *gomail.Dialer, messageBody []byte) error {
	m := gomail.NewMessage()
	var messagePayload LoginOTPPayload
	if err := json.Unmarshal(messageBody, &messagePayload); err != nil {
		return err
	}
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", messagePayload.Email)
	m.SetHeader("Subject", "One Time Password")
	m.SetBody("text/plain", fmt.Sprintf("Hello, to login enter the code below:\n%s", messagePayload.Code))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
func sendConfirmationLink(d *gomail.Dialer, messageBody []byte) error {
	m := gomail.NewMessage()
	var messagePayload ConfirmPayload
	if err := json.Unmarshal(messageBody, &messagePayload); err != nil {
		return err
	}
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", messagePayload.Email)
	m.SetHeader("Subject", "Confirmation link")
	m.SetBody("text/plain", fmt.Sprintf("Hello, use this link to confirm your account:\nhttp://localhost:8000/v1/auth/confirm?token=%s", messagePayload.Token))
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

var (
	rabbitmqUsername = os.Getenv("RABBITMQ_USERNAME")
	rabbitmqPassword = os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqHost     = os.Getenv("RABBITMQ_HOST")
	rabbitmqVHost    = os.Getenv("RABBITMQ_VHOST")
)

func main() {

	conn, err := rabbitmq.ConnectRabbitMQ(rabbitmqUsername, rabbitmqPassword, rabbitmqHost, rabbitmqVHost)
	if err != nil {
		panic(err)
	}
	var blocking chan struct{}
	defer conn.Close()
	emailClient, err := rabbitmq.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer emailClient.Close()
	queue, err := emailClient.CreateQueue("", true, false)
	if err != nil {
		panic(err)
	}
	if err := emailClient.CreateBinding(queue.Name, "email.*", "email_exchange"); err != nil {
		panic(err)
	}
	d := gomail.NewDialer("smtp.gmail.com", 587, emailFrom, os.Getenv("EMAIL_PASSWORD"))
	messageBus, err := emailClient.Consume(queue.Name, "email", false)
	if err != nil {
		panic(err)
	}
	for message := range messageBus {
		if message.RoutingKey == "email.login" {
			err = sendLoginOTP(d, message.Body)
			if err != nil {
				panic(err)
			}
		} else if message.RoutingKey == "email.registration" {
			err = sendConfirmationLink(d, message.Body)
			if err != nil {
				panic(err)
			}
		} else {
			panic("invalid routing key")
		}
		if err = message.Ack(false); err != nil {
			panic(err)
		}
	}
	<-blocking
}
