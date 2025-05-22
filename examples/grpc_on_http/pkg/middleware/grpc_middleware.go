package grpc_middleware

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GrpcMiddleware struct {
	apiKey string
}

func New(apiKey string) *GrpcMiddleware {
	return &GrpcMiddleware{
		apiKey: apiKey,
	}
}

// UnaryInterceptor for gRPC server to authenticate API key
func (m *GrpcMiddleware) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if keys := md["api-key"]; len(keys) > 0 && keys[0] == m.apiKey {
			// Process request if API key is valid
			return handler(ctx, req)
		}
	}

	return nil, status.Errorf(status.Code(errors.New("PermissionDenied")), "Permission Denied")
}
