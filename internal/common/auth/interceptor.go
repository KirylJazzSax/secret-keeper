package auth

import (
	"context"
	"fmt"

	"github.com/KirylJazzSax/secret-keeper/internal/common/errors"
	"github.com/KirylJazzSax/secret-keeper/internal/common/token"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/samber/do"
	"google.golang.org/grpc"
)

func AuthFunc(ctx context.Context) (context.Context, error) {
	// test method
	method, _ := grpc.Method(ctx)
	fmt.Println(method)
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

	grpc_ctxtags.Extract(ctx).Set("user", payload)

	newCtx := context.WithValue(ctx, "user", payload)

	return newCtx, nil
}
