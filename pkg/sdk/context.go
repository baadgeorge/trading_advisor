package sdk

import (
	"context"
	"final/internal/configs"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// финкция возвращает контекст для запросов к API
func createRequestContext(token string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.DefaultRequestTimeout)

	authHeader := fmt.Sprintf("Bearer %s", token)
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeader)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-tracking-id", uuid.New().String())
	ctx = metadata.AppendToOutgoingContext(ctx, "x-app-name", configs.AppName)

	return ctx, cancel
}
