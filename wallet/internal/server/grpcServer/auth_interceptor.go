package grpcserver

import (
	"context"

	"github.com/workshops/wallet/internal/middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtWrapper *auth.JwtWrapper
}

func NewAuthInterceptor(wrapper *auth.JwtWrapper) *AuthInterceptor {
	return &AuthInterceptor{jwtWrapper: wrapper}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		err := interceptor.authorize(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	accessToken := values[0]
	_, err := interceptor.jwtWrapper.ValidateToken(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}
	return nil
}
