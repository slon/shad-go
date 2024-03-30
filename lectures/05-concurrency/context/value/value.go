package main

import (
	"context"
	"fmt"
)

type myKey struct{} // use private type to restrict access to this package

func WithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, myKey{}, user)
}

// Export type-safe interface for users of this value
func ContextUser(ctx context.Context) (string, bool) {
	v := ctx.Value(myKey{})
	s, ok := v.(string)
	return s, ok
}

// OMIT

func main() {
	ctx := context.Background()

	user, ok := ContextUser(ctx)
	fmt.Println(ok, user)

	ctx = WithUser(ctx, "petya")
	user, ok = ContextUser(ctx)
	fmt.Println(ok, user)
}
