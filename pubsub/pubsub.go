//go:build !change

package pubsub

import "context"

// MsgHandler is a callback function that processes messages delivered to subscribers.
type MsgHandler func(msg any)

type Subscription interface {
	// Unsubscribe will remove interest in the given subject.
	Unsubscribe()
}

type PubSub interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subj string, cb MsgHandler) (Subscription, error)

	// Publish publishes the msg argument to the given subject.
	Publish(subj string, msg any) error

	// Close will shutdown pub-sub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}
