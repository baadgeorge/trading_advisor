package sdk

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"someshit/internal/configs"
)

//const token = "t.LCgR7EtlcOlfRM1epjYTV0Z9bF2gqeGZhUH831L8vUBb0-LH6EuIXEW5o6k4XKmpA7vPVw39hU02de1Mhcv-yw"

// createRequestContext returns context for API calls with timeout and auth headers attached.
func createRequestContext(token string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.DefaultRequestTimeout)

	authHeader := fmt.Sprintf("Bearer %s", token)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-app-name", configs.AppName)

	return ctx, cancel
}

// createRequestContext returns context for streams with auth headers attached.
func createStreamContext(token string) context.Context {
	ctx := context.TODO()

	authHeader := fmt.Sprintf("Bearer %s", token)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-app-name", configs.AppName)

	return ctx
}
