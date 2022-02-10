//go:build !solution
// +build !solution

package dupcall

import "context"

type Call struct {
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error)
