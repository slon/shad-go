// +build !solution

package pubsub

import "context"

var _ Subscription = (*MySubscription)(nil)

type MySubscription struct{}

func (s *MySubscription) Unsubscribe() {
	panic("implement me")
}

var _ PubSub = (*MyPubSub)(nil)

type MyPubSub struct{}

func NewPubSub() PubSub {
	panic("implement me")
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	panic("implement me")
}

func (p *MyPubSub) Publish(subj string, msg interface{}) error {
	panic("implement me")
}

func (p *MyPubSub) Close(ctx context.Context) error {
	panic("implement me")
}
