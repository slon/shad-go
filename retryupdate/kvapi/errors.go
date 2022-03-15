//go:build !change

package kvapi

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

var (
	_ error = (*APIError)(nil)
	_ error = (*ConflictError)(nil)
	_ error = (*AuthError)(nil)
)

var ErrKeyNotFound = errors.New("key not found")

type (
	APIError struct {
		Method string

		Err error
	}

	ConflictError struct {
		ProvidedVersion, ExpectedVersion uuid.UUID
	}

	AuthError struct {
		Msg string
	}
)

func (a *APIError) Error() string {
	return fmt.Sprintf("api: %q error: %v", a.Method, a.Err)
}

func (a *APIError) Unwrap() error {
	return a.Err
}

func (a *ConflictError) Error() string {
	return fmt.Sprintf("api: conflict: expected_version=%d, provided_version=%d", a.ExpectedVersion, a.ProvidedVersion)
}

func (a *AuthError) Error() string {
	return fmt.Sprintf("api: auth: %s", a.Msg)
}
