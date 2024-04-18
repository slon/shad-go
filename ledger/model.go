package ledger

import (
	"context"
	"errors"
)

type (
	ID    string
	Money int64
)

var ErrNoMoney = errors.New("no money")
var ErrNegativeAmount = errors.New("negative amount")

type Ledger interface {
	CreateAccount(ctx context.Context, id ID) error
	GetBalance(ctx context.Context, id ID) (Money, error)
	Deposit(ctx context.Context, id ID, amount Money) error
	Withdraw(ctx context.Context, id ID, amount Money) error
	Transfer(ctx context.Context, from, to ID, amount Money) error
	Close() error
}
