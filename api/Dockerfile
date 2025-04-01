FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ./go.mod ./go.sum ./

RUN go mod download


COPY . .

COPY sql/migrations /sql/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./main
COPY --from=builder /go/bin/goose /go/bin/goose
COPY --from=builder /sql/migrations /sql/migrations

EXPOSE 8000

CMD ["/bin/sh", "-c", "/go/bin/goose -dir /sql/migrations postgres \"user=$POSTGRES_USER dbname=postgres host=db password=$POSTGRES_PASSWORD sslmode=disable\" up && ./main"]