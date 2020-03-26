package pubsub

import (
	"context"
	"fmt"
	"sync"
)

func ExamplePubSub() {
	p := NewPubSub()
	defer func() { _ = p.Close(context.Background()) }()

	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err := p.Subscribe("single", func(msg interface{}) {
		fmt.Println("new message:", msg)
		// Output: new message: blah-blah
		wg.Done()
	})
	if err != nil {
		panic(err)
	}

	err = p.Publish("single", "blah-blah")
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
