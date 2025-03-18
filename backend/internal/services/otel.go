package services

import (
	"context"

	"google.golang.org/grpc"
)

type OtelService interface {
	InitializeTracerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error)
	InitializeMeterProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error)
	InitializeLoggerProvider(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error)
}