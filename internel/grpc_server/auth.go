package grpcserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/ak-karimzai/bank-api/internel/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *GRPCServer) authUser(
	ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	tokenFileds := strings.Split(authHeader, " ")
	if len(tokenFileds) < 2 {
		return nil, fmt.Errorf("invalid authorizaion format")
	}

	authType := strings.ToLower(tokenFileds[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := tokenFileds[1]
	payload, err := server.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil
}
