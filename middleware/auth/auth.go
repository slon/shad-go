//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
)

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	panic("not implemented")
}

var ErrInvalidToken = errors.New("invalid token")

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	panic("not implemented")
}
