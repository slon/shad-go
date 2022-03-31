package shopfront

import "context"

type (
	ItemID int64
	UserID int64
)

type Item struct {
	ViewCount int
	Viewed    bool
}

type Counters interface {
	GetItems(ctx context.Context, ids []ItemID, userID UserID) ([]Item, error)

	RecordView(ctx context.Context, id ItemID, userID UserID) error
}
