package auth

import (
	"context"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/samber/do"
)

func AuthFunc(ctx context.Context) (context.Context, error) {
	t, err := do.Invoke[token.Maker](nil)
	if err != nil {
		return nil, err
	}

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	payload, err := t.VerifyToken(token)
	if err != nil {
		return nil, errors.UnAuthErr()
	}

	tags.Extract(ctx).Set("user", payload)

	newCtx := context.WithValue(ctx, "user", payload)

	return newCtx, nil
}
