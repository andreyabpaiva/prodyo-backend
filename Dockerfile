# syntax=docker/dockerfile:1
FROM golang:1.25-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api ./cmd/api/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api ./api
COPY --from=builder /app/cmd/migrations ./cmd/migrations

ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=prodyo_user
ENV DB_PASSWORD=prodyo_password
ENV DB_NAME=prodyo_db

CMD ["./api"]
