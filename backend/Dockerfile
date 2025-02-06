FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
COPY /app/.envs/.local/config.env .
RUN go build -o main /app/cmd/server/main.go

EXPOSE 3030

CMD ["./main"]
