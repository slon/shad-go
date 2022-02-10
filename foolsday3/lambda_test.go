//go:build !race
// +build !race

package foolsday3

import (
	"context"
	"testing"
	"time"
)

func TestLambda(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var doNotPrint bool
	var validateLambdaFunc = func(end time.Time) bool {
		return time.Now() == end
	}
	result := lambda(ctx)
	end := time.Now()
	if validateLambdaFunc(end) == true {
		t.Logf("[%s] Great! Your function is very fast!", end.Format("15:04:05.999999"))
		if doNotPrint != true {
			t.Log("Congrats!")
			return
		}
		t.Log(result)
		t.FailNow()
	}
	t.Logf("[%s] result of your slow function:", end.Format("15:04:05.999999"))
	t.Log(result)
	t.FailNow()
}
