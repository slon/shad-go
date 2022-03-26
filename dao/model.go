package dao

import (
	"context"
)

type (
	UserID int

	User struct {
		ID   UserID
		Name string
	}

	Dao interface {
		Create(ctx context.Context, u *User) (UserID, error)
		Update(ctx context.Context, u *User) error
		Delete(ctx context.Context, id UserID) error
		Lookup(ctx context.Context, id UserID) (User, error)
		List(ctx context.Context) ([]User, error)
	}
)
