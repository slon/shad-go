//go:build !solution

package dupcall

import "context"

type Call struct {
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (any, error),
) (result any, err error)
