# Build stage
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .

# install migrate
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

RUN go build -o main /app/cmd/server/main.go

# Run stage
FROM alpine:3.19
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY --from=builder /app/internal/mysql/migrations /app/migrations
COPY --from=builder /app/afya_credit.png .

EXPOSE 3030

CMD ["./main"]
