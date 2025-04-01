package main

import (
	"context"
	"fmt"
	api "github.com/VeyelutD/go-api-microservice/internal"
	"github.com/VeyelutD/go-api-microservice/internal/auth"
	"github.com/VeyelutD/go-api-microservice/internal/db"
	"github.com/VeyelutD/go-api-microservice/internal/email"
	"github.com/VeyelutD/go-api-microservice/internal/rabbitmq"
	"github.com/VeyelutD/go-api-microservice/internal/users"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"os"
)

func RegisterV1Routes(queries *db.Queries, db *pgx.Conn, rabbitmqClient *rabbitmq.RabbitClient) http.Handler {
	router := http.NewServeMux()
	emailService := email.NewService(rabbitmqClient)
	userService := users.NewService(queries)
	authService := auth.NewService(queries, db, emailService, userService)
	authHandler := auth.NewHandler(authService)
	router.HandleFunc("POST /v1/auth/register", authHandler.Register)
	router.HandleFunc("POST /v1/auth/send-otp", authHandler.SendOTP)
	router.HandleFunc("POST /v1/auth/verify-otp", authHandler.VerifyOTP)
	router.HandleFunc("GET /v1/auth/confirm", authHandler.ConfirmAccount)
	return router
}

var (
	rabbitmqUsername = os.Getenv("RABBITMQ_USERNAME")
	rabbitmqPassword = os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqHost     = os.Getenv("RABBITMQ_HOST")
	rabbitmqVHost    = os.Getenv("RABBITMQ_VHOST")
	postgresHost     = os.Getenv("POSTGRES_HOST")
	postgresPort     = os.Getenv("POSTGRES_PORT")
	postgresUser     = os.Getenv("POSTGRES_USER")
	postgresPassword = os.Getenv("POSTGRES_PASSWORD")
)

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", postgresUser, postgresPassword, postgresHost, postgresPort)
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("An error occurred connecting to postgres, exiting", "error", err)
		os.Exit(1)
	}
	defer conn.Close()
	rabbitConn, err := rabbitmq.ConnectRabbitMQ(rabbitmqUsername, rabbitmqPassword, rabbitmqHost, rabbitmqVHost)
	if err != nil {
		slog.Error("An error occurred connecting to rabbitmq, exiting", "error", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	rabbitClient, err := rabbitmq.NewRabbitMQClient(rabbitConn)
	if err != nil {
		slog.Error("An error occurred creating a new rabbitmq client, exiting", "error", err)
		os.Exit(1)
	}
	defer rabbitClient.Close()
	queries := db.New(conn)
	apiServer := api.NewServer(":8000")
	dbp, err := conn.Acquire(context.Background())
	if err != nil {
		slog.Error("An error occurred getting a db connection, exiting", "error", err)
		os.Exit(1)
	}
	defer dbp.Release()
	router := RegisterV1Routes(queries, dbp.Conn(), rabbitClient)
	err = apiServer.Run(router)
	if err != nil {
		slog.Error("An error occurred running the server, exiting", "error", err)
		os.Exit(1)
	}
}
