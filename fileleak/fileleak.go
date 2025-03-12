//go:build !solution

package fileleak

type testingT interface {
	Errorf(msg string, args ...any)
	Cleanup(func())
}

func VerifyNone(t testingT) {
	panic("implement me")
}
