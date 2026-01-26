package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Quickaxe-Martina/link_shortening_service/internal/storage"
	"google.golang.org/grpc/metadata"
)

var ErrNoAuthHeader = errors.New("no authorization metadata")

func GetOrCreateUserFromGRPCContext(
	ctx context.Context,
	store storage.Storage,
	secretKey string,
	tokenExp time.Duration,
) (storage.User, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return store.CreateUser(ctx)
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return store.CreateUser(ctx)
	}

	token := extractBearerToken(authHeaders[0])
	if token == "" {
		return store.CreateUser(ctx)
	}

	userID := GetUserID(token, secretKey)
	if userID == -1 {
		return store.CreateUser(ctx)
	}

	return storage.User{ID: userID}, nil
}

func extractBearerToken(h string) string {
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return h
}
