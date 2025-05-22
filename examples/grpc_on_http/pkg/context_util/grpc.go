package context_util

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func NewGrpcAuthContext(ctx context.Context, apiKey string) context.Context {
	// Add metadata
	md := metadata.New(map[string]string{
		"api-key": apiKey,
	})

	// Create context with metadata
	return metadata.NewOutgoingContext(ctx, md)
}
